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

func (trxApplyService *TrxApplyService) CheckTransactionLifeTime(trx *types.Transaction) { 

}

func (trxApplyService *TrxApplyService) CheckTransactionUnique(trx *types.Transaction) { 

}


func (trxApplyService *TrxApplyService) CheckTransactionMatchChain(trx *types.Transaction) { 

}

func (trxApplyService *TrxApplyService) SaveTransactionExpiration(trx *types.Transaction) { 
	var transactionExpiration = &role.TransactionExpiration{TrxHash:trx.Hash(), Expiration:trx.Expiration}
    role.SetTransactionExpirationObjectRole(trxApplyService.stateDb, trx.Hash(), transactionExpiration)
}

func (trxApplyService *TrxApplyService)ApplyTransaction(trx *types.Transaction) error {

	/* check account validate,include contract account */
	/* check signature */
	trxApplyService.CheckTransactionLifeTime(trx)
	trxApplyService.CheckTransactionUnique(trx)
	trxApplyService.CheckTransactionMatchChain(trx)
	trxApplyService.SaveTransactionExpiration(trx)

	/* call evm... */

	return nil
}
