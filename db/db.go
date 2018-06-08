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
	//"fmt"

	"github.com/bottos-project/bottos/db/platform/codedb"
	"github.com/bottos-project/bottos/db/platform/kvdb"
	"github.com/bottos-project/bottos/db/platform/optiondb"
	log "github.com/cihub/seelog"
)

//DBService is db service struct for other package
type DBService struct {
	kvRepo    kvdb.KvDBRepo
	codeRepo  codedb.CodeDbRepo
	optDbRepo optiondb.OptionDbRepo
}

//NewDbService is to create a new db service with kv databse, codedb, and option db for optionally
func NewDbService(path string, codedbPath string, optPath string) *DBService {
	kv, err := kvdb.NewKVDatabase(path)
	if err != nil {
		return nil
	}
	db, err := codedb.NewCodeDbRepository(codedbPath)
	log.Info(err)
	if optPath == "" {
		return &DBService{kvRepo: kv, codeRepo: db, optDbRepo: nil}
	}
	optiondb := optiondb.NewOptionDbRepository(optPath)
	return &DBService{kvRepo: kv, codeRepo: db, optDbRepo: optiondb}

}

//DBApi is listing all the interface that DBService provides.
type DBApi interface {
	Lock()
	UnLock()
	//kv database interface
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	Close()
	Flush() error
	//code db interface can rollback
	StartUndoSession()
	CreatObjectIndex(objectName string, indexName string) error
	SetObject(objectName string, objectValue interface{}) error
	SetObjectByIndex(objectName string, indexName string, indexValue interface{}, objectValue interface{}) error
	SetObjectByMultiIndexs(objectName string, indexName []string, indexValue []interface{}, objectValue interface{}) error
	GetObject(objectName string, key string) (interface{}, error)
	GetObjectByIndex(objectName string, indexName string, indexValue string) (interface{}, error)
	GetObjectByMultiIndexs(objectName string, indexName []string, indexValue []interface{}) (interface{}, error)
	GetAllObjectKeys(objectName string) ([]string, error)
	GetAllObjects(keyName string) ([]string, error)
	GetAllObjectsSortByIndex(indexName string) ([]string, error)
	DeleteObject(objectName string, key string) (string, error)
	Commit()
	Rollback()
	Reset() //TODO

	//optiondb interface
	Insert(collection string, value interface{}) error
	Find(collection string, key string, value interface{}) (interface{}, error)
}
