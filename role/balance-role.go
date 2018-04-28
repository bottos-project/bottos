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
 * file description:  balance role
 * @Author: Gong Zibin
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"
	"fmt"

	"github.com/bottos-project/core/db"
)

const BalanceObjectName string = "balance"
const StakedBalanceObjectName string = "staked_balance"
 
type Balance struct {
	AccountName		string			`json:"account_name"`
	Balance			uint64			`json:"balance"`
}

type StakedBalance struct {
	AccountName			string			`json:"account_name"`
	StakedBalance		uint64			`json:"staked_balance"`
	
	// TODO
}

func CreateBalanceRole(ldb *db.DBService) error {
	return nil
}

func SetBalanceRole(ldb *db.DBService, accountName string, value *Balance) error {
	key := accountName
	jsonvalue, _ := json.Marshal(value)
	return ldb.SetObject(BalanceObjectName, key, string(jsonvalue))
}

func GetBalanceRole(ldb *db.DBService, accountName string) (*Balance, error) {
	key := accountName
	value, err := ldb.GetObject(BalanceObjectName, key)
	res := &Balance{}
	json.Unmarshal([]byte(value), res)
	fmt.Println("Get", key, value)
	return res, err
}

func CreateStakedBalanceRole(ldb *db.DBService) error {
	return nil
}

func SetStakedBalanceRole(ldb *db.DBService, accountName string, value *StakedBalance) error {
	key := accountName
	jsonvalue, _ := json.Marshal(value)
	return ldb.SetObject(StakedBalanceObjectName, key, string(jsonvalue))
}

func GetStakedBalanceRoleByName(ldb *db.DBService, name string) (*StakedBalance, error) {
	key := name
	value, err := ldb.GetObject(StakedBalanceObjectName, key)
	res := &StakedBalance{}
	json.Unmarshal([]byte(value), res)
	//fmt.Println("Get", key, value)
	return res, err
}