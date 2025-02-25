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
	"github.com/rkumar0099/algorand/gossip"
	"github.com/rkumar0099/algorand/message"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
	"github.com/syndtr/goleveldb/leveldb"
	mt "github.com/wealdtech/go-merkletree"
	"google.golang.org/grpc"
)

type Oracle struct {
	pubkey      *crypto.PublicKey
	privkey     *crypto.PrivateKey
	peers       map[uint64][]*OraclePeer
	round       uint64
	epoch       uint64
	lock        *sync.Mutex
	txLock      *sync.Mutex
	txs         []*message.PendingRequest
	commit      map[uint64]map[cmn.Hash][]byte
	reveal      map[uint64]map[cmn.Hash]*Reveal
	results     map[uint64]map[uint64][][]byte
	finalBlks   map[uint64]*FinalBlock
	count       uint64
	lastRound   uint64
	db          *leveldb.DB
	running     bool
	confirmedId int
	server      *OracleServiceServer
	grpcServer  *grpc.Server
	connPool    map[gossip.NodeId]*grpc.ClientConn
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

func New(peers []gossip.NodeId) *Oracle {
	oracle := &Oracle{
		peers:       make(map[uint64][]*OraclePeer),
		lock:        &sync.Mutex{},
		round:       0,
		epoch:       0,
		txLock:      &sync.Mutex{},
		commit:      make(map[uint64]map[cmn.Hash][]byte),
		reveal:      make(map[uint64]map[cmn.Hash]*Reveal),
		results:     make(map[uint64]map[uint64][][]byte),
		finalBlks:   make(map[uint64]*FinalBlock),
		count:       0,
		running:     false,
		confirmedId: 0,
		connPool:    make(map[gossip.NodeId]*grpc.ClientConn),
	}

	oracle.pubkey, oracle.privkey, _ = crypto.NewKeyPair()
	oracle.grpcServer = grpc.NewServer()
	oracle.server = NewServer(
		gossip.NewNodeId("127.0.0.1:9001"),
		func(uint64) (*ResOPP, error) { return nil, nil },
	)
	oracle.server.Register(oracle.grpcServer)
	oracle.addPeers(peers)
	os.RemoveAll("../oracle/blockchain")
	oracle.db, _ = leveldb.OpenFile("../oracle/blockchain", nil)
	return oracle
}

func (o *Oracle) addPeers(peers []gossip.NodeId) {
	for _, id := range peers {
		conn, err := id.Dial()
		if err == nil {
			o.connPool[id] = conn
		}
	}
}

func (o *Oracle) Address() []byte {
	return o.pubkey.Bytes()
}

func (o *Oracle) AddBlock(blk *msg.Block) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if blk.Round > o.lastRound {
		data, _ := blk.Serialize()
		o.db.Put(cmn.Uint2Bytes(blk.Round), data, nil)
		o.lastRound = blk.Round
	}
}

func (o *Oracle) preparePool() {
	o.epoch += 1
	for _, conn := range o.connPool {
		res, _ := SendOPP(conn, o.epoch)
		if res != nil {
			seed := o.sortitionSeed(1)
			role := role(params.OraclePeer, o.epoch, params.ORACLE)
			m := constructSeed(seed, role)
			opp := &message.OraclePeerProposal{
				Proof:  res.Proof,
				Vrf:    res.VRF,
				Pubkey: res.Pubkey,
				Epoch:  o.epoch,
				Weight: res.Weight,
			}
			if err := opp.Verify(m); err != nil {
				o.peers[o.epoch] = append(o.peers[o.epoch], newOraclePeer(o.epoch))
			}
		}
	}
	log.Println("The number of oracle peers for epoch: ", o.epoch, len(o.peers[o.epoch]))
	// send req to peer to send the result of proposal, see if it's selected
}

func (o *Oracle) GetPeers() ([]*OraclePeer, uint64) {
	return o.peers[o.epoch], o.epoch
}

func (o *Oracle) AddEWTx(tx *message.PendingRequest) {
	o.txLock.Lock()
	defer o.txLock.Unlock()
	log.Println("EW transaction proposed")
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
	// get block by manage
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

func (o *Oracle) completeEpoch() {
	// once oracle peers are ready, ask all of them to fetch jobs and store the result in a merkle
	// patricia tree
}

func (o *Oracle) Run() {
	for {

		time.Sleep(10 * time.Second)
		o.preparePool()

		var (
			tx  *message.PendingRequest
			txs []*message.PendingRequest
		)

		for len(txs) < 20 && len(o.txs) > 0 {
			tx, o.txs = o.txs[0], o.txs[1:]
			txs = append(txs, tx)
		}

		if len(txs) > 0 {
			log.Println("Oracle run: ", o.epoch, len(txs), len(o.peers[o.epoch]))
			o.running = true
			//go o.process(o.epoch, o.peers[o.epoch], txs)
		}
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
			log.Println("Commit phase started")
			o.commit[epoch] = make(map[cmn.Hash][]byte)
			for _, p := range peers {
				p.commit(epoch, txs, o.commit[epoch])
			}
			//time.Sleep(5 * time.Second)
			log.Println("Commit phase completed")
		} else if i == 2 {
			// reveal phase
			o.reveal[epoch] = make(map[cmn.Hash]*Reveal)

			for _, p := range peers {
				p.reveal(epoch, jobIds, o.reveal[epoch]) // oracle peers reveal about all jobs for current epoch
			}

			//time.Sleep(1 * time.Second)
			log.Println("Reveal phase completed")
			o.verifyCommitReveal(epoch)
			log.Println("Verifying phase completed")

		} else if i == 3 {
			// select a proposer using nonce different nonce, we can also consider max priority approach
			// of algorand to choose a final proposal from a list of different oracle block proposers
			log.Println("Proposing block")
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
			log.Println("Disputing the block")
			log.Println("Dispute phase completed")
		} else {
			// submit the block to all blockchain peers using service
			blk := o.finalBlks[epoch]
			log.Println("Confirmed EW transactions ", len(blk.Results))
			//o.lm.AddOracleBlk(blk)
			//o.writeBlk(blk)
			delete(o.finalBlks, epoch)
			o.running = false
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

/*
func (o *Oracle) writeBlk(blk *FinalBlock) {
	log.Println(blk)
	f, err := os.OpenFile("../logs/oracle.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("[Error] Can't open file to write oracle blk")
		return
	}
	log := "Oracle block for epoch " + strconv.Itoa(int(blk.Epoch)) + " Finalized\n"
	f.WriteString(log)
	f.Close()
	f, err = os.OpenFile("../logs/externalData.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	results := blk.Results
	log = ""
	for _, data := range results {
		o.confirmedId += 1
		res := &Response{}
		proto.Unmarshal(data, res)
		log += "Job with Id " + strconv.Itoa(o.confirmedId) + " access data point: " + string(res.Data) + "\n"
		//f.WriteString(log)
		//f.Write(res.Data)
		//f.WriteString("\n")
	}
	f.WriteString(log)
	f.Close()

}
*/
