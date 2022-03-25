package peer

import (
	"github.com/rkumar0099/algorand/oracle"
	"github.com/rkumar0099/algorand/params"
)

func (p *Peer) proposeOraclePeer() (*oracle.ResOPP, error) {
	seed := p.oracle.SortitionSeed(1)
	role := role(params.OraclePeer, p.oracleEpoch, params.ORACLE)
	vrf, proof, usr := p.sortition(seed, role, params.ExpectedOraclePeers, p.tokenOwn())
	if usr > 0 {
		opp := &oracle.ResOPP{
			Proof:  proof,
			VRF:    vrf,
			Pubkey: p.pubkey.Bytes(),
			Weight: p.tokenOwn(),
		}
		return opp, nil
	}
	return nil, nil
}

func (p *Peer) handleOracleBlk(data []byte) {
	p.oracleFinalContributions <- data
}
