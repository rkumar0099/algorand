package peer

import (
	"bytes"
	"errors"
	"log"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"
	"github.com/rkumar0099/algorand/params"
)

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
