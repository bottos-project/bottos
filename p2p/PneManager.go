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
 * file description: the interface for WASM execution
 * @Author: eripi
 * @Date:   2017-12-08
 * @Last Modified by:
 * @Last Modified time:
 */

package p2pserver

import (
	"sync"

	"github.com/bottos-project/bottos/common"
)

//PneManager peer neighbor exchange manager
type PneManager struct {
	qPeers *common.Queue
	q      sync.RWMutex

	neighbor []string
	n        sync.RWMutex
	nIndex   uint32
}

//NewPneQueue new pne queue
func NewPneQueue() *PneManager {
	return &PneManager{qPeers: common.NewQueue()}
}

//AddPnePeer add peer id
func (pne *PneManager) AddPnePeer(id uint64) {
	pne.q.Lock()
	defer pne.q.Unlock()

	pne.qPeers.Push(id)
}

//NextPnePeer pop a peer id
func (pne *PneManager) NextPnePeer() (uint64, bool) {
	pne.q.Lock()
	defer pne.q.Unlock()

	value := pne.qPeers.Pop()
	if value == nil {
		return 0, false
	}

	return value.(uint64), true
}

//AddNeighbor add new neighbor
func (pne *PneManager) AddNeighbor(n []string) {
	pne.n.Lock()
	pne.n.Unlock()

	//filter neighbor which is exist
	fn := common.Filter(n, pne.neighbor)

	for i := range fn {
		pne.neighbor = append(pne.neighbor, fn[i])
	}
}

//DelNeighbor delete neighbor
func (pne *PneManager) DelNeighbor(addr string) {
	pne.n.Lock()
	pne.n.Unlock()

	//find addr and remove
	for i := range pne.neighbor {
		if pne.neighbor[i] == addr {
			pne.neighbor = append(pne.neighbor[:i], pne.neighbor[i+1:]...)
		}
	}
}

//NextPneNeighbors get neighbors to discover
func (pne *PneManager) NextPneNeighbors() []string {
	pne.n.Lock()
	pne.n.Unlock()

	len := len(pne.neighbor)
	var index uint32
	if pne.nIndex+NEIGHBOR_DISCOVER_COUNT > uint32(len) {
		index = pne.nIndex
		pne.nIndex = 0
		return pne.neighbor[index:len]
	}

	index = pne.nIndex
	pne.nIndex += NEIGHBOR_DISCOVER_COUNT
	return pne.neighbor[index : index+NEIGHBOR_DISCOVER_COUNT]

}
