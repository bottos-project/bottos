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
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package produceractor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

var ProducerActorPid *actor.PID

type ProducerActor struct {
	props *actor.Props
}

func ContructProducerActor() *ProducerActor {
	return &ProducerActor{}
}

func NewProducerActor() *actor.PID {

	props := actor.FromProducer(func() actor.Actor { return ContructProducerActor() })

	ProducerActorPid, err := actor.SpawnNamed(props, "ProducerActor")

	if err == nil {
		return ProducerActorPid
	} else {
		panic(fmt.Errorf("ProducerActor SpawnNamed error: ", err))
	}
}

func (ProducerActor *ProducerActor) handleSystemMsg(msg interface{}) {

	switch msg.(type) {

	case *actor.Started:
		//log.Info("ProducerActor received started msg")

	case *actor.Stopping:
		//log.Warn("ProducerActor received stopping msg")

	case *actor.Restart:
		//log.Warn("ProducerActor received restart msg")

	case *actor.Restarting:
		//log.Warn("ProducerActor received restarting msg")
	}
}

func (ProducerActor *ProducerActor) Receive(context actor.Context) {

	ProducerActor.handleSystemMsg(context.Message)

	switch msg := context.Message().(type) {

	//case *types.Transaction:

	}
}
