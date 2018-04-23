

package transaction

import (
	"time"
	"sync"
	"fmt"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/action/message"
	"github.com/AsynkronIT/protoactor-go/actor"
	
)


var (
	trxExpirationCheckInterval    = time.Minute     // Time interval for check expiration pending transactions
	trxExpirationTime             = time.Minute     // Pending Trx max time , to be delete
)



type TrxPool struct {
	pending     map[common.Hash]*types.Transaction       
	expiration  map[common.Hash]time.Time    // to be delete
	
	mu           sync.RWMutex
	quit chan struct{}
}


func InitTrxPool() *TrxPool {
	
	// Create the transaction pool
	pool := &TrxPool{
		pending:      make(map[common.Hash]*types.Transaction),
		expiration:   make(map[common.Hash]time.Time),
		quit:         make(chan struct{}),
	}

	go pool.expirationCheckLoop()

	return pool
}


// expirationCheckLoop is periodically check exceed time transaction, then remove it
func (pool *TrxPool) expirationCheckLoop() {	
	expire := time.NewTicker(trxExpirationCheckInterval)
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


// expirationCheckLoop is periodically check exceed time transaction, then remove it
func (pool *TrxPool) addTransaction(trx *types.Transaction) {	
	pool.mu.Lock()
	trxHash := trx.Hash()
	pool.pending[trxHash] = trx
	//pool.expiration = time.Now()

	pool.mu.Unlock()
}


// expirationCheckLoop is periodically check exceed time transaction, then remove it
func (pool *TrxPool) AddTransaction(trx *types.Transaction) {
	pool.addTransaction(trx)
}



func (pool *TrxPool) Stop() {
	
	close(pool.quit)

	fmt.Println("Transaction pool stopped")
}

func (pool *TrxPool)CheckTransactionBaseConditionFromFront(){

	/* check max pending trx num */	
}


func (pool *TrxPool)CheckTransactionBaseConditionFromP2P(){	

}



// HandlTransactionFromFront handles a transaction from front
func (pool *TrxPool)HandleTransactionFromFront(context actor.Context, trx *types.Transaction) {
	
    pool.CheckTransactionBaseConditionFromFront()
	//start db session
	ApplyTransaction(trx)

	pool.addTransaction(trx)

	//revert db session

	//tell P2P actor to notify trx	

	context.Respond(true)
}


// HandlTransactionFromP2P handles a transaction from P2P
func (pool *TrxPool)HandleTransactionFromP2P(context actor.Context, trx *types.Transaction) {

	pool.CheckTransactionBaseConditionFromP2P()

	// start db session
	ApplyTransaction(trx)	

	pool.addTransaction(trx)

	//revert db session	
}



func (pool *TrxPool)HandlePushTransactionReq(context actor.Context, TrxSender message.TrxSenderType, trx *types.Transaction){

	if (message.TrxSenderTypeFront == TrxSender){ 
		pool.HandleTransactionFromFront(context, trx)
	} else if (message.TrxSenderTypeP2P == TrxSender) {
		pool.HandleTransactionFromP2P(context, trx)
	}	
}



func (pool *TrxPool)GetAllPendingTransactions(context actor.Context) {

	pool.mu.Lock()

	rsp := &message.GetAllPendingTrxRsp{}


	for txHash := range pool.pending {

		rsp.Trxs = append(rsp.Trxs, pool.pending[txHash])		
	}

	context.Respond(rsp)

	
	pool.mu.Unlock()
}


func (pool *TrxPool)RemoveTransactions(trxs []*types.Transaction){

	for _, trx := range trxs {
		delete(pool.pending, trx.Hash())
	}

}


func (pool *TrxPool)RemoveSingleTransaction(trx *types.Transaction){

	delete(pool.pending, trx.Hash())
}


func (pool *TrxPool)GetPendingTransaction(trxHash common.Hash) *types.Transaction {	

	return nil;
}