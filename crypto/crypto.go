package crypto

import (
	"crypto"
	"crypto/rand"
	"fmt"

	"github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/vrf"
	"golang.org/x/crypto/ed25519"
)

type PublicKey struct {
	Pk ed25519.PublicKey
}

func (pub *PublicKey) Bytes() []byte {
	return pub.Pk
}

func (pub *PublicKey) Address() common.Address {
	return common.BytesToAddress(pub.Pk)
}

func (pub *PublicKey) VerifySign(m, sign []byte) error {
	signature := sign[ed25519.PublicKeySize:]
	if ok := ed25519.Verify(pub.Pk, m, signature); !ok {
		return fmt.Errorf("signature invalid")
	}
	return nil
}

func (pub *PublicKey) VerifyVRF(proof, m []byte) error {
	_, err := vrf.ECVRF_verify(pub.Pk, proof, m)
	if err != nil {
		return err
	}
	return nil
}

type PrivateKey struct {
	sk ed25519.PrivateKey
}

func (priv *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{priv.sk.Public().(ed25519.PublicKey)}
}

func (priv *PrivateKey) Sign(m []byte) ([]byte, error) {
	sign, err := priv.sk.Sign(rand.Reader, m, crypto.Hash(0))
	if err != nil {
		return nil, err
	}
	pubkey := priv.sk.Public().(ed25519.PublicKey)
	return append(pubkey, sign...), nil
}

func (priv *PrivateKey) Evaluate(m []byte) (value, proof []byte, err error) {
	proof, err = vrf.ECVRF_prove(priv.PublicKey().Pk, priv.sk, m)
	if err != nil {
		return
	}
	value = vrf.ECVRF_proof2hash(proof)
	return
}

func NewKeyPair() (*PublicKey, *PrivateKey, error) {
	pk, sk, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	return &PublicKey{pk}, &PrivateKey{sk}, nil
}

func RecoverPubkey(sign []byte) *PublicKey {
	pubkey := sign[:ed25519.PublicKeySize]
	return &PublicKey{pubkey}
}

func MyKey() (a string, b string, err error) {
	a = "Rabindar"
	b = "Kumar"
	err = nil
	return
}
