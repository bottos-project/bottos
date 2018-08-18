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
 * file description:  native contract abi definition
 * @Author: Gong Zibin
 * @Date:   2018-08-18
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"github.com/bottos-project/bottos/contract/abi"
	"fmt"
)

func NewNativeContractABI() (*abi.ABI1, error) {
	abiDef := abi.ABIDef{}
	abiDef.Structs = append(abiDef.Structs,
		abi.NewABIStruct("NewAccount", "",
			abi.ABIDefFeild{"name", "string"},
			abi.ABIDefFeild{"pubkey", "string"},
		))
	abiDef.Structs = append(abiDef.Structs,
		abi.NewABIStruct("Transfer", "",
			abi.ABIDefFeild{"from", "string"},
			abi.ABIDefFeild{"to", "string"},
			abi.ABIDefFeild{"value", "uint64"},
		))
	abiDef.Structs = append(abiDef.Structs,
		abi.NewABIStruct("SetDelegate", "",
			abi.ABIDefFeild{"name", "string"},
			abi.ABIDefFeild{"pubkey", "string"},
		))
	abiDef.Structs = append(abiDef.Structs,
		abi.NewABIStruct("GrantCredit", "",
			abi.ABIDefFeild{"name", "string"},
			abi.ABIDefFeild{"spender", "string"},
			abi.ABIDefFeild{"limit", "uint64"},
		))
	abiDef.Structs = append(abiDef.Structs,
		abi.NewABIStruct("CancelCredit", "",
			abi.ABIDefFeild{"name", "string"},
			abi.ABIDefFeild{"spender", "string"},
		))
	abiDef.Structs = append(abiDef.Structs,
		abi.NewABIStruct("TransferFrom", "",
			abi.ABIDefFeild{"from", "string"},
			abi.ABIDefFeild{"to", "string"},
			abi.ABIDefFeild{"value", "uint64"},
		))
	abiDef.Structs = append(abiDef.Structs,
		abi.NewABIStruct("DeployCode", "",
			abi.ABIDefFeild{"contract", "string"},
			abi.ABIDefFeild{"vm_type", "uint8"},
			abi.ABIDefFeild{"vm_version", "uint8"},
			abi.ABIDefFeild{"contract_code", "bytes"},
		))
	abiDef.Structs = append(abiDef.Structs,
		abi.NewABIStruct("DeployABI", "",
			abi.ABIDefFeild{"contract", "string"},
			abi.ABIDefFeild{"contract_abi", "bytes"},
		))

	abiDef.Methods = append(abiDef.Methods, abi.NewABIMethod("newaccount", "NewAccount"))
	abiDef.Methods = append(abiDef.Methods, abi.NewABIMethod("transfer", "Transfer"))
	abiDef.Methods = append(abiDef.Methods, abi.NewABIMethod("setdelegate", "SetDelegate"))
	abiDef.Methods = append(abiDef.Methods, abi.NewABIMethod("grantcredit", "GrantCredit"))
	abiDef.Methods = append(abiDef.Methods, abi.NewABIMethod("cancelcredit", "CancelCredit"))
	abiDef.Methods = append(abiDef.Methods, abi.NewABIMethod("transferfrom", "TransferFrom"))
	abiDef.Methods = append(abiDef.Methods, abi.NewABIMethod("deploycode", "DeployCode"))
	abiDef.Methods = append(abiDef.Methods, abi.NewABIMethod("deployabi", "DeployABI"))

	fmt.Println(abiDef)
	return abi.NewABIFromDef(&abiDef)
}
