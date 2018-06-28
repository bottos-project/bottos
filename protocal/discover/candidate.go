package discover

import (
	"container/list"
	"encoding/json"
	"errors"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocal/common"
	log "github.com/cihub/seelog"
	"sync"
	"time"
)

const (
	//TIME_PEER_INFO_EXCHANGE get peer info timer,  second
	TIMER_PEER_EXCHANGE = 5
	MAX_REQ_COUNT       = 10
)

type candidate struct {
	peer  *p2p.Peer
	count uint16
}

type candidates struct {
	cs     *list.List
	qindex *common.Queue
	l      sync.RWMutex

	p *pne
	k *keeplive
}

func makeCandidates(p *pne) *candidates {
	cs := candidates{
		cs:     list.New(),
		qindex: common.NewQueue(),
		p:      p,
	}

	var i uint16
	for i = 1; i <= MAX_PEER_COUNT; i++ {
		cs.qindex.Push(i)
	}

	return &cs
}

func (c *candidates) start() {
	go c.exchangeTimer()
}

func (c *candidates) setKeeplive(k *keeplive) {
	c.k = k
}

func (c *candidates) exchangeTimer() {
	log.Debug("exchangeTimer start")

	exchangeTimer := time.NewTimer(TIMER_PEER_EXCHANGE * time.Second)

	defer func() {
		log.Debug("exchangeTimer stop")
		exchangeTimer.Stop()
	}()

	for {
		select {
		case <-exchangeTimer.C:
			c.exchange()
			exchangeTimer.Reset(TIMER_PEER_EXCHANGE * time.Second)
		}
	}
}

func (c *candidates) exchange() {
	c.l.Lock()
	defer c.l.Unlock()

	var next *list.Element
	for e := c.cs.Front(); e != nil; {
		candi := e.Value.(*candidate)
		if candi.count >= MAX_REQ_COUNT {
			log.Debugf("exchange max req count index: %d", candi.peer.Index)
			next = e.Next()
			candi.peer.Stop()
			c.deleteCandidate(e, true)
			e = next
		} else {
			candi.count++
			if candi.peer.State == p2p.PEER_STATE_INIT {
				c.sendPeerInfoReq(candi)
			} else if candi.peer.State == p2p.PEER_STATE_HANDSHAKE {
				c.sendHandshakeReq(candi)
			}
			e = e.Next()
		}
	}
}

func (c *candidates) isCandidateFull() bool {
	c.l.Lock()
	defer c.l.Unlock()

	if c.qindex.Length() == 0 {
		return true
	}
	return false
}

func (c *candidates) addCandidate(peer *p2p.Peer) error {
	c.l.Lock()
	defer c.l.Unlock()

	index := c.qindex.Pop()
	if index == nil {
		log.Error("candidates full")
		return errors.New("candidates full")
	}

	log.Debugf("AddCandidate index: %d", index.(uint16))
	peer.Index = index.(uint16)
	peer.State = p2p.PEER_STATE_INIT
	candi := &candidate{peer: peer, count: 0}

	c.cs.PushBack(candi)
	c.sendPeerInfoReq(candi)
	return nil
}

func (c *candidates) pushPeerIndex(index uint16) {
	c.l.Lock()
	defer c.l.Unlock()

	c.qindex.Push(index)
}

func (c *candidates) deleteCandidate(e *list.Element, bRetureIndex bool) {
	candi := e.Value.(*candidate)
	index := candi.peer.Index

	log.Debugf("deleteCandidate index: %d", index)

	c.cs.Remove(e)
	if bRetureIndex {
		c.qindex.Push(index)
	}

}

func (c *candidates) processPeerInfoReq(index uint16, date []byte) {
	c.l.Lock()
	defer c.l.Unlock()

	e := c.getCandidate(index)
	if e == nil {
		log.Debug("ProcessPeerInfoReq candi not exist index: %d", index)
		return
	}

	candi := e.Value.(*candidate)
	c.sendPeerInfoRsp(candi)
}

func (c *candidates) processPeerInfoRsp(index uint16, date []byte) {
	c.l.Lock()
	defer c.l.Unlock()

	e := c.getCandidate(index)
	if e == nil {
		log.Debugf("ProcessPeerInfoRsp candi not exist index: %d", index)
		return
	}

	candi := e.Value.(*candidate)

	var rsp PeerInfoRsp
	err := json.Unmarshal(date, &rsp)
	if err != nil {
		log.Error("ProcessPeerInfoRsp Unmarshal error")
		return
	}

	if rsp.Info.IsIncomplete() {
		log.Error("ProcessPeerInfoRsp rsp info error")
		return
	}

	candi.peer.Info.Id = rsp.Info.Id
	//check peer addr and port if the connection is our init
	if !candi.peer.In {
		if candi.peer.Info.Addr != rsp.Info.Addr &&
			candi.peer.Info.Port != rsp.Info.Port {
			log.Errorf("ProcessPeerInfoRsp wrong peer info addr: %s, port: %s", rsp.Info.Addr, rsp.Info.Port)
			return
		}
	} else {
		candi.peer.Info.Addr = rsp.Info.Addr
		candi.peer.Info.Port = rsp.Info.Port
	}
	candi.peer.State = p2p.PEER_STATE_HANDSHAKE

	c.sendHandshakeReq(candi)
}

func (c *candidates) processHandshakeReq(index uint16, date []byte) {
	c.l.Lock()
	defer c.l.Unlock()

	e := c.getCandidate(index)
	if e == nil {
		log.Debugf("ProcessHandshakeReq candi not exist index: %d", index)
		return
	}

	candi := e.Value.(*candidate)

	if candi.peer.State != p2p.PEER_STATE_HANDSHAKE {
		log.Debug("ProcessHandshakeReq not in hand shake state")
		return
	}

	c.sendHandshakeRsp(candi)
}

func (c *candidates) processHandshakeRsp(index uint16, date []byte) {
	c.l.Lock()
	defer c.l.Unlock()

	var ec *list.Element
	ec = c.getCandidate(index)
	if ec == nil {
		log.Debug("ProcessPeerInfoReq candi not exist ")
		return
	}

	candi := ec.Value.(*candidate)

	if candi.peer.State != p2p.PEER_STATE_HANDSHAKE {
		log.Debug("ProcessHandshakeReq not in hand shake state")
		return
	}

	//check peer

	/*check duplicate candidate*/
	for e := c.cs.Front(); e != nil; e = e.Next() {
		temp := e.Value.(*candidate)
		if temp.peer.Info.Equal(candi.peer.Info) && temp.peer.Index != candi.peer.Index {
			temp.peer.Stop()
			c.deleteCandidate(e, true)
			break
		}
	}

	//send response ack
	c.sendHandshakeRspAck(candi)

	//add peer
	err := p2p.Runner.AddPeer(candi.peer)
	if err == nil {
		c.p.pushPeerIndex(candi.peer.Index)
		c.k.initCounter(candi.peer.Index)
	} else {
		candi.peer.Stop()
	}

	/*remove from canidata*/
	c.deleteCandidate(ec, false)

}

func (c *candidates) processHandshakeRspAck(index uint16, date []byte) {
	c.l.Lock()
	defer c.l.Unlock()

	var ec *list.Element
	ec = c.getCandidate(index)
	if ec == nil {
		log.Debug("ProcessHandshakeRspAck candi not exist ")
		return
	}

	candi := ec.Value.(*candidate)

	/*check duplicate candidate*/
	for e := c.cs.Front(); e != nil; e = e.Next() {
		temp := e.Value.(*candidate)
		if temp.peer.Info.Equal(candi.peer.Info) && temp.peer.Index != candi.peer.Index {
			temp.peer.Stop()
			c.deleteCandidate(e, true)
			break
		}
	}

	//add peer
	err := p2p.Runner.AddPeer(candi.peer)
	if err == nil {
		c.p.pushPeerIndex(candi.peer.Index)
		c.k.initCounter(candi.peer.Index)
	} else {
		candi.peer.Stop()
	}

	/*remove from canidata*/
	c.deleteCandidate(ec, false)
}

func (c *candidates) getCandidate(index uint16) *list.Element {
	for e := c.cs.Front(); e != nil; e = e.Next() {
		candi := e.Value.(*candidate)
		if candi.peer.Index == index {
			return e
		}
	}

	return nil
}

func (c *candidates) sendPeerInfoReq(candi *candidate) {
	head := p2p.Head{ProtocalType: pcommon.P2P_PACKET,
		PacketType: PEER_INFO_REQ,
	}

	packet := p2p.Packet{H: head}

	candi.peer.Send(packet)
}

func (c *candidates) sendPeerInfoRsp(candi *candidate) {
	info := p2p.PeerInfo{
		Id:   p2p.LocalPeerInfo.Id,
		Addr: p2p.LocalPeerInfo.Addr,
		Port: p2p.LocalPeerInfo.Port,
	}

	rsp := PeerInfoRsp{
		Info: info,
	}

	data, err := json.Marshal(rsp)
	if err != nil {
		log.Error("sendPeerInfoRsp Marshal data error ")
		return
	}

	head := p2p.Head{ProtocalType: pcommon.P2P_PACKET,
		PacketType: PEER_INFO_RSP,
	}

	packet := p2p.Packet{
		H:    head,
		Data: data,
	}

	candi.peer.Send(packet)
}

func (c *candidates) sendHandshakeReq(candi *candidate) {
	//hold bigger peer send hand shake
	if p2p.LocalPeerInfo.Bigger(candi.peer.Info) < 1 {
		log.Debugf("sendHandshakeReq local is small")
		return
	}

	head := p2p.Head{ProtocalType: pcommon.P2P_PACKET,
		PacketType: PEER_HANDSHAKE_REQ,
	}

	packet := p2p.Packet{H: head}

	candi.peer.Send(packet)
}

func (c *candidates) sendHandshakeRsp(candi *candidate) {
	head := p2p.Head{ProtocalType: pcommon.P2P_PACKET,
		PacketType: PEER_HANDSHAKE_RSP,
	}

	packet := p2p.Packet{H: head}

	candi.peer.Send(packet)
}

func (c *candidates) sendHandshakeRspAck(candi *candidate) {
	head := p2p.Head{ProtocalType: pcommon.P2P_PACKET,
		PacketType: PEER_HANDSHAKE_RSP_ACK,
	}

	packet := p2p.Packet{H: head}

	candi.peer.Send(packet)
}
