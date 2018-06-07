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
	"fmt"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/contractdb"
	"github.com/bottos-project/bottos/role"

	"crypto/sha256"
	"encoding/hex"

	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/crypto-go/crypto"
	proto "github.com/golang/protobuf/proto"
)

var (
	trxExpirationCheckInterval = 2 * time.Second // Time interval for check expiration pending transactions
)

// TrxPoolInst is local var of TrxPool
var TrxPoolInst *TrxPool

// TrxPool is definition of trx pool
type TrxPool struct {
	pending     map[common.Hash]*types.Transaction
	roleIntf    role.RoleInterface
	contractDB  *contractdb.ContractDB
	netActorPid *actor.PID

	mu   sync.RWMutex
	quit chan struct{}
}

// InitTrxPool is init trx pool process when system start
func InitTrxPool(env *env.ActorEnv, netActorPid *actor.PID) *TrxPool {

	TrxPoolInst := &TrxPool{
		pending:     make(map[common.Hash]*types.Transaction),
		roleIntf:    env.RoleIntf,
		contractDB:  env.ContractDB,
		netActorPid: netActorPid,

		quit: make(chan struct{}),
	}

	CreateTrxApplyService(env)

	go TrxPoolInst.expirationCheckLoop()

	return TrxPoolInst
}

func (trxPool *TrxPool) expirationCheckLoop() {

	expire := time.NewTicker(trxExpirationCheckInterval)
	defer expire.Stop()

	for {
		select {
		case <-expire.C:
			trxPool.mu.Lock()

			var currentTime = common.Now()
			for trxHash := range trxPool.pending {
				if currentTime >= (trxPool.pending[trxHash].Lifetime) {
					fmt.Println("remove expirate trx, hash is: ", trxHash, "curtime", currentTime, "lifeTime", trxPool.pending[trxHash].Lifetime)
					delete(trxPool.pending, trxHash)
				}
			}

			trxPool.mu.Unlock()

		case <-trxPool.quit:
			return
		}
	}
}

func (trxPool *TrxPool) addTransaction(trx *types.Transaction) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	trxHash := trx.Hash()
	trxPool.pending[trxHash] = trx
}

// Stop is processing when system stop
func (trxPool *TrxPool) Stop() {

	close(trxPool.quit)

	fmt.Println("Transaction pool stopped")
}

// CheckTransactionBaseCondition is checking trx
func (trxPool *TrxPool) CheckTransactionBaseCondition(trx *types.Transaction) (bool, bottosErr.ErrCode) {

	if config.DEFAULT_MAX_PENDING_TRX_IN_POOL <= (uint64)(len(trxPool.pending)) {
		return false, bottosErr.ErrTrxPendingNumLimit
	}

	if !trxPool.VerifySignature(trx) {
		return false, bottosErr.ErrTrxSignError
	}

	return true, bottosErr.ErrNoError
}

// HandleTransactionCommon is processing trx
func (trxPool *TrxPool) HandleTransactionCommon(context actor.Context, trx *types.Transaction) {

	if checkResult, err := trxPool.CheckTransactionBaseCondition(trx); true != checkResult {
		context.Respond(err)
		return
	}

	result, err, _ := trxApplyServiceInst.ApplyTransaction(trx)
	if !result {
		context.Respond(err)
		return
	}

	trxPool.addTransaction(trx)

	notify := &message.NotifyTrx{
		Trx: trx,
	}
	trxPool.netActorPid.Tell(notify)

	context.Respond(bottosErr.ErrNoError)
}

// HandleTransactionFromFront is handling trx from front
func (trxPool *TrxPool) HandleTransactionFromFront(context actor.Context, trx *types.Transaction) {

	trxPool.HandleTransactionCommon(context, trx)
}

// HandleTransactionFromP2P is handling trx from P2P
func (trxPool *TrxPool) HandleTransactionFromP2P(context actor.Context, trx *types.Transaction) {

	trxPool.HandleTransactionCommon(context, trx)
}

// HandlePushTransactionReq is entry of trx req
func (trxPool *TrxPool) HandlePushTransactionReq(context actor.Context, TrxSender message.TrxSenderType, trx *types.Transaction) {

	if message.TrxSenderTypeFront == TrxSender {
		trxPool.HandleTransactionFromFront(context, trx)
	} else if message.TrxSenderTypeP2P == TrxSender {
		trxPool.HandleTransactionFromP2P(context, trx)
	}
}

// GetAllPendingTransactions is interface to get all pending trxs in trx pool
func (trxPool *TrxPool) GetAllPendingTransactions(context actor.Context) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	rsp := &message.GetAllPendingTrxRsp{}
	for trxHash := range trxPool.pending {
		rsp.Trxs = append(rsp.Trxs, trxPool.pending[trxHash])
	}

	context.Respond(rsp)
}

// RemoveTransactions is interface to remove trxs in trx pool
func (trxPool *TrxPool) RemoveTransactions(trxs []*types.Transaction) {

	for _, trx := range trxs {
		delete(trxPool.pending, trx.Hash())
	}
}

// RemoveSingleTransaction is interface to remove single trx in trx pool
func (trxPool *TrxPool) RemoveSingleTransaction(trx *types.Transaction) {

	delete(trxPool.pending, trx.Hash())
}

func (trxPool *TrxPool) getPubKey(accountName string) ([]byte, error) {

	account, err := trxPool.roleIntf.GetAccount(accountName)
	if nil != err {
		return nil, fmt.Errorf("get account failed")
	}

	return account.PublicKey, nil
}

// VerifySignature is verify signature from trx whether it is valid
func (trxPool *TrxPool) VerifySignature(trx *types.Transaction) bool {

	trxToVerify := &types.BasicTransaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       trx.Param,
		SigAlg:      trx.SigAlg,
	}

	serializeData, err := proto.Marshal(trxToVerify)
	if nil != err {
		return false
	}

	senderPubKey, err := trxPool.getPubKey(trx.Sender)
	if nil != err {
		return false
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(serializeData)))
	//h.Write([]byte(config.Param.ChainId))	
	hashData := h.Sum(nil)

	verifyResult := crypto.VerifySign(senderPubKey, hashData, trx.Signature)

	return verifyResult
}
