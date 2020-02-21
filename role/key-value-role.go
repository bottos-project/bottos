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
 * file description:  key-value role
 * @Author: Gong Zibin
 * @Date:   2018-4-20
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"github.com/bottos-project/bottos/db"
	log "github.com/cihub/seelog"
)

const (
	//KeyValueObjectName is the table name of keyvalue object
	KeyValueObjectName string = "keyvalue"
)

// CreateKeyValueRole is create account role
func CreateKeyValueRole(ldb *db.DBService) error {
	return nil
}

func constructKey(contract string, object string, key string) string {
	return contract + object + key
}

//SetStrValue
func SetStrValue(ldb *db.DBService, contract string, object string, key string, value string) error {
	newkey := constructKey(contract, object, key)
	err := setKeyValueRole(ldb, newkey, value)

	log.Infof("SetStrValue, contract: %v, object: %v, key: %v, value: %v", contract, object, key, value)

	if err != nil {
		return log.Errorf("SetStrValue error, contract: %v, object: %v, key: %v, value: %v", contract, object, key, value)
	}
	return nil
}

//GetStrValue
func GetStrValue(ldb *db.DBService, contract string, object string, key string) (string, error) {
	newkey := constructKey(contract, object, key)
	value, err := getKeyValueRole(ldb, newkey)

	log.Infof("GetStrValue, contract: %v, object: %v, key: %v, value: %x", contract, object, key, value)

	if err != nil {
		return "", log.Errorf("GetStrValue error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return string(value), nil
}

//RemoveStrValue
func RemoveKeyValue(ldb *db.DBService, contract string, object string, key string) error {
	newkey := constructKey(contract, object, key)
	err := removeKeyValueRole(ldb, newkey)

	log.Infof("RemoveKeyValue, contract: %v, object: %v, key: %v", contract, object, key)

	if err != nil {
		return log.Errorf("RemoveKeyValue error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return nil
}

//SetBinValue
func SetBinValue(ldb *db.DBService, contract string, object string, key string, value []byte) error {
	newkey := constructKey(contract, object, key)
	err := setKeyValueRole(ldb, newkey, string(value))

	log.Infof("SetBinValue, contract: %v, object: %v, key: %v, value: %x", contract, object, key, value)
	if err != nil {
		return log.Errorf("SetBinValue error, contract: %v, object: %v, key: %v, value: %x", contract, object, key, value)
	}
	return nil
}

//GetBinValue
func GetBinValue(ldb *db.DBService, contract string, object string, key string) ([]byte, error) {
	newkey := constructKey(contract, object, key)
	value, err := getKeyValueRole(ldb, newkey)

	log.Infof("GetStrValue, contract: %v, object: %v, key: %v, value: %x", contract, object, key, value)

	if err != nil {
		return []byte{}, log.Errorf("GetStrValue error, contract: %v, object: %v, key: %v", contract, object, key)
	}
	return []byte(value), nil
}

func IsBinExist(ldb *db.DBService, contract string, object string, key string) bool {
	newkey := constructKey(contract, object, key)
	value, err := getKeyValueRole(ldb, newkey)

	log.Infof("IsBinExist, contract: %v, object: %v, key: %v, hex string value: %s", contract, object, key, value)

	if err != nil {
		log.Infof("IsBinExist, false")
		return false
	} else {
		log.Infof("IsBinExist, true")
		return true
	}
}

func setKeyValueRole(ldb *db.DBService, key string, value string) error {
	return ldb.SetObject(KeyValueObjectName, key, value)
}

func getKeyValueRole(ldb *db.DBService, key string) (string, error) {
	value, err := ldb.GetObject(KeyValueObjectName, key)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func removeKeyValueRole(ldb *db.DBService, key string) error {
	_, err := ldb.DeleteObject(KeyValueObjectName, key)
	return err
}
