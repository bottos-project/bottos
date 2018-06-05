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

// Package exec provides functions for executing WebAssembly bytecode.

/*
 * file description: the interface for WASM execution
 * @Author: Stewart Li
 * @Date:   2018-02-08
 * @Last Modified by:    Stewart Li
 * @Last Modified time:  2018-05-16
 */

package p2pserver

import (
	"errors"
	"fmt"
	"net"
)

//Peer info
type Peer struct {
	peerAddr  string
	servPort  int
	peerId    uint32
	publicKey string

	peerSock *net.UDPAddr
	conn     net.Conn

	syncState    uint32
	neighborNode []*Peer
}

//NewPeer is to create new peer
func NewPeer(addrName string, servPort int, conn net.Conn) *Peer {
	return &Peer{
		peerAddr:  addrName,
		servPort:  servPort,
		peerId:    0,
		conn:      conn,
		syncState: 0,
	}
}
//GetPeerAddr is getting peer address
func (p *Peer) GetPeerAddr() string {
	return p.peerAddr
}
//SetPeerAddr is setting peer address
func (p *Peer) SetPeerAddr(addr string) {
	p.peerAddr = addr
}
//SetPeerState is setting peer state
func (p *Peer) SetPeerState(state uint32) {
	p.syncState = state
}
//GetPeerState is getting peer state
func (p *Peer) GetPeerState() uint32 {
	return p.syncState
}
//GetId is to get id
func (p *Peer) GetId() uint64 {
	if p.peerId == 0 {
		addr_port := p.peerAddr + ":" + fmt.Sprint(p.servPort)
		p.peerId = Hash(addr_port)
	}

	return uint64(p.peerId)
}
//SendTo is to send to message buffer
func (p *Peer) SendTo(buf []byte, isSync bool) error {
	len, err := p.conn.Write(buf)
	if err != nil {
		return errors.New("*ERROR* Failed to send data !!!")
	} else if len <= 0 {
		return errors.New("*ERROR* Failed to send data !!!")
	}

	return nil
}
