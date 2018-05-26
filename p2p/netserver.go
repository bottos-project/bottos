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
	//"io"
	"fmt"
	"net"
	"sync"
	"time"
	"errors"
	//"reflect"
	//"unsafe"
	//"crypto/sha1"
	//"hash/fnv"
	"encoding/json"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/common/types"
)

type NetServer struct {
	config          *P2PConfig
	port            int
	addr            string

	notify          *NotifyManager
	listener        net.Listener

	seed_peer       []string

	neighborList    []*net.UDPAddr
	serverAddr      *net.UDPAddr
	socket          *net.UDPConn
	//todo
	publicKey       string

	time_interval   *time.Timer
	netLock         sync.RWMutex
}

func NewNetServer() *NetServer {
	fmt.Println("NewNetServer config.Param.PeerList = ", config.Param.PeerList)
	return &NetServer{
		//config:        config,
		addr:          config.Param.ServAddr,
		seed_peer:     config.Param.PeerList,
		port:          config.Param.P2PPort,
		notify:        NewNotifyManager(),
		time_interval: time.NewTimer(TIME_INTERVAL * time.Second),
	}
}

//for UT
func NewNetServerTst(config *P2PConfig) *NetServer {
	if config == nil {
		fmt.Println("*ERROR* Parmeter is empty !!!")
		return nil
	}

	return &NetServer{
		config:        config,
		seed_peer:     config.PeerLst,
		addr:          config.ServAddr,
		port:          config.ServPort,
		notify:        NewNotifyManager(),
		time_interval: time.NewTimer(TIME_INTERVAL * time.Second),
	}
}

//start listener
func (serv *NetServer) Start() error {
	fmt.Println("netServer::Start()")

	go serv.Listening()

	return nil
}

//run accept
func (serv *NetServer) Listening() {
	fmt.Println("NetServer::Listening()")
	listener, err := net.Listen("tcp", ":"+fmt.Sprint(serv.port))
	if err != nil {
		fmt.Println("*ERROR* Failed to listen at port: "+fmt.Sprint(serv.port))
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

//run accept
func (serv *NetServer) HandleMessage(conn net.Conn) {

	data := make([]byte, 4096)
	var msg message

	len , err := conn.Read(data)
	if err != nil {
		fmt.Println("*WRAN* Can't read data from remote peer !!!")
		return
	}

	err = json.Unmarshal(data[0:len] , &msg)
	if err != nil {
		fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
		return
	}

	switch msg.MsgType {
	case REQUEST:
		//receive a connection request from other peer passively
		rsp := message {
			Src:      serv.addr,
			Dst:      msg.Src,
			MsgType:  RESPONSE,
		}

		data , err := json.Marshal(rsp)
		if err != nil{
			fmt.Println("*WRAN* Failed to package the response message : ", err)
		}

		//create a new conn to response the remote peer
		remote_conn , err := net.Dial("tcp", msg.Src+":"+fmt.Sprint(serv.port))
		if err != nil {
			fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ",err)
			return
		}

		len , err := remote_conn.Write(data)
		if err != nil {
			fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
			return
		} else if len < 0 {
			fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
			return
		}

		serv.AppendList(remote_conn , msg)

	case RESPONSE:
		//a response from my proactive connect
		//if the remote peer hadn't existed at local , add it into local
		fmt.Println("NetServer::HandleMessage() response to = ", msg.Src)
		if serv.notify.IsExist(msg.Src , false) {
			return
		}

		remote_conn , err := net.Dial("tcp", msg.Src+":"+fmt.Sprint(serv.port))
		if err != nil {
			fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ",err)
			return
		}

		serv.AppendList(remote_conn , msg)

	case CRX_BROADCAST:
		//Todo receive crx_boardcast from other peer , and set it to txpool
		fmt.Println("NetServer::HandleMessage()")

		var new_crx types.Transaction
		err = json.Unmarshal(msg.Content , &new_crx)
		if err != nil {
			fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
			return
		}

		if serv.notify.trxActorPid != nil {
			fmt.Println("NetServer::HandleMessage() send new_crx: ",new_crx)
			serv.notify.trxActorPid.Tell(new_crx)
		}

	case BLK_BROADCAST:
		//Todo receive blk_boardcast from other peer
		fmt.Println("NetServer::HandleMessage()")

	}

	return
}

func (serv *NetServer) ActiveSeeds() error {
	fmt.Println("p2pServer::ActiveSeeds()")
	for {
		select {
		case <- serv.time_interval.C:
			serv.ConnectSeeds()
			serv.WatchStatus()
			serv.ResetTimer()
		}
	}
}

func (serv *NetServer) AppendList(conn net.Conn , msg message) error {
	//package remote peer info as "peer" struct and add it into peer list
	fmt.Println("NetServer::AppendList")
	peer := NewPeer(msg.Src , serv.port , conn)
	peer.SetPeerState(ESTABLISH)
	serv.notify.AddPeer(peer)

	return nil
}

//reset time to start timer for a new round
func  (serv *NetServer) ResetTimer ()  {
	serv.time_interval.Stop()
	serv.time_interval.Reset(time.Second * TIME_INTERVAL)
}

//connect seed during start p2p server
func (serv *NetServer) ConnectSeeds() error {

	fmt.Println("p2pServer::ConnectSeeds")
	for _ , peer := range serv.seed_peer {
		//check if the new peer is in peer list
		if serv.notify.IsExist(peer , false) {
			continue
		}

		var msg = message {
			Src:      serv.addr,
			Dst:      peer,
			MsgType:  REQUEST,
		}

		req , err := json.Marshal(msg)
		if err != nil {
			return err
		}

		//connect remote seed peer , if it's successful , add it into remote_list
		go serv.Connect(peer , req , false)
	}

	return nil
}

//to connect specified peer
func (serv *NetServer) ConnectTo (conn net.Conn , msg []byte , isExist bool) error {
	fmt.Println("p2pServer::ConnectTo")
	if conn == nil {
		return errors.New("*ERROR* Invalid parameter !!!")
	}

	len , err := conn.Write(msg)
	if err != nil {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
		return err
	} else if len < 0 {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
		return err
	}

	return nil
}

//to connect to certain peer proactively
func (serv *NetServer) Connect(addr string , msg []byte , isExist bool) error {
	addr_port := addr+":"+fmt.Sprint(serv.port)
	conn , err := net.Dial("tcp", addr_port)
	if err != nil {
		fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ",err)
		return err
	}

	len , err := conn.Write(msg)
	if err != nil {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
		return err
	} else if len < 0 {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
		return err
	}

	return nil
}


//to connect certain peer with udp
func (serv *NetServer) ConnectUDP(addr string , msg []byte , isExist bool) error {
	fmt.Println("p2pServer::ConnectSeed() addr = ",addr)

	addr_port := addr+":"+fmt.Sprint(serv.port)
	remoteAddr, err := net.ResolveUDPAddr("udp4", addr_port)
	if err != nil {
		return errors.New("*ERROR* Failed to create a remote server addr !!!")
	}

	/*
	//test connection with remote peer
	var msg = message {
		src:      serv.addr,
		dst:      addr,
		msg_type: request,
	}

	req , err := json.Marshal(msg)
	if err != nil {
		return err
	}
	*/

	_ , err = serv.socket.WriteToUDP(msg , remoteAddr)
	if err != nil { //todo check len
		fmt.Println("*ERROR* Failed to send Test message to remote peer !!! ",err)
		return errors.New("*ERROR* Failed to send Test message to remote peer !!!")
	}

	/*
	//package remote peer info as "peer" struct and add it into peer list
	peer := NewPeer(addr)
	peer_identify := Hash(addr_port)
	serv.peerMap[uint64(peer_identify)] = peer
	*/

	return nil
}

func (serv *NetServer) WatchStatus() {
	fmt.Println("NetServer::WatchStatus")

	for key, peer := range serv.notify.peerMap {
		fmt.Println("<------------ NetServer::WatchStatus() current status: key = ", key, " , peer = ", peer.peerAddr)
	}
	//serv.notify.BoardcastTrx(nil , false)
}
