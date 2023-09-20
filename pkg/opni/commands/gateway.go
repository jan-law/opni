//go:build !minimal && !cli

package commands

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/rancher/opni/pkg/config"
	"github.com/rancher/opni/pkg/config/v1beta1"
	"github.com/rancher/opni/pkg/dashboard"
	"github.com/rancher/opni/pkg/features"
	"github.com/rancher/opni/pkg/gateway"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/machinery"
	"github.com/rancher/opni/pkg/management"
	cliutil "github.com/rancher/opni/pkg/opni/util"
	"github.com/rancher/opni/pkg/plugins"
	"github.com/rancher/opni/pkg/plugins/hooks"
	"github.com/rancher/opni/pkg/tracing"
	"github.com/rancher/opni/pkg/update/noop"
	"github.com/rancher/opni/pkg/util/waitctx"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"k8s.io/client-go/rest"

	_ "github.com/rancher/opni/pkg/oci/kubernetes"
	_ "github.com/rancher/opni/pkg/oci/noop"
	_ "github.com/rancher/opni/pkg/plugins/apis"
	_ "github.com/rancher/opni/pkg/storage/crds"
	_ "github.com/rancher/opni/pkg/storage/etcd"
	_ "github.com/rancher/opni/pkg/storage/jetstream"
)

func BuildGatewayCmd() *cobra.Command {
	lg := logger.New()
	var configLocation string

	run := func() error {
		tracing.Configure("gateway")

		objects := cliutil.LoadConfigObjectsOrDie(configLocation, lg)

		ctx, cancel := context.WithCancel(waitctx.Background())

		inCluster := true
		restconfig, err := rest.InClusterConfig()
		if err != nil {
			if errors.Is(err, rest.ErrNotInCluster) {
				inCluster = false
			} else {
				lg.Error("failed to create config", "fatal", err)
				os.Exit(1)
			}
		}

		var fCancel context.CancelFunc
		if inCluster {
			features.PopulateFeatures(ctx, restconfig)
			fCancel = features.FeatureList.WatchConfigMap()
		} else {
			fCancel = cancel
		}

		machinery.LoadAuthProviders(ctx, objects)
		var gatewayConfig *v1beta1.GatewayConfig
		found := objects.Visit(
			func(config *v1beta1.GatewayConfig) {
				if gatewayConfig == nil {
					gatewayConfig = config
				}
			},
			func(ap *v1beta1.AuthProvider) {
				// noauth is a special case, we need to start a noauth server but only
				// once - other auth provider consumers such as plugins can load
				// auth providers themselves, but we don't want them to start their
				// own noauth server.
				if ap.Name == "noauth" {
					server := machinery.NewNoauthServer(ctx, ap)
					waitctx.Go(ctx, func() {
						if err := server.ListenAndServe(ctx); err != nil {
							lg.Warn("noauth server exited with error", logger.Err(err))
						}
					})
				}
			},
		)
		if !found {
			lg.Error("config file does not contain a GatewayConfig object", "config", configLocation)
			os.Exit(1)
		}

		lg.Info("loading plugins", "dir", gatewayConfig.Spec.Plugins.Dir)
		pluginLoader := plugins.NewPluginLoader(plugins.WithLogger(lg.WithGroup("gateway")))

		lifecycler := config.NewLifecycler(objects)
		g := gateway.NewGateway(ctx, gatewayConfig, pluginLoader,
			gateway.WithLifecycler(lifecycler),
			gateway.WithExtraUpdateHandlers(noop.NewSyncServer()),
		)

		m := management.NewServer(ctx, &gatewayConfig.Spec.Management, g, pluginLoader,
			management.WithCapabilitiesDataSource(g),
			management.WithHealthStatusDataSource(g),
			management.WithLifecycler(lifecycler),
		)

		g.MustRegisterCollector(m)

		// start web server
		d, err := dashboard.NewServer(&gatewayConfig.Spec.Management)
		if err != nil {
			lg.Error("failed to initialize web dashboard", logger.Err(err))
		} else {
			pluginLoader.Hook(hooks.OnLoadingCompleted(func(int) {
				waitctx.AddOne(ctx)
				defer waitctx.Done(ctx)
				if err := d.ListenAndServe(ctx); err != nil {
					lg.Warn("dashboard server exited with error", logger.Err(err))
				}
			}))
		}

		pluginLoader.Hook(hooks.OnLoadingCompleted(func(numLoaded int) {
			lg.Info("loaded plugins", "count", numLoaded)
		}))

		pluginLoader.Hook(hooks.OnLoadingCompleted(func(int) {
			waitctx.AddOne(ctx)
			defer waitctx.Done(ctx)
			if err := m.ListenAndServe(ctx); err != nil {
				lg.Warn("management server exited with error", logger.Err(err))
			}
		}))

		pluginLoader.Hook(hooks.OnLoadingCompleted(func(int) {
			waitctx.AddOne(ctx)
			defer waitctx.Done(ctx)
			if err := g.ListenAndServe(ctx); err != nil {
				lg.Warn("gateway server exited with error", logger.Err(err))
			}
		}))

		pluginLoader.LoadPlugins(ctx, gatewayConfig.Spec.Plugins.Dir, plugins.GatewayScheme)

		style := chalk.Yellow.NewStyle().
			WithBackground(chalk.ResetColor).
			WithTextStyle(chalk.Bold)
		reloadC := make(chan struct{})
		go func() {
			c, err := lifecycler.ReloadC()

			fNotify := make(<-chan struct{})
			if inCluster {
				fNotify = features.FeatureList.NotifyChange()
			}

			if err != nil {
				lg.Error("failed to get reload channel from lifecycler", logger.Err(err))
				os.Exit(1)
			}
			select {
			case <-c:
				lg.Info(style.Style("--- received reload signal ---"))
				close(reloadC)
			case <-fNotify:
				lg.Info(style.Style("--- received feature update signal ---"))
				close(reloadC)
			}
		}()

		<-reloadC
		lg.Info(style.Style("waiting for servers to shut down"))
		fCancel()
		cancel()
		waitctx.WaitWithTimeout(ctx, 60*time.Second, 10*time.Second)

		atomic.StoreUint32(&plugin.Killed, 0)
		lg.Info(style.Style("--- reloading ---"))
		return nil
	}

	serveCmd := &cobra.Command{
		Use:   "gateway",
		Short: "Run the Opni Monitoring Gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer waitctx.RecoverTimeout()
			for {
				if err := run(); err != nil {
					return err
				}
			}
		},
	}

	serveCmd.Flags().StringVar(&configLocation, "config", "", "Absolute path to a config file")
	return serveCmd
}

func init() {
	AddCommandsToGroup(OpniComponents, BuildGatewayCmd())
}
