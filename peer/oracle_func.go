package peer

import (
	"time"

	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
)

func (p *Peer) proposeOraclePeer() {
	for {
		time.Sleep(1 * time.Minute)
		p.oracleEpoch += 1
		seed := p.sortitionSeed(p.oracleEpoch)
		role := role(params.OracleProposer, p.oracleEpoch, params.ORACLE)
		vrf, proof, usr := p.sortition(seed, role, params.ExpectedOraclePeers, p.tokenOwn())
		if usr > 0 {
			opp := &msg.OraclePeerProposal{
				Pubkey: p.pubkey.Address().Bytes(),
				Proof:  proof,
				Vrf:    vrf,
				Epoch:  p.oracleEpoch,
				Weight: p.tokenOwn(),
			}
			p.oracle.AddOPP(opp)
		}
	}
}

func (p *Peer) handleOracleBlk(data []byte) {
	p.oracleFinalContributions <- data
}
