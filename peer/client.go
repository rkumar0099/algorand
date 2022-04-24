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

	} else if t >= 5 && t < 7 {
		// external world tx
		data, _ := req.Serialize()

		pr := &msg.PendingRequest{
			Type: t,
			Addr: req.Addr,
			Data: data,
		}

		go p.oracle.AddEWTx(pr)
	}

	re := &client.ResEmpty{}
	return re, nil
}

func (p *Peer) sendRes(req *client.ResTx) (*client.ResEmpty, error) {
	return &client.ResEmpty{}, nil
}
