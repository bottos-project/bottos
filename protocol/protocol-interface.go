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
 * @Author: eripi
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package protocol

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
)

type ProtocolInstance interface {
	ProtocolInterface
	Start()
	SetChainActor(tpid *actor.PID)
	SetTrxActor(tpid *actor.PID)
	SetProducerActor(tpid *actor.PID)

	ProcessNewTrx(notify *message.NotifyTrx)
	ProcessNewBlock(notify *message.NotifyBlock)
}

type ProtocolInterface interface {
	GetBlockSyncState() bool
}
