package blockchain

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/rkumar0099/algorand/common"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/pool"
)

type Blockchain struct {
	Last            *msg.Block
	Genesis         *msg.Block
	db              *leveldb.DB
	blockCache      *pool.RingBuffer
	blockCacheIndex map[common.Hash]*msg.Block
	cacheLock       *sync.Mutex
}

const cacheSize = 5

func NewBlockchain(persistBlkStorage *leveldb.DB) *Blockchain {
	bc := &Blockchain{
		db:              persistBlkStorage,
		blockCache:      pool.NewRingBuffer(cacheSize),
		blockCacheIndex: make(map[common.Hash]*msg.Block),
		cacheLock:       &sync.Mutex{},
	}
	emptyHash := common.Sha256([]byte{})
	bc.Genesis = &msg.Block{
		Round:      0,
		Seed:       emptyHash.Bytes(),
		ParentHash: emptyHash.Bytes(),
		Author:     common.HashToAddr(emptyHash).Bytes(),
		StateHash:  nil,
	}
	bc.Add(bc.Genesis)
	return bc
}

func (bc *Blockchain) Get(hash common.Hash) (*msg.Block, error) {
	bc.cacheLock.Lock()
	blk, ok := bc.blockCacheIndex[hash]
	bc.cacheLock.Unlock()
	if ok {
		return blk, nil
	}
	data, err := bc.db.Get(hash.Bytes(), nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Chain] block %s not found: %s", hash.Hex(), err.Error()))
	}
	blk = &msg.Block{}
	err = blk.Deserialize(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("[Chain] block %s not found: %s", hash.Hex(), err.Error()))
	}
	return blk, nil
}

func (bc *Blockchain) GetByRound(round uint64) *msg.Block {
	last := bc.Last
	for round > 0 {
		if last.Round == round {
			return last
		}
		last, _ = bc.Get(common.BytesToHash(last.ParentHash))
		round--
	}
	return last
}

func (bc *Blockchain) Add(blk *msg.Block) {
	data, _ := blk.Serialize()
	bc.db.Put(blk.Hash().Bytes(), data, nil)
	if bc.Last == nil || blk.Round > bc.Last.Round {
		bc.Last = blk
	}
	bc.CacheBlock(blk)
}

func (bc *Blockchain) CacheBlock(blk *msg.Block) {
	bc.cacheLock.Lock()
	defer bc.cacheLock.Unlock()
	hash := blk.Hash()
	if _, ok := bc.blockCacheIndex[hash]; !ok {
		bc.blockCacheIndex[hash] = blk
		removedBlk := bc.blockCache.Push(hash)
		if removedBlk != nil {
			delete(bc.blockCacheIndex, removedBlk.(common.Hash))
		}
	}
}

func (bc *Blockchain) ResolveFork(fork *msg.Block) {
	bc.Last = fork
}

func (bc *Blockchain) PrintState(round uint64) {
	//for round > 0 {
	blk := bc.GetByRound(round)
	log.Printf("State Hash round %d is %s\n", round, common.BytesToHash(blk.StateHash).Hex())
	//round -= 1
	//}
}
