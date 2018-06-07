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
	"fmt"
	"net"
	"sync"
	"time"
	"errors"
	"encoding/json"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/common/types"
	msgDef "github.com/bottos-project/bottos/action/message"
)

const (
	SYN_BLK_NUM = 10
)

var isSynced bool = true
var syncLock sync.RWMutex

func GetSyncStatus() bool {
	syncLock.RLock()
	defer syncLock.RUnlock()

	return isSynced
}

type NetServer struct {
	config          *P2PConfig
	port            int
	addr            string

	notify          *NotifyManager
	listener        net.Listener

	seedPeer        []string
	connPeerNum     int

	neighborList    []*net.UDPAddr
	serverAddr      *net.UDPAddr
	socket          *net.UDPConn
	//todo publicKey to identify credit peer
	publicKey       string

	timeInterval   *time.Timer
	syncLock        sync.RWMutex

	actorEnv        *env.ActorEnv
	isSync          bool

	sync.RWMutex
}

func NewNetServer() *NetServer {

	return &NetServer{
		addr:          config.Param.ServAddr,
		seedPeer:      config.Param.PeerList,
		port:          config.Param.P2PPort,
		notify:        NewNotifyManager(),
		connPeerNum:   0,
		actorEnv:      nil,
		isSync:        false,
		timeInterval:  time.NewTimer(TIME_INTERVAL * time.Second),
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
		seedPeer:      config.PeerLst,
		addr:          config.ServAddr,
		port:          config.ServPort,
		notify:        NewNotifyManager(),
		connPeerNum:   0,
		actorEnv:      nil,
		timeInterval:  time.NewTimer(TIME_INTERVAL * time.Second),
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

		go serv.handleMessage(conn)
	}
}

func (serv *NetServer) requestSyncLock() bool {
	serv.Lock()
	defer serv.Unlock()
	if serv.isSync {
		return true
	} else {
		serv.isSync = true
		return false
	}
}


func (serv *NetServer) releaseSyncLock() {
	serv.Lock()
	defer serv.Unlock()

	serv.isSync = false
}

//run accept
func (serv *NetServer) handleMessage(conn net.Conn) {

	data := make([]byte, 4096)
	var msg CommonMessage

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
		serv.handleRequest(msg)

	case RESPONSE:
		serv.handleResponse(msg)

	case CRX_BROADCAST:
		serv.handleCrxBroadcast(msg)

	case BLK_BROADCAST:
		serv.handleBlkBroadcast(msg)

	case BLOCK_INFO:
		serv.handleBlkInfo(msg)

	case BLOCK_REQ:
		serv.handleBlkReq(msg)

	case BLOCK_RES:
		serv.handleBlkRes(msg)

	case DEFAULT:
		SuperPrint(YELLO_PRINT , "DEFAULT ")
	}

	return
}

func (serv *NetServer) activeTimedTask() error {
	fmt.Println("p2pServer::ActiveSeeds()")
	for {
		select {
		case <- serv.timeInterval.C:
			serv.connectSeeds()
			serv.watchStatus()
			serv.broadcastBlkInfo()
			serv.resetTimer()
		}
	}
}

func (serv *NetServer) appendList(conn net.Conn , msg CommonMessage) error {
	//package remote peer info as "peer" struct and add it into peer list
	fmt.Println("NetServer::AppendList")
	peer := NewPeer(msg.Src , serv.port , conn)
	peer.SetPeerState(ESTABLISH)
	serv.notify.addPeer(peer)

	return nil
}

//reset time to start timer for a new round
func  (serv *NetServer) resetTimer()  {
	serv.timeInterval.Stop()
	serv.timeInterval.Reset(time.Second * TIME_INTERVAL)
}

//connect seed during start p2p server
func (serv *NetServer) connectSeeds() error {

	for _ , peer := range serv.seedPeer {
		//check if the new peer is in peer list
		if serv.notify.isExist(peer , false) {
			continue
		}

		var msg = CommonMessage {
			Src:      serv.addr,
			Dst:      peer,
			MsgType:  REQUEST,
		}

		req , err := json.Marshal(msg)
		if err != nil {
			return err
		}

		//connect remote seed peer , if it's successful , add it into remote peer list
		go serv.Send(peer , req , false)
	}

	return nil
}

//to connect specified peer
func (serv *NetServer) SendTo (conn net.Conn , msg []byte , isExist bool) error {
	fmt.Println("p2pServer::SendTo")
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
func (serv *NetServer) Send (addr string , msg []byte , isExist bool) error {
	addrPort := addr+":"+fmt.Sprint(serv.port)
	conn , err := net.Dial("tcp", addrPort)
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

	conn.Close()

	return nil
}

//to connect certain peer with udp
func (serv *NetServer) ConnectUDP(addr string , msg []byte , isExist bool) error {

	addrPort := addr+":"+fmt.Sprint(serv.port)
	remoteAddr, err := net.ResolveUDPAddr("udp4", addrPort)
	if err != nil {
		return errors.New("*ERROR* Failed to create a remote server addr !!!")
	}

	_ , err = serv.socket.WriteToUDP(msg , remoteAddr)
	if err != nil { //todo check len
		fmt.Println("*ERROR* Failed to send Test message to remote peer !!! ",err)
		return errors.New("*ERROR* Failed to send Test message to remote peer !!!")
	}

	return nil
}

func (serv *NetServer) watchStatus() {

	blockNum  := serv.actorEnv.Chain.LastConsensusBlockNum()
	headerNum := serv.actorEnv.Chain.HeadBlockNum()

	SuperPrint(BLUE_PRINT,"NetServer::WatchStatus() blockNum: ",blockNum," , headerNum: ", headerNum)

	for _, peer := range serv.notify.peerMap {
		SuperPrint(BLUE_PRINT,"*** NetServer::WatchStatus() current status: peer = ",peer.peerAddr," ***")
	}

}

func  (serv *NetServer) broadCastImpl(m interface{} , msgType uint8) error {
	fmt.Println("P2PServer::BroadCastImpl")

	contentByte , err := json.Marshal(m)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the trx message to broadcast : ", err)
		return err
	}

	msg := CommonMessage {
		Src:     serv.addr,
		MsgType: msgType, // the type to notify other peers new crx
		Content: contentByte,
	}

	msgByte , err := json.Marshal(msg)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the trx message to broadcast : ", err)
		return err
	}

	serv.notify.broadcastByte(msgByte , false)

	return nil
}

// broadcast current blk info(height and so on) to peerMap
func (serv *NetServer) broadcastBlkInfo()  {

	peerMap := serv.notify.getPeerMap()
	for _ , peer := range peerMap {
		go serv.sendBklInfo(peer)
	}
}

func (serv *NetServer) sendBklInfo (peer *Peer) {
	if peer.syncState != ESTABLISH {
		return
	}

	if serv.actorEnv == nil {
		return
	}

	blockNum  := serv.actorEnv.Chain.LastConsensusBlockNum()
	headerNum := serv.actorEnv.Chain.HeadBlockNum()

	//no generated blk and return
	if blockNum <= 0 && headerNum <= 0 {
		return
	}

	blockInfo := BlockInfo {
		BlockNum:  serv.actorEnv.Chain.LastConsensusBlockNum(),
		HeaderNum: serv.actorEnv.Chain.HeadBlockNum(),
	}

	blockInfoByte , err := json.Marshal(blockInfo)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the blockinfo to broadcast : ", err)
		return
	}

	msg := CommonMessage {
		Src:     serv.addr,
		Dst:     peer.peerAddr,
		MsgType: BLOCK_INFO, // the type to notify other peers new crx
		Content: blockInfoByte,
	}

	msgByte , err := json.Marshal(msg)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the trx message to broadcast : ", err)
		return
	}

	//SuperPrint(PURPLISH_RED_PRINT , "NetServer::SendBklInfo() send msg: ", msg , " , blockInfoByte = ", blockInfoByte)

	peer.SendTo(msgByte , false)

	return
}

func (serv *NetServer) syncBlock(srcAddr string , blockInfo *BlockInfo) error {
	//if true means it is synchronsizing else to start synchronsize
	//it enable just one goruntine is running for the function
	if serv.requestSyncLock() {
		return nil
	}
	defer serv.releaseSyncLock()

	//Get block info at local
	//blockNum  := actorEnv.Chain.LastConsensusBlockNum()
	headerNum := actorEnv.Chain.HeadBlockNum()
	gap       := blockInfo.BlockNum - headerNum
	if gap <= 0 {
		syncLock.Lock()
		defer syncLock.Unlock()
		isSynced = true
		return nil
	}

	//set synced = false
	syncLock.Lock()
	isSynced = false
	syncLock.Unlock()
	//SuperPrint(PURPLISH_RED_PRINT , "NetServer::syncBlock() BLOCK_INFO block_num = " , block_num , " , block_info = ",block_info)

	//if local header_num < remote header_num , request remote peer to sync
	//todo blockNum < blockInfo.BlockNum
	for i := headerNum + 1; i <=  blockInfo.HeaderNum; i++ {
		//use block id to require block from other peer
		serv.reqBlock(srcAddr , i)
	}

	return nil
}

func (serv *NetServer) reqBlock (addr string , blockId uint32) error {
	blockReq := BlockReq {
		BlockNum: blockId,
	}

	blockReqByte , err := json.Marshal(blockReq)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the blockinfo to broadcast : ", err)
		return err
	}

	msg := CommonMessage {
		Src:     serv.addr,
		Dst:     addr,
		MsgType: BLOCK_REQ, // the type to notify other peers new crx
		Content: blockReqByte,
	}

	msgByte , err := json.Marshal(msg)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the trx message to broadcast : ", err)
		return err
	}

	serv.Send(addr , msgByte , false )

	SuperPrint(PURPLISH_RED_PRINT , "NetServer::ReqBlock()  sync blk req: ", msg )

	return nil
}


func (serv *NetServer) handleRequest (msg CommonMessage) {
	//receive a connection request from other peer passively
	rsp := CommonMessage {
		Src:      serv.addr,
		Dst:      msg.Src,
		MsgType:  RESPONSE,
	}

	data , err := json.Marshal(rsp)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the response message : ", err)
		return
	}

	//create a new conn to response the remote peer
	remoteConn , err := net.Dial("tcp", msg.Src+":"+fmt.Sprint(serv.port))
	if err != nil {
		fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ",err)
		return
	}

	len , err := remoteConn.Write(data)
	if err != nil {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
		return
	} else if len < 0 {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
		return
	}

	//remoteConn.SetDeadline(time.Now().Add(20 * time.Second))
	serv.appendList(remoteConn , msg)
}

func (serv *NetServer) handleResponse (msg CommonMessage) {
	//a response from my proactive connect
	//if the remote peer hadn't existed at local , add it into local
	if serv.notify.isExist(msg.Src , false) {
		return
	}


	remoteConn , err := net.Dial("tcp", msg.Src+":"+fmt.Sprint(serv.port))
	if err != nil {
		fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ",err)
		return
	}
	//remoteConn.SetDeadline(time.Now().Add(20 * time.Second))
	serv.appendList(remoteConn , msg)
}


func (serv *NetServer) handleCrxBroadcast (msg CommonMessage) {
	//Receive crx boardcast from other peer , and set it to txpool
	var newCrx types.Transaction
	err := json.Unmarshal(msg.Content , &newCrx)
	if err != nil {
		fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
		return
	}

	recvTrx := msgDef.PushTrxReq{
		Trx:       &newCrx,
		TrxSender: msgDef.TrxSenderTypeP2P,
	}
	SuperPrint(YELLO_PRINT , "<<<<<<<<<<<<<<<<<<<<<< NetServer::HandleMessage from:",msg.Src," newCrx = ",newCrx)

	if serv.notify.trxActorPid != nil {
		fmt.Println("NetServer::HandleMessage() send new crx to trxActor: ",recvTrx)
		serv.notify.trxActorPid.Tell(&recvTrx)
	}

	//todo broadcast to other peers
}

func (serv *NetServer) handleBlkBroadcast (msg CommonMessage) {
	//Receive blk boardcast from other peer
	var newBlk types.Block
	err := json.Unmarshal(msg.Content , &newBlk)
	if err != nil {
		fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
		return
	}

	//build a new message struct (ReceiveBlock) to send to chainactor
	recvBlk := msgDef.ReceiveBlock{
		Block:   &newBlk,
	}
	SuperPrint(YELLO_PRINT,"<<<<<<<<<<<<<<<<<<<<<< NetServer::HandleMessage from:",msg.Src," newBlk = ",newBlk)

	if serv.notify.chainActorPid != nil {
		SuperPrint(YELLO_PRINT,"NetServer::HandleMessage() send new crx to chainActor")
		serv.notify.chainActorPid.Tell(&recvBlk)
	}

	//todo broadcast to other peers
}

func (serv *NetServer) handleBlkInfo (msg CommonMessage) {
	//Receive broadcast from other peers
	var blockInfo BlockInfo
	err := json.Unmarshal(msg.Content , &blockInfo)
	if err != nil {
		fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
		return
	}

	//SuperPrint(PURPLISH_RED_PRINT , "NetServer::HandleMessage() blockInfo: ", blockInfo ," , msg= " , msg)
	go serv.syncBlock(msg.Src , &blockInfo)
}

func (serv *NetServer) handleBlkReq (msg CommonMessage) {
	var blockReq BlockReq
	err := json.Unmarshal(msg.Content , &blockReq)
	if err != nil {
		fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
		return
	}

	blk := actorEnv.Chain.GetBlockByNumber(blockReq.BlockNum)
	blkByte , err := json.Marshal(blk)
	if err != nil{
		fmt.Println("*WRAN* Failed to package the trx message to broadcast : ", err)
		return
	}

	var req = CommonMessage {
		Src:      serv.addr,
		Dst:      msg.Src,
		MsgType:  BLOCK_RES,
		Content:  blkByte,
	}

	reqByte , err := json.Marshal(req)
	if err != nil {
		return
	}

	SuperPrint(PURPLISH_RED_PRINT, "NetServer::SyncBlock()  BLOCK_REQ  send back : ", msg.Src)
	serv.Send(msg.Src , reqByte , false)
}

func (serv *NetServer) handleBlkRes (msg CommonMessage) {
	//SuperPrint(PURPLISH_RED_PRINT,"NetServer::HandleMessage() BLOCK_RES")
	var newBlk types.Block
	err := json.Unmarshal(msg.Content , &newBlk)
	if err != nil {
		fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
		return
	}

	//build a new message struct (ReceiveBlock) to send to chainactor
	recvBlk := msgDef.ReceiveBlock{
		Block:   &newBlk,
	}

	if serv.notify.chainActorPid != nil {
		SuperPrint(PURPLISH_RED_PRINT , "NetServer::HandleMessage() send BLOCK_RES to actor: ",newBlk)
		serv.notify.chainActorPid.Tell(&recvBlk)
	}
}

func (serv *NetServer) matchMinConnection() bool {
	return int(serv.notify.getPeerCnt())+1 >= MIN_NODE_NUM
}

func (serv *NetServer) syncFinished() bool {
	//check all peers about blk height
	return true
}

//goruntine,
func (serv *NetServer) initSync()  {
	if serv.matchMinConnection() {
		syncLock.Lock()
		isSynced = true
		syncLock.Unlock()

		return
	}

	time.Sleep(INIT_SYNC_WAIT * time.Second)
	syncLock.Lock()
	isSynced = true
	syncLock.Unlock()

	return
}