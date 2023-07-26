package gateway

import (
	"context"
	"log/slog"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/kralicky/totem"
	agentv1 "github.com/rancher/opni/pkg/agent"
	corev1 "github.com/rancher/opni/pkg/apis/core/v1"
	streamv1 "github.com/rancher/opni/pkg/apis/stream/v1"
	"github.com/rancher/opni/pkg/auth/cluster"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/storage"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type agentInfo struct {
	grpc.ClientConnInterface

	labels map[string]string
	id     string
}

func (a agentInfo) GetLabels() map[string]string {
	return a.labels
}

func (a agentInfo) GetId() string {
	return a.id
}

type DelegateServer struct {
	streamv1.UnsafeDelegateServer
	mu           sync.RWMutex
	activeAgents map[string]agentInfo
	logger       *slog.Logger
	clusterStore storage.ClusterStore
}

func NewDelegateServer(clusterStore storage.ClusterStore, lg *slog.Logger) *DelegateServer {
	return &DelegateServer{
		activeAgents: make(map[string]agentInfo),
		clusterStore: clusterStore,
		logger:       lg.WithGroup("delegate"),
	}
}

func (d *DelegateServer) HandleAgentConnection(ctx context.Context, clientSet agentv1.ClientSet) {
	d.mu.Lock()
	id := cluster.StreamAuthorizedID(ctx)

	cluster, err := d.clusterStore.GetCluster(ctx, &corev1.Reference{Id: id})
	if err != nil {
		d.logger.Error("internal error: failed to look up connecting agent", logger.Err(err), "id", id)
		return
	}

	d.activeAgents[id] = agentInfo{
		ClientConnInterface: clientSet.ClientConn(),
		labels:              cluster.GetLabels(),
		id:                  id,
	}
	d.logger.Debug("agent connected", "id", id)
	d.mu.Unlock()

	<-ctx.Done()

	d.mu.Lock()
	delete(d.activeAgents, id)
	d.logger.Debug("agent disconnected", "id", id)
	d.mu.Unlock()
}

func (d *DelegateServer) Request(ctx context.Context, req *streamv1.DelegatedMessage) (*totem.RPC, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	targetId := req.GetTarget().GetId()
	lg := d.logger.With(
		"target", targetId,
		"request", req.GetRequest().QualifiedMethodName(),
	)
	lg.Debug("delegating rpc request")
	target, ok := d.activeAgents[targetId]
	if ok {
		fwdResp := &totem.RPC{}
		err := target.Invoke(ctx, totem.Forward, req.GetRequest(), fwdResp)
		if err != nil {
			d.logger.Error("delegating rpc request failed", logger.Err(err))
			return nil, err
		}

		resp := &totem.RPC{}
		err = proto.Unmarshal(fwdResp.GetResponse().GetResponse(), resp)
		if err != nil {
			d.logger.Error("delegating rpc request failed", logger.Err(err))
			return nil, err
		}

		return resp, nil
	}

	err := status.Error(codes.NotFound, "target not found")
	lg.Warn("delegating rpc request failed", logger.Err(err))
	return nil, err
}

func (d *DelegateServer) Broadcast(ctx context.Context, req *streamv1.BroadcastMessage) (*streamv1.BroadcastReplyList, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	sp := storage.NewSelectorPredicate[agentInfo](req.GetTargetSelector())

	var targets []agentInfo
	for _, aa := range d.activeAgents {
		if sp(aa) {
			targets = append(targets, aa)
		}
	}

	if len(targets) == 0 {
		return nil, status.Error(codes.NotFound, "no targets found")
	}

	eg, ctx := errgroup.WithContext(ctx)
	reply := &streamv1.BroadcastReplyList{
		Responses: make([]*streamv1.BroadcastReply, len(targets)),
	}

	for i, target := range targets {
		i, target := i, target
		eg.Go(func() error {
			item := &totem.RPC{}
			if err := target.Invoke(ctx, totem.Forward, req.GetRequest(), item); err != nil {
				return err
			}

			reply.Responses[i] = &streamv1.BroadcastReply{
				Ref: &corev1.Reference{
					Id: target.id,
				},
				Reply: item,
			}

			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}

	return reply, nil
}
