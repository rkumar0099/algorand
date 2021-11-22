package oracle

import (
	"bytes"

	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"
	"github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
	"github.com/syndtr/goleveldb/leveldb"
	mt "github.com/wealdtech/go-merkletree"
)

type Oracle struct {
	pubkey    *crypto.PublicKey
	privkey   *crypto.PrivateKey
	peers     map[uint64][]*OraclePeer
	round     uint64
	genesis   bool
	epoch     uint64
	prevepoch uint64
	lock      *sync.Mutex
	txLock    *sync.Mutex
	txs       []*message.PendingRequest
	commit    map[uint64]map[cmn.Hash][]byte
	reveal    map[uint64]map[cmn.Hash]*Reveal
	results   map[uint64]map[uint64][][]byte
	finalBlks map[uint64]*FinalBlock
	db        *leveldb.DB
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

func New() *Oracle {
	oracle := &Oracle{
		peers:     make(map[uint64][]*OraclePeer),
		lock:      &sync.Mutex{},
		round:     0,
		epoch:     0,
		prevepoch: 0,
		txLock:    &sync.Mutex{},
		commit:    make(map[uint64]map[cmn.Hash][]byte),
		reveal:    make(map[uint64]map[cmn.Hash]*Reveal),
		results:   make(map[uint64]map[uint64][][]byte),
		finalBlks: make(map[uint64]*FinalBlock),
	}
	//oracle.blockchain = blockchain.NewBlockchain(oracle.store)
	oracle.pubkey, oracle.privkey, _ = crypto.NewKeyPair()
	oracle.db, _ = leveldb.OpenFile("../oracle/blockDB", nil)
	return oracle
}

func (o *Oracle) Address() []byte {
	return o.pubkey.Bytes()
}

func (o *Oracle) AddOPP(opp *message.OraclePeerProposal) {
	o.lock.Lock()
	defer o.lock.Unlock()
	round := opp.Round
	if o.prevepoch == o.epoch {
		o.epoch += 1
	}
	//log.Println("Current epoch ", o.epoch)
	if _, ok := o.peers[o.epoch]; !ok {
		o.peers[o.epoch] = make([]*OraclePeer, 0)
		log.Println("Make peer array again")
	}
	seed := o.sortitionSeed(round)
	role := role(params.OraclePeer, round, params.ORACLE)
	if err := opp.Verify(opp.Proof, constructSeed(seed, role)); err != nil {
		return
	}

	o.peers[o.epoch] = append(o.peers[o.epoch], newOraclePeer(o.epoch))
	//log.Println(len(o.peers[o.epoch]))
}

func (o *Oracle) AddBlk(blk *message.Block) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if o.round < blk.Round {
		data, _ := blk.Serialize()
		o.db.Put(cmn.Uint2Bytes(blk.Round), data, nil)
	}
}

func (o *Oracle) GetPeers() ([]*OraclePeer, uint64) {
	return o.peers[o.epoch], o.epoch
}

func (o *Oracle) GetBlkByRound(round uint64) *message.Block {
	blk := &message.Block{}
	data, _ := o.db.Get(cmn.Uint2Bytes(round), nil)
	blk.Deserialize(data)
	return blk
}

func (o *Oracle) AddEWTx(tx *message.PendingRequest) {
	o.txLock.Lock()
	defer o.txLock.Unlock()
	o.txs = append(o.txs, tx)
}

// sortitionSeed returns the selection seed with a refresh interval R.
func (o *Oracle) sortitionSeed(round uint64) []byte {
	realR := round - 1
	mod := round % params.R
	if realR < mod {
		realR = 0
	} else {
		realR -= mod
	}
	blk := o.GetBlkByRound(realR)
	return blk.Seed
}

// sortition runs cryptographic selection procedure and returns vrf,proof and amount of selected sub-users.
func (o *Oracle) sortition(seed, role []byte, expectedNum int, weight uint64) (vrf, proof []byte, selected int) {
	vrf, proof, _ = o.privkey.Evaluate(constructSeed(seed, role))
	selected = cmn.SubUsers(expectedNum, weight, vrf)
	return
}

// role returns the role bytes from current round and step
func role(iden string, round uint64, event string) []byte {
	return bytes.Join([][]byte{
		[]byte(iden),
		common.Uint2Bytes(round),
		[]byte(event),
	}, nil)
}

// constructSeed construct a new bytes for vrf generation.
func constructSeed(seed, role []byte) []byte {
	return bytes.Join([][]byte{seed, role}, nil)
}

func (o *Oracle) RunOracle() *FinalBlock {
	//for {
	//time.Sleep(30 * time.Second)
	epoch := o.epoch
	peers := o.peers[epoch]
	//log.Println(len(peers))

	var (
		tx  *message.PendingRequest
		txs []*message.PendingRequest
	)
	for len(txs) < 20 && len(o.txs) > 0 {
		tx, o.txs = o.txs[0], o.txs[1:]
		txs = append(txs, tx)
	}
	return o.process(epoch, peers, txs)
	//o.prevepoch = o.epoch
	//}
}

func (o *Oracle) process(epoch uint64, peers []*OraclePeer, txs []*message.PendingRequest) *FinalBlock {
	log.Println(len(peers), len(txs), epoch)
	var jobIds []uint64
	for _, tx := range txs {
		jobIds = append(jobIds, tx.Id)
	}
	for i := 1; i <= 5; i++ {
		if i%2 != 0 {
			if i == 1 {
				o.commit[epoch] = make(map[cmn.Hash][]byte)
				for _, p := range peers {
					go p.commit(epoch, txs, o.commit[epoch])
				}
				time.Sleep(5 * time.Second)
				log.Println("Commit phase completed")
			} else if i == 3 {
				nonce := uint64(0)
				for {
					blkChan := make(chan *FinalBlock, 1)
					nonce += 1
					for _, p := range peers {
						go p.proposeBlock(nonce, o.sortitionSeed, blkChan, epoch, o.results[epoch])
					}
					time.Sleep(1 * time.Second)
					if len(blkChan) == 1 {
						o.finalBlks[epoch] = <-blkChan
						break
					}
				}
				log.Println("Propose block phase compeleted")
			} else {
				// submit final block for this epoch process to all the peers
				log.Println("Step 5 reached for epcch ", epoch)
				return o.finalBlks[epoch]
			}
		} else if i == 2 {
			o.reveal[epoch] = make(map[cmn.Hash]*Reveal)
			cut := rand.Intn(len(txs))
			choice := rand.Intn(2)
			for _, p := range peers {
				if choice == 0 {
					go p.reveal(epoch, jobIds[0:cut+1], o.reveal[epoch])
				} else {
					go p.reveal(epoch, jobIds[cut+1:], o.reveal[epoch])
				}
			}
			time.Sleep(1 * time.Second)
			o.verifyCommitReveal(epoch)
			time.Sleep(8 * time.Second)
			log.Println("Reveal phase and verifying phase completed")
		}
	}
	return nil
}

func (o *Oracle) verifyCommitReveal(epoch uint64) {
	commit := o.commit[epoch]
	reveal := o.reveal[epoch]
	o.results[epoch] = make(map[uint64][][]byte)
	finalResult := o.results[epoch]
	for addr, root := range commit {
		if val, ok := reveal[addr]; ok && bytes.Equal(root, val.tree.Root()) {
			for Id, res := range val.results {
				proof, err := val.tree.GenerateProof(res)
				if err != nil {
					log.Println("Error generating proof for tx with Id ", Id)
					panic(err)
				}
				verification, err := mt.VerifyProof(res, proof, val.tree.Root())
				if err != nil {
					log.Println("Error generating the verification for res")
					panic(err)
				}
				if !verification {
					log.Println("failed to verify the proof for res")
				}
				finalResult[Id] = append(finalResult[Id], res)
			}
		} else {
			log.Println("Commit but not reveal")
		}
	}
	delete(o.commit, epoch)
	delete(o.reveal, epoch)
}
