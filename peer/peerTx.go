package peer

import (
	"bytes"
	"log"
	"math/rand"

	"github.com/golang/protobuf/proto"
	cmn "github.com/rkumar0099/algorand/common"
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
	st := p.executeTxSet(txSet, p.lastState, p.transactionStorage)
	p.lastState = st.RootHash()
	st.Commit()
	p.finalContributions <- txSet
	//st := p.executeTxSet(txSet)
	//p.lastState = st.RootHash()
	//st.Commit()
	//log.Println("Final contributions are: ", len(p.finalContributions))
	return nil
}

func (p *Peer) handleContribution(data []byte) ([]byte, error) {
	txSet := &msg.ProposedTx{}
	err := proto.Unmarshal(data, txSet)
	if err != nil {
		return nil, err
	}

	//log.Println("Got contribution of size ", len(txSet.Txs))
	st := p.executeTxSet(txSet, p.lastState, p.transactionStorage)
	res := &msg.StateHash{
		Epoch:     txSet.Epoch,
		StateHash: st.RootHash(),
	}
	resBytes, _ := proto.Marshal(res)
	return resBytes, nil
}

func (p *Peer) executeTxSet(txSet *msg.ProposedTx, rootHash []byte, store kvstore.KVStore) *mpt.Trie {
	/*blk := p.lastBlock()
	sh := blk.StateHash */
	rootNode := mpt.HashNode(rootHash)
	st := mpt.New(&rootNode, store)
	p.execute(st, txSet.Txs, txSet.Epoch)
	return st
}

func (p *Peer) execute(st *mpt.Trie, txs []*msg.Transaction, epoch uint64) {
	for _, tx := range txs {
		switch tx.Type {
		case transaction.TOPUP:
			transaction.Topup(st, tx, cmn.Bytes2Uint(tx.Data))
		case transaction.TRNASFER:
			transaction.Transfer(st, tx, cmn.Bytes2Uint(tx.Data))
		default:
			log.Printf("Received invalid transaction")
		}
	}
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
	p.chain.Genesis.StateHash = st.RootHash()
	st.Commit()
	p.startState = st.RootHash()
	p.lastState = st.RootHash()
}

func (p *Peer) TopupTransaction(value uint64) {
	//txTopup := transaction.TxTopup{
	//	Value: value,
	//}
	//data, _ := proto.Marshal(&txTopup)
	tx := &msg.Transaction{
		From:  p.pubkey.Address().Bytes(),
		To:    nil,
		Nonce: rand.Uint64(),
		Type:  transaction.TOPUP,
		Data:  cmn.Uint2Bytes(value),
	}
	p.manage.AddTransaction(tx)

	//tx.Sign(p.privkey)
	//txMsg, _ := tx.Serialize()
	//p.gossip(msg.TRANSACTION, txMsg)
}

func (p *Peer) GetBalance() uint64 {
	blk := p.lastBlock()
	hn := mpt.HashNode(blk.StateHash)
	st := mpt.New(&hn, p.storage)
	addr := bytes.Join([][]byte{
		p.Address().Bytes(),
		[]byte("value"),
	}, nil)
	val, _ := st.Get(addr)
	return cmn.Bytes2Uint(val)
}

func (p *Peer) GetLastState() []byte {
	return p.lastBlock().StateHash
}

func (p *Peer) GetStartState() []byte {
	return p.chain.Genesis.StateHash
}
