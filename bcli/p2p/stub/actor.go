package stub

import (
	"github.com/AsynkronIT/protoactor-go/actor"
)

type DumActor struct {
}

var dumactor *DumActor

func NewDumActor() *actor.PID {

	dumactor = &DumActor{}

	props := actor.FromProducer(func() actor.Actor { return dumactor })

	pid, err := actor.SpawnNamed(props, "DumActor")
	if err == nil {
		return pid
	} else {
		return nil
	}

	return nil
}

func (n *DumActor) Receive(context actor.Context) {
	n.handleSystemMsg(context)
}

func (n *DumActor) handleSystemMsg(context actor.Context) {

}
