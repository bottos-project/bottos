package protocol

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
)

type ProtocolInstance interface {
	ProtocolInterface
	Start()
	SetChainActor(tpid *actor.PID)
	SetTrxActor(tpid *actor.PID)
	SetProducerActor(tpid *actor.PID)

	ProcessNewTrx(notify *message.NotifyTrx)
	ProcessNewBlock(notify *message.NotifyBlock)
}

type ProtocolInterface interface {
	GetBlockSyncState() bool
}
