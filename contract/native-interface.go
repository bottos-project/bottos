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
	berr "github.com/bottos-project/bottos/common/errors"
	"math/big"
	"regexp"
)

const MaxDelegateLocationLen int = 32
const MaxDelegateDescriptionLen int = 128

// BLOCKS_PER_ROUND define block num per round
const BLOCKS_PER_ROUND uint32 = 29

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

//StakeParam for stake, unstke and claim contract
type StakeParam struct {
	Amount *big.Int `json:"amount"`
}

//NativeContractInterface is native contract interface
type NativeContractInterface interface {
	IsNativeContract(contract string, method string) bool
	ExecuteNativeContract(*Context) berr.ErrCode
}

//NativeContractMethod is native contract method
type NativeContractMethod func(*Context) berr.ErrCode

//NativeContract is native contract handler
type NativeContract struct {
	Handler map[string]NativeContractMethod
	re *regexp.Regexp
}

//NewNativeContractHandler is native contract handler to handle different contracts
func NewNativeContractHandler() (NativeContractInterface, error) {
	nc := &NativeContract{
		Handler: make(map[string]NativeContractMethod),
	}
	nc.re = regexp.MustCompile(config.ACCOUNT_NAME_REGEXP)

	nc.Handler["newaccount"] = nc.newAccount
	nc.Handler["transfer"] = nc.transfer
	nc.Handler["setdelegate"] = nc.setDelegate
	nc.Handler["grantcredit"] = nc.grantCredit
	nc.Handler["cancelcredit"] = nc.cancelCredit
	nc.Handler["transferfrom"] = nc.transferFrom
	nc.Handler["deploycode"] = nc.deployCode
	nc.Handler["deployabi"] = nc.deployAbi
	nc.Handler["stake"] = nc.stake
	nc.Handler["unstake"] = nc.unstake
	nc.Handler["claim"] = nc.claim
	nc.Handler["regdelegate"] = nc.regDelegate
	nc.Handler["unregdelegate"] = nc.unregDelegate
	nc.Handler["votedelegate"] = nc.voteDelegate

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
func (nc *NativeContract) ExecuteNativeContract(ctx *Context) berr.ErrCode {
	contract := ctx.Trx.Contract
	method := ctx.Trx.Method
	if nc.IsNativeContract(contract, method) {
		if handler, ok := nc.Handler[method]; ok {
			contErr := handler(ctx)
			return contErr
		}
		return berr.ErrContractUnknownMethod
	}
	return berr.ErrContractUnknownContract

}
