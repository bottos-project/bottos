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
 * @Author: may luo
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package produceractor

import (
	"fmt"
	"time"
	//	"unsafe"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/action/env"
	//	"github.com/bottos-project/core/action/message"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/producer"
	"github.com/bottos-project/core/role"
)

type ProducerActor struct {
	roleIntf role.RoleInterface
	ins      producer.ReporterRepo
}

func NewProducerActor(env *env.ActorEnv) *actor.PID {

	ins := producer.New(env.Chain, env.RoleIntf)
	props := actor.FromProducer(func() actor.Actor {
		return &ProducerActor{env.RoleIntf, ins}
	})

	pid, err := actor.SpawnNamed(props, "ProducerActor")

	if err != nil {
		return nil
	}

	return pid
}

func (p *ProducerActor) handleSystemMsg(context actor.Context) {
	switch msg := context.Message().(type) {

	case *actor.Started:
		fmt.Printf("ProducerActor received started msg", msg)
		context.SetReceiveTimeout(500 * time.Millisecond)

	case *actor.ReceiveTimeout:
		fmt.Println("\n\n\n\n ")
		p.working()
		context.SetReceiveTimeout(500 * time.Millisecond)

	case *actor.Stopping:
		fmt.Printf("ProducerActor received stopping msg")

	case *actor.Restart:
		fmt.Printf("ProducerActor received restart msg")

	case *actor.Restarting:
		fmt.Printf("ProducerActor received restarting msg")
	}

}

func (p *ProducerActor) Receive(context actor.Context) {

	p.handleSystemMsg(context)
}
func (p *ProducerActor) working() {
	fmt.Println("begin to producer block ")

	if p.ins.IsReady() {
		start := common.MeasureStart()
		trxs := GetAllPendingTrx()
		if len(trxs) == 0 {
			//fmt.Println("trxs is null,continue produce block")
		}
		block := &types.Block{}
		pendingBlockSize := uint32(10) //unsafe.Sizeof(block)
		coreStat, err := p.roleIntf.GetCoreState()
		if err != nil {
			fmt.Println("GetGlobalPropertyRole failed")
			return
		}
		var pendingTrx = []*types.Transaction{}
		var pendingBlockTrx = []*types.Transaction{}
		for _, trx := range trxs {
			dtag := new(types.Transaction)
			dtag = trx
			if uint64(common.Elapsed(start)) > config.DEFAULT_BLOCK_TIME_LIMIT ||
				pendingBlockSize > coreStat.Config.MaxBlockSize {
				pendingTrx = append(pendingTrx, dtag)
				fmt.Println("Warning reach max size")
				continue
			}

			//p.myDB.StartUndoSession()
			pass, _ := VerifyTransactions(trx)
			if pass == false {
				fmt.Println("ApplyTransaction failed")
				continue
			}
			pendingBlockSize += uint32(20) //unsafe.Sizeof(trx)

			if pendingBlockSize > coreStat.Config.MaxBlockSize {
				fmt.Println("Warning pending block size reach MaxBlockSize")
				pendingTrx = append(pendingTrx, dtag)
				continue
			}
			//	p.myDB.Commit()

			pendingBlockTrx = append(pendingBlockTrx, dtag)
		}
		//fmt.Println("start package block")
		block = p.ins.Woker(trxs)
		if block != nil {
			fmt.Printf("Apply block: hash: %x, delegate: %s, number:%v, trxn:%v\n", block.Hash(), block.Header.Delegate, block.GetNumber(), len(block.Transactions))

			ApplyBlock(block)
			//TODO brocast
		}
	}
}
