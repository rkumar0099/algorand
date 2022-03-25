package peer

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/rkumar0099/algorand/logs"
	"github.com/rkumar0099/algorand/manage"
	"github.com/rkumar0099/algorand/mpt/kvstore"
	"github.com/rkumar0099/algorand/oracle"
	"google.golang.org/grpc"

	"github.com/rkumar0099/algorand/blockchain"
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

	chain       *blockchain.Blockchain // blockchain
	quitCh      chan struct{}
	hangForever chan struct{}

	node          *gossip.Node // node to communicate msgs over network
	grpcServer    *grpc.Server
	msgAgent      *msg.MsgAgent
	ServiceServer *service.Server

	votePool     *pool.VotePool
	proposalPool *pool.ProposalPool

	pt     *msg.ProposedTx
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
	//peer.tempTxStorage = kvstore.NewMemKVStore()

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

// processMain performs the main processing of algorand algorithm
func (p *Peer) processMain() {
	if cmn.MetricsRound == p.round() {
		cmn.ProposerSelectedHistogram.Update(cmn.ProposerSelectedCounter.Count())
		cmn.ProposerSelectedCounter.Clear()
		cmn.MetricsRound = p.round() + 1
	}
	currRound := p.round() + 1
	block := p.blockProposal(false) // not modified
	//log.Printf("[algorand] [%s] init BA with block #%d %s (%d txs), is empty? %v", p.Id.String(), block.Round, block.Hash(), len(block.Txs), block.Signature == nil)

	// 2. init BA with block with the highest priority.
	consensusType, block := p.BA(currRound, block) // not modified

	// 3. reach consensus on a FINAL or TENTATIVE new block.
	if consensusType == params.FINAL_CONSENSUS {
		//log.Printf("[algorand] [%s] reach final consensus at round %d, block (%d txs) hash %s, is empty? %v", p.Id.String(), currRound, len(block.Txs), block.Hash(), block.Signature == nil)
	} else {
		//log.Printf("[algorand] [%s] reach tentative consensus at round %d, block (%d txs) hash %s, is empty? %v", p.Id.String(), currRound, len(block.Txs), block.Hash(), block.Signature == nil)
	}

	// all peers are guaranteed to receive the same block
	if len(block.Txs) > 0 {
		parentBlk := p.lastBlock()
		txSet := &msg.ProposedTx{Epoch: p.txEpoch, Txs: block.Txs}
		st := p.executeTxSet(txSet, p.lastState, p.permanentTxStorage)
		//st := p.recentMPT

		if bytes.Equal(block.StateHash, st.RootHash()) {
			// txs finalized
			log.Printf("%s finalized\n", p.pt.Hash().String())
			//p.manage.AddVote(p.pt.Hash())
			st.Commit()
			p.lastState = block.StateHash
		} else {
			block.StateHash = parentBlk.StateHash
			// txs not finalized, put all txs back to manage
		}
	}

	p.chain.Add(block) // add blk to blockchain
	go p.lm.AddFinalBlk(block.Hash(), block.Round)
	//go p.oracle.AddBlk(block)
	//time.Sleep(5 * time.Second) // sleep to allow all peers add block to their storage, we change params.R and simulate the algorand with changing seed
}

// run performs the all procedures of Algorand algorithm in infinite loop.
func (p *Peer) Run() {
	// sleep 1 second for all peers ready.
	time.Sleep(1 * time.Second)
	s := fmt.Sprintf("[alogrand] [%s] found %d peers\n", p.Id.String(), p.node.GetNeighborList().Len())
	log.Print(s)
	p.lm.AddLog(s)

	// propose block
	//go p.proposeOraclePeer() // run this func continuously to propose oracle peer every oracle epoch
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
