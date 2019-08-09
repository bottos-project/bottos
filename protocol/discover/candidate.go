// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  producer actor
 * @Author: eripi
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package discover

import (
	"container/list"
	"github.com/bottos-project/bottos/bpl"
	"errors"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
	"sync"
	"time"
)

//DO NOT EDIT
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
	log.Debug("PROTOCOL exchangeTimer start")

	exchangeTimer := time.NewTimer(TIMER_PEER_EXCHANGE * time.Second)

	defer func() {
		log.Debug("PROTOCOL exchangeTimer stop")
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
			log.Debugf("PROTOCOL exchange max req count index: %d", candi.peer.Index)
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
		log.Error("PROTOCOL candidates full")
		return errors.New("candidates full")
	}

	log.Debugf("PROTOCOL addCandidate index: %d", index.(uint16))
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

	log.Debugf("PROTOCOL deleteCandidate index: %d", index)

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
		log.Debugf("PROTOCOL processPeerInfoReq candi not exist index: %d", index)
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
		log.Debugf("PROTOCOL processPeerInfoRsp candi not exist index: %d", index)
		return
	}

	candi := e.Value.(*candidate)

	var rsp PeerInfoRsp
	err := bpl.Unmarshal(date, &rsp)
	if err != nil {
		log.Error("PROTOCOL processPeerInfoRsp Unmarshal error")
		return
	}

	if rsp.Info.IsIncomplete() {
		log.Error("PROTOCOL processPeerInfoRsp rsp info error")
		return
	}

	if rsp.Info.ChainId != p2p.LocalPeerInfo.ChainId {
		log.Error("PROTOCOL not on the same chain, drop candidate")
		c.deleteCandidate(e, true)
		return
	}

	//check peer addr and port if the connection is our init
	if !candi.peer.In &&
		candi.peer.Info.Addr != rsp.Info.Addr &&
		candi.peer.Info.Port != rsp.Info.Port {
		log.Errorf("PROTOCOL processPeerInfoRsp wrong peer info addr: %s, port: %s", rsp.Info.Addr, rsp.Info.Port)
		return
	}

	candi.peer.Info = rsp.Info
	candi.peer.State = p2p.PEER_STATE_HANDSHAKE

	c.sendHandshakeReq(candi)
}

func (c *candidates) processHandshakeReq(index uint16, date []byte) {
	c.l.Lock()
	defer c.l.Unlock()

	e := c.getCandidate(index)
	if e == nil {
		log.Debugf("PROTOCOL processHandshakeReq candi not exist index: %d", index)
		return
	}

	candi := e.Value.(*candidate)

	if candi.peer.State != p2p.PEER_STATE_HANDSHAKE {
		log.Debug("PROTOCOL processHandshakeReq not in hand shake state")
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
		log.Debug("PROTOCOL processPeerInfoReq candi not exist ")
		return
	}

	candi := ec.Value.(*candidate)

	if candi.peer.State != p2p.PEER_STATE_HANDSHAKE {
		log.Debug("PROTOCOL processHandshakeReq not in hand shake state")
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

	//add peer
	bReturnIndex := false
	err := p2p.Runner.AddPeer(candi.peer)
	if err == nil {
		//send response ack
		c.sendHandshakeRspAck(candi)

		c.p.pushPeerIndex(candi.peer.Index)
		c.k.initCounter(candi.peer.Index)
	} else {
		candi.peer.Stop()
		bReturnIndex = true
	}

	/*remove from canidata*/
	c.deleteCandidate(ec, bReturnIndex)

}

func (c *candidates) processHandshakeRspAck(index uint16, date []byte) {
	c.l.Lock()
	defer c.l.Unlock()

	var ec *list.Element
	ec = c.getCandidate(index)
	if ec == nil {
		log.Debug("PROTOCOL processHandshakeRspAck candi not exist ")
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
	bReturnIndex := false
	err := p2p.Runner.AddPeer(candi.peer)
	if err == nil {
		c.p.pushPeerIndex(candi.peer.Index)
		c.k.initCounter(candi.peer.Index)
	} else {
		candi.peer.Stop()
		bReturnIndex = true
	}

	/*remove from canidata*/
	c.deleteCandidate(ec, bReturnIndex)
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
	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_INFO_REQ,
	}

	packet := p2p.Packet{H: head}

	candi.peer.Send(packet)
}

func (c *candidates) sendPeerInfoRsp(candi *candidate) {
	info := p2p.PeerInfo{
		Id:      p2p.LocalPeerInfo.Id,
		Addr:    p2p.LocalPeerInfo.Addr,
		Port:    p2p.LocalPeerInfo.Port,
		ChainId: p2p.LocalPeerInfo.ChainId,
	}

	rsp := PeerInfoRsp{
		Info: info,
	}

	data, err := bpl.Marshal(rsp)
	if err != nil {
		log.Error("PROTOCOL sendPeerInfoRsp Marshal data error ")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_INFO_RSP,
	}

	packet := p2p.Packet{
		H:    head,
		Data: data,
	}

	candi.peer.Send(packet)
}

func (c *candidates) sendHandshakeReq(candi *candidate) {
	//bigger peer send hand shake
	if p2p.LocalPeerInfo.Bigger(candi.peer.Info) < 1 {
		log.Debugf("PROTOCOL sendHandshakeReq local is small")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_HANDSHAKE_REQ,
	}

	packet := p2p.Packet{H: head}

	candi.peer.Send(packet)
}

func (c *candidates) sendHandshakeRsp(candi *candidate) {
	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_HANDSHAKE_RSP,
	}

	packet := p2p.Packet{H: head}

	candi.peer.Send(packet)
}

func (c *candidates) sendHandshakeRspAck(candi *candidate) {
	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_HANDSHAKE_RSP_ACK,
	}

	packet := p2p.Packet{H: head}

	candi.peer.Send(packet)
}
