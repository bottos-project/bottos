package discover

import (
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
	"sync/atomic"
	"time"
)

type keeplive struct {
	counter [MAX_PEER_COUNT + 1]int32

	c *candidates
}

const (
	//TIME_KEEP_LIVE ping/pong timer, second
	TIMER_KEEP_LIVE = 20
	//TIMER_CHECK time out second
	TIMER_CHECK = 200
)

func makeKeeplive(c *candidates) *keeplive {
	return &keeplive{c: c}
}

func (k *keeplive) start() {
	for i := 0; i < MAX_PEER_COUNT; i++ {
		k.counter[i] = -1
	}

	go k.keepliveTimer()
	go k.checkTimer()
}

func (k *keeplive) initCounter(index uint16) {
	atomic.StoreInt32(&k.counter[index], 0)
}

func (k *keeplive) keepliveTimer() {
	log.Debug("keepliveTimer")

	keep := time.NewTimer(TIMER_KEEP_LIVE * time.Second)

	defer func() {
		log.Debug("keepliveTimer stop")
		keep.Stop()
	}()

	for {
		select {
		case <-keep.C:
			k.sendPing()
			keep.Reset(TIMER_KEEP_LIVE * time.Second)
		}
	}
}

func (k *keeplive) checkTimer() {
	log.Debug("checkTimer")

	check := time.NewTimer(TIMER_CHECK * time.Second)

	defer func() {
		log.Debug("checkTimer stop")
		check.Stop()
	}()

	for {
		select {
		case <-check.C:
			k.checkPeer()
			check.Reset(TIMER_CHECK * time.Second)
		}
	}
}

func (k *keeplive) checkPeer() {
	for i := 0; i < MAX_PEER_COUNT; i++ {
		if k.counter[i] != -1 {
			if k.counter[i] == 0 {
				atomic.StoreInt32(&k.counter[i], -1)
				if p2p.Runner.DelPeer(uint16(i)) {
					k.c.pushPeerIndex(uint16(i))
				}
			} else {
				atomic.StoreInt32(&k.counter[i], 0)
			}
		}
	}
}

func (k *keeplive) processPing(index uint16, date []byte) {
	k.sendPong(index)
}

func (k *keeplive) processPong(index uint16, date []byte) {
	k.counterPeer(index)
}

func (k *keeplive) counterPeer(index uint16) {
	atomic.AddInt32(&k.counter[index], 1)
}

func (k *keeplive) sendPing() {
	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_PING,
	}

	packet := p2p.Packet{H: head}

	ping := p2p.MsgPacket{
		Index: nil,
		P:     packet,
	}

	p2p.Runner.SendBroadcast(ping)

}

func (k *keeplive) sendPong(index uint16) {
	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_PONG,
	}

	packet := p2p.Packet{H: head}

	pong := p2p.MsgPacket{
		Index: []uint16{index},
		P:     packet,
	}

	p2p.Runner.SendUnicast(pong)

}
