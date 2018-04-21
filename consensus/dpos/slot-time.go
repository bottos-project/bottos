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
 * file description:  file introduction for commom tip
 * @Author: May Luo
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */
package dpos

import (
	"time"

	"github.com/bottos-project/core/config"
)

func GetSlotAtTime(current time.Time) uint32 {
	firstSlotTime := GetSlotTime(1)

	if current.Unix() < firstSlotTime {
		return 0
	}
	return uint32(current.Unix()-firstSlotTime)/config.DEFAULT_BLOCK_INTERVAL + 1
}

/**/
func GetLastBlockTimeStamp() int64 {

	return 0 // microseconds
}

func GetHeadBlockNum() int64 {
	return 0
}
func GetGenesisTime() int64 {
	return 1
}
func GetHeadBlockTime() int64 {
	return 1
}
func GetHeadBlockTimeSinceEpoch() int64 {
	return 1
}

func GetSlotTime(slot_num uint32) int64 {
	if slot_num == 0 {
		return GetLastBlockTimeStamp()
	}
	interval := config.DEFAULT_BLOCK_INTERVAL

	if GetHeadBlockNum() == 0 {
		//TODO	return GenesisTimestamp() + int64(slot_num*interval)
		return 1
	}
	GetHeadBlockTime()
	head_block_abs_slot := GetHeadBlockTimeSinceEpoch() / int64(interval)
	head_slot_time := head_block_abs_slot * int64(interval)
	return head_slot_time + int64(slot_num*interval)

}
