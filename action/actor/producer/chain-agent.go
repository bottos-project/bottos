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
 * @Author:
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package produceractor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/action/message"
	"github.com/bottos-project/core/common/types"
)

var chainActorPid *actor.PID

func (p *ProducerActor) SetChainActorPid(tpid *actor.PID) {
	chainActorPid = tpid
}

func ApplyBlock(block *types.Block) {

	applyBlock := &message.InsertBlockReq{block}
	chainActorPid.Tell(applyBlock)

	fmt.Println("send to chain to apply block")

	return
}
