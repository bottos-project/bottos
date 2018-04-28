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
 * file description:  account role
 * @Author: Gong Zibin
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"
	_"fmt"

	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/common"
)
 
const (
	AccountObjectName string = "account"
)

type Account struct {
	AccountName			string			`json:"account_name"`
	PublicKey			[]byte			`json:"public_key"`
	VMType				byte			`json:"vm_type"`
	VMVersion           byte			`json:"vm_version"`
	CodeVersion			common.Hash		`json:"code_version"`
	CreateTime 			uint64			`json:"create_date"`
	ContractCode		[]byte			`json:"contract_code"`
	ContractAbi			[]byte			`json:"abi"`
}
 
func CreateAccountRole(ldb *db.DBService) error {
	return nil
}

func accountNameToKey(name string) string {
	return name
}

func SetAccountRole(ldb *db.DBService, accountName string, value *Account) error {
	key := accountNameToKey(accountName)
	jsonvalue, _ := json.Marshal(value)
	return ldb.SetObject(AccountObjectName, key, string(jsonvalue))
}

func GetAccountRole(ldb *db.DBService, accountName string) (*Account, error) {
	key := accountNameToKey(accountName)
	value, err := ldb.GetObject(AccountObjectName, key)
	res := &Account{}
	json.Unmarshal([]byte(value), res)
	//fmt.Println("Get", key, value)
	return res, err
}
 