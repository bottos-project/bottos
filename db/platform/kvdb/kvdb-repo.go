// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description: database for key-value
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package kvdb

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var OpenFileLimit = 64

type KVDatabase struct {
	fn string      // filename for reporting
	db *leveldb.DB // LevelDB instance
}

func NewKVDatabase(file string) (*KVDatabase, error) {
	// open a kvdatabase
	db, err := leveldb.OpenFile(file, &opt.Options{OpenFilesCacheCapacity: OpenFileLimit})

	if err != nil {
		return nil, err
	}
	return &KVDatabase{
		fn: file,
		db: db,
	}, nil
}

func (k *KVDatabase) CallPut(key []byte, value []byte) error {

	return k.db.Put(key, value, nil)
}

func (k *KVDatabase) CallGet(key []byte) ([]byte, error) {
	return k.db.Get(key, nil)
}

func (k *KVDatabase) CallDelete(key []byte) error {

	return k.db.Delete(key, nil)
}

func (k *KVDatabase) CallNewIterator() iterator.Iterator {
	return k.db.NewIterator(nil, nil)
}

func (k *KVDatabase) CallFlush() error {

	return nil
}

func (k *KVDatabase) CallClose() {

	k.db.Close()
	fmt.Println("flushed and closed db:", k.fn)
}
