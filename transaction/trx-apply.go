package transaction

import (
	"sync"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/role"
)

type TrxApplyService struct {
	stateDb		*db.DBService
}


var trxApplyServiceInst *TrxApplyService
var once sync.Once


func CreateTrxApplyService(dbInstance *db.DBService) *TrxApplyService {	
	once.Do(func() {
        trxApplyServiceInst = &TrxApplyService { stateDb : dbInstance}
	})
	
	return trxApplyServiceInst
}

func GetTrxApplyService() *TrxApplyService {	
	return trxApplyServiceInst
}

func (trxApplyService *TrxApplyService) CheckTransactionLifeTime(trx *types.Transaction) bool { 
	return true
}

func (trxApplyService *TrxApplyService) CheckTransactionUnique(trx *types.Transaction) bool { 	
	transactionExpiration, _ := role.GetTransactionExpirationObjectByHash(trxApplyService.stateDb, trx.Hash())
	if (nil != transactionExpiration) {
		return false
	}

	return true
}


func (trxApplyService *TrxApplyService) CheckTransactionMatchChain(trx *types.Transaction) bool { 
	return true
}

func (trxApplyService *TrxApplyService) SaveTransactionExpiration(trx *types.Transaction) { 
	var transactionExpiration = &role.TransactionExpiration{TrxHash:trx.Hash(), Expiration:trx.Expiration}
    role.SetTransactionExpirationObjectRole(trxApplyService.stateDb, trx.Hash(), transactionExpiration)
}

func (trxApplyService *TrxApplyService)ApplyTransaction(trx *types.Transaction) (bool, error) {

	/* check account validate,include contract account */
	/* check signature */
	if (!trxApplyService.CheckTransactionLifeTime(trx)) {
		return false, nil
	}

	if (!trxApplyService.CheckTransactionUnique(trx)) {
		return false, nil
	}
	
	if (!trxApplyService.CheckTransactionMatchChain(trx)) {
		return false, nil
	}
	trxApplyService.SaveTransactionExpiration(trx)

	/* call evm... */

	return true , nil
}
