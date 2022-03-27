package peer

import (
	"github.com/rkumar0099/algorand/client"
	msg "github.com/rkumar0099/algorand/message"
)

func (p *Peer) handleTx(req *client.ReqTx) (*client.ResEmpty, error) {
	t := req.Type
	if t > 0 && t < 5 {
		// normal tx
		tx := &msg.Transaction{
			Type: t,
			Addr: req.Addr,
			Data: req.Data,
		}

		p.manage.AddTransaction(tx)

	} else if t > 4 && t < 6 {
		// external world tx

		pr := &msg.PendingRequest{
			Type: req.Type,
			Addr: req.Addr,
			Data: req.Data,
		}
		p.oracle.AddEWTx(pr)
	}

	re := &client.ResEmpty{}
	return re, nil
}

func (p *Peer) sendRes(req *client.ResTx) (*client.ResEmpty, error) {
	return &client.ResEmpty{}, nil
}
