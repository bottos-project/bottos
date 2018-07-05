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

	"github.com/bottos-project/bottos/db"
)

// DelegateObjectName is definition of delegate object name
const DelegateObjectName string = "delegate"

// DelegateObjectKeyName is definition of delegate object key name
const DelegateObjectKeyName string = "account_name"

// DelegateObjectIndexName is definition of delegate object index name
const DelegateObjectIndexName string = "signing_key"

// Delegate is definition of delegate
type Delegate struct {
	AccountName           string `json:"account_name"`
	LastSlot              uint64 `json:"last_slot"`
	ReportKey             string `json:"report_key"`
	TotalMissed           int64  `json:"total_missed"`
	LastConfirmedBlockNum uint32 `json:"last_confirmed_block_num"`
}

// CreateDelegateRole is to save initial delegate
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

// SetDelegateRole is to save delegate
func SetDelegateRole(ldb *db.DBService, key string, value *Delegate) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(DelegateObjectName, key, string(jsonvalue))
}

// GetDelegateRoleByAccountName is to get delegate by account name
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

// GetDelegateRoleBySignKey is to get delegate by sign key
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

// GetAllDelegates is to get all delegates
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

// FilterOutgoingDelegate is to filter outgoing delegate
func FilterOutgoingDelegate(ldb *db.DBService) []string {
	objects, err := ldb.GetAllObjects(DelegateObjectName)
	if err != nil {
		return nil
	}
	var accounts = make([]string, 0, len(objects))
	for _, object := range objects {
		res := &Delegate{}
		err = json.Unmarshal([]byte(object), res)
		if err != nil {
			return nil
		}
		if res.ReportKey == "xxxxxx" { //TODO
			continue
		}
		accounts = append(accounts, res.AccountName)
	}
	return accounts

}
