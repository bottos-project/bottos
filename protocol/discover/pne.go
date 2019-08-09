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
	"strings"
	"sync"
	"time"

	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
)

//DO NOT EDIT
const (
	TIME_FAST_PNE_EXCHANGE = 8
	//TIME_PNE_EXCHANGE time to exchange peer neighbor info, minute
	TIME_PNE_EXCHANGE = 30
)

type pne struct {
	qPeers *common.Queue
	lock   sync.RWMutex

	n *neighbors

	seeds []p2p.PeerInfo
}

func makePne(config *config.P2PConfig) *pne {

	pne := &pne{
		qPeers: common.NewQueue(),
		n:      makeNeighbors(),
	}

	pne.parseSeeds(config)

	return pne
}

func (p *pne) start() {
	go p.pneTimer()
}

func (p *pne) pushPeerIndex(index uint16) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.qPeers.Push(index)
}

func (p *pne) parseSeeds(config *config.P2PConfig) {

	var peers []p2p.PeerInfo

	for _, element := range config.PeerList {
		peerCfg := strings.Split(element, ":")
		if len(peerCfg) != 2 {
			log.Errorf("PROTOCOL parse peer addr of config.PeerList failed, peeraddr: %v", element)
			continue
		}
		addr := peerCfg[0]
		port := peerCfg[1]

		var peer p2p.PeerInfo
		peer.Addr = addr
		peer.Port = port
		peers = append(peers, peer)
		log.Debugf("PROTOCOL parseSeeds: %s:%s", addr, port)
		p.n.addNeighbor(peers)
	}
}

func (p *pne) pneTimer() {
	log.Debug("PROTOCOL pneTimer")

	tripple := 0
	exchange := time.NewTimer(TIME_FAST_PNE_EXCHANGE * time.Second)

	defer func() {
		log.Debug("PROTOCOL pneTimer stop")
		exchange.Stop()
	}()

	for {
		select {
		case <-exchange.C:
			if tripple < 3 {
				log.Debugf("PROTOCOL pneTimer send pne request")
				p.sendPneRequest(0)
				tripple++
				exchange.Reset(TIME_FAST_PNE_EXCHANGE * time.Second)
			} else {
				index := p.nextPeer()
				if index != 0 {
					log.Debugf("PROTOCOL pneTimer peer index: %d", index)
					p.sendPneRequest(index)
				}
				exchange.Reset(TIME_PNE_EXCHANGE * time.Second)
			}
		}
	}
}

//NextPnePeer pop a peer id
func (p *pne) nextPeer() uint16 {
	p.lock.Lock()
	defer p.lock.Unlock()

	value := p.qPeers.Pop()
	if value == nil {
		return 0
	}

	return value.(uint16)
}

func (p *pne) sendPneRequest(index uint16) {
	//check peer is exist or not
	if index != 0 {
		ok := p2p.Runner.IsPeerExist(index)
		if !ok {
			return
		}
	}

	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_NEIGHBOR_REQ,
	}

	packet := p2p.Packet{H: head}

	if index > 0 {
		send := p2p.UniMsgPacket{
			Index: index,
			P:     packet,
		}

		p2p.Runner.SendUnicast(send)

		// add back to queue
		p.pushPeerIndex(index)
	} else {
		send := p2p.BcastMsgPacket{
			Indexs: nil,
			P:      packet,
		}

		p2p.Runner.SendBroadcast(send)
	}

}

func (p *pne) sendPneResponse(index uint16) {
	peers := p2p.Runner.GetPeers()
	if len(peers) == 0 {
		return
	}

	resp := PeerNeighborRsp{
		Neighbor: peers,
	}

	data, err := bpl.Marshal(resp)
	if err != nil {
		log.Errorf("PROTOCOL pne response addrs Marshal error:%s", err)
		return
	}

	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_NEIGHBOR_RSP,
	}

	packet := p2p.Packet{H: head,
		Data: data,
	}

	send := p2p.UniMsgPacket{
		Index: index,
		P:     packet,
	}

	p2p.Runner.SendUnicast(send)
}

func (p *pne) processPneNeighborReq(index uint16, date []byte) {
	//check peer is exist or not
	ok := p2p.Runner.IsPeerExist(index)
	if !ok {
		return
	}

	p.sendPneResponse(index)
}

func (p *pne) processPneNeighborRsp(index uint16, date []byte) {
	//check peer is exist or not
	ok := p2p.Runner.IsPeerExist(index)
	if !ok {
		return
	}

	var rsp PeerNeighborRsp
	err := bpl.Unmarshal(date, &rsp)
	if err != nil {
		log.Errorf("PROTOCOL ProcessPneNeighborRsp Unmarshal error")
		return
	}

	//add neighbor
	p.n.addNeighbor(rsp.Neighbor)
}

type neighbors struct {
	neighbor *list.List
	lock     sync.RWMutex
	nIndex   uint16
}

func makeNeighbors() *neighbors {
	return &neighbors{neighbor: list.New()}
}

//AddNeighbor add new neighbor
func (n *neighbors) addNeighbor(peers []p2p.PeerInfo) {
	n.lock.Lock()
	defer n.lock.Unlock()

	//filter neighbor which is exist
	for j := range peers {
		//filter neighbor of ourself
		if p2p.LocalPeerInfo.Equal(peers[j]) {
			continue
		}

		//filter neighbor of peer
		ok := p2p.Runner.IsPeerInfoExist(peers[j])
		if ok {
			return
		}

		//filter neighbor which is exist
		find := false
		for e := n.neighbor.Front(); e != nil; e = e.Next() {
			peerinfo := e.Value.(p2p.PeerInfo)
			if peerinfo.Equal(peers[j]) {
				find = true
				break
			}
		}

		if !find {
			n.neighbor.PushBack(peers[j])
		}
	}
}

//DelNeighbor delete neighbor
func (n *neighbors) delNeighbor(peer p2p.PeerInfo) {
	n.lock.Lock()
	defer n.lock.Unlock()

	//find addr and remove
	var count uint16
	for e := n.neighbor.Front(); e != nil; e = e.Next() {
		peerinfo := e.Value.(p2p.PeerInfo)
		if peerinfo.Equal(peer) {
			n.neighbor.Remove(e)

			if n.nIndex > count {
				n.nIndex--
			}
			return
		}

		count++
	}
}

//NextPneNeighbors get neighbors to discover
func (n *neighbors) nextPneNeighbors() []p2p.PeerInfo {
	n.lock.Lock()
	defer n.lock.Unlock()

	len := n.neighbor.Len()
	if len <= 0 {
		return nil
	}

	var start uint16
	var end uint16
	if n.nIndex+NEIGHBOR_DISCOVER_COUNT > uint16(len) {
		start = n.nIndex
		end = uint16(len)

		n.nIndex = 0
	} else {
		start = n.nIndex
		end = start + NEIGHBOR_DISCOVER_COUNT

		n.nIndex = end + 1
	}

	var count uint16
	var peers []p2p.PeerInfo
	for e := n.neighbor.Front(); e != nil; e = e.Next() {
		if count >= start && count <= end {
			peers = append(peers, e.Value.(p2p.PeerInfo))
		}
		count++
	}

	return peers
}
