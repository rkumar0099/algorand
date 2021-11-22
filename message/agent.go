package message

import (
	"github.com/golang/protobuf/proto"
	cmn "github.com/rkumar0099/algorand/common"
)

type MsgAgent struct {
	msgChan  chan []byte
	handlers map[int]func([]byte)
}

func NewMsgAgent(msgChan chan []byte) *MsgAgent {
	return &MsgAgent{
		msgChan:  msgChan,
		handlers: make(map[int]func([]byte)),
	}
}

func (agent *MsgAgent) Register(msgType int, handler func([]byte)) {
	agent.handlers[msgType] = handler
}

func (agent *MsgAgent) Handle() {
	for {
		msgBytes := <-agent.msgChan
		msg := Msg{}
		err := proto.Unmarshal(msgBytes, &msg)
		if err != nil {
			cmn.Log.Errorf("[MsgAgent] Cannot deserialize message: %s", err.Error())
			continue
		}
		if handler, ok := agent.handlers[int(msg.Type)]; ok {
			handler(msg.Data)
		} else {
			cmn.Log.Errorf("[MsgAgent] Invalid message type: %d", msg.Type)
		}
	}
}
