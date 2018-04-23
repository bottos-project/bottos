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
 * file description:  actor entry
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package actionactor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	apiactor "github.com/bottos-project/core/action/actor/api"
	chainactor "github.com/bottos-project/core/action/actor/chain"
	netactor "github.com/bottos-project/core/action/actor/net"
	produceractor "github.com/bottos-project/core/action/actor/producer"
	trxactor "github.com/bottos-project/core/action/actor/transaction"

	"github.com/bottos-project/core/action/env"
)

var apiActorPid *actor.PID
var netActorPid *actor.PID
var trxActorPid *actor.PID
var chainActorPid *actor.PID

type MultiActor struct {
	apiActorPid      *actor.PID
	netActorPid      *actor.PID
	trxActorPid      *actor.PID
	chainActorPid    *actor.PID
	producerActorPid *actor.PID
}

func InitActors(env *env.ActorEnv) *MultiActor {

	mActor := &MultiActor{
		apiactor.NewApiActor(),
		netactor.NewNetActor(),
		trxactor.NewTrxActor(),
		chainactor.NewChainActor(env),
		produceractor.NewProducerActor(env),
	}
	registerActorMsgTbl(mActor)
	return mActor
}

func registerActorMsgTbl(m *MultiActor) {

	fmt.Println("RegisterActorMsgTbl")

	apiactor.SetTrxActorPid(m.trxActorPid) // api --> trx

	trxactor.SetApiActorPid(m.apiActorPid) // trx --> api

	produceractor.SetChainActorPid(m.chainActorPid) // producer --> chain

	chainactor.SetTrxActorPid(m.trxActorPid) //chain --> trx

}

func (m *MultiActor) GetTrxActorPID() *actor.PID {
	return m.trxActorPid
}
