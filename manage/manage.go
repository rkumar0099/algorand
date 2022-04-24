package manage

import (
	"bytes"
	"fmt"
	"log"
	"net"
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
	finalized     map[cmn.Hash]bool
	added         map[cmn.Hash]bool
	off           map[cmn.Hash]bool
	offLock       *sync.Mutex
}

func New(peers []gossip.NodeId, peerAddresses [][]byte) *Manage {
	nodeId := gossip.NewNodeId("127.0.0.1:9000")
	m := &Manage{
		Id:            nodeId,
		node:          gossip.New(nodeId, "transaction"),
		txLock:        &sync.Mutex{},
		txs:           make([]*msg.Transaction, 0),
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
		finalized:  make(map[cmn.Hash]bool),
		added:      make(map[cmn.Hash]bool),
		off:        make(map[cmn.Hash]bool),
		offLock:    &sync.Mutex{},
	}
	//m.lm = lm
	m.grpcServer = grpc.NewServer()
	m.client = client.New(m.Id, m.sendReqHandler, m.sendResHandler)

	m.ServiceServer = service.NewServer(
		nodeId,
		func([]byte) ([]byte, error) { return nil, nil },
		func([]byte) error { return nil },
	)
	m.node.Register(m.grpcServer)
	m.ServiceServer.Register(m.grpcServer)
	m.client.Register(m.grpcServer)
	m.addPeers(peers)
	lis, err := net.Listen("tcp", m.Id.String())
	if err == nil {
		log.Println("[Debug] [Manage] Manage listening at port 9000")
		go m.grpcServer.Serve(lis)
	}
	os.Remove("../logs/manage.txt")
	//os.RemoveAll("../database/blockchain")
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
	if _, ok := m.added[h]; !ok {
		// add the tx
		m.added[h] = true
		m.txs = append(m.txs, tx)
		m.txLog(tx)
	}
}

func (m *Manage) Run() {
	go m.propose()
}

func (m *Manage) propose() {
	for {
		time.Sleep(1 * time.Second)
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
	pt.PtHash = h.Bytes()
	s := fmt.Sprintf("[Manage] Contribution proposed. Epoch: %d, Hash: %s\n", m.epoch, h.String())
	m.log += s
	m.sendFinalContribution(pt)

	for {
		time.Sleep(1 * time.Second)
		if _, ok := m.finalized[h]; ok {
			break
		}
	}
	// proceed to next contribution
}

func (m *Manage) sendResponse(responses []*msg.TxRes) {
	for _, response := range responses {

		if response.SendRes {
			// send the response to Client
			res := &client.ResTx{}
			res.Deserialize(response.Res)
			id := gossip.NewNodeId(res.Addr)
			conn, err := id.Dial()
			if err == nil {
				_, err = client.SendResTx(conn, res)
				if err == nil {
					log.Printf("[Debug] [Manage] Sending Response to Client: %s\n", id.String())
				}
			}
		}

		if response.StoreData {
			err := m.db.Put(response.DataAddr, response.Data, nil)
			if err != nil {
				log.Println("[Debug] Unable to store response Data")
			}
		}
	}
}

func (m *Manage) AddRes(hash cmn.Hash, res []*msg.TxRes) {
	m.resLock.Lock()
	defer m.resLock.Unlock()
	if _, ok := m.finalized[hash]; !ok {
		// finalize the set
		s := fmt.Sprintf("Received response for contribution: %s, Len: %d\n", hash.String(), len(res))
		m.log += s
		go m.sendResponse(res)
		m.finalized[hash] = true
	}
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

func (m *Manage) GetClient(credentials []byte) ([]byte, error) {
	data, err := m.db.Get(credentials, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
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

func (m *Manage) GetLog() string {
	return m.log
}

func (m *Manage) txLog(tx *msg.Transaction) {
	s := fmt.Sprintf("[Manage] Tx type: %d\n", tx.Type)
	m.log += s
}

func (m *Manage) AddOffTopup(hash cmn.Hash, tx *msg.Transaction) {
	m.offLock.Lock()
	defer m.offLock.Unlock()
	if _, ok := m.off[hash]; !ok {
		m.log += "Off top up tx added"
		m.AddTransaction(tx)
		m.off[hash] = true
	}
}
