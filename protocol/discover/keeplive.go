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
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
	"sync/atomic"
	"time"
)

//DO NOT EDIT
const (
	//TIME_KEEP_LIVE ping/pong timer, second
	TIMER_KEEP_LIVE = 10
	//TIMER_CHECK time out second
	TIMER_CHECK = 40
)

type keeplive struct {
	counter [MAX_PEER_COUNT + 1]int32

	c *candidates
	p *pne
}

func makeKeeplive(c *candidates, p *pne) *keeplive {
	return &keeplive{c: c, p: p}
}

func (k *keeplive) start() {
	for i := 0; i < MAX_PEER_COUNT; i++ {
		k.counter[i] = -1
	}

	go k.keepliveTimer()
	go k.checkTimer()
}

func (k *keeplive) initCounter(index uint16) {
	atomic.StoreInt32(&k.counter[index], 0)
}

func (k *keeplive) keepliveTimer() {
	log.Debug("PROTOCOL keepliveTimer")

	keep := time.NewTimer(TIMER_KEEP_LIVE * time.Second)

	defer func() {
		log.Debug("PROTOCOL keepliveTimer stop")
		keep.Stop()
	}()

	for {
		select {
		case <-keep.C:
			k.sendPing()
			keep.Reset(TIMER_KEEP_LIVE * time.Second)
		}
	}
}

func (k *keeplive) checkTimer() {
	log.Debug("PROTOCOL checkTimer")

	check := time.NewTimer(TIMER_CHECK * time.Second)

	defer func() {
		log.Debug("PROTOCOL checkTimer stop")
		check.Stop()
	}()

	for {
		select {
		case <-check.C:
			k.checkPeer()
			check.Reset(TIMER_CHECK * time.Second)
		}
	}
}

func (k *keeplive) checkPeer() {
	for i := 0; i < MAX_PEER_COUNT; i++ {
		if k.counter[i] != -1 {
			if k.counter[i] == 0 {
				info := p2p.Runner.GetPeer(uint16(i))
				var set []p2p.PeerInfo
				set = append(set, *info)

				if p2p.Runner.DelPeer(uint16(i)) {
					log.Infof("PROTOCOL peer %s:%s disconnect, add back to connect neighbors", info.Addr, info.Port)
					k.c.pushPeerIndex(uint16(i))
					k.p.n.addNeighbor(set)
					atomic.StoreInt32(&k.counter[i], -1)
				}
			} else {
				atomic.StoreInt32(&k.counter[i], 0)
			}
		}
	}
}

func (k *keeplive) processPing(index uint16, date []byte) {
	k.sendPong(index)
}

func (k *keeplive) processPong(index uint16, date []byte) {
	k.counterPeer(index)
}

func (k *keeplive) counterPeer(index uint16) {
	atomic.AddInt32(&k.counter[index], 1)
}

func (k *keeplive) sendPing() {
	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_PING,
	}

	packet := p2p.Packet{H: head}

	ping := p2p.BcastMsgPacket{
		Indexs: nil,
		P:      packet,
	}

	p2p.Runner.SendBroadcast(ping)

}

func (k *keeplive) sendPong(index uint16) {
	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_PONG,
	}

	packet := p2p.Packet{H: head}

	pong := p2p.UniMsgPacket{
		Index: index,
		P:     packet,
	}

	p2p.Runner.SendUnicast(pong)

}
