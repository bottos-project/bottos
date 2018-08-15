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
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/p2p"
	log "github.com/cihub/seelog"
	"net"
	"time"
)

//Discover p2p protocol
type Discover struct {
	c      *candidates
	p      *pne
	k      *keeplive
	sendup p2p.SendupCb
}

//DO NOT EDIT
const (
	//TIME_DISCOVER connect to unknow peer, second
	TIME_DISCOVER = 5
	//NEIGHBOR_DISCOVER_COUNT
	NEIGHBOR_DISCOVER_COUNT = 10
)

//MakeDiscover create instance
func MakeDiscover(config *config.Parameter) *Discover {
	d := &Discover{}
	d.p = makePne(config)
	d.c = makeCandidates(d.p)
	d.k = makeKeeplive(d.c, d.p)

	d.c.setKeeplive(d.k)

	return d
}

//Start start...
func (d *Discover) Start() {
	d.c.start()
	d.p.start()
	d.k.start()

	go d.discoverTimer()
}

//SetSendupCallback  set sendup callback
func (d *Discover) SetSendupCallback(cb p2p.SendupCb) {
	d.sendup = cb
}

//Dispatch process peer message
func (d *Discover) Dispatch(index uint16, p *p2p.Packet) {
	//log.Debugf("discovery recv packet %d, from peer: %d", p.H.PacketType, index)

	switch p.H.PacketType {
	case PEER_INFO_REQ:
		d.c.processPeerInfoReq(index, p.Data)
	case PEER_INFO_RSP:
		d.c.processPeerInfoRsp(index, p.Data)
	case PEER_HANDSHAKE_REQ:
		d.c.processHandshakeReq(index, p.Data)
	case PEER_HANDSHAKE_RSP:
		d.c.processHandshakeRsp(index, p.Data)
	case PEER_HANDSHAKE_RSP_ACK:
		d.c.processHandshakeRspAck(index, p.Data)
	case PEER_NEIGHBOR_REQ:
		d.p.processPneNeighborReq(index, p.Data)
	case PEER_NEIGHBOR_RSP:
		d.p.processPneNeighborRsp(index, p.Data)
	case PEER_PING:
		d.k.processPing(index, p.Data)
	case PEER_PONG:
		d.k.processPong(index, p.Data)
	default:
		log.Errorf("discover Dispatch packet type:%d error", p.H.PacketType)
	}

}

//NewConnCb accept a connection with a peer
func (d *Discover) NewConnCb(conn net.Conn, sendup p2p.SendupCb) {
	//new candidate peer
	info := p2p.PeerInfo{}
	p := p2p.CreatePeer(info, conn, true, sendup)

	error := d.c.addCandidate(p)
	if error != nil {
		p.Stop()
		return
	}

	p.Start()
}

//newConn create a connection with a peer
func (d *Discover) newConn(peer p2p.PeerInfo) error {
	addrPort := peer.Addr + ":" + peer.Port
	conn, err := net.DialTimeout("tcp", addrPort, 2*time.Second)
	if err != nil {
		log.Debugf("connect to peer %s:%s error:%s", peer.Addr, peer.Port, err)
		return err
	}

	p := p2p.CreatePeer(peer, conn, false, d.sendup)

	err = d.c.addCandidate(p)
	if err != nil {
		p.Stop()
		return err
	}

	p.Start()
	return nil
}

func (d *Discover) discoverTimer() {
	log.Debug("discoverTimer")

	dicover := time.NewTimer(TIME_DISCOVER * time.Second)

	defer func() {
		log.Debug("discoverTimer stop")
		dicover.Stop()
	}()

	for {
		select {
		case <-dicover.C:
			if d.c.isCandidateFull() {
				dicover.Reset(TIME_DISCOVER * time.Second)
				continue
			}

			neighbors := d.p.n.nextPneNeighbors()
			if neighbors == nil {
				dicover.Reset(TIME_DISCOVER * time.Second)
				continue
			}

			/*try to connect peer*/
			for _, peer := range neighbors {
				err := d.newConn(peer)
				if err != nil {
					continue
				}

				d.p.n.delNeighbor(peer)
			}

			dicover.Reset(TIME_DISCOVER * time.Second)
		}
	}
}
