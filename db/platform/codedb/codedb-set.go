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
	"fmt"

	log "github.com/cihub/seelog"

	"github.com/tidwall/buntdb"
)

func (k *MultindexDB) CallSetObject(objectName string, key string, objectValue string) error {
	var mykey string
	strValue := fmt.Sprintf("%v", objectValue)

	mykey = objectName + key
	return k.commonSetObject(objectName, mykey, strValue)
}

/** commonSetObject is to set object that contract set value in db in common way.
* in this way, all the seting transaction of db should be record in undo objects
 */
func (k *MultindexDB) commonSetObject(objectName string, myKey string, strValue string) error {
	var err error

	//get old values
	var oldValue string

	err = k.db.View(func(tx *buntdb.Tx) error {
		oldValue, err = tx.Get(myKey)
		return err
	})
	if err != nil {
		oldValue = ""
	}
	//myKey = objectName +key
	undo := &UndoObjectValue{objectName,
		&DbItem{myKey, oldValue},
		&DbItem{myKey, strValue},
		k.revision}

	//set object to db
	err = k.innersetObject(myKey, strValue)
	if err != nil {
		log.Info("DB innersetObject failed", myKey, err)
		return err
	}

	//push to undo object
	k.PushObject(objectName, undo)
	return err

}

/**
only for undo
*/
func (k *MultindexDB) undoCallDeleteObject(objectName string, key string, isDbFlush bool) (string, error) {
	var objectValue string
	var err error
	var myKey string
	if isDbFlush == true {
		/** dbloadFlush is to set object that bottos process going to stop.
		* Then flush the undo object to db.
		 */
		myKey = objectName + key

	} else {
		/** rollbackRecord is to set object that when db is rollback
		* for example when A -> B ->C, when rollback, should reset C->B->A
		 */
		myKey = key
	}

	k.db.Update(func(tx *buntdb.Tx) error {
		objectValue, err = tx.Delete(myKey)
		return err
	})

	return objectValue, err

}

/**
* only for undo
 */
func (k *MultindexDB) undoCallsetObject(objectName string, key string, objectValue string, isDbFlush bool) error {
	var myKey string
	strValue := fmt.Sprintf("%v", objectValue)
	if isDbFlush == true {
		/** dbloadFlush is to set object that bottos process going to stop.
		* Then flush the undo object to db.
		 */
		myKey = objectName + key

	} else {
		/** rollbackRecord is to set object that when db is rollback
		* for example when A -> B ->C, when rollback, should reset C->B->A
		 */
		myKey = key
	}
	return k.innersetObject(myKey, strValue)
}

func (k *MultindexDB) innersetObject(key string, strValue string) error {
	var err error
	err = k.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, strValue, nil)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
