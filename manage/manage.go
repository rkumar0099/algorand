package manage

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rkumar0099/algorand/client"
	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/gossip"
	"github.com/rkumar0099/algorand/message"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/service"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
)

type Manage struct {
	Id     gossip.NodeId
	node   *gossip.Node
	txs    []*msg.Transaction
	txLock *sync.Mutex

	lastHash      []byte
	epoch         uint64
	connPool      map[gossip.NodeId]*grpc.ClientConn
	proposedTxSet chan *msg.ProposedTx
	ServiceServer *service.Server
	grpcServer    *grpc.Server
	db            *leveldb.DB
	lastRound     uint64
	resLock       *sync.Mutex
	response      map[common.Hash][]*msg.TxRes
	lastTxHash    []byte
	client        *client.ClientServiceServer
	blkLock       *sync.Mutex
	log           string
	recHash       []byte
}

func New(peers []gossip.NodeId, peerAddresses [][]byte) *Manage {
	nodeId := gossip.NewNodeId("127.0.0.1:9000")
	m := &Manage{
		Id:     nodeId,
		node:   gossip.New(nodeId, "transaction"),
		txLock: &sync.Mutex{},
		txs:    make([]*msg.Transaction, 0),

		epoch:         0,
		connPool:      make(map[gossip.NodeId]*grpc.ClientConn),
		proposedTxSet: make(chan *msg.ProposedTx, 1),

		blkLock: &sync.Mutex{},

		lastRound:  0,
		resLock:    &sync.Mutex{},
		response:   make(map[common.Hash][]*msg.TxRes),
		lastTxHash: []byte(""),
		log:        "",
		recHash:    []byte(""),
	}
	//m.lm = lm
	m.client = client.New(m.Id, m.sendReqHandler, m.sendResHandler)
	m.grpcServer = grpc.NewServer()
	m.ServiceServer = service.NewServer(
		nodeId,
		func([]byte) ([]byte, error) { return nil, nil },
		func([]byte) error { return nil },
	)
	m.node.Register(m.grpcServer)
	m.ServiceServer.Register(m.grpcServer)
	m.client.Register(m.grpcServer)
	m.addPeers(peers)

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

func (m *Manage) sendReqHandler(req *client.ReqTx) (*client.ResEmpty, error) {
	res := &client.ResEmpty{}
	return res, nil
}

func (m *Manage) sendResHandler(res *client.ResTx) (*client.ResEmpty, error) {
	info := &client.ResEmpty{}
	return info, nil
}

func (m *Manage) AddTransaction(tx *msg.Transaction) {
	m.txLock.Lock()
	defer m.txLock.Unlock()
	h := tx.Hash()
	if !bytes.Equal(h.Bytes(), m.lastTxHash) {
		m.lastTxHash = h.Bytes()
		m.txs = append(m.txs, tx)
		s := "[Debug] [Manage] Tx Added to Pool\n"
		log.Print(s)
		m.log += s
	}
}

func (m *Manage) Run() {
	go m.propose()
}

func (m *Manage) propose() {
	for {
		time.Sleep(10 * time.Second)
		if len(m.txs) > 0 {
			m.epoch += 1
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
	go m.sendFinalContribution(pt)

	for {
		time.Sleep(5 * time.Second)
		if bytes.Equal(h.Bytes(), m.recHash) {
			break
		}
	}

	// proceed to next contribution
}

func (m *Manage) sendResponse(responses []*msg.TxRes) {
	for _, res := range responses {
		r := &client.ResTx{}
		r.Deserialize(res.Data)
		id := gossip.NewNodeId("127.0.0.1:9020")
		conn, err := id.Dial()
		if err == nil {
			_, err = client.SendResTx(conn, r)
			if err == nil {
				log.Printf("[Debug] [Manage] Sending Response to Client: %s\n", id.String())
			}
		}
	}
}

func (m *Manage) AddRes(hash cmn.Hash, res []*msg.TxRes) {
	m.resLock.Lock()
	defer m.resLock.Unlock()
	h := hash.Bytes()
	if !bytes.Equal(h, m.recHash) {
		m.recHash = h
		m.sendResponse(res)
	}
	log.Println("[Debug] [Manage] Got TXs Response from Peer")
}

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
