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
	"github.com/bottos-project/bottos/version"
)

// ProducerActor is to define actor for producer
type ProducerActor struct {
	roleIntf         role.RoleInterface
	ins              producer.ReporterRepo
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
		log.Error("PRODUCER received started msg ", msg)

		context.SetReceiveTimeout(time.Duration(config.PRODUCER_TIME_OUT) * time.Millisecond)

	case *actor.ReceiveTimeout:
		elapse := p.working()
		context.SetReceiveTimeout(time.Duration(elapse) * time.Millisecond)

	case *actor.Stopping:
		log.Error("PRODUCER received stopping msg")

	case *actor.Restart:
		log.Error("PRODUCER received restart msg")

	case *actor.Restarting:
		log.Error("PRODUCER received restarting msg")

	case *actor.Stop:
		log.Error("PRODUCER received Stop msg")

	case *actor.Stopped:
		log.Error("PRODUCER received Stopped msg")

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

	log.Error("PRODUCER received Unknown msg")
}
func (p *ProducerActor) working() uint32 {

	if p.ins.IsReady() {
		p.db.Lock()
		defer p.db.UnLock()
		start := common.MeasureStart()
		p.pendingTxSession = p.db.GetSession()
		if p.pendingTxSession != nil {
			log.Debug("PRODUCER p.pendingTxSession need to reset")
			p.db.ResetSession()
		}
		p.pendingTxSession = p.db.BeginUndo(config.PRIMARY_TRX_SESSION)

		trxs := GetAllPendingTrx()
		log.Debug("PRODUCER get trx times", common.Elapsed(start))
		pendingTrxlen := len(trxs)
		block := &types.Block{}
		pendingBlockSize := uint32(unsafe.Sizeof(block))
		coreStat, err := p.roleIntf.GetCoreState()
		if err != nil {
			log.Error("PRODUCER GetCoreState failed,begin rollback", err)
			p.db.ResetSession()
			return config.PRODUCER_TIME_OUT
		}
		var pendingBlockTrx = []*types.Transaction{}
		var removeTrx = []*types.Transaction{}
		for _, trx := range trxs {
			dtag := new(types.Transaction)
			dtag = trx
			if uint64(common.Elapsed(start)) > config.DEFAULT_BLOCK_TIME_LIMIT {
				log.Info("Warning producing block is too slow", common.Elapsed(start))
				break
			}
			p.db.BeginUndo(config.SUB_TRX_SESSION)
			applyStart := common.MeasureStart()
			pass, _ := verifyTransactions(trx)
			if pass == false {
				log.Error("PRODUCER verify transactions failed, trx %x", trx.Hash())
				p.db.ResetSubSession()
				removeTrx = append(removeTrx, trx)
				continue
			}
			log.Info("apply start elapse", common.Elapsed(applyStart))
			data, _ := bpl.Marshal(trx)
			pendingBlockSize += uint32(unsafe.Sizeof(data))
			log.Info("pendingBlockSize ", pendingBlockSize)

			if pendingBlockSize > coreStat.Config.MaxBlockSize {
				p.db.ResetSubSession()
				log.Info("Warning pending block size reach MaxBlockSize")
				break
			}
			p.db.Squash()
			pendingBlockTrx = append(pendingBlockTrx, dtag)
			log.Info("pack apply elapse", common.Elapsed(applyStart))
		}

		removeTransaction(removeTrx)
		block = p.ins.Woker(pendingBlockTrx)
		p.db.ResetSession()
		trxs = nil
		if block != nil {
			log.Errorf("PRODUCER block, hash: %x, delegate: %s, num:%v, trxn:%v, pendingTrxn:%v, blockTime:%s, blockSize %v\n",
				block.Hash(), block.Header.Delegate, block.GetNumber(), len(block.BlockTransactions), pendingTrxlen, time.Unix(int64(block.Header.Timestamp), 0), pendingBlockSize)

			if config.BtoConfig.Delegate.Solo == false {
				ConsensusProducedBlock(block)
			} else {
				AddBlock(block)
			}

		} else {
			return config.PRODUCER_TIME_OUT
		}
		return p.ins.CalcNextReportTime(block)
	}
	return config.PRODUCER_TIME_OUT
}
