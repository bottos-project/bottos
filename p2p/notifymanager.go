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
	//"fmt"
	"sync"
	//"reflect"
	"strings"
	"github.com/AsynkronIT/protoactor-go/actor"
)

//its function to sync the trx , blk and peer info with other p2p other
type NotifyManager struct {

	p2p              *P2PServer
	stopSync         chan bool

	trxActorPid      *actor.PID
	chainActorPid    *actor.PID
	producerActorPid *actor.PID

	peerMap          map[uint64]*Peer
	//for reading/writing peerlist
	sync.RWMutex
}

func NewNotifyManager() *NotifyManager {
	return &NotifyManager {
		peerMap:          make(map[uint64]*Peer),
		trxActorPid:      nil,
		chainActorPid:    nil,
		producerActorPid: nil,
	}
}

func (notify *NotifyManager) Start() {
	//fmt.Println("NotifyManager::Start")

	//for{}
}

func (notify *NotifyManager) BroadcastByte (buf []byte, isSync bool) {
	notify.RLock()
	defer notify.RUnlock()

	for _ , peer := range notify.peerMap {
		//fmt.Println("NotifyManager::BroadcastByte() - node: ",peer.conn , "node's type = ",reflect.TypeOf(peer))
		if peer.GetPeerState() == ESTABLISH {
			peer.SendTo(buf , false)
		}
	}

	return
}

func (notify *NotifyManager) AddPeer(peer *Peer) {
	notify.Lock()
	defer notify.Unlock()

	if _ , ok := notify.peerMap[peer.GetId()]; !ok {
		notify.peerMap[peer.GetId()] = peer
	}
}


//sync blk info with other peer
func (notify *NotifyManager) BroadcastBlk() {
	//fmt.Println("NotifyManager::BroadcastBlk")
}

//sync blk's hash info with other peer
func (notify *NotifyManager) SyncHash() {
	//fmt.Println("NotifyManager::SyncHash")
}

//sync peer info with other peer
func (notify *NotifyManager) SyncPeer() {
	//fmt.Println("NotifyManager::SyncPeer")
}

func (notify *NotifyManager) IsExist(addr string , isExist bool) bool {
	for _ , peer := range notify.peerMap {
		if res := strings.Compare(peer.peerAddr , addr); res == 0 {
			return true
		}
	}

	return false
}























