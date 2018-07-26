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
	"bytes"
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
	"math"
	"sort"
	"sync"
	"time"
)

//DO NOT EDIT
const (
	TIMER_FAST_SYNC_LAST_BLOCK_NUMBER   = 1
	TIMER_NORMAL_SYNC_LAST_BLOCK_NUMBER = 4
	TIMER_CHECK_SYNC_LAST_BLOCK_NUMBER  = 20

	TIMER_SYNC_STATE_CHECK = 5

	TIMER_HEADER_SYNC = 2
	TIMER_BLOCK_SYNC  = 2

	TIMER_CATCHUP   = 2
	CATCHUP_COUNTER = 10

	TIMER_HEADER_UPDATE_CHECK = 1

	SYNC_BLOCK_BUNDLE = 28
)

//DO NOT EDIT
const (
	STATE_SYNCING = 0
	STATE_CATCHUP = 1
	STATE_NORMAL  = 2
)

//DO NOT EDIT
const (
	SET_SYNC_NULL   = 0
	SET_SYNC_HEADER = 1
	SET_SYNC_BLOCK  = 2
)

//DO NOT EDIT
const (
	CATCHUP_COMPLETE = 0
	CATCHUP_DOING    = 1
)

type peerBlockInfo struct {
	index     uint16
	lastLib   uint32
	lastBlock uint32

	counter int16
}

type syncConfig struct {
	nodeType bool
}

type syncset []peerBlockInfo

func (s syncset) Len() int {
	return len(s)
}

func (s syncset) Less(i, j int) bool {
	return s[i].lastLib > s[j].lastLib
}

func (s syncset) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type synchronizes struct {
	peers map[uint16]*peerBlockInfo
	lock  sync.Mutex

	libLocal   uint32
	libRemote  uint32
	lastLocal  uint32
	lastRemote uint32
	state      uint16
	once       bool //have synchronized one time or not when start up

	infoc        chan *peerBlockInfo
	updatec      chan chainNumber
	blockc       chan *blockUpdate
	headerc      chan *headerUpdate
	headercTimer *time.Timer
	headerCache  *headerUpdate

	set *syncSet
	c   *catchup

	config  syncConfig
	chain   *actor.PID
	chainIf chain.BlockChainInterface
}

func makeSynchronizes(nodeType bool, chainIf chain.BlockChainInterface) *synchronizes {
	return &synchronizes{
		peers:   make(map[uint16]*peerBlockInfo),
		infoc:   make(chan *peerBlockInfo),
		updatec: make(chan chainNumber),
		blockc:  make(chan *blockUpdate),
		headerc: make(chan *headerUpdate),
		state:   STATE_SYNCING,
		once:    false,
		set:     makeSyncSet(),
		c:       makeCatchup(),
		config:  syncConfig{nodeType: nodeType},
		chainIf: chainIf,
	}
}

func (s *synchronizes) setActor(tid *actor.PID) {
	s.chain = tid
}

func (s *synchronizes) start() {
	go s.exchangeRoutine()
	go s.checkRoutine()
	go s.syncSetRoutine()
	go s.catchupRoutine()
}

func (s *synchronizes) exchangeRoutine() {
	log.Debug("protocol syncBlockNumberTimer start")

	syncTimer := time.NewTimer(TIMER_FAST_SYNC_LAST_BLOCK_NUMBER * time.Second)
	checkTimer := time.NewTimer(TIMER_CHECK_SYNC_LAST_BLOCK_NUMBER * time.Second)

	defer func() {
		log.Debug("protocol syncBlockNumberTimer stop")
		syncTimer.Stop()
	}()

	for {
		select {
		case <-syncTimer.C:
			s.sendLastBlockNumberReq()
			if s.state == STATE_NORMAL {
				syncTimer.Reset(TIMER_NORMAL_SYNC_LAST_BLOCK_NUMBER * time.Second)
			} else {
				syncTimer.Reset(TIMER_FAST_SYNC_LAST_BLOCK_NUMBER * time.Second)
			}
		case info := <-s.infoc:
			s.recvBlockNumberInfo(info)
		case <-checkTimer.C:
			s.syncBlockNumberCheck()
			checkTimer.Reset(TIMER_CHECK_SYNC_LAST_BLOCK_NUMBER * time.Second)
		}
	}
}

func (s *synchronizes) checkRoutine() {
	checkTimer := time.NewTimer(TIMER_SYNC_STATE_CHECK * time.Second)
	s.headercTimer = time.NewTimer(TIMER_HEADER_UPDATE_CHECK * time.Second)

	for {
		select {
		case number := <-s.updatec:
			s.updateLocalLib(number.LibNumber)
			s.updateLocalNumber(number.BlockNumber)
			if s.state == STATE_SYNCING {
				log.Debugf("protocol local lib update in sync status : %d", s.libLocal)
				s.set.endc <- s.libLocal
			}
		case block := <-s.blockc:
			s.recvBlock(block)
		case <-checkTimer.C:
			s.syncStateCheck()
			checkTimer.Reset(TIMER_SYNC_STATE_CHECK * time.Second)
		case header := <-s.headerc:
			s.recvBlockHeader(header)
		case <-s.headercTimer.C:
			s.checkHeader()
		}
	}
}

func (s *synchronizes) syncSetRoutine() {
	s.set.syncHeaderTimer = time.NewTimer(TIMER_HEADER_SYNC * time.Second)
	s.set.syncBlockTimer = time.NewTimer(TIMER_BLOCK_SYNC * time.Second)

	for {
		select {
		case rsp := <-s.set.syncheaderc:
			if s.set.recvBlockHeader(rsp) {
				s.set.syncHeaderTimer.Stop()
				s.syncBundleBlock()
			}
		case update := <-s.set.syncblockc:
			s.syncRecvBlock(update)
		case number := <-s.set.beginc:
			s.set.begincCheck(number)
		case number := <-s.set.endc:
			s.set.endcCheck(number)
		case <-s.set.syncHeaderTimer.C:
			if s.set.state == SET_SYNC_HEADER {
				s.syncBlockHeader()
			}
		case <-s.set.syncBlockTimer.C:
			s.setSyncStateCheck()
		}
	}
}

func (s *synchronizes) catchupRoutine() {
	check := time.NewTimer(TIMER_CATCHUP * time.Second)

	for {
		select {
		case <-check.C:
			s.catchupCheck()
		case update := <-s.c.catchupc:
			s.catchupRecvBlock(update)
		case <-s.c.stopc:
			s.c.catchupReset()
		}
	}
}

func (s *synchronizes) recvBlock(update *blockUpdate) {
	number := update.block.GetNumber()

	if number <= s.libLocal {
		log.Debugf("protocol drop block: %d is smaller than local number", number)
		return
	}

	if s.state == STATE_NORMAL {
		if number > s.lastLocal+1 {
			log.Debugf("protocol lose block , need catch up with this peer")
			s.state = STATE_CATCHUP
			s.catchupWithPeer(update.index, number)
			return
		}

		if s.sendupBlock(update.block) == chain.InsertBlockSuccess {
			s.broadcastRcvNewBlock(update)

			libNumber := s.chainIf.LastConsensusBlockNum()
			s.updateLocalLib(libNumber)
			blocknumber := s.chainIf.HeadBlockNum()
			s.updateLocalNumber(blocknumber)

			log.Debugf("protocol sendup block success in normal, lib: %d head: %d", libNumber, blocknumber)
		}
		return
	} else if s.state == STATE_CATCHUP {
		if number > s.lastLocal+1 {
			log.Debugf("protocol drop block: %d when in catch up status", number)
			return
		}

		if s.sendupBlock(update.block) == chain.InsertBlockSuccess {
			s.broadcastRcvNewBlock(update)

			libNumber := s.chainIf.LastConsensusBlockNum()
			s.updateLocalLib(libNumber)
			blocknumber := s.chainIf.HeadBlockNum()
			s.updateLocalNumber(blocknumber)

			log.Debugf("protocol sendup block success in catch up, lib: %d head: %d", libNumber, blocknumber)
		}
		return
	} else if s.state == STATE_SYNCING {
		log.Debugf("protocol recv block %d in syncing status", number)
		s.set.syncblockc <- update
	}
}

func (s *synchronizes) syncRecvBlock(update *blockUpdate) {
	if s.set.state != SET_SYNC_BLOCK {
		log.Debugf("protocol drop block: %d when sync header or finish", update.block.GetNumber())
		return
	}

	if update.block.GetNumber() > s.set.end ||
		update.block.GetNumber() < s.set.begin {
		log.Infof("protocol drop block out of sync range")
		return
	}

	log.Infof("protocol sync process block: %d, index: %d", update.block.Header.Number, update.index)

	for i := 0; i < SYNC_BLOCK_BUNDLE; i++ {
		if s.set.headers[i] != nil &&
			s.set.isBlockHeadSame(s.set.headers[i], update.block.Header) {
			s.set.blocks[i] = update.block
			break
		}
	}

	if s.set.setSyncStateJudge() {
		s.sendupBundleBlock()
	}

}

func (s *synchronizes) recvBlockHeader(update *headerUpdate) {
	number := update.header.GetNumber()
	if number <= s.lastLocal {
		log.Debugf("protocol drop block header: %d is smaller than local", number)
		return
	}

	if s.state == STATE_NORMAL && number == s.lastLocal+1 {
		s.cacheHeader(update)
		return
	}

	log.Infof("protocol drop block header: %d , wait for catchup", number)
}

func (s *synchronizes) cacheHeader(update *headerUpdate) {
	s.headerCache = update
	s.headercTimer.Reset(TIMER_HEADER_UPDATE_CHECK * time.Second)
}

func (s *synchronizes) checkHeader() {
	if s.headerCache != nil {
		if s.headerCache.header.Number == s.lastLocal+1 {
			s.sendBlockReq(s.headerCache.index, s.headerCache.header.Number, BLOCK_REQ)
		}

		s.headerCache = nil
	}
}

func (s *synchronizes) recvBlockNumberInfo(info *peerBlockInfo) {
	s.lock.Lock()
	defer s.lock.Unlock()

	peer, ok := s.peers[info.index]
	if ok {
		peer.lastBlock = info.lastBlock
		peer.lastLib = info.lastLib
		peer.counter++
	} else {
		info.counter = 1
		s.peers[info.index] = info
	}

	s.updateRemoteLib(info.lastLib, false)
	s.updateRemoteNumber(info.lastBlock, false)
}

func (s *synchronizes) syncBlockNumberCheck() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for key, info := range s.peers {
		if info.counter == 0 {
			delete(s.peers, key)
		} else {
			info.counter = 0
		}
	}
}

func (s *synchronizes) getPeers() syncset {
	s.lock.Lock()
	s.lock.Unlock()

	var peerset syncset
	for _, info := range s.peers {
		peerset = append(peerset, *info)
	}

	return peerset
}

func (s *synchronizes) syncStateCheck() {
	var remoteLib uint32
	var remoteNumber uint32
	var index uint16

	peerset := s.getPeers()
	catchindex := s.c.index
	var catchremote uint32

	for _, info := range peerset {
		if info.lastLib > remoteLib {
			remoteLib = info.lastLib
		}

		if info.lastBlock > remoteNumber {
			remoteNumber = info.lastBlock
			index = info.index
		}

		if catchindex != 0 && info.index == catchindex {
			catchremote = info.lastBlock
		}
	}

	if remoteNumber == catchremote {
		index = catchindex
	}

	//remote block lib be smaller, wo should reset it
	if remoteLib < s.libRemote {
		log.Errorf("protocol syncStateCheck remote lib number change smaller")
		if remoteLib > 0 {
			s.updateRemoteLib(remoteLib, true)
			s.set.endc <- remoteLib
		}

		//judge by the next time, if no peer exist, sync is always false
		return
	} else if remoteLib > s.libRemote {
		log.Errorf("protocol syncStateCheck remote lib number change bigger")
		s.updateRemoteLib(remoteLib, false)
	}

	//remote block lib be smaller, wo should reset it
	if remoteNumber < s.lastRemote {
		log.Errorf("protocol syncStateCheck remote number change smaller")
		if remoteNumber > 0 {
			s.updateRemoteNumber(remoteNumber, true)
		}

		//judge by the next time, if no peer exist, sync is always false
		return
	} else if remoteNumber > s.lastRemote {
		log.Errorf("protocol syncStateCheck remote number change bigger")
		s.updateRemoteNumber(remoteNumber, false)
	}

	s.syncStateJudge(index)
}

func (s *synchronizes) syncStateJudge(index uint16) {
	if s.libLocal < s.libRemote {
		log.Debugf("protocol syncStateJudge lib small than remote, need sync")

		if !s.once {
			s.syncBlockHeader()
			s.once = true
			return
		}

		if s.lastLocal >= s.lastRemote {
			log.Debugf("protocol syncStateJudge head bigger than remote, sync wait")
			return
		}

		if s.state == STATE_NORMAL ||
			s.state == STATE_CATCHUP {
			log.Debugf("protocol syncStateJudge start syncing")
			s.state = STATE_SYNCING
			s.syncBlockHeader()
			s.c.stopc <- 1
		} else {
			if s.set.state == SET_SYNC_NULL {
				log.Debugf("protocol continue syncing")
				s.syncBlockHeader()
			} else {
				log.Debugf("protocol in syncing statue:%d", s.set.state)
			}
		}
	} else {
		if s.lastLocal < s.lastRemote {
			log.Debugf("protocol syncStateJudge catch up")
			s.state = STATE_CATCHUP
			s.catchupWithPeer(index, s.lastRemote)
		} else {
			s.state = STATE_NORMAL
		}
	}

}

func (s *synchronizes) updateLocalLib(lib uint32) {
	if lib < s.libLocal {
		log.Errorf("protocol update local lib number error now:%d update:%d", s.libLocal, lib)
		return
	} else if lib == s.libLocal {
		return
	}

	log.Debugf("protocol update local lib number:%d", lib)
	s.libLocal = lib
}

func (s *synchronizes) updateLocalNumber(number uint32) {
	if number < s.lastLocal {
		log.Errorf("protocol update local block number error now:%d update:%d", s.lastLocal, number)
		return
	} else if number == s.lastLocal {
		return
	}

	log.Debugf("protocol update local block number:%d", number)
	s.lastLocal = number
}

func (s *synchronizes) updateRemoteLib(lib uint32, force bool) {
	if !force && lib <= s.libRemote {
		return
	}

	log.Debugf("protocol peer max lib number:%d", lib)
	s.libRemote = lib
}

func (s *synchronizes) updateRemoteNumber(number uint32, force bool) {
	if !force && number <= s.lastRemote {
		return
	}

	log.Debugf("protocol peer max block number:%d", number)
	s.lastRemote = number
}

func (s *synchronizes) sendLastBlockNumberReq() {
	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: LAST_BLOCK_NUMBER_REQ,
	}

	packet := p2p.Packet{H: head}

	msg := p2p.BcastMsgPacket{Indexs: nil,
		P: packet}

	p2p.Runner.SendBroadcast(msg)
}

func (s *synchronizes) sendLastBlockNumberRsp(index uint16) {
	rsp := chainNumber{LibNumber: s.libLocal, BlockNumber: s.lastLocal}

	data, err := json.Marshal(rsp)
	if err != nil {
		log.Error("protocol sendGetLastRsp Marshal data error ")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: LAST_BLOCK_NUMBER_RSP,
	}

	packet := p2p.Packet{H: head, Data: data}

	msg := p2p.UniMsgPacket{Index: index,
		P: packet}

	p2p.Runner.SendUnicast(msg)
}

func (s *synchronizes) syncBlockHeader() {
	if s.libRemote <= s.libLocal {
		return
	}

	s.set.reset()

	if s.libLocal+SYNC_BLOCK_BUNDLE > s.libRemote {
		s.set.begin = s.libLocal + 1
		s.set.end = s.libRemote
	} else {
		s.set.begin = s.libLocal + 1
		s.set.end = s.libLocal + SYNC_BLOCK_BUNDLE

	}

	s.set.state = SET_SYNC_HEADER

	log.Debugf("protocol syncBlockHeader begin: %d, end:%d", s.set.begin, s.set.end)

	s.sendBlockHeaderReq(s.set.begin, s.set.end)

	s.set.syncHeaderTimer.Reset(TIMER_HEADER_SYNC * time.Second)
}

func (s *synchronizes) sendBlockHeaderReq(begin uint32, end uint32) {
	header := blockHeaderReq{Begin: begin, End: end}

	data, err := json.Marshal(header)
	if err != nil {
		log.Error("protocol sendBlockHeaderReq Marshal number error ")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_HEADER_REQ,
	}

	packet := p2p.Packet{H: head, Data: data}

	peerset := s.getPeers()
	for _, info := range peerset {
		if info.lastLib >= end {
			msg := p2p.UniMsgPacket{Index: info.index,
				P: packet}

			log.Debugf("protocol sendBlockHeaderReq index: %d", info.index)
			p2p.Runner.SendUnicast(msg)
			break
		}
	}
}

func (s *synchronizes) syncBundleBlock() {
	if s.set.end < s.set.begin {
		return
	}

	var numbers []uint32
	lenght := s.set.end + 1 - s.set.begin
	for i := 0; i < int(lenght) && i < SYNC_BLOCK_BUNDLE; i++ {
		if s.set.blocks[i] == nil {
			numbers = append(numbers, s.set.begin+uint32(i))
		}
	}

	if len(numbers) == 0 {
		return
	}

	peerset := s.getPeers()

	sort.Sort(peerset)

	i := 0
	for _, number := range numbers {
		if i == len(peerset) {
			i = 0
		}

		for i < len(peerset) {
			if peerset[i].lastLib >= number {
				s.sendBlockReq(peerset[i].index, number, BLOCK_REQ)
				i++
				break
			} else {
				i++
			}
		}
	}

	s.set.syncBlockTimer.Reset(TIMER_BLOCK_SYNC * time.Second)
}

func (s *synchronizes) sendBlockReq(index uint16, number uint32, ptype uint16) {

	data, err := json.Marshal(number)
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

	log.Debugf("protocol sendBlockReq block %d, type: %d, index: %d", number, ptype, index)
	p2p.Runner.SendUnicast(msg)
}

func (s *synchronizes) setSyncStateCheck() {
	if s.set.state != SET_SYNC_BLOCK {
		return
	}

	sync := s.set.setSyncStateJudge()
	if sync {
		s.sendupBundleBlock()
	} else {
		s.syncBundleBlock()
	}
}

func (s *synchronizes) sendupBundleBlock() {
	log.Debugf("protocol sync bundle of block finish")

	if s.set.end < s.set.begin {
		return
	}

	if s.set.begin <= s.libLocal {
		log.Errorf("lib local is change bigger, wait next time")
		s.set.reset()
		return
	}

	j := 0
	for i := s.set.begin; i <= s.set.end; i++ {
		if s.sendupBlock(s.set.blocks[j]) != chain.InsertBlockSuccess {
			s.set.blocks[j] = nil
			s.syncBundleBlock()
			return
		}
		j++
	}

	s.libLocal = s.set.end
	s.lastLocal = s.set.end
	log.Debugf("protocol update local lib and number: %d", s.libLocal)

	s.set.reset()

	if s.libLocal < s.libRemote {
		s.syncBlockHeader()
	}
}

func (s *synchronizes) sendupBlock(block *types.Block) uint32 {

	start := common.MeasureStart()
	log.Debugf("protocol send up block :%d", block.Header.Number)

	for i := 0; i < 5; i++ {
		msg := &message.ReceiveBlock{Block: block}

		result, err := s.chain.RequestFuture(msg, 500*time.Millisecond).Result()
		if err != nil {
			log.Errorf("protocol send block request error:%s", err)
			time.Sleep(10 * time.Millisecond)
			continue
		}

		rsp := result.(*message.ReceiveBlockResp)

		if rsp.ErrorNo != chain.InsertBlockSuccess {
			log.Errorf("protocol block insert error: %d", rsp.ErrorNo)
		}

		log.Error("elapsed time 1 ", common.Elapsed(start))

		return rsp.ErrorNo
	}

	log.Error("protocol block insert timeout with five times")

	log.Error("elapsed time 2 ", common.Elapsed(start))
	return 0xff
}

//if node is super node , only broadcast block hearder to some peer
func (s *synchronizes) broadcastRcvNewBlock(update *blockUpdate) {
	if s.config.nodeType {
		s.broadcastNewBlockHeader(update, false)
	} else {
		s.broadcastNewBlock(update, false)
		s.broadcastNewBlockHeader(update, true)
	}
}

func (s *synchronizes) broadcastNewBlock(update *blockUpdate, all bool) {
	buf, err := json.Marshal(update.block)
	if err != nil {
		log.Errorf("protocol block send marshal error")
	}

	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_UPDATE,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	var indexs []uint16
	if all {
		indexs = append(indexs, update.index)
	} else {
		indexs := s.getBcastFilterPeers(update.index)
		if indexs == nil {
			return
		}
	}

	msg := p2p.BcastMsgPacket{Indexs: indexs,
		P: packet}

	p2p.Runner.SendBroadcast(msg)
}

func (s *synchronizes) broadcastNewBlockHeader(update *blockUpdate, all bool) {
	buf, err := json.Marshal(update.block.Header)
	if err != nil {
		log.Errorf("protocol block send marshal error")
	}

	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_HEADER_UPDATE,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	var indexs []uint16
	if all {
		indexs = append(indexs, update.index)
	} else {
		indexs = s.getBcastFilterPeers(update.index)
		if indexs == nil {
			return
		}
	}

	msg := p2p.BcastMsgPacket{Indexs: indexs,
		P: packet}

	p2p.Runner.SendBroadcast(msg)
}

func (s *synchronizes) getBcastFilterPeers(index uint16) []uint16 {
	peers := p2p.Runner.GetPeersData()
	if len(peers) == 0 {
		return nil
	}

	peers = append(peers, p2p.PeerData{Id: p2p.LocalPeerInfo.Id})

	sort.Sort(peers)

	k := 0
	for ; k < len(peers); k++ {
		if peers[k].Id == p2p.LocalPeerInfo.Id {
			break
		}
	}

	number := int(math.Sqrt(float64(len(peers))))

	total := len(peers)
	var filter []p2p.PeerData

	if k+1+number < total {
		if k == 0 {
			filter = append(peers[0:0], peers[number:]...)
		} else {
			filter = append(peers[0:k+1], peers[k+1+number:]...)
		}
	} else if k+1+number == total {
		filter = append(peers[0 : k+1])
	} else {
		if k+1 < total {
			filter = append(peers[k+number-total+number-1 : k+1])
		} else {
			filter = append(peers[number:])
		}
	}

	var indexs []uint16
	for _, peer := range filter {
		indexs = append(indexs, peer.Index)
	}

	indexs = append(indexs, index)

	return indexs
}

func (s *synchronizes) catchupCheck() {
	if s.c.state == CATCHUP_COMPLETE {
		return
	}

	s.c.counter++
	if s.c.counter >= CATCHUP_COUNTER {
		log.Debugf("protocol catchup counter error")
		s.c.catchupReset()
	} else {
		log.Debugf("protocol catchup resend get block: %d", s.c.current)
		s.sendBlockReq(s.c.index, s.c.current, BLOCK_CATCH_REQUEST)
	}
}

func (s *synchronizes) catchupRecvBlock(update *blockUpdate) {
	if s.c.index != update.index {
		return
	}

	if update.block == nil ||
		update.block.Header == nil {
		log.Errorf("protocol catchup with peer index:%d , block:%d finish", s.c.index, s.c.current-1)
		s.c.catchupReset()
		return
	}

	if update.block.Header.Number != s.c.current {
		log.Errorf("protocol catch up recevie wrong block numbe:%d", update.block.Header.Number)
		return
	}

	result := s.sendupBlock(update.block)
	if result == chain.InsertBlockSuccess {
		s.c.current++
		s.c.counter = 0

		s.lastLocal = update.block.Header.Number
		log.Debugf("protocol catchup update local number: %d", s.lastLocal)
		log.Debugf("protocol catchup get next block: %d", s.c.current)

		s.sendBlockReq(s.c.index, s.c.current, BLOCK_CATCH_REQUEST)
	} else if result == chain.InsertBlockErrorNotLinked {
		if s.c.current > s.c.begin {
			log.Errorf("protocol catchup no link, start catchup from last: %d", s.lastLocal)
			s.c.begin = s.lastLocal + 1
			s.c.current = s.c.begin
			s.c.counter = 0
			s.sendBlockReq(s.c.index, s.c.current, BLOCK_CATCH_REQUEST)
		} else if s.c.begin == s.lastLocal+1 {
			log.Errorf("protocol catchup no link, start catchup from lib: %d", s.libLocal)
			s.c.begin = s.libLocal + 1
			s.c.current = s.c.begin
			s.c.counter = 0
			s.sendBlockReq(s.c.index, s.c.current, BLOCK_CATCH_REQUEST)
		} else {
			log.Errorf("protocol catchup with peer:%d error", s.c.index)
			s.c.catchupReset()
		}
	} else {
		log.Errorf("protocol catchup with peer error, reset and wait next time")
		s.c.catchupReset()
	}

}

func (s *synchronizes) catchupWithPeer(index uint16, number uint32) {
	log.Debugf("protocol catch up with peer:%d, number:%d", index, number)

	if s.c.state == CATCHUP_COMPLETE {
		s.c.begin = s.lastLocal + 1
		s.c.current = s.c.begin
		s.c.counter = 0
		s.c.state = CATCHUP_DOING
		s.c.index = index

		s.sendBlockReq(index, s.c.begin, BLOCK_CATCH_REQUEST)
		return
	} else if s.c.state == CATCHUP_DOING {
		if index != s.c.index {
			s.c.index = index
			s.c.counter = 0
			s.sendBlockReq(index, s.c.current, BLOCK_CATCH_REQUEST)
			return
		}
	} else {
		panic("protocol wrong state")
		return
	}
}

type syncSet struct {
	syncheaderc     chan *blockHeaderRsp
	syncblockc      chan *blockUpdate
	syncHeaderTimer *time.Timer
	syncBlockTimer  *time.Timer
	beginc          chan uint32
	endc            chan uint32

	headers [SYNC_BLOCK_BUNDLE]*types.Header
	blocks  [SYNC_BLOCK_BUNDLE]*types.Block

	begin uint32
	end   uint32

	state uint16
}

func makeSyncSet() *syncSet {
	return &syncSet{
		syncheaderc: make(chan *blockHeaderRsp),
		syncblockc:  make(chan *blockUpdate),
		beginc:      make(chan uint32),
		endc:        make(chan uint32),
		state:       SET_SYNC_NULL}
}

func (set *syncSet) recvBlockHeader(rsp *blockHeaderRsp) bool {
	if set.state != SET_SYNC_HEADER {
		log.Errorf("protocol recvBlockHeader state error")
		return false
	}

	if set.end < set.begin {
		return false
	}

	if uint32(len(rsp.set)) != (set.end + 1 - set.begin) {
		log.Errorf("protocol recvBlockHeader rsp length error")
		return false
	}

	check := false
	j := 0
	for i := set.begin; i <= set.end; i++ {
		if rsp.set[j].GetNumber() != i {
			log.Errorf("protocol recvBlockHeader rsp info error number:%d", rsp.set[j].GetNumber())
			check = true
			break
		}

		set.headers[j] = &rsp.set[j]
		j++
	}

	if check {
		set.resetHeader()
		return false
	}

	set.state = SET_SYNC_BLOCK
	return true
}

//endcCheck peer max lib change small if some peer is disconnect
func (set *syncSet) endcCheck(number uint32) {
	if set.state == SET_SYNC_NULL {
		log.Debugf("protocol sync status null")
		return
	}

	//remote lib change small , we should reset and wait for sync judge
	if number < set.end {
		log.Debugf("protocol endcCheck reset end: %d, lib: %d", set.end, number)
		set.reset()
	}
}

//begincCheck local lib change bigger when produce a block in p2p sync state
func (set *syncSet) begincCheck(number uint32) {
	if set.state == SET_SYNC_NULL {
		log.Debugf("protocol sync status null")
		return
	}

	//remote lib change small , we should reset and wait for sync judge
	if number >= set.begin {
		log.Debugf("protocol begincCheck reset begin: %d, lib: %d", set.begin, number)
		set.reset()
	}
}

func (set *syncSet) setSyncStateJudge() bool {
	if set.end < set.begin {
		return true
	}

	lenght := set.end + 1 - set.begin
	for i := 0; i < int(lenght) && i < SYNC_BLOCK_BUNDLE; i++ {
		if set.blocks[i] == nil {
			return false
		}
	}

	return true
}

func (set *syncSet) resetHeader() {
	for i := 0; i < SYNC_BLOCK_BUNDLE; i++ {
		set.headers[i] = nil
	}
}

func (set *syncSet) resetBlock() {
	for i := 0; i < SYNC_BLOCK_BUNDLE; i++ {
		set.blocks[i] = nil
	}
}

func (set *syncSet) reset() {
	set.state = SET_SYNC_NULL
	set.end = 0
	set.begin = 0
	set.resetHeader()
	set.resetBlock()
}

func (set *syncSet) isBlockHeadSame(a *types.Header, b *types.Header) bool {
	if a.Number == b.Number &&
		a.Version == b.Version &&
		a.Timestamp == b.Timestamp &&
		bytes.Equal(a.MerkleRoot, b.MerkleRoot) &&
		bytes.Equal(a.PrevBlockHash, b.PrevBlockHash) {
		return true
	}

	return false
}

type catchup struct {
	catchupc chan *blockUpdate
	stopc    chan int

	index   uint16
	begin   uint32
	current uint32
	counter uint16
	state   uint16
}

func makeCatchup() *catchup {
	return &catchup{
		catchupc: make(chan *blockUpdate),
		stopc:    make(chan int),
	}
}

func (c *catchup) catchupReset() {
	c.index = 0
	c.state = CATCHUP_COMPLETE
	c.begin = 0
	c.current = 0
	c.counter = 0
}
