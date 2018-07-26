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
 * @Author: eripi
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package block

import (
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
	"time"
)

//Block sync block with peer and send up to block actor
type Block struct {
	actor   *actor.PID
	chainIf chain.BlockChainInterface
	s       *synchronizes

	init bool
}

const (
	//WAIT_TIME wait for actors ready
	WAIT_TIME = 20
)

//MakeBlock new instance
func MakeBlock(chain chain.BlockChainInterface, nodeType bool) *Block {
	return &Block{s: makeSynchronizes(nodeType, chain),
		chainIf: chain,
		init:    false}
}

//SetActor set chain actor id
func (b *Block) SetActor(tid *actor.PID) {
	b.actor = tid
	b.s.setActor(tid)
}

//Start start
func (b *Block) Start() {
	go b.waitActorReady()
}

func (b *Block) waitActorReady() {
	time.Sleep(WAIT_TIME * time.Second)

	log.Debug("protocol wait actor ready")

	blocknumber := b.chainIf.HeadBlockNum()
	libNumber := b.chainIf.LastConsensusBlockNum()

	log.Debugf("protocol timer local block number:%d, %d", libNumber, blocknumber)
	if blocknumber < libNumber {
		panic("protocol wrong lib number")
		return
	}

	b.s.updateLocalLib(libNumber)
	b.s.updateLocalNumber(blocknumber)
	b.s.start()
	b.init = true
}

//Dispatch peer message process
func (b *Block) Dispatch(index uint16, p *p2p.Packet) {
	if !b.init {
		return
	}

	//log.Debugf("block recv packet %d, from peer: %d", p.H.PacketType, index)

	switch p.H.PacketType {
	case BLOCK_REQ:
		b.processBlockReq(index, p.Data, BLOCK_REQ)
	case BLOCK_UPDATE:
		b.processBlockInfo(index, p.Data)
	case LAST_BLOCK_NUMBER_REQ:
		b.processLastBlockNumberReq(index, p.Data)
	case LAST_BLOCK_NUMBER_RSP:
		b.processLastBlockNumberRsp(index, p.Data)
	case BLOCK_HEADER_REQ:
		b.processBlockHeaderReq(index, p.Data)
	case BLOCK_HEADER_RSP:
		b.processBlockHeaderRsp(index, p.Data)
	case BLOCK_HEADER_UPDATE:
		b.processBlockHeaderUpdate(index, p.Data)
	case BLOCK_CATCH_REQUEST:
		b.processBlockReq(index, p.Data, BLOCK_CATCH_REQUEST)
	case BLOCK_CATCH_RESPONSE:
		b.processBlockCatchRsp(index, p.Data)
	}
}

//GetSyncState get current synchronize status
func (b *Block) GetSyncState() bool {
	if b.s.state == STATE_NORMAL {
		return true
	}
	return false
}

//SendNewBlock send out a new block
func (b *Block) SendNewBlock(notify *message.NotifyBlock) {
	b.sendPacket(true, notify.Block, nil)
}

func (b *Block) sendPacket(broadcast bool, data interface{}, peers []uint16) {
	last := chainNumber{LibNumber: b.chainIf.LastConsensusBlockNum(),
		BlockNumber: b.chainIf.HeadBlockNum()}
	b.s.updatec <- last

	buf, err := json.Marshal(data)
	if err != nil {
		log.Errorf("protocol block send marshal error")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_UPDATE,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	if broadcast {
		msg := p2p.BcastMsgPacket{Indexs: peers,
			P: packet}
		p2p.Runner.SendBroadcast(msg)
	} else {
		msg := p2p.UniMsgPacket{Index: peers[0],
			P: packet}
		p2p.Runner.SendUnicast(msg)
	}
}

func (b *Block) processLastBlockNumberReq(index uint16, data []byte) {
	b.s.sendLastBlockNumberRsp(index)
}

func (b *Block) processLastBlockNumberRsp(index uint16, data []byte) {
	var last chainNumber
	err := json.Unmarshal(data, &last)
	if err != nil {
		log.Errorf("protocol processLastBlockNumberRsp Unmarshal error:%s", err)
		return
	}

	info := &peerBlockInfo{
		index:     index,
		lastLib:   last.LibNumber,
		lastBlock: last.BlockNumber,
	}

	b.s.infoc <- info
}

func (b *Block) processBlockHeaderReq(index uint16, data []byte) {
	var req blockHeaderReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Errorf("protocol processBlockHeaderReq Unmarshal err:%s", err)
		return
	}

	if req.Begin > req.End ||
		req.End-req.Begin >= SYNC_BLOCK_BUNDLE_MAX {
		log.Errorf("protocol processBlockHeaderReq wrong lenght")
		return
	}

	var rsp blockHeaderRsp
	j := 0
	for i := req.Begin; i <= req.End; i++ {
		head := b.chainIf.GetHeaderByNumber(i)
		if head == nil {
			log.Errorf("protocol processBlockHeaderReq header:%d not exist", i)
			return
		}

		rsp.set = append(rsp.set, *head)
		j++
	}

	b.sendBlockHeaderRsp(index, &rsp)
}

func (b *Block) processBlockHeaderRsp(index uint16, data []byte) {
	var rsp blockHeaderRsp
	err := json.Unmarshal(data, &rsp.set)
	if err != nil {
		log.Errorf("protocol processBlockHeaderRsp Unmarshal error:%s", err)
		return
	}

	log.Debugf("protocol processBlockHeaderRsp index: %d", index)

	b.s.set.syncheaderc <- &rsp
}

func (b *Block) processBlockReq(index uint16, data []byte, ptype uint16) {
	var blocknumber uint32
	err := json.Unmarshal(data, &blocknumber)
	if err != nil {
		log.Errorf("protocol processBlockReq Unmarshal error:%s", err)
		return
	}

	block := b.chainIf.GetBlockByNumber(blocknumber)
	if block == nil {
		log.Debugf("protocol get block:%d return nil ptype:%d", blocknumber, ptype)

		if ptype == BLOCK_REQ {
			return
		}
	} else {
		log.Debugf("protocol get block number:%d ptype:%d", block.Header.Number, ptype)
	}

	if ptype == BLOCK_REQ {
		b.sendBlockRsp(index, block, BLOCK_UPDATE)
	} else if ptype == BLOCK_CATCH_REQUEST {
		b.sendBlockRsp(index, block, BLOCK_CATCH_RESPONSE)
	} else {
		log.Errorf("protocol processBlockReq error ptype")
	}
}

func (b *Block) processBlockInfo(index uint16, data []byte) {
	var block types.Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		log.Errorf("protocol processBlockInfo Unmarshal error:%s", err)
		return
	}

	update := &blockUpdate{index: index, block: &block}

	b.s.blockc <- update
}

func (b *Block) processBlockCatchRsp(index uint16, data []byte) {
	var block types.Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		log.Errorf("protocol processBlockInfo Unmarshal error:%s", err)
		return
	}

	update := blockUpdate{index: index, block: &block}
	b.s.c.catchupc <- &update
}

func (b *Block) sendBlockHeaderRsp(index uint16, rsp *blockHeaderRsp) {
	data, err := json.Marshal(rsp.set)
	if err != nil {
		log.Error("protocol sendGetBlock Marshal number error ")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_HEADER_RSP,
	}

	packet := p2p.Packet{H: head, Data: data}

	msg := p2p.UniMsgPacket{Index: index,
		P: packet}

	p2p.Runner.SendUnicast(msg)
}

func (b *Block) sendBlockRsp(index uint16, block *types.Block, ptype uint16) {
	data, err := json.Marshal(block)
	if err != nil {
		log.Error("protocol sendGetBlock Marshal number error ")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: ptype,
	}

	packet := p2p.Packet{H: head, Data: data}

	msg := p2p.UniMsgPacket{Index: index,
		P: packet}

	p2p.Runner.SendUnicast(msg)
}

func (b *Block) processBlockHeaderUpdate(index uint16, data []byte) {
	var header types.Header
	err := json.Unmarshal(data, &header)
	if err != nil {
		log.Errorf("protocol processBlockHeaderUpdate Unmarshal error:%s", err)
		return
	}

	update := headerUpdate{index: index, header: &header}

	b.s.headerc <- &update
}
