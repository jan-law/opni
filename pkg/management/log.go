package management

import (
	"fmt"

	controlv1 "github.com/rancher/opni/pkg/apis/control/v1"
	managementv1 "github.com/rancher/opni/pkg/apis/management/v1"
	"github.com/rancher/opni/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (m *Server) StreamAgentLogs(
	req *managementv1.StreamAgentLogsRequest,
	stream managementv1.Management_StreamAgentLogsServer,
) error {
	cc, err := grpc.DialContext(stream.Context(), m.config.GRPCListenAddress,
		grpc.WithBlock(),
		grpc.WithContextDialer(util.DialProtocol),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial grpc server: %w", err)
	}

	logClient := controlv1.NewLogClient(cc)

	logReq := &controlv1.LogStreamRequest{
		Since:   req.Request.Since,
		Until:   req.Request.Until,
		Filters: req.Request.Filters,
	}

	agentStream, err := logClient.StreamLogs(stream.Context(), logReq)
	if err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			log, err := agentStream.Recv()
			if err != nil {
				return err
			}
			if err := stream.Send(log); err != nil {
				return err
			}
		}
	}
}
