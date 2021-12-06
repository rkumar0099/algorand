package peer

import (
	"time"

	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/params"
)

func (p *Peer) proposeOraclePeer() {
	time.Sleep(10 * time.Second)
	for {
		time.Sleep(5 * time.Second)
		p.oracleEpoch += 1
		seed := p.oracle.SortitionSeed(p.oracleEpoch)
		role := role(params.OraclePeer, p.oracleEpoch, params.ORACLE)
		vrf, proof, usr := p.sortition(seed, role, params.ExpectedOraclePeers, p.tokenOwn())
		if usr > 0 {
			opp := &msg.OraclePeerProposal{
				Pubkey: p.pubkey.Bytes(),
				Proof:  proof,
				Vrf:    vrf,
				Epoch:  p.oracleEpoch,
				Weight: p.tokenOwn(),
			}
			//log.Println(opp)
			p.oracle.AddOPP(opp)
		}
	}
}

func (p *Peer) handleOracleBlk(data []byte) {
	p.oracleFinalContributions <- data
}
