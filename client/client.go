package client

import (
	"github.com/rkumar0099/algorand/gossip"
	"google.golang.org/grpc"
)

// know the list of peers

type ClientServiceServer struct {
	UnimplementedRPCClientServiceServer
	Id               gossip.NodeId
	SendTxHandler    func(ReqTx) (ResEmpty, error)
	SendResTxHandler func(ResTx) (ResEmpty, error)
}

func New(
	id gossip.NodeId,
	sendTxHandler func(ReqTx) (ResEmpty, error),
	sendResTxHandler func(ResTx) (ResEmpty, error),
) *ClientServiceServer {
	css := &ClientServiceServer{
		UnimplementedRPCClientServiceServer{},
		id,
		sendTxHandler,
		sendResTxHandler,
	}
	return css
}

func (css *ClientServiceServer) Register(s *grpc.Server) {
	RegisterRPCClientServiceServer(s, css)
}

func (css *ClientServiceServer) SendTx(ReqTx) (ResEmpty, error) {
	re := ResEmpty{}
	return re, nil
}

func (css *ClientServiceServer) SendResTx(ResTx) (ResEmpty, error) {
	re := ResEmpty{}
	return re, nil
}
