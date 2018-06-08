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
 * file description:  producer entry
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package producer

import (
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/role"
	log "github.com/cihub/seelog"
)

//Reporter is the producer
type Reporter struct {
	core     chain.BlockChainInterface
	roleIntf role.RoleInterface
	state    ReportState
}

//ReporterRepo is the interface of reporters
type ReporterRepo interface {
	Woker(Trxs []*types.Transaction) *types.Block
	IsReady() bool
}

//New is to create new reporter
func New(b chain.BlockChainInterface, roleIntf role.RoleInterface) ReporterRepo {
	stat := ReportState{0, "", "", false, 0, false}
	return &Reporter{core: b, roleIntf: roleIntf, state: stat}
}

//Woker is an actor of repoter
func (p *Reporter) Woker(trxs []*types.Transaction) *types.Block {

	accountName := p.state.ScheduledReporter
	block, err := p.reportBlock(p.state.ScheduledTime, accountName, trxs)
	if err != nil {
		return nil // errors.New("report Block failed")
	}

	return block
}
func (p *Reporter) reportBlock(blockTime uint64, accountName string, trxs []*types.Transaction) (*types.Block, error) {
	head := types.NewHeader()
	head.PrevBlockHash = p.core.HeadBlockHash().Bytes()
	head.Number = p.core.HeadBlockNum() + 1
	head.Timestamp = blockTime
	head.Delegate = []byte(accountName)
	block := types.NewBlock(head, trxs)
	block.Header.DelegateSign = block.Sign("123").Bytes()

	// If this block is last in a round, calculate the schedule for the new round
	if block.Header.Number%config.BLOCKS_PER_ROUND == 0 {
		newSchedule := p.roleIntf.ElectNextTermDelegates(block)
		log.Infof("next term delgates", newSchedule)
		currentState, err := p.roleIntf.GetCoreState()
		if err != nil {
			return nil, err
		}
		block.Header.DelegateChanges = common.Filter(currentState.CurrentDelegates, newSchedule)
	}
	return block, nil
}
