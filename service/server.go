package service

import (
	context "context"
	"encoding/hex"

	"github.com/rkumar0099/algorand/gossip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedRPCServiceServer
	nodeId                       gossip.NodeId
	getDataByHashHandler         func([]byte) ([]byte, error)
	sendContributionHandler      func([]byte) ([]byte, error)
	sendFinalContributionHandler func([]byte) error
}

func NewServer(
	nodeId gossip.NodeId,
	getDataByHashHandler func([]byte) ([]byte, error),
	sendContributionHandler func([]byte) ([]byte, error),
	sendFinalContributionHandler func([]byte) error,
) *Server {
	return &Server{
		UnimplementedRPCServiceServer{},
		nodeId,
		getDataByHashHandler,
		sendContributionHandler,
		sendFinalContributionHandler,
	}
}

func (server *Server) Register(grpcServer *grpc.Server) {
	RegisterRPCServiceServer(grpcServer, server)
}

func (server *Server) GetDataByHash(ctx context.Context, in *ReqHash) (*ResData, error) {
	/*value, err := server.memKv.Get(hex.EncodeToString(in.Hash))
	if err == nil {
		if blk, ok := value.(*msg.Block); ok {
			data, _ := blk.Serialize()
			return &ResData{Data: data}, nil
		} else {
			return nil, errors.New(fmt.Sprintf("Hash %s is not a block,", hex.EncodeToString(in.Hash)))
		}
	}
	data, err := server.persistKv.Get(in.Hash)*/
	data, err := server.getDataByHashHandler(in.Hash)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "[from %s] Hash %s not found: %s", server.nodeId.String(), hex.EncodeToString(in.Hash), err.Error())
	}
	return &ResData{Data: data}, err
}

func (s *Server) SendFinalContribution(ctx context.Context, req *ReqContribution) (*ResContribution, error) {
	err := s.sendFinalContributionHandler(req.Data)
	return &ResContribution{}, err
}
