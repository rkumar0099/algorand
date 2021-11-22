package manage

import (
	"bytes"
	"log"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/gossip"
	"github.com/rkumar0099/algorand/logs"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/mpt/kvstore"
	"github.com/rkumar0099/algorand/mpt/mpt"
	"github.com/rkumar0099/algorand/params"
	"github.com/rkumar0099/algorand/service"
	"github.com/rkumar0099/algorand/transaction"
	"google.golang.org/grpc"
)

type Manage struct {
	Id            gossip.NodeId
	node          *gossip.Node
	txs           []*msg.Transaction
	lock          *sync.Mutex
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
}

func New(peers []gossip.NodeId, peerAddresses [][]byte, lm *logs.LogManager) *Manage {
	nodeId := gossip.NewNodeId("127.0.0.1:9000")
	m := &Manage{
		Id:   nodeId,
		node: gossip.New(nodeId, "transaction"),
		lock: &sync.Mutex{},
		txs:  make([]*msg.Transaction, 0),

		stateHash:     make(map[cmn.Hash]uint64),
		epoch:         0,
		connPool:      make(map[gossip.NodeId]*grpc.ClientConn),
		proposedTxSet: make(chan *msg.ProposedTx, 1),
		count:         0,
		lm:            lm,
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
	lis, _ := net.Listen("tcp", m.Id.String())
	go m.grpcServer.Serve(lis)
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
	m.lock.Lock()
	defer m.lock.Unlock()
	m.txs = append(m.txs, tx)
	//log.Println("Transaction added")
}

func (m *Manage) Run() {
	go m.propose()
}

func (m *Manage) propose() {
	for {
		time.Sleep(30 * time.Second)
		m.epoch += 1
		go m.proposeTxs()
	}
}

func (m *Manage) proposeTxs() {
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
	//log.Println(len(m.txs))

	pt := &msg.ProposedTx{
		Epoch: m.epoch,
		Txs:   txs,
	}

	//m.proposedTxSet <- pt
	st, sh := m.executeTxSet(pt)
	m.sendContributionToPeers(pt)
	time.Sleep(10 * time.Second)
	res := m.validate(sh)
	log.Println("Response for epoch ", m.epoch, res)

	if res {
		m.sendFinalContribution(pt)
		time.Sleep(5 * time.Second)
		st.Commit()
		m.recentMPT = st
		log.Printf("[Manage] Epoch %d", m.epoch)
		log.Println(st.RootHash())
	} else {
		copy(m.txs, txs)
	}
}

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

func (m *Manage) sendContribution(conn *grpc.ClientConn, data []byte) {
	res, err := service.SendContribution(conn, data)
	if err == nil {
		m.addStateHash(cmn.BytesToHash(res.StateHash))
	}
}

func (m *Manage) addStateHash(hash cmn.Hash) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.stateHash[hash] += 1
}

func (m *Manage) validate(hash cmn.Hash) bool {
	count := m.stateHash[hash]
	ans := uint64(float64(2)/float64(3) + float64(params.UserAmount) + 0.5)
	return uint64(count) >= ans
}

func (m *Manage) sendFinalContribution(txSet *msg.ProposedTx) {
	data, _ := proto.Marshal(txSet)
	for _, conn := range m.connPool {
		go service.SendFinalContribution(conn, data)
	}
	m.lm.AddFinalTxLog(txSet.Txs)
}

func (m *Manage) LastState() []byte {
	return m.recentMPT.RootHash()
}
