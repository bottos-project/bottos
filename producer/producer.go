package consensus

import (
	"fmt"
	"time"

	"github.com/bottos-project/core/action/actor/producer"
	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/consensus/dpos"
)

//var trxactorPid *actor.PID
type hello struct{ Who string }
type helloActor struct{}

func Working() {
	loop()

}
func loop() {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				reportBlockLoop()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
func IsEligible() bool {
	return true
}
func IsReady() bool {
	slotTime := dpos.GetSlotTime(1)
	fmt.Println(slotTime)
	if slotTime >= common.NowToSeconds() {
		return true
	}
	return false
}
func IsMyTurn() bool {
	return true

}
func reportBlockLoop() {
	if IsEligible() && IsReady() && IsMyTurn() {
		now := time.Now()
		slot := dpos.GetSlotAtTime(now)
		scheduledTime := dpos.GetSlotTime(slot)
		fmt.Println(scheduledTime)
		block, err := reportBlock()
		if err != nil {
			return
		}
		fmt.Println("brocasting block", block)
	}
}

//func reportBlock(reportTime time.Time, reportor role.Delegate) *types.Block {
func reportBlock() (*types.Block, error) {
	chain := chain.GetChain()
	head := types.NewHeader()
	head.PrevBlockHash = chain.HeadBlockHash().Bytes()
	head.Number = chain.HeadBlockNum()
	head.Timestamp = chain.HeadBlockTime()
	head.Producer = []byte("my")
	block := types.NewBlock(head, nil)
	block.Header.ProducerSign = block.Sign("123").Bytes()
	produceractor.ApplyBlock(block)
	return block, nil

}
