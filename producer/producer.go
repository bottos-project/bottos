package producer

import (
	"fmt"
	"time"
	"sync"

	"github.com/bottos-project/bottos/core/event"
	"github.com/bottos-project/bottos/core/common"
	"github.com/bottos-project/bottos/core/common/types"
)

type Producer struct {
	em *event.TypeMux
	bc *common.BlockChain
	cnt int

	txSub event.Subscription

	mu sync.RWMutex
	txcnt int
}

func NewProducer(em *event.TypeMux, bc *common.BlockChain) *Producer {
	sub := em.Subscribe(common.TxPreEvent{})
	producer := &Producer{em, bc, 0, sub, sync.RWMutex{}, 0}

	return producer
}

func (p *Producer) TxRecvLoop() {
	for obj := range p.txSub.Chan() {
		switch ev := obj.(type) {
		case common.TxPreEvent:
			fmt.Printf("Producer : recv a tx Id=%s\n", ev.Tx.Id)
			p.mu.Lock()
			p.txcnt += 1;
			p.mu.Unlock()
		}
	}
}

func (p *Producer) ProducerLoop() {

	fmt.Println("ProducerLoop : Start")

	go p.TxRecvLoop()

	for {
		fmt.Println("\n\nProducerLoop : wating for my turn...")
		time.Sleep(3000 * time.Millisecond)
		
		fmt.Println("ProducerLoop : producing...")
		p.cnt += 1
		p.mu.Lock()
		fmt.Printf("ProducerLoop : produced, txn=%d\n", p.txcnt)
		if p.txcnt > 0 {
			p.txcnt = 0
		}
		p.mu.Unlock()
		
	    header := types.NewHeader([]byte{}, p.cnt)
		b := types.NewBlock(header)

		p.bc.PushBlock(b)
		p.em.Post(common.NewMinedBlockEvent{Block: b})

		time.Sleep(100 * time.Millisecond)
	}
}
