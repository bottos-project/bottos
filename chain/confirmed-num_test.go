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
 * file description: account role test
 * @Author: Gong Zibin
 * @Date:   2017-12-13
 * @Last Modified by:
 * @Last Modified time:
 */
package chain

import (
	"fmt"
	"sort"
	"testing"

	"github.com/bottos-project/bottos/config"
)

func TestBlockChain_ConfirmedSort(t *testing.T) {
	delegateNum := config.BLOCKS_PER_ROUND
	lastConfirmedNums := make(ConfirmedNum, delegateNum)
	var i uint32
	for i = 0; i < delegateNum; i++ {
		lastConfirmedNums[i] = delegateNum - i
	}

	fmt.Println(lastConfirmedNums)

	consensusIndex := (100 - int(config.CONSENSUS_BLOCKS_PERCENT)) * len(lastConfirmedNums) / 100
	sort.Sort(lastConfirmedNums)
	fmt.Println(lastConfirmedNums)
	fmt.Println(lastConfirmedNums[consensusIndex])
}
