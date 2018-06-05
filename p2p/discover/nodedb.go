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

// Contains the node database.


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
	"encoding/binary"
	"os"
	"sync"
	"time"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"

)




type nodeDB struct {
	ldb    *leveldb.DB       
	iter iterator.Iterator 
	nid NodeID 
	sync sync.Once     
	quit   chan struct{} 
}

var (
	nodeDBID      = NodeID{}       
	nodeDBExpiration = 24 * time.Hour 
	nodeDBCycle   = time.Hour  
	
	nodeDBVersionKey = []byte("version") 
	nodeDBItemPrefix = []byte("n:")     

	nodeDBDiscoverRoot      = ":discover"
	nodeDBDiscoverPing      = nodeDBDiscoverRoot + ":lastpingop"
	nodeDBDiscoverPong      = nodeDBDiscoverRoot + ":lastpongop"
	nodeDBDiscoverFindFails = nodeDBDiscoverRoot + ":fail"
)


func newNodeDBInfo(path string, version int, nid NodeID) (*nodeDB, error) {
	if path == "" {
		return newMemoryNodeDBInfo(nid)
	}
	return newPersistentNodeDBInfo(path, version, nid)
}


func newMemoryNodeDBInfo(nid NodeID) (*nodeDB, error) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}
	return &nodeDB{
		ldb:  db,
		nid: nid,
		quit: make(chan struct{}),
	}, nil
}


func newPersistentNodeDBInfo(path string, version int, nid NodeID) (*nodeDB, error) {
	opts := &opt.Options{OpenFilesCacheCapacity: 5}
	db, err := leveldb.OpenFile(path, opts)
	if _, iscorrupted := err.(*errors.ErrCorrupted); iscorrupted {
		db, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		return nil, err
	}

	currentVer := make([]byte, binary.MaxVarintLen64)
	currentVer = currentVer[:binary.PutVarint(currentVer, int64(version))]

	blob, err := db.Get(nodeDBVersionKey, nil)
	switch err {
	case leveldb.ErrNotFound:
		if err := db.Put(nodeDBVersionKey, currentVer, nil); err != nil {
			db.Close()
			return nil, err
		}

	case nil:
		if !bytes.Equal(blob, currentVer) {
			db.Close()
			if err = os.RemoveAll(path); err != nil {
				return nil, err
			}
			return newPersistentNodeDBInfo(path, version, nid)
		}
	}
	return &nodeDB{
		ldb:  db,
		nid: nid,
		quit: make(chan struct{}),
	}, nil
}


func makeNodeKey(id NodeID, field string) []byte {
	if bytes.Equal(id[:], nodeDBID[:]) {
		return []byte(field)
	}
	return append(nodeDBItemPrefix, append(id[:], field...)...)
}


func splitNodeKey(key []byte) (id NodeID, field string) {

	if !bytes.HasPrefix(key, nodeDBItemPrefix) {
		return NodeID{}, string(key)
	}

	item := key[len(nodeDBItemPrefix):]
	copy(id[:], item[:len(id)])
	field = string(item[len(id):])

	return id, field
}


func (db *nodeDB) fetchNodeInt64(key []byte) int64 {
	blob, err := db.ldb.Get(key, nil)
	if err != nil {
		return 0
	}
	val, read := binary.Varint(blob)
	if read <= 0 {
		return 0
	}
	return val
}


func (db *nodeDB) storeNodeInt64(key []byte, n int64) error {
	blob := make([]byte, binary.MaxVarintLen64)
	blob = blob[:binary.PutVarint(blob, n)]

	return db.ldb.Put(key, blob, nil)
}


func (db *nodeDB) node(id NodeID) *Node {
	node := new(Node)
	return node
}
