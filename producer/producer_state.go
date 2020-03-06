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
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	log "github.com/cihub/seelog"
)

//CalcNextReportTime calculate next report time then wait up this producer
func (p *Reporter) CalcNextReportTime(block *types.Block) uint32 {
	if config.BtoConfig.Delegate.Solo == true {
		log.Debug("PRODUCER wait for next report time ", 50)
		return 1000
	}

	slot := uint32(block.Header.Number) % config.BLOCKS_PER_ROUND

	var elapseTime uint32

	if 0 == slot {
		log.Debug("PRODUCER wait for next report time ", 1000)
		return 1000
	} else {
		elapseTime = config.DEFAULT_BLOCK_INTERVAL * (config.BLOCKS_PER_ROUND - slot)
	}

	log.Debug("PRODUCER wait for next report time ", elapseTime*1000)
	return elapseTime * 1000

}
