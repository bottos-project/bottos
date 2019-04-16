package codedb

import (
	"errors"

	log "github.com/cihub/seelog"
)

func (m *MultindexDB) CallAddObject(object string) {
	m.CallLock()
	defer m.CallUnLock()
	m.undoList[object] = NewUndoObject(object)
}

func (m *MultindexDB) CallRollback() error {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB rollback failed ", m.db, "len(m.undoList) ", len(m.undoList))
		return errors.New("Invalid Param")
	}

	log.Info("DB begin rollback")
	for _, v := range m.undoList {
		v.objectRollback(m)
	}
	m.revision--
	return nil
}

func (m *MultindexDB) CallRollbackAll() error {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB rollbackall failed ", m.db, "len(m.undoList) ", len(m.undoList))
		return errors.New("Invalid Param")
	}
	log.Info("DB begin rollback all")
	for k, v := range m.undoList {
		log.Info(k, v)
		v.objectRollbackAll(m)
	}
	m.revision = m.commitRevision
	log.Info("DB rollbackall after revision", m.revision)
	return nil
}

func (m *MultindexDB) CallCommit(revision uint64) error {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB commit failed ", m.db, "len(m.undoList) ", len(m.undoList))
		return errors.New("Invalid Param")
	}
	//log.Info("db begin CallCommit m.revision", revision)
	for _, v := range m.undoList {
		v.objectCommit(revision)
	}
	myUndoList := make(map[string]*UndoObject)
	for key, val := range m.undoList {
		myUndoList[key] = val
	}
	m.undoList = myUndoList
	m.commitRevision = revision
	log.Info("DB CallCommit m.revision", revision)

	return nil
}

func (m *MultindexDB) CallGetRevision() uint64 {
	return m.revision
}

func (m *MultindexDB) CallSetRevision(myRevision uint64) {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB set revision failed ", m.db, "len(m.undoList) ", len(m.undoList))
		return
	}
	for _, v := range m.undoList {
		v.objectSetRevision(myRevision)
	}
	m.revision = myRevision
}

func (m *MultindexDB) PushObject(objectName string, value interface{}) {
	m.CallLock()
	defer m.CallUnLock()
	if m.undoList[objectName] == nil {
		//log.Info("Error invalid param,m.undoList not exist")
		return
	}
	undoobject := m.undoList[objectName]
	if undoobject.item == nil {
		log.Error("DB invalid param undoobject item is nil")
		return
	}
	//special for chain commit block
	if m.sessionEx != nil && objectName == "chain_state" {
		log.Info("DB sessionEx set object ", objectName, value)
		m.sessionEx.pushSessionObject(objectName, value, m.sessionEx.sessionRevision)
	} else if m.session != nil && m.subsession != nil {
		log.Info("DB subsession set object ", objectName, value)
		m.subsession.pushSessionObject(objectName, value, m.session.sessionRevision)

	} else if m.session != nil && m.subsession == nil {
		log.Info("DB session set object ", objectName, value)
		m.session.pushSessionObject(objectName, value, m.session.sessionRevision)
	}

}
