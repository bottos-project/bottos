package block

import (
	"bytes"
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocal/common"
	log "github.com/cihub/seelog"
	"sort"
	"time"
)

const (
	TIMER_FAST_SYNC_LAST_BLOCK_NUMBER   = 1
	TIMER_NORMAL_SYNC_LAST_BLOCK_NUMBER = 4
	//SYNC_LAST_BLOCK_NUMBER_COUNTER counter of no last block number response to set a peer expired
	SYNC_LAST_BLOCK_NUMBER_COUNTER = 10

	TIMER_SYNC_STATE_CHECK = 7

	TIMER_HEADER_SYNC = 2
	TIMER_BLOCK_SYNC  = 2

	SYNC_BLOCK_BUNDLE = 10
)

const (
	SET_SYNC_NULL   = 0
	SET_SYNC_HEADER = 1
	SET_SYNC_BLOCK  = 2
	SET_SYNC_FINISH = 3
)

type blockNumberInfo struct {
	index uint16
	last  uint32

	counter int16
}

type syncset []blockNumberInfo

func (s syncset) Len() int {
	return len(s)
}

func (s syncset) Less(i, j int) bool {
	return s[i].last > s[j].last
}
func (s syncset) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type synchronizes struct {
	peers map[uint16]*blockNumberInfo

	numberc chan *blockNumberInfo
	sendc   chan uint32
	recvc   chan *blockUpdate

	syncc     chan *blockHeaderRsp
	remoteupc chan uint32

	lastLocal  uint32
	lastRemote uint32
	//synchronize status , true is synchronized,  false is unsynchronized
	state bool
	//have synchronized one time or not when start up
	once bool

	syncHeaderTimer *time.Timer
	syncBlockTimer  *time.Timer

	set   *blockset
	chain *actor.PID
}

func MakeSynchronizes() *synchronizes {
	return &synchronizes{
		peers:     make(map[uint16]*blockNumberInfo),
		numberc:   make(chan *blockNumberInfo, 10),
		sendc:     make(chan uint32),
		recvc:     make(chan *blockUpdate),
		syncc:     make(chan *blockHeaderRsp),
		remoteupc: make(chan uint32),
		state:     false,
		once:      false,
		set:       makeBlockSet(),
	}
}

func (s *synchronizes) SetActor(tid *actor.PID) {
	s.chain = tid
}

func (s *synchronizes) start() {
	go s.syncBlockNumberTimer()
	go s.recvRoutine()
	go s.syncRoutine()
}

func (s *synchronizes) syncBlockNumberTimer() {
	log.Debug("syncBlockNumberTimer start")

	syncTimer := time.NewTimer(TIMER_FAST_SYNC_LAST_BLOCK_NUMBER * time.Second)

	defer func() {
		log.Debug("syncBlockNumberTimer stop")
		syncTimer.Stop()
	}()

	for {
		select {
		case <-syncTimer.C:
			s.sendLastBlockNumberReq()
			if s.state {
				syncTimer.Reset(TIMER_NORMAL_SYNC_LAST_BLOCK_NUMBER * time.Second)
			} else {
				syncTimer.Reset(TIMER_FAST_SYNC_LAST_BLOCK_NUMBER * time.Second)
			}
		}
	}
}

func (s *synchronizes) recvRoutine() {
	checkTimer := time.NewTimer(TIMER_SYNC_STATE_CHECK * time.Second)

	for {
		select {
		case info := <-s.numberc:
			s.recvPeerBlockNumberInfo(info)
		case number := <-s.sendc:
			s.sendUpdateLocalNumber(number)
		case update := <-s.recvc:
			s.recvBlock(update)
		case <-checkTimer.C:
			s.syncStateCheck()
			checkTimer.Reset(TIMER_SYNC_STATE_CHECK * time.Second)
		}
	}
}

func (s *synchronizes) syncRoutine() {
	s.syncHeaderTimer = time.NewTimer(TIMER_HEADER_SYNC * time.Second)
	s.syncBlockTimer = time.NewTimer(TIMER_BLOCK_SYNC * time.Second)

	for {
		select {
		case rsp := <-s.syncc:
			if s.set.recvBlockHeader(rsp) {
				s.syncHeaderTimer.Stop()
				s.syncBundleBlock()
			}
		case number := <-s.remoteupc:
			s.set.updateRemoteNumber(number)
		case <-s.syncHeaderTimer.C:
			if s.set.state == SET_SYNC_HEADER {
				s.syncBlockHeader()
			}
		case <-s.syncBlockTimer.C:
			s.setSyncStateCheckTimer()
		}
	}
}

func (s *synchronizes) recvPeerBlockNumberInfo(info *blockNumberInfo) {
	info.counter = 0

	peer := s.peers[info.index]
	if peer == nil {
		s.peers[info.index] = info
	} else {
		if peer.last >= info.last {
			peer.counter = 0
			return
		} else {
			peer.last = info.last
			peer.counter = 0
		}
	}

	s.updateRemoteNumber(info.last, false)
}

//sendNumber produce a block and set out, so update local block number
func (s *synchronizes) sendUpdateLocalNumber(number uint32) {
	s.updateLocalNumber(number)
}

func (s *synchronizes) recvBlock(update *blockUpdate) {
	number := update.block.GetNumber()

	if number <= s.lastLocal {
		log.Debugf("drop block: %d is smaller than local", number)
		return
	}

	if (s.state && number == s.lastLocal+1) ||
		(!s.state && s.set.state == SET_SYNC_NULL && s.lastLocal+1 == number) {
		s.updateLocalNumber(number)
		s.updateRemoteNumber(number, false)

		msg := message.ReceiveBlock{
			Block: update.block,
		}
		s.chain.Tell(msg)
		return
	}

	info := blockNumberInfo{index: update.index, last: number}
	s.recvPeerBlockNumberInfo(&info)

	if s.state == false {
		if s.set.state == SET_SYNC_NULL {
			s.syncStateJudge()
			log.Debugf("drop block: %d when in sync null status", number)
		} else if s.set.state == SET_SYNC_BLOCK {
			setsync := s.set.recvBlock(update.block)
			if setsync {
				s.sendupBundleBlock()
			}
		} else {
			log.Debugf("drop block: %d when in sync status", number)
		}
	} else {
		s.syncStateJudge()
		log.Info("drop block %d to begin sync", number)

	}
}

func (s *synchronizes) syncStateCheck() {
	var remote uint32

	for key, info := range s.peers {
		info.counter++

		if info.counter >= SYNC_LAST_BLOCK_NUMBER_COUNTER {
			delete(s.peers, key)
		}

		if info.last > remote {
			remote = info.last
		}
	}

	//remote block number be smaller, wo should reset all
	if remote < s.lastRemote {
		log.Errorf("syncStateCheck remote block number change smaller")
		if remote > 0 {
			s.updateRemoteNumber(remote, true)

		}
	} else if remote > s.lastRemote {
		log.Errorf("syncStateCheck remote block number change bigger")
		s.updateRemoteNumber(remote, false)
	}

	s.syncStateJudge()
}

func (s *synchronizes) syncStateJudge() {
	if s.lastLocal < s.lastRemote {
		if !s.once {
			s.syncBlockHeader()
			s.once = true
		} else if s.state {
			log.Debugf("syncStateJudge not sync")
			s.state = false
			s.syncBlockHeader()
		}
	} else {
		s.state = true
	}
}

func (s *synchronizes) updateLocalNumber(number uint32) {
	if number <= s.lastLocal {
		return
	}

	log.Debugf("update local block number:%d", number)
	s.lastLocal = number
}

func (s *synchronizes) updateRemoteNumber(number uint32, force bool) {
	if !force && number <= s.lastRemote {
		return
	}

	log.Debugf("peer max block number:%d", number)
	s.lastRemote = number
}

func (s *synchronizes) sendLastBlockNumberReq() {
	head := p2p.Head{ProtocalType: pcommon.BLOCK_PACKET,
		PacketType: LAST_BLOCK_NUMBER_REQ,
	}

	packet := p2p.Packet{H: head}

	msg := p2p.MsgPacket{Index: nil,
		P: packet}

	p2p.Runner.SendBroadcast(msg)
}

func (s *synchronizes) sendLastBlockNumberRsp(index uint16) {
	rsp := s.lastLocal

	data, err := json.Marshal(rsp)
	if err != nil {
		log.Error("sendGetLastRsp Marshal data error ")
		return
	}

	head := p2p.Head{ProtocalType: pcommon.BLOCK_PACKET,
		PacketType: LAST_BLOCK_NUMBER_RSP,
	}

	packet := p2p.Packet{H: head, Data: data}

	msg := p2p.MsgPacket{Index: []uint16{index},
		P: packet}

	p2p.Runner.SendUnicast(msg)
}

func (s *synchronizes) syncBlockHeader() {
	if s.lastRemote <= s.lastLocal {
		return
	}

	if s.lastLocal+SYNC_BLOCK_BUNDLE > s.lastRemote {
		s.set.begin = s.lastLocal + 1
		s.set.end = s.lastRemote
	} else {
		s.set.begin = s.lastLocal + 1
		s.set.end = s.lastLocal + SYNC_BLOCK_BUNDLE

	}

	s.set.state = SET_SYNC_HEADER
	s.sendBlockHeaderReq(s.set.begin, s.set.end)

	s.syncHeaderTimer.Reset(TIMER_HEADER_SYNC * time.Second)
}

func (s *synchronizes) sendBlockHeaderReq(begin uint32, end uint32) {
	header := blockHeaderReq{Begin: begin, End: end}

	data, err := json.Marshal(header)
	if err != nil {
		log.Error("sendBlockHeaderReq Marshal number error ")
		return
	}

	head := p2p.Head{ProtocalType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_HEADER_REQ,
	}

	packet := p2p.Packet{H: head, Data: data}

	for _, info := range s.peers {
		if info.last >= end {
			msg := p2p.MsgPacket{Index: []uint16{info.index},
				P: packet}

			p2p.Runner.SendUnicast(msg)
			break
		}
	}
}

func (s *synchronizes) syncBundleBlock() {
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

	var peerset syncset
	for _, info := range s.peers {
		peerset = append(peerset, *info)
	}

	sort.Sort(peerset)

	i := 0
	for _, number := range numbers {
		if i == len(peerset) {
			i = 0
		}

		for i < len(peerset) {
			if peerset[i].last >= number {
				s.sendBlockReq(peerset[i].index, number)
				i++
				break
			} else {
				i++
			}
		}
	}

	s.syncBlockTimer.Reset(TIMER_BLOCK_SYNC * time.Second)
}

func (s *synchronizes) sendBlockReq(index uint16, number uint32) {

	data, err := json.Marshal(number)
	if err != nil {
		log.Error("sendGetBlock Marshal number error ")
		return
	}

	head := p2p.Head{ProtocalType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_REQ,
	}

	packet := p2p.Packet{H: head, Data: data}

	msg := p2p.MsgPacket{Index: []uint16{index},
		P: packet}

	log.Debugf("sendBlockReq block：%d", number)
	p2p.Runner.SendUnicast(msg)
}

func (s *synchronizes) setSyncStateCheckTimer() {
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
	log.Debugf("sync bundle of block finish")

	s.set.state = SET_SYNC_FINISH

	j := 0
	for i := s.set.begin; i <= s.set.end; i++ {
		msg := message.ReceiveBlock{Block: s.set.blocks[j]}
		s.chain.Tell(msg)
		j++
	}

	s.lastLocal = s.set.end

	s.set.state = SET_SYNC_NULL
	s.set.end = 0
	s.set.begin = 0
	s.set.resetHeader()
	s.set.resetBlock()

	if s.lastLocal < s.lastRemote {
		s.syncBlockHeader()
	}
}

type blockset struct {
	headers [SYNC_BLOCK_BUNDLE]*types.Header
	blocks  [SYNC_BLOCK_BUNDLE]*types.Block

	begin uint32
	end   uint32

	state uint16
}

func makeBlockSet() *blockset {
	return &blockset{state: SET_SYNC_NULL}
}

func (set *blockset) recvBlockHeader(rsp *blockHeaderRsp) bool {
	if set.state != SET_SYNC_HEADER {
		log.Errorf("recvBlockHeader state error")
		return false
	}

	if uint32(len(rsp.set)) != (set.end + 1 - set.begin) {
		log.Errorf("recvBlockHeader rsp length error")
		return false
	}

	check := true
	j := 0
	for i := set.begin; i <= set.end; i++ {
		if rsp.set[j].GetNumber() != i {
			log.Errorf("recvBlockHeader rsp info error number:%d", rsp.set[j].GetNumber())
			check = false
			break
		}

		set.headers[j] = &rsp.set[j]
		j++
	}

	if !check {
		set.resetHeader()
		return false
	}

	set.state = SET_SYNC_BLOCK
	return true
}

func (set *blockset) recvBlock(block *types.Block) bool {
	log.Infof("recvBlock block:%d", block.GetNumber())

	if block.GetNumber() > set.end {
		log.Infof("drop block bigger than current sync end")
		return false
	}

	for i := 0; i < SYNC_BLOCK_BUNDLE; i++ {
		if set.headers[i] != nil &&
			set.isBlockHeadSame(set.headers[i], block.Header) {
			set.blocks[i] = block
			break
		}
	}

	//chech if set sync status
	return set.setSyncStateJudge()
}

//updateRemoteNumber update peer max block number if some peer is disconnect
func (set *blockset) updateRemoteNumber(number uint32) {
	if set.end > number {
		log.Debugf("update syn set max block number: %d", number)
		set.end = number
	}
}

func (set *blockset) resetHeader() {
	for i := 0; i < SYNC_BLOCK_BUNDLE; i++ {
		set.headers[i] = nil
	}
}

func (set *blockset) resetBlock() {
	for i := 0; i < SYNC_BLOCK_BUNDLE; i++ {
		set.blocks[i] = nil
	}
}

func (set *blockset) isBlockHeadSame(a *types.Header, b *types.Header) bool {
	if a.Number == b.Number &&
		a.Version == b.Version &&
		a.Timestamp == b.Timestamp &&
		bytes.Equal(a.MerkleRoot, b.MerkleRoot) &&
		bytes.Equal(a.PrevBlockHash, b.PrevBlockHash) {
		return true
	} else {
		return false
	}
}

func (set *blockset) setSyncStateJudge() bool {

	lenght := set.end + 1 - set.begin
	for i := 0; i < int(lenght) && i < SYNC_BLOCK_BUNDLE; i++ {
		if set.blocks[i] == nil {
			return false
		}
	}

	return true
}
