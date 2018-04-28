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
 * file description: producer
 * @Author: May Luo
 * @Date:   2017-12-11
 * @Last Modified by:
 * @Last Modified time:
 */
package producer

import (
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/config"
)

func (p *Reporter) GetSlotAtTime(current uint64) uint32 {
	firstSlotTime := p.GetSlotTime(1)

	if current < firstSlotTime {
		return 0
	}
	return uint32(current-firstSlotTime)/config.DEFAULT_BLOCK_INTERVAL + 1
}

func (p *Reporter) GetSlotTime(slotNum uint32) uint64 {

	if slotNum == 0 {
		return 0
	}
	interval := config.DEFAULT_BLOCK_INTERVAL

	object, err := p.roleIntf.GetChainState()
	if err != nil {
		return 0
	}
	genesisTime := p.core.GenesisTimestamp()
	if object.LastBlockNum == 0 {

		return genesisTime + uint64(slotNum*interval)
	}
	headBlockAbsSlot := common.GetSecondSincEpoch(object.LastBlockTime, genesisTime) / uint64(interval)
	headSlotTime := headBlockAbsSlot * uint64(interval)
	return headSlotTime + uint64(slotNum*interval)
}
