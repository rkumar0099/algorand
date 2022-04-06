package peer

import (
	"log"

	"github.com/rkumar0099/algorand/client"
	msg "github.com/rkumar0099/algorand/message"
)

func (p *Peer) HandleTx(req *client.ReqTx) (*client.ResEmpty, error) {
	t := req.Type
	//p.rec = true
	data, _ := req.Serialize()
	if t < 5 {
		// normal tx
		log.Println("[Debug] [Peer Tx] Received Tx from client")
		tx := &msg.Transaction{
			Type: t,
			Addr: req.Addr,
			Data: data,
		}
		go p.manage.AddTransaction(tx)

	} else if t > 4 && t < 6 {
		// external world tx

		pr := &msg.PendingRequest{
			Type: req.Type,
			Addr: req.Addr,
			Data: req.Data,
		}
		go p.oracle.AddEWTx(pr)
	}

	re := &client.ResEmpty{}
	return re, nil
}

func (p *Peer) sendRes(req *client.ResTx) (*client.ResEmpty, error) {
	return &client.ResEmpty{}, nil
}
