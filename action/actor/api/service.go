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
)

type ApiService struct {
	trxActorPid *actor.PID
	chainActorPid *actor.PID
}

func NewApiService() api.CoreApiHandler {
	apiService := &ApiService{}
	return apiService
}

func (a *ApiService) PushTrx(ctx context.Context, trx *types.Transaction, resp *api.PushResponse) error {
	return nil
}

func (a *ApiService) QueryTrx(ctx context.Context, req *api.QueryTrxRequest, resp *api.QueryTrxResponse) error {
	msgReq := message.QueryTrxReq{
		TxHash: common.BytesToHash(req.TxHash),
	}
	res, err := a.chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Trx = nil
		resp.Errcode = 0 // TODO
		return err
	}

	response := res.(*message.QueryTrxResp)
	resp.Trx = response.Tx
	resp.Errcode = 0
	return nil

}

func (a *ApiService) QueryBlock(ctx context.Context, req *api.QueryBlockRequest, resp *api.QueryBlockResponse) error {
	msgReq := message.QueryBlockReq{
		BlockHash: common.BytesToHash(req.BlockHash),
	}
	res, err := a.chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Block = nil
		resp.Errcode = 0 // TODO
		return err
	}

	response := res.(*message.QueryBlockResp)
	resp.Block = response.Block
	resp.Errcode = 0
	return nil
}