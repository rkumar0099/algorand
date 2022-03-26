package client

import (
	context "context"

	"github.com/rkumar0099/algorand/gossip"
	"google.golang.org/grpc"
)

// know the list of peers

type ClientServiceServer struct {
	UnimplementedRPCClientServiceServer
	Id               gossip.NodeId
	SendTxHandler    func(*ReqTx) (*ResEmpty, error)
	SendResTxHandler func(*ResTx) (*ResEmpty, error)
}

func New(
	id gossip.NodeId,
	sendTxHandler func(*ReqTx) (*ResEmpty, error),
	sendResTxHandler func(*ResTx) (*ResEmpty, error),
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

func (css *ClientServiceServer) SendTx(ctx context.Context, req *ReqTx) (*ResEmpty, error) {
	return css.SendTxHandler(req)
}

func (css *ClientServiceServer) SendResTx(ctx context.Context, res *ResTx) (*ResEmpty, error) {
	return css.SendResTxHandler(res)
}

func SendReqTx(conn *grpc.ClientConn, req *ReqTx) (*ResEmpty, error) {
	c := NewRPCClientServiceClient(conn)
	return c.SendTx(context.Background(), req)
}

func SendResTx(conn *grpc.ClientConn, res *ResTx) (*ResEmpty, error) {
	c := NewRPCClientServiceClient(conn)
	return c.SendResTx(context.Background(), res)
}
