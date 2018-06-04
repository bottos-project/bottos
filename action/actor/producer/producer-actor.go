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
	"unsafe"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/producer"
	"github.com/bottos-project/bottos/role"
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

	if p.ins.IsReady() {
		start := common.MeasureStart()
		trxs := GetAllPendingTrx()
		block := &types.Block{}
		pendingBlockSize := uint32(unsafe.Sizeof(block))
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
			pass, _ := VerifyTransactions(trx)
			if pass == false {
				fmt.Println("ApplyTransaction failed")
				continue
			}
			pendingBlockSize += uint32(unsafe.Sizeof(trx))

			if pendingBlockSize > coreStat.Config.MaxBlockSize {
				fmt.Println("Warning pending block size reach MaxBlockSize")
				pendingTrx = append(pendingTrx, dtag)
				continue
			}
			pendingBlockTrx = append(pendingBlockTrx, dtag)
		}
		block = p.ins.Woker(trxs)
		if block != nil {
			fmt.Printf("Generate block: hash: %x, delegate: %s, number:%v, trxn:%v,blockTime:%s\n", block.Hash(), block.Header.Delegate, block.GetNumber(), len(block.Transactions), time.Unix(int64(block.Header.Timestamp), 0))

			ApplyBlock(block)
			fmt.Printf("Broadcast block: block num:%v, trxn:%v, delegate: %s, hash: %x\n", block.GetNumber(), len(block.Transactions), block.Header.Delegate, block.Hash())
		}
	}
}
