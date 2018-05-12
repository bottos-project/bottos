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
 * file description:  trx agent
 * @Author:
 * @Date:   2017-12-13
 * @Last Modified by:
 * @Last Modified time:
 */

package apiactor

import (
	"time"
	"context"
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/api"
	"github.com/bottos-project/core/action/message"
	"github.com/bottos-project/core/action/env"
)

type ApiService struct {
	env *env.ActorEnv
}

func NewApiService(env *env.ActorEnv) api.CoreApiHandler {
	apiService := &ApiService{env:env}
	return apiService
}

var 	chainActorPid *actor.PID
func SetChainActorPid(tpid *actor.PID) {
	chainActorPid = tpid
}


var   trxactorPid *actor.PID

func SetTrxActorPid(tpid *actor.PID) {
	trxactorPid = tpid
}


func (a *ApiService) PushTrx(ctx context.Context, trx *types.Transaction, resp *api.PushTrxResponse) error {
	if trx == nil {
		//rsp.retCode = ??
		return nil
	}

	reqMsg := &message.PushTrxReq{
		Trx: trx,
		TrxSender : message.TrxSenderTypeFront,
	}
	
	handlerErr, err := trxactorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result() // await result

	if (nil != err) {
		resp.Errcode = 100
		resp.Msg = "message handle failed"    

		return nil
	}

	fmt.Println("handle result is ",handlerErr)

	if (nil == handlerErr) {
		resp.Result = &api.PushTrxResponse_Result{}
		resp.Result.TrxHash = trx.Hash().ToHexString()
		resp.Result.Trx = trx
		resp.Msg = "trx receive succ"
		resp.Errcode = 0
	} else {
		resp.Result = &api.PushTrxResponse_Result{}
		resp.Result.TrxHash = trx.Hash().ToHexString()
		resp.Result.Trx = trx
		//resp.Msg = handlerErr.(string)
		resp.Msg = "to be add detail error describtion"
		resp.Errcode = 100
	}

	return nil
}


func (a *ApiService) QueryTrx(ctx context.Context, req *api.QueryTrxRequest, resp *api.QueryTrxResponse) error {
	msgReq := &message.QueryTrxReq{
		TrxHash: common.HexToHash(req.TrxHash),
	}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Errcode = 1
		return nil
	}

	response := res.(*message.QueryTrxResp)
	if response.Trx == nil {
		resp.Errcode = 2
		resp.Msg = "Transaction not Found"
		return nil
	}

	resp.Result = response.Trx
	resp.Errcode = 0
	return nil
}

func (a *ApiService) QueryBlock(ctx context.Context, req *api.QueryBlockRequest, resp *api.QueryBlockResponse) error {
	msgReq := &message.QueryBlockReq{
		BlockHash: common.HexToHash(req.BlockHash),
		BlockNumber: req.BlockNum,
	}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Errcode = 1
		return nil
	}

	response := res.(*message.QueryBlockResp)
	if response.Block == nil {
		resp.Errcode = 2
		resp.Msg = "Block not Found"
		return nil
	}

	resp.Result = &api.QueryBlockResponse_Result{}
	hash := response.Block.Hash()
	resp.Result.PrevBlockHash = response.Block.GetPrevBlockHash().ToHexString()
	resp.Result.BlockNum = response.Block.GetNumber()
	resp.Result.BlockHash = hash.ToHexString()
	resp.Result.CursorBlockLabel = hash.Label()
	resp.Result.BlockTime = response.Block.GetTimestamp()
	resp.Result.TrxMerkleRoot = response.Block.ComputeMerkleRoot().ToHexString()
	resp.Result.Delegate = string(response.Block.GetDelegate())
	resp.Result.DelegateSign = response.Block.GetDelegateSign().ToHexString()

	resp.Errcode = 0
	return nil
}

func (h *ApiService) QueryChainInfo(ctx context.Context, req *api.QueryChainInfoRequest, resp *api.QueryChainInfoResponse) error {
	msgReq := &message.QueryChainInfoReq{}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Errcode = 1
		return nil
	}

	response := res.(*message.QueryChainInfoResp)
	if response.Error != nil {
		resp.Errcode = 2
		return nil
	}

	resp.Result = &api.QueryChainInfoResponse_Result{}
	resp.Result.HeadBlockNum = response.HeadBlockNum
	resp.Result.LastConsensusBlockNum = response.LastConsensusBlockNum
	resp.Result.HeadBlockHash = response.HeadBlockHash.ToHexString()
	resp.Result.HeadBlockTime = response.HeadBlockTime
	resp.Result.HeadBlockDelegate = response.HeadBlockDelegate
	resp.Result.CursorLabel = response.HeadBlockHash.Label()
	resp.Errcode = 0
	return nil
}

func (h *ApiService) QueryAccount(ctx context.Context, req *api.QueryAccountRequest, resp *api.QueryAccountResponse) error {
	name := req.AccountName
	account, err := h.env.RoleIntf.GetAccount(name)
	if err != nil {
		resp.Errcode = 1
		resp.Msg = "Account Not Found"
		return nil
	}

	balance, err := h.env.RoleIntf.GetBalance(name)
	if err != nil {
		resp.Errcode = 1
		resp.Msg = "Balance Not Found"
		return nil
	}

	stakedBalance, err := h.env.RoleIntf.GetStakedBalance(name)
	if err != nil {
		resp.Errcode = 1
		resp.Msg = "Staked Balance Not Found"
		return nil
	}

	resp.Result = &api.QueryAccountResponse_Result{}
	resp.Result.AccountName = name
	resp.Result.Pubkey = string(account.PublicKey);
	resp.Result.Balance = balance.Balance
	resp.Result.StakedBalance = stakedBalance.StakedBalance
	resp.Errcode = 0

	return nil
}


func (h *ApiService) QueryObject(ctx context.Context, req *api.QueryObjectReq, resp *api.QueryObjectResponse) error {
	contract := req.Contract
	object := req.Object
	key := req.Key
	value, err := h.env.ContractDB.GetStrValue(contract, object, key)
	if err != nil {
		resp.Errcode = 1
		resp.Msg = "KeyValue Not Found"
		return nil
	}

	resp.Result = &api.QueryObjectResponse_Result{}
	resp.Result.Contract = contract
	resp.Result.Object = object
	resp.Result.Key = key
	resp.Result.Value = value
	resp.Errcode = 0

	return nil
}
