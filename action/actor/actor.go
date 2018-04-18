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
)

var apiActorPid *actor.PID
var netActorPid *actor.PID
var trxActorPid *actor.PID
var chainActorPid *actor.PID
var producerActorPid *actor.PID

func InitActors() {

	fmt.Println("InitActors")

	apiActorPid = apiactor.NewApiActor()

	netActorPid = netactor.NewNetActor()

	trxActorPid = trxactor.NewTrxActor()

	chainActorPid = chainactor.NewChainActor()

	producerActorPid = produceractor.NewProducerActor()

	RegisterActorMsgTbl()

}

func RegisterActorMsgTbl() {

	fmt.Println("RegisterActorMsgTbl")

	apiactor.SetTrxActorPid(trxActorPid) // api --> trx

	trxactor.SetApiActorPid(apiActorPid) // trx --> api

}
