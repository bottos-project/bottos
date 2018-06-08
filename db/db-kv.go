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
 * file description: database interface
 * @Author: May Luo
 * @Date:   2017-12-04
 * @Last Modified by:
 * @Last Modified time:
 */

package db

import (
	"github.com/syndtr/goleveldb/leveldb/iterator"
)
//Put is to put record by key and value.
func (d *DBService) Put(key []byte, value []byte) error {
	return d.kvRepo.CallPut(key, value)
}

//Get is to get record by key and return value.
func (d *DBService) Get(key []byte) ([]byte, error) {
	return d.kvRepo.CallGet(key)
}

//Delete is to delete record by key.
func (d *DBService) Delete(key []byte) error {

	return d.kvRepo.CallDelete(key)
}

//Flush is to flush record to db.
func (d *DBService) Flush() error {
	return d.kvRepo.CallFlush()
}

//Close is to close db.
func (d *DBService) Close() {

	d.kvRepo.CallClose()
}

//Seek is to Seek db.
func (d *DBService) Seek(prefixKey []byte) ([]string, error) {

	return d.kvRepo.CallSeek(prefixKey)
}

// NewIterator is to create an iterator
func (d *DBService) NewIterator() iterator.Iterator {

	return d.kvRepo.CallNewIterator()
}
