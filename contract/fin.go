package contract

import (
	"log"

	"github.com/rkumar0099/algorand/api"
)

func Testing_FIN_Feed() {
	ac := api.New()
	log.Println("[Debug] [SessionData] [FIN] API object created to communicate with blockchain")
	ac.CreateAccount()

	// write smart contract logic here

}
