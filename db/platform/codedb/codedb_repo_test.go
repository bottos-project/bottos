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
 * file description: code database test
 * @Author: May Luo
 * @Date:   2017-12-04
 * @Last Modified by:
 * @Last Modified time:
 */
package codedb

import (
	//	"encoding/json"
	//"fmt"
	"testing"

	"github.com/tidwall/buntdb"
	log "github.com/cihub/seelog"
)

func TestCodeDbRepository_CallCreatObjectIndex(t *testing.T) {
	ins, _ := NewCodeDbRepository("./b.db")
	log.Info("abc")

	rtx, err := ins.db.Begin(true)
	if err != nil {
		log.Info("gdfddd")
	}
	rtx.CreateIndex("account_name", "account*", buntdb.IndexJSON("account_name"))
	rtx.Set("accountebc", `{"account_name":"ebc","vm_type":"123","vm_version":123,"code_version":"1","creation_date":"20171121","code":"{dfdfd,dfdfd,dfdfd}"}`, nil)
	//	rtx.Ascend("account_name", func(key, value string) bool {
	//		fmt.Printf("%s: %s\n", key, value)
	//		return true
	//	})
	//rtx.Commit()
	//	//rtx.Rollback()
	rtx.Ascend("account_name", func(key, value string) bool {
		log.Infof("ddd%s: %s\n", key, value)
		return true
	})
	value, _ := rtx.Get("accountebc")
	log.Info(value)
	//	defer func() {
	//		if err != nil {
	//			// The caller returned an error. We must rollback.
	//			_ = tx.Rollback()
	//			return
	//		}
	//		if writable {
	//			// Everything went well. Lets Commit()
	//			err = tx.Commit()
	//		} else {
	//			// read-only transaction can only roll back.
	//			err = tx.Rollback()
	//		}
	//	}()
	//	tx.funcd = true
	//	defer func() {
	//		tx.funcd = false
	//	}()
	//	err = fn(tx)

	//fmt.Printf(rtx)
	//	mapD := map[string]string{"account_name": "goood", "lettuce": "abc"}
	//	mapB, _ := json.Marshal(mapD)
	//	fmt.Println(string(mapB))
	//ins.db.CreateIndex("account_name", "account*", buntdb.IndexJSON("account_name"))
	//	fmt.Printf("gdf")
	//	//	ins.db.CreateIndex("age", "*", buntdb.IndexJSON("age"))
	//	//	ins.db.Update(func(tx *buntdb.Tx) error {
	//	//		tx.Set("account1", `{"account_name":"ebc","vm_type":"123","vm_version":123,"code_version":"1","creation_date":"20171121","code":"{dfdfd,dfdfd,dfdfd}"}`, nil)
	//	//		tx.Set("account2", `{"account_name":"abc","vm_type":"223","vm_version":123,"code_version":"1","creation_date":"20171121","code":"{dfdfd,dfdfd,dfdfd}"}`, nil)
	//	//		tx.Set("account3", `{"account_name":"fbc","vm_type":"323","vm_version":123,"code_version":"1","creation_date":"20171121","code":"{dfdfd,dfdfd,dfdfd}"}`, nil)
	//	//		//		tx.Set("2", `{"name":{"first":"Janet","last":"Prichard"},"age":47}`, nil)
	//	//		//		tx.Set("3", `{"name":{"first":"Carol","last":"Anderson"},"age":52}`, nil)
	//	//		//		tx.Set("4", `{"name":{"first":"Alan","last":"Cooper"},"age":28}`, nil)
	//	//		return nil
	//	//	})
	//	ins.db.View(func(tx *buntdb.Tx) error {
	//		fmt.Println("Order by account_name")
	//		tx.Ascend("account_name", func(key, value string) bool {
	//			fmt.Printf("%s: %s\n", key, value)
	//			return true
	//		})
	//		//		fmt.Println("Order by age")
	//		//		tx.Ascend("age", func(key, value string) bool {
	//		//			fmt.Printf("%s: %s\n", key, value)
	//		//			return true
	//		//		})
	//		//		fmt.Println("Order by age range 30-50")
	//		//		tx.AscendRange("age", `{"age":30}`, `{"age":50}`, func(key, value string) bool {
	//		//			fmt.Printf("%s: %s\n", key, value)
	//		//			return true
	//		//		})
	//		return nil
	//	})

	//	ins.db.View(func(tx *buntdb.Tx) error {
	//		fmt.Println("rollback after")
	//		tx.Ascend("account_name", func(key, value string) bool {
	//			fmt.Printf("%s: %s\n", key, value)
	//			return true
	//		})
	//		//		fmt.Println("Order by age")
	//		//		tx.Ascend("age", func(key, value string) bool {
	//		//			fmt.Printf("%s: %s\n", key, value)
	//		//			return true
	//		//		})
	//		//		fmt.Println("Order by age range 30-50")
	//		//		tx.AscendRange("age", `{"age":30}`, `{"age":50}`, func(key, value string) bool {
	//		//			fmt.Printf("%s: %s\n", key, value)
	//		//			return true
	//		//		})
	//		return nil
	//	})

}
