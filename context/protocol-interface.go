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

package context

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
)

type ProtocolInstance interface {
	ProtocolInterface
	Start()
	SetChainActor(tpid *actor.PID)
	SetTrxPreHandleActor(tpid *actor.PID)
	SetConsensusActor(tpid *actor.PID)

	SendNewTrx(notify *message.NotifyTrx)
	SendNewBlock(notify *message.NotifyBlock)
	SendPrevote(notify *message.SendPrevote)
	SendPrecommit(notify *message.SendPrecommit)
	SendCommit(notify *message.SendCommit)
}
//PeersInfo peersinfo for rest
type PeersInfo struct {
	LastLib uint64
	LastBlock uint64
	Addr string
	Port string
	Account string
	NodeType string
	ChainId string
	Version uint32
	IsActive bool
}
type ProtocolInterface interface {
	GetBlockSyncState() bool
	GetBlockSyncDistance() uint64
	GetPeerInfo()(uint64,[]*PeersInfo)
	UpdatePeerStateToActive(addr string) bool
	UpdatePeerStateToInActive(addr string, timeout uint32) bool
	QueryPeerState(addr string) bool
}
