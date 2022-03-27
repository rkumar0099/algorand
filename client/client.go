package client

import (
	context "context"

	"github.com/golang/protobuf/proto"
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

func (c *Create) Serialize() ([]byte, error) {
	return proto.Marshal(c)
}

func (c *Create) Deserialize(data []byte) error {
	return proto.Unmarshal(data, c)
}

func (logIn *LogIn) Serialize() ([]byte, error) {
	return proto.Marshal(logIn)
}

func (logIn *LogIn) Deserialize(data []byte) error {
	return proto.Unmarshal(data, logIn)
}

func (logOut *LogOut) Serialize() ([]byte, error) {
	return proto.Marshal(logOut)
}

func (logOut *LogOut) Deserialize(data []byte) error {
	return proto.Unmarshal(data, logOut)
}

func (topUp *TopUp) Serialize() ([]byte, error) {
	return proto.Marshal(topUp)
}

func (topUp *TopUp) Deserialize(data []byte) error {
	return proto.Unmarshal(data, topUp)
}

func (transfer *Transfer) Serialize() ([]byte, error) {
	return proto.Marshal(transfer)
}

func (transfer *Transfer) Deserialize(data []byte) error {
	return proto.Unmarshal(data, transfer)
}

func (pf *PriceFeed) Serialize() ([]byte, error) {
	return proto.Marshal(pf)
}

func (pf *PriceFeed) Deserialize(data []byte) error {
	return proto.Unmarshal(data, pf)
}

func (sd *SessionData) Serialize() ([]byte, error) {
	return proto.Marshal(sd)
}

func (sd *SessionData) Deserialize(data []byte) error {
	return proto.Unmarshal(data, sd)
}

func (res *ResTx) Serialize() ([]byte, error) {
	return proto.Marshal(res)
}

func (res *ResTx) Deserialize(data []byte) error {
	return proto.Unmarshal(data, res)
}
