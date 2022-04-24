package oracle

import (
	"sync"

	"github.com/rkumar0099/algorand/crypto"
)

type response struct {
	val_bytes []byte
	val_float float64
}

func newOraclePeer(epoch uint64) *OraclePeer {
	op := &OraclePeer{
		results: make(map[uint64][]byte),
		epoch:   epoch,
		lock:    &sync.Mutex{},
	}
	op.pubkey, op.privkey, _ = crypto.NewKeyPair()
	return op
}

/*

func (op *OraclePeer) commit(epoch uint64, txs []*message.PendingRequest, res map[cmn.Hash][]byte) {
	for _, tx := range txs {
		req := &client.ReqTx{}
		req.Deserialize(t.Data)
		pf := &client.PriceFeed{}
		pf.Deserialize(req.Data)
		val_bytes, val, err := op.fetchData(pf.Type)
		res := &Response{}
		res.val_bytes = val_bytes
		res.val_float = val
		op.responses[tx.Id] = res

		op.data = append(op.data, val_bytes)
	}
}


func (op *OraclePeer) processTx(tx *msg.PendingRequest) (float64, error) {
	req := &client.ReqTx{}
	req.Deserialize(tx.Data)
	pf := &client.PriceFeed{}
	pf.Deserialize(req.Data)
	return op.FetchData(pf.Type)
}

func (op *OraclePeer) FetchData(pricefeed uint64) (float64, error) {
	switch pricefeed {
	case 1:
		f, err := os.OpenFile("../pricefeeds/eth.txt", os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return 0.0, err
		}
		buffer := make([]byte, 128)
		f.Read(buffer)
		val, _ := strconv.ParseFloat(string(buffer), 2)
		return val, nil

	case 2:
		f, err := os.OpenFile("../pricefeeds/bit.txt", os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return 0.0, err
		}
		buffer := make([]byte, 128)
		f.Read(buffer)
		val, _ := strconv.ParseFloat(string(buffer), 2)
		return val, nil

	case 3:
		f, err := os.OpenFile("../pricefeeds/algo.txt", os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return 0.0, err
		}
		buffer := make([]byte, 128)
		f.Read(buffer)
		val, _ := strconv.ParseFloat(string(buffer), 2)
		return val, nil

	default:
		log.Printf("[Oracle] Invalid Pricefeed req\n")
		return 0.0, nil

	}

}

func (op *OraclePeer) fetchFile(url string, Id uint64, reqType uint64) {
	f, err := os.OpenFile("../oracle/data.txt", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	buffer := make([]byte, 1024)

	n, _ := f.Read(buffer)
	f.Close()
	response := &Response{
		Data: buffer[0:n],
		Type: int32(reqType),
	}
	data, _ := proto.Marshal(response)
	dataCommit := bytes.Join([][]byte{
		data,
		cmn.Uint2Bytes(Id),
	}, nil)

	op.data = append(op.data, dataCommit) // make data[][] array and make mpt from this array after all txs done
	op.results[Id] = data

}

func (op *OraclePeer) fetchURL(url string, Id uint64, reqType uint64) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("error accessing data")
			return
		}

		response := op.convert(respBytes, int(reqType))
		response.Id = Id
		data, _ := proto.Marshal(response)
		dataCommit := bytes.Join([][]byte{
			data,
			cmn.Uint2Bytes(Id),
		}, nil)

		op.data = append(op.data, dataCommit) // make data[][] array and make mpt from this array after all txs done
		op.results[Id] = data

	}

}

func (op *OraclePeer) convert(data []byte, reqType int) *Response {
	response := &Response{}
	switch reqType {
	case CURRENCY_EXCHANGE:
		res := &CurrencyExchange{}
		json.Unmarshal(data, res)
		resByte, _ := json.Marshal(res)
		response.Data = resByte
		response.Type = int32(reqType)

	case FLIGHT:
		// write your logic here

	default:
		log.Println("Invalid request")
	}
	return response
}

func (op *OraclePeer) reveal(epoch uint64, jobIds []uint64, res map[cmn.Hash]*Reveal) {
	//log.Println(epoch, jobIds)

	reveal := &Reveal{
		tree:    op.tree,
		addr:    op.pubkey.Bytes(),
		results: make(map[uint64][]byte),
	}

	for _, id := range jobIds {
		reveal.results[id] = op.results[id]
		delete(op.results, id)
	}

	res[cmn.BytesToHash(op.pubkey.Bytes())] = reveal

}

// role returns the role bytes from current round and step
func opRole(pk []byte, iden string, round uint64, event string) []byte {
	return bytes.Join([][]byte{
		[]byte(iden),
		common.Uint2Bytes(round),
		[]byte(event),
	}, nil)
}

// sortition runs cryptographic selection procedure and returns vrf,proof and amount of selected sub-users.
func (op *OraclePeer) sortition(seed, role []byte, expectedNum int, weight uint64) (vrf, proof []byte, selected int) {
	vrf, proof, _ = op.privkey.Evaluate(constructSeed(seed, role))
	selected = cmn.SubUsers(expectedNum, weight, vrf)
	return
}

type SortitionSeed func(uint64) []byte

func (op *OraclePeer) proposeBlock(nonce uint64, ss SortitionSeed, blkChan chan *FinalBlock,
	epoch uint64, res map[uint64][][]byte) {
	//salt, _ := op.privkey.Sign(cmn.Uint2Bytes(epoch))
	role := role(params.OracleBlockProposer, epoch, params.ORACLE_BLK_PROPOSAL)
	seed := ss(epoch)
	finalSeed := bytes.Join([][]byte{
		seed,
		cmn.Uint2Bytes(nonce),
	}, nil)
	proof, vrf, subusers := op.sortition(finalSeed, role, 20, uint64(10000))
	if subusers > 0 {
		//log.Println("Selected to propose block", op.pubkey.Bytes())
		blk := &FinalBlock{
			Proof:  proof,
			VRF:    vrf,
			Pubkey: op.pubkey.Bytes(),
			Epoch:  epoch,
		}

		blk.Results = make(map[uint64][]byte)
		for Id, results := range res {
			ind := int(len(results) / 2)
			blk.Results[Id] = results[ind]
		}

		if len(blkChan) < 1 {
			blkChan <- blk
		}
	}
}

*/
