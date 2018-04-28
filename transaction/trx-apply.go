package transaction

import (
	"github.com/bottos-project/core/action/env"
	"sync"
	"fmt"
	"time"

	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/chain"
)

type TrxApplyService struct {
	roleIntf role.RoleInterface
        core        chain.BlockChainInterface
}

var trxApplyServiceInst *TrxApplyService
var once sync.Once

func CreateTrxApplyService(env *env.ActorEnv) *TrxApplyService {
	once.Do(func() {
		trxApplyServiceInst = &TrxApplyService{roleIntf: env.RoleIntf, core: env.Chain}
	})

	return trxApplyServiceInst
}

func GetTrxApplyService() *TrxApplyService {
	return trxApplyServiceInst
}

func (trxApplyService *TrxApplyService) CheckTransactionLifeTime(trx *types.Transaction) bool {
	
	curTime := trxApplyService.core.HeadBlockTime()

	//for test:
	trx.Lifetime = curTime  + 600

	if (curTime >= trx.Lifetime) {
		fmt.Println("lifetime ", time.Unix((int64)(trx.Lifetime), 0),"have past, head time ", time.Unix((int64)(curTime), 0), "trx hash: ", trx.Hash())
		return false
	}

	

	if (trx.Lifetime >= (curTime + config.DEFAULT_MAX_LIFE_TIME)) {
		fmt.Println("lifetime ", time.Unix((int64)(trx.Lifetime), 0),"too far, head time ", time.Unix((int64)(curTime), 0), "trx hash: ", trx.Hash())
		return false
	}

	return true
}

func (trxApplyService *TrxApplyService) CheckTransactionUnique(trx *types.Transaction) bool {
	transactionExpiration, _ := trxApplyService.roleIntf.GetTransactionExpiration(trx.Hash())
	if nil != transactionExpiration {
		fmt.Println("check unique error ", trx.Hash())
		fmt.Println("transactionExpiration is  ", transactionExpiration)

		return false
	}

	return true
}

func (trxApplyService *TrxApplyService) CheckTransactionMatchChain(trx *types.Transaction) bool {



	//blockHistory, _ = role.GetBlockHistoryByNumber(stateDb, trx.Cursor)
	//if (nil != blockHistory) && blockHistory.BlockHash == trx.CursorLabel) {
	//	return true
	//}

	return true
}

func (trxApplyService *TrxApplyService) SaveTransactionExpiration(trx *types.Transaction) {
	var transactionExpiration = &role.TransactionExpiration{TrxHash: trx.Hash(), Expiration: trx.Lifetime}
	trxApplyService.roleIntf.SetTransactionExpiration(trx.Hash(), transactionExpiration)
}

func (trxApplyService *TrxApplyService) ApplyTransaction(trx *types.Transaction) (bool, error) {
	/* check account validate,include contract account */
	/* check signature */
	if !trxApplyService.CheckTransactionLifeTime(trx) {
		fmt.Println("check lift time error, trx: ", trx.Hash())
		return false, nil
	}

	if !trxApplyService.CheckTransactionUnique(trx) {
		fmt.Println("check trx unique error, trx: ", trx.Hash())
		return false, nil
	}

	if !trxApplyService.CheckTransactionMatchChain(trx) {
		fmt.Println("check chain match error, trx: ", trx.Hash())
		return false, nil
	}

	trxApplyService.SaveTransactionExpiration(trx)

	/* call evm... */

	fmt.Println("trx : ", trx.Hash(),trx,"apply success")

	return true, nil
}
