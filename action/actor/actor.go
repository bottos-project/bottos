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
	apiactor "github.com/bottos-project/bottos/action/actor/api"
	chainactor "github.com/bottos-project/bottos/action/actor/chain"
	netactor "github.com/bottos-project/bottos/action/actor/net"
	produceractor "github.com/bottos-project/bottos/action/actor/producer"
	trxactor "github.com/bottos-project/bottos/action/actor/transaction"

	"github.com/bottos-project/bottos/action/env"
)

var apiActorPid *actor.PID
var netActorPid *actor.PID
var trxActorPid *actor.PID
var chainActorPid *actor.PID
//MultiActor actor group
type MultiActor struct {
	apiActorPid      *actor.PID
	netActorPid      *actor.PID
	trxActorPid      *actor.PID
	chainActorPid    *actor.PID
	producerActorPid *actor.PID
}

func (m *MultiActor) GetTrxActor() *actor.PID {
	return m.trxActorPid
}
//GetNetActor get net actor PID
func (m *MultiActor) GetNetActor() *actor.PID {
	return m.netActorPid
}

//InitActors init all actor
func InitActors(env *env.ActorEnv) *MultiActor {

	mActor := &MultiActor{
		apiactor.NewApiActor(),
		netactor.NewNetActor(env),
		trxactor.NewTrxActor(),
		chainactor.NewChainActor(env),
		produceractor.NewProducerActor(env),
	}
	registerActorMsgTbl(mActor)
	return mActor
}

func registerActorMsgTbl(m *MultiActor) {

	fmt.Println("RegisterActorMsgTbl")

	apiactor.SetTrxActorPid(m.trxActorPid)          // api --> trx
	apiactor.SetChainActorPid(m.chainActorPid)
	trxactor.SetApiActorPid(m.apiActorPid)          // trx --> api
	produceractor.SetChainActorPid(m.chainActorPid) // producer --> chain
	produceractor.SetTrxActorPid(m.trxActorPid)     // producer --> trx
	produceractor.SetNetActorPid(m.netActorPid)     // producer --> chain
	chainactor.SetTrxActorPid(m.trxActorPid)        //chain --> trx

	netactor.SetTrxActorPid(m.trxActorPid)          //p2p --> trx
	netactor.SetChainActorPid(m.chainActorPid)      //p2p --> chain
}

//GetTrxActorPID get trx actor pid
func (m *MultiActor) GetTrxActorPID() *actor.PID {
	return m.trxActorPid
}

//ActorsStop stop all actor
func (m *MultiActor) ActorsStop() {
	m.chainActorPid.Stop()
	m.producerActorPid.Stop()
	m.apiActorPid.Stop()
	m.netActorPid.Stop()
	m.trxActorPid.Stop()

}
