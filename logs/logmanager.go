package logs

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/oracle"
)

const bufferCap = 100
const txBufferCap = 40

var (
	txcount      int    = 0
	finaltxcount int    = 0
	blkcount     int    = 0
	processcount int    = 0
	startRound   uint64 = 0
	endRound     uint64 = 0
)

type LogManager struct {
	txLock      *sync.Mutex
	blkLock     *sync.Mutex
	processLock *sync.Mutex
	balanceLock *sync.Mutex
	buffer      chan *message.Msg
	txBuffer    chan cmn.Hash
	finalBlk    map[uint64]cmn.Hash
	balances    map[string]uint64
}

func New() *LogManager {
	lm := &LogManager{
		txLock:      &sync.Mutex{},
		blkLock:     &sync.Mutex{},
		processLock: &sync.Mutex{},
		balanceLock: &sync.Mutex{},
		buffer:      make(chan *message.Msg, bufferCap),
		txBuffer:    make(chan cmn.Hash, txBufferCap),
		finalBlk:    make(map[uint64]cmn.Hash),
		balances:    make(map[string]uint64),
	}
	lm.init()
	go lm.writeBalances()
	return lm
}

func (lm *LogManager) init() {
	os.Remove("../logs/process.txt")
	os.Create("../logs/process.txt")

	os.Remove("../logs/block.txt")
	os.Create("../logs/block.txt")

	os.Remove("../logs/balance.txt")
	os.Create("../logs/balance.txt")

	//os.Remove("../logs/oracle.txt")
	//os.Create("../logs/oracle.txt")

	os.Remove("../logs/txs.txt")
	os.Create("../logs/txs.txt")

	//os.Remove("../logs/externalData.txt")
	//os.Create("../logs/externalData.txt")

	os.Remove("../logs/peers.txt")
	os.Create("../logs/peers.txt")

}

func (lm *LogManager) AddProcessLog(msg *message.Msg) {
	lm.processLock.Lock()
	defer lm.processLock.Unlock()
	if len(lm.buffer) < bufferCap {
		lm.buffer <- msg
		return
	}
	newChan := make(chan *message.Msg, bufferCap)
	for len(lm.buffer) > 0 {
		newChan <- <-lm.buffer
	}
	lm.buffer <- msg
	go lm.writeProcessLog(newChan)
}

func (lm *LogManager) writeProcessLog(logs chan *message.Msg) {
	latestLog := ""
	for len(logs) > 0 {
		msg := <-logs
		pid := msg.PID
		switch msg.Type {
		case message.BLOCK:
			blk := &message.Block{}
			if err := blk.Deserialize(msg.Data); err != nil {
				latestLog += time.Now().String() + " " + pid + ": Invalid block data: " + err.Error() + "\n"
			} else {
				latestLog += time.Now().String() + " : send block " + blk.Hash().Hex() + "\n"
			}
			//log.Println(latestLog)
		case message.BLOCK_PROPOSAL:
			bp := &message.Proposal{}
			if err := bp.Deserialize(msg.Data); err != nil {
				latestLog += "[" + time.Now().String() + "] " + ": Invalid block proposal: " + err.Error() + "\n"
			} else {
				latestLog += "[" + time.Now().String() + "] " + fmt.Sprintf("[%s] propose block for round %d\n", pid, bp.Round)
				//latestLog += pid + ": propose block " + common.BytesToHash(bp.Hash).Hex() + "\n"
			}
			//log.Println(latestLog)
		case message.FORK_PROPOSAL:
			bp := &message.Proposal{}
			if err := bp.Deserialize(msg.Data); err != nil {
				latestLog += "[" + time.Now().String() + "] " + string(pid) + ": Invalid fork proposal: " + err.Error() + "\n"
			} else {
				latestLog += "[" + time.Now().String() + "] " + string(pid) + ": propose fork " + common.BytesToHash(bp.Hash).Hex() + "\n"
			}

		case message.VOTE:
			vote := &message.VoteMessage{}
			if err := vote.Deserialize(msg.Data); err != nil {
				latestLog += "[" + time.Now().String() + "] " + string(pid) + ": Invalid vote " + err.Error() + "\n"
			} else {
				latestLog += "[" + time.Now().String() + "] " + fmt.Sprintf("%d: vote in %d:%d %s\n", pid, vote.Round, 5, common.BytesToHash(vote.Hash).Hex())
			}
		}
	}
	f, err := os.OpenFile("../logs/process.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	f.WriteString(latestLog)
	f.Close()
}

func (lm *LogManager) AddTxLog(hash cmn.Hash) {
	lm.txLock.Lock()
	defer lm.txLock.Unlock()
	if len(lm.txBuffer) < txBufferCap {
		lm.txBuffer <- hash
		return
	}
	newChan := make(chan cmn.Hash, txBufferCap)
	for len(lm.txBuffer) > 0 {
		newChan <- <-lm.txBuffer
	}
	lm.txBuffer <- hash
	go lm.writeTxLog(newChan)
}

func (lm *LogManager) writeTxLog(txs chan cmn.Hash) {
	latestLog := ""
	for len(txs) > 0 {
		hash := <-txs
		txcount += 1
		latestLog += "Transaction " + strconv.Itoa(txcount) + " proposed. Hash: " + hash.Hex() + "\n"
	}
	f, err := os.OpenFile("../logs/txs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	f.WriteString(latestLog)
	f.Close()
}

func (lm *LogManager) AddFinalTxLog(txs []*message.Transaction) {
	lm.txLock.Lock()
	defer lm.txLock.Unlock()
	latestLog := ""
	for _, tx := range txs {
		finaltxcount += 1
		latestLog += "Transaction " + strconv.Itoa(finaltxcount) + " finalized, Hash: " + tx.Hash().Hex() + "\n"
	}
	f, _ := os.OpenFile("../logs/txs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	f.WriteString(latestLog)
	f.Close()
}

func (lm *LogManager) AddFinalBlk(blkHash cmn.Hash, round uint64) {
	lm.blkLock.Lock()
	defer lm.blkLock.Unlock()
	if _, ok := lm.finalBlk[round]; !ok {
		lm.finalBlk[round] = blkHash
		blkcount += 1
		if blkcount%10 == 0 {
			go lm.writeBlkLog(lm.finalBlk, startRound, endRound)
			startRound = uint64(blkcount)
		}
		endRound += 1
	}
}

func (lm *LogManager) writeBlkLog(finalBlk map[uint64]cmn.Hash, sr uint64, er uint64) {
	latestLog := ""
	for sr <= er {
		latestLog += "Block " + strconv.Itoa(int(sr)) + " added to blockchain, Hash: " + finalBlk[sr].Hex() + "\n"
		delete(finalBlk, sr)
		sr += 1
	}

	f, err := os.OpenFile("../logs/block.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	f.WriteString(latestLog)
	f.Close()
}

func (lm *LogManager) AddBalance(addr string, bal uint64) {
	lm.balanceLock.Lock()
	defer lm.balanceLock.Unlock()
	lm.balances[addr] = bal
}

func (lm *LogManager) writeBalances() {
	for {
		time.Sleep(10 * time.Second)
		latestLog := ""
		for addr, bal := range lm.balances {
			latestLog += time.Now().String() + ", " + addr + ", balance = " + strconv.Itoa(int(bal)) + " \n"
		}
		f, _ := os.OpenFile("../logs/balance.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		f.WriteString(latestLog)
		f.Close()
	}
}

func (lm *LogManager) AddOracleBlk(blk *oracle.FinalBlock) {

}
