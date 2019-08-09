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
	berr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/config"
)

const MaxDelegateLocationLen int = 32
const MaxDelegateDescriptionLen int = 128
const NativeContractExecTime uint64 = 100

//NativeContractInterface is native contract interface
type NativeContractInterface interface {
	IsNativeContractMethod(contract string, method string) bool
	ExecuteNativeContract(*Context) (berr.ErrCode, uint64, uint64)
}

//NativeContractMethod is native contract method
type NativeContractMethod func(*Context) berr.ErrCode

type NativeContractMethodConfig struct {
	handler  NativeContractMethod
	timeCost uint64
}

//NativeContract is native contract handler
type NativeContract struct {
	//Config map[string]NativeContractMethod
	Config map[string]NativeContractMethodConfig
}

//NewNativeContractHandler is native contract handler to handle different contracts
func NewNativeContractHandler() (NativeContractInterface, error) {
	nc := &NativeContract{
		Config: make(map[string]NativeContractMethodConfig),
	}

	nc.Config["newaccount"] = NativeContractMethodConfig{nc.newAccount, NativeContractExecTime}
	nc.Config["transfer"] = NativeContractMethodConfig{nc.transfer, NativeContractExecTime}
	nc.Config["grantcredit"] = NativeContractMethodConfig{nc.grantCredit, NativeContractExecTime}
	nc.Config["cancelcredit"] = NativeContractMethodConfig{nc.cancelCredit, NativeContractExecTime}
	nc.Config["transferfrom"] = NativeContractMethodConfig{nc.transferFrom, NativeContractExecTime}
	nc.Config["deploycontract"] = NativeContractMethodConfig{nc.deployContract, NativeContractExecTime}
	nc.Config["stake"] = NativeContractMethodConfig{nc.stake, NativeContractExecTime}
	nc.Config["unstake"] = NativeContractMethodConfig{nc.unstake, NativeContractExecTime}
	nc.Config["claim"] = NativeContractMethodConfig{nc.claim, NativeContractExecTime}
	nc.Config["regdelegate"] = NativeContractMethodConfig{nc.regDelegate, NativeContractExecTime}
	nc.Config["unregdelegate"] = NativeContractMethodConfig{nc.unregDelegate, NativeContractExecTime}
	nc.Config["votedelegate"] = NativeContractMethodConfig{nc.voteDelegate, NativeContractExecTime}
	nc.Config["newmsignaccount"] = NativeContractMethodConfig{nc.newMsignAccount, NativeContractExecTime}
	nc.Config["pushmsignproposal"] = NativeContractMethodConfig{nc.pushMsignProposal, NativeContractExecTime}
	nc.Config["approvemsignproposal"] = NativeContractMethodConfig{nc.approveMsignProposal, NativeContractExecTime}
	nc.Config["unapprovemsign"] = NativeContractMethodConfig{nc.unapproveMsignProposal, NativeContractExecTime}
	nc.Config["execmsignproposal"] = NativeContractMethodConfig{nc.execMsignProposal, NativeContractExecTime}
	nc.Config["cancelmsignproposal"] = NativeContractMethodConfig{nc.cancelMsignProposal, NativeContractExecTime}

	// genesis
	nc.Config["setdelegate"] = NativeContractMethodConfig{nc.setDelegate, NativeContractExecTime}
	nc.Config["unsetdelegate"] = NativeContractMethodConfig{nc.unsetDelegate, NativeContractExecTime}
	nc.Config["cancelgsperm"] = NativeContractMethodConfig{nc.cancelGsPermission, NativeContractExecTime}
	nc.Config["blkprodtrans"] = NativeContractMethodConfig{nc.blockProducingTransfer, NativeContractExecTime}
	nc.Config["settransitvote"] = NativeContractMethodConfig{nc.setTransitVote, NativeContractExecTime}
	nc.Config["claimreward"] = NativeContractMethodConfig{nc.claimReward, NativeContractExecTime}
	nc.Config["newstkaccount"] = NativeContractMethodConfig{nc.newStkAccount, NativeContractExecTime}

	return nc, nil
}

//IsNativeContractMethod is to check if the contract is native
func (nc *NativeContract) IsNativeContractMethod(contract string, method string) bool {
	if contract == config.BOTTOS_CONTRACT_NAME {
		if _, ok := nc.Config[method]; ok {
			return true
		}
	}
	return false
}

//ExecuteNativeContract is to call native contract
func (nc *NativeContract) ExecuteNativeContract(ctx *Context) (berr.ErrCode, uint64, uint64) {
	contract := ctx.Trx.Contract
	method := ctx.Trx.Method
	if nc.IsNativeContractMethod(contract, method) {
		if handler, ok := nc.Config[method]; ok {
			contErr := handler.handler(ctx)
			return contErr, uint64(len(ctx.Trx.Param)), handler.timeCost
		}
		return berr.ErrContractUnknownMethod, 0, 0
	}
	return berr.ErrContractUnknownContract, 0, 0

}
