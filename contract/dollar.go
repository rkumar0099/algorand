package contract

import (
	"log"

	"github.com/rkumar0099/algorand/api"
	"github.com/rkumar0099/algorand/gossip"
)

func Testing_Price_feed() {
	addr := gossip.NewNodeId("127.0.0.1:9020")
	ac := api.New(addr.String())
	log.Printf("[Debug] [Price Feed] Algorand client %s created to communicate with blockchain", ac.Address)

	// write smart contract logic here
}
