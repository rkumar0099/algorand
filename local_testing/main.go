package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/gossip"
	"github.com/rkumar0099/algorand/logs"
	"github.com/rkumar0099/algorand/manage"
	"github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/oracle"
	"github.com/rkumar0099/algorand/params"
	"github.com/rkumar0099/algorand/peer"
	"github.com/urfave/cli"
)

func main() {
	app := initApp()
	err := app.Run(os.Args)
	if err != nil {
		cmn.Log.Fatal(err)
	}
}

func initApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Algorand"
	app.Version = "0.1"
	app.Author = "Ziliang"
	app.Usage = "Algorand simulation demo for different scenario."

	app.Commands = []cli.Command{
		{
			Name:    "regular",
			Aliases: []string{"r"},
			Usage:   "run regular Algorand algorithm",
			Action:  regularRun,
			Flags: []cli.Flag{
				cli.Uint64Flag{
					Name:  "num,n",
					Value: 50,
					Usage: "amount of users",
				},
				cli.Uint64Flag{
					Name:  "token,t",
					Value: 1000,
					Usage: "token balance per users",
				},
				cli.Uint64Flag{
					Name:  "malicious,m",
					Value: 0,
					Usage: "amount of malicious users. Malicious user will use default strategy.",
				},
				cli.IntFlag{
					Name:  "mtype,i",
					Value: 0,
					Usage: "malicious type: 0 Honest, 1 block proposal misbehaving; 2 vote empty block in BA*; 3 vote nothing",
				},
				cli.IntFlag{
					Name:  "latency,l",
					Value: 0,
					Usage: "max network latency(milliseconds). Each user will simulate a random latency between 0 and ${value}",
				},
			},
		},
	}

	return app
}

func regularRun(c *cli.Context) {
	params.UserAmount = c.Uint64("num")
	params.TokenPerUser = c.Uint64("token")
	params.Malicious = c.Uint64("malicious")
	params.NetworkLatency = c.Int("latency")

	var (
		nodes     []*peer.Peer
		i         = 0
		addrPeers [][]byte
		//stores    []kvstore.KVStore
		neighbors []gossip.NodeId
		num       = 0
	)

	for ; uint64(i) < params.UserAmount-params.Malicious; i++ {
		Id := fmt.Sprintf("127.0.0.1:%d", 8000+i)
		neighbors = append(neighbors, gossip.NewNodeId(Id))
		node := peer.New(Id, params.Honest)
		num += 1
		addrPeers = append(addrPeers, node.Address().Bytes())
		nodes = append(nodes, node)
	}

	for _, p := range nodes {
		go p.Start(neighbors, addrPeers)
	}
	time.Sleep(5 * time.Second)

	lm := logs.New()
	m := manage.New(neighbors, addrPeers, lm)
	oracle := oracle.New()
	go oracle.Run()

	for _, p := range nodes {
		p.AddManage(m, lm, oracle)
		go p.Run()
	}
	time.Sleep(5 * time.Second)

	go m.Run()
	go proposeEWTxs(nodes)
	go proposeTopUpTransactions(nodes)
	go proposeTransferTransactions(nodes)

	time.Sleep(1 * time.Minute)
	//showBalances(nodes)

	/*
		count := 50
		for count > 0 {
			go proposeTransactions(nodes)
			log.Println("Proposing Transactions")
			time.Sleep(5 * time.Second)
			count -= 1
		}
	*/

}

func proposeEWTxs(peers []*peer.Peer) {
	for {
		time.Sleep(5 * time.Second)
		for i := 0; i < 10; i++ {
			ind := rand.Intn(50)
			peers[ind].ExternalWorldTransaction("", oracle.CURRENCY_EXCHANGE)
		}
	}
}

func proposeTopUpTransactions(peers []*peer.Peer) {
	for {
		time.Sleep(5 * time.Second)
		for i := 0; i < 5; i++ {
			ind := rand.Intn(50)
			peers[ind].TopupTransaction(uint64(10))
		}
	}
}

func proposeTransferTransactions(peers []*peer.Peer) {
	for {
		time.Sleep(5 * time.Second)
		for i := 0; i < 5; i++ {
			ind := rand.Intn(50)
			to := rand.Intn(50)
			for to == ind {
				to = rand.Intn(50)
			}
			peers[ind].TransferTransaction(uint64(10), peers[to].Address())
		}
	}
}

func showBalances(peers []*peer.Peer) {
	for _, p := range peers {
		log.Printf("Balance of peer %s is %d\n", p.Id.String(), p.GetBalance())
	}
}

func sendData(node *gossip.Node) {
	for {
		time.Sleep(1 * time.Second)
		blk := &message.Block{
			Round: uint64(1),
		}
		data, _ := blk.Serialize()
		msg := &message.Msg{
			PID:  "123",
			Type: message.BLOCK,
			Data: data,
		}
		msgBytes, _ := proto.Marshal(msg)
		node.Gossip(msgBytes)
	}
}
