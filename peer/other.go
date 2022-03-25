package peer

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/rkumar0099/algorand/common"
	cmn "github.com/rkumar0099/algorand/common"
	"github.com/rkumar0099/algorand/gossip"
	"github.com/rkumar0099/algorand/logs"
	"github.com/rkumar0099/algorand/manage"
	msg "github.com/rkumar0099/algorand/message"
	"github.com/rkumar0099/algorand/oracle"
	"github.com/rkumar0099/algorand/params"
	"github.com/rkumar0099/algorand/service"
	"google.golang.org/grpc"
)

func (p *Peer) Start(neighbors []gossip.NodeId, addr [][]byte) error {
	lis, err := net.Listen("tcp", p.Id.String())
	if err != nil {
		s := fmt.Sprintf("[algorand] [%s] Cannot start peer: %s\n", p.Id.String(), err.Error())
		log.Print(s)
		p.lm.AddLog(s)

		//log.Printf("[algorand] [%s] Cannot start peer: %s", p.Id.String(), err.Error())
		return errors.New(fmt.Sprintf("Cannot start algorand peer: %s", err.Error()))
	}
	go p.grpcServer.Serve(lis)
	time.Sleep(1 * time.Second)
	p.node.Join(neighbors)
	go p.msgAgent.Handle()

	p.initialize(addr, p.permanentTxStorage)
	//log.Printf("Addresses init: %s\n", cmn.Sha256(p.lastState).String())
	return nil
}

func (p *Peer) AddManage(m *manage.Manage, lm *logs.LogManager) {
	p.manage = m
	p.lm = lm
	p.chain.SetManage(m)
}

func (p *Peer) AddOracle(oracle *oracle.Oracle) {
	p.oracle = oracle
	blk := p.chain.GetByRound(0)
	go p.lm.AddFinalBlk(blk.Hash(), 0)
	go p.oracle.AddBlk(blk) // add initial algorand block to oracle
}

func (p *Peer) StartServices(bootNode []gossip.NodeId) error {
	lis, err := net.Listen("tcp", p.Id.String())
	if err != nil {
		s := fmt.Sprintf("[algorand] [%s] Cannot start peer: %s\n", p.Id.String(), err.Error())
		log.Print(s)
		p.lm.AddLog(s)
		//log.Printf("[algorand] [%s] Cannot start peer: %s", p.Id.String(), err.Error())
		return errors.New(fmt.Sprintf("Cannot start algorand peer: %s", err.Error()))
	}
	go p.grpcServer.Serve(lis)
	p.node.Join(bootNode)
	p.msgAgent.Handle()
	return nil
}

func (p *Peer) Stop() {
	close(p.quitCh)
	close(p.hangForever)
}

func (p *Peer) GetGrpcServer() *grpc.Server {
	return p.grpcServer
}

func (p *Peer) getDataByHashHandler(hash []byte) ([]byte, error) {
	// to do
	return nil, nil
}

func (p *Peer) GetLastState() string {
	return cmn.Sha256(p.lastState).String()
}

// round returns the latest round number.
func (p *Peer) round() uint64 {
	return p.lastBlock().Round
}

func (p *Peer) lastBlock() *msg.Block {
	return p.chain.Last
}

// weight returns the weight of the given address.
func (p *Peer) weight(address cmn.Address) uint64 {
	return params.TokenPerUser
}

// tokenOwn returns the token amount (weight) owned by self node.
func (p *Peer) tokenOwn() uint64 {
	return p.weight(p.Address())
}

func (p *Peer) emptyBlock(round uint64, prevHash cmn.Hash) *msg.Block {
	prevBlk, err := p.getBlock(prevHash)
	if err != nil {
		s := fmt.Sprintf("node %s hang forever because cannot get previous block\n", p.Id.String())
		log.Print(s)
		p.lm.AddLog(s)
		//log.Printf("node %d hang forever because cannot get previous block", p.Id)
		<-p.hangForever
	}
	return &msg.Block{
		Round:      round,
		ParentHash: prevHash.Bytes(),
		StateHash:  prevBlk.StateHash,
	}
}

func (p *Peer) Address() cmn.Address {
	return cmn.BytesToAddress(p.pubkey.Bytes())
}

// forkLoop periodically resolves fork
func (p *Peer) forkLoop() {
	forkInterval := time.NewTicker(params.ForkResolveInterval)

	for {
		select {
		case <-p.quitCh:
			return
		case <-forkInterval.C:
			p.processForkResolve()
		}
	}
}

func (p *Peer) gossip(typ int, data []byte) {
	message := &msg.Msg{
		PID:  p.Id.String(),
		Type: int32(typ),
		Data: data,
	}
	switch typ {
	case msg.BLOCK:
		s := fmt.Sprintf("[debug] [%s] gossip block\n", p.Id.String())
		log.Print(s)
		p.lm.AddLog(s)
	case msg.BLOCK_PROPOSAL:
		s := fmt.Sprintf("[debug] [%s] gossip block proposal\n", p.Id.String())
		log.Print(s)
		p.lm.AddLog(s)
	case msg.VOTE:
		s := fmt.Sprintf("[debug] [%s] gossip vote\n", p.Id.String())
		log.Print(s)
		p.lm.AddLog(s)

	}
	msgBytes, err := proto.Marshal(message)
	if err != nil {
		log.Printf("[alogrand] [%s] cannot gossip message: %s", err.Error())
	} else {
		p.node.Gossip(msgBytes)
	}
	p.lm.AddProcessLog(message)
}

func (p *Peer) emptyHash(round uint64, prev common.Hash) common.Hash {
	return p.emptyBlock(round, prev).Hash()
}

func (p *Peer) handleBlock(data []byte) {
	blk := &msg.Block{}
	if err := blk.Deserialize(data); err != nil {
		//log.Printf("[algorand] [%s] Received invalid block message: %s", p.Id.String(), err.Error())
		return
	}
	p.chain.CacheBlock(blk)
}

func (p *Peer) getBlock(hash common.Hash) (*msg.Block, error) {
	// find  locally
	blk, err := p.chain.Get(hash)
	if err != nil {
		s := fmt.Sprintf("[algorand] [%s] cannot get block %s locally: %s, try to find from other peers\n", p.Id.String(), hash.Hex(), err.Error())
		log.Print(s)
		p.lm.AddLog(s)
		//log.Printf("[algorand] [%s] cannot get block %s locally: %s, try to find from other peers", p.Id.String(), hash.Hex(), err.Error())
	} else {
		return blk, nil
	}

	// find from other peers
	neighborList := p.node.GetNeighborList()
	neighbors := neighborList.GetNeighborsId()
	chanBlk := make(chan *msg.Block, 1)
	blkReqDone := false
	// make query in batch of 10

	go func() {
		for i := 0; i < len(neighbors) && !blkReqDone; i += 10 {
			for j := i; j < i+10 && j < len(neighbors) && !blkReqDone; j++ {
				go func(nodeId gossip.NodeId) {
					conn, err := neighborList.GetConn(nodeId)
					if err != nil {
						//log.Printf("[alogrand] [%s] cannot get block from [%s]: %s", p.Id.String(), nodeId.String(), err.Error())
						return
					}
					blk, err := service.GetBlock(conn, hash.Bytes())
					if err != nil {
						//log.Printf("[alogrand] [%s] cannot get block from [%s]: %s", p.Id.String(), nodeId.String(), err.Error())
						return
					}
					select {
					case chanBlk <- blk:
						s := fmt.Sprintf("[debug] [%s] got block %s from [%s]\n", p.Id.String(), hash.Hex(), nodeId.String())
						log.Print(s)
						p.lm.AddLog(s)
						//log.Printf("[debug] [%s] got block %s from [%s]", p.Id.String(), hash.Hex(), nodeId.String())
						blkReqDone = true
					default:
					}
				}(neighbors[j])
			}
			time.Sleep(time.Second)
		}
	}()

	select {
	case <-time.NewTimer(10 * time.Second).C:
		blk = nil
	case blk = <-chanBlk:
	}

	if blk == nil {
		return nil, errors.New(fmt.Sprintf("[algorand] [%s] cannot get block %s from other peers: %s\n", p.Id.String(), hash.Hex(), err.Error()))
	}
	p.chain.CacheBlock(blk)
	s := fmt.Sprintf("[debug] [%s] cached block %s\n", p.Id.String(), hash.Hex())
	log.Print(s)
	p.lm.AddLog(s)
	//log.Printf("[debug] [%s] cached block %s", p.Id.String(), hash.Hex())
	return blk, nil
}

func (p *Peer) sendBalance() {
	for {
		time.Sleep(5 * time.Second)
		p.lm.AddBalance(p.Id.String(), p.GetBalance())
	}
}
