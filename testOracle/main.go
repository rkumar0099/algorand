package main

import (
	"bufio"
	"os"
	"strconv"

	cmn "github.com/rkumar0099/algorand/common"
)

type Blk struct {
	Data []byte `json:"data"`
	Name string `json:"name"`
}

type BTC_TO_CUR struct {
	Time string `json:"time"`
	//Base  string  `json:"asset_id_base"`
	//Quote string  `json:"asset_id_quote"`
	Rate float64 `json:"rate"`
}

type Res struct {
	data []byte
}

var (
	assets = "https://rest.coinapi.io/v1/exchangerate/BTC/USD?apikey=34CD7A19-89E3-4004-A2D3-9B330C62D8FB"
)

func main() {

	//testOracle()
	runOracle()
	//print(val)

}

func runOracle() (float64, error) {
	f, err := os.OpenFile("../pricefeeds/bit.txt", os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		return 0.0, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	var float_val float64 = 0.0
	for scanner.Scan() {
		val := scanner.Text()
		float_val, _ = strconv.ParseFloat(val, 64)
		//log.Println(val)
	}
	var b []byte
	b = cmn.Float2Bytes(float_val)
	print(b)
	//print(cmn.Bytes2Float(b))
	return float_val, nil
}

/*

func testOracle() {
	//res := &oracle.CurrencyExchange{}
	buffer := make([]byte, 1024)
	f, _ := os.OpenFile("./data.txt", os.O_RDONLY|os.O_CREATE, 0644)
	n, _ := f.Read(buffer)

	response := &oracle.Response{
		Data: buffer[0:n],
		Type: 2,
		Id:   4,
	}

	data, _ := proto.Marshal(response)
	log.Println(data)
}
*/

/*

func testFile() {
	var blocks map[int]int = make(map[int]int)
	os.Remove("./data.bin")
	f, err := os.OpenFile("./data.bin", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		blk := Blk{
			Data: []byte("Rabindar"),
			Name: "Kumar",
		}
		data, _ := json.Marshal(blk)
		blocks[0] = len(data)
		f.Write(data)
	}
	f.Close()
	f, _ = os.Open("./data.bin")
	readBlk := &Blk{}
	buffer := make([]byte, blocks[0])
	f.Read(buffer)
	json.Unmarshal(buffer, readBlk)
	fmt.Println(readBlk)

}

func testAddBlk() {

	/*
		for i := 1; i <= 100000; i++ {
			blk := &message.Block{
				Round: uint64(i),
			}
			o.AddBlk(blk)
		}
		time.Sleep(20 * time.Second)
		for i := 1; i <= 100; i++ {
			round := uint64(rand.Intn(10000))
			log.Println(o.GetBlkByRound(round))
		}
	}


func testProposeOraclePeer() {
	o := oracle.New()
	for i := uint64(1); i <= uint64(5); i++ {
		seed := o.SortitionSeed(i)
		vrf, proof, u := o.Sortition(seed, oracle.Role(params.OraclePeer, i, params.ORACLE), params.ExpectedOraclePeers, uint64(10000))
		if u > 0 {
			//log.Println(u)
			op := &message.OraclePeerProposal{
				Round:  i,
				Vrf:    vrf,
				Proof:  proof,
				PubKey: o.Address(),
			}
			o.AddOPP(op)
		}
	}

	for i := 0; i < 2; i++ {
		pr := &message.PendingRequest{
			URL: assets,
			Id:  uint64(i),
		}
		o.AddEWTx(pr)
	}
	res := o.RunOracle()
	time.Sleep(2 * time.Minute)
	fmt.Println(res)
	/*
		peers, epoch := o.GetPeers()
		log.Println("Number of peers for given epoch are ", len(peers))
		nonce := uint64(0)
		for {
			blkChan := make(chan *oracle.FinalBlock, 1)
			nonce += 1
			for _, p := range peers {
				go p.Propose(nonce, o.SortitionSeed, blkChan, epoch)
			}
			time.Sleep(2 * time.Second)
			if len(blkChan) == 1 {
				log.Println("Found proposer")
				break
			}
			///log.Println(len(blkChan))
		}
	}

*/
