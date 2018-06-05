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
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */
package contractdb

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bottos-project/bottos/db"
)

type TestParam struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

func TestCreateAndUpdateStrValue(t *testing.T) {
	ins := db.NewDbService("./file", "./file/db.db", "")
	cdb := NewContractDB(ins)

	var contract string = "testcontract"
	var object string = "testobject"
	var account string = "abc"
	tp := &TestParam{
		From:  "abc",
		To:    "bcd",
		Value: 9999,
	}
	str, _ := json.Marshal(tp)
	cdb.SetStrValue(contract, object, account, string(str))
	getStr, _ := cdb.GetStrValue(contract, object, account)
	fmt.Printf("str=%v, getStr=%v\n", string(str), getStr)

	tp.Value = 1234
	str, _ = json.Marshal(tp)
	cdb.SetStrValue(contract, object, account, string(str))
	getStr, _ = cdb.GetStrValue(contract, object, account)
	fmt.Printf("str=%v, getStr=%v\n", string(str), getStr)
}

func TestDifferentStrValue(t *testing.T) {
	ins := db.NewDbService("./file1", "./file1/db.db", "")
	cdb := NewContractDB(ins)

	var contract string = "testcontract"
	var object string = "testobject"
	tp1 := &TestParam{
		From:  "abc",
		To:    "bcd",
		Value: 9999,
	}
	tp2 := &TestParam{
		From:  "xxx",
		To:    "yyy",
		Value: 1234,
	}
	str1, err := json.Marshal(tp1)
	str2, err := json.Marshal(tp2)
	if err != nil {
		return
	}

	cdb.SetStrValue(contract, object, "abc", string(str1))
	cdb.SetStrValue(contract, object, "xxx", string(str2))
	getStr1, err := cdb.GetStrValue(contract, object, "abc")
	getStr2, err := cdb.GetStrValue(contract, object, "xxx")

	fmt.Printf("str1=%v, getStr1=%v\n", string(str1), getStr1)
	fmt.Printf("str2=%v, getStr2=%v\n", string(str2), getStr2)
}

func TestRemoveStrValue(t *testing.T) {
	ins := db.NewDbService("./file2", "./file2/db.db", "")
	cdb := NewContractDB(ins)

	var contract string = "testcontract"
	var object string = "testobject"
	var account string = "abc"
	tp := &TestParam{
		From:  "abc",
		To:    "bcd",
		Value: 9999,
	}
	str, err := json.Marshal(tp)
	if err != nil {
		return
	}

	cdb.SetStrValue(contract, object, account, string(str))
	getStr, err := cdb.GetStrValue(contract, object, account)
	fmt.Printf("str=%v, getStr=%v\n", string(str), getStr)

	fmt.Println("remove string")
	err = cdb.RemoveStrValue(contract, object, account)
	getStr, err = cdb.GetStrValue(contract, object, account)
	fmt.Println("get string error: ", err)
}
