package peer

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/rkumar0099/algorand/logs"
	"github.com/rkumar0099/algorand/manage"
	"github.com/rkumar0099/algorand/mpt/kvstore"
	"github.com/rkumar0099/algorand/oracle"
	"google.golang.org/grpc"

	"github.com/golang/protobuf/proto"

	"github.com/rkumar0099/algorand/blockchain"
	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"
	msg "github.com/rkumar0099/algorand/message"

	"github.com/rkumar0099/algorand/params"
	"github.com/rkumar0099/algorand/pool"
	"github.com/rkumar0099/algorand/service"

	"github.com/rkumar0099/algorand/gossip"
)

type Peer struct {
	maliciousType int // true: honest user; false: malicious user.

	Id      gossip.NodeId      // peer Id
	privkey *crypto.PrivateKey // peer private key
	pubkey  *crypto.PublicKey  // peer public key

	permanentTxStorage kvstore.KVStore // Finalized txs
	tempTxStorage      kvstore.KVStore // Proposed txs

	chain       *blockchain.Blockchain // blockchain
	quitCh      chan struct{}
	hangForever chan struct{}

	node          *gossip.Node // node to communicate msgs over network
	grpcServer    *grpc.Server
	msgAgent      *msg.MsgAgent
	ServiceServer *service.Server

	votePool     *pool.VotePool
	proposalPool *pool.ProposalPool

	manage *manage.Manage   // manage package to help in txs, logs, and oracle
	oracle *oracle.Oracle   // oracle
	lm     *logs.LogManager // log manager

	txEpoch uint64

	finalContributions       chan *msg.ProposedTx
	oracleFinalContributions chan []byte //make struct for oracle response

	lastState  []byte
	startState []byte

	oracleEpoch uint64
}

//const memPoolCap = 512
const votePoolCap = params.ExpectedCommitteeMembers * 2
const cleanTimeout = time.Minute * 2

func New(addr string, maliciousType int) *Peer {
	peer := &Peer{
		Id:                       gossip.NewNodeId(addr),
		quitCh:                   make(chan struct{}),
		hangForever:              make(chan struct{}),
		grpcServer:               grpc.NewServer(),
		maliciousType:            params.Honest,
		txEpoch:                  0,
		oracleEpoch:              0,
		finalContributions:       make(chan *msg.ProposedTx, 10),
		oracleFinalContributions: make(chan []byte, 10),
	}
	// gossip node
	peer.node = gossip.New(peer.Id, "algorand")

	// pubkey and privkey
	rand.Seed(time.Now().UnixNano())
	peer.pubkey, peer.privkey, _ = crypto.NewKeyPair()

	peer.permanentTxStorage = kvstore.NewMemKVStore()
	peer.tempTxStorage = kvstore.NewMemKVStore()

	peer.chain = blockchain.NewBlockchain()

	// proposal pool
	peer.proposalPool = pool.NewProposalPool(cleanTimeout, peer.proposalVerifier)

	// vote pool

	peer.votePool = pool.NewVotePool(votePoolCap, peer.voteVerifier, cleanTimeout)

	// register the same grpc server for service and gossip
	peer.ServiceServer = service.NewServer(peer.Id, peer.getDataByHashHandler, peer.handleContribution, peer.handleFinalContribution)
	peer.node.Register(peer.grpcServer)
	peer.ServiceServer.Register(peer.grpcServer)

	// msg agent and register handlers
	peer.msgAgent = msg.NewMsgAgent(peer.node.GetMsgChan())
	peer.msgAgent.Register(msg.VOTE, peer.handleVote)
	peer.msgAgent.Register(msg.BLOCK, peer.handleBlock)
	peer.msgAgent.Register(msg.BLOCK_PROPOSAL, peer.handleBlockProposal)
	peer.msgAgent.Register(msg.FORK_PROPOSAL, peer.handleForkProposal)
	return peer
}

func (p *Peer) Start(neighbors []gossip.NodeId, addr [][]byte) error {
	lis, err := net.Listen("tcp", p.Id.String())
	if err != nil {
		log.Printf("[algorand] [%s] Cannot start peer: %s", p.Id.String(), err.Error())
		return errors.New(fmt.Sprintf("Cannot start algorand peer: %s", err.Error()))
	}
	go p.grpcServer.Serve(lis)
	p.node.Join(neighbors)
	go p.msgAgent.Handle()

	p.initialize(addr, p.permanentTxStorage)
	p.initialize(addr, p.tempTxStorage)
	return nil
}

func (p *Peer) AddManage(m *manage.Manage, lm *logs.LogManager, oracle *oracle.Oracle) {
	p.manage = m
	p.oracle = oracle
	p.lm = lm
	p.chain.SetManage(m)
	blk := p.chain.GetByRound(0)
	//go p.lm.AddFinalBlk(blk.Hash(), 0)
	go p.oracle.AddBlk(blk) // add initial algorand block to oracle
}

func (p *Peer) StartServices(bootNode []gossip.NodeId) error {
	lis, err := net.Listen("tcp", p.Id.String())
	if err != nil {
		log.Printf("[algorand] [%s] Cannot start peer: %s", p.Id.String(), err.Error())
		return errors.New(fmt.Sprintf("Cannot start algorand peer: %s", err.Error()))
	}
	go p.grpcServer.Serve(lis)
	p.node.Join(bootNode)
	p.msgAgent.Handle()
	return nil
}

func (p *Peer) Stop() {
	close(p.quitCh)
	close(p.hangForever)
}

func (p *Peer) GetGrpcServer() *grpc.Server {
	return p.grpcServer
}

func (p *Peer) getDataByHashHandler(hash []byte) ([]byte, error) {
	// to do
	return nil, nil
}

// round returns the latest round number.
func (p *Peer) round() uint64 {
	return p.lastBlock().Round
}

func (p *Peer) lastBlock() *msg.Block {
	return p.chain.Last
}

// weight returns the weight of the given address.
func (p *Peer) weight(address cmn.Address) uint64 {
	return params.TokenPerUser
}

// tokenOwn returns the token amount (weight) owned by self node.
func (p *Peer) tokenOwn() uint64 {
	return p.weight(p.Address())
}

func (p *Peer) emptyBlock(round uint64, prevHash cmn.Hash) *msg.Block {
	prevBlk, err := p.getBlock(prevHash)
	if err != nil {
		log.Printf("node %d hang forever because cannot get previous block", p.Id)
		<-p.hangForever
	}
	return &msg.Block{
		Round:      round,
		ParentHash: prevHash.Bytes(),
		StateHash:  prevBlk.StateHash,
	}
}

func (p *Peer) Address() cmn.Address {
	return cmn.BytesToAddress(p.pubkey.Bytes())
}

// forkLoop periodically resolves fork
func (p *Peer) forkLoop() {
	forkInterval := time.NewTicker(params.ForkResolveInterval)

	for {
		select {
		case <-p.quitCh:
			return
		case <-forkInterval.C:
			p.processForkResolve()
		}
	}
}

func (p *Peer) gossip(typ int, data []byte) {
	message := &msg.Msg{
		PID:  p.Id.String(),
		Type: int32(typ),
		Data: data,
	}
	switch typ {
	case msg.BLOCK:
		log.Printf("[debug] [%s] gossip block", p.Id.String())
	case msg.BLOCK_PROPOSAL:
		log.Printf("[debug] [%s] gossip block proposal", p.Id.String())
	case msg.VOTE:
		log.Printf("[debug] [%s] gossip vote", p.Id.String())
	}
	msgBytes, err := proto.Marshal(message)
	if err != nil {
		log.Printf("[alogrand] [%s] cannot gossip message: %s", err.Error())
	} else {
		p.node.Gossip(msgBytes)
	}
	p.lm.AddProcessLog(message)
}

func (p *Peer) emptyHash(round uint64, prev common.Hash) common.Hash {
	return p.emptyBlock(round, prev).Hash()
}

func (p *Peer) handleBlock(data []byte) {
	blk := &msg.Block{}
	if err := blk.Deserialize(data); err != nil {
		//log.Printf("[algorand] [%s] Received invalid block message: %s", p.Id.String(), err.Error())
		return
	}
	p.chain.CacheBlock(blk)
}

func (p *Peer) getBlock(hash common.Hash) (*msg.Block, error) {
	// find  locally
	blk, err := p.chain.Get(hash)
	if err != nil {
		log.Printf("[algorand] [%s] cannot get block %s locally: %s, try to find from other peers", p.Id.String(), hash.Hex(), err.Error())
	} else {
		return blk, nil
	}

	// find from other peers
	neighborList := p.node.GetNeighborList()
	neighbors := neighborList.GetNeighborsId()
	chanBlk := make(chan *msg.Block, 1)
	blkReqDone := false
	// make query in batch of 10

	go func() {
		for i := 0; i < len(neighbors) && !blkReqDone; i += 10 {
			for j := i; j < i+10 && j < len(neighbors) && !blkReqDone; j++ {
				go func(nodeId gossip.NodeId) {
					conn, err := neighborList.GetConn(nodeId)
					if err != nil {
						log.Printf("[alogrand] [%s] cannot get block from [%s]: %s", p.Id.String(), nodeId.String(), err.Error())
						return
					}
					blk, err := service.GetBlock(conn, hash.Bytes())
					if err != nil {
						log.Printf("[alogrand] [%s] cannot get block from [%s]: %s", p.Id.String(), nodeId.String(), err.Error())
						return
					}
					select {
					case chanBlk <- blk:
						log.Printf("[debug] [%s] got block %s from [%s]", p.Id.String(), hash.Hex(), nodeId.String())
						blkReqDone = true
					default:
					}
				}(neighbors[j])
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case <-time.NewTimer(10 * time.Second).C:
		blk = nil
	case blk = <-chanBlk:
	}

	if blk == nil {
		return nil, errors.New(fmt.Sprintf("[algorand] [%s] cannot get block %s from other peers: %s", p.Id.String(), hash.Hex(), err.Error()))
	}
	p.chain.CacheBlock(blk)
	log.Printf("[debug] [%s] cached block %s", p.Id.String(), hash.Hex())
	return blk, nil
}

func (p *Peer) sendBalance() {
	for {
		time.Sleep(5 * time.Second)
		p.lm.AddBalance(p.Id.String(), p.GetBalance())
	}
}

// processMain performs the main processing of algorand algorithm.
func (p *Peer) processMain() {
	if cmn.MetricsRound == p.round() {
		cmn.ProposerSelectedHistogram.Update(cmn.ProposerSelectedCounter.Count())
		cmn.ProposerSelectedCounter.Clear()
		cmn.MetricsRound = p.round() + 1
	}
	currRound := p.round() + 1
	block := p.blockProposal(false) // not modified
	log.Printf("[algorand] [%s] init BA with block #%d %s (%d txs), is empty? %v", p.Id.String(), block.Round, block.Hash(), len(block.Txs), block.Signature == nil)

	// 2. init BA with block with the highest priority.
	consensusType, block := p.BA(currRound, block) // not modified

	// 3. reach consensus on a FINAL or TENTATIVE new block.
	if consensusType == params.FINAL_CONSENSUS {
		log.Printf("[algorand] [%s] reach final consensus at round %d, block (%d txs) hash %s, is empty? %v", p.Id.String(), currRound, len(block.Txs), block.Hash(), block.Signature == nil)
	} else {
		log.Printf("[algorand] [%s] reach tentative consensus at round %d, block (%d txs) hash %s, is empty? %v", p.Id.String(), currRound, len(block.Txs), block.Hash(), block.Signature == nil)
	}

	if len(block.Txs) > 0 {
		parentBlk := p.lastBlock()
		txSet := &msg.ProposedTx{Epoch: p.txEpoch, Txs: block.Txs}
		st := p.executeTxSet(txSet, parentBlk.StateHash, p.permanentTxStorage)
		//st := p.recentMPT

		if bytes.Equal(block.StateHash, st.RootHash()) {
			st.Commit()
		} else {
			block.StateHash = parentBlk.StateHash
		}
	}

	p.chain.Add(block)
	go p.lm.AddFinalBlk(block.Hash(), block.Round)
	go p.oracle.AddBlk(block)
	//time.Sleep(5 * time.Second) // sleep to allow all peers add block to their storage, we change params.R and simulate the algorand with changing seed
}

// run performs the all procedures of Algorand algorithm in infinite loop.
func (p *Peer) Run() {
	// sleep 1 second for all peers ready.
	time.Sleep(5 * time.Second)
	log.Printf("[alogrand] [%s] found %d peers", p.Id.String(), p.node.GetNeighborList().Len())

	// propose block
	go p.proposeOraclePeer() // run this func continuously to propose oracle peer every oracle epoch
	go p.forkLoop()
	for {
		select {
		case <-p.quitCh:
			return
		default:
			p.processMain() // process algorand functions
		}
	}
}
