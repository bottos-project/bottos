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
 * file description:  transaction actor
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package trxactor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/common/types"
)

var TrxActorPid *actor.PID

type TrxActor struct {
	props *actor.Props
}

func ContructTrxActor() *TrxActor {
	return &TrxActor{}
}

func NewTrxActor() *actor.PID {

	props := actor.FromProducer(func() actor.Actor { return ContructTrxActor() })

	TrxActorPid, err := actor.SpawnNamed(props, "TrxActor")

	if err == nil {
		return TrxActorPid
	} else {
		panic(fmt.Errorf("TrxActor SpawnNamed error: ", err))
	}
}

func (TrxActor *TrxActor) handleSystemMsg(context actor.Context) {

	switch msg := context.Message().(type) {

	case *actor.Started:
		//log.Info("TrxActor received started msg")
		fmt.Println("TrxActor received started msg ", msg)

	case *actor.Stopping:
		//log.Warn("TrxActor received stopping msg")
		fmt.Println("TrxActor received stopping msg ", msg)

	case *actor.Restart:
		//log.Warn("TrxActor received restart msg")
		fmt.Println("TrxActor received restart msg ", msg)

	case *actor.Restarting:
		//log.Warn("TrxActor received restarting msg")
		fmt.Println("TrxActor received restarting msg ", msg)
	}

}

func (TrxActor *TrxActor) Receive(context actor.Context) {

	fmt.Println("trxactor received msg: ", context)

	TrxActor.handleSystemMsg(context)

	switch msg := context.Message().(type) {

	case *types.Transaction:
		fmt.Println("transaction action is ", msg.Action)
		context.Respond("trx rsp from trx actor")

		//default:
		//fmt.Println("trx actor receive default msg ", msg)

	}
}
