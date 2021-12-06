package oracle

import (
	"bytes"
	"os"

	"log"
	"sync"
	"time"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"

	//"github.com/rkumar0099/algorand/logs"
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
	epoch     uint64
	lock      *sync.Mutex
	txLock    *sync.Mutex
	txs       []*message.PendingRequest
	commit    map[uint64]map[cmn.Hash][]byte
	reveal    map[uint64]map[cmn.Hash]*Reveal
	results   map[uint64]map[uint64][][]byte
	finalBlks map[uint64]*FinalBlock
	count     uint64
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
		txLock:    &sync.Mutex{},
		commit:    make(map[uint64]map[cmn.Hash][]byte),
		reveal:    make(map[uint64]map[cmn.Hash]*Reveal),
		results:   make(map[uint64]map[uint64][][]byte),
		finalBlks: make(map[uint64]*FinalBlock),
		count:     0,
	}

	oracle.pubkey, oracle.privkey, _ = crypto.NewKeyPair()
	os.RemoveAll("../oracle/blockDB")
	oracle.db, _ = leveldb.OpenFile("../oracle/blockDB", nil)
	return oracle
}

func (o *Oracle) Address() []byte {
	return o.pubkey.Bytes()
}

func (o *Oracle) AddOPP(opp *message.OraclePeerProposal) {
	o.lock.Lock()
	defer o.lock.Unlock()
	epoch := opp.Epoch

	//log.Println("Current epoch ", o.epoch)
	if _, ok := o.peers[epoch]; !ok {
		o.peers[epoch] = make([]*OraclePeer, 0)
		//log.Println("Make peer array again")
	}
	seed := o.sortitionSeed(epoch)
	role := role(params.OraclePeer, epoch, params.ORACLE)
	m := constructSeed(seed, role)
	if err := opp.Verify(m); err != nil {
		return
	}
	o.peers[epoch] = append(o.peers[epoch], newOraclePeer(epoch))
	//log.Println("The number of oracle peers ", len(o.peers[epoch]))
}

func (o *Oracle) AddBlk(blk *message.Block) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if o.round < blk.Round {
		data, _ := blk.Serialize()
		o.db.Put(cmn.Uint2Bytes(blk.Round), data, nil)
		o.round = blk.Round
	}
}

func (o *Oracle) GetPeers() ([]*OraclePeer, uint64) {
	return o.peers[o.epoch], o.epoch
}

func (o *Oracle) GetBlkByRound(round uint64) *message.Block {
	blk := &message.Block{}
	data, _ := o.db.Get(cmn.Uint2Bytes(round), nil)
	err := blk.Deserialize(data)
	if err != nil {
		return nil
	}
	return blk
}

func (o *Oracle) AddEWTx(tx *message.PendingRequest) {
	o.txLock.Lock()
	defer o.txLock.Unlock()
	o.count += 1
	tx.Id = o.count
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

func (o *Oracle) Run() {
	for {

		time.Sleep(10 * time.Second)
		o.epoch += 1

		var (
			tx  *message.PendingRequest
			txs []*message.PendingRequest
		)

		for len(txs) < 10 && len(o.txs) > 0 {
			tx, o.txs = o.txs[0], o.txs[1:]
			txs = append(txs, tx)
		}
		log.Println("Oracle run: ", o.epoch, len(txs), len(o.peers[o.epoch]))
		finalBlk := o.process(o.epoch, o.peers[o.epoch], txs)
		o.finalBlks[o.epoch] = finalBlk
		time.Sleep(2 * time.Minute)

	}
}

func (o *Oracle) process(epoch uint64, peers []*OraclePeer, txs []*message.PendingRequest) *FinalBlock {
	//log.Println(len(peers), len(txs), epoch)
	var jobIds []uint64
	for _, tx := range txs {
		jobIds = append(jobIds, tx.Id)
	}

	for i := 1; i <= 5; i++ {

		if i == 1 {
			// commit phase
			o.commit[epoch] = make(map[cmn.Hash][]byte)
			for _, p := range peers {
				go p.commit(epoch, txs, o.commit[epoch])
			}
			time.Sleep(5 * time.Second)
			log.Println("Commit phase completed")
		} else if i == 2 {
			// reveal phase
			o.reveal[epoch] = make(map[cmn.Hash]*Reveal)

			for _, p := range peers {
				go p.reveal(epoch, jobIds, o.reveal[epoch]) // oracle peers reveal about all jobs for current epoch
			}

			time.Sleep(1 * time.Second)
			log.Println("Reveal phase completed")
			o.verifyCommitReveal(epoch)
			log.Println("Verifying phase completed")

		} else if i == 3 {
			// select a proposer using nonce different nonce, we can also consider max priority approach
			// of algorand to choose a final proposal from a list of different oracle block proposers

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
					delete(o.results, epoch)
					break
				}
			}
			log.Println("Propose block phase compeleted")
			log.Println("Final block is ", o.finalBlks[epoch])
		} else if i == 4 {
			// dispute phase, out of scope for this term
			log.Println("Dispute phase")
		} else {
			// submit the block to all blockchain peers using service
			//blk := o.finalBlks[epoch]
			//o.lm.AddOracleBlk(blk)
			delete(o.finalBlks, epoch)
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
				dataCommit := bytes.Join([][]byte{
					res,
					cmn.Uint2Bytes(Id),
				}, nil)
				proof, err := val.tree.GenerateProof(dataCommit)
				if err != nil {
					log.Println("Error generating proof for tx with Id ", Id)
					panic(err)
				}
				verification, err := mt.VerifyProof(dataCommit, proof, val.tree.Root())
				if err != nil {
					log.Println("Error generating the verification for tx with Id ", Id)
					panic(err)
				}
				if !verification {
					log.Println("failed to verify the proof for tx with Id ", Id)
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
