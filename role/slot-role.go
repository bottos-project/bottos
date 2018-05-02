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
package role

import (
	"fmt"

	"github.com/bottos-project/core/config"
)

func (r *Role) GetSlotAtTime(current uint64) uint32 {
	firstSlotTime := r.GetSlotTime(1)

	if current < firstSlotTime {
		return 0
	}
	return uint32(current-firstSlotTime)/config.DEFAULT_BLOCK_INTERVAL + 1
}

func (r *Role) GetSlotTime(slotNum uint32) uint64 {

	if slotNum == 0 {
		return 0
	}
	interval := config.DEFAULT_BLOCK_INTERVAL

	object, err := r.GetChainState()
	if err != nil {
		return 0
	}

	if object.LastBlockNum == 0 {
		genesisTime := config.Genesis.GenesisTime
		fmt.Println("LastBlockNum  ======0", genesisTime)
		return genesisTime + uint64(slotNum*interval)
	}
	headBlockAbsSlot := object.LastBlockTime / uint64(interval)
	headSlotTime := headBlockAbsSlot * uint64(interval)
	return headSlotTime + uint64(slotNum*interval)
}
