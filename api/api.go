// All external clients can communicate with blockchain using this package

// If anyone wants to create an account on blockchain, they can make use of create_account function exposed by this api

// Once the account is created, they will get pub/priv key pair, they will use this pair to
// perform other transactions

package api

import (
	"crypto"

	"github.com/rkumar0099/algorand/client"
	"github.com/rkumar0099/algorand/gossip"
	"google.golang.org/grpc"
)

type API struct {
	client  *client.ClientServiceServer // to communicate with blockchain over grpc
	Pubkey  *crypto.PublicKey
	Privkey *crypto.PrivateKey
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

func (a *API) CreateAccount() bool {
	// create account by sending the transaction to blockchain
	// over 2/3 of peers must execute the transaction in order to create a new account

	return true
}

func (a *API) LogIn(pubkey *crypto.PublicKey, privkey *crypto.PrivateKey) bool {
	// send the credentials to blockchain to see if there is account for this address
	return true
}

func (a *API) Transfer(to crypto.PublicKey, amt uint) bool {
	// perform transfer, return true if tx successful
	return true
}

func (ac *API) Topup(amount uint) bool {
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
