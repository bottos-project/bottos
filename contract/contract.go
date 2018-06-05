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
 * file description:  context definition
 * @Author: Gong Zibin
 * @Date:   2017-01-20
 * @Last Modified by:
 * @Last Modified time:
 */
package contract

import (
	"math/big"

	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/msgpack"
	"github.com/bottos-project/bottos/role"
)

//NewNativeContract is to create a new native contract
func NewNativeContract(roleIntf role.RoleInterface) (NativeContractInterface, error) {
	intf, err := NewNativeContractHandler()
	if err != nil {
		return nil, err
	}
	roleIntf.SetScheduleDelegate(&role.ScheduleDelegate{big.NewInt(2)})

	return intf, nil
}

func newTransaction(contract string, method string, param []byte) *types.Transaction {
	trx := &types.Transaction{
		Sender:   contract,
		Contract: contract,
		Method:   method,
		Param:    param,
	}

	return trx
}

//NativeContractInitChain is to init
func NativeContractInitChain(roleIntf role.RoleInterface, ncIntf NativeContractInterface) ([]*types.Transaction, error) {
	err := CreateNativeContractAccount(roleIntf)
	if err != nil {
		return nil, err
	}

	var trxs []*types.Transaction

	// construct trxs
	var i int
	for i = 0; i < len(config.Genesis.InitDelegates); i++ {
		name := config.Genesis.InitDelegates[i].Name

		// 1, new account trx
		nps := &NewAccountParam{
			Name:   name,
			Pubkey: config.Genesis.InitDelegates[i].PublicKey,
		}
		nparam, _ := msgpack.Marshal(nps)
		trx := newTransaction(config.BOTTOS_CONTRACT_NAME, "newaccount", nparam)
		trxs = append(trxs, trx)

		// 2, transfer trx
		tps := &TransferParam{
			From:  config.BOTTOS_CONTRACT_NAME,
			To:    name,
			Value: uint64(config.Genesis.InitDelegates[i].Balance),
		}
		tparam, _ := msgpack.Marshal(tps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "transfer", tparam)
		trxs = append(trxs, trx)

		// 3, set delegate
		sps := &SetDelegateParam{
			Name:   name,
			Pubkey: config.Genesis.InitDelegates[i].PublicKey,
		}
		sparam, _ := msgpack.Marshal(sps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "setdelegate", sparam)
		trxs = append(trxs, trx)
	}

	// init CoreState delegates
	coreState, _ := roleIntf.GetCoreState()
	for i = 0; i < int(config.BLOCKS_PER_ROUND); i++ {
		name := config.Genesis.InitDelegates[i].Name

		coreState.CurrentDelegates = append(coreState.CurrentDelegates, name)
	}
	roleIntf.SetCoreState(coreState)

	//fmt.Println("NativeContractInitChain: ", coreState)

	return trxs, nil
}

//CreateNativeContractAccount is to create native contract account
func CreateNativeContractAccount(roleIntf role.RoleInterface) error {
	// account
	_, err := roleIntf.GetAccount(config.BOTTOS_CONTRACT_NAME)
	if err == nil {
		return nil
	}

	bto := &role.Account{
		AccountName: config.BOTTOS_CONTRACT_NAME,
		CreateTime:  config.Genesis.GenesisTime,
	}
	roleIntf.SetAccount(bto.AccountName, bto)

	// balance
	var initSupply uint64
	initSupply, err = safemath.Uint64Mul(config.BOTTOS_INIT_SUPPLY, config.BOTTOS_SUPPLY_MUL)
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
		AccountName: bto.AccountName,
	}
	roleIntf.SetStakedBalance(bto.AccountName, stakedBalance)

	return nil
}
