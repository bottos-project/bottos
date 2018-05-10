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
	"github.com/bottos-project/core/contract/contractdb"
	wasm "github.com/bottos-project/core/vm/wasm/exec"
)

type TrxApplyService struct {
	roleIntf role.RoleInterface
	ContractDB *contractdb.ContractDB
	core        chain.BlockChainInterface
	ncIntf		contract.NativeContractInterface
}

var trxApplyServiceInst *TrxApplyService
var once sync.Once

func CreateTrxApplyService(env *env.ActorEnv) *TrxApplyService {
	once.Do(func() {
		trxApplyServiceInst = &TrxApplyService{roleIntf: env.RoleIntf, ContractDB: env.ContractDB, core: env.Chain, ncIntf:env.NcIntf}
	})

	return trxApplyServiceInst
}

func GetTrxApplyService() *TrxApplyService {
	return trxApplyServiceInst
}

func (trxApplyService *TrxApplyService) CheckTransactionLifeTime(trx *types.Transaction) bool {
	
	//curTime := trxApplyService.core.HeadBlockTime()
	//for test:
	curTime := common.Now()

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

	blockHistory, err := trxApplyService.roleIntf.GetBlockHistory(trx.CursorNum)
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
	
	account, getAccountErr := trxApplyService.roleIntf.GetAccount(trx.Sender)
	if(nil != getAccountErr || nil == account) {
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

	var exeErr error
	var exeResult bool

	applyContext := &contract.Context{RoleIntf:trxApplyService.roleIntf, ContractDB: trxApplyService.ContractDB, Trx: trx}

	if (trxApplyService.ncIntf.IsNativeContract(trx.Contract, trx.Method) ) {
		exeErr = trxApplyService.ncIntf.ExecuteNativeContract(applyContext)
	} else {
		/* call evm... */		
		_, exeErr = wasm.GetInstance().Start(applyContext, 1, false)
	}

    if (nil == exeErr) {
		exeResult = true
        fmt.Println("trx : ", trx.Hash(),trx,"apply success")
	}else {
		exeResult = false
		fmt.Println("trx : ", trx.Hash(),trx,"apply failed")
	}

	return exeResult, exeErr
}
