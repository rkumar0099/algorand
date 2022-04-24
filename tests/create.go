package tests

import (
	"log"
	"time"

	"github.com/rkumar0099/algorand/api"
	"github.com/rkumar0099/algorand/crypto"
)

var a1 *api.API = api.New("127.0.0.1:9020")
var a2 *api.API = api.New("127.0.0.1:9021")
var a3 *api.API = api.New("127.0.0.1:9022")

func TestCreate(username string, password string) (*crypto.PublicKey, *crypto.PrivateKey) {
	pk, sk, msg := a1.CreateAccount(username, password)
	log.Println(msg)
	return pk, sk
}

func TestPriceFeed() {
	res := make(chan bool, 3)
	go a1.PriceFeed(1, res)
	go a2.PriceFeed(2, res)
	go a3.PriceFeed(3, res)

	for len(res) < 3 {
		time.Sleep(1 * time.Second)
	}

}
