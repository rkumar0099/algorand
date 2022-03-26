package contract

import (
	"log"

	"github.com/rkumar0099/algorand/api"
	"github.com/rkumar0099/algorand/gossip"
)

func Testing_Flight_Feed() {
	Id := gossip.NewNodeId("127.0.0.1:9022")
	ac := api.New(Id.String())
	log.Printf("[Debug] [Flight] Algorand client %s created to communicate with blockchain", ac.Address)

	// write your smart contract logic here

}
