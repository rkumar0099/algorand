// All external clients can communicate with blockchain using this package

// If anyone wants to create an account on blockchain, they can make use of create_account function exposed by this api

// Once the account is created, they will get pub/priv key pair, they will use this pair to
// perform other transactions

package api

import (
	"fmt"
	"log"

	"github.com/rkumar0099/algorand/client"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"
	"github.com/rkumar0099/algorand/gossip"
	"google.golang.org/grpc"
)

type API struct {
	client  *client.ClientServiceServer // to communicate with blockchain over grpc
	pubkey  *crypto.PublicKey
	privkey *crypto.PrivateKey
}

func New() *API {
	a := &API{}
	id := gossip.NewNodeId("127.0.0.1:9020")
	a.client = client.New(id, a.sendReqHandler, a.sendResHandler)
	a.client.Register(grpc.NewServer())
	return &API{}
}

func (a *API) sendReqHandler(req *client.ReqTx) (*client.ResEmpty, error) {
	return &client.ResEmpty{}, nil
}

func (a *API) sendResHandler(res *client.ResTx) (*client.ResEmpty, error) {
	// handle the response received from blockchain network for the req sent

	return &client.ResEmpty{}, nil
}

func (a *API) CreateAccount(username string, password string) bool {
	// create account by sending the transaction to blockchain
	// over 2/3 of peers must execute the transaction in order to create a new account
	pk, sk, _ := crypto.NewKeyPair()
	a.pubkey = pk
	a.privkey = sk

	passHash := cmn.Sha256([]byte(password))

	c := &client.Create{
		Username: username,
		Password: passHash.Bytes(),
		Pubkey:   pk.Bytes(),
	}
	data, _ := c.Serialize()
	req := &client.ReqTx{
		Type: 1,
		Addr: "127.0.0.1:9020",
		Data: data,
	}

	// send the tx to all peers
	a.sendReq(req)

	return true
}

func (a *API) LogIn(username string, password string, pubkey *crypto.PublicKey, privkey *crypto.PrivateKey) bool {
	// send the credentials to blockchain to see if there is account for this address
	// private key is important to sign the transactions
	passHash := cmn.Sha256([]byte(password))
	l := &client.LogIn{
		Username: username,
		Password: passHash.Bytes(),
		Pubkey:   pubkey.Bytes(),
	}
	data, _ := l.Serialize()
	req := &client.ReqTx{
		Type: 2,
		Addr: "127.0.0.1:9020",
		Data: data,
	}

	a.sendReq(req)

	return true
}

func (a *API) LogOut(username string, password string, pubkey *crypto.PublicKey) {
	// logout this user
}

func (ac *API) Topup(amount uint) bool {
	return true
}

func (a *API) Transfer(to crypto.PublicKey, amt uint) bool {
	// perform transfer, return true if tx successful
	return true
}

// Our system supports different price feeds
// 1 - BTC/USD
// 2 - ETH/USD
// And so on
func (a *API) PriceFeed(TYPE int) bool {
	return true
}

// Our system also aims to support SessionData type whose value may not be integer but a schema
// future work
func (a *API) SessionData(form int) bool {
	return true
}

func (a *API) sendReq(req *client.ReqTx) {
	for i := 0; i < 50; i++ {
		id := gossip.NewNodeId(fmt.Sprintf("127.0.0.1:%d", 8000+i))
		conn, _ := id.Dial()
		_, err := client.SendReqTx(conn, req)
		if err == nil {
			log.Println("[Debug] [API] Tx send successfully")
		}
	}
}
