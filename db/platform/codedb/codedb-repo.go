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
	"sync"

	log "github.com/cihub/seelog"
	"github.com/tidwall/buntdb"
)

//CodeDbRepository is to build code db
type MultindexDB struct {
	fn             string     // filename for reporting
	db             *buntdb.DB // LevelDB instance
	undoFlag       bool
	signal         sync.RWMutex // the gatekeeper for all fields
	globalSignal   sync.Mutex   // global signal for all state
	undoList       map[string]*UndoObject
	session        *UndoSession
	subsession     *UndoSession
	sessionEx      *UndoSession
	revision       uint64
	commitRevision uint64
	ai             *AutoInc
}

//NewCodeDbRepository is to create new code db
func NewMultindexDB(file string) (*MultindexDB, error) {
	codedb, err := buntdb.Open(file)
	if err != nil {
		log.Error("DB open code database failed", file)
		return nil, errors.New("buntdb open file failed")
	}
	mdb := &MultindexDB{
		fn:       file,
		db:       codedb,
		undoFlag: false,
	}
	mdb.undoList = make(map[string]*UndoObject)
	mdb.session = nil
	mdb.subsession = nil
	mdb.sessionEx = nil
	mdb.revision = uint64(0)
	mdb.commitRevision = uint64(0)
	mdb.ai = New(10000000, 1)
	return mdb, nil
}
func (m *MultindexDB) CallClose() {
	m.db.Close()
	m.undoList = nil
	m.session = nil
	m.subsession = nil
	m.sessionEx = nil
	m.revision = uint64(0)
	m.commitRevision = uint64(0)
	m.ai = nil

}

func (m *MultindexDB) CallGlobalLock() {
	m.globalSignal.Lock()
}

func (m *MultindexDB) CallGlobalUnLock() {
	m.globalSignal.Unlock()
}

func (m *MultindexDB) CallLock() {
	m.signal.Lock()
}

func (m *MultindexDB) CallUnLock() {
	m.signal.Unlock()
}
