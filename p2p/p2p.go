// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

// This program is free software: you can distribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Bottos.  If not, see <http://www.gnu.org/licenses/>.

// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * file description: p2p
 * @Author: Stewart Li
 * @Date:   2018-02-08
 * @Last Modified by:
 * @Last Modified time:
 */

package p2p

import (
	"fmt"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
	log "github.com/cihub/seelog"
	"net"
)

var LocalPeerInfo PeerInfo
var Runner *P2PServer

type P2PServer struct {
	c      *collection
	connCb NewconnCb

	sendc  chan UniMsgPacket
	bsendc chan BcastMsgPacket
}

type SendupCb func(index uint16, p *Packet)
type NewconnCb func(conn net.Conn)

func MakeP2PServer(p *config.Parameter) *P2PServer {
	LocalPeerInfo.Addr = p.ServAddr
	LocalPeerInfo.Port = p.P2PPort
	LocalPeerInfo.ChainId = p.ChainId

	id := LocalPeerInfo.Addr + LocalPeerInfo.Port
	LocalPeerInfo.Id = common.DoubleSha256([]byte(id)).ToHexString()

	Runner = &P2PServer{
		c:      createCollection(),
		sendc:  make(chan UniMsgPacket, 30),
		bsendc: make(chan BcastMsgPacket, 30),
	}

	return Runner
}

func (s *P2PServer) Start() {
	/*start listen*/
	go s.listenRoutine()
	go s.sendRoutine()
}

func (s *P2PServer) SetCallback(conn NewconnCb) {
	s.connCb = conn
}

func (s *P2PServer) SendUnicast(packet UniMsgPacket) {
	s.sendc <- packet
}

func (s *P2PServer) SendBroadcast(packet BcastMsgPacket) {
	s.bsendc <- packet
}

func (s *P2PServer) AddPeer(peer *Peer) error {
	return s.c.addPeer(peer)
}

func (s *P2PServer) GetPeer(index uint16) *PeerInfo {
	return s.c.getPeer(index)
}

func (s *P2PServer) DelPeer(index uint16) bool {
	return s.c.delPeer(index)
}

func (s *P2PServer) IsPeerExist(index uint16) bool {
	return s.c.isPeerExist(index)
}

func (s *P2PServer) IsPeerInfoExist(info PeerInfo) bool {
	return s.c.isPeerInfoExist(info)
}

func (s *P2PServer) GetPeers() []PeerInfo {
	return s.c.getPeers()
}

func (s *P2PServer) GetPeersData() PeerDataSet {
	return s.c.getPeersData()
}

func (s *P2PServer) listenRoutine() {
	l, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprint(LocalPeerInfo.Port))
	if err != nil {
		log.Errorf("start p2p server listen error: %s", err)
		panic(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error("NetServer::Listening() Failed to accept")
			continue
		}

		/*accpent ten new conection per second*/

		go s.connCb(conn)
	}

	return
}

func (s *P2PServer) sendRoutine() {
	for {
		select {
		case packet := <-s.bsendc:
			s.msend(&packet)
			continue
		default:
			select {
			case packet := <-s.bsendc:
				s.msend(&packet)
			case packet := <-s.sendc:
				s.send(&packet)
			}
		}
	}
}

func (s *P2PServer) send(packet *UniMsgPacket) {
	s.c.send(packet)
}

func (s *P2PServer) msend(packet *BcastMsgPacket) {
	s.c.sendBroadcast(packet)
}
