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
 * file description:  api actor
 * @Author:
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package apiactor

import (
	"fmt"

	"github.com/AsynkronIT/protoactor-go/actor"
)

//ApiActorPid is actor pid
var ApiActorPid *actor.PID

//ApiActor is actor props
type ApiActor struct {
	props *actor.Props
}

//ContructApiActor new an actor
func ContructApiActor() *ApiActor {
	return &ApiActor{}
}

//NewApiActor spawn a named actor
func NewApiActor() *actor.PID {
	props := actor.FromProducer(func() actor.Actor { return ContructApiActor() })

	var err error
	ApiActorPid, err = actor.SpawnNamed(props, "ApiActor")

	if err != nil {
		panic(fmt.Errorf("ApiActor SpawnNamed error: ", err))
	} else {
		return ApiActorPid
	}
}

func handleSystemMsg(context actor.Context) {

	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Printf("ApiActor received started msg", msg)
	case *actor.Stopping:
		fmt.Printf("ApiActor received stopping msg")
	case *actor.Restart:
		fmt.Printf("ApiActor received restart msg")
	case *actor.Restarting:
		fmt.Printf("ApiActor received restarting msg")
	}
}

//Receive process msg
func (apiActor *ApiActor) Receive(context actor.Context) {

	handleSystemMsg(context)

	switch msg := context.Message().(type) {
	}
}
