package transaction

import (
	"encoding/json"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
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

	msg := p2p.MsgPacket{Index: peers,
		P: packet}

	if broadcast {
		p2p.Runner.SendBroadcast(msg)
	} else {
		p2p.Runner.SendUnicast(msg)
	}
}

func (t *Transaction) processTrxInfo(index uint16, p *p2p.Packet) {
	var trx types.Transaction
	err := json.Unmarshal(p.Data, &trx)

	msg := message.ReceiveTrx{Trx: &trx}
	if err != nil {
		t.actor.Tell(&msg)
	}
}
