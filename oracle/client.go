package oracle

import (
	"context"

	"google.golang.org/grpc"
)

func SendOPP(conn *grpc.ClientConn, epoch uint64) (*ResOPP, error) {
	c := NewRPCOracleServiceClient(conn)
	req := &ReqOPP{Epoch: epoch}
	return c.SendOPP(context.Background(), req)
}
