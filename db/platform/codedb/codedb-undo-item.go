package codedb

import (
	"fmt"
	"encoding/json"

	log "github.com/cihub/seelog"
	"github.com/tidwall/buntdb"
)

const UndoObjectName string = "undo"
const UndoObjectKeyName string = "undo_key"
const DB_REVISION_KEY string = "dbrevision"

type UndoObjectValue struct {
	UndoObjectKey string  `json:"undo_object_key"`
	OldUndoValue  *DbItem `json:"old_undo_value"`
	NewUndoValue  *DbItem `json:"new_undo_value"`
	ItemRevision  uint64  `json:"item_revision"`
}
type DbItem struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type UndoRecord struct {
	UndoKey string           `json:"undo_key"`
	UndoVal *UndoObjectValue `json:"undo_val"`
}
type RevisionRecord struct {
	DbRevision     uint64 `json:"db_revision"`
	CommitRevision uint64 `json:"commit_revision"`
}

func (k *MultindexDB) createDBUndoObject() {
	k.CallCreatObjectIndex(UndoObjectName, UndoObjectKeyName, UndoObjectKeyName)
}

func (k *MultindexDB) setDBUndoObjectValue(object *UndoObjectValue) error {
	var mykey string
	mykey = fmt.Sprintf("%v", k.ai.Id())
	myrecord := &UndoRecord{
		UndoKey: mykey,
		UndoVal: object,
	}
	jsonvalue, err := json.Marshal(myrecord)
	if err != nil {
		return err
	}
	return k.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(UndoObjectName+mykey, string(jsonvalue), nil)
		log.Info(UndoObjectName+mykey, string(jsonvalue))
		return err
	})
}

func (k *MultindexDB) setDBRevision() error {
	value := RevisionRecord{
		DbRevision:     k.revision,
		CommitRevision: k.commitRevision,
	}
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return k.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(DB_REVISION_KEY, string(jsonvalue), nil)
		return err
	})
}

func (k *MultindexDB) getDBRevision() (*RevisionRecord, error) {
	var objectValue string
	var err error

	k.db.View(func(tx *buntdb.Tx) error {
		objectValue, err = tx.Get(DB_REVISION_KEY)
		return err
	})
	res := &RevisionRecord{}
	err = json.Unmarshal([]byte(objectValue), res)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (k *MultindexDB) getDBAllUndoObjectValue() ([]*UndoObjectValue, error) {

	var objectValue []*UndoObjectValue
	var err error

	k.db.View(func(tx *buntdb.Tx) error {
		return tx.Descend("undo_key", func(key, value string) bool {
			res := &UndoRecord{}
			err = json.Unmarshal([]byte(value), res)
			if err != nil {
				return false
			}
			//log.Info("key", key, "value ", res)
			objectValue = append(objectValue, res.UndoVal)
			return true
		})
	})
	return objectValue, nil

}

func (k *MultindexDB) deleteAllUndoObjectValue() {

	var delkeys []string
	k.db.View(func(tx *buntdb.Tx) error {
		return tx.Descend("undo_key", func(k, v string) bool {
			delkeys = append(delkeys, k)

			return true
		})
	})
	for _, key := range delkeys {
		fmt.Println(key)
		k.db.Update(func(tx *buntdb.Tx) error {
			_, err := tx.Delete(key)
			return err
		})
	}

	k.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(DB_REVISION_KEY)
		return err
	})
}
