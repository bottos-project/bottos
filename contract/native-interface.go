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
	"github.com/bottos-project/bottos/config"
)

//NewAccountParam struct for name and pubkey
type NewAccountParam struct {
	Name   string `json:"name"`
	Pubkey string `json:"pubkey"`
}

//TransferParam struct for transfer
type TransferParam struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

//SetDelegateParam struct for delegate
type SetDelegateParam struct {
	Name   string `json:"name"`
	Pubkey string `json:"pubkey"`
	// TODO CONFIG
}

//GrantCreditParam struct to grand credit
type GrantCreditParam struct {
	Name    string `json:"name"`
	Spender string `json:"spender"`
	Limit   uint64 `json:"limit"`
}

//CancelCreditParam struct to cancel credit
type CancelCreditParam struct {
	Name    string `json:"name"`
	Spender string `json:"spender"`
}

//TransferFromParam struct to transfer credit
type TransferFromParam struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

//DeployCodeParam struct to deploy code
type DeployCodeParam struct {
	Name         string `json:"contract"`
	VMType       byte   `json:"vm_type"`
	VMVersion    byte   `json:"vm_version"`
	ContractCode []byte `json:"contract_code"`
}

//DeployABIParam struct to deploy abi
type DeployABIParam struct {
	Name        string `json:"contract"`
	ContractAbi []byte `json:"contract_abi"`
}

//NativeContractInterface is native contract interface
type NativeContractInterface interface {
	IsNativeContract(contract string, method string) bool
	ExecuteNativeContract(*Context) ContractError
}

//NativeContractMethod is native contract method
type NativeContractMethod func(*Context) ContractError

//NativeContract is native contract handler
type NativeContract struct {
	Handler map[string]NativeContractMethod
}

//NewNativeContractHandler is native contract handler to handle different contracts
func NewNativeContractHandler() (NativeContractInterface, error) {
	nc := &NativeContract{
		Handler: make(map[string]NativeContractMethod),
	}

	nc.Handler["newaccount"] = newAccount
	nc.Handler["transfer"] = transfer
	nc.Handler["setdelegate"] = setDelegate
	nc.Handler["grantcredit"] = grantCredit
	nc.Handler["cancelcredit"] = cancelCredit
	nc.Handler["transferfrom"] = transferFrom
	nc.Handler["deploycode"] = deployCode
	nc.Handler["deployabi"] = deployAbi

	return nc, nil
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
