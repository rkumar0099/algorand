package message

import (
	"bytes"
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/crypto"
	"github.com/rkumar0099/algorand/params"
)

const (
	// message type
	VOTE = iota
	BLOCK_PROPOSAL
	FORK_PROPOSAL
	BLOCK
	TRANSACTION
	CONTRIBUTION
	STATEHASH
)

func (t *Transaction) Serialize() ([]byte, error) {
	return proto.Marshal(t)
}

func (t *Transaction) Deserialize(data []byte) error {
	return proto.Unmarshal(data, t)
}

func (t *Transaction) RecoverPubkey() *crypto.PublicKey {
	return crypto.RecoverPubkey(t.Signature)
}

func (t *Transaction) VerifySign() error {
	pubkey := t.RecoverPubkey()
	if !bytes.Equal(pubkey.Address().Bytes(), t.From) {
		return errors.New("Transaction is not signed by transaction sender")
	}
	data := bytes.Join([][]byte{
		t.From,
		t.To,
		common.Uint2Bytes(t.Nonce),
		common.Uint2Bytes(t.Type),
		t.Data,
	}, nil)
	return pubkey.VerifySign(data, t.Signature)
}

func (t *Transaction) Sign(privKey *crypto.PrivateKey) error {
	data := bytes.Join([][]byte{
		t.From,
		t.To,
		common.Uint2Bytes(t.Nonce),
		common.Uint2Bytes(t.Type),
		t.Data,
	}, nil)
	sig, err := privKey.Sign(data)
	t.Signature = sig
	return err
}

func (t *Transaction) Hash() common.Hash {
	sig := t.Signature
	t.Signature = nil
	data, _ := t.Serialize()
	t.Signature = sig
	return common.Sha256(data)
}

func (v *VoteMessage) Serialize() ([]byte, error) {
	return proto.Marshal(v)
}

func (v *VoteMessage) Deserialize(data []byte) error {
	return proto.Unmarshal(data, v)
}

func (v *VoteMessage) VerifySign() error {
	pubkey := v.RecoverPubkey()
	data := bytes.Join([][]byte{
		common.Uint2Bytes(v.Round),
		[]byte(v.Event),
		v.VRF,
		v.Proof,
		v.ParentHash,
		v.Hash,
	}, nil)
	return pubkey.VerifySign(data, v.Signature)
}

func (v *VoteMessage) Sign(priv *crypto.PrivateKey) ([]byte, error) {
	data := bytes.Join([][]byte{
		common.Uint2Bytes(v.Round),
		[]byte(v.Event),
		v.VRF,
		v.Proof,
		v.ParentHash,
		v.Hash,
	}, nil)
	sign, err := priv.Sign(data)
	if err != nil {
		return nil, err
	}
	v.Signature = sign
	return sign, nil
}

func (v *VoteMessage) RecoverPubkey() *crypto.PublicKey {
	return crypto.RecoverPubkey(v.Signature)
}

func (b *Proposal) Serialize() ([]byte, error) {
	return proto.Marshal(b)
}

func (b *Proposal) Deserialize(data []byte) error {
	return proto.Unmarshal(data, b)
}

func (b *Proposal) PublicKey() *crypto.PublicKey {
	return &crypto.PublicKey{Pk: b.Pubkey}
}

func (b *Proposal) Address() common.Address {
	return common.BytesToAddress(b.Pubkey)
}

func (b *Proposal) Verify(weight uint64, m []byte) error {
	// verify vrf
	pubkey := b.PublicKey()
	if err := pubkey.VerifyVRF(b.Proof, m); err != nil {
		return err
	}

	// verify priority
	subusers := common.SubUsers(params.ExpectedBlockProposers, weight, b.VRF)
	if bytes.Compare(common.MaxPriority(b.VRF, subusers), b.Prior) != 0 {
		return errors.New("max priority mismatch")
	}

	return nil
}

func (blk *Block) Serialize() ([]byte, error) {
	return proto.Marshal(blk)
}

func (blk *Block) Deserialize(data []byte) error {
	return proto.Unmarshal(data, blk)
}

func (blk *Block) Hash() common.Hash {
	data, _ := blk.Serialize()
	return common.Sha256(data)
}

func (blk *Block) RecoverPubkey() *crypto.PublicKey {
	return crypto.RecoverPubkey(blk.Signature)
}

func (pr *PendingRequest) Hash() common.Hash {
	data := bytes.Join([][]byte{
		common.Uint2Bytes(pr.Nonce),
		[]byte(pr.URL),
	}, nil)
	return common.Sha256(data)
}

func (pt *ProposedTx) Serialize() ([]byte, error) {
	return proto.Marshal(pt)
}

func (pt *ProposedTx) Deserialize(data []byte) error {
	return proto.Unmarshal(data, pt)
}

func (pt *ProposedTx) Hash() common.Hash {
	data, _ := pt.Serialize()
	return common.Sha256(data)
}

func (opp *OraclePeerProposal) pubkey() *crypto.PublicKey {
	return &crypto.PublicKey{Pk: opp.Pubkey}
}

func (opp *OraclePeerProposal) Verify(message []byte) error {
	pubkey := opp.pubkey()
	if err := pubkey.VerifyVRF(opp.Proof, message); err != nil {
		return err
	}
	//log.Printf("Oracle peer proposal successfully verified\n")
	return nil
}
