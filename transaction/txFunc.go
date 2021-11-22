package transaction

import (
	"bytes"

	"github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/mpt/mpt"
)

func Topup(state *mpt.Trie, meta *message.Transaction, value uint64) bool {
	stateAddr := bytes.Join([][]byte{
		meta.From,
		[]byte("value"),
	}, nil)

	//fmt.Println("The stateAddr in Topup function is ", stateAddr)
	data, err := state.Get(stateAddr)
	var oriValue uint64 = 0
	if err != nil {
		common.Log.Errorf("[Topup] cannot get state: %s", err.Error())
	}
	if err == nil {
		oriValue = common.Bytes2Uint(data)
	}
	newData := common.Uint2Bytes(oriValue + value)
	err = state.Put(stateAddr, newData)
	if err != nil {
		common.Log.Errorf("[Topup] cannot put state: %s", err.Error())
	}
	return true
}

func Transfer(state *mpt.Trie, meta *message.Transaction, value uint64) bool {
	fromValueAddr := bytes.Join([][]byte{
		meta.From,
		[]byte("value"),
	}, nil)
	data, err := state.Get(fromValueAddr)
	var fromValue uint64 = 0
	if err == nil {
		fromValue = common.Bytes2Uint(data)
	}
	if value > fromValue {
		return false
	}
	fromValue -= value
	toValueAddr := bytes.Join([][]byte{
		meta.To,
		[]byte("value"),
	}, nil)
	data, err = state.Get(toValueAddr)
	var toValue uint64 = 0
	if err == nil {
		toValue = common.Bytes2Uint(data)
	}
	toValue += value
	state.Put(fromValueAddr, common.Uint2Bytes(fromValue))
	state.Put(toValueAddr, common.Uint2Bytes(toValue))
	return true
}
