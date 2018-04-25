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


func (a *ApiService) PushTx(ctx context.Context, trx *types.Transaction, resp *api.PushTxResponse) error {
	if trx == nil {
		//rsp.retCode = ??
		return nil
	}
	
	reqMsg := &message.PushTrxReq{
		Trx: trx,
		TrxSender : message.TrxSenderTypeFront,
	}
	_, err := trxactorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result() // await result
	
	resp.Tx = trx

	if (nil == err) {
		//copy(resp.TxHash, trx.Hash().Bytes())
		resp.TxHash = trx.Hash().Bytes()		
		resp.Errcode = 0
	} else {
		resp.Errcode = 100
	}

	return nil
}


func (a *ApiService) QueryTx(ctx context.Context, req *api.QueryTxRequest, resp *api.QueryTxResponse) error {
	msgReq := &message.QueryTrxReq{
		TxHash: common.HexToHash(req.TxHash),
	}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Tx = nil
		resp.Errcode = 0
		return nil
	}

	response := res.(*message.QueryTrxResp)
	if response.Tx == nil {
		resp.Errcode = 2
		resp.Msg = "Transaction not Found"
		return nil
	}

	resp.Tx = response.Tx
	resp.Errcode = 0
	return nil
}

func (a *ApiService) QueryBlock(ctx context.Context, req *api.QueryBlockRequest, resp *api.QueryBlockResponse) error {
	msgReq := &message.QueryBlockReq{
		BlockHash: common.HexToHash(req.BlockHash),
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

	resp.BlockHash = response.Block.Hash().ToHexString()
	resp.BlockNumber = response.Block.GetNumber()
	resp.Errcode = 0
	return nil
}

func (h *ApiService) QueryChainInfo(ctx context.Context, in *api.QueryChainInfoRequest, out *api.QueryChainInfoResponse) error {
	return nil
}

func (h *ApiService) QueryAccount(ctx context.Context, in *api.QueryAccountRequest, out *api.QueryAccountResponse) error {
	return nil
}