package pool

import (
	"bytes"
	"sync"
	"time"

	msg "github.com/rkumar0099/algorand/message"
)

type ProposalPool struct {
	pool          map[uint64]*msg.Proposal
	lock          *sync.RWMutex
	cleanDuration time.Duration
	verifier      ProposalVerifier
}

type ProposalVerifier func(*msg.Proposal) bool

func NewProposalPool(cleanDuration time.Duration, verifier ProposalVerifier) *ProposalPool {
	return &ProposalPool{
		pool:          make(map[uint64]*msg.Proposal),
		lock:          &sync.RWMutex{},
		cleanDuration: cleanDuration,
		verifier:      verifier,
	}
}

func (pp *ProposalPool) Update(bp *msg.Proposal, typ int) {
	pp.lock.Lock()
	defer pp.lock.Unlock()
	maxProposal, ok := pp.pool[bp.Round]
	if ok && ((typ == msg.BLOCK_PROPOSAL && bytes.Compare(bp.Prior, maxProposal.Prior) <= 0) || (typ == msg.FORK_PROPOSAL && bp.Round <= maxProposal.Round)) {
		return
	}
	if pp.verifier(bp) {
		pp.pool[bp.Round] = bp
		if !ok {
			time.AfterFunc(pp.cleanDuration, func() {
				pp.lock.Lock()
				defer pp.lock.Unlock()
				delete(pp.pool, bp.Round)
			})
		}
	}
}

func (pp *ProposalPool) GetMaxProposal(round uint64) *msg.Proposal {
	pp.lock.RLock()
	defer pp.lock.RUnlock()
	return pp.pool[round]
}
