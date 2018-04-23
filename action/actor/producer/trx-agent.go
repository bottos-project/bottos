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
	"fmt"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/action/message"
	//	"github.com/bottos-project/core/common/types"
)

var trxActorPid *actor.PID

func SetTrxActorPid(tpid *actor.PID) {
	trxActorPid = tpid
}

func GetAllPendingTrx() {
	getTrxsReq := &message.GetAllPendingTrxReq{}
	fmt.Println("trxActorPid", trxActorPid)
	getTrxsResult, getTrxsErr := trxActorPid.RequestFuture(getTrxsReq, 500*time.Millisecond).Result()

	if nil == getTrxsErr {
		fmt.Println("get all trx req exec result:")
		fmt.Println("rusult is =======", getTrxsResult)
		fmt.Println("error  is =======", getTrxsErr)
	} else {
		fmt.Println("get all trx req exec error")
	}
	//var trxs = []*types.Transaction{}
	//for i=1; i< len(getTrxsResult)
	//return getTrxsResult
}
