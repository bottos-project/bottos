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
 * file description:  producer actor
 * @Author:
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package produceractor

import (
	"time"

	log "github.com/cihub/seelog"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/transaction"
)

var trxPoolActorPid *actor.PID

// SetTrxPoolActorPid is to set transaction actor PID for use
func SetTrxPoolActorPid(tpid *actor.PID) {
	trxPoolActorPid = tpid
}

// GetAllPendingTrx is to retrieve the pending transactions
func GetAllPendingTrx() []*types.Transaction {
	getTrxsReq := &message.GetAllPendingTrxReq{}
	getTrxsResult, getTrxsErr := trxPoolActorPid.RequestFuture(getTrxsReq, 900*time.Millisecond).Result()

	if nil == getTrxsErr {
	} else {
		log.Errorf("PRODUCER get all trx req exec error %v", getTrxsErr)
	}

	if nil == getTrxsResult {
		return nil
	}
	mesg := getTrxsResult.(*message.GetAllPendingTrxRsp)
	log.Debugf("PRODUCER pending transaction number %v", len(mesg.Trxs))
	var trxs = []*types.Transaction{}
	for i := 0; i < len(mesg.Trxs); i++ {
		dbtag := new(types.Transaction)
		dbtag = mesg.Trxs[i]

		trxs = append(trxs, dbtag)
	}

	return trxs
}

// VerifyTransactions is to verify local and received transactons
func verifyTransactions(trx *types.Transaction, version uint32) (bool, error, *types.ResourceReceipt) {
	if trx.Version != version {
		log.Errorf("PRODUCER verify trx failed, hash:%x, trx.Version:%v, my version:%v", trx.Hash(), trx.Version, version)
		return false, nil, nil
	}
	trxApply := transaction.NewTrxApplyService()
	pass, _, _, resourceReceipt, _ := trxApply.ExecuteTransaction(trx,true)

	log.Infof("PRODUCER verify trx result %v,%x,%v", pass, trx.Hash(), resourceReceipt)
	return pass, nil, resourceReceipt
}

func removeTransaction(trxs []*types.Transaction) {
	removeTrxs := &message.RemovePendingTrxsReq{Trxs: trxs}
	trxPoolActorPid.Tell(removeTrxs)
	log.Infof("PRODUCER remove transaction %v", len(trxs))
}
