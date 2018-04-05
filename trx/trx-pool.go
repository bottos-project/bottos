package trx

import (
	"fmt"
	"time"
	"github.com/bottos-project/bottos/core/common"
	"github.com/bottos-project/bottos/core/common/types"
	"github.com/bottos-project/bottos/core/event"
)

type TxPool struct{
	poolId string 
	em *event.TypeMux
} 

func CreateTxPool(em *event.TypeMux, bc *common.BlockChain) (*TxPool, error) {
	txpool := TxPool{"test", em}

	return &txpool, nil
}

// validate and queue transactions.
func (txpool *TxPool) Add(tx *types.Transaction) error {
	fmt.Println("TxPoolLoop : recv a tx")
	// mutex.Lock()
	// TODO process
	// mutex.Unlock()
	txpool.em.Post(common.TxPreEvent{Tx: tx})
	return nil
}

func (txpool *TxPool) TxPoolLoop() {
	fmt.Println("TxPoolLoop : Start")
	for {
		// TODO ChainStateEvent
		time.Sleep(100 * time.Millisecond)
	}
}
