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
package role

import (
	//	"encoding/json"
	"testing"
	"fmt"
	log "github.com/cihub/seelog"

	"github.com/bottos-project/bottos/db"
)

func Test_SetBinValue(t *testing.T) {
	ins := db.NewDbService("./file3", "./file3/db.db")
	err := CreateKeyValueRole(ins)
	if err != nil {
		log.Error(err)
	}
	value := make([]byte, 100)
	value = []byte{220, 0, 2, 206, 0, 0, 0, 2, 206, 0, 0, 0, 3}
	err = SetBinValue(ins, "dbtest3@usermng1111", "usermng1111", "testTableNametestKeyName4", value)
	if err != nil {
		log.Error("GetBinValue", err)
	}
	myvalue, err1 := GetBinValue(ins, "dbtest3@usermng1111", "usermng1111", "testTableNametestKeyName4")
	if err1 != nil {
		log.Error("GetBinValue", err1)
	}
	fmt.Println("myvalue ", myvalue)

	myvalue2, err2 := GetBinValue(ins, "dbtest3@usermng1111", "usermng1111", "testTableNametestKeyName4")
	if err2 != nil {
		log.Error("GetBinValue", err2)
	}
	fmt.Println("myvalue ", myvalue2)
	fmt.Println("success")

}
