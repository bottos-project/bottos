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
 * file description: apply transaction process
 * @Author: Wesley
 * @Date:   2017-12-15
 * @Last Modified by:
 * @Last Modified time:
 */

package transaction

import (
	//"fmt"
	"sync"
	"time"

	"github.com/bottos-project/bottos/action/env"

	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common"
	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/contract/contractdb"
	"github.com/bottos-project/bottos/role"
	wasm "github.com/bottos-project/bottos/vm/wasm/exec"
	log "github.com/cihub/seelog"
)

// TrxApplyService is to define a service for apply a transaction
type TrxApplyService struct {
	roleIntf   role.RoleInterface
	ContractDB *contractdb.ContractDB
	core       chain.BlockChainInterface
	ncIntf     contract.NativeContractInterface
}

var trxApplyServiceInst *TrxApplyService
var once sync.Once

// CreateTrxApplyService is to new a TrxApplyService
func CreateTrxApplyService(env *env.ActorEnv) *TrxApplyService {
	once.Do(func() {
		trxApplyServiceInst = &TrxApplyService{roleIntf: env.RoleIntf, ContractDB: env.ContractDB, core: env.Chain, ncIntf: env.NcIntf}
	})

	return trxApplyServiceInst
}

// GetTrxApplyService is to get trxApplyService
func GetTrxApplyService() *TrxApplyService {
	return trxApplyServiceInst
}

// CheckTransactionLifeTime is to check life time of a transaction
func (trxApplyService *TrxApplyService) CheckTransactionLifeTime(trx *types.Transaction) bool {

	chainState, _ := trxApplyService.roleIntf.GetChainState()
	curTime := chainState.LastBlockTime

	if curTime >= trx.Lifetime {
		log.Errorf("lifetime %v have past, head time %v, trx hash: %x",time.Unix((int64)(trx.Lifetime), 0), time.Unix((int64)(curTime), 0), trx.Hash())
		return false
	}

	if trx.Lifetime >= (curTime + config.DEFAULT_MAX_LIFE_TIME) {
		log.Errorf("lifetime %v too far, head time %v, trx hash: %x",time.Unix((int64)(trx.Lifetime), 0), time.Unix((int64)(curTime), 0), trx.Hash())
		return false
	}

	return true
}

// CheckTransactionUnique is to check whether a transaction is unique
func (trxApplyService *TrxApplyService) CheckTransactionUnique(trx *types.Transaction) bool {

	transactionExpiration, _ := trxApplyService.roleIntf.GetTransactionExpiration(trx.Hash())
	if nil != transactionExpiration {
		log.Errorf("check unique error, trx: %x", trx.Hash())
		log.Error("transactionExpiration is  ", transactionExpiration)

		return false
	}

	return true
}

// CheckTransactionMatchChain is to check whether it is the right chain
func (trxApplyService *TrxApplyService) CheckTransactionMatchChain(trx *types.Transaction) bool {

	blockHistory, err := trxApplyService.roleIntf.GetBlockHistory(trx.CursorNum)
	if nil != err || nil == blockHistory {
		log.Error("get block history error")
		return false
	}

	var chainCursorLabel uint32 = (uint32)(blockHistory.BlockHash[common.HashLength-1]) + (uint32)(blockHistory.BlockHash[common.HashLength-2])<<8 + (uint32)(blockHistory.BlockHash[common.HashLength-3])<<16 + (uint32)(blockHistory.BlockHash[common.HashLength-4])<<24

	if chainCursorLabel != trx.CursorLabel {
		log.Errorf("check chain match error, trx cursorlabel %v, chain cursollabel %v, trx: %x", trx.CursorLabel, chainCursorLabel, trx.Hash())
		return false
	}

	return true
}

// SaveTransactionExpiration is to save the expiration of a transaction
func (trxApplyService *TrxApplyService) SaveTransactionExpiration(trx *types.Transaction) {

	var transactionExpiration = &role.TransactionExpiration{TrxHash: trx.Hash(), Expiration: trx.Lifetime}
	trxApplyService.roleIntf.SetTransactionExpiration(trx.Hash(), transactionExpiration)
}

// ApplyTransaction is to handle a transaction, include parameters checking
func (trxApplyService *TrxApplyService) ApplyTransaction(trx *types.Transaction) (bool, bottosErr.ErrCode, *types.HandledTransaction) {

	account, getAccountErr := trxApplyService.roleIntf.GetAccount(trx.Sender)
	if nil != getAccountErr || nil == account {
		log.Errorf("check account error, trx: %x", trx.Hash())
		return false, bottosErr.ErrTrxAccountError, nil
	}

	if !trxApplyService.CheckTransactionLifeTime(trx) {		
		return false, bottosErr.ErrTrxLifeTimeError, nil
	}

	if !trxApplyService.CheckTransactionUnique(trx) {
		return false, bottosErr.ErrTrxUniqueError, nil
	}

	if !trxApplyService.CheckTransactionMatchChain(trx) {
		return false, bottosErr.ErrTrxChainMathError, nil
	}

	trxApplyService.SaveTransactionExpiration(trx)

	result, bottosError, derivedTrxList := trxApplyService.ProcessTransaction(trx, 0)

	if false == result {
		log.Errorf("process trx error, trx: %x", trx.Hash())
		return false, bottosError, nil
	}

	handleTrx := &types.HandledTransaction{
		Transaction: trx,
		DerivedTrx:  derivedTrxList,
	}

	return true, bottosErr.ErrNoError, handleTrx
}

// ProcessTransaction is to handle a transaction without parameters checking
func (trxApplyService *TrxApplyService) ProcessTransaction(trx *types.Transaction, deepLimit uint32) (bool, bottosErr.ErrCode, []*types.DerivedTransaction) {

	if deepLimit >= config.DEFAUL_MAX_CONTRACT_DEPTH {
		return false, bottosErr.ErrTrxContractDepthError, nil
	}

	var derivedTrx []*types.DerivedTransaction

	bottoserr := bottosErr.ErrNoError

	applyContext := &contract.Context{RoleIntf: trxApplyService.roleIntf, ContractDB: trxApplyService.ContractDB, Trx: trx}

	if trxApplyService.ncIntf.IsNativeContract(trx.Contract, trx.Method) {
		contErr := trxApplyService.ncIntf.ExecuteNativeContract(applyContext)
		bottoserr = contract.ConvertErrorCode(contErr)
		if bottosErr.ErrNoError == bottoserr {
			return true, bottosErr.ErrNoError, nil
		}

		log.Error("process trx, failed bottos error: ", bottoserr)
		return false, bottoserr, nil

	}
	// else branch
	trxList, exeErr := wasm.GetInstance().Start(applyContext, 1, false)
	if nil != exeErr {
		log.Error("process trx failed, error: ",exeErr)
		return false, bottosErr.ErrTrxContractHanldeError, nil
	}

	log.Trace("derived trx list len is ", len(trxList))
	for _, subTrx := range trxList {
		log.Trace(subTrx)
	}

	if (uint32(len(trxList))) >= config.DEFAUL_MAX_SUB_CONTRACT_NUM {
		return false, bottosErr.ErrTrxSubContractNumError, nil
	}

	for _, subTrx := range trxList {
		result, bottosErr, subDerivedTrx := trxApplyService.ProcessTransaction(subTrx, deepLimit+1)
		if false == result {
			return false, bottosErr, nil
		}

		handleTrx := &types.DerivedTransaction{
			Transaction: subTrx,
			DerivedTrx:  subDerivedTrx,
		}
		derivedTrx = append(derivedTrx, handleTrx)
	}
	return true, bottosErr.ErrNoError, derivedTrx

}
