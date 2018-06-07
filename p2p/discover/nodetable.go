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
// Package discover implements the Node Discovery Protocol.
//
// The Node Discovery protocol provides a way to find RLPx nodes that
// can be connected to. It uses a Kademlia-like protocol to maintain a
// distributed database of the IDs and endpoints of all listening
// nodes.

/*
 * file description: the interface for WASM execution
 * @Author: Richard
 * @Date:   2018-02-10
 * @Last Modified by:
 * @Last Modified time:
 */

package discover

import (
	"crypto/rand"
	"encoding/binary"
	"net"
	"sync"
	"time"
)

const (
	alpha      = 3
	bucketSize = 16
	hashBits   = 8
	nBuckets   = hashBits + 1

	maxBondingPingPongs = 16
	maxFindnodeFailures = 5
)

// Table is definition of node table
type Table struct {
	mutex         sync.Mutex
	buckets       [nBuckets]*bucket
	nursery       []*Node
	db            *nodeDB
	bondmu        sync.Mutex
	bonding       map[NodeID]*bondproc
	bondslots     chan struct{}
	nodeAddedHook func(*Node)
	net           transport
	self          *Node
}

type bondproc struct {
	err  error
	n    *Node
	done chan struct{}
}

type transport interface {
	ping(NodeID, *net.UDPAddr) error
	waitping(NodeID) error
	findnode(toid NodeID, addr *net.UDPAddr, target NodeID) ([]*Node, error)
	close()
}

type bucket struct {
	lastLookup time.Time
	entries    []*Node
}

func newTableInfo(t transport, ourID NodeID, ourAddr *net.UDPAddr, nodeDBPath string) *Table {

	db, err := newNodeDBInfo(nodeDBPath, Version, ourID)
	if err != nil {

		db, _ = newNodeDBInfo("", Version, ourID)
	}
	tab := &Table{
		net:       t,
		db:        db,
		self:      newNodeInfo(ourID, ourAddr.IP, uint64(ourAddr.Port), uint64(ourAddr.Port)),
		bonding:   make(map[NodeID]*bondproc),
		bondslots: make(chan struct{}, maxBondingPingPongs),
	}
	for i := 0; i < cap(tab.bondslots); i++ {
		tab.bondslots <- struct{}{}
	}
	for i := range tab.buckets {
		tab.buckets[i] = new(bucket)
	}
	return tab
}

// Self is to get table it self
func (tab *Table) Self() *Node {
	return tab.self
}

// ReadRandomNodes is to get random nodes
func (tab *Table) ReadRandomNodes(buf []*Node, buckets [][]*Node) (n int) {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()
	for _, b := range tab.buckets {
		if len(b.entries) > 0 {
			buckets = append(buckets, b.entries[:])
		}
	}
	if len(buckets) == 0 {
		return 0
	}

	for i := uint64(len(buckets)) - 1; i > 0; i-- {
		j := randUint(i)
		buckets[i], buckets[j] = buckets[j], buckets[i]
	}

	var i, j int
	for ; i < len(buf); i, j = i+1, (j+1)%len(buckets) {
		b := buckets[j]
		buf[i] = &(*b[0])
		buckets[j] = b[1:]
		if len(b) == 1 {
			buckets = append(buckets[:j], buckets[j+1:]...)
		}
		if len(buckets) == 0 {
			break
		}
	}
	return i + 1
}

func randUint(max uint64) uint64 {
	if max == 0 {
		return 0
	}
	var b []byte
	rand.Read(b[:])
	return binary.BigEndian.Uint64(b[:]) % max
}

// Close is to stop p2p when system stop
func (tab *Table) Close() {
	tab.net.close()
	tab.db.close()
}

// Bootstrap is to boost trap
func (tab *Table) Bootstrap(nodes []*Node) {
	tab.mutex.Lock()

	tab.nursery = make([]*Node, 0, len(nodes))
	for _, n := range nodes {
		cpy := *n
		tab.nursery = append(tab.nursery, &cpy)
	}
	tab.mutex.Unlock()
	tab.refresh()
}

// Lookup is to get node talbe by node id
func (tab *Table) Lookup(targetID NodeID) []*Node {
	var (
		target         = string(targetID[:])
		asked          = make(map[NodeID]bool)
		seen           = make(map[NodeID]bool)
		reply          = make(chan []*Node, alpha)
		pendingQueries = 0
	)
	asked[tab.self.ID] = true

	tab.mutex.Lock()

	tab.buckets[logdistance(tab.self.hash, target)].lastLookup = time.Now()

	result := tab.closest(target, bucketSize)
	tab.mutex.Unlock()

	if len(result.entries) == 0 {
		tab.refresh()
		return nil
	}

	for {

		for i := 0; i < len(result.entries) && pendingQueries < alpha; i++ {
			n := result.entries[i]
			if !asked[n.ID] {
				asked[n.ID] = true
				pendingQueries++
				go func() {

					r, err := tab.net.findnode(n.ID, n.addr(), targetID)
					if err != nil {

						fails := tab.db.failFind(n.ID) + 1
						tab.db.updateFindFails(n.ID, fails)

						if fails >= maxFindnodeFailures {

							tab.del(n)
						}
					}
					reply <- tab.bondall(r)
				}()
			}
		}
		if pendingQueries == 0 {

			break
		}

		for _, n := range <-reply {
			if n != nil && !seen[n.ID] {
				seen[n.ID] = true
				result.push(n, bucketSize)
			}
		}
		pendingQueries--
	}
	return result.entries
}

func (tab *Table) refresh() {
	seed := true

	tab.mutex.Lock()
	for _, bucket := range tab.buckets {
		if len(bucket.entries) > 0 {
			seed = false
			break
		}
	}
	tab.mutex.Unlock()

	if !seed {

		var target NodeID
		rand.Read(target[:])

		result := tab.Lookup(target)
		if len(result) == 0 {

			seed = true
		}
	}

	if seed {

		seeds := tab.db.querySeedNodes(10)

		nodes := append(tab.nursery, seeds...)

		bonded := tab.bondall(nodes)
		if len(bonded) > 0 {
			tab.Lookup(tab.self.ID)
		}

	}
}

func (tab *Table) closest(target string, nresults int) *nodesByDistance {

	close := &nodesByDistance{target: target}
	for _, b := range tab.buckets {
		for _, n := range b.entries {
			close.push(n, nresults)
		}
	}
	return close
}

func (tab *Table) len() (n int) {
	for _, b := range tab.buckets {
		n += len(b.entries)
	}
	return n
}
