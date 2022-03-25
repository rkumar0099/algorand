package oracle

import (
	context "context"

	"github.com/rkumar0099/algorand/gossip"
	grpc "google.golang.org/grpc"
)

type OracleServiceServer struct {
	UnimplementedRPCOracleServiceServer
	nodeId         gossip.NodeId
	sendOPPHandler func(uint64) (*ResOPP, error)
}

func NewServer(id gossip.NodeId, oppHandler func(uint64) (*ResOPP, error)) *OracleServiceServer {
	return &OracleServiceServer{
		UnimplementedRPCOracleServiceServer{},
		id,
		oppHandler,
	}
}

func (server *OracleServiceServer) Register(grpcServer *grpc.Server) {
	RegisterRPCOracleServiceServer(grpcServer, server)
}

func (s *OracleServiceServer) SendOPP(ctx context.Context, req *ReqOPP) (*ResOPP, error) {
	return s.sendOPPHandler(req.Epoch)
}
