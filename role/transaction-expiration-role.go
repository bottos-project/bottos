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
	_"fmt"

	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/common"
)

const TransactionExpirationObjectName string = "transaction_expiration"
 
type TransactionExpiration struct {
	TrxHash			common.Hash		`json:"trx_hash"`
	Expiration		uint64			`json:"expiration"`
}
  
func CreateTransactionExpirationObjectRole(ldb *db.DBService) error {
	return nil
}

func hashToKey(hash common.Hash) string {
	return hash.ToString()
}

func SetTransactionExpirationObjectRole(ldb *db.DBService, hash common.Hash, value *TransactionExpiration) error {
	key := hashToKey(hash)
	jsonvalue, _ := json.Marshal(value)
	return ldb.SetObject(TransactionExpirationObjectName, key, string(jsonvalue))
}

func GetTransactionExpirationObjectByHash(ldb *db.DBService, hash common.Hash) (*TransactionExpiration, error) {
	key := hashToKey(hash)
	value, err := ldb.GetObject(TransactionExpirationObjectName, key)
	res := &TransactionExpiration{}
	json.Unmarshal([]byte(value), res)
	//fmt.Println("Get", key, value)
	return res, err
}
  