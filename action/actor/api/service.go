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

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"

	bottosErr "github.com/bottos-project/bottos/common/errors"
	log "github.com/cihub/seelog"
	)

//ApiService is actor service
type ApiService struct {
	env *env.ActorEnv
}

//NewApiService new api service
func NewApiService(env *env.ActorEnv) api.ChainHandler {
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


func ConvertApiTrxToIntTrx(trx *api.Transaction) (*types.Transaction, error) {
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

func ConvertIntTrxToApiTrx(trx *types.Transaction) *api.Transaction {
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
func (a *ApiService) SendTransaction(ctx context.Context, trx *api.Transaction, resp *api.SendTransactionResponse) error {
	if trx == nil {
		//rsp.retCode = ??
		return nil
	}

	intTrx, err := ConvertApiTrxToIntTrx(trx)
	if err != nil {
		return nil
	}

	reqMsg := &message.PushTrxReq{
		Trx: intTrx,
	}

	handlerErr, err := trxactorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result() // await result

	if nil != err {
		resp.Errcode = uint32(bottosErr.ErrActorHandleError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrActorHandleError)

		log.Errorf("trx: %x actor process failed", intTrx.Hash(),)

		return nil
	}

	if bottosErr.ErrNoError == handlerErr {
		resp.Result = &api.SendTransactionResponse_Result{}
		resp.Result.TrxHash = intTrx.Hash().ToHexString()
		resp.Result.Trx = ConvertIntTrxToApiTrx(intTrx)
		resp.Msg = "trx receive succ"
		resp.Errcode = 0
	} else {
		resp.Result = &api.SendTransactionResponse_Result{}
		resp.Result.TrxHash = intTrx.Hash().ToHexString()
		resp.Result.Trx = ConvertIntTrxToApiTrx(intTrx)
		//resp.Msg = handlerErr.(string)GetCodeString
		//resp.Msg = "to be add detail error description"
		var tempErr bottosErr.ErrCode
		tempErr = handlerErr.(bottosErr.ErrCode)

		resp.Errcode = (uint32)(tempErr)
		resp.Msg = bottosErr.GetCodeString((bottosErr.ErrCode)(resp.Errcode))
	}

	log.Infof("trx: %v %s", resp.Result.TrxHash, resp.Msg)

	return nil
}

//GetTransaction query trx
func (a *ApiService) GetTransaction(ctx context.Context, req *api.GetTransactionRequest, resp *api.GetTransactionResponse) error {
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

	resp.Result = ConvertIntTrxToApiTrx(response.Trx)
	resp.Errcode = uint32(bottosErr.ErrNoError)
	return nil
}

//GetBlock query block
func (a *ApiService) GetBlock(ctx context.Context, req *api.GetBlockRequest, resp *api.GetBlockResponse) error {
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

	resp.Result = &api.GetBlockResponse_Result{}
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

//GetInfo query chain info
func (a *ApiService) GetInfo(ctx context.Context, req *api.GetInfoRequest, resp *api.GetInfoResponse) error {
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

	resp.Result = &api.GetInfoResponse_Result{}
	resp.Result.HeadBlockNum = response.HeadBlockNum
	resp.Result.LastConsensusBlockNum = response.LastConsensusBlockNum
	resp.Result.HeadBlockHash = response.HeadBlockHash.ToHexString()
	resp.Result.HeadBlockTime = response.HeadBlockTime
	resp.Result.HeadBlockDelegate = response.HeadBlockDelegate
	resp.Result.CursorLabel = response.HeadBlockHash.Label()
	resp.Errcode = 0
	return nil
}

//GetAccount query account info
func (a *ApiService) GetAccount(ctx context.Context, req *api.GetAccountRequest, resp *api.GetAccountResponse) error {
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

	resp.Result = &api.GetAccountResponse_Result{}
	resp.Result.AccountName = name
	resp.Result.Pubkey = common.BytesToHex(account.PublicKey)
	resp.Result.Balance = balance.Balance.String()
	resp.Result.StakedBalance = stakedBalance.StakedBalance.String()
	resp.Errcode = 0

	return nil
}

//GetKeyValue query contract object
func (a *ApiService) GetKeyValue(ctx context.Context, req *api.GetKeyValueRequest, resp *api.GetKeyValueResponse) error {
	contract := req.Contract
	object := req.Object
	key := req.Key
	value, err := a.env.RoleIntf.GetStrValue(contract, object, key)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiObjectNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiObjectNotFound)
		return nil
	}

	resp.Result = &api.GetKeyValueResponse_Result{}
	resp.Result.Contract = contract
	resp.Result.Object = object
	resp.Result.Key = key
	resp.Result.Value = common.BytesToHex([]byte(value))
	resp.Errcode = 0

	return nil
}

//GetAbi query contract abi info
func (a *ApiService) GetAbi(ctx context.Context, req *api.GetAbiRequest, resp *api.GetAbiResponse) error {
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
		return nil
	}

	return nil
}

//GetTransferCredit query trx credit info
func (a *ApiService) GetTransferCredit(ctx context.Context, req *api.GetTransferCreditRequest, resp *api.GetTransferCreditResponse) error {
	name := req.Name
	spender := req.Spender
	credit, err := a.env.RoleIntf.GetTransferCredit(name, spender)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrTransferCreditNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrTransferCreditNotFound)
		return nil
	}

	resp.Result = &api.GetTransferCreditResponse_Result{}
	resp.Result.Name = credit.Name
	resp.Result.Spender = credit.Spender
	resp.Result.Limit = credit.Limit.String()

	return nil
}
