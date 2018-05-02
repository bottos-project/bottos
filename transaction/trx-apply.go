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
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/contract"
)

type TrxApplyService struct {
	roleIntf role.RoleInterface
	core        chain.BlockChainInterface
	ncIntf		contract.NativeContractInterface
}

var trxApplyServiceInst *TrxApplyService
var once sync.Once

func CreateTrxApplyService(env *env.ActorEnv) *TrxApplyService {
	once.Do(func() {
		trxApplyServiceInst = &TrxApplyService{roleIntf: env.RoleIntf, core: env.Chain, ncIntf:env.NcIntf}
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

	blockHistory, err := trxApplyService.roleIntf.GetBlockHistory(trx.Cursor)
	if (nil != err || nil == blockHistory) {
		return false
	} 

	var  chainCursorLabel uint32  = (uint32)(blockHistory.BlockHash[common.HashLength-1]) + (uint32)(blockHistory.BlockHash[common.HashLength-2])<<8 + (uint32)(blockHistory.BlockHash[common.HashLength-3])<<16 + (uint32)(blockHistory.BlockHash[common.HashLength-4])<<24

	if ( chainCursorLabel != trx.CursorLabel )  {
		fmt.Println("check chain match error,trx cursorlabel ", trx.CursorLabel, "chain cursollabel ", chainCursorLabel, "trx: ", trx.Hash())
		return false
	}

	return true
}

func (trxApplyService *TrxApplyService) SaveTransactionExpiration(trx *types.Transaction) {
	var transactionExpiration = &role.TransactionExpiration{TrxHash: trx.Hash(), Expiration: trx.Lifetime}
	trxApplyService.roleIntf.SetTransactionExpiration(trx.Hash(), transactionExpiration)
}

func (trxApplyService *TrxApplyService) ApplyTransaction(trx *types.Transaction) (bool, error) {
	/* check account validate,include contract account */
	/* check signature */	
	
	return true, nil

	account, error := trxApplyService.roleIntf.GetAccount(trx.Sender.Name)
	if(nil != error || nil == account) {
		fmt.Println("check account error, trx: ", trx.Hash())		
		return false, fmt.Errorf("check account error")
	}

	if !trxApplyService.CheckTransactionLifeTime(trx) {
		fmt.Println("check lift time error, trx: ", trx.Hash())
		return false, fmt.Errorf("check lift time error")
	}

	if !trxApplyService.CheckTransactionUnique(trx) {
		fmt.Println("check trx unique error, trx: ", trx.Hash())
		return false, fmt.Errorf("check trx unique error")
	}

	if !trxApplyService.CheckTransactionMatchChain(trx) {
		fmt.Println("check chain match error, trx: ", trx.Hash())
		return false, fmt.Errorf("check chain match error")
	}

	trxApplyService.SaveTransactionExpiration(trx)

	if (trxApplyService.ncIntf.IsNativeContract(trx.Contract.Name, trx.Method.Name) ) {

		applyContext := &contract.Context{RoleIntf:trxApplyService.roleIntf, Trx: trx}
		trxApplyService.ncIntf.ExecuteNativeContract(applyContext)
	} else {
        /* call evm... */
	}

	fmt.Println("trx : ", trx.Hash(),trx,"apply success")

	return true, nil
}
