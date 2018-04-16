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

func (self *KVDatabase) Put(key []byte, value []byte) error {

	return self.db.Put(key, value, nil)
}

func (self *KVDatabase) Get(key []byte) ([]byte, error) {
	return self.db.Get(key, nil)
}

func (self *KVDatabase) Delete(key []byte) error {

	return self.db.Delete(key, nil)
}

func (self *KVDatabase) NewIterator() iterator.Iterator {
	return self.db.NewIterator(nil, nil)
}

func (self *KVDatabase) Flush() error {
	return nil
}

func (self *KVDatabase) Close() {

	self.db.Close()
	fmt.Println("flushed and closed db:", self.fn)
}
