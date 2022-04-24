package oracle

import (
	"sync"

	"github.com/rkumar0099/algorand/crypto"
	msg "github.com/rkumar0099/algorand/message"
	mt "github.com/wealdtech/go-merkletree"
)

const (
	CURRENCY_EXCHANGE = iota
	FLIGHT
)

// fetch data from given url, and send rate at time t to end client.

type CurrencyExchange struct {
	Time string  `json:"time"`
	Rate float64 `json:"rate"`
}

type FlightDetails struct {
	Base    string `json:"base"`
	Price   string `json:"price"`
	Time    string `json:"time"`
	Details string `json:"details"`
}

// clients can only send external world transactions related to the type of access reponses defined here
// because size of data is large. We consider only important details of access response to have small size data

/*
type Response struct {

}


type  struct {
	Time string `json:"time"`
	Rate string `json:"rate"`
}


var URL map[string]
*/

type OraclePeer struct {
	pubkey  *crypto.PublicKey
	privkey *crypto.PrivateKey
	//store   kvstore.KVStore
	epoch   uint64
	results map[uint64][]byte
	data    [][]byte
	tree    *mt.MerkleTree
	//db      *leveldb.DB
	lock *sync.Mutex
}

type OracleBlock struct {
	Epoch uint64
	Res   []*msg.PendingRequestRes
}

type Reveal struct {
	addr    []byte
	tree    *mt.MerkleTree
	results map[uint64][]byte
}

type FinalBlock struct {
	Results map[uint64][]byte
	Proof   []byte
	VRF     []byte
	Pubkey  []byte
	Epoch   uint64
}
