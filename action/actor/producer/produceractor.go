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
	"log"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/action/env"
	"github.com/bottos-project/core/action/message"
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/producer"
)

type ProducerActor struct {
	myDB *db.DBService
	ins  producer.ReporterRepo
}

func NewProducerActor(env *env.ActorEnv) *actor.PID {

	ins := producer.New(env.Chain)
	props := actor.FromProducer(func() actor.Actor {
		return &ProducerActor{env.Db, ins}
	})

	pid, err := actor.SpawnNamed(props, "ProducerActor")

	if err != nil {
		return nil
	}

	return pid
}

func (p *ProducerActor) handleSystemMsg(context actor.Context) {
	switch msg := context.Message().(type) {

	case *actor.Started:
		log.Printf("ProducerActor received started msg", msg)
		context.SetReceiveTimeout(500 * time.Millisecond)

	case *actor.ReceiveTimeout:
		p.working()
		context.SetReceiveTimeout(500 * time.Millisecond)

	case *actor.Stopping:
		log.Printf("ProducerActor received stopping msg")

	case *actor.Restart:
		log.Printf("ProducerActor received restart msg")

	case *actor.Restarting:
		log.Printf("ProducerActor received restarting msg")
	}

}

func (p *ProducerActor) Receive(context actor.Context) {

	p.handleSystemMsg(context)

	switch msg := context.Message().(type) {

	case *message.GetAllPendingTrxRsp:
		fmt.Println("receive pending........", msg)

	}
}
func (p *ProducerActor) working() {
	if p.ins.IsReady() {
		fmt.Println("Ready to generate block")
		trxs := GetAllPendingTrx()
		p.myDB.StartUndoSession()

		p.ins.VerifyTrxs(trxs)

		block := p.ins.Woker(trxs)

		if block != nil {
			fmt.Println("apply block", block)
			ApplyBlock(block)
			//TODO brocast
		}
	}
}
