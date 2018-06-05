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
 * @Last Modified time: 2018-06-05
 */

package p2pserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	msgDef "github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
)

//NetServer p2p server
type NetServer struct {
	config *P2PConfig
	port   int
	addr   string

	notify   *NotifyManager
	listener net.Listener

	seedPeer []string

	neighborList []*net.UDPAddr
	serverAddr   *net.UDPAddr
	socket       *net.UDPConn
	//todo
	publicKey string

	timeInterval *time.Timer
	netLock       sync.RWMutex
}

//NewNetServer new a p2p server
func NewNetServer() *NetServer {

	return &NetServer{
		addr:          config.Param.ServAddr,
		seedPeer:     config.Param.PeerList,
		port:          config.Param.P2PPort,
		notify:        NewNotifyManager(),
		timeInterval: time.NewTimer(TIME_INTERVAL * time.Second),
	}
}

//NewNetServerTst new a p2p server with test configure
func NewNetServerTst(config *P2PConfig) *NetServer {
	if config == nil {
		fmt.Println("*ERROR* Parmeter is empty !!!")
		return nil
	}

	return &NetServer{
		config:        config,
		seedPeer:     config.PeerLst,
		addr:          config.ServAddr,
		port:          config.ServPort,
		notify:        NewNotifyManager(),
		timeInterval: time.NewTimer(TIME_INTERVAL * time.Second),
	}
}

//Start start listener
func (serv *NetServer) Start() error {

	go serv.Listening()

	return nil
}

//Listening listen and run accept
func (serv *NetServer) Listening() {
	listener, err := net.Listen("tcp", ":"+fmt.Sprint(serv.port))
	if err != nil {
		fmt.Println("*ERROR* Failed to listen at port: " + fmt.Sprint(serv.port))
		return
	}

	defer listener.Close()

	//main loop for listen
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("NetServer::Listening() Failed to accept")
			continue
		}

		go serv.HandleMessage(conn)
	}
}

//HandleMessage process message
func (serv *NetServer) HandleMessage(conn net.Conn) {

	data := make([]byte, 4096)
	var msg message

	len, err := conn.Read(data)
	if err != nil {
		fmt.Println("*WRAN* Can't read data from remote peer !!!")
		return
	}

	err = json.Unmarshal(data[0:len], &msg)
	if err != nil {
		fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
		return
	}

	switch msg.MsgType {
	case REQUEST:
		//receive a connection request from other peer passively
		rsp := message{
			Src:     serv.addr,
			Dst:     msg.Src,
			MsgType: RESPONSE,
		}

		data, err := json.Marshal(rsp)
		if err != nil {
			fmt.Println("*WRAN* Failed to package the response message : ", err)
			return
		}

		//create a new conn to response the remote peer
		remoteConn, err := net.Dial("tcp", msg.Src+":"+fmt.Sprint(serv.port))
		if err != nil {
			fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ", err)
			return
		}

		len, err := remoteConn.Write(data)
		if err != nil {
			fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ", err)
			return
		} else if len < 0 {
			fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ", err)
			return
		}

		serv.AppendList(remoteConn, msg)

	case RESPONSE:
		//a response from my proactive connect
		//if the remote peer hadn't existed at local , add it into local
		if serv.notify.IsExist(msg.Src, false) {
			return
		}

		remoteConn, err := net.Dial("tcp", msg.Src+":"+fmt.Sprint(serv.port))
		if err != nil {
			fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ", err)
			return
		}

		serv.AppendList(remoteConn, msg)

	case CRX_BROADCAST:
		//Receive crx_boardcast from other peer , and set it to txpool

		var newCrx types.Transaction
		err = json.Unmarshal(msg.Content, &newCrx)
		if err != nil {
			fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
			return
		}

		recvTrx := msgDef.PushTrxReq{
			Trx:       &newCrx,
			TrxSender: msgDef.TrxSenderTypeP2P,
		}

		fmt.Printf("%c[%d;%d;%dm%s %v %v %v %c[0m ", 0x1B, 123, 40, 33, "<<<<<<<<<<<<<<<<<<<<<< NetServer::HandleMessage from:", msg.Src, " new_crx = ", newCrx, 0x1B)

		if serv.notify.trxActorPid != nil {
			fmt.Println("NetServer::HandleMessage() send new_crx to trxActor: ", recvTrx)
			serv.notify.trxActorPid.Tell(&recvTrx)
		}

	case BLK_BROADCAST:
		//Receive blk_boardcast from other peer

		var newBlk types.Block
		err = json.Unmarshal(msg.Content, &newBlk)
		if err != nil {
			fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
			return
		}

		//build a new message struct (ReceiveBlock) to send to chainactor
		recvBlk := msgDef.ReceiveBlock{
			Block: &newBlk,
		}

		if serv.notify.chainActorPid != nil {
			fmt.Println("NetServer::HandleMessage() send new_crx to chainActor: ", recvBlk)
			serv.notify.chainActorPid.Tell(recvBlk)
		}
	}

	return
}

//ActiveSeeds active seeds
func (serv *NetServer) ActiveSeeds() error {

	for {
		select {
		case <-serv.timeInterval.C:
			serv.ConnectSeeds()
			serv.WatchStatus()
			serv.ResetTimer()
		}
	}
}

//AppendList add a establish peer
func (serv *NetServer) AppendList(conn net.Conn, msg message) error {
	//package remote peer info as "peer" struct and add it into peer list
	peer := NewPeer(msg.Src, serv.port, conn)
	peer.SetPeerState(ESTABLISH)
	serv.notify.AddPeer(peer)

	return nil
}

//ResetTimer reset time to start timer for a new round
func (serv *NetServer) ResetTimer() {
	serv.timeInterval.Stop()
	serv.timeInterval.Reset(time.Second * TIME_INTERVAL)
}

//ConnectSeeds connect seed during start p2p server
func (serv *NetServer) ConnectSeeds() error {

	for _, peer := range serv.seedPeer {
		//check if the new peer is in peer list
		if serv.notify.IsExist(peer, false) {
			continue
		}

		var msg = message{
			Src:     serv.addr,
			Dst:     peer,
			MsgType: REQUEST,
		}

		req, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		//connect remote seed peer , if it's successful , add it into remote_list
		go serv.Send(peer, req, false)

	}

	return nil
}

//SendTo to connect specified peer
func (serv *NetServer) SendTo(conn net.Conn, msg []byte, isExist bool) error {

	if conn == nil {
		return errors.New("*ERROR* Invalid parameter !!!")
	}

	len, err := conn.Write(msg)
	if err != nil {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ", err)
		return err
	} else if len < 0 {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ", err)
		return err
	}

	return nil
}

//Send to connect to certain peer proactively
func (serv *NetServer) Send(addr string, msg []byte, isExist bool) error {
	addPort := addr + ":" + fmt.Sprint(serv.port)
	conn, err := net.Dial("tcp", addPort)

	if err != nil {
		fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ", err)
		return err
	}

	len, err := conn.Write(msg)
	if err != nil {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ", err)
		return err
	} else if len < 0 {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ", err)
		return err
	}

	return nil
}

//ConnectUDP to connect certain peer with udp
func (serv *NetServer) ConnectUDP(addr string, msg []byte, isExist bool) error {

	addPort := addr + ":" + fmt.Sprint(serv.port)
	remoteAddr, err := net.ResolveUDPAddr("udp4", addPort)
	if err != nil {
		return errors.New("*ERROR* Failed to create a remote server addr !!!")
	}

	_, err = serv.socket.WriteToUDP(msg, remoteAddr)
	if err != nil { //todo check len
		fmt.Println("*ERROR* Failed to send Test message to remote peer !!! ", err)
		return errors.New("*ERROR* Failed to send Test message to remote peer !!!")
	}

	return nil
}

//WatchStatus debug peers state
func (serv *NetServer) WatchStatus() {
	for key, peer := range serv.notify.peerMap {
		fmt.Println("*** NetServer::WatchStatus() current status: key = ", key, " , peer = ", peer.peerAddr, " ***")
	}

}

//Ping send ping message to a peer
func (serv *NetServer) Ping() {

	peerMap := serv.notify.GetPeerMap()
	for _, peer := range peerMap {
		if peer.syncState != ESTABLISH {
			continue
		}
		//Todo send msg to test connection with other peers
	}

}
