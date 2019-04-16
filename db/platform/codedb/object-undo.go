package codedb

import (
	log "github.com/cihub/seelog"
	"github.com/golang-collections/collections/stack"
)

type UndoObject struct {
	objectName     string
	item           *stack.Stack
	objectRevision uint64
}

func NewUndoObject(myobject string) *UndoObject {
	item := stack.New()
	return &UndoObject{objectName: myobject, item: item, objectRevision: uint64(0)}
}

func (u *UndoObject) objectFlush(mdb *MultindexDB) {

	for u.item.Len() != 0 {
		head := u.item.Peek()
		object := head.(*UndoObjectValue)
		log.Info("flush", object)
		mdb.setDBUndoObjectValue(object)
		u.item.Pop()
	}

}

func (u *UndoObject) objectRollback(mdb *MultindexDB) {
	if u.item.Len() == 0 {
		return
	}
	for u.item.Len() != 0 {
		head := u.item.Peek()
		object := head.(*UndoObjectValue)
		if object.ItemRevision < mdb.revision {
			break
		}

		if object.OldUndoValue.Val == "" {
			log.Info("DB undo delete ", object.ItemRevision, " ", u.objectName, object.NewUndoValue.Key)
			mdb.undoCallDeleteObject(u.objectName, object.NewUndoValue.Key, false)
		} else {
			log.Info("DB undo reset ", object.ItemRevision, " ", u.objectName, " ", object.OldUndoValue.Key, " ", object.OldUndoValue.Val)
			err := mdb.undoCallsetObject(u.objectName, object.OldUndoValue.Key, object.OldUndoValue.Val, false)
			if err != nil {
				log.Info("undo reset failed")
			}
		}
		u.item.Pop()
	}

}
func (u *UndoObject) objectRollbackAll(mdb *MultindexDB) {
	for u.item.Len() != 0 {
		head := u.item.Peek()
		object := head.(*UndoObjectValue)
		if object.OldUndoValue.Val == "" {
			log.Info("undo delete ", object.ItemRevision, " ", u.objectName, object.NewUndoValue.Key)
			mdb.undoCallDeleteObject(u.objectName, object.NewUndoValue.Key, false)
		} else {
			log.Info("undo reset ", object.ItemRevision, " ", u.objectName, " ", object.OldUndoValue.Key, " ", object.OldUndoValue.Val)
			err := mdb.undoCallsetObject(u.objectName, object.OldUndoValue.Key, object.OldUndoValue.Val, false)
			if err != nil {
				log.Info("undo reset failed")
			}
		}
		u.item.Pop()
	}

}

func (u *UndoObject) objectCommit(revision uint64) {
	if u.item.Len() == 0 {
		return
	}

	var pre = &UndoObject{}
	pre = u

	tmpObject := NewUndoObject(u.objectName)

	for pre.item.Len() != 0 {
		head := pre.item.Peek()
		object := head.(*UndoObjectValue)
		if object.ItemRevision > revision {
			tmpObject.item.Push(object)
		}
		pre.item.Pop()
	}

	for u.item.Len() != 0 {
		u.item.Pop()
	}
	for tmpObject.item.Len() != 0 {
		head := tmpObject.item.Peek()
		leftobject := head.(*UndoObjectValue)
		u.item.Push(leftobject)
		tmpObject.item.Pop()
	}
}

func (u *UndoObject) objectSetRevision(revision uint64) {
	if u.item.Len() != 0 {
		log.Error("DB set revision when has undo objects")
		return
	}
	u.objectRevision = revision
}

func (u *UndoObject) objectFree() {
	for u.item.Len() != 0 {
		u.item.Pop()
	}
	u.objectRevision = uint64(0)
}

func (u *UndoObject) objectPush(sesObj *UndoObject) {
	if u == nil || sesObj == nil {

		return
	}

	if u.objectName != sesObj.objectName {
		return
	}
	var size int
	size = sesObj.item.Len()
	if size == 0 {
		return
	}
	tagItem := NewUndoObject(u.objectName)
	log.Info("object push", sesObj.objectName, "size", size)
	for sesObj.item.Len() != 0 {
		prepeek := sesObj.item.Peek()
		preObj := prepeek.(*UndoObjectValue)
		tagItem.item.Push(preObj)
		sesObj.item.Pop()
	}
	//	log.Info("object push", sesObj.objectName, "tagItem len", tagItem.item.Len())
	if tagItem.item.Len() != size {
		log.Error("DB invalid tag item")
		return
	}
	//	log.Info("object push", sesObj.objectName, "tagItem", tagItem)
	for tagItem.item.Len() != 0 {
		peek := tagItem.item.Peek()
		obj := peek.(*UndoObjectValue)
		//		log.Info("object push", sesObj.objectName, "value", obj)
		u.item.Push(obj)
		tagItem.item.Pop()
	}
}
