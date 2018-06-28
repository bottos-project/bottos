package block

import (
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocal/common"
	log "github.com/cihub/seelog"
	"time"
)

type Block struct {
	actor   *actor.PID
	chainIf chain.BlockChainInterface
	s       *synchronizes
}

const (
	WAIT_TIME              = 30
	WAIT_LOCAL_BLOCK_TIMER = 3
)

func MakeBlock(chain chain.BlockChainInterface) *Block {
	return &Block{s: MakeSynchronizes(),
		chainIf: chain}
}

func (b *Block) SetActor(tid *actor.PID) {
	b.actor = tid
	b.s.SetActor(tid)
}

func (b *Block) Start() {
	go b.waitLastBlockTimer()
}

func (b *Block) waitLastBlockTimer() {
	time.Sleep(WAIT_TIME * time.Second)

	waitTimer := time.NewTimer(WAIT_LOCAL_BLOCK_TIMER * time.Second)

	log.Debug("waitLastBlockTimer start")

	defer func() {
		log.Debug("waitLastBlockTimer stop")
		waitTimer.Stop()
	}()

	for {
		select {
		case <-waitTimer.C:
			blocknumber := b.chainIf.HeadBlockNum()
			log.Debugf("timer local block number:%d", blocknumber)
			b.s.updateLocalNumber(blocknumber)
			b.s.start()
			break
		}
	}

}

func (b *Block) Dispatch(index uint16, p *p2p.Packet) {
	switch p.H.PacketType {
	case BLOCK_REQ:
		b.processBlockReq(index, p.Data)
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
	}
}

func (b *Block) Send(broadcast bool, data interface{}, peers []uint16) {
	block := data.(*types.Block)
	b.s.sendc <- block.GetNumber()

	buf, err := json.Marshal(data)
	if err != nil {
		log.Errorf("block send marshal error")
	}

	head := p2p.Head{ProtocalType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_UPDATE,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	msg := p2p.MsgPacket{Index: peers,
		P: packet}

	if broadcast {
		p2p.Runner.SendBroadcast(msg)
	} else {
		p2p.Runner.SendUnicast(msg)
	}
}

func (b *Block) GetSyncState() bool {
	return b.s.state
}

func (b *Block) processLastBlockNumberReq(index uint16, data []byte) {
	b.s.sendLastBlockNumberRsp(index)
}

func (b *Block) processLastBlockNumberRsp(index uint16, data []byte) {
	var blocknumber uint32
	err := json.Unmarshal(data, &blocknumber)
	if err != nil {
		log.Errorf("processLastBlockNumberRsp Unmarshal error:%s", err)
		return
	}

	info := blockNumberInfo{
		index: index,
		last:  blocknumber,
	}

	b.s.numberc <- &info
}

func (b *Block) processBlockHeaderReq(index uint16, data []byte) {
	var req blockHeaderReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		log.Errorf("processBlockHeaderReq Unmarshal err:%s", err)
		return
	}

	if req.Begin > req.End ||
		req.End-req.Begin >= SYNC_BLOCK_BUNDLE {
		log.Errorf("processBlockHeaderReq wrong lenght")
		return
	}

	var rsp blockHeaderRsp
	j := 0
	for i := req.Begin; i <= req.End; i++ {
		head := b.chainIf.GetHeaderByNumber(i)
		if head == nil {
			log.Errorf("processBlockHeaderReq header not exist")
			return
		} else {
			rsp.set = append(rsp.set, *head)
		}

		j++
	}

	b.sendBlockHeaderRsp(index, &rsp)
}

func (b *Block) processBlockHeaderRsp(index uint16, data []byte) {
	var rsp blockHeaderRsp
	err := json.Unmarshal(data, &rsp.set)
	if err != nil {
		log.Errorf("processBlockInfo Unmarshal error:%s", err)
	}

	b.s.syncc <- &rsp
}

func (b *Block) processBlockReq(index uint16, data []byte) {
	var blocknumber uint32
	err := json.Unmarshal(data, &blocknumber)
	if err != nil {
		log.Errorf("processLastBlockRsp Unmarshal error:%s", err)
		return
	}

	block := b.chainIf.GetBlockByNumber(blocknumber)

	b.sendBlockRsp(index, block)
}

func (b *Block) processBlockInfo(index uint16, data []byte) {
	var block types.Block
	err := json.Unmarshal(data, &block)
	if err != nil {
		log.Errorf("processBlockInfo Unmarshal error:%s", err)
	}

	update := blockUpdate{index: index, block: &block}

	b.s.recvc <- &update
}

func (b *Block) sendBlockHeaderRsp(index uint16, rsp *blockHeaderRsp) {
	data, err := json.Marshal(rsp.set)
	if err != nil {
		log.Error("sendGetBlock Marshal number error ")
		return
	}

	head := p2p.Head{ProtocalType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_HEADER_RSP,
	}

	packet := p2p.Packet{H: head, Data: data}

	msg := p2p.MsgPacket{Index: []uint16{index},
		P: packet}

	p2p.Runner.SendUnicast(msg)
}

func (b *Block) sendBlockRsp(index uint16, block *types.Block) {
	data, err := json.Marshal(block)
	if err != nil {
		log.Error("sendGetBlock Marshal number error ")
		return
	}

	head := p2p.Head{ProtocalType: pcommon.BLOCK_PACKET,
		PacketType: BLOCK_UPDATE,
	}

	packet := p2p.Packet{H: head, Data: data}

	msg := p2p.MsgPacket{Index: []uint16{index},
		P: packet}

	p2p.Runner.SendUnicast(msg)
}
