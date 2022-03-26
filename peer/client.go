package peer

import (
	"github.com/rkumar0099/algorand/client"
)

func (p *Peer) handleTx(req *client.ReqTx) (*client.ResEmpty, error) {
	re := &client.ResEmpty{}
	return re, nil
}

func (p *Peer) sendRes(req *client.ResTx) (*client.ResEmpty, error) {
	return &client.ResEmpty{}, nil
}
