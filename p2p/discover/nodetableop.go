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
	"net"
	"sort"
	"time"
)


func (tab *Table) bondall(nodes []*Node) (result []*Node) {
	rc := make(chan *Node, len(nodes))
	for i := range nodes {
		go func(n *Node) {
			nn, _ := tab.bond(false, n.ID, n.addr(), uint64(n.TCP))
			rc <- nn
		}(nodes[i])
	}
	for _ = range nodes {
		if n := <-rc; n != nil {
			result = append(result, n)
		}
	}
	return result
}


func (tab *Table) bond(pinged bool, id NodeID, addr *net.UDPAddr, tcpPort uint64) (*Node, error) {
	
	node, fails := tab.db.node(id), 0
	if node != nil {
		fails = tab.db.failFind(id)
	}
	var result error
	if node == nil || fails > 0 {
		

		tab.bondmu.Lock()
		w := tab.bonding[id]
		if w != nil {
			tab.bondmu.Unlock()
			<-w.done
		} else {
			w = &bondproc{done: make(chan struct{})}
			tab.bonding[id] = w
			tab.bondmu.Unlock()
			tab.pingpong(w, pinged, id, addr, tcpPort)
			tab.bondmu.Lock()
			delete(tab.bonding, id)
			tab.bondmu.Unlock()
		}
		result = w.err
		if result == nil {
			node = w.n
		}
	}
	if node != nil {
		tab.mutex.Lock()
		defer tab.mutex.Unlock()

		b := tab.buckets[logdistance(tab.self.hash, node.hash)]
		if !b.bump(node) {
			tab.pingreplace(node, b)
		}
		tab.db.updateFindFails(id, 0)
	}
	return node, result
}

func (tab *Table) pingpong(w *bondproc, pinged bool, id NodeID, addr *net.UDPAddr, tcpPort uint64) {
	<-tab.bondslots
	defer func() { tab.bondslots <- struct{}{} }()

	if w.err = tab.ping(id, addr); w.err != nil {
		close(w.done)
		return
	}
	if !pinged {
		
		tab.net.waitping(id)
	}
	
	w.n = newNodeInfo(id, addr.IP, uint64(addr.Port), tcpPort)
	tab.db.updNode(w.n)
	close(w.done)
}

func (tab *Table) pingreplace(new *Node, b *bucket) {
	if len(b.entries) == bucketSize {
		oldest := b.entries[bucketSize-1]
		if err := tab.ping(oldest.ID, oldest.addr()); err == nil {
			
			return
		}
	} else {
		
		b.entries = append(b.entries, nil)
	}
	copy(b.entries[1:], b.entries)
	b.entries[0] = new
	if tab.nodeAddedHook != nil {
		tab.nodeAddedHook(new)
	}
}


func (tab *Table) ping(id NodeID, addr *net.UDPAddr) error {

	tab.db.updateLastPing(id, time.Now())
	if err := tab.net.ping(id, addr); err != nil {
		return err
	}

	tab.db.updateLastPong(id, time.Now())
	tab.db.ensureExp()

	return nil
}


func (tab *Table) add(entries []*Node) {
outer:
	for _, n := range entries {
		if n.ID == tab.self.ID {

			continue
		}
		bucket := tab.buckets[logdistance(tab.self.hash, n.hash)]
		for i := range bucket.entries {
			if bucket.entries[i].ID == n.ID {

				continue outer
			}
		}
		if len(bucket.entries) < bucketSize {
			bucket.entries = append(bucket.entries, n)
			if tab.nodeAddedHook != nil {
				tab.nodeAddedHook(n)
			}
		}
	}
}

func (tab *Table) del(node *Node) {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	bucket := tab.buckets[logdistance(tab.self.hash, node.hash)]
	for i := range bucket.entries {
		if bucket.entries[i].ID == node.ID {
			bucket.entries = append(bucket.entries[:i], bucket.entries[i+1:]...)
			return
		}
	}
}

func (b *bucket) bump(n *Node) bool {
	for i := range b.entries {
		if b.entries[i].ID == n.ID {
			
			copy(b.entries[1:], b.entries[:i])
			b.entries[0] = n
			return true
		}
	}
	return false
}

type nodesByDistance struct {
	entries []*Node
	target  string
}


func (h *nodesByDistance) push(n *Node, maxElems int) {
	ix := sort.Search(len(h.entries), func(i int) bool {
		return distancecmp(h.target, h.entries[i].hash, n.hash) > 0
	})
	if len(h.entries) < maxElems {
		h.entries = append(h.entries, n)
	}
	if ix == len(h.entries) {
		
	} else {
		
		copy(h.entries[ix+1:], h.entries[ix:])
		h.entries[ix] = n
	}
}
