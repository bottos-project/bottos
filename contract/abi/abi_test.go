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
 * @Date:   2017-01-15
 * @Last Modified by:
 * @Last Modified time:
 */

package abi

import (
	"fmt"
	"testing"
)

func TestTransfer(t *testing.T) {
	j := `
	{
		"types": [],
		"structs": [
			{
				"name": "NewAccount",
				"base": "",
				"fields": {
					"name": "string",
					"pubkey": "string"
				}
			},
			{
				"name": "Transfer",
				"base": "",
				"fields": {
					"from": "string",
					"to": "string",
					"value": "uint64"
				}
			},
			{
				"name": "SetDelegate",
				"base": "",
				"fields": {
					"name": "string",
					"pubkey": "string"
				}
			},
			{
				"name": "GrantCredit",
				"base": "",
				"fields": {
					"name": "string",
					"spender": "string",
					"limit": "uint64"
				}
			},
			{
				"name": "CancelCredit",
				"base": "",
				"fields": {
					"name": "string",
					"spender": "string"
				}
			},
			{
				"name": "TransferFrom",
				"base": "",
				"fields": {
					"from": "string",
					"to": "string",
					"value": "uint64"
				}
			},
			{
				"name": "DeployCode",
				"base": "",
				"fields": {
					"contract": "string",
					"vm_type": "uint8",
					"vm_version": "uint8",
					"contract_code": "bytes"
				}
			},
			{
				"name": "DeployABI",
				"base": "",
				"fields": {
					"contract": "string",
					"contract_abi": "bytes"
				}
			}
		],
		"actions": [
			{
				"action_name": "newaccount",
				"type": "NewAccount"
			},
			{
				"action_name": "transfer",
				"type": "Transfer"
			},
			{
				"action_name": "setdelegate",
				"type": "SetDelegate"
			},
			{
				"action_name": "grantcredit",
				"type": "GrantCredit"
			},
			{
				"action_name": "cancelcredit",
				"type": "CancelCredit"
			},
			{
				"action_name": "transferfrom",
				"type": "TransferFrom"
			},
			{
				"action_name": "deploycode",
				"type": "DeployCode"
			},
			{
				"action_name": "deployabi",
				"type": "DeployABI"
			}
		],
		"tables": []
	}
	`
	a, err := ParseAbi([]byte(j))
	fmt.Println(a)
	if err != nil {
		fmt.Printf("Abi Parse Error Str: %v", j)
		return
	}

	for i := range a.Structs {
		fmt.Println(a.Structs[i].Fields.GetStringPair())
	}

}
