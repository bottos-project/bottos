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
 * file description:  api agent
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package trxactor

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

var apiactorPid *actor.PID
var netActor *actor.PID

// SetApiActorPid is to save api actor
func SetApiActorPid(apid *actor.PID) {
	apiactorPid = apid
}

// SetNetActorPid is to save net actor
func SetNetActorPid(pid *actor.PID) {
	netActor = pid
}

func sendTrxRsp(trxRsp uint64, pid *actor.PID) {

	pid.Tell("pushTrxRsp")
	/*
		pushTrxReq := &types.Transaction{
			RefBlockNum: 11,
			Sender:      22,
		}

			trxactorPid.Tell(pushTrxReq)

			f := trxactorPid.RequestFuture(pushTrxReq, 5000*time.Millisecond)
			es, err := f.Result() // waits for pid to reply

			fmt.Println("this is es err", es, err)
	*/
	//result, _ := trxactorPid.RequestFuture(pushTrxReq, 500*time.Millisecond).Result() // await result

	//fmt.Println(result)

}
