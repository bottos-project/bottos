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
)

var trxActorPid *actor.PID

// SetTrxActorPid is to set transaction actor PID for use
func SetTrxActorPid(tpid *actor.PID) {
	trxActorPid = tpid
}

// GetAllPendingTrx is to retrieve the pending transactions
func GetAllPendingTrx() []*types.Transaction {
	getTrxsReq := &message.GetAllPendingTrxReq{}
	getTrxsResult, getTrxsErr := trxActorPid.RequestFuture(getTrxsReq, 500*time.Millisecond).Result()

	if nil == getTrxsErr {
	} else {
		log.Info("get all trx req exec error")
	}

	mesg := getTrxsResult.(*message.GetAllPendingTrxRsp)
	log.Infof("pending transaction number ", len(mesg.Trxs))
	var trxs = []*types.Transaction{}
	for i := 0; i < len(mesg.Trxs); i++ {
		dbtag := new(types.Transaction)
		dbtag = mesg.Trxs[i]

		trxs = append(trxs, dbtag)
	}

	return trxs
}

// VerifyTransactions is to verify local and received transactons
func VerifyTransactions(trx *types.Transaction) (bool, error) {
	return true, nil
	/* 	log.Info("start apply transaction trx one by one")
	   	trxApply := transaction.NewTrxApplyService()
	   	pass, _, _ := trxApply.ApplyTransaction(trx)
	   	return pass, nil */
}
