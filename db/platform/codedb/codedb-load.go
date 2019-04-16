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
	log "github.com/cihub/seelog"
)

func (m *MultindexDB) CallUndoFlush() {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB undo flush failed ", m.db, "len(m.undoList) ", len(m.undoList))
		return
	}
	m.createDBUndoObject()
	for _, v := range m.undoList {
		v.objectFlush(m)
	}
	m.setDBRevision()
}

func (m *MultindexDB) CallLoadStateDB() {
	m.CallLock()
	defer m.CallUnLock()
	m.createDBUndoObject()

	record, err := m.getDBRevision()
	if err != nil {
		log.Error("DB get Db revision without value")
		m.revision = uint64(0)
		m.commitRevision = uint64(0)
		return
	}
	m.revision = record.DbRevision
	m.commitRevision = record.CommitRevision

	undoobjects, err := m.getDBAllUndoObjectValue()
	if err != nil {
		log.Critical("DB getDBAllUndoObjectValue failed ", err)
		return
	}

	for _, obj := range undoobjects {
		if m.undoList[obj.UndoObjectKey] == nil {
			log.Critical("DB failed load object ", obj.UndoObjectKey)
			return
		}
		m.undoList[obj.UndoObjectKey].item.Push(obj)
		if m.undoList[obj.UndoObjectKey].objectRevision < obj.ItemRevision {
			m.undoList[obj.UndoObjectKey].objectRevision = obj.ItemRevision

		}

	}
	log.Info("DB revision", m.revision)
}

func (m *MultindexDB) CallReleaseUndoInfo() {
	m.deleteAllUndoObjectValue()
}