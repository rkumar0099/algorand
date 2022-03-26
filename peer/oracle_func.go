package peer

import (
	"log"

	"github.com/rkumar0099/algorand/oracle"
	"github.com/rkumar0099/algorand/params"
)

func (p *Peer) proposeOraclePeer(epoch uint64) (*oracle.ResOPP, error) {
	log.Println("[Debug] Proposing Oracle Peer")
	seed := p.oracle.SortitionSeed(1)
	role := role(params.OraclePeer, epoch, params.ORACLE)
	vrf, proof, usr := p.sortition(seed, role, params.ExpectedOraclePeers, p.tokenOwn())
	opp := &oracle.ResOPP{}
	opp.Weight = 0
	if usr > 0 {
		log.Printf("[Debug] %s selected as oracle peer\n", p.Id.String())
		opp.Proof = proof
		opp.VRF = vrf
		opp.Pubkey = p.pubkey.Bytes()
		opp.Weight = p.tokenOwn()
	}
	return opp, nil
}

func (p *Peer) handleOracleBlk(data []byte) {
	p.oracleFinalContributions <- data
}
