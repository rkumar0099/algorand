package blockchain

import (
	"sync"

	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/manage"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/pool"
)

type Blockchain struct {
	Last            *msg.Block
	Genesis         *msg.Block
	blockCache      *pool.RingBuffer
	blockCacheIndex map[common.Hash]*msg.Block
	cacheLock       *sync.Mutex
	hashes          []cmn.Hash
	manage          *manage.Manage
}

const cacheSize = 10

func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		blockCache:      pool.NewRingBuffer(cacheSize),
		blockCacheIndex: make(map[common.Hash]*msg.Block),
		cacheLock:       &sync.Mutex{},
		hashes:          make([]cmn.Hash, 0),
	}

	emptyHash := common.Sha256([]byte{})
	bc.Genesis = &msg.Block{
		Round:      0,
		Seed:       emptyHash.Bytes(),
		ParentHash: emptyHash.Bytes(),
		Author:     common.HashToAddr(emptyHash).Bytes(),
	}
	bc.addGenesis(bc.Genesis)
	return bc
}

func (bc *Blockchain) addGenesis(blk *msg.Block) {
	//data, _ := blk.Serialize()
	//go bc.manage.AddBlk(blk.Hash(), data)
	bc.hashes = append(bc.hashes, blk.Hash())
	bc.Last = blk
	bc.CacheBlock(blk)
}

func (bc *Blockchain) Get(hash common.Hash) (*msg.Block, error) {
	bc.cacheLock.Lock()
	blk, ok := bc.blockCacheIndex[hash]
	bc.cacheLock.Unlock()
	if ok {
		return blk, nil
	}
	blk = bc.manage.GetBlkByHash(hash)
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
	hash := blk.Hash()
	if len(bc.hashes) >= 10 {
		_, bc.hashes = bc.hashes[0], bc.hashes[1:]
		//delete(bc.persistKv, removedHash)
	}
	bc.hashes = append(bc.hashes, hash)
	data, _ := blk.Serialize()
	go bc.manage.AddBlk(hash, data)
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

func (bc *Blockchain) SetManage(m *manage.Manage) {
	bc.manage = m
	data, _ := bc.Genesis.Serialize()
	bc.manage.AddBlk(bc.Genesis.Hash(), data)
}
