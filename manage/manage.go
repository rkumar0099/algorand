package manage

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/gossip"
	"github.com/rkumar0099/algorand/logs"
	"github.com/rkumar0099/algorand/message"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/mpt/kvstore"
	"github.com/rkumar0099/algorand/mpt/mpt"
	"github.com/rkumar0099/algorand/service"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
)

type Manage struct {
	Id            gossip.NodeId
	node          *gossip.Node
	txs           []*msg.Transaction
	txLock        *sync.Mutex
	shLock        *sync.Mutex
	votesLock     *sync.Mutex
	epoch         uint64
	stateHash     map[cmn.Hash]uint64
	connPool      map[gossip.NodeId]*grpc.ClientConn
	proposedTxSet chan *msg.ProposedTx
	ServiceServer *service.Server
	grpcServer    *grpc.Server
	count         int
	storage       kvstore.KVStore
	recentMPT     *mpt.Trie
	lm            *logs.LogManager
	blkLock       *sync.Mutex
	db            *leveldb.DB
	lastHash      []byte
	lastRound     uint64
	contributions map[common.Hash]*msg.ProposedTx
	votes         map[common.Hash]uint64
	confirm       uint64
}

func New(peers []gossip.NodeId, peerAddresses [][]byte, lm *logs.LogManager) *Manage {
	nodeId := gossip.NewNodeId("127.0.0.1:9000")
	m := &Manage{
		Id:     nodeId,
		node:   gossip.New(nodeId, "transaction"),
		txLock: &sync.Mutex{},
		shLock: &sync.Mutex{},
		txs:    make([]*msg.Transaction, 0),

		stateHash:     make(map[cmn.Hash]uint64),
		epoch:         0,
		connPool:      make(map[gossip.NodeId]*grpc.ClientConn),
		proposedTxSet: make(chan *msg.ProposedTx, 1),
		count:         0,
		lm:            lm,
		blkLock:       &sync.Mutex{},
		contributions: make(map[common.Hash]*msg.ProposedTx),
		votes:         make(map[cmn.Hash]uint64),
		votesLock:     &sync.Mutex{},
		confirm:       0,
		lastRound:     0,
	}
	m.storage = kvstore.NewMemKVStore() // manage storage to execute Tx set and generate Hash
	m.grpcServer = grpc.NewServer()
	m.ServiceServer = service.NewServer(
		nodeId,
		func([]byte) ([]byte, error) { return nil, nil },
		func([]byte) ([]byte, error) { return nil, nil },
		func([]byte) error { return nil },
	)
	m.node.Register(m.grpcServer)
	m.ServiceServer.Register(m.grpcServer)
	m.addPeers(peers)
	m.initializeMPT(peerAddresses)

	//lis, _ := net.Listen("tcp", m.Id.String())
	//go m.grpcServer.Serve(lis)

	os.Remove("../logs/manage.txt")
	os.RemoveAll("../database/blockchain")
	m.db, _ = leveldb.OpenFile("../database/blockchain", nil)
	return m
}

func (m *Manage) addPeers(peers []gossip.NodeId) {
	for _, id := range peers {
		conn, err := id.Dial()
		if err == nil {
			m.connPool[id] = conn
		}
	}
}

func (m *Manage) initializeMPT(addr [][]byte) {
	st := mpt.New(nil, m.storage)
	for _, val := range addr {
		stateAddr := bytes.Join([][]byte{
			val,
			[]byte("value"),
		}, nil)
		st.Put(stateAddr, cmn.Uint2Bytes(0))
	}
	m.recentMPT = st
	st.Commit()
}

func (m *Manage) AddTransaction(tx *msg.Transaction) {
	m.txLock.Lock()
	defer m.txLock.Unlock()
	if err := tx.VerifySign(); err != nil {
		log.Printf("[algorand] Received invalid transaction: %s", err.Error())
		return
	}
	m.txs = append(m.txs, tx)
	m.lm.AddTxLog(tx.Hash())
}

func (m *Manage) Run() {
	go m.propose()
}

func (m *Manage) propose() {
	for {
		time.Sleep(5 * time.Second)
		if len(m.txs) >= 20 {
			m.process()
		}
	}
}

func (m *Manage) process() {
	var (
		txs []*msg.Transaction
		tx  *msg.Transaction
	)
	for len(m.txs) > 0 {
		tx, m.txs = m.txs[0], m.txs[1:]
		txs = append(txs, tx)
		if len(txs) >= 20 {
			break
		}
	}

	pt := &msg.ProposedTx{
		Epoch: m.epoch,
		Txs:   txs,
	}

	h := pt.Hash()
	m.contributions[h] = pt
	m.sendFinalContribution(pt)
	//m.confirm += 1
	time.Sleep(10 * time.Second)
	if m.votes[h] >= 33 {
		// confirm contribution
		m.confirm += 1
	} else {
		// put the contribution back into the txs
		pt = m.contributions[h]
		m.txs = append(m.txs, pt.Txs...)
	}

	delete(m.contributions, h)
	delete(m.votes, h)

	// proceed to next contribution
}

func (m *Manage) GetConfirmedContributions() uint64 {
	return m.confirm
}

/*
func (m *Manage) proposeTxs() {
	var (
		txs []*msg.Transaction
		tx  *msg.Transaction
	)
	for len(m.txs) > 0 {
		tx, m.txs = m.txs[0], m.txs[1:]
		txs = append(txs, tx)
		if len(txs) > 40 {
			break
		}
	}

	pt := &msg.ProposedTx{
		Epoch: m.epoch,
		Txs:   txs,
	}

	m.sendFinalContribution(pt)
	// wait for the contribution to finish
	//log.Println(len(m.txs))

	/*
		//m.proposedTxSet <- pt
		//st, sh := m.executeTxSet(pt)
		//m.sendContributionToPeers(pt)
		//time.Sleep(5 * time.Second)
		//res, _ := m.validate(sh)
		//go m.writeLog(sh, m.stateHash[sh], res, c)

		if res {
			m.sendFinalContribution(pt)
			// wait for the block to finalized, and remove the transactions from queue iff only
			// they are present in the blk

			time.Sleep(5 * time.Second)
			st.Commit()
			m.recentMPT = st
		} else {
			copy(m.txs, txs)
		}

}
*/
/*
func (m *Manage) executeTxSet(txSet *msg.ProposedTx) (*mpt.Trie, cmn.Hash) {
	txs := txSet.Txs
	st := m.recentMPT
	for _, tx := range txs {
		switch tx.Type {
		case transaction.TOPUP:
			transaction.Topup(st, tx, cmn.Bytes2Uint(tx.Data))
		case transaction.TRNASFER:
			transaction.Transfer(st, tx, cmn.Bytes2Uint(tx.Data))
		default:
			log.Printf("Received invalid transaction")
		}
	}
	hash := cmn.BytesToHash(st.RootHash())
	m.stateHash[hash] = 0
	return st, hash
}

func (m *Manage) sendContributionToPeers(txSet *msg.ProposedTx) {
	data, _ := proto.Marshal(txSet)
	for _, conn := range m.connPool {
		go m.sendContribution(conn, data)
		//go service.SendContribution(conn, data)
	}
}


*/

func (m *Manage) GetByRound(round uint64) *msg.Block {
	blk := &message.Block{}
	data, _ := m.db.Get(cmn.Uint2Bytes(round), nil)
	err := blk.Deserialize(data)
	if err != nil {
		return nil
	}
	return blk
}

func (m *Manage) sendFinalContribution(txSet *msg.ProposedTx) {
	data, _ := proto.Marshal(txSet)
	for _, conn := range m.connPool {
		go service.SendFinalContribution(conn, data)
	}
	//m.lm.AddFinalTxLog(txSet.Txs)
}

func (m *Manage) AddVote(contribution common.Hash) {
	m.votesLock.Lock()
	defer m.votesLock.Unlock()
	m.votes[contribution] += 1
	log.Printf("Confirmed votes for %s: %d\n", contribution.String(), m.votes[contribution])
}

func (m *Manage) LastState() []byte {
	return m.recentMPT.RootHash()
}

func (m *Manage) AddBlk(hash cmn.Hash, data []byte) {
	m.blkLock.Lock()
	defer m.blkLock.Unlock()

	// find blk with hash key or round key
	if !bytes.Equal(m.lastHash, hash.Bytes()) {
		m.db.Put(hash.Bytes(), data, nil)
		blk := &message.Block{}
		blk.Deserialize(data)
		m.db.Put(cmn.Uint2Bytes(blk.Round), data, nil)
		m.lastHash = hash.Bytes()
		m.lastRound = blk.Round
	}
}

func (m *Manage) GetBlkByHash(hash cmn.Hash) *msg.Block {
	blk := &message.Block{}
	data, err := m.db.Get(hash.Bytes(), nil)
	if err != nil {
		s := fmt.Sprintf("[Error] Can't find block by hash, %s\n", err.Error())
		fmt.Print(s)
		return nil
	}
	err = blk.Deserialize(data)
	if err != nil {
		s := fmt.Sprintf("[Error] Can't Deserialize blk data found by hash, %s\n", err.Error())
		fmt.Print(s)
		return nil
	}
	return blk
}

func (m *Manage) GetBlkByRound(round uint64) *msg.Block {
	blk := &message.Block{}
	data, err := m.db.Get(cmn.Uint2Bytes(round), nil)
	if err != nil {
		s := fmt.Sprintf("[Error] Can't find block by Round, %s\n", err.Error())
		fmt.Print(s)
		return nil
	}
	err = blk.Deserialize(data)
	if err != nil {
		s := fmt.Sprintf("[Error] Can't Deserialize blk data found by Round, %s\n", err.Error())
		fmt.Print(s)
		return nil
	}
	return blk
}

func (m *Manage) writeLog(hash cmn.Hash, count uint64, res bool, c uint64) {
	f, err := os.OpenFile("../logs/manage.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	latestLog := ""
	if res {
		latestLog += "True. "
	} else {
		latestLog += "False. "
	}
	latestLog += "Count is " + strconv.Itoa(int(count)) + "Res is " + strconv.Itoa(int(c)) + " for hash " + hash.Hex() + "\n"
	f.WriteString(latestLog)
	f.Close()
}
