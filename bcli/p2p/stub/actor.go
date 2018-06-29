package stub

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
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
	switch msg := context.Message().(type) {
	case *message.ReceiveBlock:
		n.HandleReceiveBlock(context, msg)
	}
}

func (n *DumActor) handleSystemMsg(context actor.Context) {

}

func (n *DumActor) HandleReceiveBlock(ctx actor.Context, req *message.ReceiveBlock) {
	rsp := &message.ReceiveBlockResp{
		BlockNum: req.Block.GetNumber(),
		ErrorNo:  0,
	}

	ctx.Sender().Request(rsp, ctx.Self())

}
