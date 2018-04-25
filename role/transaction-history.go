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
	 //"fmt"
 
	 "github.com/bottos-project/core/db"
	 "github.com/bottos-project/core/common"
 )
 
 const TransactionHistoryObjectName string = "transaction_history"
  
 type TransactionHistory struct {
	TxHash			common.Hash		`json:"trx_hash"`
	BlockHash		common.Hash		`json:"block_hash"`
 }
   
 func CreateTransactionHistoryObjectRole(ldb *db.DBService) error {
	 return nil
 }

 func AddTransactionHistory(ldb *db.DBService, txhash common.Hash, blockhash common.Hash) error {
	value := &TransactionHistory{
		TxHash: txhash,
		BlockHash: blockhash,
	}
	return SetTransactionHistoryObjectRole(ldb, txhash, value)
 }
 
 func SetTransactionHistoryObjectRole(ldb *db.DBService, txhash common.Hash, value *TransactionHistory) error {
	 key := hashToKey(txhash)
	 jsonvalue, _ := json.Marshal(value)
	 return ldb.SetObject(TransactionHistoryObjectName, key, string(jsonvalue))
 }

 func GetBlockHashByTxHash(ldb *db.DBService, txhash common.Hash) (common.Hash, error) {
	 history, err := GetTransactionHistoryByHash(ldb, txhash)
	 if err != nil {
		 return common.Hash{}, err
	 }
	 return history.BlockHash, nil
 }
 
func GetTransactionHistoryByHash(ldb *db.DBService, hash common.Hash) (*TransactionHistory, error) {
	key := hash.ToHexString()
	//fmt.Println("GetTransactionHistoryByHash key: ", key)
	value, err := ldb.GetObject(TransactionHistoryObjectName, key)
	if err != nil {
		return nil, err
	}
	res := &TransactionHistory{}
	json.Unmarshal([]byte(value), res)
	return res, err
}
   