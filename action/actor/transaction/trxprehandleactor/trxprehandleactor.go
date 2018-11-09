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

package trxprehandleactor

import (
	log "github.com/cihub/seelog"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/router"
	"github.com/bottos-project/bottos/action/message"
	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/transaction"
)

//TrxPreHandleActorPid trx actor pid
var TrxPreHandleActorPid *actor.PID

var trxActorPid *actor.PID

const maxConcurrency = 100

var trxPool *transaction.TrxPool

//TrxActor trx actor props
// type TrxActor struct {
// 	props *actor.Props
// }

// //ContructTrxActor new a trx actor
// func ContructTrxActor() *TrxActor {
// 	return &TrxActor{}
// }

func handleSystemMsg(context actor.Context) bool {
	switch context.Message().(type) {
	case *actor.Started:
		log.Info("TrxPoolActor received started msg")
	case *actor.Stopping:
		log.Info("TrxPoolActor received stopping msg")
	case *actor.Restart:
		log.Info("TrxPoolActor received restart msg")
	case *actor.Restarting:
		log.Info("TrxPoolActor received restarting msg")
	default:
		return false
	}

	return true
}

func preHandleCommon(trx *types.Transaction) (bool, bottosErr.ErrCode) {
	if checkResult, err := trxPool.CheckTransactionBaseCondition(trx); true != checkResult {
		return false, err
	}
	if !trxPool.VerifySignature(trx) {
		log.Errorf("trx %v VerifySignature error\n", trx.Hash())
		return false, bottosErr.ErrTrxSignError
	}

	return true, bottosErr.ErrNoError
}

func preHandlePushTrxReq(msg *message.PushTrxReq, ctx actor.Context) {

	preHandleResult, err := preHandleCommon(msg.Trx)

	if !preHandleResult {
		ctx.Respond(err)
	} else {
		trxActorPid.Tell(msg)
		ctx.Respond(bottosErr.ErrNoError)
	}
}

func preHandleReceiveTrx(msg *message.ReceiveTrx, ctx actor.Context) {

	preHandleResult, err := preHandleCommon(msg.Trx)

	if !preHandleResult {
		ctx.Respond(err)
	} else {
		trxActorPid.Tell(msg)
		ctx.Respond(bottosErr.ErrNoError)
	}
}

func doWork(ctx actor.Context) {

	if handleSystemMsg(ctx) {
		return
	}

	switch msg := ctx.Message().(type) {
	case *message.PushTrxReq:

		log.Infof("rcv trx %x in PushTrxReq\n", msg.Trx.Hash())

		preHandlePushTrxReq(msg, ctx)

	case *message.ReceiveTrx:

		log.Infof("rcv trx %x in ReceiveTrx\n", msg.Trx.Hash())

		preHandleReceiveTrx(msg, ctx)

	default:
		log.Info("trx actor: Unknown msg")
	}

}

//NewTrxPreHandleActor spawn a named actor
func NewTrxPreHandleActor() *actor.PID {

	// props := actor.FromProducer(func() actor.Actor { return ContructTrxActor() })

	// var err error
	// TrxPreHandleActorPid, err = actor.SpawnNamed(props, "TrxActor")

	// if err != nil {
	// 	panic(log.Errorf("TrxActor SpawnNamed error: %v", err))
	// } else {
	// 	return TrxPreHandleActorPid
	// }

	TrxPreHandleActorPid := actor.Spawn(router.NewRoundRobinPool(maxConcurrency).WithFunc(doWork))

	return TrxPreHandleActorPid
}

//SetTrxPool set trx pool
func SetTrxPool(pool *transaction.TrxPool) {
	trxPool = pool
}

func SetTrxActor(trxactorPid *actor.PID) {
	trxActorPid = trxactorPid
}

// func handleSystemMsg(context actor.Context) bool {
// 	switch context.Message().(type) {
// 	case *actor.Started:
// 		log.Info("TrxActor received started msg")
// 	case *actor.Stopping:
// 		log.Info("TrxActor received stopping msg")
// 	case *actor.Restart:
// 		log.Info("TrxActor received restart msg")
// 	case *actor.Restarting:
// 		log.Info("TrxActor received restarting msg")
// 	default:
// 		return false
// 	}

// 	return true
// }

//Receive process message
var trxcnt uint64 = 0

// func (t *TrxActor) Receive(context actor.Context) {

// 	if handleSystemMsg(context) {
// 		return
// 	}

// 	switch msg := context.Message().(type) {
// 	case *message.PushTrxReq:
// 		trxcnt += 1
// 		//log.Error("TrxActor received trx: ", trxcnt)
// 		trxPool.HandleTransactionFromFront(context, msg.Trx)

// 	case *message.NotifyTrx:

// 		trxPool.HandleTransactionFromP2P(context, msg.Trx)

// 	// case *message.GetAllPendingTrxReq:

// 	// 	trxPool.GetAllPendingTransactions(context)

// 	// case *message.RemovePendingTrxsReq:

// 	// 	trxPool.RemoveTransactions(msg.Trxs)

// 	default:
// 		log.Info("trx actor: Unknown msg")
// 	}
// }
