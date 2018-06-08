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
 * @Last Modified by:
 * @Last Modified time:
 */

package p2pserver

import (
	"errors"
	"fmt"
	"net"
)

//Peer peer info
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

//NewPeer new a peer
func NewPeer(addrName string, servPort int, conn net.Conn) *Peer {
	return &Peer{
		peerAddr:  addrName,
		servPort:  servPort,
		peerId:    0,
		conn:      conn,
		syncState: 0,
	}
}

//GetPeerAddr get peer addr
func (p *Peer) GetPeerAddr() string {
	return p.peerAddr
}

//SetPeerAddr set peer addr
func (p *Peer) SetPeerAddr(addr string) {
	p.peerAddr = addr
}

//SetPeerState set peer sync state
func (p *Peer) SetPeerState(state uint32) {
	p.syncState = state
}

//GetPeerState get peer sync state
func (p *Peer) GetPeerState() uint32 {
	return p.syncState
}

//GetId get peer id from peer address
func (p *Peer) GetId() uint64 {
	if p.peerId == 0 {
		addrPort := p.peerAddr + ":" + fmt.Sprint(p.servPort)
		p.peerId = Hash(addrPort)
	}

	return uint64(p.peerId)
}

//SendTo create connection and send
func (p *Peer) SendTo(buf []byte, isSync bool) error {

	conn, err := net.Dial("tcp", p.peerAddr+":"+fmt.Sprint(p.servPort))
	if err != nil {
		SuperPrint(RED_PRINT, "*ERROR* Failed to create a connection for remote server !!! err: ", err.Error())
		return err
	}

	len, err := conn.Write(buf)
	if err != nil {
		SuperPrint(RED_PRINT, "*ERROR* Failed to send data !!! len: ", len, err.Error())
		return errors.New("*ERROR* Failed to send data !!!")
	} else if len <= 0 {
		SuperPrint(RED_PRINT, "*ERROR* Failed to send data !!! len: ", len, err.Error())
		return errors.New("*ERROR* Failed to send data !!!")
	}

	conn.Close()

	return nil
}
