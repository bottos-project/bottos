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

	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/db"
	"math/big"
)

// BalanceObjectName is definition of object name of balance
const BalanceObjectName string = "balance"
// StakedBalanceObjectName is definition of object name of stake balance
const StakedBalanceObjectName string = "staked_balance"

// Balance is definition of balance
type Balance struct {
	AccountName string `json:"account_name"`
	Balance     *big.Int `json:"balance"`
}

// StakedBalance is definition of stake balance
type StakedBalance struct {
	AccountName   string `json:"account_name"`
	StakedBalance *big.Int `json:"staked_balance"`
	UnstakingBalance *big.Int `json:"unstaking_balance"`
	LastUnstakingTime uint64 `json:"unstaking_time"`
}

// UnstakingBalance is definition of unstaking balance
type UnstakingBalance struct {
	AccountName   string `json:"account_name"`
	UnstakingBalance *big.Int `json:"unstaking_balance"`
	UnstakingTime uint64 `json:"unstaking_time"`
}

// CreateBalanceRole is to create balance role
func CreateBalanceRole(ldb *db.DBService) error {
	return nil
}

// SetBalanceRole is to set balance
func SetBalanceRole(ldb *db.DBService, accountName string, value *Balance) error {
	key := accountName
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ldb.SetObject(BalanceObjectName, key, string(jsonvalue))
}

// GetBalanceRole is to get balance
func GetBalanceRole(ldb *db.DBService, accountName string) (*Balance, error) {
	key := accountName
	value, err := ldb.GetObject(BalanceObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &Balance{Balance:big.NewInt(0)}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func safeAdd(result *big.Int, a *big.Int, b *big.Int) (*big.Int, error) {
	//var c uint64
	result, err := safemath.U256Add(result, a, b)
	if err != nil {
		return result, err
	}
	return result, nil
}

func safeSub(result *big.Int, a *big.Int, b *big.Int) (*big.Int, error) {
	result, err := safemath.U256Sub(result, a, b)
	if err != nil {
		return result, err
	}
	return result, nil
}

// SafeAdd is safe function to add balance
func (balance *Balance) SafeAdd(amount *big.Int) error {
	//var a, c uint64
	result := big.NewInt(0)
	result, err := safeAdd(result, balance.Balance, amount)
	if err != nil {
		return err
	}

	balance.Balance.Set(result)
	return nil
}

// SafeSub is safe function to sub balance
//var a, c big.Int
func (balance *Balance) SafeSub(amount *big.Int) error {
	result := big.NewInt(0)
	//a = 
	result, err := safeSub(result, balance.Balance, amount)
	if err != nil {
		return err
	}

	balance.Balance.Set(result)
	return nil
}

// CreateStakedBalanceRole is to create stake balance role
func CreateStakedBalanceRole(ldb *db.DBService) error {
	return nil
}

// SetStakedBalanceRole is to set stake balance role
func SetStakedBalanceRole(ldb *db.DBService, accountName string, value *StakedBalance) error {
	key := accountName
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(StakedBalanceObjectName, key, string(jsonvalue))
}

// GetStakedBalanceRoleByName is to get stake balance
func GetStakedBalanceRoleByName(ldb *db.DBService, name string) (*StakedBalance, error) {
	key := name
	value, err := ldb.GetObject(StakedBalanceObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &StakedBalance{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SafeAdd is safe function to add  stake balance
func (balance *StakedBalance) SafeAdd(amount *big.Int) error {
	//var a, c uint64
	result := big.NewInt(0)
	//a = balance.StakedBalance
	result, err := safeAdd(result, balance.StakedBalance, amount)
	if err != nil {
		return err
	}

	balance.StakedBalance.Set(result)
	return nil
}

// SafeSub is safe function to sub stake balance
func (balance *StakedBalance) SafeSub(amount *big.Int) error {
	//var a, c uint64
	result := big.NewInt(0)
	result, err := safeSub(result, balance.StakedBalance, amount)
	if err != nil {
		return err
	}

	balance.StakedBalance.Set(result)
	return nil
}

// UnstakingAmount
func (balance *StakedBalance) UnstakingAmount(amount *big.Int) error {
	result := big.NewInt(0)
	result, err := safeSub(result, balance.StakedBalance, amount)
	if err != nil {
		return err
	}
	result1 := big.NewInt(0)
	result1, err = safeAdd(result1, balance.UnstakingBalance, amount)
	if err != nil {
		return err
	}
	balance.StakedBalance.Set(result)
	balance.UnstakingBalance.Set(result1)
	return nil
}

// UnstakingAmount
func (balance *StakedBalance) Claim(amount *big.Int) error {
	result := big.NewInt(0)
	result, err := safeSub(result, balance.UnstakingBalance, amount)
	if err != nil {
		return err
	}
	balance.UnstakingBalance.Set(result)
	return nil
}
