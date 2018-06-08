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
	"sync"
	//"reflect"
	"strings"

	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/cihub/seelog"
)

//NotifyManager its function to sync the trx , blk and peer info with other p2p other
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

//NewNotifyManager create manager
func NewNotifyManager() *NotifyManager {
	return &NotifyManager{
		peerMap:          make(map[uint64]*Peer),
		trxActorPid:      nil,
		chainActorPid:    nil,
		producerActorPid: nil,
	}
}

//Start start
func (notify *NotifyManager) Start() {
}

func (notify *NotifyManager) broadcastByte(buf []byte, isSync bool) {

	log.Info("NotifyManager::broadcastByte")
	peerMap := notify.getPeerMap()

	for _, peer := range peerMap {
		if peer.GetPeerState() == ESTABLISH {
			peer.SendTo(buf, false)
		}
	}

	return
}

func (notify *NotifyManager) addPeer(peer *Peer) {
	notify.Lock()
	defer notify.Unlock()

	if _, ok := notify.peerMap[peer.GetId()]; !ok {
		notify.peerMap[peer.GetId()] = peer
	}
}

func (notify *NotifyManager) delPeer(peer *Peer) {
	notify.Lock()
	defer notify.Unlock()

	if _, ok := notify.peerMap[peer.GetId()]; !ok {
		delete(notify.peerMap, peer.GetId())
	}
}

func (notify *NotifyManager) getPeer(addr string) {

}

func (notify *NotifyManager) getPeerMap() map[uint64]*Peer {
	notify.RLock()
	defer notify.RUnlock()

	return notify.peerMap
}

func (notify *NotifyManager) getPeerCnt() uint32 {
	notify.RLock()
	defer notify.RUnlock()

	if len(notify.peerMap) == 0 {
		return 0
	}

	var cnt uint32
	for _, peer := range notify.peerMap {
		if ESTABLISH == peer.GetPeerState() {
			cnt++
		}
	}
	return cnt
}

//BroadcastBlk todo sync blk info with other peer
func (notify *NotifyManager) BroadcastBlk() {
	//fmt.Println("NotifyManager::BroadcastBlk")
}

//SyncHash todo sync blk's hash info with other peer
func (notify *NotifyManager) SyncHash() {
	//fmt.Println("NotifyManager::SyncHash")
}

//SyncPeer todo sync peer info with other peer
func (notify *NotifyManager) SyncPeer() {
	//fmt.Println("NotifyManager::SyncPeer")
}

func (notify *NotifyManager) isExist(addr string, isExist bool) bool {
	for _, peer := range notify.peerMap {
		if res := strings.Compare(peer.peerAddr, addr); res == 0 {
			return true
		}
	}

	return false
}

//GetPeerInfo get peer info
func (notify *NotifyManager) GetPeerInfo(id uint64) string {
	notify.Lock()
	defer notify.Unlock()

	log.Debug("GetPeerInfo")

	peer := notify.peerMap[id]
	if peer == nil {
		return ""
	}

	return peer.GetPeerAddr()

}

//GetPeersAddr get peer addr
func (notify *NotifyManager) GetPeersAddr() []string {
	notify.Lock()
	defer notify.Unlock()

	var addrs []string
	for _, value := range notify.peerMap {
		addrs = append(addrs, value.peerAddr)
	}

	return addrs
}
