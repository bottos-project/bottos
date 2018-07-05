package transaction

import (
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
	"time"
)

type Transaction struct {
	actor *actor.PID
}

func MakeTransaction() *Transaction {
	return &Transaction{}
}

func (t *Transaction) Start() {

}

func (t *Transaction) SetActor(tid *actor.PID) {
	t.actor = tid
}

func (t *Transaction) Dispatch(index uint16, p *p2p.Packet) {
	switch p.H.PacketType {
	case TRX_UPDATE:
		t.processTrxInfo(index, p)
	}
}

func (t *Transaction) SendNewTrx(notify *message.NotifyTrx) {
	t.sendPacket(true, notify.Trx, nil)
}

func (t *Transaction) sendPacket(broadcast bool, data interface{}, peers []uint16) {
	buf, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Transaction send marshal error")
	}

	head := p2p.Head{ProtocolType: pcommon.TRX_PACKET,
		PacketType: TRX_UPDATE,
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

func (t *Transaction) processTrxInfo(index uint16, p *p2p.Packet) {
	var trx types.Transaction

	err := json.Unmarshal(p.Data, &trx)
	if err != nil {
		log.Errorf("processTrxInfo Unmarshal error")
		return
	}

	if t.sendupTrx(&trx) {
		t.sendPacket(true, &trx, []uint16{index})
	}
}

func (t *Transaction) sendupTrx(trx *types.Transaction) bool {
	for i := 0; i < 5; i++ {
		msg := &message.ReceiveTrx{Trx: trx}
		handlerErr, err := t.actor.RequestFuture(msg, 500*time.Millisecond).Result()
		if err != nil {
			log.Errorf("send block request error:%s", err)
			time.Sleep(10000)
			continue
		}

		if handlerErr == bottosErr.ErrNoError {
			log.Errorf("send block request response error:%d", handlerErr)
			return true
		}

		return false
	}

	return false
}
