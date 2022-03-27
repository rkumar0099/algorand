package peer

import (
	"bytes"
	"log"
	"math/rand"

	"github.com/golang/protobuf/proto"
	"github.com/rkumar0099/algorand/client"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/message"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/mpt/kvstore"
	"github.com/rkumar0099/algorand/mpt/mpt"
	"github.com/rkumar0099/algorand/transaction"
)

func (p *Peer) handleFinalContribution(data []byte) error {
	//log.Println("Got final contribution, len: ", len(p.finalContributions))
	txSet := &msg.ProposedTx{}
	err := proto.Unmarshal(data, txSet)
	if err != nil {
		return err
	}

	p.finalContributions <- txSet
	return nil
}

func (p *Peer) handleContribution(data []byte) ([]byte, error) {
	txSet := &msg.ProposedTx{}
	err := proto.Unmarshal(data, txSet)
	if err != nil {
		return nil, err
	}

	//log.Println("Got contribution of size ", len(txSet.Txs))
	st, _ := p.executeTxSet(txSet, p.lastState, p.permanentTxStorage)
	res := &msg.StateHash{
		Epoch:     txSet.Epoch,
		StateHash: st.RootHash(),
	}
	resBytes, _ := proto.Marshal(res)
	return resBytes, nil
}

func (p *Peer) executeTxSet(txSet *msg.ProposedTx, rootHash []byte, store kvstore.KVStore) (*mpt.Trie, []*msg.TxRes) {
	rootNode := mpt.HashNode(rootHash)
	st := mpt.New(&rootNode, store)
	res := p.execute(st, txSet.Txs, txSet.Epoch)
	return st, res
}

func (p *Peer) execute(st *mpt.Trie, txs []*msg.Transaction, epoch uint64) []*msg.TxRes {
	var responses []*msg.TxRes
	for _, tx := range txs {
		switch tx.Type {
		case transaction.CREATE:
			r := p.createAccount(st, tx.Data, tx.Addr)
			res := &msg.TxRes{}
			res.TxHash = tx.Hash().Bytes()
			d, _ := r.Serialize()
			res.Data = d
			responses = append(responses, res)
		case transaction.LOGIN:
			p.logIn(st, tx.Data)
		case transaction.LOGOUT:
			p.logOut(st, tx.Data)
		case transaction.TOPUP:
			p.topUp(st, tx.Data)
			transaction.Topup(st, tx, cmn.Bytes2Uint(tx.Data))
		case transaction.TRNASFER:
			p.transfer(st, tx.Data)
			transaction.Transfer(st, tx, cmn.Bytes2Uint(tx.Data))
		default:
			log.Printf("Received invalid transaction")
		}
	}
	return responses
}

func (p *Peer) createAccount(st *mpt.Trie, data []byte, addr string) *client.ResTx {
	info := &client.Create{}
	info.Deserialize(data)
	user := &User{}
	user.Type = 2 // 1 = Blockchain Peer, 2 = Normal Client
	user.Username = info.Username
	user.PassHash = info.Password
	user.Pubkey = info.Pubkey
	user.Balance = 100
	user_data, _ := proto.Marshal(user)
	err := st.Put(info.Pubkey, user_data)
	res := &client.ResTx{}
	if err == nil {
		log.Println("[Debug] [Create User] Successfully created user")
		res.Status = true
		res.Msg = "Account Created Successfully"
		res.Balance = user.Balance
		res.Pubkey = user.Pubkey
		res.Addr = addr
	} else {
		res.Status = false
		res.Msg = "[Error] Can't create account"
		res.Balance = 0
		res.Pubkey = info.Pubkey
		res.Addr = addr
	}
	return res

}

func (p *Peer) logIn() {

}

func (p *Peer) logOut() {

}

func (p *Peer) topUp() {

}

func (p *Peer) transfer() {

}

func (p *Peer) initialize(addr [][]byte, store kvstore.KVStore) {
	st := mpt.New(nil, store)
	for _, val := range addr {
		stateAddr := bytes.Join([][]byte{
			val,
			[]byte("value"),
		}, nil)
		st.Put(stateAddr, cmn.Uint2Bytes(0))
	}
	//log.Println("State trie is init")
	p.chain.Genesis.StateHash = st.RootHash()
	p.startState = st.RootHash()
	p.lastState = st.RootHash()
	st.Commit()
}

func (p *Peer) TopupTransaction(value uint64) {
	tx := &msg.Transaction{
		From:  p.pubkey.Address().Bytes(),
		To:    nil,
		Nonce: rand.Uint64(),
		Type:  transaction.TOPUP,
		Data:  cmn.Uint2Bytes(value),
	}

	tx.Sign(p.privkey)
	p.manage.AddTransaction(tx)
}

func (p *Peer) TransferTransaction(value uint64, to cmn.Address) {
	tx := &msg.Transaction{
		From:  p.pubkey.Address().Bytes(),
		To:    to.Bytes(),
		Nonce: rand.Uint64(),
		Type:  transaction.TRNASFER,
		Data:  cmn.Uint2Bytes(value),
	}
	tx.Sign(p.privkey)
	p.manage.AddTransaction(tx)
}

func (p *Peer) GetBalance() uint64 {
	//blk := p.lastBlock()
	hn := mpt.HashNode(p.lastState)
	st := mpt.New(&hn, p.permanentTxStorage)
	addr := bytes.Join([][]byte{
		p.Address().Bytes(),
		[]byte("value"),
	}, nil)
	val, _ := st.Get(addr)
	return cmn.Bytes2Uint(val)
}

func (p *Peer) ExternalWorldTransaction(url string, reqType int) {
	pr := &message.PendingRequest{
		URL:  url,
		Type: uint64(reqType),
	}
	p.oracle.AddEWTx(pr)
}
