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
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/context"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/version"

	bottosErr "github.com/bottos-project/bottos/common/errors"
	log "github.com/cihub/seelog"
)

var (
	trxExpirationCheckInterval = 60 * time.Second // Time interval for check expiration pending transactions
	trxCacheCheckInterval      = 2 * time.Second  // Time interval for check cache pending transactions
)

// TrxPoolInst is local var of TrxPool
var TrxPoolInst *TrxPool

type CachedTransaction struct {
	hash common.Hash
	msg  interface{} //*message.PushTrxForP2PReq or *message.ReceiveTrx
}

// TrxPool is definition of trx pool
type TrxPool struct {
	cache       []*CachedTransaction
	cacheMap    map[common.Hash]*CachedTransaction
	pending     map[common.Hash]*types.Transaction
	roleIntf    role.RoleInterface
	protocol    context.ProtocolInterface
	netActorPid *actor.PID
	trxActorPid *actor.PID

	dbInst   *db.DBService

	mu   sync.RWMutex
	cacheMutex sync.RWMutex
	quit chan struct{}
}

// InitTrxPool is init trx pool process when system start
func InitTrxPool(dbInstance *db.DBService, roleIntf role.RoleInterface, nc contract.NativeContractInterface, protocol context.ProtocolInterface, netActorPid *actor.PID) *TrxPool {

	TrxPoolInst = &TrxPool{
		cache:       make([]*CachedTransaction, 0, 100),
		cacheMap:    make(map[common.Hash]*CachedTransaction),
		pending:     make(map[common.Hash]*types.Transaction),
		roleIntf:    roleIntf,
		protocol:    protocol,
		netActorPid: netActorPid,
		trxActorPid: trxActorPid,
		dbInst:      dbInstance,
		quit: make(chan struct{}),
	}

	CreateTrxApplyService(roleIntf, nc)

	go TrxPoolInst.cacheCheckLoop()
	go TrxPoolInst.expirationCheckLoop()

	return TrxPoolInst
}

func (trxPool *TrxPool) cacheCheckLoop() {

	expire := time.NewTicker(trxCacheCheckInterval)
	defer expire.Stop()

	for {
		select {
		case <-expire.C:

			trxPool.cacheMutex.Lock()

			if len(trxPool.cache) > 0 && trxPool.protocol.GetBlockSyncState() {
				log.Infof("TRX trx num in cache before check: %v", len(trxPool.cache))
				for _, c := range trxPool.cache {
					log.Infof("TRX remove cache trx %x", c.hash)
					trxPool.trxActorPid.Tell(c.msg)
				}
				trxPool.cache = trxPool.cache[0:0]
				trxPool.cacheMap = make(map[common.Hash]*CachedTransaction)
				log.Infof("TRX trx num in cache after check: %v", len(trxPool.cache))
			}

			trxPool.cacheMutex.Unlock()

		case <-trxPool.quit:
			return
		}
	}
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

func (trxPool *TrxPool) IsCacheEmpty() bool {

	trxPool.cacheMutex.Lock()
	defer trxPool.cacheMutex.Unlock()

	return len(trxPool.cache) == 0
}

func (trxPool *TrxPool) IsTransactionInCache(trxHash common.Hash) bool {

	trxPool.cacheMutex.Lock()
	defer trxPool.cacheMutex.Unlock()

	_, exist := trxPool.cacheMap[trxHash]
	return exist
}

func (trxPool *TrxPool) AddTransactionToCache(trxMsg interface{}) {
	trxPool.cacheMutex.Lock()
	defer trxPool.cacheMutex.Unlock()

	var trxHash common.Hash
	switch msg := trxMsg.(type) {
	case *message.PushTrxForP2PReq:
		trxHash = msg.P2PTrx.Transaction.Hash()
	case *message.ReceiveTrx:
		trxHash = msg.P2PTrx.Transaction.Hash()
	default:
		log.Errorf("add trx to cache, unknown msg type %v", msg)
		return
	}

	c := &CachedTransaction{hash: trxHash, msg: trxMsg}
	trxPool.cache = append(trxPool.cache, c)
	trxPool.cacheMap[c.hash] = c
	log.Infof("TRX add trx %x to cache, num in cache %v", c.hash, len(trxPool.cache))
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

	if trxPool.IsTransactionInCache(trx.Hash()) {
		log.Infof("TRX check exist error, already in cache, trx %x ", trx.Hash())

		return false, bottosErr.ErrTrxAlreadyInCache
	}

	if trxPool.isTransactionExist(trx) {
		log.Infof("TRX check exist error, already in pool, trx %x ", trx.Hash())

		return false, bottosErr.ErrTrxAlreadyInPool
	}

	if len(trxPool.cache) > 0 && (config.DEFAULT_MAX_PENDING_TRX_IN_POOL/10) <= (uint64)(len(trxPool.cache)) {
		log.Errorf("TRX check cache num reach max error, trx %x", trx.Hash())
		return false, bottosErr.ErrTrxCacheNumLimit
	}

	if config.DEFAULT_MAX_PENDING_TRX_IN_POOL <= (uint64)(len(trxPool.pending)) {
		log.Errorf("TRX check pool num reach max error, trx %x", trx.Hash())
		return false, bottosErr.ErrTrxPendingNumLimit
	}

	chainState, _ := trxPool.roleIntf.GetChainState()
	myVersion := version.GetVersionByBlockNum(chainState.LastBlockNum)
	if myVersion != nil && trx.Version > myVersion.VersionNumber {
		log.Errorf("VERSION handle CheckTransactionBaseCondition failed, trx.hash %x, trx.version %v, my version %v", trx.Hash(), version.GetStringVersion(trx.Version), myVersion.VersionString)
		
		return false, bottosErr.ErrTrxVersionError
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

// SendP2PTrx sends p2p trx
func (trxPool *TrxPool) SendP2PTrx(p2pTrx *types.P2PTransaction){
	notify := &message.NotifyTrx{
		P2PTrx: p2pTrx,
	}
	trxPool.netActorPid.Tell(notify)
}

// HandleTransactionFromFront is handling trx from front
func (trxPool *TrxPool) HandleTransactionFromFront(context actor.Context, p2pTrx *types.P2PTransaction) {
	log.Infof("TRX rcv trx from front, trx %x, sender %v, contract %v, method %v, TTL %v", p2pTrx.Transaction.Hash(), p2pTrx.Transaction.Sender, p2pTrx.Transaction.Contract, p2pTrx.Transaction.Method, p2pTrx.TTL)

	err := trxPool.HandleTransactionCommon(context, p2pTrx.Transaction)
	
	if bottosErr.ErrNoError != err {
		log.Errorf("TRX handle trx from front failed, trx %x, error %v", p2pTrx.Transaction.Hash(), err)
		trxApplyServiceInst.AddTrxErrorCode(p2pTrx.Transaction.Hash(), err)
	} else {
		trxPool.SendP2PTrx(p2pTrx)
	}
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
		log.Infof("TRX rm trx %x", trx.Hash())
		delete(trxPool.pending, trx.Hash())
	}

	log.Infof("TRX after rm trx num in pool %v", len(trxPool.pending))
}

// RemoveBlockTransactions is interface to remove trxs in trx pool
func (trxPool *TrxPool) RemoveBlockTransactions(trxs []*types.BlockTransaction) {

	trxPool.mu.Lock()
	defer trxPool.mu.Unlock()

	for _, trx := range trxs {
		log.Infof("TRX rm block trx %x", trx.Transaction.Hash())
		delete(trxPool.pending, trx.Transaction.Hash())
	}

	log.Infof("TRX after rm block trx num in pool %v", len(trxPool.pending))
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

