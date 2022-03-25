package service

import (
	context "context"
	"encoding/hex"

	"github.com/golang/protobuf/proto"

	msg "github.com/rkumar0099/algorand/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetBlock(conn *grpc.ClientConn, hash []byte) (*msg.Block, error) {
	client := NewRPCServiceClient(conn)
	req := &ReqHash{
		Hash: hash,
	}
	res, err := client.GetDataByHash(context.Background(), req)
	if err != nil {
		return nil, err
	}
	blk := &msg.Block{}
	err = proto.Unmarshal(res.Data, blk)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "[service] hash %s is not a block", hex.EncodeToString(hash))
	}
	return blk, nil
}

func SendFinalContribution(conn *grpc.ClientConn, data []byte) error {
	client := NewRPCServiceClient(conn)
	req := &ReqContribution{
		Data: data,
	}
	_, err := client.SendFinalContribution(context.Background(), req)
	return err
}

/*
func Connect(addr string) (StoreShareClient, error) {
	conn, err := grpc.Dial(addr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[RCP] Cannot connect to %s: %s", addr, err.Error()))
	}
	client := NewStoreShareClient(conn)
	return client, nil
}

func GetBlock(addrs []string, hash []byte) (*msg.Block, error) {
	for i := range addrs {
		conn, err := grpc.Dial(addrs[i], grpc.WithInsecure())
		if err != nil {
			common.Log.Errorf("[RPC] Cannot dial %s: %s", addrs[i], err.Error())
			continue
		}
		client := NewStoreShareClient(conn)
		res, err := client.GetBlock(context.Background(), &ReqHash{Hash: hash})
		conn.Close()
		if err != nil {
			common.Log.Errorf("[RPC] Cannot get %s from %s: %s", hex.EncodeToString(hash), addrs[i], err.Error())
			continue
		}
		if bytes.Equal(common.Sha256(res.Data).Bytes(), hash) {
			blk := &msg.Block{}
			err = blk.Deserialize(res.Data)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("[RPC] cannot deserialize: %s", hex.EncodeToString(hash)))
			}
			return blk, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("[RPC] Not found: %s", hex.EncodeToString(hash)))
}
*/
