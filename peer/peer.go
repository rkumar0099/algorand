package peer

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
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

	Id      gossip.NodeId
	privkey *crypto.PrivateKey
	pubkey  *crypto.PublicKey

	permanentTxStorage kvstore.KVStore
	tempTxStorage      kvstore.KVStore

	chain       *blockchain.Blockchain
	quitCh      chan struct{}
	hangForever chan struct{}

	node          *gossip.Node
	grpcServer    *grpc.Server
	msgAgent      *msg.MsgAgent
	ServiceServer *service.Server

	votePool     *pool.VotePool
	proposalPool *pool.ProposalPool

	manage *manage.Manage
	oracle *oracle.Oracle
	lm     *logs.LogManager

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
	go p.oracle.AddBlk(blk)
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

// seed returns the vrf-based seed of block r.
func (p *Peer) vrfSeed(round uint64) (seed, proof []byte, err error) {
	if round == 0 {
		return p.chain.Genesis.Seed, nil, nil
	}
	lastBlock := p.chain.GetByRound(round - 1)
	// last block is not genesis, verify the seed r-1.
	if round != 1 {
		lastParentBlock, err := p.getBlock(cmn.BytesToHash(lastBlock.ParentHash))
		if err != nil {
			log.Printf("[Error] node %s hang forever: cannot get last block", p.Id.String())
			<-p.hangForever
		}
		if lastBlock.Proof != nil {
			// vrf-based seed
			pubkey := crypto.RecoverPubkey(lastBlock.Signature)
			m := bytes.Join([][]byte{lastParentBlock.Seed, cmn.Uint2Bytes(lastBlock.Round)}, nil)
			err = pubkey.VerifyVRF(lastBlock.Proof, m)
		} else if bytes.Compare(lastBlock.Seed, cmn.Sha256(
			bytes.Join([][]byte{
				lastParentBlock.Seed,
				cmn.Uint2Bytes(lastBlock.Round)},
				nil)).Bytes()) != 0 {
			// hash-based seed
			err = errors.New("hash seed invalid")
		}
		if err != nil {
			// seed r-1 invalid
			return cmn.Sha256(bytes.Join([][]byte{lastBlock.Seed, cmn.Uint2Bytes(lastBlock.Round + 1)}, nil)).Bytes(), nil, nil
		}
	}

	seed, proof, err = p.privkey.Evaluate(bytes.Join([][]byte{lastBlock.Seed, cmn.Uint2Bytes(lastBlock.Round + 1)}, nil))
	return
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

// sortitionSeed returns the selection seed with a refresh interval R.
func (p *Peer) sortitionSeed(round uint64) []byte {
	realR := round - 1
	mod := round % params.R
	if realR < mod {
		realR = 0
	} else {
		realR -= mod
	}

	return p.chain.GetByRound(realR).Seed
}

func (p *Peer) Address() cmn.Address {
	return cmn.BytesToAddress(p.pubkey.Bytes())
}

// run performs the all procedures of Algorand algorithm in infinite loop.
func (p *Peer) Run() {
	// sleep 1 second for all peers ready.
	time.Sleep(5 * time.Second)
	log.Printf("[alogrand] [%s] found %d peers", p.Id.String(), p.node.GetNeighborList().Len())

	// propose block
	go p.proposeOraclePeer()
	go p.forkLoop()
	for {
		select {
		case <-p.quitCh:
			return
		default:
			p.processMain()
		}
	}
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

// processMain performs the main processing of algorand algorithm.
func (p *Peer) processMain() {
	if cmn.MetricsRound == p.round() {
		cmn.ProposerSelectedHistogram.Update(cmn.ProposerSelectedCounter.Count())
		cmn.ProposerSelectedCounter.Clear()
		cmn.MetricsRound = p.round() + 1
	}
	currRound := p.round() + 1
	block := p.blockProposal(false)
	log.Printf("[algorand] [%s] init BA with block #%d %s (%d txs), is empty? %v", p.Id.String(), block.Round, block.Hash(), len(block.Txs), block.Signature == nil)

	// 2. init BA with block with the highest priority.
	consensusType, block := p.BA(currRound, block)

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

// processForkResolve performs a special algorand processing to resolve fork.
func (p *Peer) processForkResolve() {
	// force quit the hanging in BA if any.
	close(p.hangForever)

	if cmn.MetricsRound == p.round() {
		cmn.ProposerSelectedHistogram.Update(cmn.ProposerSelectedCounter.Count())
		cmn.ProposerSelectedCounter.Clear()
		cmn.MetricsRound = p.round() + 1
	}
	// propose fork
	longest := p.blockProposal(true)
	// init BA with a highest priority fork
	_, fork := p.BA(longest.Round, longest)
	// commit fork
	p.chain.ResolveFork(fork)

	p.hangForever = make(chan struct{})
}

// proposeBlock proposes a new block.
func (p *Peer) proposeBlock() *msg.Block {
	currRound := p.round() + 1
	parentBlk := p.lastBlock()
	seed, proof, err := p.vrfSeed(currRound)
	if err != nil {
		return p.emptyBlock(currRound, p.lastBlock().Hash())
	}

	// random data field to simulate different version of block.
	blk := &msg.Block{
		Round:      currRound,
		Seed:       seed,
		ParentHash: p.lastBlock().Hash().Bytes(),
		Author:     p.pubkey.Address().Bytes(),
		Time:       time.Now().Unix(),
		Proof:      proof,
		Data:       nil,
		StateHash:  parentBlk.StateHash,
	}

	//var txSet *msg.ProposedTx
	if len(p.finalContributions) > 0 {
		txSet := <-p.finalContributions
		blk.Txs = txSet.Txs
		st := p.executeTxSet(txSet, blk.StateHash, p.permanentTxStorage)
		blk.StateHash = st.RootHash()
		blk.TxEpoch = txSet.Epoch
	}

	bhash := blk.Hash()
	sign, _ := p.privkey.Sign(bhash.Bytes())
	blk.Signature = sign
	log.Printf("[alogrand] [%s] propose a new block with %d txs: #%d %s, stateHash: %s, parent: %s", p.Id.String(), len(blk.Txs), blk.Round, blk.Hash(), hex.EncodeToString(blk.StateHash), hex.EncodeToString(blk.ParentHash))
	return blk
}

func (p *Peer) proposeFork() *msg.Block {
	longest := p.lastBlock()
	return p.emptyBlock(p.round()+1, longest.Hash())
}

// blockProposal performs the block proposal procedure.
func (p *Peer) blockProposal(resolveFork bool) *msg.Block {
	round := p.round() + 1
	vrf, proof, subusers := p.sortition(p.sortitionSeed(round), role(params.Proposer, round, params.PROPOSE), params.ExpectedBlockProposers, p.tokenOwn())
	// have been selected.
	//log.Printf("node %d get %d sub-users in block proposal", alg.id, subusers)

	if subusers > 0 {
		cmn.ProposerSelectedCounter.Inc(1)
		var (
			newBlk       *msg.Block
			proposalType int
		)

		if !resolveFork {
			newBlk = p.proposeBlock()
			proposalType = msg.BLOCK_PROPOSAL
		} else {
			newBlk = p.proposeFork()
			proposalType = msg.FORK_PROPOSAL
		}

		proposal := &msg.Proposal{
			Round:  newBlk.Round,
			Hash:   newBlk.Hash().Bytes(),
			Prior:  cmn.MaxPriority(vrf, subusers),
			VRF:    vrf,
			Proof:  proof,
			Pubkey: p.pubkey.Bytes(),
		}
		// p.proposalPool.Update(proposal, proposalType)
		p.chain.CacheBlock(newBlk)
		blkMsg, _ := newBlk.Serialize()
		proposal.Block = blkMsg
		proposalMsg, _ := proposal.Serialize()

		p.gossip(proposalType, proposalMsg)
	}

	// wait for λstepvar + λpriority time to identify the highest priority.
	timeoutForPriority := time.NewTimer(params.LamdaStepvar + params.LamdaPriority)
	<-timeoutForPriority.C

	// timeout for block gossiping.
	timeoutForBlockFlying := time.NewTimer(params.LamdaBlock)
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-timeoutForBlockFlying.C:
			// empty block
			return p.emptyBlock(round, p.lastBlock().Hash())
		case <-ticker.C:
			// get the block with the highest priority
			pp := p.proposalPool.GetMaxProposal(round)
			if pp == nil {
				continue
			}
			blk, _ := p.getBlock(cmn.BytesToHash(pp.Hash))
			if blk != nil {
				return blk
			}
		}
	}
}

// sortition runs cryptographic selection procedure and returns vrf,proof and amount of selected sub-users.
func (p *Peer) sortition(seed, role []byte, expectedNum int, weight uint64) (vrf, proof []byte, selected int) {
	vrf, proof, _ = p.privkey.Evaluate(constructSeed(seed, role))
	selected = cmn.SubUsers(expectedNum, weight, vrf)
	return
}

// verifySort verifies the vrf and returns the amount of selected sub-users.
func (p *Peer) verifySort(vrf, proof, seed, role []byte, expectedNum int) int {
	if err := p.pubkey.VerifyVRF(proof, constructSeed(seed, role)); err != nil {
		return 0
	}

	return cmn.SubUsers(expectedNum, p.tokenOwn(), vrf)
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

// committeeVote votes for `value`.
func (p *Peer) committeeVote(round uint64, event string, expectedNum int, hash cmn.Hash) error {
	if p.maliciousType == params.EvilVoteNothing {
		// vote nothing
		return nil
	}

	vrf, proof, j := p.sortition(p.sortitionSeed(round), role(params.Committee, round, event), expectedNum, p.tokenOwn())

	if j > 0 {
		// Gossip vote message
		voteMsg := &msg.VoteMessage{
			Round:      round,
			Event:      event,
			VRF:        vrf,
			Proof:      proof,
			ParentHash: p.chain.Last.Hash().Bytes(),
			Hash:       hash.Bytes(),
		}

		_, err := voteMsg.Sign(p.privkey)
		if err != nil {
			return err
		}

		data, err := voteMsg.Serialize()
		if err != nil {
			return err
		}
		p.gossip(msg.VOTE, data)
	}
	return nil
}

// BA runs BA* for the next round, with a proposed block.
func (p *Peer) BA(round uint64, block *msg.Block) (int8, *msg.Block) {
	var (
		newBlk *msg.Block
		hash   cmn.Hash
	)
	if p.maliciousType == params.EvilVoteEmpty {
		hash = p.emptyHash(round, cmn.BytesToHash(block.ParentHash))
		p.reduction(round, hash)
	} else {
		hash = p.reduction(round, block.Hash())
	}
	hash = p.binaryBA(round, hash)
	prevHash := p.lastBlock().Hash()
	emptyBlk := p.emptyBlock(round, prevHash)
	p.chain.CacheBlock(emptyBlk)
	r, _ := p.votePool.CountVotes(round, params.FINAL, params.FinalThreshold, params.ExpectedFinalCommitteeMembers, params.LamdaStep)
	if hash == p.emptyHash(round, prevHash) {
		// empty block
		newBlk = emptyBlk
	} else {
		var err error
		newBlk, err = p.getBlock(hash)
		if err != nil {
			log.Printf("[Algorand] [%s] hang forever becaue BA error: %s", p.Id.String(), err.Error())
			<-p.hangForever
		}
	}
	if r == hash {
		return params.FINAL_CONSENSUS, newBlk
	} else {
		return params.TENTATIVE_CONSENSUS, newBlk
	}
}

// The two-step reduction.
func (p *Peer) reduction(round uint64, hash cmn.Hash) cmn.Hash {
	// step 1: gossip the block hash
	p.committeeVote(round, params.REDUCTION_ONE, params.ExpectedCommitteeMembers, hash)

	// other users might still be waiting for block proposals,
	// so set timeout for λblock + λstep
	hash1, err := p.votePool.CountVotes(round, params.REDUCTION_ONE, params.ThresholdOfBAStep, params.ExpectedCommitteeMembers, params.LamdaBlock+params.LamdaStep)

	// step 2: re-gossip the popular block hash
	empty := p.emptyHash(round, p.chain.Last.Hash())

	if err == cmn.ErrCountVotesTimeout {
		p.committeeVote(round, params.REDUCTION_TWO, params.ExpectedCommitteeMembers, empty)
	} else {
		p.committeeVote(round, params.REDUCTION_TWO, params.ExpectedCommitteeMembers, hash1)
	}

	hash2, err := p.votePool.CountVotes(round, params.REDUCTION_TWO, params.ThresholdOfBAStep, params.ExpectedCommitteeMembers, params.LamdaStep)
	if err == cmn.ErrCountVotesTimeout {
		return empty
	}
	return hash2
}

// binaryBA executes until consensus is reached on either the given `hash` or `empty_hash`.
func (p *Peer) binaryBA(round uint64, hash cmn.Hash) cmn.Hash {
	var (
		step = 1
		r    = hash
		err  error
		coin int64
	)
	empty := p.emptyHash(round, p.chain.Last.Hash())
	defer func() {
		log.Printf("[algorand] [%s] complete binaryBA with %d steps", p.Id.String(), step)
	}()
	for step < params.MAXSTEPS {
		p.committeeVote(round, strconv.Itoa(step), params.ExpectedCommitteeMembers, r)
		r, err = p.votePool.CountVotes(round, strconv.Itoa(step), params.ThresholdOfBAStep, params.ExpectedCommitteeMembers, params.LamdaStep)
		if err != nil {
			r = hash
		} else if r != empty {
			for s := step + 1; s <= step+3; s++ {
				p.committeeVote(round, strconv.Itoa(s), params.ExpectedCommitteeMembers, r)
			}
			if step == 1 {
				p.committeeVote(round, params.FINAL, params.ExpectedFinalCommitteeMembers, r)
			}
			return r
		}
		step++

		p.committeeVote(round, strconv.Itoa(step), params.ExpectedCommitteeMembers, r)
		r, err = p.votePool.CountVotes(round, strconv.Itoa(step), params.ThresholdOfBAStep, params.ExpectedCommitteeMembers, params.LamdaStep)
		if err != nil {
			r = empty
		} else if r == empty {
			for s := step + 1; s <= step+3; s++ {
				p.committeeVote(round, strconv.Itoa(s), params.ExpectedCommitteeMembers, r)
			}
			return r
		}
		step++

		p.committeeVote(round, strconv.Itoa(step), params.ExpectedCommitteeMembers, r)
		r, coin, err = p.votePool.CountVotesAndCoin(round, strconv.Itoa(step), params.ThresholdOfBAStep, params.ExpectedCommitteeMembers, params.LamdaStep)
		if err != nil {
			if coin == 0 {
				r = hash
			} else {
				r = empty
			}
		}
	}

	//log.Printf("reach the maxstep hang forever")
	// hang forever
	<-p.hangForever
	return common.Hash{}
}

func (p *Peer) voteVerifier(vote *msg.VoteMessage, expectedNum int) int {
	if err := vote.VerifySign(); err != nil {
		return 0
	}

	prevHash := common.BytesToHash(vote.ParentHash)
	if prevHash != p.chain.Last.Hash() {
		return 0
	}

	return p.verifySort(vote.VRF, vote.Proof, p.sortitionSeed(vote.Round), role(params.Committee, vote.Round, vote.Event), expectedNum)
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

func (p *Peer) handleBlockProposal(data []byte) {
	bp := &msg.Proposal{}
	if err := bp.Deserialize(data); err != nil {
		//log.Printf("[algorand] [%s] Received invalid proposal message: %s", p.Id.String(), err.Error())
		return
	}
	p.proposalPool.Update(bp, msg.BLOCK_PROPOSAL)
}

func (p *Peer) handleForkProposal(data []byte) {
	bp := &msg.Proposal{}
	if err := bp.Deserialize(data); err != nil {
		log.Printf("[algorand] [%s] Received invalid proposal message: %s", p.Id.String(), err.Error())
		return
	}
	p.proposalPool.Update(bp, msg.FORK_PROPOSAL)
}

func (p *Peer) proposalVerifier(bp *msg.Proposal) bool {
	if err := bp.Verify(p.weight(bp.Address()), constructSeed(p.sortitionSeed(bp.Round), role(params.Proposer, bp.Round, params.PROPOSE))); err != nil {
		log.Printf("[algorand] [%s] Received invaild proposal: 5s", p.Id.String(), err.Error())
		return false
	}
	blk := &msg.Block{}
	if err := blk.Deserialize(bp.Block); err != nil {
		log.Printf("[algorand] [%s] Received proposal with invalid block: %s", p.Id.String(), err.Error())
		return false
	}
	p.chain.CacheBlock(blk)

	return true
}

func (p *Peer) handleVote(data []byte) {
	vote := &msg.VoteMessage{}
	if err := vote.Deserialize(data); err != nil {
		log.Printf("[algorand] [%s] Received invalid vote: %s", p.Id.String(), err.Error())
		return
	}
	p.votePool.HandleVote(vote)
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
