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
 * file description:  native contract
 * @Author: Gong Zibin
 * @Date:   2017-01-15
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"math/big"

	"fmt"
	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/abi"
		"github.com/bottos-project/bottos/role"
)


//NativeContract is native contract handler
type NativeContract struct {
	Handler map[string]NativeContractMethod
	Abi     *abi.ABI1
}


//NewNativeContract is to create a new native contract
func NewNativeContract(roleIntf role.RoleInterface) (NativeContractInterface, error) {
	nc := &NativeContract{
		Handler: make(map[string]NativeContractMethod),
	}

	var err error
	nc.Abi, err = NewNativeContractABI()
	if err != nil {
		return nil, fmt.Errorf("Native abi error: %v", err)
	}

	nc.Handler["newaccount"] = newAccount
	nc.Handler["transfer"] = transfer
	nc.Handler["setdelegate"] = setDelegate
	nc.Handler["grantcredit"] = grantCredit
	nc.Handler["cancelcredit"] = cancelCredit
	nc.Handler["transferfrom"] = transferFrom
	nc.Handler["deploycode"] = deployCode
	nc.Handler["deployabi"] = deployAbi

	roleIntf.SetScheduleDelegate(&role.ScheduleDelegate{CurrentTermTime: big.NewInt(2)})

	return nc, nil
}

//NativeContractInit is to init
func (nc *NativeContract) NativeContractInit(role role.RoleInterface) ([]*types.Transaction, error) {
	err := nc.CreateNativeContractAccount(role)
	if err != nil {
		return nil, err
	}

	a := nc.GetABI()
	var trxs []*types.Transaction
	for i := 0; i < len(config.Genesis.InitDelegates); i++ {
		name := config.Genesis.InitDelegates[i].Name

		// 1, new account trx
		param, _ := a.Pack("newaccount", name, config.Genesis.InitDelegates[i].PublicKey)
		trx := newNativeTransaction("newaccount", param)
		trxs = append(trxs, trx)
		fmt.Println("newaccount: ", trx)

		// 2, transfer trx
		param, err := a.Pack("transfer", config.BOTTOS_CONTRACT_NAME, name, uint64(config.Genesis.InitDelegates[i].Balance))
		trx = newNativeTransaction("transfer", param)
		trxs = append(trxs, trx)
		fmt.Println("transfer: ", trx, err)

		// 3, set delegate
		param, _ = a.Pack("setdelegate", name, config.Genesis.InitDelegates[i].PublicKey)
		trx = newNativeTransaction("setdelegate", param)
		trxs = append(trxs, trx)
		fmt.Println("setdelegate: ", trx)
	}

	// init CoreState delegates
	coreState, _ := role.GetCoreState()
	for i := 0; i < int(config.BLOCKS_PER_ROUND); i++ {
		name := config.Genesis.InitDelegates[i].Name
		coreState.CurrentDelegates = append(coreState.CurrentDelegates, name)
	}
	role.SetCoreState(coreState)

	return trxs, nil
}

func newNativeTransaction(method string, param []byte) *types.Transaction {
	trx := &types.Transaction{
		Version:     1,
		CursorNum:   0,
		CursorLabel: 0,
		Sender:      config.BOTTOS_CONTRACT_NAME,
		Contract:    config.BOTTOS_CONTRACT_NAME,
		Method:      method,
	}
	trx.Param = make([]byte, len(param))
	copy(trx.Param, param)

	return trx
}

//CreateNativeContractAccount is to create native contract account
func (nc *NativeContract) CreateNativeContractAccount(roleIntf role.RoleInterface) error {
	// account
	_, err := roleIntf.GetAccount(config.BOTTOS_CONTRACT_NAME)
	if err == nil {
		return nil
	}

	a := nc.GetABI()
	bto := &role.Account{
		AccountName: config.BOTTOS_CONTRACT_NAME,
		CreateTime:  config.Genesis.GenesisTime,
		PublicKey:   config.Genesis.GenesisKey,
		ContractAbi: []byte(a.ToJson(false)),
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

func (nc *NativeContract) GetABI() *abi.ABI1 {
	return nc.Abi
}

//IsNativeContract is to check if the contract is native
func (nc *NativeContract) IsNativeContract(contract string, method string) bool {
	if contract == config.BOTTOS_CONTRACT_NAME {
		if _, ok := nc.Handler[method]; ok {
			return true
		}
	}
	return false
}

//ExecuteNativeContract is to call native contract
func (nc *NativeContract) ExecuteNativeContract(ctx *Context) ContractError {
	contract := ctx.Trx.Contract
	method := ctx.Trx.Method
	if nc.IsNativeContract(contract, method) {
		if handler, ok := nc.Handler[method]; ok {
			contErr := handler(ctx)
			return contErr
		}
		return ERROR_CONT_UNKNOWN_METHOD
	}
	return ERROR_CONT_UNKNOWN_CONTARCT
}
