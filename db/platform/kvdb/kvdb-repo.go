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
	

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	log "github.com/cihub/seelog"
)

//OpenFileLimit is to limiting the size of open leveldb
var OpenFileLimit = 64

//KVDatabase struct
type KVDatabase struct {
	fn string      // filename for reporting
	db *leveldb.DB // LevelDB instance
}

//NewKVDatabase is to create a new kv database
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

//CallPut is to put object by key and value
func (k *KVDatabase) CallPut(key []byte, value []byte) error {

	return k.db.Put(key, value, nil)
}

//CallGet is to get object by key and return value
func (k *KVDatabase) CallGet(key []byte) ([]byte, error) {
	return k.db.Get(key, nil)
}

//CallDelete is to delete object by key
func (k *KVDatabase) CallDelete(key []byte) error {

	return k.db.Delete(key, nil)
}

//CallNewIterator is to interate object
func (k *KVDatabase) CallNewIterator() iterator.Iterator {
	return k.db.NewIterator(nil, nil)
}

//CallNewIteratorPrefix is to iterate prefix
func (k *KVDatabase) CallNewIteratorPrefix() iterator.Iterator {
	return k.db.NewIterator(nil, nil)
}

//CallFlush is to flush object
func (k *KVDatabase) CallFlush() error {

	return nil
}

//CallClose is to close object
func (k *KVDatabase) CallClose() {

	k.db.Close()
	log.Info("flushed and closed db:", k.fn)
}

//CallSeek is to seek object
func (k *KVDatabase) CallSeek(prefixKey []byte) ([]string, error) {
	var valueList []string
	iter := k.db.NewIterator(util.BytesPrefix(prefixKey), nil)
	for iter.Next() {
		//ptrKey := iter.Key()
		value := iter.Value()
		log.Infof("CallSeek: %x\n", value)
		valueList = append(valueList, string(value))
		log.Info("CallSeek1: ", valueList)
	}
	iter.Release()
	err := iter.Error()
	return valueList, err
}
