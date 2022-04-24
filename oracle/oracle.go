package oracle

import (
	"bufio"
	"bytes"
	"net"
	"os"
	"strconv"

	"log"
	"sync"
	"time"

	"github.com/rkumar0099/algorand/client"
	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"
	"github.com/rkumar0099/algorand/gossip"
	"github.com/rkumar0099/algorand/message"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
)

type Oracle struct {
	id          gossip.NodeId
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
	log         string

	client *client.ClientServiceServer
}

func New(peers []gossip.NodeId) *Oracle {
	oracle := &Oracle{
		id:          gossip.NewNodeId("127.0.0.1:9005"),
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
		log:         "",
	}

	oracle.pubkey, oracle.privkey, _ = crypto.NewKeyPair()
	oracle.grpcServer = grpc.NewServer()
	oracle.server = NewServer(
		oracle.id,
		func(uint64) (*ResOPP, error) { return nil, nil },
	)
	oracle.client = client.New(oracle.id, oracle.sendReqHandler, oracle.sendResHandler)
	oracle.server.Register(oracle.grpcServer)
	oracle.client.Register(oracle.grpcServer)
	oracle.addPeers(peers)
	lis, err := net.Listen("tcp", oracle.id.String())
	if err == nil {
		log.Println("[Debug] [Oracle] Oracle listening at port 9005")
		go oracle.grpcServer.Serve(lis)
	}
	//os.RemoveAll("../oracle/blockchain")
	oracle.db, _ = leveldb.OpenFile("../oracle/blockchain", nil)
	return oracle
}

func (o *Oracle) sendReqHandler(req *client.ReqTx) (*client.ResEmpty, error) {
	res := &client.ResEmpty{}
	return res, nil
}

func (o *Oracle) sendResHandler(res *client.ResTx) (*client.ResEmpty, error) {
	info := &client.ResEmpty{}
	return info, nil
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
		if res.Weight > 0 {
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
			if err := opp.Verify(m); err == nil {
				o.peers[o.epoch] = append(o.peers[o.epoch], newOraclePeer(o.epoch))
			}
		}
	}
	log.Println("The number of oracle peers for epoch: ", o.epoch, len(o.peers[o.epoch]))
}

func (o *Oracle) GetPeers() ([]*OraclePeer, uint64) {
	return o.peers[o.epoch], o.epoch
}

func (o *Oracle) AddEWTx(tx *message.PendingRequest) {
	o.txLock.Lock()
	defer o.txLock.Unlock()
	log.Println("[Debug] [Oracle] EWT added into the pool")
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

/*

func (o *Oracle) medianizer(values []float64) float64 {
	sort.Slice(values, func(i, j int) bool {
		return values[i] <= values[j]
	})


		for _, v := range a {
			fmt.Println(v)
		}

	ind := (len(values) + 1) / 2
	return values[ind]
}
*/

func (o *Oracle) Run() {
	for {

		time.Sleep(10 * time.Second)
		if len(o.txs) == 0 {
			continue
		}

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
			log.Printf("[Oracle] Running Epoch %d, Peers %d, Txs %d\n", o.epoch, len(o.peers[o.epoch]), len(o.txs))
			o.process(o.epoch, o.peers[o.epoch], txs)
		}
	}
}

func (o *Oracle) fetch_pf(feed_type uint64) (float64, error) {
	var path string = ""
	switch feed_type {
	case 1:
		log.Println("[Debug][Oracle] ETH Pricefeed request")
		path = "../pricefeeds/eth.txt"

	case 2:
		log.Println("[Debug][Oracle] BIT Pricefeed request")
		path = "../pricefeeds/bit.txt"

	case 3:
		log.Println("[Debug][Oracle] ALGO Pricefeed request")
		path = "../pricefeeds/algo.txt"

	default:
		log.Println("[Debug][Oracle] INVALID Pricefeed request")
		return 0.0, nil
	}

	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("[Debug] [Oracle] Can't open the file")
		return 0.0, err
	}
	scanner := bufio.NewScanner(f)
	var val float64
	scanner.Scan()
	val, err = strconv.ParseFloat(scanner.Text(), 64)
	return val, err
}

func (o *Oracle) process(epoch uint64, peers []*OraclePeer, txs []*message.PendingRequest) bool {

	for _, tx := range txs {
		req := &client.ReqTx{}
		req.Deserialize(tx.Data)
		pf := &client.PriceFeed{}
		pf.Deserialize(req.Data)

		log.Printf("[Debug] [Oracle] PriceFeed: %d\n", pf.Type)

		point, _ := o.fetch_pf(pf.Type)
		log.Printf("[Debug] [Oracle] %f\n", point)
		res := &client.ResTx{}
		res.Addr = req.Addr
		res.Pubkey = req.Pubkey
		res.Data = cmn.Float2Bytes(point)

		id := gossip.NewNodeId(res.Addr)
		conn, err := id.Dial()
		if err == nil {
			_, err = client.SendResTx(conn, res)
			if err == nil {
				log.Printf("[Debug] [Oracle] Sending Response to Client: %s\n", id.String())
			}
		}
	}

	/*

		//log.Println(len(peers), len(txs), epoch)
		var responses []*msg.PendingRequestRes
		for _, tx := range txs {
			//id := tx.Id
			var values []float64

			for _, p := range peers {
				val, _ := p.processTx(tx)
				values = append(values, val)
			}

			final_point := o.medianizer(values)
			req := &client.ReqTx{}
			req.Deserialize(tx.Data)

			res := &client.ResTx{}
			res.Addr = req.Addr
			res.Pubkey = req.Pubkey
			res.Data = cmn.Float2Bytes(final_point)

			prr := &msg.PendingRequestRes{}
			prr.TxHash = tx.Hash().Bytes()
			prr.Data, _ = res.Serialize()
			responses = append(responses, prr)

		}

		blk := OracleBlock{}
		blk.Epoch = o.epoch
		blk.Res = responses
		data, _ := blk.Serialize()
		h := cmn.BytesToHash(data)
		SendBlk(data)

		// send the data to all peers and wait for the finalization

		// run the commit phase, asking all peers to commit the txs
		// run the reveal phase, asking all peers to reveal the random txs
		// compare the reveal value with commit value using merkle patricia tree roothash
		// if not match, punish the staker. If match, the value is added to the list of finalized values for this tx
		// run the medianizer phase, this will finalize the median value and punish the staker who proposed deviated value
		// compile all response with their corresponding tx response. Make a block. Propose the block
		// blk is send to blockchain, blockchain peer will verify all the information in blockchain
		// If the blk is finalized by blockchain peers, txs are confirmed, if not all txs process again in the next round
		// Once blk is finalized, corresponding responses are send to the clients
		// move to next epoch. repeat.

		o.commit(epoch, txs, peers)
		o.reveal(epoch, txs, peers)
		o.verify(epoch, txs, peers)
		o.punishers(epoch, txs, peers)
		o.propose(epoch, txs, peers)
		o.sendBlk()
		// wait for the response
		// if blk is finalized, confirm the corresponding txs

	*/

	return true

}

/*

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

func (o *Oracle) commit(epoch uint64, txs []*msg.PendingRequest, peers []*OraclePeer) {
	for _, p := range peers {

	}
}

*/
