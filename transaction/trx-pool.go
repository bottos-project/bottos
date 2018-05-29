

package transaction

import (
	"time"
	"sync"
	"fmt"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/action/message"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/contractdb"

	proto "github.com/golang/protobuf/proto"
    "github.com/bottos-project/crypto-go/crypto"
    "crypto/sha256"
	"encoding/hex"
	bottosErr "github.com/bottos-project/bottos/common/errors"
	
)


var (
	trxExpirationCheckInterval    = 2*time.Second     // Time interval for check expiration pending transactions
)

var TrxPoolInst *TrxPool

type TrxPool struct {
	pending     map[common.Hash]*types.Transaction       
	roleIntf	role.RoleInterface
	contractDB  *contractdb.ContractDB

	mu           sync.RWMutex
	quit chan struct{}
}

func InitTrxPool(env *env.ActorEnv) *TrxPool {	
	// Create the transaction pool
	TrxPoolInst := &TrxPool{
		pending:      make(map[common.Hash]*types.Transaction),
		roleIntf:     env.RoleIntf,
		contractDB:   env.ContractDB,
		
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

			var currentTime = common.Now()
			for trxHash := range self.pending {				
				if (currentTime >= (self.pending[trxHash].Lifetime)) {				
					fmt.Println("remove expirate trx, hash is: ", trxHash,"curtime",currentTime,"lifeTime",self.pending[trxHash].Lifetime )	
					delete(self.pending, trxHash)					
				}				
			}
			
			self.mu.Unlock()

		case <-self.quit:
			return
		}
	}
}

func (self *TrxPool) addTransaction(trx *types.Transaction) {		

	self.mu.Lock()
	defer self.mu.Unlock()

	trxHash := trx.Hash()
	self.pending[trxHash] = trx

	//fmt.Println("add trx to pool succ, hash is: ",trx.Hash())
}

func (self *TrxPool) Stop() {	

	close(self.quit)

	fmt.Println("Transaction pool stopped")
}

func (self *TrxPool)CheckTransactionBaseConditionFromFront(trx *types.Transaction) (bool, bottosErr.ErrCode){

	if (config.DEFAULT_MAX_PENDING_TRX_IN_POOL <= (uint64)(len(self.pending))) {
		//return false, fmt.Errorf("check max pending trx num error")
		return false, bottosErr.ErrTrxPendingNumLimit		
	}

	/* check account validate,include contract account */
	
	if (!self.VerifySignature(trx)) {
		//return false, fmt.Errorf("check signature error")
		return false, bottosErr.ErrTrxSignError		
	}

	return true, bottosErr.ErrNoError
}

func (self *TrxPool)CheckTransactionBaseConditionFromP2P(){	

}

// HandlTransactionFromFront handles a transaction from front
func (self *TrxPool)HandleTransactionFromFront(context actor.Context, trx *types.Transaction) {

	fmt.Println("receive trx, detail: \n",trx)
	
	if checkResult, err := self.CheckTransactionBaseConditionFromFront(trx); true != checkResult {
		fmt.Println("check base condition  error, trx: ", trx.Hash())
		context.Respond(err)		
		return
	}
	//pool.stateDb.StartUndoSession()

	result , err , handleTrx := trxApplyServiceInst.ApplyTransaction(trx)
	if (!result) {
		fmt.Println("apply trx  error, trx: ", trx.Hash())
		context.Respond(err)	
		return
	}

	fmt.Println("handle trx succ, result: \n", handleTrx)

	self.addTransaction(trx)
	//pool.stateDb.Rollback()

	//tell P2P actor to notify trx	

	context.Respond(bottosErr.ErrNoError)
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
	for trxHash := range self.pending {
        fmt.Println("get all trx, add trx to rsp, hash is: ", self.pending[trxHash].Hash())
		rsp.Trxs = append(rsp.Trxs, self.pending[trxHash])		
	}

	context.Respond(rsp)
}

func (self *TrxPool)RemoveTransactions(trxs []*types.Transaction){

	for _, trx := range trxs {
		fmt.Println("remove trx, hash is: ", trx.Hash())
		delete(self.pending, trx.Hash())
	}
}

func (self *TrxPool)RemoveSingleTransaction(trx *types.Transaction){
    fmt.Println("remove single trx, hash is: ", trx.Hash())
	delete(self.pending, trx.Hash())
}

func (self *TrxPool)GetPendingTransaction(trxHash common.Hash) *types.Transaction {	

	return nil;
}


func (self *TrxPool)getPubKey(accountName string) ([]byte, error) {

	account ,err := self.roleIntf.GetAccount(accountName)
	if (nil != err) {
		return account.PublicKey, nil
	} else {
		return nil, fmt.Errorf("get account failed")
	}
	
	//for debug
	//pub_key, _ := hex.DecodeString("0488c8087c7fd0e1f0281c025902a444364a15e6732c65ff1c8b6673da977097447c1fd0c529482521a9883b0d1ce37e151b4572d4ecd996fefedcf0f6901508aa") 
	//return pub_key, nil
}



func (self *TrxPool) VerifySignature(trx *types.Transaction) bool {
	
	return true
	trxToVerify := &types.BasicTransaction {
			Version    :trx.Version    , 
			CursorNum  :trx.CursorNum  ,
			CursorLabel:trx.CursorLabel,
			Lifetime   :trx.Lifetime   ,
			Sender     :trx.Sender     ,
			Contract   :trx.Contract   ,
			Method     :trx.Method     ,
			Param      :trx.Param      ,
			SigAlg     :trx.SigAlg     ,
			//Signature  :[] byte{},
	}

	serializeData, err := proto.Marshal(trxToVerify)
	if nil != err {
		return false
	}
	
	senderPubKey ,err:= self.getPubKey(trx.Sender)
	if nil != err {
		return false
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(serializeData)))
	hashData := h.Sum(nil)

	verifyResult := crypto.VerifySign(senderPubKey, hashData, trx.Signature)
		
	fmt.Println("VerifySignature, result",verifyResult)

	return verifyResult
       
}

