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
 * file description:  transaction pool actor
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package trxpoolactor

import (
	log "github.com/cihub/seelog"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/transaction"
)

//TrxPoolActorPid trx actor pid
var TrxPoolActorPid *actor.PID

var trxPool *transaction.TrxPool

//TrxPoolActor trx actor props
type TrxPoolActor struct {
	props *actor.Props
}

//ContructTrxPoolActor new a trx actor
func ContructTrxPoolActor() *TrxPoolActor {
	return &TrxPoolActor{}
}

//NewTrxPoolActor spawn a named actor
func NewTrxPoolActor() *actor.PID {

	props := actor.FromProducer(func() actor.Actor { return ContructTrxPoolActor() })

	var err error
	TrxPoolActorPid, err = actor.SpawnNamed(props, "TrxPoolActor")

	if err != nil {
		panic(log.Errorf("TrxPoolActor SpawnNamed error: %v", err))
	} else {
		return TrxPoolActorPid
	}
}

//SetTrxPool set trx pool
func SetTrxPool(pool *transaction.TrxPool) {
	trxPool = pool
}

func handleSystemMsg(context actor.Context) bool {
	switch context.Message().(type) {
	case *actor.Started:
		log.Error("TrxPoolActor received started msg")
	case *actor.Stopping:
		log.Error("TrxPoolActor received stopping msg")
	case *actor.Restart:
		log.Error("TrxPoolActor received restart msg")
	case *actor.Restarting:
		log.Error("TrxPoolActor received restarting msg")
	case *actor.Stop:
		log.Error("TrxPoolActor received Stop msg")
	case *actor.Stopped:
		log.Error("TrxPoolActor received Stopped msg")
	default:
		return false
	}

	return true
}

//Receive process message
var trxcnt uint64 = 0

func (t *TrxPoolActor) Receive(context actor.Context) {

	if handleSystemMsg(context) {
		return
	}

	switch msg := context.Message().(type) {
	// case *message.PushTrxReq:
	// 	trxcnt += 1
	// 	//log.Error("TrxPoolActor received trx: ", trxcnt)
	// 	trxPool.HandleTransactionFromFront(context, msg.Trx)

	// case *message.NotifyTrx:

	// 	trxPool.HandleTransactionFromP2P(context, msg.Trx)

	case *message.GetAllPendingTrxReq:

		trxPool.GetAllPendingTransactions(context)

	case *message.RemovePendingTrxsReq:

		trxPool.RemoveTransactions(msg.Trxs)

	case *message.RemovePendingBlockTrxsReq:

		trxPool.RemoveBlockTransactions(msg.Trxs)

	default:
		log.Errorf("trx pool actor: Unknown msg ", msg)
	}
}
