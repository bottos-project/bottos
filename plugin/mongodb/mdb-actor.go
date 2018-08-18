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
 * file description:  mdb actor
 * @Author:
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package mongodb

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	log "github.com/cihub/seelog"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/common/types"
)

type MdbActor struct {
	Mdb *MongoDBPlugin
}

//NewMdbActor spawn a named actor
func NewMdbActor(roleIntf role.RoleInterface, db *db.OptionDBService) *actor.PID {
	mdb := NewMongoDBPlugin(roleIntf, db)
	props := actor.FromProducer(func() actor.Actor {
		return &MdbActor{Mdb: mdb}
	})
	pid, err := actor.SpawnNamed(props, "MdbActor")

	if err != nil {
		log.Errorf("mdb actor fail")
		return nil
	}

	return pid
}

func handleSystemMsg(context actor.Context) bool {
	switch context.Message().(type) {
	case *actor.Started:
		log.Info("MdbActor received started msg")
	case *actor.Stopping:
		log.Info("MdbActor received stopping msg")
	case *actor.Restart:
		log.Info("MdbActor received restart msg")
	case *actor.Restarting:
		log.Info("MdbActor received restarting msg")
	case *actor.Stop:
		log.Info("MdbActor received Stop msg")
	case *actor.Stopped:
		log.Info("MdbActor received Stopped msg")
	default:
		return false
	}

	return true
}

func (actor *MdbActor) Receive(context actor.Context) {
	if handleSystemMsg(context) {
		return
	}

	switch msg := context.Message().(type) {
	case *types.Block:
		actor.HandleReceiveBlock(context, msg)
	default:
		log.Error("MdbActor received Unknown msg")
	}
}

//HandleNewProducedBlock new block msg
func (actor *MdbActor) HandleReceiveBlock(ctx actor.Context, block *types.Block) {
	actor.Mdb.ApplyBlock(block)
}
