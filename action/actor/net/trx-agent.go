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
 * @Author: Stewart Li
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package netactor

//import "github.com/bottos-project/bottos/transaction"
import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/common/types"
)

var trxActorPid *actor.PID

func SetTrxActorPid(tpid *actor.PID) {
	//trxActorPid = tpid
	p2p.SetTrxActor(tpid)
}

/*
func SetChainActorPid(cpid *actor.PID) {
	p2p.SetChainActor(cpid)
}
*/

//Get Trx from TxPool , and the trx will be boardcasted by p2p component
//func ReceiveNewTrx() []*types.Transaction {
func GetAllPendingTrx() []*types.Transaction {
	/*
	getTrxsReq := &message.GetAllPendingTrxReq{}
	getTrxsResult, getTrxsErr := trxActorPid.RequestFuture(getTrxsReq, 500*time.Millisecond).Result()
	if nil == getTrxsErr {
	} else {
		fmt.Println("get all trx req exec error") //TODO
	}

	mesg := getTrxsResult.(*message.GetAllPendingTrxRsp)
	fmt.Println("pending transaction number ", len(mesg.Trxs))
	var trxs = []*types.Transaction{}
	for i := 0; i < len(mesg.Trxs); i++ {
		dbtag := new(types.Transaction)
		dbtag = mesg.Trxs[i]

		trxs = append(trxs, dbtag)
	}

	return trxs
	*/
	return nil
}

//Send new Trx from other peers
func SendNewTrx(trx *types.Transaction) (bool, error) {
	/*
	if trx == nil {
		return false , errors.New("*ERROR* Failed to send the data from netactor !!!")
	}
	trxActorPid.Tell(trx)
	*/
	return true , nil
}


