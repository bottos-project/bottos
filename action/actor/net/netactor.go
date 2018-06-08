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
 * file description:  net actor
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package netactor

import (
	log "github.com/cihub/seelog"
	//"encoding/json"
	"github.com/bottos-project/bottos/action/message"
	//"github.com/bottos-project/core/common/types"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	p2pserv "github.com/bottos-project/bottos/p2p"
)

var NetActorPid *actor.PID = nil

var trxActorPid *actor.PID = nil
var chainActorPid *actor.PID = nil

var p2p *p2pserv.P2PServer = nil

var actorEnv *env.ActorEnv

type NetActor struct {
	props *actor.Props
}

func ContructNetActor() *NetActor {
	return &NetActor{}
}

func NewNetActor(env *env.ActorEnv) *actor.PID {
	actorEnv = env

	p2p = p2pserv.NewServ()
	p2p.SetActorEnv(env)
	go p2p.Start()

	props := actor.FromProducer(func() actor.Actor { return ContructNetActor() })

	var err error
	NetActorPid, err = actor.SpawnNamed(props, "NetActor")

	if err == nil {
		return NetActorPid
	} else {
		panic(log.Errorf("NetActor SpawnNamed error: ", err))
	}

	return nil
}

//main loop
func (NetActor *NetActor) handleSystemMsg(context actor.Context) {
	switch msg := context.Message().(type) {

	case *actor.Started:
		log.Infof("NetActor received started msg", msg)

	case *actor.Stopping:
		log.Info("NetActor received stopping msg")

	case *actor.Restart:
		log.Info("NetActor received restart msg")

	case *actor.Restarting:
		log.Info("NetActor received restarting msg")

	case *message.NotifyTrx:
		log.Infof("%c[%d;%d;%dm%v: %v %c[0m ", 0x1B, 123, 40, 35, "<======================== NetActor received Transaction msg  , msg.Trx: ", msg.Trx, 0x1B)
		go p2p.BroadCast(msg.Trx, p2pserv.TRANSACTION)

	case *message.NotifyBlock:
		log.Infof("%c[%d;%d;%dm%v: %v %c[0m ", 0x1B, 123, 40, 32, "<======================== NetActor received Block msg , msg.Block: ", msg.Block, 0x1B)
		go p2p.BroadCast(msg.Block, p2pserv.BLOCK)

	}

}

func (NetActor *NetActor) Receive(context actor.Context) {
	NetActor.handleSystemMsg(context)

	switch msg := context.Message().(type) {

	}
}

func SetChainActorPid(tpid *actor.PID) {
	//chainActorPid = tpid
	p2p.SetChainActor(tpid)
}

func GetChainActorPid() *actor.PID {
	return chainActorPid
}

func SetTrxActorPid(tpid *actor.PID) {
	//trxActorPid = tpid
	p2p.SetTrxActor(tpid)
}
