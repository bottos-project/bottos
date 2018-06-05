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

// Contains the node database, storing previously seen nodes and any collected
// metadata about them for QoS purposes.


/*
 * file description: the interface for WASM execution
 * @Author: Richard
 * @Date:   2018-02-10
 * @Last Modified by:
 * @Last Modified time:
 */

package discover

import (
	"bytes"
	"time"
	"github.com/syndtr/goleveldb/leveldb/util"
)



func (db *nodeDB) updNode(node *Node) error {
	blob := []byte("test")
	
	return db.ldb.Put(makeNodeKey(node.ID, nodeDBDiscoverRoot), blob, nil)
}


func (db *nodeDB) delNode(id NodeID) error {
	deleter := db.ldb.NewIterator(util.BytesPrefix(makeNodeKey(id, "")), nil)
	for deleter.Next() {
		if err := db.ldb.Delete(deleter.Key(), nil); err != nil {
			return err
		}
	}
	return nil
}


func (db *nodeDB) ensureExp() {
	db.sync.Do(func() { go db.exp() })
}


func (db *nodeDB) exp() {
	tick := time.Tick(nodeDBCycle)
	for {
		select {
		case <-tick:
			if err := db.expNodes(); err != nil {
		       return 
			}
		case <-db.quit:
			return
		}
	}
}


func (db *nodeDB) expNodes() error {
	threshold := time.Now().Add(-nodeDBExpiration)


	item := db.ldb.NewIterator(nil, nil)
	defer item.Release()

	for item.Next() {

		id, field := splitNodeKey(item.Key())
		if field != nodeDBDiscoverRoot {
			continue
		}

		if bytes.Compare(id[:], db.nid[:]) != 0 {
			if seen := db.lastPongOp(id); seen.After(threshold) {
				continue
			}
		}

		db.delNode(id)
	}
	return nil
}


func (db *nodeDB) lastPingOp(id NodeID) time.Time {
	return time.Unix(db.fetchNodeInt64(makeNodeKey(id, nodeDBDiscoverPing)), 0)
}

func (db *nodeDB) updateLastPing(id NodeID, instance time.Time) error {
	return db.storeNodeInt64(makeNodeKey(id, nodeDBDiscoverPing), instance.Unix())
}

func (db *nodeDB) lastPongOp(id NodeID) time.Time {
	return time.Unix(db.fetchNodeInt64(makeNodeKey(id, nodeDBDiscoverPong)), 0)
}


func (db *nodeDB) updateLastPong(id NodeID, instance time.Time) error {
	return db.storeNodeInt64(makeNodeKey(id, nodeDBDiscoverPong), instance.Unix())
}


func (db *nodeDB) failFind(id NodeID) int {
	return int(db.fetchNodeInt64(makeNodeKey(id, nodeDBDiscoverFindFails)))
}


func (db *nodeDB) updateFindFails(id NodeID, fails int) error {
	return db.storeNodeInt64(makeNodeKey(id, nodeDBDiscoverFindFails), int64(fails))
}


func (db *nodeDB) querySeedNodes(n int) []*Node {

	if db.iter == nil {
		db.iter = db.ldb.NewIterator(nil, nil)
	}

	nodes := make([]*Node, 0, n)
	for len(nodes) < n && db.iter.Next() {

		id, field := splitNodeKey(db.iter.Key())
		if field != nodeDBDiscoverRoot {
			continue
		}

		if bytes.Compare(id[:], db.nid[:]) == 0 {
			db.delNode(id)
			continue
		}
	
		if node := db.node(id); node != nil {
			nodes = append(nodes, node)
		}
	}

	if len(nodes) == 0 {
		db.iter.Release()
		db.iter = nil
	}
	return nodes
}

func (db *nodeDB) close() {
	if db.iter != nil {
		db.iter.Release()
	}
	close(db.quit)
	db.ldb.Close()
}
