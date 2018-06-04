//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  trx pool
 * @Author: Wesley
 * @Date:   2017-12-15
 * @Last Modified by:
 * @Last Modified time:
 */

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
}

func (self *TrxPool) Stop() {	

	close(self.quit)

	fmt.Println("Transaction pool stopped")
}

func (self *TrxPool)CheckTransactionBaseConditionFromFront(trx *types.Transaction) (bool, bottosErr.ErrCode){

	if (config.DEFAULT_MAX_PENDING_TRX_IN_POOL <= (uint64)(len(self.pending))) {
		return false, bottosErr.ErrTrxPendingNumLimit		
	}
	
	if (!self.VerifySignature(trx)) {
		return false, bottosErr.ErrTrxSignError		
	}

	return true, bottosErr.ErrNoError
}

func (self *TrxPool)CheckTransactionBaseConditionFromP2P(){	

}

func (self *TrxPool)HandleTransactionFromFront(context actor.Context, trx *types.Transaction) {

	if checkResult, err := self.CheckTransactionBaseConditionFromFront(trx); true != checkResult {
		fmt.Println("check base condition  error, trx: ", trx.Hash())
		context.Respond(err)		
		return
	}

	result , err , _ := trxApplyServiceInst.ApplyTransaction(trx)
	if (!result) {
		fmt.Println("apply trx  error, trx: ", trx.Hash())
		context.Respond(err)	
		return
	}

	self.addTransaction(trx)

	context.Respond(bottosErr.ErrNoError)
}

func (self *TrxPool)HandleTransactionFromP2P(context actor.Context, trx *types.Transaction) {

	self.CheckTransactionBaseConditionFromP2P()

	trxApplyServiceInst.ApplyTransaction(trx)	

	self.addTransaction(trx)
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
		rsp.Trxs = append(rsp.Trxs, self.pending[trxHash])		
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


func (self *TrxPool)getPubKey(accountName string) ([]byte, error) {

	account ,err := self.roleIntf.GetAccount(accountName)
	if (nil == err) {
		return account.PublicKey, nil
	} else {
		return nil, fmt.Errorf("get account failed")
	}
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

	return verifyResult       
}

