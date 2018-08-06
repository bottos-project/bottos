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

	t.sendupTrx(&trx)

}

func (t *Transaction) sendupTrx(trx *types.Transaction) bool {
	for i := 0; i < 5; i++ {
		msg := &message.ReceiveTrx{Trx: trx}
		handlerErr, err := t.actor.RequestFuture(msg, 500*time.Millisecond).Result()
		if err != nil {
			log.Errorf("send trx request error:%s", err)
			time.Sleep(10 * time.Millisecond)
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
