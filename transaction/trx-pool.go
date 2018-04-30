

package transaction

import (
	"time"
	"sync"
	"fmt"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/action/message"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/action/env"
	"github.com/bottos-project/core/config"
)


var (
	trxExpirationCheckInterval    = time.Minute     // Time interval for check expiration pending transactions
	trxExpirationTime             = time.Minute     // Pending Trx max time , to be delete
)

var TrxPoolInst *TrxPool


type TrxPool struct {
	pending     map[common.Hash]*types.Transaction       
	expiration  map[common.Hash]time.Time    // to be delete
	roleIntf	role.RoleInterface
	
	mu           sync.RWMutex
	quit chan struct{}
}


func InitTrxPool(env *env.ActorEnv) *TrxPool {
	
	// Create the transaction pool
	TrxPoolInst := &TrxPool{
		pending:      make(map[common.Hash]*types.Transaction),
		expiration:   make(map[common.Hash]time.Time),
		roleIntf:     env.RoleIntf,
		
		quit:         make(chan struct{}),		
	}

	CreateTrxApplyService(env)

	go TrxPoolInst.expirationCheckLoop()

	return TrxPoolInst
}


// expirationCheckLoop is periodically check exceed time transaction, then remove it
func (self *TrxPool) expirationCheckLoop() {	
	expire := time.NewTicker(trxExpirationCheckInterval)
	defer expire.Stop()

	for {
		select {
		case <-expire.C:
			self.mu.Lock()

			var currentTime = time.Now()
			for txHash := range self.expiration {

				if (currentTime.After(self.expiration[txHash])) {
					delete(self.expiration, txHash)
					delete(self.pending, txHash)					
				}
				
			}
			self.mu.Unlock()

		case <-self.quit:
			return
		}
	}
}


// expirationCheckLoop is periodically check exceed time transaction, then remove it
func (self *TrxPool) addTransaction(trx *types.Transaction) {	
	self.mu.Lock()
	defer self.mu.Unlock()

	trxHash := trx.Hash()
	self.pending[trxHash] = trx
}

func (self *TrxPool) Stop() {
	
	close(self.quit)

	fmt.Println("Transaction pool stopped")
}

func (self *TrxPool)CheckTransactionBaseConditionFromFront() bool {

	if (config.DEFAULT_MAX_PENDING_TRX_IN_POOL <= (uint64)(len(self.pending))) {
		return false
	}
	return true
}


func (self *TrxPool)CheckTransactionBaseConditionFromP2P(){	

}



// HandlTransactionFromFront handles a transaction from front
func (self *TrxPool)HandleTransactionFromFront(context actor.Context, trx *types.Transaction) {
	fmt.Println("receive trx: ",trx, "hash: ", trx.Hash())

	fmt.Printf("%s",trx.Param)
	
	if (!self.CheckTransactionBaseConditionFromFront()) {
		fmt.Println("check base condition  error, trx: ", trx.Hash())

		return
	}
	//pool.stateDb.StartUndoSession()

	result , _ := trxApplyServiceInst.ApplyTransaction(trx)
	if (!result) {
		fmt.Println("apply trx  error, trx: ", trx.Hash())
		return
	}

	self.addTransaction(trx)
	//pool.stateDb.Rollback()

	//tell P2P actor to notify trx	

	context.Respond(true)
}


// HandlTransactionFromP2P handles a transaction from P2P
func (self *TrxPool)HandleTransactionFromP2P(context actor.Context, trx *types.Transaction) {

	self.CheckTransactionBaseConditionFromP2P()

	// start db session
	trxApplyServiceInst.ApplyTransaction(trx)	

	self.addTransaction(trx)

	//revert db session	
}



func (self *TrxPool)HandlePushTransactionReq(context actor.Context, TrxSender message.TrxSenderType, trx *types.Transaction){

	if (message.TrxSenderTypeFront == TrxSender){ 
		self.HandleTransactionFromFront(context, trx)
	} else if (message.TrxSenderTypeP2P == TrxSender) {
		self.HandleTransactionFromP2P(context, trx)
	}	
}



func (self *TrxPool)GetAllPendingTransactions(context actor.Context) {

	self.mu.Lock()

	defer self.mu.Unlock()

	rsp := &message.GetAllPendingTrxRsp{}


	for txHash := range self.pending {

		rsp.Trxs = append(rsp.Trxs, self.pending[txHash])		
	}

	context.Respond(rsp)
}


func (self *TrxPool)RemoveTransactions(trxs []*types.Transaction){

	for _, trx := range trxs {
		delete(self.pending, trx.Hash())
	}

}


func (self *TrxPool)RemoveSingleTransaction(trx *types.Transaction){

	delete(self.pending, trx.Hash())
}


func (self *TrxPool)GetPendingTransaction(trxHash common.Hash) *types.Transaction {	

	return nil;
}