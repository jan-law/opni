package backend

import (
	"context"
	"os"

	opnicorev1 "github.com/rancher/opni/pkg/apis/core/v1"
	managementv1 "github.com/rancher/opni/pkg/apis/management/v1"
	"github.com/rancher/opni/pkg/logger"
)

func (b *LoggingBackend) updateClusterMetadata(ctx context.Context, event *managementv1.WatchEvent) error {
	newName, oldName := event.Cluster.Metadata.Labels[opnicorev1.NameLabel], event.Previous.Metadata.Labels[opnicorev1.NameLabel]
	if newName == oldName {
		b.Logger.Debug("cluster was not renamed", "oldName", oldName, "newName", newName)
		return nil
	}

	b.Logger.Debug("newName", newName, "oldName", oldName)

	if err := b.ClusterDriver.StoreClusterMetadata(ctx, event.Cluster.GetId(), newName); err != nil {
		b.Logger.Debug("could not update cluster metadata", logger.Err(err), "cluster", event.Cluster.Id)
		return nil
	}

	return nil
}

func (b *LoggingBackend) watchClusterEvents(ctx context.Context) {
	clusterClient, err := b.MgmtClient.WatchClusters(ctx, &managementv1.WatchClustersRequest{})
	if err != nil {
		b.Logger.Error("failed to watch clusters, existing", logger.Err(err))
		os.Exit(1)
	}

	b.Logger.Info("watching cluster events")

outer:
	for {
		select {
		case <-clusterClient.Context().Done():
			b.Logger.Info("context cancelled, stoping cluster event watcher")
			break outer
		default:
			event, err := clusterClient.Recv()
			if err != nil {
				b.Logger.Error("failed to receive cluster event", logger.Err(err))
				continue
			}

			b.watcher.HandleEvent(event)
		}
	}
}

func (b *LoggingBackend) reconcileClusterMetadata(ctx context.Context, clusters []*opnicorev1.Cluster) (retErr error) {
	for _, cluster := range clusters {
		err := b.ClusterDriver.StoreClusterMetadata(ctx, cluster.GetId(), cluster.Metadata.Labels[opnicorev1.NameLabel])
		if err != nil {
			b.Logger.Warn("could not update cluster metadata", logger.Err(err), "cluster", cluster.Id)
			retErr = err
		}
	}
	return
}
