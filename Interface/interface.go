package Interface

import (
	"fmt"

	"github.com/rkumar0099/algorand/gossip"
	"google.golang.org/grpc"
)

// know the list of peers

type AlgorandClient struct {
	nodes    []gossip.NodeId
	connPool map[gossip.NodeId]*grpc.ClientConn
	Address  string
}

func New(addr string) *AlgorandClient {
	ac := &AlgorandClient{
		Address:  addr,
		connPool: make(map[gossip.NodeId]*grpc.ClientConn),
	}

	for i := 0; i < 50; i++ {
		Id := gossip.NewNodeId(fmt.Sprintf("127.0.0.1:%d", 8000+i))
		ac.nodes = append(ac.nodes, Id)
		conn, err := Id.Dial()
		if err == nil {
			ac.connPool[Id] = conn
		}
	}
	return ac
}

func (ac *AlgorandClient) Transfer(to string, amount uint) {

}

func (ac *AlgorandClient) Topup(amount uint) {

}

// Our system supports different price feeds
// 1 - BTC/USD
// 2 - ETH/USD
// And so on
func (ac *AlgorandClient) PriceFeed(form int) {

}

// Our system also aims to support SessionData type whose value may not be integer but a schema
// future work
func (ac *AlgorandClient) SessionData(form int) {

}
