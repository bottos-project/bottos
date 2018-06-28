package consensus

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/p2p"
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

func (c *Consensus) Dispatch(index uint16, p *p2p.Packet) {

}

func (c *Consensus) Send(broadcast bool, m interface{}, peers []uint16) {

}

func (c *Consensus) Start() {

}
