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

package consensus

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/p2p"
	"github.com/bottos-project/bottos/bpl"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
)

type Consensus struct {
	actor *actor.PID
}

func MakeConsensus() *Consensus {
	return &Consensus{}
}

func (c *Consensus) SetActor(tid *actor.PID) {
	c.actor = tid
}

func (c *Consensus) Start() {

}

func (c *Consensus) Dispatch(index uint16, p *p2p.Packet) {
	switch p.H.PacketType {
	case ConsensusPreVote:
		c.processPrevote(index, p.Data)
	case ConsensusPreCommit:
		c.processPrecommit(index, p.Data)
	case ConsensusCommit:
		c.processCommit(index, p.Data)
	}
}

func (c *Consensus) SendPrevote(notify *message.SendPrevote) {
	buf, err := bpl.Marshal(notify.BlockState)
	if err != nil {
		log.Errorf("protocol block send marshal error")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.CONSENSUS_PACKET,
		PacketType: ConsensusPreVote,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	msg := p2p.BcastMsgPacket{Indexs: nil,
		P: packet}
	p2p.Runner.SendBroadcast(msg)

}

func (c *Consensus) SendPrecommit(notify *message.SendPrecommit) {
	buf, err := bpl.Marshal(notify.BlockState)
	if err != nil {
		log.Errorf("protocol block send marshal error")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.CONSENSUS_PACKET,
		PacketType: ConsensusPreCommit,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	msg := p2p.BcastMsgPacket{Indexs: nil,
		P: packet}
	p2p.Runner.SendBroadcast(msg)

}

func (c *Consensus) SendCommit(notify *message.SendCommit) {
	buf, err := bpl.Marshal(notify.BftHeaderState)
	if err != nil {
		log.Errorf("protocol block send marshal error")
		return
	}

	head := p2p.Head{ProtocolType: pcommon.CONSENSUS_PACKET,
		PacketType: ConsensusCommit,
	}

	packet := p2p.Packet{H: head,
		Data: buf,
	}

	msg := p2p.BcastMsgPacket{Indexs: nil,
		P: packet}
	p2p.Runner.SendBroadcast(msg)
}



