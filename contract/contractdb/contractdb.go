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
	log "github.com/cihub/seelog"
)

//ContractObjectNamePrefix is the contract object name prefix
const ContractObjectNamePrefix string = "obj-"

//ContractDB is struct for contract
type ContractDB struct {
	Db *db.DBService
}

//NewContractDB create a new contractDB
func NewContractDB(db *db.DBService) *ContractDB {
	cdb := &ContractDB{Db: db}

	return cdb
}

//SetStrValue is to create a record if key not exsit
func (cdb *ContractDB) SetStrValue(contract string, object string, key string, value string) error {
	objName := cdb.getObjectName(contract, object)
	err := cdb.setStrValue(objName, key, value)

	log.Info("SetStrValue: ", contract, object, key, objName, value)
	if err != nil {
		return fmt.Errorf("SetStr error, contract: %v, object: %v, key: %v, value: %v", contract, object, key, value)
	}
	return nil
}

//GetStrValue is to get contract by object and return contract value
func (cdb *ContractDB) GetStrValue(contract string, object string, key string) (string, error) {
	objName := cdb.getObjectName(contract, object)
	value, err := cdb.getStrValue(objName, key)

	log.Info("GetStrValue: ", contract, object, key, objName, value)

	if err != nil {
		return "", fmt.Errorf("GetStr error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return string(value), nil
}

//RemoveStrValue is to Remove contrace value by object
func (cdb *ContractDB) RemoveStrValue(contract string, object string, key string) error {
	objName := cdb.getObjectName(contract, object)
	err := cdb.removeStrValue(objName, key)

	log.Info("RemoveStrValue: ", contract, object, key, objName)

	if err != nil {
		return fmt.Errorf("RemoveStr error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return nil
}

//SetBinValue is to create a record if key not exsit
func (cdb *ContractDB) SetBinValue(contract string, object string, key string, value []byte) error {
	objName := cdb.getObjectName(contract, object)
	err := cdb.setBinValue(objName, key, value)

	log.Info("SetBinValue: ", contract, object, key, objName, value)
	if err != nil {
		return fmt.Errorf("SetBin error, contract: %v, object: %v, key: %v, value: %x", contract, object, key, value)
	}
	return nil
}

//GetBinValue is to get contract by object and return contract value
func (cdb *ContractDB) GetBinValue(contract string, object string, key string) ([]byte, error) {
	objName := cdb.getObjectName(contract, object)
	value, err := cdb.getBinValue(objName, key)

	log.Info("GetBinValue: ", contract, object, key, objName, value)

	if err != nil {
		return []byte{}, fmt.Errorf("GetBin error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return value, nil
}

//RemoveBinValue is to Remove contrace value by object
func (cdb *ContractDB) RemoveBinValue(contract string, object string, key string) error {
	objName := cdb.getObjectName(contract, object)
	err := cdb.removeBinValue(objName, key)

	log.Info("RemoveBinValue: ", contract, object, key, objName)

	if err != nil {
		return fmt.Errorf("RemoveBin error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return nil
}

func (cdb *ContractDB) getObjectName(contract string, object string) string {
	return ContractObjectNamePrefix + contract + object
}
