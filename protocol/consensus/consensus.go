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

package consensus

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/p2p"
)

type Consensus struct {
	actor *actor.PID
}

func MakeConsensus() *Consensus {
	return &Consensus{}
}

func (c *Consensus) SetActor(tid *actor.PID) {
	c.actor = tid
}

func (c *Consensus) Dispatch(index uint16, p *p2p.Packet) {

}

func (c *Consensus) Send(broadcast bool, m interface{}, peers []uint16) {

}

func (c *Consensus) Start() {

}
