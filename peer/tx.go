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
	log.Println("[Debug] [Tx] Received Final Contribution")
	p.finalContributions <- txSet
	return nil
}

func (p *Peer) executeTxSet(txSet *msg.ProposedTx, rootHash []byte, store kvstore.KVStore) (*mpt.Trie, []*msg.TxRes) {
	rootNode := mpt.HashNode(rootHash)
	st := mpt.New(&rootNode, store)
	res := p.execute(st, txSet.Txs, txSet.Epoch)
	log.Printf("[Debug] [Peer] [Execute Tx Set] Len: %d\n", len(res))
	return st, res
}

func (p *Peer) execute(st *mpt.Trie, txs []*msg.Transaction, epoch uint64) []*msg.TxRes {
	var responses []*msg.TxRes = make([]*msg.TxRes, 0)
	for _, tx := range txs {
		log.Printf("[Debug] [Peer] [Tx Type] Tx type: %d\n", tx.Type)
		switch tx.Type {
		case transaction.CREATE:
			res := p.createAccount(st, tx.Data, tx.Addr)
			res.TxHash = tx.Hash().Bytes()
			responses = append(responses, res)
		case transaction.LOGIN:
			res := p.login(st, tx.Data, tx.Addr)
			res.TxHash = tx.Hash().Bytes()
			responses = append(responses, res)

		case transaction.LOGOUT:

		case transaction.TOPUP:
			res := p.topUp(st, tx.Data, tx.Addr)
			res.TxHash = tx.Hash().Bytes()
			responses = append(responses, res)

		case transaction.TRNASFER:
			res := p.transfer(st, tx.Data, tx.Addr)
			res.TxHash = tx.Hash().Bytes()
			responses = append(responses, res)

		case transaction.OFFTOPUP:
			res := p.offTopUp(st, tx.From, cmn.Bytes2Uint(tx.Data))
			res.TxHash = tx.Hash().Bytes()
			responses = append(responses, res)

		default:
			log.Printf("Received invalid transaction")
		}
	}

	return responses
}

func (p *Peer) createAccount(st *mpt.Trie, data []byte, addr string) *msg.TxRes {
	response := &msg.TxRes{}
	response.SendRes = true
	req := &client.ReqTx{}
	res := &client.ResTx{}
	req.Deserialize(data)
	info := &client.Create{}
	info.Deserialize(req.Data)
	res.Addr = addr
	_, err := p.manage.GetClient(req.Pubkey)
	if err == nil {
		res.Status = false
		res.Msg = "Account already exists"
		response.Res, _ = res.Serialize()
		response.StoreData = false
		return response
	}

	user := &User{}
	user.Type = 2 // 1 = Blockchain Peer, 2 = Normal Client
	user.Username = info.Username
	user.PassHash = info.Password
	user.Pubkey = req.Pubkey
	user.Balance = 100
	user.Online = false
	user_data, _ := proto.Marshal(user)

	log.Println("[Debug] [Create User] Successfully created user")
	res.Status = true
	res.Msg = "Account Created Successfully"
	res.Pubkey = user.Pubkey

	response.Res, _ = res.Serialize()
	response.Data = user_data
	response.StoreData = true
	response.DataAddr = req.Pubkey

	st.Put(req.Pubkey, user_data)

	return response

}

func (p *Peer) login(st *mpt.Trie, data []byte, addr string) *msg.TxRes {
	response := &msg.TxRes{}
	response.SendRes = true
	response.StoreData = false
	req := &client.ReqTx{}
	res := &client.ResTx{}
	info := &client.LogIn{}
	req.Deserialize(data)
	info.Deserialize(req.Data)
	res.Addr = addr
	data, err := p.manage.GetClient(req.Pubkey)
	if err != nil {
		res.Status = false
		res.Msg = "User doesn't exists"
		response.Res, _ = res.Serialize()
		return response
	}
	user := &User{}
	proto.Unmarshal(data, user)
	if bytes.Equal(info.Password, user.PassHash) {
		// login successful
		user.Online = true
		user_data, _ := proto.Marshal(user)
		st.Put(req.Pubkey, user_data)
		response.StoreData = true
		response.Data = user_data
		response.DataAddr = req.Pubkey
		res.Status = true
		res.Msg = "Login successful"

	} else {
		// login unsuccessful
		res.Status = false
		res.Msg = "Invalid credentials"
	}
	response.Res, _ = res.Serialize()
	return response
}

func (p *Peer) logOut(st *mpt.Trie, data []byte, addr string) *msg.TxRes {
	response := &msg.TxRes{}
	req := &client.ReqTx{}
	req.Deserialize(data)
	res := &client.ResTx{}
	res.Addr = addr
	response.SendRes = true
	response.StoreData = false
	res.Status = true
	res.Msg = "Log out successful"
	return response
}

func (p *Peer) topUp(st *mpt.Trie, data []byte, addr string) *msg.TxRes {
	response := &msg.TxRes{}
	response.SendRes = true
	req := &client.ReqTx{}
	res := &client.ResTx{}
	info := &client.TopUp{}
	req.Deserialize(data)
	info.Deserialize(req.Data)
	res.Addr = addr
	data, err := p.manage.GetClient(req.Pubkey)
	if err != nil {
		res.Status = false
		res.Msg = "Invalid credentials"
		response.Res, _ = res.Serialize()
		response.StoreData = false
		return response
	}

	user := &User{}
	proto.Unmarshal(data, user)
	if user.Balance >= info.Amount {
		// req successful
		user.Balance += info.Amount
		user_data, _ := proto.Marshal(user)
		st.Put(req.Pubkey, user_data)
		response.StoreData = true
		response.Data = user_data
		response.DataAddr = req.Pubkey
		res.Status = true
		res.Msg = "Successful topup"
	} else {
		res.Status = false
		res.Msg = "UnSuccessful topup. Insufficient balance"
		response.StoreData = false
	}

	response.Res, _ = res.Serialize()
	response.Data, _ = proto.Marshal(user)
	return response
}

func (p *Peer) transfer(st *mpt.Trie, data []byte, addr string) *msg.TxRes {
	// check status of user. true if online false otherwise
	response := &msg.TxRes{}
	response.SendRes = true
	req := &client.ReqTx{}
	res := &client.ResTx{}
	info := &client.Transfer{}

	req.Deserialize(data)
	info.Deserialize(req.Data)
	res.Addr = addr
	data, err := p.manage.GetClient(req.Pubkey)
	if err != nil {
		res.Status = false
		res.Msg = "Invalid credentials"
		response.Res, _ = res.Serialize()
		response.StoreData = false
		return response
	}
	user := &User{}
	proto.Unmarshal(data, user)
	data, err = p.manage.GetClient(info.To)
	if err == nil && !bytes.Equal(info.To, info.From) && user.Balance >= info.Amount {
		// req success
		user.Balance -= info.Amount
		user_data, _ := proto.Marshal(user)
		st.Put(req.Pubkey, user_data)
		response.StoreData = true
		response.Data = user_data
		response.DataAddr = req.Pubkey
		p.TopupTransaction(info.To, info.Amount)
		response.Data, _ = proto.Marshal(user)
		res.Status = true
		res.Msg = "Successful transfer"

	} else {
		res.Status = false
		res.Msg = "Unsuccessful req"
		response.StoreData = false
	}

	response.Res, _ = res.Serialize()
	response.Data, _ = proto.Marshal(user)
	return response
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

func (p *Peer) TopupTransaction(addr []byte, value uint64) {
	tx := &msg.Transaction{
		From:  addr,
		To:    nil,
		Nonce: rand.Uint64(),
		Type:  transaction.OFFTOPUP,
		Data:  cmn.Uint2Bytes(value),
	}

	tx.Sign(p.privkey)
	p.manage.AddTransaction(tx)
}

func (p *Peer) offTopUp(st *mpt.Trie, addr []byte, amount uint64) *msg.TxRes {
	// handle it correctly
	response := &msg.TxRes{}
	response.SendRes = false
	data, err := p.manage.GetClient(addr)
	if err != nil {
		// user does not exists
		response.StoreData = false
		response.Type = transaction.OFFTOPUP
		return response
	}
	user := &User{}
	proto.Unmarshal(data, user)
	user.Balance += amount
	user_data, _ := proto.Marshal(user)
	st.Put(addr, user_data)
	response.Data, _ = proto.Marshal(user)
	response.StoreData = true
	return response
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
