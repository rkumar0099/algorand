package tests

import (
	"log"

	"github.com/rkumar0099/algorand/crypto"
)

func TestLogin(username string, password string, pk *crypto.PublicKey) {
	status, msg := a.LogIn(username, password, pk)
	log.Println(status, msg)
}
