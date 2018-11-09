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

package protocol

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/context"
	"github.com/bottos-project/bottos/p2p"
	"github.com/bottos-project/bottos/protocol/block"
	"github.com/bottos-project/bottos/protocol/common"
	"github.com/bottos-project/bottos/protocol/consensus"
	"github.com/bottos-project/bottos/protocol/discover"
	"github.com/bottos-project/bottos/protocol/transaction"
	log "github.com/cihub/seelog"
	"net"
)

type protocol struct {
	d *discover.Discover
	t *transaction.Transaction
	b *block.Block
	c *consensus.Consensus
}

func MakeProtocol(config *config.Parameter, chain chain.BlockChainInterface) context.ProtocolInstance {
	runner := p2p.MakeP2PServer(config)

	p := &protocol{
		d: discover.MakeDiscover(config),
		t: transaction.MakeTransaction(),
		b: block.MakeBlock(chain, true), //we should get node type here
		c: consensus.MakeConsensus(),
	}

	sendup := func(index uint16, packet *p2p.Packet) {
		if packet.H.ProtocolType == common.P2P_PACKET {
			p.d.Dispatch(index, packet)
		} else if packet.H.ProtocolType == common.TRX_PACKET {
			p.t.Dispatch(index, packet)
		} else if packet.H.ProtocolType == common.BLOCK_PACKET {
			p.b.Dispatch(index, packet)
		} else if packet.H.ProtocolType == common.CONSENSUS_PACKET {
			p.c.Dispatch(index, packet)
		} else {
			log.Errorf("wrong packet type")
		}
	}

	p.d.SetSendupCallback(sendup)

	conn := func(conn net.Conn) {
		p.d.NewConnCb(conn, sendup)
	}

	runner.SetCallback(conn)

	return p
}

func (p *protocol) Start() {
	p2p.Runner.Start()
	p.d.Start()
	p.c.Start()
	p.b.Start()
	p.t.Start()
}

func (p *protocol) GetBlockSyncState() bool {
	return p.b.GetSyncState()
}

func (p *protocol) SetChainActor(tpid *actor.PID) {
	p.b.SetActor(tpid)
}

func (p *protocol) SetTrxPreHandleActor(tpid *actor.PID) {
	p.t.SetActor(tpid)
}

func (p *protocol) SetConsensusActor(tpid *actor.PID) {
	p.c.SetActor(tpid)
}

func (p *protocol) SendNewTrx(notify *message.NotifyTrx) {
	p.t.SendNewTrx(notify)
}

func (p *protocol) SendNewBlock(notify *message.NotifyBlock) {
	p.b.SendNewBlock(notify)
	p.b.UpdateHeadNumber()
	p.b.UpdateNumber()
}
