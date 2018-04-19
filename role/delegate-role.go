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
	"fmt"

	"github.com/bottos-project/core/db"
)

//TODO type
const DelegateObjectName string = "delegate"
const DelegateObjectKeyName string = "account_name"
const DelegateObjectIndexName string = "signing_key"

type Delegate struct {
	AccountName           string `json:"account_name"`
	LastSlot              uint64 `json:"last_slot"`
	SigningKey            string `json:"signing_key"`
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
	jsonvalue, _ := json.Marshal(value)
	return ldb.SetObject(DelegateObjectName, key, string(jsonvalue))
}

func GetDelegateRoleByAccountName(ldb *db.DBService, key string) (*Delegate, error) {
	value, err := ldb.GetObject(DelegateObjectName, key)
	res := &Delegate{}
	json.Unmarshal([]byte(value), res)
	fmt.Println("Get", key, value)
	return res, err

}
func GetDelegateRoleBySignKey(ldb *db.DBService, keyValue string) (*Delegate, error) {

	value, err := ldb.GetObjectByIndex(DelegateObjectName, DelegateObjectIndexName, keyValue)
	res := &Delegate{}
	json.Unmarshal([]byte(value), res)
	fmt.Println("Get", keyValue, value)
	return res, err
}
