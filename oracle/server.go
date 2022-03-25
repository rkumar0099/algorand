package oracle

import (
	context "context"

	"github.com/rkumar0099/algorand/gossip"
	grpc "google.golang.org/grpc"
)

type OracleServiceServer struct {
	UnimplementedRPCServiceServer
	nodeId    gossip.NodeId
	OPPResult func() ([]byte, error)
}

func NewServer(id gossip.NodeId, opp func() ([]byte, error)) *OracleServiceServer {
	return &OracleServiceServer{
		UnimplementedRPCServiceServer{},
		id,
		opp,
	}
}

func (server *OracleServiceServer) Register(grpcServer *grpc.Server) {
	RegisterRPCServiceServer(grpcServer, server)
}

func (s *OracleServiceServer) OPP(ctx context.Context) ([]byte, error) {

}
