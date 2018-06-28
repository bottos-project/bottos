package protocal

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/p2p"
	"github.com/bottos-project/bottos/protocal/block"
	"github.com/bottos-project/bottos/protocal/common"
	"github.com/bottos-project/bottos/protocal/consensus"
	"github.com/bottos-project/bottos/protocal/discover"
	"github.com/bottos-project/bottos/protocal/transaction"
	log "github.com/cihub/seelog"
	"net"
)

type protocal struct {
	d *discover.Discover
	t *transaction.Transaction
	b *block.Block
	c *consensus.Consensus
}

func MakeProtocal(config *config.Parameter, chain chain.BlockChainInterface) ProtocalInstance {
	runner := p2p.MakeP2PServer(config)

	p := &protocal{
		d: discover.MakeDiscover(config),
		t: transaction.MakeTransaction(),
		b: block.MakeBlock(chain),
		c: consensus.MakeConsensus(),
	}

	sendup := func(index uint16, packet *p2p.Packet) {
		if packet.H.ProtocalType == common.P2P_PACKET {
			p.d.Dispatch(index, packet)
		} else if packet.H.ProtocalType == common.TRX_PACKET {
			p.t.Dispatch(index, packet)
		} else if packet.H.ProtocalType == common.BLOCK_PACKET {
			p.b.Dispatch(index, packet)
		} else if packet.H.ProtocalType == common.CONSENSUS_PACKET {
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

func (p *protocal) Start() {
	p2p.Runner.Start()
	p.d.Start()
	p.c.Start()
	p.b.Start()
	p.t.Start()
}

func (p *protocal) GetBlockSyncState() bool {
	return p.b.GetSyncState()
}

func (p *protocal) Send(ptype uint16, broadcast bool, data interface{}, peers []uint16) {
	if ptype == common.TRX_PACKET {
		p.t.Send(broadcast, data, peers)
	} else if ptype == common.BLOCK_PACKET {
		p.b.Send(broadcast, data, peers)
	} else if ptype == common.CONSENSUS_PACKET {
		p.c.Send(broadcast, data, peers)
	} else {
		log.Errorf("wrong packet type")
	}
}

func (p *protocal) SetChainActor(tpid *actor.PID) {
	p.b.SetActor(tpid)
}

func (p *protocal) SetTrxActor(tpid *actor.PID) {
	p.t.SetActor(tpid)
}

func (p *protocal) SetProducerActor(tpid *actor.PID) {
	p.c.SetActor(tpid)
}
