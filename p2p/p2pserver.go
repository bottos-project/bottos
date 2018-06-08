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
 * @Last Modified by:
 * @Last Modified time:
 */

package p2pserver

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	"io/ioutil"
	"sync"
	"time"
)

var actorEnv *env.ActorEnv

//
type P2PServer struct {
	serv      *NetServer
	p2pConfig *P2PConfig

	p2pLock sync.RWMutex
}

type P2PConfig struct {
	ServAddr string
	ServPort int
	PeerLst  []string
}

//parse json configuration
func ReadFile(filename string) *P2PConfig {

	if filename == "" {
		fmt.Println("*ERROR* parmeter is null")
		return &P2PConfig{}
	}
	var pc P2PConfig

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("*ERROR* Failed to read the config: ", filename)
		return &P2PConfig{}
	}

	str := string(bytes)

	if err := json.Unmarshal([]byte(str), &pc); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return &P2PConfig{}
	}

	return &pc
}

//
func NewServ() *P2PServer {
	fmt.Println("NewServ()")

	//config file for test
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
	} else {
		p2pserv = &P2PServer{
			serv:      NewNetServerTst(p2pconfig),
			p2pConfig: p2pconfig,
		}
	}

	return p2pserv
}

func (p2p *P2PServer) Init() error {
	fmt.Println("p2pServer::Init()")
	return nil
}

//it is the entry of p2p
func (p2p *P2PServer) Start() error {
	fmt.Println("p2pServer::Start()")

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

	//start udp server
	go p2p.serv.StartUdpServer()

	//connect to other seed nodes
	go p2p.serv.activeTimedTask()
	if TST == 1 {
		go p2p.reportBlock()
	}

	//peer neighbor exchange
	go p2p.serv.PneStart()

	//get all seeds or wait for 3 seconds
	//go p2p.serv.initSync()
	// Todo ping/pong
	go p2p.RunHeartBeat()

	return nil
}

//run a heart beat to watch the network status
func (p2p *P2PServer) RunHeartBeat() error {
	fmt.Println("p2pServer::RunHeartBeat()")
	return nil
}

func (p2p *P2PServer) SetTrxActor(trxActorPid *actor.PID) {
	p2p.serv.notify.trxActorPid = trxActorPid
}

func (p2p *P2PServer) SetChainActor(chainActorPid *actor.PID) {
	p2p.serv.notify.chainActorPid = chainActorPid
}

func (p2p *P2PServer) SetChainActorPid(tpid *actor.PID) {
	p2p.serv.notify.chainActorPid = tpid
}

func (p2p *P2PServer) SetActorEnv(env *env.ActorEnv) {
	p2p.serv.actorEnv = env
	actorEnv = env
}

//A interface for call from other component
func (p2p *P2PServer) BroadCast(m interface{}, call_type uint8) error {
	fmt.Println("p2pServer::BroadCast()")
	var res error
	switch call_type {
	case TRANSACTION:
		res = p2p.serv.broadCastImpl(m, CRX_BROADCAST)

	case BLOCK:
		res = p2p.serv.broadCastImpl(m, BLK_BROADCAST)

	}

	return res
}

type RsaKeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// Key represents a crypto key that can be compared to another key
type Key interface {
	// Bytes returns a serialized, storeable representation of this key
	Bytes() ([]byte, error)

	// Equals checks whether two PubKeys are the same
	Equals(Key) bool
}

// PrivKey represents a private key that can be used to generate a public key,
// sign data, and decrypt data that was encrypted with a public key
type PrivKey interface {
	Key

	// Cryptographically sign the given bytes
	Sign([]byte) ([]byte, error)

	// Return a public key paired with this private key
	GetPublic() PubKey
}

type PubKey interface {
	Key

	// Verify that 'sig' is the signed hash of 'data'
	Verify(data []byte, sig []byte) (bool, error)
}

/*
// Generates a keypair
func GenerateKeyPairWithReader(typ, bits int, src io.Reader) (PrivKey, PubKey, error) {

	privateKey, err := rsa.GenerateKey(src, bits)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return &RsaKeyPair{privateKey:privateKey}, &RsaKeyPair{ publicKey:publicKey}, nil
}
*/

func (p2p *P2PServer) reportBlock() {
	var timeInterval *time.Timer = time.NewTimer(3 * time.Second)
	var blockNum uint32
	var headerNum uint32

	for {
		select {
		case <-timeInterval.C:

			blockNum = actorEnv.Chain.LastConsensusBlockNum()
			headerNum = actorEnv.Chain.HeadBlockNum()
			SuperPrint(AUQA_PRINT, "P2PServer::reportBlock() blockNum: ", blockNum, " , headerNum: ", headerNum)

			timeInterval.Stop()
			timeInterval.Reset(time.Second * 3)
		}
	}

}

