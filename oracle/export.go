package oracle

import (
	"bytes"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/params"
)

// sortitionSeed returns the selection seed with a refresh interval R.
func (o *Oracle) SortitionSeed(round uint64) []byte {
	realR := round - 1
	mod := round % params.R
	if realR < mod {
		realR = 0
	} else {
		realR -= mod
	}

	blk := o.GetBlkByRound(realR)
	//log.Println(blk, realR)
	return blk.Seed
}

// sortition runs cryptographic selection procedure and returns vrf,proof and amount of selected sub-users.
func (o *Oracle) Sortition(seed, role []byte, expectedNum int, weight uint64) (vrf, proof []byte, selected int) {
	vrf, proof, _ = o.privkey.Evaluate(constructSeed(seed, role))
	selected = cmn.SubUsers(expectedNum, weight, vrf)
	return
}

// role returns the role bytes from current round and step
func Role(iden string, round uint64, event string) []byte {
	return bytes.Join([][]byte{
		[]byte(iden),
		common.Uint2Bytes(round),
		[]byte(event),
	}, nil)
}

// constructSeed construct a new bytes for vrf generation.
func ConstructSeed(seed, role []byte) []byte {
	return bytes.Join([][]byte{seed, role}, nil)
}
