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
	"time"
	"unsafe"

	log "github.com/cihub/seelog"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/producer"
	"github.com/bottos-project/bottos/role"
	proto "github.com/golang/protobuf/proto"
)

// ProducerActor is to define actor for producer
type ProducerActor struct {
	roleIntf role.RoleInterface
	ins      producer.ReporterRepo
}

// NewProducerActor is to create actor for producer
func NewProducerActor(env *env.ActorEnv) *actor.PID {

	ins := producer.New(env.Chain, env.RoleIntf, env.Protocol)
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
		log.Infof("ProducerActor received started msg: %s", msg)
		context.SetReceiveTimeout(500 * time.Millisecond)

	case *actor.ReceiveTimeout:
		p.working()
		context.SetReceiveTimeout(500 * time.Millisecond)

	case *actor.Stopping:
		log.Info("ProducerActor received stopping msg")

	case *actor.Restart:
		log.Info("ProducerActor received restart msg")

	case *actor.Restarting:
		log.Info("ProducerActor received restarting msg")
	}

}

// Receive is to receive and handle message
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
			log.Info("GetGlobalPropertyRole failed")
			return
		}
		var pendingTrx = []*types.Transaction{}
		var pendingBlockTrx = []*types.Transaction{}
		var removeTrx = []*types.Transaction{}
		for _, trx := range trxs {
			dtag := new(types.Transaction)
			dtag = trx
			if uint64(common.Elapsed(start)) > config.DEFAULT_BLOCK_TIME_LIMIT ||
				pendingBlockSize > coreStat.Config.MaxBlockSize {
				pendingTrx = append(pendingTrx, dtag)
				log.Info("Warning reach max size")
				continue
			}
			pass, _ := verifyTransactions(trx)
			if pass == false {
				log.Info("ApplyTransaction failed")
				removeTrx = append(removeTrx, trx)
				continue
			}
			data, _ := proto.Marshal(trx)
			pendingBlockSize += uint32(unsafe.Sizeof(data))

			if pendingBlockSize > coreStat.Config.MaxBlockSize {
				log.Info("Warning pending block size reach MaxBlockSize")
				pendingTrx = append(pendingTrx, dtag)
				continue
			}
			pendingBlockTrx = append(pendingBlockTrx, dtag)
		}
		removeTransaction(removeTrx)
		block = p.ins.Woker(pendingBlockTrx)
		trxs = nil
		if block != nil {
			log.Infof("Generate block: hash: %x, delegate: %s, number:%v, trxn:%v,blockTime:%s\n", block.Hash(), block.Header.Delegate, block.GetNumber(), len(block.Transactions), time.Unix(int64(block.Header.Timestamp), 0))
			ApplyBlock(block)
		}
	}
}
