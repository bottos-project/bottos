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
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/action/env"
	"github.com/bottos-project/core/action/message"
	"github.com/bottos-project/core/chain"
)

var ChainActorPid *actor.PID

type ChainActor struct {
	props *actor.Props
}

func ContructChainActor() *ChainActor {
	return &ChainActor{}
}

func NewChainActor(env *env.ActorEnv) *actor.PID {
	var err error

	props := actor.FromProducer(func() actor.Actor { return ContructChainActor() })

	ChainActorPid, err = actor.SpawnNamed(props, "ChainActor")

	if err == nil {
		return ChainActorPid
	} else {
		panic(fmt.Errorf("ChainActor SpawnNamed error: ", err))
	}
}

func (self *ChainActor) handleSystemMsg(context actor.Context) {

	switch msg := context.Message().(type) {

	case *actor.Started:
		log.Printf("BlockActor received started msg", msg)

	case *actor.Stopping:
		log.Printf("BlockActor received stopping msg")

	case *actor.Restart:
		log.Printf("BlockActor received restart msg")

	case *actor.Restarting:
		log.Printf("BlockActor received restarting msg")
	}

}

func (self *ChainActor) Receive(context actor.Context) {

	self.handleSystemMsg(context)

	switch msg := context.Message().(type) {
	case *message.InsertBlockReq:
		self.HandleBlockMessage(context, msg)
	}
}

func (self *ChainActor) HandleBlockMessage(ctx actor.Context, req *message.InsertBlockReq) {
	err := chain.GetChain().InsertBlock(req.Block)
	if ctx.Sender() != nil {
		resp := &message.InsertBlockRsp{
			Hash:  req.Block.Hash(),
			Error: err,
		}
		ctx.Sender().Request(resp, ctx.Self())
	}
}
