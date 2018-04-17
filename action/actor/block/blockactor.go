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
 * file description:  block actor
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package blockactor

import (
	"fmt"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
)

var BlockActorPid *actor.PID

type BlockActor struct {
	props *actor.Props
}

func ContructBlockActor() *BlockActor {
	return &BlockActor{}
}

func NewBlockActor() *actor.PID {

	props := actor.FromProducer(func() actor.Actor { return ContructBlockActor() })

	BlockActorPid, err := actor.SpawnNamed(props, "BlockActor")

	if err == nil {
		return BlockActorPid
	} else {
		panic(fmt.Errorf("BlockActor SpawnNamed error: ", err))
	}
}

func (BlockActor *BlockActor) handleSystemMsg(context actor.Context) {

	switch msg := context.Message().(type) {

	case *actor.Started:
		log.Printf("BlockActor received started msg", msg)

	case *actor.Stopping:
		log.Printf("BlockActor received stopping msg")

	case *actor.Restart:
		log.Printf("BlockActor received restart msg")

	case *actor.Restarting:
		log.Printf("BlockActor received restarting msg")
	}

}

func (BlockActor *BlockActor) Receive(context actor.Context) {

	BlockActor.handleSystemMsg(context)

	switch msg := context.Message().(type) {

	//case *types.Transaction:

	}
}
