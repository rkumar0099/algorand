package pool

import (
	"bytes"
	"sync"
	"time"

	"github.com/golang-collections/collections/set"
	"github.com/rkumar0099/algorand/common"
	msg "github.com/rkumar0099/algorand/message"
)

type VoteInfo struct {
	votes     chan *msg.VoteMessage
	voters    *set.Set
	startTime time.Time
}

type VotePool struct {
	pool          map[string]*VoteInfo
	bufferSize    int
	cleanDuration time.Duration
	lock          *sync.RWMutex
	verifier      VoteVerifier
}

type VoteVerifier func(*msg.VoteMessage, int) int

func NewVoteInfo(bufferSize int) *VoteInfo {
	return &VoteInfo{
		votes:  make(chan *msg.VoteMessage, bufferSize),
		voters: set.New(),
	}
}

func NewVotePool(bufferSize int, verifier VoteVerifier, cleanDuration time.Duration) *VotePool {
	return &VotePool{
		pool:          make(map[string]*VoteInfo),
		bufferSize:    bufferSize,
		cleanDuration: cleanDuration,
		lock:          &sync.RWMutex{},
		verifier:      verifier,
	}
}

func (vp *VotePool) HandleVote(vote *msg.VoteMessage) {
	vp.lock.Lock()
	defer vp.lock.Unlock()
	voteKey := common.GetVoteKey(vote.Round, vote.Event)
	voteInfo, ok := vp.pool[voteKey]
	if !ok {
		voteInfo = NewVoteInfo(vp.bufferSize)
		vp.pool[voteKey] = voteInfo
		vp.cleanAfter(voteKey, vp.cleanDuration)
	}
	if !voteInfo.voters.Has(string(vote.RecoverPubkey().Pk)) {
		select {
		case voteInfo.votes <- vote:
		default:
			//log.Printf("[VotePool] vote buffer full for vote key: %s", common.BytesToHash(vote.Hash).Hex())
		}
	}
}

func (vp *VotePool) CountVotes(round uint64, event string, threshold float64, expectedNum int, timeout time.Duration) (common.Hash, error) {
	vp.lock.Lock()
	voteKey := common.GetVoteKey(round, event)
	voteInfo, ok := vp.pool[voteKey]
	if !ok {
		voteInfo = NewVoteInfo(vp.bufferSize)
		vp.pool[voteKey] = voteInfo
		vp.cleanAfter(voteKey, vp.cleanDuration)
	}
	vp.lock.Unlock()
	count := make(map[common.Hash]int)
	timer := time.NewTimer(timeout).C
	for {
		select {
		case <-timer:
			return common.Hash{}, common.ErrCountVotesTimeout
		case vote := <-voteInfo.votes:
			hash := common.BytesToHash(vote.Hash)
			count[hash] += vp.verifier(vote, expectedNum)
			if count[hash] > int(threshold*float64(expectedNum)) {
				return hash, nil
			}
		}
	}
}
func (vp *VotePool) CountVotesAndCoin(round uint64, event string, threshold float64, expectedNum int, timeout time.Duration) (common.Hash, int64, error) {
	vp.lock.Lock()
	voteKey := common.GetVoteKey(round, event)
	voteInfo, ok := vp.pool[voteKey]
	if !ok {
		voteInfo = NewVoteInfo(vp.bufferSize)
		vp.pool[voteKey] = voteInfo
		vp.cleanAfter(voteKey, vp.cleanDuration)
	}
	vp.lock.Unlock()
	count := make(map[common.Hash]int)
	timer := time.NewTimer(timeout).C
	minHash := make([]byte, common.HashLength)
	for i := range minHash {
		minHash[i] = 255
	}
	for {
		select {
		case <-timer:
			return common.Hash{}, int64(minHash[0]) % 2, common.ErrCountVotesTimeout
		case vote := <-voteInfo.votes:
			hash := common.BytesToHash(vote.Hash)
			count[hash] += vp.verifier(vote, expectedNum)
			if count[hash] > int(threshold*float64(expectedNum)) {
				return hash, 0, nil
			}
			if bytes.Compare(vote.VRF, minHash) < 0 {
				minHash = vote.VRF
			}
		}
	}
}

func (vp *VotePool) cleanAfter(voteKey string, timeout time.Duration) {
	time.AfterFunc(timeout, func() {
		vp.lock.Lock()
		defer vp.lock.Unlock()
		delete(vp.pool, voteKey)
	})
}

func (vp *VotePool) Clean(voteKey string) {
	vp.lock.Lock()
	defer vp.lock.Unlock()
	delete(vp.pool, voteKey)
}
