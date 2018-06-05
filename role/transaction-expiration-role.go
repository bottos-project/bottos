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

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/db"
)

const transactionExpirationName string = "transaction_expiration"

//TransactionExpiration is transaction expiration struct
type TransactionExpiration struct {
	TrxHash    common.Hash `json:"trx_hash"`
	Expiration uint64      `json:"expiration"`
}

//CreateTransactionExpirationRole is creating transaction expiration role
func CreateTransactionExpirationRole(ldb *db.DBService) error {
	return nil
}

func hashToKey(hash common.Hash) string {
	return hash.ToString()
}

//SetTransactionExpirationRole is setting transaction expiration role
func SetTransactionExpirationRole(ldb *db.DBService, hash common.Hash, value *TransactionExpiration) error {
	key := hashToKey(hash)
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(transactionExpirationName, key, string(jsonvalue))
}

//GetTransactionExpirationRoleByHash is getting transaction expiration role hash
func GetTransactionExpirationRoleByHash(ldb *db.DBService, hash common.Hash) (*TransactionExpiration, error) {
	key := hashToKey(hash)
	value, err := ldb.GetObject(transactionExpirationName, key)
	if err != nil {
		return nil, err
	}

	res := &TransactionExpiration{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
