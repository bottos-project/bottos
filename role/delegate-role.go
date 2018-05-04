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
 * file description:  delegate role
 * @Author: May Luo
 * @Date:   2017-12-02
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"
	//"fmt"

	"github.com/bottos-project/core/db"
)

//TODO type
const DelegateObjectName string = "delegate"
const DelegateObjectKeyName string = "account_name"
const DelegateObjectIndexName string = "signing_key"

type Delegate struct {
	AccountName           string `json:"account_name"`
	LastSlot              uint64 `json:"last_slot"`
	ReportKey             string `json:"report_key"`
	TotalMissed           int64  `json:"total_missed"`
	LastConfirmedBlockNum uint32 `json:"last_confirmed_block_num"`
}

func CreateDelegateRole(ldb *db.DBService) error {
	err := ldb.CreatObjectIndex(DelegateObjectName, DelegateObjectKeyName, DelegateObjectKeyName)
	if err != nil {
		return err
	}
	err = ldb.CreatObjectIndex(DelegateObjectName, DelegateObjectIndexName, DelegateObjectIndexName)
	if err != nil {
		return err
	}
	return nil
}

func SetDelegateRole(ldb *db.DBService, key string, value *Delegate) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(DelegateObjectName, key, string(jsonvalue))
}

func GetDelegateRoleByAccountName(ldb *db.DBService, key string) (*Delegate, error) {
	value, err := ldb.GetObject(DelegateObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &Delegate{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil

}
func GetDelegateRoleBySignKey(ldb *db.DBService, keyValue string) (*Delegate, error) {

	value, err := ldb.GetObjectByIndex(DelegateObjectName, DelegateObjectIndexName, keyValue)
	if err != nil {
		return nil, err
	}

	res := &Delegate{}
	json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func GetAllDelegates(ldb *db.DBService) []*Delegate {
	objects, err := ldb.GetAllObjects(DelegateObjectName)
	if err != nil {
		return nil
	}
	var dgates = []*Delegate{}
	for _, object := range objects {
		res := &Delegate{}
		err = json.Unmarshal([]byte(object), res)
		if err != nil {
			return nil
		}
		dgates = append(dgates, res)
	}
	return dgates

}

//func GetDelegates(ldb *db.DBService) []*Delegate {
//	var listDelegate []*Delegate
//	bc.db.View(func(tx *bolt.Tx) error {
//		// Assume bucket exists and has keys
//		b := tx.Bucket([]byte(peerBucket))
//		c := b.Cursor()

//		for k, v := c.First(); k != nil; k, v = c.Next() {
//			delegate := DeserializePeer(v)
//			listDelegate = append(listDelegate, delegate)
//		}

//		return nil
//	})
//	return listDelegate
//}

//// get number delegates
//func GetNumberDelegates(bc *Blockchain) int {
//	numberDelegate := 0
//	bc.db.View(func(tx *bolt.Tx) error {
//		// Assume bucket exists and has keys
//		b := tx.Bucket([]byte(peerBucket))
//		c := b.Cursor()

//		for k, _ := c.First(); k != nil; k, _ = c.Next() {
//			numberDelegate += 1
//		}

//		return nil
//	})
//	return numberDelegate
//}

//func UpdateDelegate(bc *Blockchain, address string, lastHeight int64) {
//	var delegate Delegates
//	bc.db.Update(func(tx *bolt.Tx) error {
//		// Assume bucket exists and has keys
//		b := tx.Bucket([]byte(peerBucket))
//		delegateData := b.Get([]byte(address))
//		if delegateData == nil {
//			return errors.New("Delegates is not found.")
//		}
//		delegate = *DeserializePeer(delegateData)
//		if delegate.LastHeight < lastHeight {
//			delegate.LastHeight = lastHeight
//			b.Put([]byte(address), delegate.SerializeDelegate())
//			log.Println("updated", address, lastHeight, delegate)
//		}
//		return nil
//	})
//}

//func InsertDelegates(bc *Blockchain, delegate *Delegates, lastHeight int64) bool {
//	isInsert := false
//	err := bc.db.Update(func(tx *bolt.Tx) error {
//		// Assume bucket exists and has keys
//		b := tx.Bucket([]byte(peerBucket))

//		delegateData := b.Get([]byte(delegate.Address))
//		if delegateData == nil {
//			if delegate.LastHeight < lastHeight {
//				delegate.LastHeight = lastHeight
//			}
//			err := b.Put([]byte(delegate.Address), delegate.SerializeDelegate())
//			if err != nil {
//				log.Panic(err)
//			}
//			isInsert = true
//			return err
//		} else {
//			tmpDelegate := *DeserializePeer(delegateData)
//			if tmpDelegate.LastHeight < lastHeight {
//				delegate.LastHeight = lastHeight
//				err := b.Put([]byte(delegate.Address), delegate.SerializeDelegate())
//				if err != nil {
//					log.Panic(err)
//				}
//				isInsert = true
//				return err
//			}
//		}
//		return nil
//	})

//	if err != nil {
//		log.Panic(err)
//	}
//	return isInsert
//}
