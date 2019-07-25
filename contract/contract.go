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
 * file description:  contract
 * @Author: Gong Zibin
 * @Date:   2017-01-15
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"math/big"

	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/common"
)



//NewNativeContract is to create a new native contract
func NewNativeContract(roleIntf role.RoleInterface) (NativeContractInterface, error) {
	intf, err := NewNativeContractHandler()
	if err != nil {
		return nil, err
	}
	roleIntf.SetScheduleDelegate(&role.ScheduleDelegate{CurrentTermTime: big.NewInt(2)})

	return intf, nil
}

//NativeContractInitChain is to init
func NativeContractInitChain(ldb *db.DBService, roleIntf role.RoleInterface, ncIntf NativeContractInterface) error {
	err := CreateNativeContractAccount(roleIntf)
	if err != nil {
		return err
	}

	// init CoreState delegates
	coreState, _ := roleIntf.GetCoreState()
	coreState.CurrentDelegates = []string{config.BOTTOS_CONTRACT_NAME}
	roleIntf.SetCoreState(coreState)

	roleIntf.InitRewardPoolRole()

	return nil
}

//CreateNativeContractAccount is to create native contract account
func CreateNativeContractAccount(roleIntf role.RoleInterface) error {
	// account
	_, err := roleIntf.GetAccount(config.BOTTOS_CONTRACT_NAME)
	if err == nil {
		return nil
	}

	a := abi.CreateNativeContractABI()
	abijson, _ := abi.AbiToJson(a)
	bto := &role.Account{
		AccountName: config.BOTTOS_CONTRACT_NAME,
		CreateTime:  config.Genesis.GenesisTime,
		PublicKey:   config.Genesis.GenesisKey,
		ContractAbi: []byte(abijson),
	}
	roleIntf.SetAccount(bto.AccountName, bto)

	// balance
	var initSupply *big.Int = big.NewInt(0)
	initSupply, err = safemath.U256Mul(initSupply, new(big.Int).SetUint64(config.BOTTOS_INIT_SUPPLY), new(big.Int).SetUint64(config.BOTTOS_SUPPLY_MUL))
	if err != nil {
		return err
	}

	balance := &role.Balance{
		AccountName: bto.AccountName,
		Balance:     initSupply,
	}
	roleIntf.SetBalance(bto.AccountName, balance)

	// staked_balance
	stakedBalance := &role.StakedBalance{
		AccountName:       bto.AccountName,
		StakedBalance:     big.NewInt(0),
		UnstakingBalance:  big.NewInt(0),
		LastUnstakingTime: 0,
	}
	roleIntf.SetStakedBalance(bto.AccountName, stakedBalance)

	return nil
}
