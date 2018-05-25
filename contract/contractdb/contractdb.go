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
 * file description:  contract db
 * @Author: Gong Zibin
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */

package contractdb

import (
	"fmt"
	"github.com/bottos-project/bottos/db"
)

const ContractObjectNamePrefix string = "obj-"

type ContractDB struct {
	Db *db.DBService
}

func NewContractDB(db *db.DBService) *ContractDB {
	cdb := &ContractDB{Db:db}

	return cdb
}

// create a record if key not exsit
func (cdb *ContractDB) SetStrValue(contract string, object string, key string, value string) error {
	objName := cdb.getObjectName(contract, object)
	err := cdb.setStrValue(objName, key, value)

	fmt.Println("SetStrValue: ", contract, object, key, objName, value)
	if err != nil {
		return fmt.Errorf("SetStr error, contract: %v, object: %v, key: %v, value: %v", contract, object, key, value)
	}
	return nil
}

func (cdb *ContractDB) GetStrValue(contract string, object string, key string) (string, error) {
	objName := cdb.getObjectName(contract, object)
	value, err := cdb.getStrValue(objName, key)

	fmt.Println("GetStrValue: ", contract, object, key, objName, value)

	if err != nil {
		return "", fmt.Errorf("GetStr error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return string(value), nil
}

func (cdb *ContractDB) RemoveStrValue(contract string, object string, key string) error {
	objName := cdb.getObjectName(contract, object)
	err := cdb.removeStrValue(objName, key)

	fmt.Println("RemoveStrValue: ", contract, object, key, objName)

	if err != nil {
		return fmt.Errorf("RemoveStr error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return nil
}

func (cdb *ContractDB) getObjectName(contract string, object string) string {
	return ContractObjectNamePrefix + contract + object;
}
