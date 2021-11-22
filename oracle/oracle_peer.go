package oracle

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"
	"github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
	mt "github.com/wealdtech/go-merkletree"
)

type OraclePeer struct {
	pubkey  *crypto.PublicKey
	privkey *crypto.PrivateKey
	//store   kvstore.KVStore
	epoch   uint64
	results map[uint64][]byte
	data    [][]byte
	tree    *mt.MerkleTree
	//db      *leveldb.DB
}

func newOraclePeer(epoch uint64) *OraclePeer {
	op := &OraclePeer{
		results: make(map[uint64][]byte),
		epoch:   epoch,
	}
	op.pubkey, op.privkey, _ = crypto.NewKeyPair()
	//op.db, err := leveldb.OpenFile("../oracle/blockDB", nil)

	return op
}

type BTC_TO_CUR struct {
	Time  string  `json:"time"`
	Base  string  `json:"asset_id_base"`
	Quote string  `json:"asset_id_quote"`
	Rate  float64 `json:"rate"`
}

func (op *OraclePeer) commit(epoch uint64, txs []*message.PendingRequest, res map[cmn.Hash][]byte) {
	//log.Println("Peer commit ", epoch, len(txs))
	for _, tx := range txs {
		op.fetchURL(tx.URL, tx.Id)
	}
	log.Println(len(op.data))
	op.tree, _ = mt.New(op.data)
	log.Println(op.tree)
	res[cmn.BytesToHash(op.pubkey.Bytes())] = op.tree.Root()
}

func (op *OraclePeer) fetchURL(url string, Id uint64) {
	resp, err := http.Get(url)
	log.Println(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("error accessing data")
			return
		}
		log.Println(respBytes)
		res := &BTC_TO_CUR{}
		json.Unmarshal(respBytes, res)
		data, _ := json.Marshal(res)
		log.Println(data)
		//fmt.Println(respBytes)
		op.data = append(op.data, data)
		//log.Println(len(op.data))
		op.results[Id] = data
		log.Println(len(op.results))
	}
}

func (op *OraclePeer) reveal(epoch uint64, jobIds []uint64, res map[cmn.Hash]*Reveal) {
	reveal := &Reveal{
		tree:    op.tree,
		addr:    op.pubkey.Bytes(),
		results: make(map[uint64][]byte),
	}
	for _, id := range jobIds {
		reveal.results[id] = op.results[id]
	}
	res[cmn.BytesToHash(op.pubkey.Bytes())] = reveal
}

// role returns the role bytes from current round and step
func opRole(pk []byte, iden string, round uint64, event string) []byte {
	return bytes.Join([][]byte{
		[]byte(iden),
		common.Uint2Bytes(round),
		[]byte(event),
	}, nil)
}

type SortitionSeed func(uint64) []byte

//func (op *OraclePeer) Propose(nonce uint64, ss SortitionSeed, blkChan chan *FinalBlock, epoch uint64) {
//	op.proposeBlock(nonce, ss, blkChan, epoch)
//}

func (op *OraclePeer) proposeBlock(nonce uint64, ss SortitionSeed, blkChan chan *FinalBlock,
	epoch uint64, res map[uint64][][]byte) {
	//salt, _ := op.privkey.Sign(cmn.Uint2Bytes(epoch))
	role := role(params.OracleBlockProposer, epoch, params.ORACLE_BLOCK_PROPOSAL)
	seed := ss(epoch)
	finalSeed := bytes.Join([][]byte{
		seed,
		cmn.Uint2Bytes(nonce),
	}, nil)
	proof, vrf, subusers := op.sortition(finalSeed, role, 20, uint64(10000))
	if subusers > 0 {
		log.Println("Selected", op.pubkey.Bytes())
		blk := &FinalBlock{
			Proof:  proof,
			VRF:    vrf,
			Pubkey: op.pubkey.Bytes(),
		}

		blk.Results = make(map[uint64][]byte)
		for Id, results := range res {
			ind := int(len(results) / 2)
			blk.Results[Id] = results[ind]
		}

		if len(blkChan) < 1 {
			blkChan <- blk
		}
	}
}

// sortition runs cryptographic selection procedure and returns vrf,proof and amount of selected sub-users.
func (op *OraclePeer) sortition(seed, role []byte, expectedNum int, weight uint64) (vrf, proof []byte, selected int) {
	vrf, proof, _ = op.privkey.Evaluate(constructSeed(seed, role))
	selected = cmn.SubUsers(expectedNum, weight, vrf)
	return
}
