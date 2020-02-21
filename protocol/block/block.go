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
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/bottos/chain"
	berr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	"github.com/bottos-project/bottos/version"
	log "github.com/cihub/seelog"
)

//Block sync block with peer and send up to block actor
type Block struct {
	actor   *actor.PID
	chainIf chain.BlockChainInterface
	S       *synchronizes

	headerc chan *headerReq
	init    bool
}

const (
	//WAIT_TIME wait for actors ready
	WAIT_TIME = 20
)

//MakeBlock new instance
func MakeBlock(chain chain.BlockChainInterface, nodeType bool) *Block {
	return &Block{S: makeSynchronizes(nodeType, chain),
		chainIf: chain,
		headerc: make(chan *headerReq),
		init:    false}
}

//SetActor set chain actor id
func (b *Block) SetActor(tid *actor.PID) {
	b.actor = tid
	b.S.setActor(tid)
}

//Start start
func (b *Block) Start() {
	go b.waitActorReady()
	go b.routine()
}

func (b *Block) waitActorReady() {
	time.Sleep(WAIT_TIME * time.Second)

	blocknumber := b.chainIf.HeadBlockNum()
	libNumber := b.chainIf.LastConsensusBlockNum()

	log.Debugf("PROTOCOL timer local block number:%d, %d", libNumber, blocknumber)
	if blocknumber < libNumber {
		panic("PROTOCOL wrong lib number")
		return
	}

	b.S.updateLocalLib(libNumber)
	b.S.updateLocalNumber(blocknumber)
	b.S.start()
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

func (b *Block) routine() {
	for {
		select {
		case r := <-b.headerc:
			b.processHeaderReq(r)
		}
	}
}

//GetSyncState get current synchronize status
func (b *Block) GetSyncState() bool {
	if b.S.state == STATE_NORMAL {
		return true
	}
	return false
}

//GetSyncState get current synchronize number
func (b *Block) GetSyncDistance() uint64 {
	if b.S.state == STATE_NORMAL {
		return 0
	}
	return  b.S.libRemote - b.S.libLocal
}

func (b *Block) UpdateHeadNumber() {
	lastNum := b.chainIf.HeadBlockNum()

	log.Debugf("PROTOCOL update  head number %d", lastNum)
	b.S.updateHeadc <- lastNum
}

func (b *Block) UpdateNumber() {
	version := version.GetVersionNumByBlockNum(b.chainIf.HeadBlockNum())
	number := &chainNumber{
		LibNumber:    b.chainIf.LastConsensusBlockNum(),
		BlockNumber:  b.chainIf.HeadBlockNum(),
		BlockVersion: version,
	}

	log.Debugf("PROTOCOL update lib number: %d head number %d,version %d", number.LibNumber, number.BlockNumber, version)
	b.S.updateLibc <- number
}

//SendNewBlock send out a new block
func (b *Block) SendNewBlock(notify *message.NotifyBlock) {
	log.Debugf("PROTOCOL boardcast new block")
	b.sendPacket(true, notify.Block, nil)
}

func (b *Block) sendPacket(broadcast bool, block *types.Block, peers []uint16) {
	last := b.chainIf.HeadBlockNum()
	b.S.updateHeadc <- last

	log.Debugf("PROTOCOL update head number %d", last)
	buf, err := block.Marshal()
	if err != nil {
		log.Errorf("PROTOCOL block send marshal error")
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
	b.S.sendLastBlockNumberRsp(index)
}

func (b *Block) processLastBlockNumberRsp(index uint16, data []byte) {
	var last chainNumber
	err := bpl.Unmarshal(data, &last)
	if err != nil {
		log.Errorf("PROTOCOL processLastBlockNumberRsp Unmarshal error:%s", err)
		return
	}

	info := &peerBlockInfo{
		Index:            index,
		LastLib:          last.LibNumber,
		LastBlock:        last.BlockNumber,
		LastBlockVersion: last.BlockVersion,
	}

	log.Debugf("PROTOCOL processLastBlockNumberRsp index %d lib:%d  head: %d head version: %d", index, last.LibNumber, last.BlockNumber, last.BlockVersion)

	b.S.infoc <- info
}

func (b *Block) processBlockHeaderReq(index uint16, data []byte) {
	var req blockHeaderReq
	err := bpl.Unmarshal(data, &req)
	if err != nil {
		log.Errorf("PROTOCOL processBlockHeaderReq Unmarshal err:%s", err)
		return
	}

	headReq := &headerReq{index: index, req: &req}
	b.headerc <- headReq
}

func (b *Block) processHeaderReq(r *headerReq) {
	if r.req.Begin > r.req.End ||
		r.req.End-r.req.Begin >= SYNC_BLOCK_BUNDLE_MAX {
		log.Errorf("PROTOCOL processBlockHeaderReq wrong length")
		return
	}

	log.Debugf("PROTOCOL processHeaderReq")

	var rsp blockHeaderRsp
	j := 0
	for i := r.req.Begin; i <= r.req.End; i++ {
		head := b.chainIf.GetHeaderByNumber(i)
		if head == nil {
			log.Errorf("PROTOCOL processBlockHeaderReq header:%d not exist", i)
			return
		}

		rsp.set = append(rsp.set, *head)
		j++
	}

	log.Debugf("PROTOCOL send head response")

	b.sendBlockHeaderRsp(r.index, &rsp)
}

func (b *Block) processBlockHeaderRsp(index uint16, data []byte) {
	var rsp blockHeaderRsp
	err := bpl.Unmarshal(data, &rsp.set)
	if err != nil {
		log.Errorf("PROTOCOL processBlockHeaderRsp Unmarshal error:%s", err)
		return
	}

	log.Debugf("PROTOCOL processBlockHeaderRsp index: %d", index)

	b.S.set.syncheaderc <- &rsp
}

func (b *Block) processBlockReq(index uint16, data []byte, ptype uint16) {
	var req syncReq
	err := bpl.Unmarshal(data, &req)
	if err != nil {
		log.Errorf("PROTOCOL processBlockReq Unmarshal error:%s", err)
		return
	}
	localVersion := version.GetVersionNumByBlockNum(req.Number)
	log.Debugf("PROTOCOL get block:%d, %d, num %d", index, ptype, req.Number)
	if localVersion != req.Version {
		log.Errorf("PROTOCOL blocknum %d version not match localVersion %d,version %d ", req.Number, localVersion, req.Version)
		return
	}
	block := b.chainIf.GetBlockByNumber(req.Number)
	if block == nil {
		log.Errorf("PROTOCOL get block:%d return nil ptype:%d", req.Number, ptype)
		return
		/*if ptype == BLOCK_REQ {
			return
		}*/
	} else {
		log.Debugf("PROTOCOL get block number:%d ptype:%d", block.Header.Number, ptype)
	}

	if ptype == BLOCK_REQ {
		b.sendBlockRsp(index, block, BLOCK_UPDATE)
	} else if ptype == BLOCK_CATCH_REQUEST {
		b.sendBlockRsp(index, block, BLOCK_CATCH_RESPONSE)
	} else {
		log.Errorf("PROTOCOL processBlockReq error ptype")
	}
}

func (b *Block) processBlockInfo(index uint16, data []byte) {
	vblock := &types.Block{}
	if err := bpl.UnmarshalUntilField(data, vblock, "Number"); err != nil { //unmarshal version field first
		log.Errorf("VERSION handle processBlockInfo unmarshal block version failed, peer index %v, err %v", index, err)
		return
	}
	if err := version.CheckBlock(vblock, "processBlockInfo"); err != berr.ErrNoError {
		return
	}

	var block types.Block
	pblock := &block
	err := pblock.Unmarshal(data)
	if err != nil {
		log.Errorf("PROTOCOL processBlockInfo Unmarshal error:%s", err)
		return
	}

	update := &blockUpdate{index: index, block: &block}

	b.S.blockc <- update
}

func (b *Block) processBlockCatchRsp(index uint16, data []byte) {
	var block types.Block
	pblock := &block
	err := pblock.Unmarshal(data)
	if err != nil {
		log.Errorf("PROTOCOL processBlockInfo Unmarshal error:%s", err)
		return
	}

	update := blockUpdate{index: index, block: &block}
	b.S.c.catchupc <- &update
}

func (b *Block) sendBlockHeaderRsp(index uint16, rsp *blockHeaderRsp) {
	data, err := bpl.Marshal(rsp.set)
	if err != nil {
		log.Error("PROTOCOL sendGetBlock Marshal number error ")
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
	data, err := block.Marshal()
	if err != nil {
		log.Error("PROTOCOL sendGetBlock Marshal number error ")
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
	err := bpl.Unmarshal(data, &header)
	if err != nil {
		log.Errorf("PROTOCOL processBlockHeaderUpdate Unmarshal error:%s", err)
		return
	}

	update := headerUpdate{index: index, header: &header}

	b.S.headerc <- &update
}
