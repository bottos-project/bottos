package codedb

import log "github.com/cihub/seelog"

type UndoSession struct {
	db              *MultindexDB
	sessionObject   map[string]*UndoObject
	sessionRevision uint64
	apply           bool
}

func (u *UndoSession) rollback() {

	if u == nil || len(u.sessionObject) == 0 {
		log.Error("DB used undo session, ", u, u.sessionObject)
		return
	}
	if u.apply == false {
		log.Info("rollback apply is false  ")
		return
	}

	for _, v := range u.sessionObject {
		//log.Info("key", k)
		v.objectRollbackAll(u.db)
	}
	u.apply = false
}

func (u *UndoSession) reset() {
	if u == nil || len(u.sessionObject) == 0 {
		log.Error("DB used undo session, ", u, u.sessionObject)
		return
	}
	if u.apply == false {
		log.Info("rollback apply is false  ")
		return
	}
	u.rollback()
	u.sessionRevision = uint64(0)
	u.apply = false
}

func (u *UndoSession) squash(tmp *UndoSession) {
	if u == nil || len(u.sessionObject) == 0 {
		log.Error("DB used undo session, ", u, u.sessionObject)
		return
	}
	if u.apply == false {
		log.Info("rollback apply is false  ")
		return
	}
	if len(u.sessionObject) == 0 {
		log.Info("sessionObject len is 0")

	}
	for k, v := range u.sessionObject {
		v.objectPush(tmp.sessionObject[k])
	}
}

func (u *UndoSession) pushSessionObject(objectName string, value interface{}, objRevision uint64) {
	if u == nil || u.sessionObject == nil {
		log.Error("DB used undo session, ", u, u.sessionObject)
		return
	}
	value.(*UndoObjectValue).ItemRevision = objRevision

	obj := u.sessionObject[objectName]
	if obj == nil {
		log.Info("error used undo session should be beginundo")
		return
	} else if obj != nil {
		obj.item.Push(value)
	}
	u.sessionObject[objectName] = obj
}

func (u *UndoSession) free() {
	if u == nil || len(u.sessionObject) == 0 {
		log.Error("DB used undo session, ", u, u.sessionObject)
		return
	}
	for _, v := range u.sessionObject {
		v.objectFree()
	}
	u.sessionRevision = uint64(0)

}
