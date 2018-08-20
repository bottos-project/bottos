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
 * file description:  native contract handler
 * @Author: Gong Zibin
 * @Date:   2017-01-15
 * @Last Modified by:
 * @Last Modified time:
 */

package contract

import (
	"fmt"
	"regexp"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/role"
)

func newAccount(ctx *Context) ContractError {
	param := &NewAccountParam{}
	if err := abi.UmarshalStruct(ctx.Trx.Param, param); err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	pubkey, err := common.HexToBytes(param.Pubkey)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	//check account
	cerr := checkAccountName(param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	if isAccountNameExist(ctx.RoleIntf, param.Name) {
		return ERROR_CONT_ACCOUNT_ALREADY_EXIST
	}

	chainState, _ := ctx.RoleIntf.GetChainState()
	// 1, create account
	account := &role.Account{
		AccountName: param.Name,
		PublicKey:   pubkey,
		CreateTime:  chainState.LastBlockTime,
	}
	ctx.RoleIntf.SetAccount(account.AccountName, account)

	// 2, create balance
	balance := &role.Balance{
		AccountName: param.Name,
		Balance:     0,
	}
	ctx.RoleIntf.SetBalance(param.Name, balance)

	// 3, create staked_balance
	stakedBalance := &role.StakedBalance{
		AccountName:   param.Name,
		StakedBalance: 0,
	}
	ctx.RoleIntf.SetStakedBalance(param.Name, stakedBalance)

	return ERROR_NONE
}

func transfer(ctx *Context) ContractError {
	var err error
	param := &TransferParam{}
	if err := abi.UmarshalStruct(ctx.Trx.Param, param); err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.From)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, param.To)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(param.From)
	if from.Balance < param.Value {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}
	to, _ := ctx.RoleIntf.GetBalance(param.To)

	err = from.SafeSub(param.Value)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = to.SafeAdd(param.Value)
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
	var err error
	param := &SetDelegateParam{}
	if err := abi.UmarshalStruct(ctx.Trx.Param, param); err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	_, err = ctx.RoleIntf.GetDelegateByAccountName(param.Name)
	if err != nil {
		// new delegate
		newdelegate := &role.Delegate{
			AccountName: param.Name,
			ReportKey:   param.Pubkey,
		}
		ctx.RoleIntf.SetDelegate(newdelegate.AccountName, newdelegate)

		//create schedule delegate vote role
		scheduleDelegate, err := ctx.RoleIntf.GetScheduleDelegate()
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}
		//create delegate vote role
		ctx.RoleIntf.CreateDelegateVotes()

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
	var err error
	param := &GrantCreditParam{}
	if err := abi.UmarshalStruct(ctx.Trx.Param, param); err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, param.Spender)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ctx.Trx.Sender)
	if cerr != ERROR_NONE {
		return cerr
	}

	// sender must be from
	if ctx.Trx.Sender != param.Name {
		return ERROR_CONT_ACCOUNT_MISMATCH
	}

	// check limit
	balance, err := ctx.RoleIntf.GetBalance(param.Name)
	if balance.Balance < param.Limit {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}

	credit := &role.TransferCredit{
		Name:    param.Name,
		Spender: param.Spender,
		Limit:   param.Limit,
	}
	err = ctx.RoleIntf.SetTransferCredit(credit.Name, credit)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func cancelCredit(ctx *Context) ContractError {
	var err error
	param := &CancelCreditParam{}
	if err := abi.UmarshalStruct(ctx.Trx.Param, param); err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, param.Spender)
	if cerr != ERROR_NONE {
		return cerr
	}

	_, err = ctx.RoleIntf.GetTransferCredit(param.Name, param.Spender)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	err = ctx.RoleIntf.DeleteTransferCredit(param.Name, param.Spender)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func transferFrom(ctx *Context) ContractError {
	var err error
	param := &TransferFromParam{}
	if err := abi.UmarshalStruct(ctx.Trx.Param, param); err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.From)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, param.To)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ctx.Trx.Sender)
	if cerr != ERROR_NONE {
		return cerr
	}

	// Note: sender is the spender
	// check limit
	credit, err := ctx.RoleIntf.GetTransferCredit(param.From, ctx.Trx.Sender)
	if err != nil {
		return ERROR_CONT_INSUFFICIENT_CREDITS
	}
	if param.Value > credit.Limit {
		return ERROR_CONT_INSUFFICIENT_CREDITS
	}

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(param.From)
	if from.Balance < param.Value {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}
	to, _ := ctx.RoleIntf.GetBalance(param.To)

	err = from.SafeSub(param.Value)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = credit.SafeSub(param.Value)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = to.SafeAdd(param.Value)
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

	if credit.Limit > 0 {
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
	var err error
	param := &DeployCodeParam{}
	if err := abi.UmarshalStruct(ctx.Trx.Param, param); err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check code
	err = checkCode(param.ContractCode)
	if err != nil {
		return ERROR_CONT_CODE_INVALID
	}

	codeHash := common.Sha256(param.ContractCode)

	account, _ := ctx.RoleIntf.GetAccount(param.Name)
	account.VMType = param.VMType
	account.VMVersion = param.VMVersion
	account.CodeVersion = codeHash
	account.ContractCode = make([]byte, len(param.ContractCode))
	copy(account.ContractCode, param.ContractCode)
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
	var err error
	param := &DeployABIParam{}
	if err := abi.UmarshalStruct(ctx.Trx.Param, param); err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check abi
	err = checkAbi(param.ContractAbi)
	if err != nil {
		return ERROR_CONT_ABI_PARSE_FAIL
	}

	account, _ := ctx.RoleIntf.GetAccount(param.Name)
	account.ContractAbi = make([]byte, len(param.ContractAbi))
	copy(account.ContractAbi, param.ContractAbi)
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
