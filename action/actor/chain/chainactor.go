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
 * file description:  chain actor
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package chainactor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/action/message"
)

var ChainActorPid *actor.PID
var actorEnv *env.ActorEnv
var trxactorPid *actor.PID

type ChainActor struct {
	props *actor.Props
}

func ContructChainActor() *ChainActor {
	return &ChainActor{}
}

func SetTrxActorPid(tpid *actor.PID) {
	trxactorPid = tpid
}

func NewChainActor(env *env.ActorEnv) *actor.PID {
	var err error

	props := actor.FromProducer(func() actor.Actor { return ContructChainActor() })

	ChainActorPid, err = actor.SpawnNamed(props, "ChainActor")
	actorEnv = env

	if err == nil {
		return ChainActorPid
	} else {
		panic(fmt.Errorf("ChainActor SpawnNamed error: ", err))
	}
}

func (self *ChainActor) handleSystemMsg(context actor.Context) {

	switch msg := context.Message().(type) {

	case *actor.Started:
		fmt.Printf("BlockActor received started msg", msg)

	case *actor.Stopping:
		fmt.Printf("BlockActor received stopping msg")

	case *actor.Restart:
		fmt.Printf("BlockActor received restart msg")

	case *actor.Restarting:
		fmt.Printf("BlockActor received restarting msg")
	}

}

func (self *ChainActor) Receive(context actor.Context) {

	self.handleSystemMsg(context)

	switch msg := context.Message().(type) {
	case *message.InsertBlockReq:
		self.HandleBlockMessage(context, msg)
	case *message.QueryTrxReq:
		self.HandleQueryTrxReq(context, msg)
	case *message.QueryBlockReq:
		self.HandleQueryBlockReq(context, msg)
	case *message.QueryChainInfoReq:
		self.HandleQueryChainInfoReq(context, msg)
	}
}

func (self *ChainActor) HandleBlockMessage(ctx actor.Context, req *message.InsertBlockReq) {
	err := actorEnv.Chain.InsertBlock(req.Block)
	if ctx.Sender() != nil {
		resp := &message.InsertBlockRsp{
			Hash:  req.Block.Hash(),
			Error: err,
		}
		ctx.Sender().Request(resp, ctx.Self())
	}
	if err == nil {
		req := &message.RemovePendingTrxsReq{Trxs: req.Block.Transactions}
		trxactorPid.Tell(req)
	}
}

func (self *ChainActor) HandleQueryTrxReq(ctx actor.Context, req *message.QueryTrxReq) {
	tx := actorEnv.TxStore.GetTransaction(req.TrxHash)
	if ctx.Sender() != nil {
		resp := &message.QueryTrxResp{}
		if tx == nil {
			resp.Error = fmt.Errorf("Transaction not found")
		} else {
			resp.Trx = tx
		}
		ctx.Sender().Request(resp, ctx.Self())
	}
}

func (self *ChainActor) HandleQueryBlockReq(ctx actor.Context, req *message.QueryBlockReq) {
	block := actorEnv.Chain.GetBlockByHash(req.BlockHash)
	if block == nil {
		block = actorEnv.Chain.GetBlockByNumber(req.BlockNumber)
	}
	if ctx.Sender() != nil {
		resp := &message.QueryBlockResp{}
		if block == nil {
			resp.Error = fmt.Errorf("Block not found")
		} else {
			resp.Block = block
		}
		ctx.Sender().Request(resp, ctx.Self())
	}
}

func (self *ChainActor) HandleQueryChainInfoReq(ctx actor.Context, req *message.QueryChainInfoReq) {
	if ctx.Sender() != nil {
		resp := &message.QueryChainInfoResp{}
		resp.HeadBlockNum = actorEnv.Chain.HeadBlockNum()
		resp.HeadBlockHash = actorEnv.Chain.HeadBlockHash()
		resp.HeadBlockTime = actorEnv.Chain.HeadBlockTime()
		resp.HeadBlockDelegate = actorEnv.Chain.HeadBlockDelegate()
		resp.LastConsensusBlockNum = actorEnv.Chain.LastConsensusBlockNum()
		ctx.Sender().Request(resp, ctx.Self())
	}
}
