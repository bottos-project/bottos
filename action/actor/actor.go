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
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/actor/api"
	"github.com/bottos-project/bottos/action/actor/chain"
	consensusactor "github.com/bottos-project/bottos/action/actor/consensus"
	"github.com/bottos-project/bottos/action/actor/net"
	"github.com/bottos-project/bottos/action/actor/producer"
	"github.com/bottos-project/bottos/action/actor/transaction/trxhandleactor"
	"github.com/bottos-project/bottos/action/actor/transaction/trxpoolactor"
	"github.com/bottos-project/bottos/action/actor/transaction/trxprehandleactor"
	//	bfttestactor "github.com/bottos-project/bottos/action/actor/bfttest"
	log "github.com/cihub/seelog"

	"github.com/bottos-project/bottos/action/env"
	restactor"github.com/bottos-project/bottos/restful/handler"
	walletrestactor "github.com/bottos-project/bottos/restful/wallet"
)

var apiActorPid *actor.PID
var netActorPid *actor.PID
var trxActorPid *actor.PID
var chainActorPid *actor.PID

//MultiActor actor group
type MultiActor struct {
	apiActorPid       *actor.PID
	netActorPid       *actor.PID
	trxActorPid       *actor.PID
	trxPoolActorPid   *actor.PID
	trxPreHandleActor *actor.PID
	chainActorPid     *actor.PID
	producerActorPid  *actor.PID
	consensusActorPid *actor.PID
}

//GetTrxActor get net actor PID
func (m *MultiActor) GetTrxActor() *actor.PID {
	return m.trxActorPid
}

//GetTrxPreHandleActor get net actor PID
func (m *MultiActor) GetTrxPreHandleActor() *actor.PID {
	return m.trxPreHandleActor
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
		trxpoolactor.NewTrxPoolActor(),
		trxprehandleactor.NewTrxPreHandleActor(env),
		chainactor.NewChainActor(env),
		produceractor.NewProducerActor(env),
		consensusactor.NewConsensusActor(env),
	}
	registerActorMsgTbl(mActor)
	return mActor
}

func registerActorMsgTbl(m *MultiActor) {

	log.Info("RegisterActorMsgTbl")
	apiactor.SetTrxPreHandleActorPid(m.trxPreHandleActor) // api --> trx
	restactor.SetTrxPreHandleActorPid(m.trxPreHandleActor) // api --> trx
	walletrestactor.SetTrxPreHandleActorPid(m.trxPreHandleActor) // api --> trx
	apiactor.SetChainActorPid(m.chainActorPid) // api --> chain
	restactor.SetChainActorPid(m.chainActorPid)	// restapi --> chain
	//walletrestactor.SetChainActorPid(m.chainActorPid)            // restapi --> chain
	trxactor.SetApiActorPid(m.apiActorPid)          // trx --> api
	produceractor.SetChainActorPid(m.chainActorPid) // producer --> chain
	produceractor.SetTrxPoolActorPid(m.trxPoolActorPid)     // producer --> trx
	produceractor.SetNetActorPid(m.netActorPid)     // producer --> chain
	produceractor.SetConsensusActorPid(m.consensusActorPid)      // producer --> consensus

	chainactor.SetTrxPoolActorPid(m.trxPoolActorPid)        // chain --> trx
	chainactor.SetNetActorPid(m.netActorPid)        // chain --> net

	netactor.SetTrxPreHandleActorPid(m.trxPreHandleActor)     //p2p --> trx
	netactor.SetChainActorPid(m.chainActorPid) //p2p --> chain
	netactor.SetConsensusActorPid(m.consensusActorPid)    //p2p --> consensus

	consensusactor.SetChainActorPid(m.chainActorPid) //consensus --> chain
	consensusactor.SetNetActorPid(m.netActorPid)     //consensus --> p2p

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
	m.trxPoolActorPid.Stop()
	m.consensusActorPid.Stop()
}
