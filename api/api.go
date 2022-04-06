// All external clients can communicate with blockchain using this package

// If anyone wants to create an account on blockchain, they can make use of create_account function exposed by this api

// Once the account is created, they will get pub/priv key pair, they will use this pair to
// perform other transactions

package api

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/rkumar0099/algorand/client"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"
	"github.com/rkumar0099/algorand/gossip"
	"google.golang.org/grpc"
)

type API struct {
	client      *client.ClientServiceServer // to communicate with blockchain over grpc
	pubkey      *crypto.PublicKey
	privkey     *crypto.PrivateKey
	res         *client.ResTx
	ready       chan bool
	server      *grpc.Server
	id          gossip.NodeId
	acctCreated bool
	login       bool
	logout      bool
}

func New(id string) *API {

	a := &API{
		id:          gossip.NewNodeId(id),
		server:      grpc.NewServer(),
		ready:       make(chan bool, 1),
		res:         &client.ResTx{},
		acctCreated: false,
		login:       false,
		logout:      true,
	}

	a.client = client.New(a.id, a.sendReqHandler, a.sendResHandler)
	a.client.Register(a.server)
	lis, err := net.Listen("tcp", a.id.String())
	if err == nil {
		go a.server.Serve(lis)
		log.Println("[Debug] [API] Listening for responses")
	}
	return a
}

func (a *API) sendReqHandler(req *client.ReqTx) (*client.ResEmpty, error) {
	return &client.ResEmpty{}, nil
}

func (a *API) sendResHandler(res *client.ResTx) (*client.ResEmpty, error) {
	// handle the response received from blockchain network for the req sent
	log.Println("[Debug] [API] Received Response from blockchain")
	//log.Printf("[Debug] [API] [RES] Id: %s, Pk: %s\n", a.id.String(),
	//cmn.BytesToHash(a.pubkey.Bytes()).String())

	a.res = res
	a.ready <- true

	return &client.ResEmpty{}, nil
}

func (a *API) CreateAccount(username string, password string) (*crypto.PublicKey, *crypto.PrivateKey, string) {
	// create account by sending the transaction to blockchain
	// over 2/3 of peers must execute the transaction in order to create a new account
	pk, sk, _ := crypto.NewKeyPair()
	a.pubkey = pk
	a.privkey = sk

	passHash := cmn.Sha256([]byte(password))

	c := &client.Create{
		Username: username,
		Password: passHash.Bytes(),
	}
	data, _ := c.Serialize()
	req := &client.ReqTx{
		Type:   0,
		Addr:   a.id.String(),
		Pubkey: pk.Bytes(),
		Data:   data,
	}

	a.empty()
	log.Printf("Len ready: %d\n", len(a.ready))
	// send the tx to all peers
	a.sendReq(req)
	a.receiveRes()

	if !a.res.Status {
		log.Println("[API] [Error] Account created unsuccessful")
		a.pubkey = nil
		a.privkey = nil
		return nil, nil, a.res.Msg
	}
	return pk, sk, a.res.Msg
}

func (a *API) LogIn(username string, password string, pubkey *crypto.PublicKey) (bool, string) {
	// send the credentials to blockchain to see if there is account for this address
	// private key is important to sign the transactions
	passHash := cmn.Sha256([]byte(password))
	info := &client.LogIn{
		Username: username,
		Password: passHash.Bytes(),
	}
	data, _ := info.Serialize()
	req := &client.ReqTx{
		Type:   1,
		Addr:   a.id.String(),
		Pubkey: pubkey.Bytes(),
		Data:   data,
	}
	a.empty()
	log.Printf("Len ready: %d\n", len(a.ready))
	a.sendReq(req)
	a.receiveRes()
	if a.res.Status {
		log.Println("[API] [LOGIN] Successful login")
		a.login = true
		a.logout = false
	}
	return a.res.Status, a.res.Msg

}

func (a *API) LogOut(username string, password string, pubkey *crypto.PublicKey) {
	// logout this user
}

func (a *API) TopUp(amount uint64) (bool, string) {
	// perform transfer, return true if tx successful
	if !a.login {
		return false, "You must login first"
	}
	info := &client.TopUp{}
	info.From = a.pubkey.Bytes()
	info.Amount = amount
	req := &client.ReqTx{}
	req.Type = 3
	req.Addr = a.id.String()
	req.Pubkey = a.pubkey.Bytes()
	req.Data, _ = info.Serialize()
	a.empty()
	a.sendReq(req)
	a.receiveRes()
	if a.res.Status {
		log.Println("[API] [TOPUP] Successful topup")
	}
	return a.res.Status, a.res.String()

}

func (a *API) Transfer(to *crypto.PublicKey, amount uint64) (bool, string) {
	if !a.login {
		return false, "You must login first"
	}
	info := &client.Transfer{}
	info.From = a.pubkey.Bytes()
	info.To = to.Bytes()
	info.Amount = amount
	req := &client.ReqTx{}
	req.Type = 4
	req.Addr = a.id.String()
	req.Pubkey = a.pubkey.Bytes()
	req.Data, _ = info.Serialize()
	a.empty()
	a.sendReq(req)
	a.receiveRes()
	if a.res.Status {
		log.Println("[API] Successful topup")
	}
	return a.res.Status, a.res.Msg

}

func (a *API) receiveRes() {
	for {
		time.Sleep(1 * time.Second)
		if len(a.ready) > 0 {
			<-a.ready
			break
		}
	}
}

func (a *API) empty() {
	for len(a.ready) > 0 {
		<-a.ready
	}
	a.res = nil
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
	i := 1
	//for i := 0; i < 50; i++ {
	id := gossip.NewNodeId(fmt.Sprintf("127.0.0.1:%d", 8000+i))
	conn, _ := id.Dial()
	_, err := client.SendReqTx(conn, req)
	if err == nil {
		log.Println("[Debug] [API] Tx send successfully")
	}
	//}
}
