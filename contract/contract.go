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
 * @Date:   2018-08-18
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"errors"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/contract/abi"
)

func getContractAbi(r role.RoleInterface, contract string) (*abi.ABI, error) {
	account, err := r.GetAccount(contract)
	if err != nil {
		return nil, errors.New("Get account fail")
	}

	Abi, err := abi.ParseAbi(account.ContractAbi)
	if err != nil {
		return nil, err
	}

	return Abi, nil
}

func ParseParam(role role.RoleInterface, contract string, method string, param []byte) (map[string]interface{}, error) {
	Abi, err := getContractAbi(role, contract)
	if  err != nil {
		return nil, errors.New("Abi not found")
	}

	decodedParam := abi.UnmarshalAbiEx(contract, Abi, method, param)
	if decodedParam == nil || len(decodedParam) <= 0 {
		return nil, errors.New("parse param fail")
	}
	return decodedParam, nil
}