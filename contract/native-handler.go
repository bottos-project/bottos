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
		"math/big"
	"fmt"
	"github.com/bottos-project/bottos/common"
	berr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/vm"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/role"
	berr "github.com/bottos-project/bottos/common/errors"
	)

func (nc *NativeContract) newAccount(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	newaccount := abi.UnmarshalAbiEx("bottos", Abi, "newaccount", ctx.Trx.Param)
	if newaccount == nil || len(newaccount) <= 0 {
		return berr.ErrContractParamParseError
	}
	
	NewaccountName   := newaccount["name"].(string)
	NewaccountPubKey := newaccount["pubkey"].(string)

	if len(NewaccountPubKey) != config.PUBKEY_LEN {

		return berr.ErrAccountPubkeyLenIllegal
	}

	//log.Errorf("test new account %s, len is %d, stand len is %d\n", NewaccountPubKey, len(NewaccountPubKey), config.PUBKEY_LEN)

	//check account
	cerr := nc.checkAccountName(NewaccountName)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if nc.isAccountNameExist(ctx.RoleIntf, NewaccountName) {
		return berr.ErrAccountAlreadyExist
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
		AccountName:        NewaccountName,
		StakedBalance:      big.NewInt(0),
		StakedSpaceBalance: big.NewInt(0),
		StakedTimeBalance:  big.NewInt(0),
		UnstakingBalance:   big.NewInt(0),
		LastUnstakingTime:  0,
	}
	ctx.RoleIntf.SetStakedBalance(NewaccountName, stakedBalance)

	// 4, create ResourceUsage
	resourceUsage := &role.ResourceUsage{
		AccountName:                NewaccountName,
		PledgedSpaceTokenUsedInWin: 0,
		PledgedTimeTokenUsedInWin:  0,
		FreeTimeTokenUsedInWin:     0,
		FreeSpaceTokenUsedInWin:    0,
		LastSpaceCursorBlock:       0,
		LastTimeCursorBlock:        0,
	}
	ctx.RoleIntf.SetResourceUsage(NewaccountName, resourceUsage)

	return berr.ErrNoError
}

func (nc *NativeContract) checkSigner(account string, expected string) bool {
	return account == expected
}

func (nc *NativeContract) transfer(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	transfer := abi.UnmarshalAbiEx("bottos", Abi, "transfer", ctx.Trx.Param)
	if transfer == nil || len(transfer) <= 0 {
		return berr.ErrContractParamParseError
	}
	
	FromWhom := transfer["from"].(string)
	ToWhom   := transfer["to"].(string)
	TransValue := transfer["value"].(*big.Int)
	
	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, FromWhom)
	if cerr != berr.ErrNoError {
		return cerr
	}

	cerr = nc.checkAccount(ctx.RoleIntf, ToWhom)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if FromWhom == ToWhom {
		return berr.ErrContractTransferToSelf
	}

	if !nc.checkSigner(FromWhom, ctx.Trx.Sender) {
		return berr.ErrContractAccountMismatch
	}

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(FromWhom)
	if -1 == from.Balance.Cmp(TransValue) {
		return berr.ErrContractInsufficientFunds
	}
	to, _ := ctx.RoleIntf.GetBalance(ToWhom)

	err := from.SafeSub(TransValue)
	if err != nil {
		return berr.ErrContractTransferOverflow
	}
	err = to.SafeAdd(TransValue)
	if err != nil {
		return berr.ErrContractTransferOverflow
	}

	err = ctx.RoleIntf.SetBalance(from.AccountName, from)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}
	err = ctx.RoleIntf.SetBalance(to.AccountName, to)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}
	
	return berr.ErrNoError
}

func (nc *NativeContract) setDelegate(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "setdelegate", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return berr.ErrContractParamParseError
	}
	
	ParamName   := param["name"].(string)
	ParamPubKey := param["pubkey"].(string)
	location := param["location"].(string)
	description := param["description"].(string)

	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, ParamName)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if len(location) > MaxDelegateLocationLen {
		return berr.ErrTrxContractHanldeError
	}

	if len(description) > MaxDelegateDescriptionLen {
		return berr.ErrTrxContractHanldeError
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
			return berr.ErrTrxContractHanldeError
		}

		delegateVote := &role.DelegateVotes{
			OwnerAccount: newdelegate.AccountName,
			Serve : role.Serve{
				Votes: big.NewInt(0),
				Position: big.NewInt(0),
				TermUpdateTime: big.NewInt(0),
				TermFinishTime: big.NewInt(0),
			},
		}
		delegateVote.OwnerAccount = newdelegate.AccountName
		delegateVote.StartNewTerm(scheduleDelegate.CurrentTermTime)
		err = ctx.RoleIntf.SetDelegateVotes(newdelegate.AccountName, delegateVote)
		if err != nil {
			return berr.ErrTrxContractHanldeError
		}
	} else {
		return berr.ErrTrxContractHanldeError
	}

	return berr.ErrNoError
}

func (nc *NativeContract) grantCredit(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "grantcredit", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return berr.ErrContractParamParseError
	}
	
	ParamName    := param["name"].(string)
	ParamSpender := param["spender"].(string)
	ParamLimit   := param["limit"].(*big.Int)

	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, ParamName)
	if cerr != berr.ErrNoError {
		return cerr
	}

	cerr = nc.checkAccount(ctx.RoleIntf, ParamSpender)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if ParamName == ParamSpender {
		return berr.ErrContractGrantToSelf
	}

	cerr = nc.checkAccount(ctx.RoleIntf, ctx.Trx.Sender)
	if cerr != berr.ErrNoError {
		return cerr
	}

	// sender must be from
	if !nc.checkSigner(ParamName, ctx.Trx.Sender) {
		return berr.ErrContractAccountMismatch
	}

	// check limit
	balance, err := ctx.RoleIntf.GetBalance(ParamName)
	if -1 == balance.Balance.Cmp(ParamLimit) {
		return berr.ErrContractInsufficientCredits
	}

	credit := &role.TransferCredit{
		Name:    ParamName,
		Spender: ParamSpender,
		Limit:   ParamLimit,  
	}
	err = ctx.RoleIntf.SetTransferCredit(credit.Name, credit)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	return berr.ErrNoError
}

func (nc *NativeContract) cancelCredit(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "cancelcredit", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return berr.ErrContractParamParseError
	}

	ParamName    := param["name"].(string)
	ParamSpender := param["spender"].(string)
	
	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, ParamName)
	if cerr != berr.ErrNoError {
		return cerr
	}

	cerr = nc.checkAccount(ctx.RoleIntf, ParamSpender)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if !nc.checkSigner(ParamName, ctx.Trx.Sender) {
		return berr.ErrContractAccountMismatch
	}

	_, err := ctx.RoleIntf.GetTransferCredit(ParamName, ParamSpender)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	err = ctx.RoleIntf.DeleteTransferCredit(ParamName, ParamSpender)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	return berr.ErrNoError
}

func (nc *NativeContract) transferFrom(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	transfer := abi.UnmarshalAbiEx("bottos", Abi, "transferfrom", ctx.Trx.Param)
	if transfer == nil || len(transfer) <= 0 {
		return berr.ErrContractParamParseError
	}
	TransFrom := transfer["from"].(string)
	TransTo := transfer["to"].(string)
	TransValue := transfer["value"].(*big.Int)
	
	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, TransFrom)
	if cerr != berr.ErrNoError {
		return cerr
	}

	cerr = nc.checkAccount(ctx.RoleIntf, TransTo)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if TransFrom == TransTo {
		return berr.ErrContractTransferToSelf
	}

	cerr = nc.checkAccount(ctx.RoleIntf, ctx.Trx.Sender)
	if cerr != berr.ErrNoError {
		return cerr
	}

	// Note: sender is the spender
	// check limit
	credit, err := ctx.RoleIntf.GetTransferCredit(TransFrom, ctx.Trx.Sender)
	if err != nil {
		return berr.ErrContractInsufficientCredits
	}
	if 1 == TransValue.Cmp(credit.Limit) {
		return berr.ErrContractInsufficientCredits
	}

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(TransFrom)
	if -1 == from.Balance.Cmp(TransValue) {
		return berr.ErrContractInsufficientFunds
	}
	to, _ := ctx.RoleIntf.GetBalance(TransTo)

	err = from.SafeSub(TransValue)
	if err != nil {
		return berr.ErrContractTransferOverflow
	}
	err = credit.SafeSub(TransValue)
	if err != nil {
		return berr.ErrContractTransferOverflow
	}
	err = to.SafeAdd(TransValue)
	if err != nil {
		return berr.ErrContractTransferOverflow
	}

	err = ctx.RoleIntf.SetBalance(from.AccountName, from)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}
	err = ctx.RoleIntf.SetBalance(to.AccountName, to)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	if 1 == credit.Limit.Cmp(big.NewInt(0)) {
		err = ctx.RoleIntf.SetTransferCredit(credit.Name, credit)
		if err != nil {
			return berr.ErrTrxContractHanldeError
		}
	} else {
		err = ctx.RoleIntf.DeleteTransferCredit(credit.Name, ctx.Trx.Sender)
		if err != nil {
			return berr.ErrTrxContractHanldeError
		}
	}

	return berr.ErrNoError
}

func (nc *NativeContract) checkCode(code []byte) error {
	// TODO
	return nil
}

func (nc *NativeContract) deployCode(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "deploycode", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return berr.ErrContractParamParseError
	}
	
	ParamContract := param["contract"].(string)
	ParamContractCode, _ := common.HexToBytes(param["contract_code"].(string))
	ParamVmType := param["vm_type"].(byte)

	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, ParamContract)
	if cerr != berr.ErrNoError {
		return cerr
	}

	// check code
	err := nc.checkCode(ParamContractCode)
	if err != nil {
		return berr.ErrContractInvalidContractCode
	}

	codeHash := common.Sha256(ParamContractCode)

	account, _ := ctx.RoleIntf.GetAccount(ParamContract)
	account.VMType = byte(ParamVmType)
	account.CodeVersion = codeHash
	account.ContractCode = make([]byte, len(ParamContractCode))
	copy(account.ContractCode, ParamContractCode)
	err = ctx.RoleIntf.SetAccount(account.AccountName, account)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	return berr.ErrNoError
}

func (nc *NativeContract) checkAbi(abiRaw []byte) error {
	_, err := abi.ParseAbi(abiRaw)
	if err != nil {
		return fmt.Errorf("ABI Parse error: %v", err)
	}
	return nil
}

func (nc *NativeContract) deployAbi(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "deployabi", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return berr.ErrContractParamParseError
	}
	
	ParamContract := param["contract"].(string)
	ParamContractAbi, _ := common.HexToBytes(param["contract_abi"].(string))
	
	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, ParamContract)
	if cerr != berr.ErrNoError {
		return cerr
	}

	// check abi
	err := nc.checkAbi(ParamContractAbi)
	if err != nil {
		return berr.ErrContractInvalidContractAbi
	}

	account, _ := ctx.RoleIntf.GetAccount(ParamContract)
	account.ContractAbi = make([]byte, len(ParamContractAbi))
	copy(account.ContractAbi, ParamContractAbi)
	err = ctx.RoleIntf.SetAccount(account.AccountName, account)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	return berr.ErrNoError
}

func (nc *NativeContract) stake(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	transfer := abi.UnmarshalAbiEx("bottos", Abi, "stake", ctx.Trx.Param)
	if transfer == nil || len(transfer) <= 0 {
		return berr.ErrContractParamParseError
	}

	amount := transfer["amount"].(*big.Int)

	// check account
	if errcode := nc.checkAccount(ctx.RoleIntf, ctx.Trx.Sender); errcode != berr.ErrNoError {
		return errcode
	}

	// amount should more than 0
	if 1 != amount.Cmp(big.NewInt(0)) {
		return berr.ErrTrxContractHanldeError
	}

	// check funds
	balance, _ := ctx.RoleIntf.GetBalance(ctx.Trx.Sender)
	if -1 == balance.Balance.Cmp(amount) {
		return berr.ErrContractInsufficientFunds
	}
	sb, _ := ctx.RoleIntf.GetStakedBalance(ctx.Trx.Sender)

	if err := balance.SafeSub(amount); err != nil {
		return berr.ErrContractTransferOverflow
	}
	if err := sb.SafeAdd(amount, target); err != nil {
		return berr.ErrContractTransferOverflow
	}

	if err := ctx.RoleIntf.SetBalance(balance.AccountName, balance); err != nil {
		return berr.ErrTrxContractHanldeError
	}
	if err := ctx.RoleIntf.SetStakedBalance(sb.AccountName, sb); err != nil {
		return berr.ErrTrxContractHanldeError
	}

	//update  AllStakedSpaceBalance& AllStakedTimeBalance
	cs, _ := ctx.RoleIntf.GetChainState()
	if err := cs.SafeAdd(amount, target); err != nil {
		return berr.ErrContractTransferOverflow
	}
	if err := ctx.RoleIntf.SetChainState(cs); err != nil {
		return berr.ErrTrxContractHanldeError
	}

	voter, _ := ctx.RoleIntf.GetVoter(ctx.Trx.Sender)
	if voter == nil {
		voter := &role.Voter{
			Owner: ctx.Trx.Sender,
			Delegate: string(""),
		}
		if err := ctx.RoleIntf.SetVoter(voter.Owner, voter); err != nil {
			return berr.ErrTrxContractHanldeError
		}
	} else {
		if voter.Delegate != "" {
			delegateVote, err := ctx.RoleIntf.GetDelegateVotes(voter.Delegate)
			if err != nil {
				return berr.ErrTrxContractHanldeError
			}
			sd, err := ctx.RoleIntf.GetScheduleDelegate()
			if err != nil {
				return berr.ErrTrxContractHanldeError
			}
			delegateVote.UpdateVotes(amount, sd.CurrentTermTime)
			if err := ctx.RoleIntf.SetDelegateVotes(delegateVote.OwnerAccount, delegateVote); err != nil {
				return berr.ErrTrxContractHanldeError
			}
		}
	}

	return berr.ErrNoError
}

func (nc *NativeContract) unstake(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	transfer := abi.UnmarshalAbiEx("bottos", Abi, "unstake", ctx.Trx.Param)
	if transfer == nil || len(transfer) <= 0 {
		return berr.ErrContractParamParseError
	}

	amount := transfer["amount"].(*big.Int)

	// check account
	if errcode := nc.checkAccount(ctx.RoleIntf, ctx.Trx.Sender); errcode != berr.ErrNoError {
		return errcode
	}

	// amount should more than 0
	if 1 != amount.Cmp(big.NewInt(0)) {
		return berr.ErrTrxContractHanldeError
	}

	// check funds
	sb, _ := ctx.RoleIntf.GetStakedBalance(ctx.Trx.Sender)
	if -1 == sb.StakedBalance.Cmp(amount) {
		return berr.ErrContractInsufficientFunds
	}
	oldStakeAmount := sb.StakedBalance

	if err := sb.UnstakingAmount(amount); err != nil {
		return berr.ErrContractTransferOverflow
	}

	chainState, _ := ctx.RoleIntf.GetChainState()
	sb.LastUnstakingTime = chainState.LastBlockTime
	if err := ctx.RoleIntf.SetStakedBalance(sb.AccountName, sb); err != nil {
		return berr.ErrTrxContractHanldeError
	}

	voter, _ := ctx.RoleIntf.GetVoter(ctx.Trx.Sender)
	if voter != nil {
		if voter.Delegate != "" {
			delegateVote, err := ctx.RoleIntf.GetDelegateVotes(voter.Delegate)
			if err != nil {
				return berr.ErrTrxContractHanldeError
			}
			sd, err := ctx.RoleIntf.GetScheduleDelegate()
			if err != nil {
				return berr.ErrTrxContractHanldeError
			}
			delta := big.NewInt(0)
			delta.Sub(sb.StakedBalance, oldStakeAmount)
			delegateVote.UpdateVotes(delta, sd.CurrentTermTime)

			if err := ctx.RoleIntf.SetDelegateVotes(delegateVote.OwnerAccount, delegateVote); err != nil {
				return berr.ErrTrxContractHanldeError
			}
		}
	}

	return berr.ErrNoError
}

func (nc *NativeContract) claim(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	transfer := abi.UnmarshalAbiEx("bottos", Abi, "claim", ctx.Trx.Param)
	if transfer == nil || len(transfer) <= 0 {
		return berr.ErrContractParamParseError
	}

	amount := transfer["amount"].(*big.Int)

	// check account
	if errcode := nc.checkAccount(ctx.RoleIntf, ctx.Trx.Sender); errcode != berr.ErrNoError {
		return errcode
	}

	// check funds
	sb, _ := ctx.RoleIntf.GetStakedBalance(ctx.Trx.Sender)
	releaseTime := sb.LastUnstakingTime + config.UNSTAKING_BALANCE_DURATION
	chainState, _ := ctx.RoleIntf.GetChainState()
	if chainState.LastBlockTime < releaseTime {
		return berr.ErrTrxContractHanldeError
	}

	if -1 == sb.UnstakingBalance.Cmp(amount) {
		return berr.ErrContractInsufficientFunds
	}

	balance, _ := ctx.RoleIntf.GetBalance(ctx.Trx.Sender)
	if err := balance.SafeAdd(amount); err != nil {
		return berr.ErrContractTransferOverflow
	}
	if err := sb.Claim(amount); err != nil {
		return berr.ErrContractTransferOverflow
	}

	if err := ctx.RoleIntf.SetBalance(balance.AccountName, balance); err != nil {
		return berr.ErrTrxContractHanldeError
	}
	if err := ctx.RoleIntf.SetStakedBalance(sb.AccountName, sb); err != nil {
		return berr.ErrTrxContractHanldeError
	}

	return berr.ErrNoError
}

func (nc *NativeContract) regDelegate(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "regdelegate", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return berr.ErrContractParamParseError
	}

	ParamName   := param["name"].(string)
	ParamPubKey := param["pubkey"].(string)
	location := param["location"].(string)
	description := param["description"].(string)

	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, ParamName)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if len(location) > MaxDelegateLocationLen {
		return berr.ErrTrxContractHanldeError
	}

	if len(description) > MaxDelegateDescriptionLen {
		return berr.ErrTrxContractHanldeError
	}

	delegate, err := ctx.RoleIntf.GetDelegateByAccountName(ParamName)
	if err != nil {
		// new delegate
		newdelegate := &role.Delegate{
			AccountName: ParamName,
			ReportKey:   ParamPubKey,
			Location: location,
			Description: description,
			Active: true,
		}
		if err := ctx.RoleIntf.SetDelegate(newdelegate.AccountName, newdelegate); err != nil {
			return berr.ErrTrxContractHanldeError
		}

		//create schedule delegate vote role
		scheduleDelegate, err := ctx.RoleIntf.GetScheduleDelegate()
		if err != nil {
			return berr.ErrTrxContractHanldeError
		}

		delegateVote := &role.DelegateVotes{
			OwnerAccount: newdelegate.AccountName,
			Serve : role.Serve{
				Votes: big.NewInt(0),
				Position: big.NewInt(0),
				TermUpdateTime: big.NewInt(0),
				TermFinishTime: big.NewInt(0),
			},
		}
		delegateVote.StartNewTerm(scheduleDelegate.CurrentTermTime)
		err = ctx.RoleIntf.SetDelegateVotes(newdelegate.AccountName, delegateVote)
		if err != nil {
			return berr.ErrTrxContractHanldeError
		}
	} else {
		delegate.ReportKey = ParamPubKey
		delegate.Active = true
		ctx.RoleIntf.SetDelegate(delegate.AccountName, delegate)
	}

	return berr.ErrNoError
}

func (nc *NativeContract) unregDelegate(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "unregdelegate", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return berr.ErrContractParamParseError
	}

	ParamName   := param["name"].(string)

	// check account
	cerr := nc.checkAccount(ctx.RoleIntf, ParamName)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if !nc.checkSigner(ParamName, ctx.Trx.Sender) {
		return berr.ErrContractAccountMismatch
	}

	delegate, err := ctx.RoleIntf.GetDelegateByAccountName(ParamName)
	if err == nil {
		// new delegate
		delegate.Active = false
		if err := ctx.RoleIntf.SetDelegate(delegate.AccountName, delegate); err != nil {
			return berr.ErrTrxContractHanldeError
		}
	} else {
		return berr.ErrTrxContractHanldeError
	}

	return berr.ErrNoError
}

func (nc *NativeContract) voteDelegate(ctx *Context) berr.ErrCode {
	Abi := abi.GetAbi()
	param := abi.UnmarshalAbiEx("bottos", Abi, "votedelegate", ctx.Trx.Param)
	if param == nil || len(param) <= 0 {
		return berr.ErrContractParamParseError
	}

	voteop := param["voteop"].(uint8)
	voterName := param["voter"].(string)
	delegateName := param["delegate"].(string)

	if !nc.checkSigner(voterName, ctx.Trx.Sender) {
		return berr.ErrContractAccountMismatch
	}

	if errcode := nc.checkAccount(ctx.RoleIntf, voterName); errcode != berr.ErrNoError {
		return errcode
	}

	voter, err := ctx.RoleIntf.GetVoter(voterName)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	sb, err := ctx.RoleIntf.GetStakedBalance(voterName)
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	// staked balance should more than 0
	if 1 != sb.StakedBalance.Cmp(big.NewInt(0)) {
		return berr.ErrContractInsufficientFunds
	}

	sd, err := ctx.RoleIntf.GetScheduleDelegate()
	if err != nil {
		return berr.ErrTrxContractHanldeError
	}

	if voteop == 1 {
		// vote
		if errcode := nc.checkAccount(ctx.RoleIntf, delegateName); errcode != berr.ErrNoError {
			return errcode
		}

		if voter.Delegate != "" {
			oldDelegateVote, err := ctx.RoleIntf.GetDelegateVotes(voter.Delegate)
			if err != nil {
				return berr.ErrTrxContractHanldeError
			}
			voteStake := big.NewInt(0).Set(sb.StakedBalance)
			voteStake.Mul(voteStake, big.NewInt(-1))
			oldDelegateVote.UpdateVotes(voteStake, sd.CurrentTermTime)

			if err := ctx.RoleIntf.SetDelegateVotes(oldDelegateVote.OwnerAccount, oldDelegateVote); err != nil {
				return berr.ErrTrxContractHanldeError
			}
		}

		delegateVote, err := ctx.RoleIntf.GetDelegateVotes(delegateName)
		if err != nil {
			return berr.ErrTrxContractHanldeError
		}
		delegateVote.UpdateVotes(sb.StakedBalance, sd.CurrentTermTime)

		voter.Delegate = delegateName
		if err := ctx.RoleIntf.SetVoter(voterName, voter); err != nil {
			return berr.ErrTrxContractHanldeError
		}

		if err := ctx.RoleIntf.SetDelegateVotes(delegateVote.OwnerAccount, delegateVote); err != nil {
			return berr.ErrTrxContractHanldeError
		}
	} else if voteop == 0 {
		// cancel vote
		// staked balance should more than 0
		if 1 != sb.StakedBalance.Cmp(big.NewInt(0)) {
			return berr.ErrContractInsufficientFunds
		}

		if voter.Delegate != "" {
			oldDelegateVote, err := ctx.RoleIntf.GetDelegateVotes(voter.Delegate)
			if err != nil {
				return berr.ErrTrxContractHanldeError
			}
			voteStake := big.NewInt(0).Set(sb.StakedBalance)
			voteStake.Mul(voteStake, big.NewInt(-1))
			oldDelegateVote.UpdateVotes(voteStake, sd.CurrentTermTime)

			voter.Delegate = ""
			if err := ctx.RoleIntf.SetVoter(voterName, voter); err != nil {
				return berr.ErrTrxContractHanldeError
			}

			if err := ctx.RoleIntf.SetDelegateVotes(oldDelegateVote.OwnerAccount, oldDelegateVote); err != nil {
				return berr.ErrTrxContractHanldeError
			}
		} else {
			return berr.ErrTrxContractHanldeError
		}
	}

	return berr.ErrNoError
}

func (nc *NativeContract) checkAccountName(name string) berr.ErrCode {
	if len(name) == 0 {
		return berr.ErrContractAccountNameIllegal
	}

	if len(name) > common.MaxNameLength {
		return berr.ErrContractAccountNameIllegal
	}

	if !nc.checkAccountNameContent(name) {
		return berr.ErrContractAccountNameIllegal
	}

	return berr.ErrNoError
}

func (nc *NativeContract) checkAccountNameContent(name string) bool {
	return nc.re.MatchString(name)
}

func (nc *NativeContract) isAccountNameExist(RoleIntf role.RoleInterface, name string) bool {
	account, err := RoleIntf.GetAccount(name)
	if err == nil {
		if account != nil && account.AccountName == name {
			return true
		}
	}
	return false
}

func (nc *NativeContract) checkAccount(RoleIntf role.RoleInterface, name string) berr.ErrCode {
	cerr := nc.checkAccountName(name)
	if cerr != berr.ErrNoError {
		return cerr
	}

	if !nc.isAccountNameExist(RoleIntf, name) {
		return berr.ErrContractAccountNotFound
	}

	return berr.ErrNoError
}
