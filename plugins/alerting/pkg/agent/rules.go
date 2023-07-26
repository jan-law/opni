package agent

import (
	"context"
	"log/slog"
	"time"

	healthpkg "github.com/rancher/opni/pkg/health"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/util"
	"github.com/rancher/opni/plugins/alerting/pkg/agent/drivers"
	"github.com/rancher/opni/plugins/alerting/pkg/apis/node"
	"github.com/rancher/opni/plugins/alerting/pkg/apis/rules"
)

var RuleSyncInterval = time.Minute * 2

type RuleStreamer struct {
	util.Initializer

	parentCtx context.Context

	lg *slog.Logger

	ruleStreamCtx  context.Context
	stopRuleStream context.CancelFunc
	ruleSyncClient rules.RuleSyncClient

	conditions healthpkg.ConditionTracker
	nodeDriver drivers.NodeDriver
}

var _ drivers.ConfigPropagator = (*RuleStreamer)(nil)

func NewRuleStreamer(
	ctx context.Context,
	lg *slog.Logger,
	ct healthpkg.ConditionTracker,
	nodeDriver drivers.NodeDriver,
) *RuleStreamer {
	return &RuleStreamer{
		parentCtx:  ctx,
		lg:         lg,
		conditions: ct,
		nodeDriver: nodeDriver,
	}
}

func (r *RuleStreamer) Initialize(ruleSyncClient rules.RuleSyncClient) {
	r.InitOnce(func() {
		r.ruleSyncClient = ruleSyncClient
	})
}

func (r *RuleStreamer) ConfigureNode(nodeId string, cfg *node.AlertingCapabilityConfig) error {
	return r.configureRuleStreamer(nodeId, cfg)
}

func (r *RuleStreamer) configureRuleStreamer(nodeId string, cfg *node.AlertingCapabilityConfig) error {
	lg := r.lg.With("nodeId", nodeId)
	lg.Debug("alerting capability updated")

	currentlyRunning := r.stopRuleStream != nil
	shouldRun := cfg.GetEnabled()

	startRuleStreamer := func() {
		ctx, ca := context.WithCancel(r.parentCtx)
		r.stopRuleStream = ca
		go r.run(ctx)
	}

	switch {
	case currentlyRunning && shouldRun:
		lg.Debug("restarting rule stream")
		r.stopRuleStream()
		startRuleStreamer()
	case currentlyRunning && !shouldRun:
		lg.Debug("stopping rule stream")
		r.stopRuleStream()
	case !currentlyRunning && shouldRun:
		lg.Debug("starting rule stream")
		startRuleStreamer()
	case !currentlyRunning && !shouldRun:
		lg.Debug("rule sync is disabled")
	}
	return nil
}

func (r *RuleStreamer) sync(ctx context.Context) {
	ruleManifest, err := r.nodeDriver.DiscoverRules(ctx)
	if err != nil {
		r.lg.Warn("failed to discover rules", logger.Err(err))
	}
	r.lg.Info("discovered rules", "count", len(ruleManifest.Rules))
	if _, err := r.ruleSyncClient.SyncRules(ctx, ruleManifest); err != nil {
		r.lg.Warn("failed to sync rules", logger.Err(err))
	}
}

func (r *RuleStreamer) run(ctx context.Context) {
	r.lg.Info("waiting for rule sync client...")
	r.WaitForInitContext(ctx)
	r.lg.Info("rule sync client acquired")
	r.sync(ctx)
	t := time.NewTicker(RuleSyncInterval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			r.sync(ctx)
		case <-ctx.Done():
			r.lg.Info("Exiting rule sync loop")
			return
		}
	}
}
