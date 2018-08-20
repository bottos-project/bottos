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
 * file description:  native contract interface
 * @Author: Gong Zibin
 * @Date:   2017-01-15
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/role"
)

//NewAccountParam struct for name and pubkey
type NewAccountParam struct {
	Name   string
	Pubkey string
}

//TransferParam struct for transfer
type TransferParam struct {
	From  string
	To    string
	Value uint64
}

//SetDelegateParam struct for delegate
type SetDelegateParam struct {
	Name   string
	Pubkey string
	// TODO CONFIG
}

//GrantCreditParam struct to grand credit
type GrantCreditParam struct {
	Name    string
	Spender string
	Limit   uint64
}

//CancelCreditParam struct to cancel credit
type CancelCreditParam struct {
	Name    string
	Spender string
}

//TransferFromParam struct to transfer credit
type TransferFromParam struct {
	From  string
	To    string
	Value uint64
}

//DeployCodeParam struct to deploy code
type DeployCodeParam struct {
	Name         string
	VMType       byte
	VMVersion    byte
	ContractCode []byte
}

//DeployABIParam struct to deploy abi
type DeployABIParam struct {
	Name        string
	ContractAbi []byte
}


//NativeContractInterface is native contract interface
type NativeContractInterface interface {
	NativeContractInit(role role.RoleInterface) ([]*types.Transaction, error)
	IsNativeContract(contract string, method string) bool
	ExecuteNativeContract(*Context) ContractError
	GetABI() *abi.ABI1
}

//NativeContractMethod is native contract method
type NativeContractMethod func(*Context) ContractError
