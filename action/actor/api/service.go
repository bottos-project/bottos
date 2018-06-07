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
	"context"
	"time"
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"

	bottosErr "github.com/bottos-project/bottos/common/errors"
)

//ApiService is actor service
type ApiService struct {
	env *env.ActorEnv
}

//NewApiService new api service
func NewApiService(env *env.ActorEnv) api.CoreApiHandler {
	apiService := &ApiService{env: env}
	return apiService
}

var chainActorPid *actor.PID

//SetChainActorPid set chain actor pid
func SetChainActorPid(tpid *actor.PID) {
	chainActorPid = tpid
}

var trxactorPid *actor.PID

//SetTrxActorPid set trx actor pid
func SetTrxActorPid(tpid *actor.PID) {
	trxactorPid = tpid
}

func convertApiTrxToIntTrx(trx *api.Transaction) (*types.Transaction, error) {
	param, err := common.HexToBytes(trx.Param)
	if err != nil {
		return nil, err
	}

	signature, err := common.HexToBytes(trx.Signature)
	if err != nil {
		return nil, err
	}

	intTrx := &types.Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       param,
		SigAlg:      trx.SigAlg,
		Signature:   signature,
	}

	return intTrx, nil
}

func convertIntTrxToApiTrx(trx *types.Transaction) *api.Transaction {
	apiTrx := &api.Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       common.BytesToHex(trx.Param),
		SigAlg:      trx.SigAlg,
		Signature:   common.BytesToHex(trx.Signature),
	}

	return apiTrx
}

//PushTrx push trx
func (a *ApiService) PushTrx(ctx context.Context, trx *api.Transaction, resp *api.PushTrxResponse) error {
	if trx == nil {
		//rsp.retCode = ??
		return nil
	}

	intTrx, err := convertApiTrxToIntTrx(trx)
	if err != nil {
		return nil
	}

	reqMsg := &message.PushTrxReq{
		Trx:       intTrx,
		TrxSender: message.TrxSenderTypeFront,
	}

	handlerErr, err := trxactorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result() // await result

	if nil != err {
		resp.Errcode = uint32(bottosErr.ErrActorHandleError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrActorHandleError)

		return nil
	}

	if bottosErr.ErrNoError == handlerErr {
		resp.Result = &api.PushTrxResponse_Result{}
		resp.Result.TrxHash = intTrx.Hash().ToHexString()
		resp.Result.Trx = convertIntTrxToApiTrx(intTrx)
		resp.Msg = "trx receive succ"
		resp.Errcode = 0
	} else {
		resp.Result = &api.PushTrxResponse_Result{}
		resp.Result.TrxHash = intTrx.Hash().ToHexString()
		resp.Result.Trx = convertIntTrxToApiTrx(intTrx)
		//resp.Msg = handlerErr.(string)GetCodeString
		//resp.Msg = "to be add detail error description"
		var tempErr bottosErr.ErrCode
		tempErr = handlerErr.(bottosErr.ErrCode)

		resp.Errcode = (uint32)(tempErr)
		resp.Msg = bottosErr.GetCodeString((bottosErr.ErrCode)(resp.Errcode))
	}
	
	fmt.Println("trx: ", resp.Result.TrxHash, resp.Msg)

	return nil
}

//QueryTrx query trx
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
		resp.Errcode = uint32(bottosErr.ErrApiTrxNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiTrxNotFound)
		return nil
	}

	resp.Result = convertIntTrxToApiTrx(response.Trx)
	resp.Errcode = uint32(bottosErr.ErrNoError)
	return nil
}

//QueryBlock query block
func (a *ApiService) QueryBlock(ctx context.Context, req *api.QueryBlockRequest, resp *api.QueryBlockResponse) error {
	msgReq := &message.QueryBlockReq{
		BlockHash:   common.HexToHash(req.BlockHash),
		BlockNumber: req.BlockNum,
	}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiBlockNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiBlockNotFound)
		return nil
	}

	response := res.(*message.QueryBlockResp)
	if response.Block == nil {
		resp.Errcode = uint32(bottosErr.ErrApiBlockNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiBlockNotFound)
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

//QueryChainInfo query chain info
func (a *ApiService) QueryChainInfo(ctx context.Context, req *api.QueryChainInfoRequest, resp *api.QueryChainInfoResponse) error {
	msgReq := &message.QueryChainInfoReq{}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
		return nil
	}

	response := res.(*message.QueryChainInfoResp)
	if response.Error != nil {
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
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

//QueryAccount query account info
func (a *ApiService) QueryAccount(ctx context.Context, req *api.QueryAccountRequest, resp *api.QueryAccountResponse) error {
	name := req.AccountName
	account, err := a.env.RoleIntf.GetAccount(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		return nil
	}

	balance, err := a.env.RoleIntf.GetBalance(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		return nil
	}

	stakedBalance, err := a.env.RoleIntf.GetStakedBalance(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		return nil
	}

	resp.Result = &api.QueryAccountResponse_Result{}
	resp.Result.AccountName = name
	resp.Result.Pubkey = common.BytesToHex(account.PublicKey)
	resp.Result.Balance = balance.Balance
	resp.Result.StakedBalance = stakedBalance.StakedBalance
	resp.Errcode = 0

	return nil
}

//QueryObject query contract object
func (a *ApiService) QueryObject(ctx context.Context, req *api.QueryObjectReq, resp *api.QueryObjectResponse) error {
	contract := req.Contract
	object := req.Object
	key := req.Key
	value, err := a.env.ContractDB.GetStrValue(contract, object, key)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiObjectNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiObjectNotFound)
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

//QueryAbi query contract abi info
func (a *ApiService) QueryAbi(ctx context.Context, req *api.QueryAbiReq, resp *api.QueryAbiResponse) error {
	contract := req.Contract
	account, err := a.env.RoleIntf.GetAccount(contract)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		return nil
	}

	if len(account.ContractAbi) > 0 {
		resp.Result = string(account.ContractAbi)
	} else {
		// TODO
		return nil
	}

	return nil
}

//QueryTransferCredit query trx credit info
func (a *ApiService) QueryTransferCredit(ctx context.Context, req *api.QueryTransferCreditRequest, resp *api.QueryTransferCreditResponse) error {
	name := req.Name
	spender := req.Spender
	credit, err := a.env.RoleIntf.GetTransferCredit(name, spender)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrTransferCreditNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrTransferCreditNotFound)
		return nil
	}

	resp.Result = &api.QueryTransferCreditResponse_Result{}
	resp.Result.Name = credit.Name
	resp.Result.Spender = credit.Spender
	resp.Result.Limit = credit.Limit

	return nil
}
