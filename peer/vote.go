package peer

import (
	"log"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
)

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

func (p *Peer) handleVote(data []byte) {
	vote := &msg.VoteMessage{}
	if err := vote.Deserialize(data); err != nil {
		log.Printf("[algorand] [%s] Received invalid vote: %s", p.Id.String(), err.Error())
		return
	}
	p.votePool.HandleVote(vote)
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
