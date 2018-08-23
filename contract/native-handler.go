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
	"fmt"
	"regexp"
	"math/big"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/role"
)

func newAccount(ctx *Context) ContractError {
	Abi := abi.GetAbi()
	newaccount := abi.UnmarshalAbiEx("bottos", Abi, "newaccount", ctx.Trx.Param)
	if newaccount == nil || len(newaccount) <= 0 {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	
	NewaccountName   := newaccount["name"].(string)
	NewaccountPubKey := newaccount["pubkey"].(string)
	
	//check account
	cerr := checkAccountName(NewaccountName)
	if cerr != ERROR_NONE {
		return cerr
	}

	if isAccountNameExist(ctx.RoleIntf, NewaccountName) {
		return ERROR_CONT_ACCOUNT_ALREADY_EXIST
	}

	chainState, _ := ctx.RoleIntf.GetChainState()
	// 1, create account
	pubkey, _ := common.HexToBytes(NewaccountPubKey)
	account := &role.Account{
		AccountName: NewaccountName,
		PublicKey:   pubkey,
		CreateTime:  chainState.LastBlockTime,
	}
	ctx.RoleIntf.SetAccount(account.AccountName, account)

	// 2, create balance
	balance := &role.Balance{
		AccountName: NewaccountName,
		Balance:     big.NewInt(0),
	}
	ctx.RoleIntf.SetBalance(NewaccountName, balance)

	// 3, create staked_balance
	stakedBalance := &role.StakedBalance{
		AccountName:   NewaccountName,
		StakedBalance: big.NewInt(0),
	}
	ctx.RoleIntf.SetStakedBalance(NewaccountName, stakedBalance)

	return ERROR_NONE
}

func transfer(ctx *Context) ContractError {
	Abi := abi.GetAbi()
	transfer := abi.UnmarshalAbiEx("bottos", Abi, "transfer", ctx.Trx.Param)
	if transfer == nil || len(transfer) <= 0 {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	
	FromWhom := transfer["from"].(string)
	ToWhom   := transfer["to"].(string)
	TransValue := transfer["value"].(*big.Int)
	
	// check account
	cerr := checkAccount(ctx.RoleIntf, FromWhom)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ToWhom)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(FromWhom)
	if -1 == from.Balance.Cmp(TransValue) {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}
	to, _ := ctx.RoleIntf.GetBalance(ToWhom)

	err := from.SafeSub(TransValue)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = to.SafeAdd(TransValue)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}

	err = ctx.RoleIntf.SetBalance(from.AccountName, from)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}
	err = ctx.RoleIntf.SetBalance(to.AccountName, to)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}
	
	return ERROR_NONE
}

func setDelegate(ctx *Context) ContractError {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "setdelegate", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	
	ParamName   := param["name"].(string)
	ParamPubKey := param["pubkey"].(string)	
	
	// check account
	cerr := checkAccount(ctx.RoleIntf, ParamName)
	if cerr != ERROR_NONE {
		return cerr
	}

	_, err := ctx.RoleIntf.GetDelegateByAccountName(ParamName)
	if err != nil {
		// new delegate
		newdelegate := &role.Delegate{
			AccountName: ParamName,
			ReportKey:   ParamPubKey,
		}
		ctx.RoleIntf.SetDelegate(newdelegate.AccountName, newdelegate)

		//create schedule delegate vote role
		scheduleDelegate, err := ctx.RoleIntf.GetScheduleDelegate()
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}

		newDelegateVotes := new(role.DelegateVotes).StartNewTerm(scheduleDelegate.CurrentTermTime)
		newDelegateVotes.OwnerAccount = newdelegate.AccountName
		err = ctx.RoleIntf.SetDelegateVotes(newdelegate.AccountName, newDelegateVotes)
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}
	} else {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func grantCredit(ctx *Context) ContractError {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "grantcredit", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	
	ParamName    := param["name"].(string)
	ParamSpender := param["spender"].(string)
	ParamLimit   := param["limit"].(*big.Int)

	// check account
	cerr := checkAccount(ctx.RoleIntf, ParamName)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ParamSpender)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ctx.Trx.Sender)
	if cerr != ERROR_NONE {
		return cerr
	}

	// sender must be from
	if ctx.Trx.Sender != ParamName {
		return ERROR_CONT_ACCOUNT_MISMATCH
	}

	// check limit
	balance, err := ctx.RoleIntf.GetBalance(ParamName)
	if -1 == balance.Balance.Cmp(ParamLimit) {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}

	credit := &role.TransferCredit{
		Name:    ParamName,
		Spender: ParamSpender,
		Limit:   ParamLimit,  
	}
	err = ctx.RoleIntf.SetTransferCredit(credit.Name, credit)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func cancelCredit(ctx *Context) ContractError {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "cancelcredit", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	ParamName    := param["name"].(string)
	ParamSpender := param["spender"].(string)
	
	// check account
	cerr := checkAccount(ctx.RoleIntf, ParamName)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ParamSpender)
	if cerr != ERROR_NONE {
		return cerr
	}

	_, err := ctx.RoleIntf.GetTransferCredit(ParamName, ParamSpender)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	err = ctx.RoleIntf.DeleteTransferCredit(ParamName, ParamSpender)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func transferFrom(ctx *Context) ContractError {
	Abi := abi.GetAbi()
	transfer := abi.UnmarshalAbiEx("bottos", Abi, "transferfrom", ctx.Trx.Param)
	if transfer == nil || len(transfer) <= 0 {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	TransFrom := transfer["from"].(string)
	TransTo := transfer["to"].(string)
	TransValue := transfer["value"].(*big.Int)
	
	// check account
	cerr := checkAccount(ctx.RoleIntf, TransFrom)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, TransTo)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ctx.Trx.Sender)
	if cerr != ERROR_NONE {
		return cerr
	}

	// Note: sender is the spender
	// check limit
	credit, err := ctx.RoleIntf.GetTransferCredit(TransFrom, ctx.Trx.Sender)
	if err != nil {
		return ERROR_CONT_INSUFFICIENT_CREDITS
	}
	if 1 == TransValue.Cmp(credit.Limit) {
		return ERROR_CONT_INSUFFICIENT_CREDITS
	}

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(TransFrom)
	if -1 == from.Balance.Cmp(TransValue) {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}
	to, _ := ctx.RoleIntf.GetBalance(TransTo)

	err = from.SafeSub(TransValue)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = credit.SafeSub(TransValue)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = to.SafeAdd(TransValue)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}

	err = ctx.RoleIntf.SetBalance(from.AccountName, from)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}
	err = ctx.RoleIntf.SetBalance(to.AccountName, to)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	if 1 == credit.Limit.Cmp(big.NewInt(0)) {
		err = ctx.RoleIntf.SetTransferCredit(credit.Name, credit)
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}
	} else {
		err = ctx.RoleIntf.DeleteTransferCredit(credit.Name, ctx.Trx.Sender)
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}
	}

	return ERROR_NONE
}

func checkCode(code []byte) error {
	// TODO
	return nil
}

func deployCode(ctx *Context) ContractError {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "deploycode", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	
	ParamContract := param["contract"].(string)
	ParamContractCode, _ := common.HexToBytes(param["contract_code"].(string))

	// check account
	cerr := checkAccount(ctx.RoleIntf, ParamContract)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check code
	err := checkCode(ParamContractCode)
	if err != nil {
		return ERROR_CONT_CODE_INVALID
	}

	codeHash := common.Sha256(ParamContractCode)

	account, _ := ctx.RoleIntf.GetAccount(ParamContract)
	account.CodeVersion = codeHash
	account.ContractCode = make([]byte, len(ParamContractCode))
	copy(account.ContractCode, ParamContractCode)
	err = ctx.RoleIntf.SetAccount(account.AccountName, account)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func checkAbi(abiRaw []byte) error {
	_, err := abi.ParseAbi(abiRaw)
	if err != nil {
		return fmt.Errorf("ABI Parse error: %v", err)
	}
	return nil
}

func deployAbi(ctx *Context) ContractError {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "deployabi", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	
	ParamContract := param["contract"].(string)
	ParamContractAbi, _ := common.HexToBytes(param["contract_abi"].(string))
	
	// check account
	cerr := checkAccount(ctx.RoleIntf, ParamContract)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check abi
	err := checkAbi(ParamContractAbi)
	if err != nil {
		return ERROR_CONT_ABI_PARSE_FAIL
	}

	account, _ := ctx.RoleIntf.GetAccount(ParamContract)
	account.ContractAbi = make([]byte, len(ParamContractAbi))
	copy(account.ContractAbi, ParamContractAbi)
	err = ctx.RoleIntf.SetAccount(account.AccountName, account)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func checkAccountName(name string) ContractError {
	if len(name) == 0 {
		return ERROR_CONT_ACCOUNT_NAME_NULL
	}

	if len(name) > config.MAX_ACCOUNT_NAME_LENGTH {
		return ERROR_CONT_ACCOUNT_NAME_TOO_LONG
	}

	if !checkAccountNameContent(name) {
		return ERROR_CONT_ACCOUNT_NAME_ILLEGAL
	}

	return ERROR_NONE
}

func checkAccountNameContent(name string) bool {
	match, err := regexp.MatchString(config.ACCOUNT_NAME_REGEXP, name)
	if err != nil {
		return false
	}
	if !match {
		return false
	}

	return true
}

func isAccountNameExist(RoleIntf role.RoleInterface, name string) bool {
	account, err := RoleIntf.GetAccount(name)
	if err == nil {
		if account != nil && account.AccountName == name {
			return true
		}
	}
	return false
}

func checkAccount(RoleIntf role.RoleInterface, name string) ContractError {
	cerr := checkAccountName(name)
	if cerr != ERROR_NONE {
		return cerr
	}

	if !isAccountNameExist(RoleIntf, name) {
		return ERROR_CONT_ACCOUNT_NOT_EXIST
	}

	return ERROR_NONE
}
