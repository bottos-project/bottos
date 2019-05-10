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
 * file description:  transaction history role
 * @Author: Gong Zibin
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"
	"errors"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/db"
	log "github.com/cihub/seelog"
	"strconv"
)

const TransactionExpirationObjectName string = "transaction_expiration"
const TransactionExpirationObjectKeyName string = "trx_hash"
const TransactionExpirationObjectIndexExpiration string = "expiration"

//TransactionExpiration is transaction expiration struct
type TransactionExpiration struct {
	TrxHash    common.Hash `json:"trx_hash"`
	Expiration uint64      `json:"expiration"`
}

//CreateTransactionExpirationRole is creating transaction expiration role
func CreateTransactionExpirationRole(ldb *db.DBService) error {
	err := ldb.CreatObjectIndex(TransactionExpirationObjectName, TransactionExpirationObjectKeyName, TransactionExpirationObjectKeyName)
	if err != nil {
		return err
	}
	err = ldb.CreatObjectIndex(TransactionExpirationObjectName, TransactionExpirationObjectIndexExpiration, TransactionExpirationObjectIndexExpiration)
	if err != nil {
		return err
	}
	ldb.AddObject(TransactionExpirationObjectName)
	return nil
}

func hashToKey(hash common.Hash) string {
	return hash.ToHexString()
}

//SetTransactionExpirationRole is setting transaction expiration role
func SetTransactionExpirationRole(ldb *db.DBService, hash common.Hash, value *TransactionExpiration) error {
	key := hashToKey(hash)
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		log.Error("ROLE Marshal failed ", err)
		return err
	}

	log.Debugf("ROLE set exp trx:%x, exp:%v", value.TrxHash, value.Expiration)

	return ldb.SetObject(TransactionExpirationObjectName, key, string(jsonvalue))
}

//GetTransactionExpirationRoleByHash is getting transaction expiration role hash
func GetTransactionExpirationRoleByHash(ldb *db.DBService, hash common.Hash) (*TransactionExpiration, error) {
	key := hashToKey(hash)
	value, err := ldb.GetObject(TransactionExpirationObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &TransactionExpiration{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		log.Error("ROLE Unmarshal failed ", err)
		return nil, err
	}

	if res.Expiration == 0 {
		log.Errorf("ROLE null expiration record, trx: %x", hash)
		return nil, errors.New("Null expiration record")
	}

	return res, nil
}

func RemoveTransactionExpirationRoleByExpiration(ldb *db.DBService, expiration uint64) error {
	var objects []string
	var err error
	objects, err = ldb.GetObjectsWithinRangeByIndex(TransactionExpirationObjectIndexExpiration, strconv.FormatUint(expiration, 10), "1")
	if err != nil {
		log.Error("GetObjectsWithinRangeByIndex failed", err)
		return err
	}
	log.Debug("ROLE objects ", objects)
	log.Debugf("ROLE remove exp :%v", expiration)

	for _, object := range objects {
		res := &TransactionExpiration{}
		err = json.Unmarshal([]byte(object), res)
		if err != nil {
			log.Error("ROLE RemoveTransactionExpirationRoleByExpiration Unmarshal failed", err)
			return err
		}
		log.Debugf("ROLE remove expired trx: %x, expiration: %v", res.TrxHash, res.Expiration)

		// Set the record to null instead of deleting the record
		res.Expiration = 0
		err = SetTransactionExpirationRole(ldb, res.TrxHash, res)
		if err != nil {
			return err
		}
	}

	return nil
}
