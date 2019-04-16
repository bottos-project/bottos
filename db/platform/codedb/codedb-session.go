package codedb

import (
	"errors"

	"github.com/bottos-project/bottos/config"
	log "github.com/cihub/seelog"
)

func (m *MultindexDB) CallBeginUndo(name string) *UndoSession {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB begin undo failed ", m.db, "len(m.undoList) ", len(m.undoList))
		return nil
	}
	if name != config.PRIMARY_TRX_SESSION &&
		name != config.SUB_TRX_SESSION &&
		name != config.ADDITIONAL_BLOCK_SESSION {
		return nil
	}

	if m.session != nil && m.subsession != nil && m.sessionEx != nil {
		log.Critical("DB only support three session at one time")
		return nil
	}
	revision := m.revision
	if name != config.ADDITIONAL_BLOCK_SESSION {
		revision++
	}
	newUndoSession := &UndoSession{}
	newUndoSession.sessionObject = make(map[string]*UndoObject)
	for key, _ := range m.undoList {
		newObjectSes := NewUndoObject(key)
		newUndoSession.sessionObject[key] = newObjectSes
		newUndoSession.apply = true
	}
	newUndoSession.db = m
	if name == config.PRIMARY_TRX_SESSION {
		newUndoSession.sessionRevision = revision
		m.session = newUndoSession
		m.revision = revision
		log.Info("DB revision ", revision, "db revision ", m.revision)
		return m.session
	} else if name == config.SUB_TRX_SESSION {
		log.Info("DB sub session revision ", revision)
		m.subsession = newUndoSession
		return m.subsession
	} else if name == config.ADDITIONAL_BLOCK_SESSION {
		newUndoSession.sessionRevision = revision
		m.revision = revision
		m.sessionEx = newUndoSession
		return m.sessionEx
	}
	return nil
}

func (m *MultindexDB) CallResetSession() error {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB resetSession invalid param", m.db, "len(m.undoList) ", len(m.undoList))
		return errors.New("Invalid param")
	}
	if m.session == nil {
		return nil
	}
	if m.session.apply == true {
		m.revision--
	}
	m.session.reset()
	m.session = nil

	return nil
}
func (m *MultindexDB) CallResetSubSession() error {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB reset sub Session invalid param", m.db, "len(m.undoList) ", len(m.undoList))
		return errors.New("Invalid param")
	}
	if m.subsession == nil {
		return nil
	}
	m.subsession.reset()
	m.subsession = nil
	return nil
}

func (m *MultindexDB) CallResetSessionEx() error {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB reset reset Session ex invalid param", m.db, "len(m.undoList) ", len(m.undoList))
		return errors.New("Invalid param")
	}

	m.sessionEx.reset()
	m.sessionEx = nil
	return nil
}

func (m *MultindexDB) CallGetSession() *UndoSession {
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB get session failed", m.db, "len(m.undoList) ", len(m.undoList))
		return nil // errors.New("Invalid param")
	}
	return m.session
}
func (m *MultindexDB) CallGetSubSession() *UndoSession {
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB get sub session failed", m.db, "len(m.undoList) ", len(m.undoList))
		return nil // errors.New("Invalid param")
	}
	return m.subsession
}
func (m *MultindexDB) CallGetSessionEx() *UndoSession {
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB get sub session ex failed", m.db, "len(m.undoList) ", len(m.undoList))
		return nil // errors.New("Invalid param")
	}
	return m.sessionEx
}

func (m *MultindexDB) CallFreeSessionEx() error {
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB get free session ex failed", m.db, "len(m.undoList) ", len(m.undoList))
		return nil // errors.New("Invalid param")
	}
	m.sessionEx.apply = false
	m.sessionEx.free()
	m.sessionEx = nil
	return nil
}
func (m *MultindexDB) CallSquash() {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB squash failed", m.db, "len(m.undoList) ", len(m.undoList))
		return // errors.New("Invalid param")
	}

	if m.CallGetSession() == nil {
		return
	}
	log.Info("Squash session ", m.session)
	m.session.squash(m.subsession)
	m.subsession.free()
	m.subsession = nil
}
func (m *MultindexDB) CallPush(session *UndoSession) {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB push failed", m.db, "len(m.undoList) ", len(m.undoList))
		return // errors.New("Invalid param")
	}
	if m.session.apply == false {
		log.Error("DB invalid for apply is false")
		return
	}

	if m.CallGetSession() != session {
		log.Error("DB session is not primary session")
		return
	}
	log.Info("Push session", session)

	for k, v := range session.sessionObject {
		m.undoList[k].objectPush(v)
	}
	m.session.apply = false
	m.session.free()
	m.session = nil
}

func (m *MultindexDB) CallPushEx(session *UndoSession) {
	m.CallLock()
	defer m.CallUnLock()
	if m.db == nil || len(m.undoList) == 0 {
		log.Error("DB push ex failed", m.db, "len(m.undoList) ", len(m.undoList))
		return // errors.New("Invalid param")
	}
	if m.sessionEx.apply == false {
		log.Error("DB push ex invalid for apply is false")
		return
	}

	if m.CallGetSessionEx() != session {
		log.Error("DB session is not external session")
		return
	}

	for k, v := range session.sessionObject {

		if k != "chain_state" {
			continue
		}
		m.undoList[k].objectPush(v)
	}
	m.sessionEx.apply = false
	m.sessionEx.free()
	m.sessionEx = nil
}
