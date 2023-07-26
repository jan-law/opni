package testbench

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"emperror.dev/errors"
	"github.com/golang/snappy"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/rancher/opni/internal/bench"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/test"
)

// Adapted from cortex benchtool

type BenchRunner struct {
	client    remote.WriteClient
	Logger    *slog.Logger
	batchChan chan bench.BatchReq
}

func NewBenchRunner(env *test.Environment, agent string, desc bench.WorkloadDesc) (*BenchRunner, error) {
	runningAgent := env.GetAgent(agent)
	if runningAgent.Agent == nil {
		return nil, fmt.Errorf("agent %s not found", agent)
	}
	u, err := url.Parse(fmt.Sprintf("http://%s/api/agent/push", runningAgent.Agent.ListenAddress()))
	if err != nil {
		return nil, err
	}
	writeClient, err := remote.NewWriteClient(agent, &remote.ClientConfig{
		URL:              &config.URL{URL: u},
		HTTPClientConfig: config.HTTPClientConfig{},
		RetryOnRateLimit: true,
		Timeout:          model.Duration(10 * time.Second),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error creating remote write client")
	}

	workload := bench.NewWriteWorkload(desc, nil)
	batchChan := make(chan bench.BatchReq, 100)
	go func() {
		if err := workload.GenerateWriteBatch(env.Context(), agent, 100, batchChan); err != nil {
			env.Logger.Error("error generating write batch: ", logger.Err(err))
		}
	}()

	return &BenchRunner{
		client:    writeClient,
		Logger:    env.Logger.WithGroup("bench").With("agent", agent),
		batchChan: batchChan,
	}, nil
}

func (b *BenchRunner) StartWorker(ctx context.Context) {
	go func() {
		b.Logger.Info("worker started")
		for {
			select {
			case batchReq := <-b.batchChan:
				if err := b.sendBatch(ctx, batchReq.Batch); err != nil {
					b.Logger.Error("failed to send batch", logger.Err(err))
				}
				batchReq.PutBack <- batchReq.Batch
				batchReq.Wg.Done()
			case <-ctx.Done():
				b.Logger.Warn("worker stopped", "error", ctx.Err())
				return
			}
		}
	}()
}

func (b *BenchRunner) sendBatch(ctx context.Context, batch []prompb.TimeSeries) error {
	req := prompb.WriteRequest{
		Timeseries: batch,
	}

	data, err := req.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal remote-write request")
	}

	compressed := snappy.Encode(nil, data)

	return b.client.Store(ctx, compressed)
}
