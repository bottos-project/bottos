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
	"github.com/tidwall/buntdb"
)

func (k *MultindexDB) CallCreatObjectIndex(objectName string, indexName string, indexJson string) error {

	return k.db.CreateIndex(indexName, objectName+"*", buntdb.IndexJSON(indexJson))

}

func (k *MultindexDB) CallCreatObjectMultiIndex(objectName string, indexName string, indexJson string, secKey string) error {

	return k.db.CreateIndex(indexName, objectName+"*", buntdb.IndexJSON(indexJson), buntdb.IndexJSON(secKey))
}
func (k *MultindexDB) CallDeleteObject(objectName string, key string) (string, error) {
	var objectValue string
	var err error

	k.db.Update(func(tx *buntdb.Tx) error {
		objectValue, err = tx.Delete(objectName + key)
		return err
	})

	return objectValue, err

}
