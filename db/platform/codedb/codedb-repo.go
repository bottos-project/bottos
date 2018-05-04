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
	"errors"
	"fmt"

	"github.com/tidwall/buntdb"
)

type CodeDbRepository struct {
	fn string     // filename for reporting
	db *buntdb.DB // LevelDB instance
	tx *buntdb.Tx
}

func NewCodeDbRepository(file string) (*CodeDbRepository, error) {
	codedb, err := buntdb.Open(file)
	if err != nil {
		return nil, err
	}
	return &CodeDbRepository{
		fn: file,
		db: codedb,
	}, nil
}
func (k *CodeDbRepository) CallStartUndoSession(writable bool) {
	k.tx, _ = k.db.Begin(true)
}

func (k *CodeDbRepository) CallCreatObjectIndex(objectName string, indexName string, indexJson string) error {
	if k.tx == nil {

		return k.db.CreateIndex(indexName, objectName+"*", buntdb.IndexJSON(indexJson))
	}

	return k.tx.CreateIndex(indexName, objectName+"*", buntdb.IndexJSON(indexJson))
}
func (k *CodeDbRepository) CallCreatObjectMultiIndexs(objectName string, indexName string, indexJson string) error {
	if k.tx == nil {
		return k.db.CreateIndex(indexName, objectName+"*", buntdb.IndexJSON(indexJson))
	}

	return k.tx.CreateIndex(indexName, objectName+"*", buntdb.IndexJSON(indexJson))
}
func (k *CodeDbRepository) CallSetObject(objectName string, key string, objectValue string) error {
	strValue := fmt.Sprintf("%v", objectValue)
	if k.tx == nil {
		return k.db.Update(func(tx *buntdb.Tx) error {

			_, _, err := tx.Set(objectName+key, strValue, nil)
			return err
		})
	}
	_, _, err := k.tx.Set(objectName+key, strValue, nil)
	return err
}

func (k *CodeDbRepository) CallDeleteObject(objectName string, key string) (string, error) {
	var objectValue string
	var err error

	k.db.Update(func(tx *buntdb.Tx) error {
		objectValue, err = tx.Delete(objectName + key)
		return err
	})

	return objectValue, err

}

func (k *CodeDbRepository) CallCommit() error {
	if k.tx == nil {
		fmt.Println("tx is not start undo session")
		return errors.New("tx is not start undo session")
	}
	return k.tx.Commit()
}
func (k *CodeDbRepository) CallRollback() error {
	if k.tx == nil {
		fmt.Println("tx is not start undo session")
		return errors.New("tx is not start undo session")
	}
	return k.tx.Rollback()
}
