package transaction

import (
	"sync"
	"fmt"

	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/role"
)

type TrxApplyService struct {
	roleIntf role.RoleInterface
}

var trxApplyServiceInst *TrxApplyService
var once sync.Once

func CreateTrxApplyService(roleIntf role.RoleInterface) *TrxApplyService {
	once.Do(func() {
		trxApplyServiceInst = &TrxApplyService{roleIntf: roleIntf}
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
	transactionExpiration, _ := trxApplyService.roleIntf.GetTransactionExpiration(trx.Hash())
	if nil != transactionExpiration {
		return false
	}

	return true
}

func (trxApplyService *TrxApplyService) CheckTransactionMatchChain(trx *types.Transaction) bool {
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
