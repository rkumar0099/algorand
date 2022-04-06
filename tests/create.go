package tests

import (
	"log"

	"github.com/rkumar0099/algorand/api"
	"github.com/rkumar0099/algorand/crypto"
)

var a *api.API = api.New("127.0.0.1:9020")

func TestCreate(username string, password string) (*crypto.PublicKey, *crypto.PrivateKey) {
	pk, sk, msg := a.CreateAccount(username, password)
	log.Println(msg)
	return pk, sk
}
