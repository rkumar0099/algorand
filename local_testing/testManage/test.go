package main

import (
	"fmt"
	"time"

	"github.com/rkumar0099/algorand/logs"
	"github.com/rkumar0099/algorand/manage"
	msg "github.com/rkumar0099/algorand/message"
)

func main() {
	lm := logs.New()
	m := manage.New(nil, nil, lm)
	blk := &msg.Block{}
	data, _ := blk.Serialize()
	fmt.Println(m, data)
	for i := 0; i < 1000; i++ {
		go m.AddBlk(blk.Hash(), data)
	}
	time.Sleep(10 * time.Second)
}
