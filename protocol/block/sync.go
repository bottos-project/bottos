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
	"math"
	"sort"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common"
	berr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	"github.com/bottos-project/bottos/version"
	log "github.com/cihub/seelog"
)

//DO NOT EDIT
const (
	TIMER_FAST_SYNC_LAST_BLOCK_NUMBER   = 2
	TIMER_NORMAL_SYNC_LAST_BLOCK_NUMBER = 4
	TIMER_CHECK_SYNC_LAST_BLOCK_NUMBER  = 20

	TIMER_SYNC_STATE_CHECK  = 5
	TIMER_SYNC_STATE_CHECK1 = 1

	TIMER_HEADER_SYNC = 2
	TIMER_BLOCK_SYNC  = 2

	TIMER_CATCHUP   = 2
	CATCHUP_COUNTER = 10

	TIMER_HEADER_UPDATE_CHECK = 1

	SYNC_BLOCK_BUNDLE     = 60
	SYNC_BLOCK_BUNDLE_MAX = 200

	SYNC_HEADER_BUNDLE = 3
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
	Index            uint16
	LastLib          uint64
	LastBlock        uint64
	LastBlockVersion uint32

	syncTimeoutCounter int16
	exchangeCounter    int16
}

type syncConfig struct {
	nodeType bool
}

type syncset []peerBlockInfo

func (s syncset) Len() int {
	return len(s)
}

func (s syncset) Less(i, j int) bool {
	return s[i].syncTimeoutCounter < s[j].syncTimeoutCounter
}

func (s syncset) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type syncsetlib syncset

func (s syncsetlib) Len() int {
	return len(s)
}

func (s syncsetlib) Less(i, j int) bool {
	return s[i].LastLib < s[j].LastLib
}

func (s syncsetlib) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type synchronizes struct {
	Peers map[uint16]*peerBlockInfo
	lock  sync.Mutex

	libLocal  uint64
	libRemote uint64

	lastLocal        uint64
	lastLocalVersion uint32

	lastRemote        uint64
	lastRemoteVersion uint32

	state uint16
	once  bool //have synchronized one time or not when start up

	infoc        chan *peerBlockInfo
	updateLibc   chan *chainNumber
	updateHeadc  chan uint64
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
		Peers:       make(map[uint16]*peerBlockInfo),
		infoc:       make(chan *peerBlockInfo),
		updateLibc:  make(chan *chainNumber),
		updateHeadc: make(chan uint64),
		blockc:      make(chan *blockUpdate),
		headerc:     make(chan *headerUpdate),
		state:       STATE_SYNCING,
		once:        false,
		set:         makeSyncSet(),
		c:           makeCatchup(),
		config:      syncConfig{nodeType: nodeType},
		chainIf:     chainIf,
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
		case number := <-s.updateLibc:
			s.updateLocalNumber(number.BlockNumber)
			s.updateLocalLib(number.LibNumber)
			if s.state == STATE_SYNCING {
				log.Debugf("protocol local lib update in sync status : %d", s.libLocal)
				s.set.beginc <- s.libLocal
			}
		case number := <-s.updateHeadc:
			s.updateLocalNumber(number)
		case block := <-s.blockc:
			s.recvBlock(block)
		case <-checkTimer.C:
			if s.syncStateCheck() {
				checkTimer.Reset(TIMER_SYNC_STATE_CHECK1 * time.Second)
			} else {
				checkTimer.Reset(TIMER_SYNC_STATE_CHECK * time.Second)
			}
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
				s.checkSyncHeaderTimeoutPeer()
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
			log.Debugf("protocol lose block , need catch up with this peer index:%d,number:%d", update.index, number)
			s.state = STATE_CATCHUP
			s.catchupWithPeer(update.index, number)
			return
		}

		if s.sendupBlock(update.block) == berr.ErrNoError {
			s.broadcastRcvNewBlock(update)

			blocknumber := s.chainIf.HeadBlockNum()
			s.updateLocalNumber(blocknumber)

			log.Debugf("protocol sendup block success in normal, head: %d", blocknumber)
		}
		return
	} else if s.state == STATE_CATCHUP {
		if number > s.lastLocal+1 {
			log.Debugf("protocol drop block: %d when in catch up status", number)
			return
		}

		if s.sendupBlock(update.block) == berr.ErrNoError {
			s.broadcastRcvNewBlock(update)

			blocknumber := s.chainIf.HeadBlockNum()
			s.updateLocalNumber(blocknumber)

			log.Debugf("protocol sendup block success in catch up, head: %d", blocknumber)
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
	myVersion := update.header.Version
	lastVersion := version.GetVersionNumByBlockNum(s.lastLocal)
	if myVersion <= lastVersion {
		log.Debugf("protocol drop block header version%d : is smaller than local %d", myVersion, lastVersion)
		return
	}
	number := update.header.GetNumber()
	if number <= s.lastLocal {
		log.Debugf("protocol drop block header: %d is smaller than local %d ", number, s.lastLocal)
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
	//check version if is lower than delete this peer
	expectRemoteVersion := version.GetVersionNumByBlockNum(info.LastBlock)
	if info.LastBlockVersion < expectRemoteVersion {
		log.Errorf("protocol remote version %d header %d smaller %d ",info.LastBlockVersion, info.LastBlock, expectRemoteVersion )
		_, ok := s.Peers[info.Index]
		if ok{
			delete(s.Peers, info.Index)
		}
		return
	}

	peer, ok := s.Peers[info.Index]
	if ok {
		peer.LastBlock = info.LastBlock
		peer.LastLib = info.LastLib
		peer.LastBlockVersion = info.LastBlockVersion
		peer.exchangeCounter++
	} else {
		info.exchangeCounter = 1
		s.Peers[info.Index] = info
	}

	s.updateRemoteLib(info.LastLib, false)
	s.updateRemoteNumber(info.LastBlock, info.LastBlockVersion, false)
}

func (s *synchronizes) syncBlockNumberCheck() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for key, info := range s.Peers {
		expectRemoteVersion := version.GetVersionNumByBlockNum(info.LastBlock)
		if info.LastBlockVersion < expectRemoteVersion {
			delete(s.Peers, key)
		}
		if info.exchangeCounter == 0 {
			delete(s.Peers, key)
		} else {
			info.exchangeCounter = 0
		}
	}
}

func (s *synchronizes) recordPeerSyncTimeout(index uint16) {
	s.lock.Lock()
	defer s.lock.Unlock()

	peer, ok := s.Peers[index]
	if ok {
		peer.syncTimeoutCounter++
	}
}

func (s *synchronizes) resetPeerSyncTimeout() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, info := range s.Peers {
		info.syncTimeoutCounter = 0
	}
}

func (s *synchronizes) getPeers() syncset {
	s.lock.Lock()
	defer s.lock.Unlock()

	var peerset syncset
	for _, info := range s.Peers {
		peerset = append(peerset, *info)
	}

	return peerset
}

func (s *synchronizes) syncStateCheck() (syncFlag bool) {
	var remoteLib uint64
	var remoteNumber uint64
	var remoteNumVersion uint32
	var index uint16

	//we can't judge where peer exist or not because we need in sync status when only one node
	peerset := s.getPeers()

	catchindex := s.c.index
	var catchremote uint64
	var catchremoteVersion uint32

	for _, info := range peerset {
		if info.LastBlockVersion < remoteNumVersion {
			continue
		} else if info.LastBlockVersion > remoteNumVersion {
			remoteLib = info.LastLib
			remoteNumber = info.LastBlock
			remoteNumVersion = info.LastBlockVersion
			catchremote = info.LastBlock
			catchremoteVersion = info.LastBlockVersion
			index = info.Index
		} else {
			if info.LastLib > remoteLib {
				remoteLib = info.LastLib
			}

			if info.LastBlock > remoteNumber {
				remoteNumber = info.LastBlock
				remoteNumVersion = info.LastBlockVersion
				index = info.Index
			}

			if catchindex != 0 && info.Index == catchindex {
				catchremote = info.LastBlock
				catchremoteVersion = info.LastBlockVersion
			}
		}
	}

	if remoteNumber == catchremote && remoteNumVersion == catchremoteVersion {
		index = catchindex
	}

	if remoteNumVersion < s.lastRemoteVersion {
		if remoteLib > 0 {
			s.updateRemoteLib(remoteLib, true)
			s.set.endc <- remoteLib
		}
		if remoteNumber > 0 {
			s.updateRemoteNumber(remoteNumber, remoteNumVersion, true)
		}

	} else {

		//remote block lib be smaller, wo should reset it
		if remoteLib < s.libRemote {
			log.Errorf("protocol syncStateCheck remote lib number change smaller")
			if remoteLib > 0 {
				s.updateRemoteLib(remoteLib, true)
				s.set.endc <- remoteLib
			}

			//judge by the next time, if no peer exist, sync could be always false
			return
		} else if remoteLib > s.libRemote {
			log.Errorf("protocol syncStateCheck remote lib number change bigger")
			s.updateRemoteLib(remoteLib, false)
		}

		//remote block lib be smaller, wo should reset it
		if remoteNumber < s.lastRemote {
			log.Errorf("protocol syncStateCheck remote number change smaller")
			if remoteNumber > 0 {
				s.updateRemoteNumber(remoteNumber, remoteNumVersion, true)
			}

			//judge by the next time, if no peer exist, sync could be always false
			return
		} else if remoteNumber > s.lastRemote {
			log.Errorf("protocol syncStateCheck remote number change bigger")
			s.updateRemoteNumber(remoteNumber, remoteNumVersion, false)
		}
	}

	flag := s.syncStateJudge(index)
	return flag
}

func (s *synchronizes) syncStateJudge(index uint16) (syncFlag bool) {
	localNumVersion := version.GetVersionNumByBlockNum(s.lastLocal)
	if localNumVersion > s.lastRemoteVersion {
		return true
	}
	if s.libLocal < s.libRemote {
		log.Debugf("protocol syncStateJudge lib small than remote, need sync %d,%d,version %d,%d", s.libLocal, s.libRemote, localNumVersion, s.lastRemoteVersion)

		if !s.once {
			s.state = STATE_SYNCING
			s.syncBlockHeader()
			s.once = true
			return true
		}

		if s.lastLocal >= s.lastRemote {
			log.Debugf("protocol syncStateJudge head bigger than remote, sync wait")
			return true
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
			return false
		}
	}
	return true

}

func (s *synchronizes) updateLocalLib(lib uint64) {
	if lib < s.libLocal {
		log.Errorf("protocol update local lib number error now:%d update:%d", s.libLocal, lib)
		return
	} else if lib == s.libLocal {
		return
	}

	log.Debugf("protocol update local lib number:%d", lib)
	s.libLocal = lib
}

func (s *synchronizes) updateLocalNumber(number uint64) {
	if number < s.lastLocal {
		log.Errorf("protocol update local block number error now:%d update:%d", s.lastLocal, number)
		return
	} else if number == s.lastLocal {
		log.Debugf("protocol update head number same %d", number)
		return
	}

	log.Debugf("protocol update local block number:%d", number)
	s.lastLocalVersion = version.GetVersionNumByBlockNum(number)
	s.lastLocal = number
}

func (s *synchronizes) updateRemoteLib(lib uint64, force bool) {
	if !force && lib <= s.libRemote {
		return
	}

	log.Debugf("protocol peer max lib number:%d", lib)
	s.libRemote = lib
}

func (s *synchronizes) updateRemoteNumber(number uint64, numVersion uint32, force bool) {
	if numVersion < s.lastRemoteVersion {
		return
	}
	if !force && number <= s.lastRemote {
		return
	}

	log.Debugf("protocol peer max block number:%d", number)
	s.lastRemote = number
	s.lastRemoteVersion = numVersion
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
	myVersion := version.GetVersionNumByBlockNum(s.lastLocal)
	rsp := chainNumber{LibNumber: s.libLocal, BlockNumber: s.lastLocal, BlockVersion: myVersion}

	data, err := bpl.Marshal(rsp)
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

	log.Debugf("protocol sendGetLastRsp lib:%d head: %d version %d ", s.libLocal, s.lastLocal, myVersion)

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

func (s *synchronizes) sendBlockHeaderReq(begin uint64, end uint64) {
	header := blockHeaderReq{Begin: begin, End: end}

	data, err := bpl.Marshal(header)
	if err != nil {
		log.Error("protocol sendBlockHeaderReq Marshal number error ")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_HEADER_REQ,
	}

	packet := p2p.Packet{H: head, Data: data}

	peerset := s.getPeers()
	if len(peerset) == 0 {
		log.Error("PROTOCOL sendBlockHeaderReq no peer")
		return
	}

	sort.Sort(peerset)

	//send to three peers which counter of time out is min
	var counter uint16
	for _, info := range peerset {
		if counter >= SYNC_HEADER_BUNDLE {
			break
		}

		if info.LastLib >= end {
			msg := p2p.UniMsgPacket{Index: info.Index,
				P: packet}

			s.set.indexHeader[counter] = info.Index
			log.Debugf("PROTOCOL sendBlockHeaderReq index: %d", s.set.indexHeader[counter])

			p2p.Runner.SendUnicast(msg)

			counter++
		}
	}

}

func (s *synchronizes) syncBundleBlock() {
	if s.set.end < s.set.begin {
		log.Errorf("PROTOCOL syncBundleBlock end %d smaller than begin %d", s.set.end, s.set.begin)
		return
	}

	var numbers []uint64
	lenght := s.set.end + 1 - s.set.begin
	for i := 0; i < int(lenght) && i < SYNC_BLOCK_BUNDLE; i++ {
		if s.set.blocks[i] == nil {
			numbers = append(numbers, s.set.begin+uint64(i))
		}
	}

	if len(numbers) == 0 {
		log.Errorf("PROTOCOL syncBundleBlock sync bundle block finish, wait for send up")
		return
	}

	peerset := s.getPeers()
	if len(peerset) == 0 {
		log.Errorf("PROTOCOL syncBundleBlock no peer")
		return
	}

	sort.Sort(peerset)

	//filter half of time out peer
	avglen := peerset.Len()
	if avglen%2 == 0 {
		avglen = avglen/2 - 1
	} else {
		avglen = avglen / 2
	}

	avg := peerset[avglen].syncTimeoutCounter
	var setlib syncsetlib
	var j int
	for j = 0; j < len(peerset); j++ {
		info := peerset[j]
		if info.syncTimeoutCounter > avg {
			break
		} else {
			setlib = append(setlib, info)
		}
	}

	sort.Sort(setlib)

	if setlib[len(setlib)-1].LastLib < numbers[len(numbers)-1] {
		//can't filter peers, because timeout peer lib is bigger
		for ; j < len(peerset); j++ {
			setlib = append(setlib, peerset[j])
		}

		sort.Sort(setlib)
	}

	if setlib[len(setlib)-1].LastLib < numbers[len(numbers)-1] {
		log.Errorf("PROTOCOL syncBundleBlock peers max lib is smaller than number")
		return
	}

	k := 0
	for _, number := range numbers {
		if k == len(setlib) {
			k = 0
		}

		for k < len(setlib) {
			if setlib[k].LastLib >= number {
				s.sendBlockReq(setlib[k].Index, number, BLOCK_REQ)
				s.set.indexs[number-s.set.begin] = setlib[k].Index
				k++
				break
			} else {
				k++
			}
		}
	}

	s.set.syncBlockTimer.Reset(TIMER_BLOCK_SYNC * time.Second)
}

func (s *synchronizes) sendBlockReq(index uint16, number uint64, ptype uint16) {
	req := syncReq{Number: number, Version: version.GetVersionNumByBlockNum(number)}
	data, err := bpl.Marshal(req)
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

func (s *synchronizes) setSyncStateCheck() {
	if s.set.state != SET_SYNC_BLOCK {
		return
	}

	sync := s.set.setSyncStateJudge()
	if sync {
		s.sendupBundleBlock()
	} else {
		s.checkSyncBlockTimeoutPeers()
		s.syncBundleBlock()
	}
}

func (s *synchronizes) checkSyncHeaderTimeoutPeer() {
	log.Debugf("PROTOCOL index %d sync head time out", s.set.indexHeader)
	for i := 0; i < SYNC_HEADER_BUNDLE; i++ {
		if s.set.indexHeader[i] != 0 {
			s.recordPeerSyncTimeout(s.set.indexHeader[i])
		}
	}
}

func (s *synchronizes) checkSyncBlockTimeoutPeers() {
	lenght := s.set.end + 1 - s.set.begin
	for i := 0; i < int(lenght) && i < SYNC_BLOCK_BUNDLE; i++ {
		if s.set.blocks[i] == nil {
			log.Debugf("PROTOCOL index %d sync block time out", s.set.indexs[i])
			s.recordPeerSyncTimeout(s.set.indexs[i])
		}
	}
}

func (s *synchronizes) sendupBundleBlock() {
	log.Debugf("PROTOCOL sync bundle of block finish")

	if s.set.end < s.set.begin {
		return
	}

	if s.set.begin <= s.libLocal {
		log.Errorf("PROTOCOL sendupBundleBlock lib local is change bigger, wait next time")
		s.set.reset()
		return
	}

	j := 0
	for i := s.set.begin; i <= s.set.end; i++ {
		if s.sendupBlock(s.set.blocks[j]) != berr.ErrNoError {
			s.set.blocks[j] = nil
			return
		}
		j++
	}

	s.libLocal = s.set.end
	s.lastLocal = s.set.end
	log.Debugf("PROTOCOL update local lib and number: %d", s.libLocal)

	s.set.reset()

	if s.libLocal < s.libRemote {
		s.syncBlockHeader()
	} else {
		log.Debugf("PROTOCOL sync finish reset peer sync counter")
		s.resetPeerSyncTimeout()
	}
}

func (s *synchronizes) sendupBlock(block *types.Block) berr.ErrCode {

	start := common.MeasureStart()
	log.Debugf("PROTOCOL send up block :%d", block.Header.Number)

	for i := 0; i < 5; i++ {
		msg := &message.ReceiveBlock{Block: block}

		result, err := s.chain.RequestFuture(msg, 500*time.Millisecond).Result()
		if err != nil {
			log.Errorf("PROTOCOL send block request error:%s", err)
			time.Sleep(10 * time.Millisecond)
			continue
		}

		rsp := result.(*message.ReceiveBlockResp)

		if rsp.ErrorNo != berr.ErrNoError {
			log.Errorf("PROTOCOL block insert error: %d", rsp.ErrorNo)
		}
		blocknumber := s.chainIf.HeadBlockNum()
		s.updateLocalNumber(blocknumber)
		libnumber := s.chainIf.LastConsensusBlockNum()
		s.updateLocalLib(libnumber)
		log.Debugf("PROTOCOL elapsed time 1 %d ", common.Elapsed(start))

		return rsp.ErrorNo
	}

	log.Error("PROTOCOL block insert timeout with five times")

	log.Debugf("PROTOCOL elapsed time 2 %d", common.Elapsed(start))

	return berr.ErrNoError
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
	buf, err := update.block.Marshal()
	if err != nil {
		log.Errorf("PROTOCOL block send marshal error")
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
	buf, err := bpl.Marshal(update.block.Header)
	if err != nil {
		log.Errorf("PROTOCOL block send marshal error")
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
		filter = append(peers[0:k], peers[k+1+number:]...)
	} else if k+1+number == total {
		filter = append(peers[0:k])
	} else {
		filter = append(peers[k+1+number-total : k])
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
		log.Debugf("PROTOCOL catchup counter error")
		s.c.catchupReset()
	} else {
		log.Debugf("PROTOCOL catchup resend get block: %d", s.c.current)
		s.sendBlockReq(s.c.index, s.c.current, BLOCK_CATCH_REQUEST)
	}
}

func (s *synchronizes) catchupRecvBlock(update *blockUpdate) {
	if s.c.index != update.index {
		return
	}

	if update.block == nil ||
		update.block.Header == nil {
		log.Errorf("PROTOCOL catchup with peer index:%d , block:%d finish", s.c.index, s.c.current-1)
		s.c.catchupReset()
		return
	}

	if update.block.Header.Number != s.c.current {
		log.Errorf("PROTOCOL catch up recevie wrong block numbe:%d", update.block.Header.Number)
		return
	}

	result := s.sendupBlock(update.block)
	if result == berr.ErrNoError {
		s.c.current++
		s.c.counter = 0

		s.lastLocal = update.block.Header.Number
		log.Debugf("PROTOCOL catchup update local number: %d", s.lastLocal)
		log.Debugf("PROTOCOL catchup get next block: %d", s.c.current)

		s.sendBlockReq(s.c.index, s.c.current, BLOCK_CATCH_REQUEST)
	} else if result == berr.ErrBlockInsertErrorNotLinked {
		if s.c.current > s.c.begin {
			log.Errorf("PROTOCOL catchup no link, start catchup from begin: %d", s.lastLocal)
			s.c.current = s.c.begin
			s.c.counter = 0
			s.sendBlockReq(s.c.index, s.c.current, BLOCK_CATCH_REQUEST)
		} else if s.c.current == s.c.begin && s.c.begin > s.libLocal+1 {
			log.Errorf("PROTOCOL catchup no link, start catchup from lib: %d", s.libLocal)
			s.c.begin = s.libLocal + 1
			s.c.current = s.c.begin
			s.c.counter = 0
			s.sendBlockReq(s.c.index, s.c.current, BLOCK_CATCH_REQUEST)
		} else {
			log.Errorf("PROTOCOL catchup with peer:%d error", s.c.index)
			s.c.catchupReset()
		}
	} else {
		log.Errorf("PROTOCOL catchup with peer error, reset and wait next time")
		s.c.catchupReset()
	}

}

func (s *synchronizes) catchupWithPeer(index uint16, number uint64) {
	log.Errorf("PROTOCOL catch up with peer:%d, number:%d,s.c.state%d,s.c.index%d", index, number, s.c.state, s.c.index)

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
			log.Debugf("PROTOCOL catch up doing with peer:%d, number:%d,state%d,s.c.index%d", index, number, s.c.state, s.c.index)
			s.sendBlockReq(index, s.c.current, BLOCK_CATCH_REQUEST)
			return
		} else {
			log.Debugf("PROTOCOL catch up doing for extra with peer:%d, number:%d,state%d,s.c.index%d", index, number, s.c.state, s.c.index)
			s.sendBlockReq(index, s.c.current, BLOCK_CATCH_REQUEST)
		}
	} else {
		log.Debugf("PROTOCOL catchupWithPeer wrong state %d", s.c.state)
		panic("PROTOCOL wrong state")
		return
	}
}

type syncSet struct {
	syncheaderc     chan *blockHeaderRsp
	syncblockc      chan *blockUpdate
	syncHeaderTimer *time.Timer
	syncBlockTimer  *time.Timer
	beginc          chan uint64
	endc            chan uint64

	indexHeader [SYNC_HEADER_BUNDLE]uint16
	headers     [SYNC_BLOCK_BUNDLE]*types.Header
	indexs      [SYNC_BLOCK_BUNDLE]uint16
	blocks      [SYNC_BLOCK_BUNDLE]*types.Block

	begin uint64
	end   uint64

	state uint16
}

func makeSyncSet() *syncSet {
	return &syncSet{
		syncheaderc: make(chan *blockHeaderRsp),
		syncblockc:  make(chan *blockUpdate),
		beginc:      make(chan uint64),
		endc:        make(chan uint64),
		state:       SET_SYNC_NULL}
}

func (set *syncSet) recvBlockHeader(rsp *blockHeaderRsp) bool {
	if set.state != SET_SYNC_HEADER {
		log.Debug("PROTOCOL recvBlockHeader state error, could have receive ack")
		return false
	}

	if set.end < set.begin {
		log.Errorf("PROTOCOL recvBlockHeader set end %d small than begin %d", set.end, set.begin)
		return false
	}

	if uint64(len(rsp.set)) != (set.end + 1 - set.begin) {
		log.Errorf("PROTOCOL recvBlockHeader rsp length error")
		return false
	}

	check := false
	j := 0
	for i := set.begin; i <= set.end; i++ {
		if rsp.set[j].GetNumber() != i {
			log.Errorf("PROTOCOL recvBlockHeader rsp info error number:%d", rsp.set[j].GetNumber())
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
func (set *syncSet) endcCheck(number uint64) {
	if set.state == SET_SYNC_NULL {
		log.Debugf("PROTOCOL sync status null")
		return
	}

	//remote lib change small , we should reset and wait for sync judge
	if number < set.end {
		log.Debugf("PROTOCOL endcCheck reset end: %d, lib: %d", set.end, number)
		set.reset()
	}
}

//begincCheck local lib change bigger when produce a block in p2p sync state
func (set *syncSet) begincCheck(number uint64) {
	if set.state == SET_SYNC_NULL {
		log.Debugf("PROTOCOL sync status null")
		return
	}

	//local lib change bigger , we should reset and wait for sync judge
	if number >= set.begin {
		log.Debugf("PROTOCOL begincCheck reset begin: %d, lib: %d", set.begin, number)
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

func (set *syncSet) resetHeaderIndex() {
	for i := 0; i < SYNC_HEADER_BUNDLE; i++ {
		set.indexHeader[i] = 0
	}
}

func (set *syncSet) resetIndex() {
	for i := 0; i < SYNC_BLOCK_BUNDLE; i++ {
		set.indexs[i] = 0
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
	set.resetHeaderIndex()
	set.resetHeader()
	set.resetIndex()
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
	begin   uint64
	current uint64
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
