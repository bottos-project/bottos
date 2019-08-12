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
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"

	"crypto/sha256"
	"encoding/hex"

	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/crypto-go/crypto"
	log "github.com/cihub/seelog"
)

var (
	trxExpirationCheckInterval = 60 * time.Second // Time interval for check expiration pending transactions
)

// TrxPoolInst is local var of TrxPool
var TrxPoolInst *TrxPool

// TrxPool is definition of trx pool
type TrxPool struct {
	pending     map[common.Hash]*types.Transaction
	roleIntf    role.RoleInterface
	netActorPid *actor.PID

	dbInst *db.DBService
	mu     sync.RWMutex
	quit   chan struct{}
}

// InitTrxPool is init trx pool process when system start
func InitTrxPool(dbInstance *db.DBService, roleIntf role.RoleInterface, nc contract.NativeContractInterface, protocol context.ProtocolInterface, netActorPid *actor.PID) *TrxPool {

	TrxPoolInst = &TrxPool{
		pending:     make(map[common.Hash]*types.Transaction),
		roleIntf:    roleIntf,
		netActorPid: netActorPid,
		dbInst:      dbInstance,
		quit: make(chan struct{}),
	}

	CreateTrxApplyService(roleIntf, nc)

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
			log.Infof("TRX trx num in pool before check: %v", len(trxPool.pending))
			for trxHash := range trxPool.pending {
				if currentTime >= (trxPool.pending[trxHash].Lifetime) {
					log.Infof("TRX remove expirate trx, trx %x, curtime %v, lifeTime %v", trxHash, currentTime, trxPool.pending[trxHash].Lifetime)
					trxPool.RemoveSingleTransactionbyHashNotLock(trxHash)
				}
			}
			log.Infof("TRX trx num in pool after check: %v", len(trxPool.pending))

			trxPool.mu.Unlock()

		case <-trxPool.quit:
			return
		}
	}
}

func (trxPool *TrxPool) isTransactionExist(trx *types.Transaction) bool {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	if nil == trxPool.pending[trx.Hash()] {
		return false
	} else {
		return true
	}
}

func (trxPool *TrxPool) addTransaction(trx *types.Transaction) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	trxHash := trx.Hash()

	trxPool.pending[trxHash] = trx

	log.Infof("TRX add add trx, num in pool %v", len(trxPool.pending))

}

// Stop is processing when system stop
func (trxPool *TrxPool) Stop() {

	close(trxPool.quit)

	log.Errorf("TRX Transaction pool stopped")
}

// CheckTransactionBaseCondition is checking trx
func (trxPool *TrxPool) CheckTransactionBaseCondition(trx *types.Transaction) (bool, bottosErr.ErrCode) {
	if isTransactionExist(trx) {
		return false, bottos.ErrTrxAlreadyInPool
	}
	if config.DEFAULT_MAX_PENDING_TRX_IN_POOL <= (uint64)(len(trxPool.pending)) {
		log.Errorf("trx %x pending num over", trx.Hash())
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
}

// HandleTransactionFromFront is handling trx from front
func (trxPool *TrxPool) HandleTransactionFromFront(context actor.Context, trx *types.Transaction) {
	log.Infof("rcv trx %x from front,sender %v, contract %v, method %v", trx.Hash(), trx.Sender, trx.Contract, trx.Method)
	trxPool.HandleTransactionCommon(context, trx)
}

// HandleTransactionFromP2P is handling trx from P2P
func (trxPool *TrxPool) HandleTransactionFromP2P(context actor.Context, p2pTrx *types.P2PTransaction) {
	log.Infof("TRX rcv trx from P2P, trx %x, sender %v, contract %v method %v, TTL %v", p2pTrx.Transaction.Hash(), p2pTrx.Transaction.Sender, p2pTrx.Transaction.Contract, p2pTrx.Transaction.Method, p2pTrx.TTL)
	err := trxPool.HandleTransactionCommon(context, p2pTrx.Transaction)	

	if bottosErr.ErrNoError != err {
		log.Errorf("TRX handle trx from node failed, trx %x, error %v", p2pTrx.Transaction.Hash(), err)
	} else {
		trxPool.SendP2PTrx(p2pTrx)
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

	log.Infof("TRX get all pending trx num in pool, total num %v", len(rsp.Trxs))

	context.Respond(rsp)
}

// GetAllPendingTransactions is interface to get all pending trxs in trx pool
func (trxPool *TrxPool) GetAllPendingTransactions4funcCall() []*types.Transaction {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	var trxs []*types.Transaction
	for trxHash := range trxPool.pending {
		trxs = append(trxs, trxPool.pending[trxHash])
	}
	return trxs
}

// RemoveTransactions is interface to remove trxs in trx pool
func (trxPool *TrxPool) RemoveTransactions(trxs []*types.Transaction) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	for _, trx := range trxs {
		delete(trxPool.pending, trx.Hash())
	}
}

// RemoveSingleTransaction is interface to remove single trx in trx pool
func (trxPool *TrxPool) RemoveSingleTransaction(trx *types.Transaction) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	log.Infof("TRX rm single trx %x", trx.Hash())

	delete(trxPool.pending, trx.Hash())

	log.Infof("TRX after rm single trx num in pool %v", len(trxPool.pending))
}

// RemoveSingleTransactionbyHash is interface to remove single trx in trx pool
func (trxPool *TrxPool) RemoveSingleTransactionbyHash(trxHash common.Hash) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	log.Infof("TRX rm trx by hash, trx %x", trxHash)

	delete(trxPool.pending, trxHash)

	log.Infof("TRX after rm trx by hash, trx num in pool %v", len(trxPool.pending))
}

// RemoveSingleTransactionbyHash is interface to remove single trx in trx pool
func (trxPool *TrxPool) RemoveSingleTransactionbyHashNotLock(trxHash common.Hash) {

	log.Infof("TRX rm trx by hash not lock, trx %x", trxHash)

	delete(trxPool.pending, trxHash)

	log.Infof("TRX after rm trx by hash not lock, trx num in pool %v", len(trxPool.pending))
}

