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
	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/db/platform/codedb"
	"github.com/bottos-project/bottos/producer"
	"github.com/bottos-project/bottos/role"
)

// ProducerActor is to define actor for producer
type ProducerActor struct {
	roleIntf role.RoleInterface
	ins      producer.ReporterRepo
	db               *db.DBService
	pendingTxSession *codedb.UndoSession
}

// NewProducerActor is to create actor for producer
func NewProducerActor(env *env.ActorEnv) *actor.PID {

	ins := producer.New(env.Chain, env.RoleIntf, env.Protocol)
	props := actor.FromProducer(func() actor.Actor {
		return &ProducerActor{
			env.RoleIntf,
			ins,
			env.Db,
			env.PendingTxSession,
		}
	})

	pid, err := actor.SpawnNamed(props, "ProducerActor")

	if err != nil {
		return nil
	}

	return pid
}

func (p *ProducerActor) handleSystemMsg(context actor.Context) bool {
	switch msg := context.Message().(type) {

	case *actor.Started:
		log.Infof("ProducerActor received started msg: %s", msg)

		context.SetReceiveTimeout(time.Duration(config.PRODUCER_TIME_OUT) * time.Millisecond)

	case *actor.ReceiveTimeout:
		elapse := p.working()
		context.SetReceiveTimeout(time.Duration(elapse) * time.Millisecond)

	case *actor.Stopping:
		log.Info("ProducerActor received stopping msg")

	case *actor.Restart:
		log.Info("ProducerActor received restart msg")

	case *actor.Restarting:
		log.Info("ProducerActor received restarting msg")

	case *actor.Stop:
		log.Info("ProducerActor received Stop msg")

	case *actor.Stopped:
		log.Info("ProducerActor received Stopped msg")

	default:
		return false
	}

	return true
}

// Receive is to receive and handle message
func (p *ProducerActor) Receive(context actor.Context) {

	if p.handleSystemMsg(context) {
		return
	}

	log.Error("ProducerActor received Unknown msg")
}
func (p *ProducerActor) working() uint32 {

	if p.ins.IsReady() {
		p.db.Lock()
		defer p.db.UnLock()
		start := common.MeasureStart()
		log.Infof("begin to producer block ")
		p.pendingTxSession = p.db.GetSession()
		if p.pendingTxSession != nil {
			log.Infof("p.pendingTxSession need to reset ")
			p.db.ResetSession()
		}
		log.Infof("begin session......... ")
		p.pendingTxSession = p.db.BeginUndo(config.PRIMARY_TRX_SESSION)

		trxs := GetAllPendingTrx()
		log.Info("get trx times", common.Elapsed(start))
		block := &types.Block{}
		pendingBlockSize := uint32(unsafe.Sizeof(block))
		coreStat, err := p.roleIntf.GetCoreState()
		if err != nil {
			log.Info("GetGlobalPropertyRole failed")
			p.db.ResetSession()
			return config.PRODUCER_TIME_OUT
		}
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
			data, _ := bpl.Marshal(trx)
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
