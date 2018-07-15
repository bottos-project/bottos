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

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"
	log "github.com/cihub/seelog"
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
func NativeContractInitChain(ldb *db.DBService, roleIntf role.RoleInterface, ncIntf NativeContractInterface) ([]*types.Transaction, error) {
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

		Abi := GetAbi()

		mapstruct := make(map[string]interface{})
		abi.Setmapval(mapstruct, "name", name)
		abi.Setmapval(mapstruct, "pubkey", config.Genesis.InitDelegates[i].PublicKey)
		nparam, err2   := abi.MarshalAbiEx(mapstruct, Abi, config.BOTTOS_CONTRACT_NAME, "newaccount")
		if err2 != nil {
			log.Info("abi.MarshalAbiEx failed for new account:", name)
			continue
		}

		trx := newTransaction(config.BOTTOS_CONTRACT_NAME, "newaccount", nparam)
		trxs = append(trxs, trx)

		// 2, transfer trx

		mapstruct2 := make(map[string]interface{})
		abi.Setmapval(mapstruct2, "from", config.BOTTOS_CONTRACT_NAME)
		abi.Setmapval(mapstruct2, "to", name)
		abi.Setmapval(mapstruct2, "value", uint64(config.Genesis.InitDelegates[i].Balance))
		tparam, err3   := abi.MarshalAbiEx(mapstruct2, Abi, config.BOTTOS_CONTRACT_NAME, "transfer")
		if err3 != nil {
			log.Info("abi.MarshalAbiEx failed for transfer with account:", name)
			continue
		}

		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "transfer", tparam)
		trxs = append(trxs, trx)

		// 3, set delegate

		mapstruct3 := make(map[string]interface{})
		abi.Setmapval(mapstruct3, "name", name)
		abi.Setmapval(mapstruct3, "pubkey", config.Genesis.InitDelegates[i].PublicKey)
		sparam, err4   := abi.MarshalAbiEx(mapstruct3, Abi, config.BOTTOS_CONTRACT_NAME, "setdelegate")
		if err4 != nil {
			log.Info("abi.MarshalAbiEx failed for setdegelage with account:", name)
			continue
		}
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

	return trxs, nil
}

var a  *abi.ABI

func GetAbi() *abi.ABI {
	return a
}

func createNativeContractABI() *abi.ABI {

	a = &abi.ABI{}
	a.Actions = append(a.Actions, abi.ABIAction{ActionName: "newaccount", Type: "NewAccount"})
	a.Actions = append(a.Actions, abi.ABIAction{ActionName: "transfer", Type: "Transfer"})
	a.Actions = append(a.Actions, abi.ABIAction{ActionName: "setdelegate", Type: "SetDelegate"})
	a.Actions = append(a.Actions, abi.ABIAction{ActionName: "grantcredit", Type: "GrantCredit"})
	a.Actions = append(a.Actions, abi.ABIAction{ActionName: "cancelcredit", Type: "CancelCredit"})
	a.Actions = append(a.Actions, abi.ABIAction{ActionName: "transferfrom", Type: "TransferFrom"})
	a.Actions = append(a.Actions, abi.ABIAction{ActionName: "deploycode", Type: "DeployCode"})
	a.Actions = append(a.Actions, abi.ABIAction{ActionName: "deployabi", Type: "DeployABI"})

	s := abi.ABIStruct{Name: "NewAccount", Fields: abi.New()}
	s.Fields.Set("name", "string")
	s.Fields.Set("pubkey", "string")
	a.Structs = append(a.Structs, s)
	s = abi.ABIStruct{Name: "Transfer", Fields: abi.New()}
	s.Fields.Set("from", "string")
	s.Fields.Set("to", "string")
	s.Fields.Set("value", "uint64")
	a.Structs = append(a.Structs, s)
	s = abi.ABIStruct{Name: "SetDelegate", Fields: abi.New()}
	s.Fields.Set("name", "string")
	s.Fields.Set("pubkey", "string")
	a.Structs = append(a.Structs, s)
	s = abi.ABIStruct{Name: "GrantCredit", Fields: abi.New()}
	s.Fields.Set("name", "string")
	s.Fields.Set("spender", "string")
	s.Fields.Set("limit", "uint64")
	a.Structs = append(a.Structs, s)
	s = abi.ABIStruct{Name: "CancelCredit", Fields: abi.New()}
	s.Fields.Set("name", "string")
	s.Fields.Set("spender", "string")
	a.Structs = append(a.Structs, s)
	s = abi.ABIStruct{Name: "TransferFrom", Fields: abi.New()}
	s.Fields.Set("from", "string")
	s.Fields.Set("to", "string")
	s.Fields.Set("value", "uint64")
	a.Structs = append(a.Structs, s)
	s = abi.ABIStruct{Name: "DeployCode", Fields: abi.New()}
	s.Fields.Set("contract", "string")
	s.Fields.Set("vm_type", "uint8")
	s.Fields.Set("vm_version", "uint8")
	s.Fields.Set("contract_code", "bytes")
	a.Structs = append(a.Structs, s)
	s = abi.ABIStruct{Name: "DeployABI", Fields: abi.New()}
	s.Fields.Set("contract", "string")
	s.Fields.Set("contract_abi", "bytes")
	a.Structs = append(a.Structs, s)

	role.AbiAttr = a
	return a
}

//CreateNativeContractAccount is to create native contract account
func CreateNativeContractAccount(roleIntf role.RoleInterface) error {
	// account
	_, err := roleIntf.GetAccount(config.BOTTOS_CONTRACT_NAME)
	if err == nil {
		return nil
	}

	pubkey, _ := common.HexToBytes(config.Param.KeyPairs[0].PublicKey)
	a := createNativeContractABI()
	abijson, _ := abi.AbiToJson(a)
	bto := &role.Account{
		AccountName: config.BOTTOS_CONTRACT_NAME,
		CreateTime:  config.Genesis.GenesisTime,
		PublicKey:   pubkey,
		ContractAbi: []byte(abijson),
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
