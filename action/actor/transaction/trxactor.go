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
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/action/message"
	"github.com/bottos-project/core/transaction"
)

var TrxActorPid *actor.PID

var trxPool *transaction.TrxPool

type TrxActor struct {
	props *actor.Props
}

func ContructTrxActor() *TrxActor {
	return &TrxActor{}
}

func NewTrxActor() *actor.PID {

	props := actor.FromProducer(func() actor.Actor { return ContructTrxActor() })

	var err error
	TrxActorPid, err = actor.SpawnNamed(props, "TrxActor")

	if err == nil {
		return TrxActorPid
	} else {
		panic(fmt.Errorf("TrxActor SpawnNamed error: ", err))
	}
}

func SetTrxPool(pool *transaction.TrxPool) {
	trxPool = pool
}

func (self *TrxActor) handleSystemMsg(context actor.Context) bool {

	switch msg := context.Message().(type) {

	case *actor.Started:
		log.Printf("TrxActor received started msg", msg)

	case *actor.Stopping:
		log.Printf("TrxActor received stopping msg")

	case *actor.Restart:
		log.Printf("TrxActor received restart msg")

	case *actor.Restarting:
		log.Printf("TrxActor received restarting msg")

	default:
		return false
	}

	return true

}

func (self *TrxActor) Receive(context actor.Context) {

	fmt.Println("trxactor received msg: ", context)

	if self.handleSystemMsg(context) {
		return
	}

	switch msg := context.Message().(type) {

	// case *types.Transaction:
	// 	fmt.Println("transaction action is ", msg.Method.Name)
	// 	context.Respond("trx rsp from trx actor")

	case *message.PushTrxReq:

		fmt.Println("==========")
		fmt.Println(">>>>>>>>>>trx actor Rcv trx, sendType: ", msg.TrxSender, "<<<<<<<<<<<")
		fmt.Println("==========")

		trxPool.HandlePushTransactionReq(context, msg.TrxSender, msg.Trx)

	case *message.GetAllPendingTrxReq:

		fmt.Println("trx actor Rcv get all trx req")

		trxPool.GetAllPendingTransactions(context)

	case *message.RemovePendingTrxsReq:

		fmt.Println("trx actor Rcv remove trxs req")

		trxPool.RemoveTransactions(msg.Trxs)

	default:
		//fmt.Println("trx actor: Unknown msg ", msg, "type", reflect.TypeOf(msg))
		fmt.Println("trx actor: Unknown msg")

	}
}
