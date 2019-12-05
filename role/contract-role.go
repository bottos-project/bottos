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
 * file description:  contract role
 * @Author: Gong Zibin
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/db"
	log "github.com/cihub/seelog"
)

const (
	//ContractObjectName is the table name of user contract
	ContractObjectName string = "contract"
)

// Contract is definition of user contract
type Contract struct {
	ContractName      string      `json:"contract_name"`
	VMType            byte        `json:"vm_type"`
	VMVersion         byte        `json:"vm_version"`
	CodeVersion       common.Hash `json:"code_version"`
	ContractCode      []byte      `json:"contract_code"`
	ContractAbi       []byte      `json:"abi"`
	DeployAccountName string      `json:"deploy_account_name"`
}

// CreateContractRole is create contract role
func CreateContractRole(ldb *db.DBService) error {
	ldb.AddObject(ContractObjectName)
	return nil
}

func contractNameToKey(name string) string {
	return name
}

// SetContractRole is common func to set role for contract
func SetContractRole(ldb *db.DBService, contractName string, value *Contract) error {
	key := contractNameToKey(contractName)
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		log.Error("ROLE Marshal failed contractName ", contractName)
		return err
	}
	return ldb.SetObject(ContractObjectName, key, string(jsonvalue))
}

// GetContractRole is common func to get role for contract
func GetContractRole(ldb *db.DBService, contractName string) (*Contract, error) {
	key := contractNameToKey(contractName)
	value, err := ldb.GetObject(ContractObjectName, key)
	if err != nil {
		return nil, err
	}

	res := &Contract{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		log.Error("ROLE Unmarshal failed contractName ", contractName)
		return nil, err
	}

	return res, nil
}
