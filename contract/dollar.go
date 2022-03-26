package contract

import (
	"log"

	"github.com/rkumar0099/algorand/api"
)

func Testing_Price_feed() {
	ac := api.New()
	log.Println("[Debug] [PriceFeed] API object created to communicate with blockchain")
	ac.CreateAccount()

	// write smart contract logic here
}
