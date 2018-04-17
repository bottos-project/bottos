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
	"github.com/bottos-project/core/db/platform/codedb"
	"github.com/bottos-project/core/db/platform/kvdb"
	"github.com/bottos-project/core/db/platform/optiondb"
)

type DBService struct {
	optDbRepo  optiondb.OptionDbRepo
	kvRepo     kvdb.KvDBRepo
	codeDbRepo codedb.CodeDbRepo
}

func NewDbService(path string, codedbPath string) *DBService {
	kv, err := kvdb.NewKVDatabase(path)
	if err != nil {
		return nil
	}
	optiondb := optiondb.NewOptionDbRepository(path)
	codedb, err := codedb.NewCodeDbRepository(codedbPath)

	return &DBService{optDbRepo: optiondb, kvRepo: kv, codeDbRepo: codedb}
}

type DBApi interface {
	//kv database interface
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Close()
	Flush() error
	//code db interface can rollback
	CreatObject(objectName string, objectValue interface{}) error
	CreatObjectIndex(objectName string, indexName string) error
	SetObject(objectName string, objectValue interface{}) error
	SetObjectByIndex(objectName string, indexName string, indexValue interface{}, objectValue interface{}) error
	SetObjectByMultiIndexs(objectName string, indexName []string, indexValue []interface{}, objectValue interface{}) error
	GetObject(objectName string) (interface{}, error)
	GetObjectByIndex(objectName string, indexName string, indexValue interface{}) (interface{}, error)
	GetObjectByMultiIndexs(objectName string, indexName []string, indexValue []interface{}) (interface{}, error)

	//mongodb interface
	Set(string, string)
}
