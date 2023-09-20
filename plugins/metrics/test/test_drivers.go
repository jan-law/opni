package test

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"

	corev1 "github.com/rancher/opni/pkg/apis/core/v1"
	"github.com/rancher/opni/pkg/config/v1beta1"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/plugins/driverutil"
	"github.com/rancher/opni/pkg/rules"
	"github.com/rancher/opni/pkg/test"
	"github.com/rancher/opni/pkg/util/notifier"
	"github.com/rancher/opni/pkg/util/waitctx"
	"github.com/rancher/opni/plugins/metrics/apis/cortexops"
	"github.com/rancher/opni/plugins/metrics/apis/node"
	"github.com/rancher/opni/plugins/metrics/apis/remoteread"
	metrics_agent_drivers "github.com/rancher/opni/plugins/metrics/pkg/agent/drivers"
	metrics_drivers "github.com/rancher/opni/plugins/metrics/pkg/gateway/drivers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	cortexVersion string
)

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		panic("could not read build info")
	}
	// https://github.com/golang/go/issues/33976
	if buildInfo.Main.Path == "" {
		cortexVersion = "(unknown)"
	} else {
		var found bool
		for _, depInfo := range buildInfo.Deps {
			if depInfo.Path == "github.com/cortexproject/cortex" {
				if depInfo.Replace != nil {
					cortexVersion = depInfo.Replace.Version
				} else {
					cortexVersion = depInfo.Version
				}
				found = true
				break
			}
		}
		if !found {
			panic("could not find cortex dependency in build info")
		}
	}

	metrics_drivers.ClusterDrivers.Register("test-environment", func(ctx context.Context, _ ...driverutil.Option) (metrics_drivers.ClusterDriver, error) {
		env := test.EnvFromContext(ctx)
		return NewTestEnvMetricsClusterDriver(env), nil
	})

	metrics_agent_drivers.NodeDrivers.Register("test-environment-prometheus", func(ctx context.Context, _ ...driverutil.Option) (metrics_agent_drivers.MetricsNodeDriver, error) {
		env := test.EnvFromContext(ctx)
		return NewTestEnvPrometheusNodeDriver(env), nil
	})

	metrics_agent_drivers.NodeDrivers.Register("test-environment-otel", func(ctx context.Context, _ ...driverutil.Option) (metrics_agent_drivers.MetricsNodeDriver, error) {
		env := test.EnvFromContext(ctx)
		return NewTestEnvOtelNodeDriver(env), nil
	})
}

type TestEnvMetricsClusterDriver struct {
	cortexops.UnsafeCortexOpsServer

	lock         sync.Mutex
	state        cortexops.InstallState
	cortexCtx    context.Context
	cortexCancel context.CancelFunc

	Env           *test.Environment
	Configuration *cortexops.ClusterConfiguration
}

func NewTestEnvMetricsClusterDriver(env *test.Environment) *TestEnvMetricsClusterDriver {
	return &TestEnvMetricsClusterDriver{
		Env:           env,
		Configuration: &cortexops.ClusterConfiguration{},
		state:         cortexops.InstallState_NotInstalled,
	}
}

func NewTestEnvPrometheusNodeDriver(env *test.Environment) *TestEnvPrometheusNodeDriver {
	return &TestEnvPrometheusNodeDriver{
		env: env,
	}
}

func NewTestEnvOtelNodeDriver(env *test.Environment) *TestEnvOtelNodeDriver {
	return &TestEnvOtelNodeDriver{
		env: env,
	}
}

func (d *TestEnvMetricsClusterDriver) Name() string {
	return "test-environment"
}

func (d *TestEnvMetricsClusterDriver) ShouldDisableNode(*corev1.Reference) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	switch d.state {
	case cortexops.InstallState_NotInstalled, cortexops.InstallState_Uninstalling:
		return status.Error(codes.Unavailable, fmt.Sprintf("Cortex cluster is not installed"))
	case cortexops.InstallState_Updating, cortexops.InstallState_Installed:
		return nil
	case cortexops.InstallState_Unknown:
		fallthrough
	default:
		// can't determine cluster status, so don't disable the node
		return nil
	}
}

func (d *TestEnvMetricsClusterDriver) GetClusterConfiguration(context.Context, *emptypb.Empty) (*cortexops.ClusterConfiguration, error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.state == cortexops.InstallState_NotInstalled {
		return nil, status.Error(codes.NotFound, "Cortex cluster is not installed")
	}
	return d.Configuration, nil
}

func (d *TestEnvMetricsClusterDriver) ConfigureCluster(_ context.Context, conf *cortexops.ClusterConfiguration) (*emptypb.Empty, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	switch d.state {
	case cortexops.InstallState_NotInstalled, cortexops.InstallState_Installed:
		d.state = cortexops.InstallState_Updating
	case cortexops.InstallState_Updating:
		return nil, status.Error(codes.FailedPrecondition, "cluster is already being updated")
	case cortexops.InstallState_Uninstalling:
		return nil, status.Error(codes.FailedPrecondition, "cluster is currently being uninstalled")
	default:
		panic("bug: unknown state")
	}

	oldCtx, oldCancel := d.cortexCtx, d.cortexCancel

	ctx, ca := context.WithCancel(waitctx.FromContext(d.Env.Context()))
	d.cortexCtx = ctx
	d.cortexCancel = ca
	d.Configuration = conf

	go func() {
		if oldCancel != nil {
			oldCancel()
			waitctx.Wait(oldCtx)
		}
		d.Env.StartCortex(ctx)
		d.lock.Lock()
		defer d.lock.Unlock()
		d.state = cortexops.InstallState_Installed
	}()

	return &emptypb.Empty{}, nil
}

func (d *TestEnvMetricsClusterDriver) GetClusterStatus(context.Context, *emptypb.Empty) (*cortexops.InstallStatus, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	return &cortexops.InstallStatus{
		State:    d.state,
		Version:  cortexVersion,
		Metadata: map[string]string{"test-environment": "true"},
	}, nil
}

func (d *TestEnvMetricsClusterDriver) UninstallCluster(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	switch d.state {
	case cortexops.InstallState_NotInstalled:
		return nil, status.Error(codes.FailedPrecondition, "cluster is not installed")
	case cortexops.InstallState_Installed, cortexops.InstallState_Updating:
		d.state = cortexops.InstallState_Uninstalling
	case cortexops.InstallState_Uninstalling:
		return nil, status.Error(codes.FailedPrecondition, "cluster is already being uninstalled")
	default:
		panic("bug: unknown state")
	}

	oldCtx, oldCancel := d.cortexCtx, d.cortexCancel

	go func() {
		if oldCancel != nil {
			oldCancel()
			waitctx.Wait(oldCtx)
		}
		d.lock.Lock()
		defer d.lock.Unlock()
		d.cortexCtx = nil
		d.cortexCancel = nil
		d.state = cortexops.InstallState_NotInstalled
	}()

	return &emptypb.Empty{}, nil
}

type TestEnvPrometheusNodeDriver struct {
	env *test.Environment

	prometheusMu     sync.Mutex
	prometheusCtx    context.Context
	prometheusCancel context.CancelFunc
}

func (d *TestEnvPrometheusNodeDriver) ConfigureRuleGroupFinder(config *v1beta1.RulesSpec) notifier.Finder[rules.RuleGroup] {
	if config.GetDiscovery().GetFilesystem() != nil {
		return rules.NewFilesystemRuleFinder(config.Discovery.Filesystem)
	}
	return nil
}

var _ metrics_agent_drivers.MetricsNodeDriver = (*TestEnvPrometheusNodeDriver)(nil)

// ConfigureNode implements drivers.MetricsNodeDriver
func (d *TestEnvPrometheusNodeDriver) ConfigureNode(nodeId string, conf *node.MetricsCapabilityConfig) error {
	lg := d.env.Logger.With(
		"node", nodeId,
		"driver", "prometheus",
	)
	lg.Debug("configuring node")

	d.prometheusMu.Lock()
	defer d.prometheusMu.Unlock()

	exists := d.prometheusCtx != nil && d.prometheusCancel != nil
	shouldExist := conf.Enabled && conf.GetSpec().GetPrometheus() != nil

	if exists && !shouldExist {
		lg.Info("stopping prometheus")
		d.prometheusCancel()
		waitctx.Wait(d.prometheusCtx)
		d.prometheusCancel = nil
		d.prometheusCtx = nil
	} else if !exists && shouldExist {
		lg.Info("starting prometheus")
		ctx, ca := context.WithCancel(d.env.Context())
		ctx = waitctx.FromContext(ctx)
		// this is the only place where UnsafeStartPrometheus is safe
		_, err := d.env.UnsafeStartPrometheus(ctx, nodeId)
		if err != nil {
			lg.Error("failed to start prometheus", logger.Err(err))
			ca()
			return err
		}
		lg.Info("started prometheus")
		d.prometheusCtx = ctx
		d.prometheusCancel = ca
	} else if exists && shouldExist {
		lg.Debug("nothing to do (already running)")
	} else {
		lg.Debug("nothing to do (already stopped)")
	}

	return nil
}

// DiscoverPrometheuses implements drivers.MetricsNodeDriver
func (*TestEnvPrometheusNodeDriver) DiscoverPrometheuses(context.Context, string) ([]*remoteread.DiscoveryEntry, error) {
	return nil, nil
}

type TestEnvOtelNodeDriver struct {
	env *test.Environment

	otelMu     sync.Mutex
	otelCtx    context.Context
	otelCancel context.CancelFunc
}

// ConfigureNode implements drivers.MetricsNodeDriver.
func (d *TestEnvOtelNodeDriver) ConfigureNode(nodeId string, conf *node.MetricsCapabilityConfig) error {
	lg := d.env.Logger.With(
		"node", nodeId,
		"driver", "otel",
	)
	lg.Debug("configuring node")

	d.otelMu.Lock()
	defer d.otelMu.Unlock()

	exists := d.otelCtx != nil && d.otelCancel != nil
	shouldExist := conf.Enabled && conf.GetSpec().GetOtel() != nil

	if exists && !shouldExist {
		lg.Info("stopping otel")
		d.otelCancel()
		waitctx.Wait(d.otelCtx)
		d.otelCancel = nil
		d.otelCtx = nil
	} else if !exists && shouldExist {
		lg.Info("starting otel")
		ctx, ca := context.WithCancel(d.env.Context())
		ctx = waitctx.FromContext(ctx)
		err := d.env.StartOTELCollectorContext(ctx, nodeId, node.CompatOTELStruct(conf.GetSpec().GetOtel()))
		if err != nil {
			lg.Error("failed to configure otel collector", logger.Err(err))
			ca()
			return fmt.Errorf("failed to configure otel collector: %w", err)
		}
		d.otelCtx = ctx
		d.otelCancel = ca
	} else if exists && shouldExist {
		lg.Debug("nothing to do (already running)")
	} else {
		lg.Debug("nothing to do (already stopped)")
	}

	return nil
}

// ConfigureRuleGroupFinder implements drivers.MetricsNodeDriver.
func (*TestEnvOtelNodeDriver) ConfigureRuleGroupFinder(config *v1beta1.RulesSpec) notifier.Finder[rules.RuleGroup] {
	if config.GetDiscovery().GetFilesystem() != nil {
		return rules.NewFilesystemRuleFinder(config.Discovery.Filesystem)
	}
	return nil
}

// DiscoverPrometheuses implements drivers.MetricsNodeDriver.
func (*TestEnvOtelNodeDriver) DiscoverPrometheuses(context.Context, string) ([]*remoteread.DiscoveryEntry, error) {
	return nil, nil
}

var _ metrics_agent_drivers.MetricsNodeDriver = (*TestEnvOtelNodeDriver)(nil)
