package logs

import (
	"fmt"
	"testing"

	"github.com/rkumar0099/algorand/gossip"
	"github.com/rkumar0099/algorand/message"
)

func TestLogManager(t *testing.T) {
	n := gossip.New(gossip.NewNodeId("127.0.0.1:9001"), "logs")
	lm := New(n)
	for i := 0; i < 49; i++ {
		msg := &message.Msg{
			PID:  "R",
			Type: message.BLOCK,
			Data: []byte("a"),
		}
		res := lm.AddProcessLog(msg)
		fmt.Println(res)
	}
}
