// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

// This program is free software: you can distribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Bottos.  If not, see <http://www.gnu.org/licenses/>.

// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package exec provides functions for executing WebAssembly bytecode.

/*
 * file description: the interface for WASM execution
 * @Author: Stewart Li
 * @Date:   2018-02-08
 * @Last Modified by: Stewart Li
 * @Last Modified time: 2018-06-04
 */

package p2pserver

import (
	"fmt"
	"sync"
	"errors"
	"encoding/json"
	"io/ioutil"
	"github.com/AsynkronIT/protoactor-go/actor"

)

//
type P2PServer struct{
	serv          *NetServer
	p2pConfig     *P2PConfig

	p2pLock        sync.RWMutex
}

type P2PConfig struct {
	ServAddr    string
	ServPort    int
	PeerLst     []string
}

//parse json configuration
func ReadFile(filename string) *P2PConfig{

	if filename == "" {
		fmt.Println("*ERROR* parmeter is null")
		return &P2PConfig{}
	}
	var pc P2PConfig

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("*ERROR* Failed to read the config: ",filename)
		return  &P2PConfig{}
	}

	str:=string(bytes)

	if err := json.Unmarshal([]byte(str), &pc) ; err != nil{
		fmt.Println("Unmarshal: ", err.Error())
		return &P2PConfig{}
	}

	return &pc
}


func NewServ() *P2PServer{

	p2pconfig := ReadFile(CONF_FILE)

	/*
	prvKey, pubKey, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	fmt.Println("prvKey = ",prvKey," , pubKey = ",pubKey)
	*/

	var p2pserv *P2PServer = nil
	if TST == 0 {
		p2pserv = &P2PServer{
			serv:      NewNetServer(),
			p2pConfig: p2pconfig,
		}
	}else{
		p2pserv = &P2PServer{
			serv:      NewNetServerTst(p2pconfig),
			p2pConfig: p2pconfig,
		}
	}

	return p2pserv
}

func (p2p *P2PServer) Init() error {
	return nil
}


//it is the entry of p2p
func (p2p *P2PServer) Start() error {

	if p2p.p2pConfig == nil {
		return errors.New("*ERROR* P2P Configuration hadn't been inited yet !!!")
	}

	if err := p2p.Init(); err != nil {
		return err
	}

	if p2p.serv != nil {
		//wait for connection from others
		p2p.serv.Start()
	}

	//connect to other seed nodes
	go p2p.serv.ActiveSeeds()

	//ping/pong
	go p2p.RunHeartBeat()

	return nil
}

//run a heart beat to watch the network status
func  (p2p *P2PServer) RunHeartBeat() error {
	return nil
}

func  (p2p *P2PServer) SetTrxActor (trxActorPid *actor.PID)  {
	p2p.serv.notify.trxActorPid = trxActorPid
}

func  (p2p *P2PServer) SetChainActor (chainActorPid *actor.PID)  {
	p2p.serv.notify.chainActorPid = chainActorPid
}

func  (p2p *P2PServer) BroadCastImpl(m interface{} , msg_type uint8) error {

	content_byte , err := json.Marshal(m)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the trx message to broadcast : ", err)
		return err
	}

	msg := message {
		Src:     p2p.p2pConfig.ServAddr,
		MsgType: msg_type, // the type to notify other peers new crx
		Content: content_byte,
	}

	msg_byte , err := json.Marshal(msg)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the trx message to broadcast : ", err)
		return err
	}

	p2p.serv.notify.BroadcastByte(msg_byte , false)

	return nil
}

//A interface for call from other component
func  (p2p *P2PServer) BroadCast (m interface{} , call_type uint8) error {

	var res error
	switch call_type{
	case TRANSACTION:
		res = p2p.BroadCastImpl(m , CRX_BROADCAST)

	case BLOCK:
		res = p2p.BroadCastImpl(m , BLK_BROADCAST)

	}

	return res
}










