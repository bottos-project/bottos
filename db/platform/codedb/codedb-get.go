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
 * file description: code database interface
 * @Author: May Luo
 * @Date:   2017-12-05
 * @Last Modified by:
 * @Last Modified time:
 */

package codedb

import (
	//"fmt"

	"github.com/tidwall/buntdb"
	log "github.com/cihub/seelog"
)

//CallGetObject is to get object by key
func (k *CodeDbRepository) CallGetObject(objectName string, key string) (string, error) {
	var objectValue string
	var err error

	k.db.View(func(tx *buntdb.Tx) error {
		objectValue, err = tx.Get(objectName + key)
		return err
	})

	return objectValue, err

}

//CallGetAllObjectKeys is to get all objects by objectName
func (k *CodeDbRepository) CallGetAllObjectKeys(objectName string) ([]string, error) {
	var objectValue []string
	var err error

	k.db.View(func(tx *buntdb.Tx) error {
		err = tx.Ascend(objectName, func(key, value string) bool {
			objectValue = append(objectValue, key)
			return true
		})
		return err
	})

	return objectValue, err

}

//CallGetAllObjects is to get all objects by keyName which is indexname
func (k *CodeDbRepository) CallGetAllObjects(keyName string) ([]string, error) {
	var objectValue []string
	var err error
	k.db.View(func(tx *buntdb.Tx) error {
		err = tx.Ascend(keyName, func(key, value string) bool {
			objectValue = append(objectValue, value)
			return true
		})
		return err
	})

	return objectValue, err

}

//CallGetObjectByIndex is to get object by one indexName
func (k *CodeDbRepository) CallGetObjectByIndex(objectName string, indexName string, indexValue string) (string, error) {
	var objectValue string

	log.Info(`{` + indexName + ":" + indexValue + `}`)
	err := k.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendGreaterOrEqual(indexName, `{`+indexName+":"+indexValue+`}`, func(key, value string) bool {
			objectValue = value
			log.Info(value)
			return true
		})
	})

	return objectValue, err

}

//CallGetAllObjectsSortByIndex is to get all objects by sort indexName
func (k *CodeDbRepository) CallGetAllObjectsSortByIndex(indexName string) ([]string, error) {
	var objectValue []string
	var err error

	err = k.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend(indexName, func(key, value string) bool {
			objectValue = append(objectValue, value)
			return true
		})
	})

	return objectValue, err
}
