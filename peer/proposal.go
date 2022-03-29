package peer

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	cmn "github.com/rkumar0099/algorand/common"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
)

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
	log.Printf("[Debug] [Peer BLK] Num Cont: %d\n", len(p.finalContributions))
	if len(p.finalContributions) > 0 {
		p.pt = <-p.finalContributions
		txSet := p.pt
		blk.Txs = txSet.Txs
		st, responses := p.executeTxSet(txSet, p.lastState, p.permanentTxStorage)
		blk.ResTx = responses
		blk.StateHash = st.RootHash()
		blk.TxEpoch = txSet.Epoch
	}

	bhash := blk.Hash()
	sign, _ := p.privkey.Sign(bhash.Bytes())
	blk.Signature = sign
	s := fmt.Sprintf("[alogrand] [%s] propose a new block with %d txs: #%d %s, stateHash: %s, parent: %s\n", p.Id.String(), len(blk.Txs), blk.Round, blk.Hash(), hex.EncodeToString(blk.StateHash), hex.EncodeToString(blk.ParentHash))
	log.Print(s)
	//p.lm.AddLog(s)
	//log.Printf
	return blk
}

func (p *Peer) handleBlockProposal(data []byte) {
	bp := &msg.Proposal{}
	if err := bp.Deserialize(data); err != nil {
		s := fmt.Sprintf("[algorand] [%s] Received invalid proposal message: %s", p.Id.String(), err.Error())
		log.Print(s)
		//p.lm.AddLog(s)
		return
	}
	p.proposalPool.Update(bp, msg.BLOCK_PROPOSAL)
}

func (p *Peer) proposalVerifier(bp *msg.Proposal) bool {
	if err := bp.Verify(p.weight(bp.Address()), constructSeed(p.sortitionSeed(bp.Round), role(params.Proposer, bp.Round, params.PROPOSE))); err != nil {
		s := fmt.Sprintf("[algorand] [%s] Received invaild proposal: %5s\n", p.Id.String(), err.Error())
		log.Print(s)
		//p.lm.AddLog(s)
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

func (p *Peer) handleForkProposal(data []byte) {
	bp := &msg.Proposal{}
	if err := bp.Deserialize(data); err != nil {
		log.Printf("[algorand] [%s] Received invalid proposal message: %s", p.Id.String(), err.Error())
		return
	}
	p.proposalPool.Update(bp, msg.FORK_PROPOSAL)
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

func (p *Peer) proposeFork() *msg.Block {
	longest := p.lastBlock()
	return p.emptyBlock(p.round()+1, longest.Hash())
}
