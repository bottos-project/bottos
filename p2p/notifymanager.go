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
	"fmt"
	"sync"
	"github.com/AsynkronIT/protoactor-go/actor"
)

//its function to sync the trx , blk and peer info with other p2p other
type NotifyManager struct {

	p2p      *P2PServer

	stopSync chan bool
	pid      *actor.PID

	peerMap  map[uint64]*Peer
	//for reading/writing peerlist
	sync.RWMutex
}

func NewNotifyManager() *NotifyManager {
	return &NotifyManager {
		peerMap:    make(map[uint64]*Peer),
	}
}

func (notify *NotifyManager) Start() {
	fmt.Println("NotifyManager::Start")

	for {
		//signal from actor
		go notify.BoardcastTrx(nil , false)
		//signal from actor
		go notify.BoradcastBlk()

		go notify.SyncHash()
		go notify.SyncPeer()

		//receive
	}
}

func (notify *NotifyManager) BoardcastTrx (buf []byte, isSync bool) {
	notify.RLock()
	defer notify.RUnlock()

	for _ , node := range notify.peerMap {
		fmt.Println("node: ",node)
	}

	return
}

//sync trx info with other peer
func (notify *NotifyManager) BoradcastBlk() {
	fmt.Println("NotifyManager::BroadcastTrx")
}

//sync blk info with other peer
func (notify *NotifyManager) BroadcastBlk() {
	fmt.Println("NotifyManager::BroadcastBlk")
}

//sync blk's hash info with other peer
func (notify *NotifyManager) SyncHash() {
	fmt.Println("NotifyManager::SyncHash")
}

//sync peer info with other peer
func (notify *NotifyManager) SyncPeer() {
	fmt.Println("NotifyManager::SyncPeer")
}

























