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
 * file description: account role test
 * @Author: Gong Zibin
 * @Date:   2017-12-13
 * @Last Modified by:
 * @Last Modified time:
 */
package role

import (
	//	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/db"
)

func TestAccount_writedb(t *testing.T) {
	ins := db.NewDbService("./file", "./file/db.db", "")
	err := CreateAccountRole(ins)
	if err != nil {
		fmt.Println(err)
	}
	value1 := &Account{
		AccountName:  "account1",
		PublicKey:    []byte("7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF"),
		VMType:       0,
		VMVersion:    1,
		CodeVersion:  common.StringToHash("26427d49aca564c5774724de0e0b2eb1a7c4f303e73ce7dcef3b52e1ab95cc4b"),
		CreateTime:   uint64(time.Now().Unix()),
		ContractCode: []byte{},
		ContractAbi:  []byte{},
	}
	value2 := &Account{
		AccountName:  "account2",
		PublicKey:    []byte("7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF"),
		VMType:       0,
		VMVersion:    1,
		CodeVersion:  common.StringToHash("26427d49aca564c5774724de0e0b2eb1a7c4f303e73ce7dcef3b52e1ab95cc4b"),
		CreateTime:   uint64(time.Now().Unix() + 1000),
		ContractCode: []byte{},
		ContractAbi:  []byte{},
	}
	err = SetAccountRole(ins, value1.AccountName, value1)
	if err != nil {
		t.Fatal("SetAccountRole", err)
	}
	err = SetAccountRole(ins, value2.AccountName, value2)
	if err != nil {
		t.Fatal("SetAccountRole", err)
	}

	var value *Account
	value, err = GetAccountRole(ins, value1.AccountName)
	if err != nil {
		t.Fatal("GetAccountRoleByName", err)
	}

	if value.AccountName != value1.AccountName {
		t.Fatal("Account Name error")
	}
	fmt.Println(value)
}
