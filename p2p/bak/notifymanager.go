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
 * @Last Modified by: Stewart Li
 * @Last Modified time:2018-05-29
 */

package p2pserver

import (
	//"fmt"
	"sync"
	//"reflect"
	"github.com/AsynkronIT/protoactor-go/actor"
	"strings"
)

//NotifyManager is function to sync the trx , blk and peer info with other p2p other
type NotifyManager struct {
	p2p      *P2PServer
	stopSync chan bool

	trxActorPid      *actor.PID
	chainActorPid    *actor.PID
	producerActorPid *actor.PID

	peerMap map[uint64]*Peer
	//for reading/writing peerlist
	sync.RWMutex
}

// NewNotifyManager is to initial notify manager
func NewNotifyManager() *NotifyManager {
	return &NotifyManager{
		peerMap:          make(map[uint64]*Peer),
		trxActorPid:      nil,
		chainActorPid:    nil,
		producerActorPid: nil,
	}
}

// Start is to start notify manager
func (notify *NotifyManager) Start() {
	//for{}
}

// BroadcastByte is to broadcast data
func (notify *NotifyManager) BroadcastByte(buf []byte, isSync bool) {
	notify.RLock()
	defer notify.RUnlock()

	for _, peer := range notify.peerMap {
		if peer.GetPeerState() == ESTABLISH {
			peer.SendTo(buf, false)
		}
	}

	return
}

// AddPeer is to add an peer to local
func (notify *NotifyManager) AddPeer(peer *Peer) {
	notify.Lock()
	defer notify.Unlock()

	if _, ok := notify.peerMap[peer.GetId()]; !ok {
		notify.peerMap[peer.GetId()] = peer
	}
}

// DelPeer is to del peer from local
func (notify *NotifyManager) DelPeer(peer *Peer) {
	notify.Lock()
	defer notify.Unlock()

	if _, ok := notify.peerMap[peer.GetId()]; !ok {
		delete(notify.peerMap, peer.GetId())
	}
}

//BroadcastBlk is to sync blk info with other peer
func (notify *NotifyManager) BroadcastBlk() {
}

//SyncHash is to sync blk's hash info with other peer
func (notify *NotifyManager) SyncHash() {
}

//SyncPeer is to sync peer info with other peer
func (notify *NotifyManager) SyncPeer() {
}

// IsExist is to judge whether an addr is exist
func (notify *NotifyManager) IsExist(addr string, isExist bool) bool {
	for _, peer := range notify.peerMap {
		if res := strings.Compare(peer.peerAddr, addr); res == 0 {
			return true
		}
	}

	return false
}

// GetPeerMap is to get peer map
func (notify *NotifyManager) GetPeerMap() map[uint64]*Peer {
	return notify.peerMap
}
