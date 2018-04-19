

package transaction

import (
	"time"
	"sync"
	"fmt"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/action/message"
	
)


var (
	expirationCheckInterval    = time.Minute     // Time interval for check expiration pending transactions
)



type TxPool struct {
	pending     map[common.Hash]*types.Transaction       
	expiration  map[common.Hash]time.Time    
	
	mu           sync.RWMutex
	quit chan struct{}
}


func InitTxPool() *TxPool {
	
	// Create the transaction pool
	pool := &TxPool{
		pending:      make(map[common.Hash]*types.Transaction),
		expiration:   make(map[common.Hash]time.Time),
		quit:         make(chan struct{}),
	}

	go pool.expirationCheckLoop()

	return pool
}


// expirationCheckLoop is periodically check exceed time transaction, then remove it
func (pool *TxPool) expirationCheckLoop() {
	
	expire := time.NewTicker(expirationCheckInterval)
	defer expire.Stop()

	for {
		select {
		case <-expire.C:
			pool.mu.Lock()

			var currentTime = time.Now()
			for txHash := range pool.expiration {

				if (currentTime.After(pool.expiration[txHash])) {
					delete(pool.expiration, txHash)
					delete(pool.pending, txHash)					
				}
				
			}
			pool.mu.Unlock()

		case <-pool.quit:
			return
		}
	}
}


func (pool *TxPool) Stop() {
	
	close(pool.quit)

	fmt.Println("Transaction pool stopped")
}

func CheckTransactionBaseConditionFromFront(){

	/* check max pending trx num */
	/* check account validate */
	/* check signature */

}


func CheckTransactionBaseConditionFromP2P(){	

}



// HandlTransactionFromFront handles a transaction from front
func HandleTransactionFromFront(trx *types.Transaction) {
	
    CheckTransactionBaseConditionFromFront()
	//start db session
	ApplyTransaction(trx)

	//add to pending

	//revert db session

	//tell P2P actor to notify trx	
}



// HandlTransactionFromP2P handles a transaction from P2P
func HandleTransactionFromP2P(trx *types.Transaction) {

	CheckTransactionBaseConditionFromP2P()

	// start db session
	ApplyTransaction(trx)
	//revert db session	
}



func HandlePushTransactionReq(TrxSender message.TrxSenderType, trx *types.Transaction){

	if (message.TrxSenderTypeFront == TrxSender){ 
		HandleTransactionFromFront(trx)
	} else if (message.TrxSenderTypeP2P == TrxSender) {
		HandleTransactionFromP2P(trx)
	}	
}
