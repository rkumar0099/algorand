package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/oracle"
	"github.com/rkumar0099/algorand/params"
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

var (
	assets = "https://rest.coinapi.io/v1/exchangerate/BTC/CAD?apikey=34CD7A19-89E3-4004-A2D3-9B330C62D8FB"
)

func main() {
	/*
		resp, _ := http.Get(assets)
		respBytes, _ := ioutil.ReadAll(resp.Body)
		res := &BTC_TO_CUR{}

		json.Unmarshal(respBytes, res)
		os.Remove("./data.txt")
		f, err := os.OpenFile("./data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			data, _ := json.Marshal(res)
			f.Write(data)
		}
	*/
	//testOracle()
	//testFile()
	//testAddBlk()
	res, _ := http.Get(assets)
	respBytes, _ := ioutil.ReadAll(res.Body)
	log.Println(string(respBytes))
	//testProposeOraclePeer()
}

func testOracle() {
	o := oracle.New()
	for i := 0; i < 5; i++ {
		pr := &message.PendingRequest{
			URL: assets,
			Id:  uint64(i),
		}
		o.AddEWTx(pr)
	}
	for i := 0; i < 2; i++ {
		opp := &message.OraclePeerProposal{}
		o.AddOPP(opp)
	}
	o.RunOracle()
	time.Sleep(1 * time.Minute)
}

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
	*/
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
	*/
}
