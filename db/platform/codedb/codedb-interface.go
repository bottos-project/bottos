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
 * file description: database for object
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package codedb

//CodeDbRepo is the interface for code db
type CodeDbRepo interface {
	CallCreatObjectIndex(objectName string, indexName string, indexJson string) error
	CallCreatObjectMultiIndex(objectName string, indexName string, indexJson string, secKey string) error
	CallSetObject(objectName string, key string, objectValue string) error
	CallGetObject(objectName string, key string) (string, error)
	CallGetObjectByIndex(objectName string, indexName string, indexValue string) (string, error)
	CallDeleteObject(objectName string, key string) (string, error)
	CallGetAllObjectKeys(objectName string) ([]string, error)
	CallGetAllObjects(keyName string) ([]string, error)
	CallGetAllObjectsFilter(keyName string) ([]string, error)
	CallGetAllObjectsSortByIndex(indexName string) ([]string, error)
	CallGetObjectsWithinRangeByIndex(indexName string, lessOrEqual string, greaterThan string) ([]string, error)
	CallGlobalLock()
	CallGlobalUnLock()
	CallClose()

	////db undo
	CallUndoFlush()
	CallAddObject(object string)
	CallRollback() error
	CallRollbackAll() error
	CallCommit(revision uint64) error
	CallGetRevision() uint64
	CallSetRevision(myRevision uint64)
	CallLoadStateDB()
	CallReleaseUndoInfo()

	//session undo
	CallBeginUndo(string) *UndoSession
	CallGetSession() *UndoSession
	CallGetSessionEx() *UndoSession
	CallResetSession() error
	CallResetSubSession() error
	CallFreeSessionEx() error
	CallPush(session *UndoSession)
	CallPushEx(session *UndoSession)
	CallSquash()
}
