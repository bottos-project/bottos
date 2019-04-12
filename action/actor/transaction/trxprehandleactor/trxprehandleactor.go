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
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/action/message"
	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/context"
	"github.com/bottos-project/bottos/transaction"
	"github.com/bottos-project/bottos/config"
)

//TrxPreHandleActorPid trx actor pid
var TrxPreHandleActorPid *actor.PID

var trxActorPid *actor.PID


const maxConcurrency = 100

var trxPool *transaction.TrxPool

var actorEnv *env.ActorEnv

var protocolInterface context.ProtocolInterface

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
		log.Error("TrxPreHandleActor received started msg")
	case *actor.Stopping:
		log.Error("TrxPreHandleActor received stopping msg")
	case *actor.Restart:
		log.Error("TrxPreHandleActor received restart msg")
	case *actor.Restarting:
		log.Error("TrxPreHandleActor received restarting msg")
	case *actor.Stop:
		log.Error("TrxPreHandleActor received Stop msg")
	case *actor.Stopped:
		log.Error("TrxPreHandleActor received Stopped msg")
		
	default:
		return false
	}

	return true
}

func preHandleCommon (trx *types.Transaction) (bool, bottosErr.ErrCode) {
	if checkResult, err := trxPool.CheckTransactionBaseCondition(trx); true != checkResult {
		return false, err
	}

	if false == actorEnv.Protocol.GetBlockSyncState() {
		log.Errorf("TRX rcv trx when block is syncing, trx %x", trx.Hash())

		return false, bottosErr.ErrTrxBlockSyncingError
	}

	sender, err := actorEnv.RoleIntf.GetAccount(trx.Sender)
	if nil != err {
		return false, bottosErr.ErrTrxAccountError
	}

	if !trx.VerifySignature(sender.PublicKey) {
		return false, bottosErr.ErrTrxSignError
	}

	return true, bottosErr.ErrNoError
}
func initP2PTrxMsg(msg *message.PushTrxReq) (msgp *message.PushTrxForP2PReq) {
	//set trx TTL
	var TTL uint16
	switch  actorEnv.RoleIntf.IsMyselfDelegate() {
	case true:
		TTL = config.TRX_IN_TTL
	case false:
		TTL = config.TRX_OUT_TTL
	}
	var p2pTrx types.P2PTransaction
	p2pTrx.Transaction = msg.Trx
	p2pTrx.TTL = TTL
	msgp = &message.PushTrxForP2PReq{P2PTrx: &p2pTrx}
	return msgp
}

func preHandlePushTrxReq(msg *message.PushTrxReq, ctx actor.Context) {

	preHandleResult, err := preHandleCommon(msg.Trx)
	
	if !preHandleResult {			
		log.Errorf("TRX pre handle trx from front failed, trx %x", msg.Trx.Hash())
		ctx.Respond(err)
	} else {
		msgP2P := initP2PTrxMsg(msg)
		trxActorPid.Tell(msgP2P)
		ctx.Respond(bottosErr.ErrNoError)
	}
}

func preHandleReceiveTrx(msg *message.ReceiveTrx, ctx actor.Context) {

	preHandleResult, _ := preHandleCommon(msg.P2PTrx.Transaction)
	
	if preHandleResult {
		trxActorPid.Tell(msg)
	} else {
		if actorEnv.RoleIntf.IsMyselfDelegate() == true{
			log.Info("TRX pre handle trx from producer node failed, trx %x", msg.P2PTrx.Transaction.Hash())
		}else{
			log.Errorf("TRX pre handle trx from service node failed, trx %x", msg.P2PTrx.Transaction.Hash())
		}

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

		log.Infof("rcv trx %x in ReceiveTrx\n", msg.P2PTrx.Transaction.Hash())

		preHandleReceiveTrx(msg, ctx)		

	default:
		log.Errorf("trx pool actor: Unknown msg ", msg)
	}

}

//NewTrxPreHandleActor spawn a named actor
func NewTrxPreHandleActor(env *env.ActorEnv) *actor.PID {

	actorEnv = env

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

