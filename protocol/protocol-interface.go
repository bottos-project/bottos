package protocol

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type ProtocolInstance interface {
	ProtocolInterface
	Start()
	Send(ptype uint16, broadcast bool, data interface{}, peers []uint16)
	SetChainActor(tpid *actor.PID)
	SetTrxActor(tpid *actor.PID)
	SetProducerActor(tpid *actor.PID)
}

type ProtocolInterface interface {
	GetBlockSyncState() bool
}
