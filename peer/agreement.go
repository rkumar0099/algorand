package peer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
)

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
			s := fmt.Sprintf("[Algorand] [%s] hang forever becaue BA error: %s\n", p.Id.String(), err.Error())
			log.Print(s)
			p.lm.AddLog(s)
			//log.Printf("[Algorand] [%s] hang forever becaue BA error: %s", p.Id.String(), err.Error())
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
		s := fmt.Sprintf("[algorand] [%s] complete binaryBA with %d steps\n", p.Id.String(), step)
		log.Print(s)
		p.lm.AddLog(s)
		//log.Printf("[algorand] [%s] complete binaryBA with %d steps", p.Id.String(), step)
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
