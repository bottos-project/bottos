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

	"github.com/bottos-project/bottos/common"
	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/role"
	wasm "github.com/bottos-project/bottos/vm/wasm/exec"
	log "github.com/cihub/seelog"

	"github.com/bottos-project/bottos/vm/duktape"
	"github.com/bottos-project/bottos/common/vm"
)

// TrxApplyService is to define a service for apply a transaction
type TrxApplyService struct {
	roleIntf   role.RoleInterface
	ncIntf     contract.NativeContractInterface

	trxHashErrorList     [config.DEFAULT_MAX_TRX_ERROR_CODE_NUM]common.Hash
	curTrxErrorCodeIndex uint64
	trxHashErrorMap      map[common.Hash]bottosErr.ErrCode

	mu sync.RWMutex
}

var trxApplyServiceInst *TrxApplyService
var once sync.Once

// CreateTrxApplyService is to new a TrxApplyService
func CreateTrxApplyService(roleIntf role.RoleInterface, nc contract.NativeContractInterface) *TrxApplyService {
	once.Do(func() {
		trxApplyServiceInst = &TrxApplyService{roleIntf: roleIntf, ncIntf: nc, curTrxErrorCodeIndex: 0, trxHashErrorMap: make(map[common.Hash]bottosErr.ErrCode)}
	})

	duktape.InitDuktapeVm(roleIntf)

	return trxApplyServiceInst
}

func CreateTempTrxApplyService(roleIntf role.RoleInterface, nc contract.NativeContractInterface) *TrxApplyService {
	trxApplyServiceInst = &TrxApplyService{roleIntf: roleIntf, ncIntf: nc, curTrxErrorCodeIndex: 0, trxHashErrorMap: make(map[common.Hash]bottosErr.ErrCode)}
	duktape.InitDuktapeVm(roleIntf)
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

	systemTime := common.Now()

	//log.Errorf("lifetime %v have past, head time %v system time %v trx hash: %x", time.Unix((int64)(trx.Lifetime), 0),  time.Unix((int64)(curTime), 0), time.Unix((int64)(systemTime), 0), trx.Hash())

	if curTime >= trx.Lifetime {
		log.Errorf("TRX check life time error, have past, trx %x, lifetime %v, head time %v, system time %v", trx.Hash(), time.Unix((int64)(trx.Lifetime), 0), time.Unix((int64)(curTime), 0), time.Unix((int64)(systemTime), 0))
		return false
	}

	if trx.Lifetime >= (curTime + config.DEFAULT_MAX_LIFE_TIME) {
		log.Errorf("TRX check life time error, too far, trx %x, lifetime %v, head time %v, system time %v", trx.Hash(), time.Unix((int64)(trx.Lifetime), 0), time.Unix((int64)(curTime), 0), time.Unix((int64)(systemTime), 0))
		return false
	}

	return true
}

// CheckTransactionUnique is to check whether a transaction is unique
func (trxApplyService *TrxApplyService) CheckTransactionUnique(trx *types.Transaction) bool {

	transactionExpiration, _ := trxApplyService.roleIntf.GetTransactionExpiration(trx.Hash())
	if nil != transactionExpiration {
		log.Errorf("check unique error, trx: %x", trx.Hash())

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
		log.Errorf("process trx error: %v trx: %x", bottosError, trx.Hash())
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

	applyContext := &contract.Context{RoleIntf: trxApplyService.roleIntf, Trx: trx}

	if trxApplyService.ncIntf.IsNativeContract(trx.Contract, trx.Method) {
		bottoserr = trxApplyService.ncIntf.ExecuteNativeContract(applyContext)
		if bottosErr.ErrNoError == bottoserr {
			return true, bottosErr.ErrNoError, nil
		}

		log.Error("process trx, failed bottos error: ", bottoserr)
		return false, bottoserr, nil

	}


	account, _ := trxApplyService.roleIntf.GetAccount(trx.Contract)
	if (vm.VmTypeJS == vm.VmType(account.VMType)) {

		exeErr, trxList := duktape.Process(trx.Contract, account.ContractCode, trx.Method, trx.Param, trx)

		if nil != exeErr {
			log.Error("process trx failed, error: ", exeErr)
			return false, bottosErr.ErrTrxContractHanldeError, nil
		}

		for i, subTrx:= range trxList {
			log.Infof("go in trx apply sub trx:slice[%d] = %v", i, subTrx)

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

		return true, bottosErr.ErrNoError ,derivedTrx
	} else {	
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

func (trxApplyService *TrxApplyService) AddTrxErrorCode(trxHash common.Hash, errCode bottosErr.ErrCode) {

	trxApplyService.mu.Lock()
	defer trxApplyService.mu.Unlock()

	var nextTrxErrorCodeIndex = (trxApplyService.curTrxErrorCodeIndex + 1)%config.DEFAULT_MAX_TRX_ERROR_CODE_NUM
	delete(trxApplyService.trxHashErrorMap, trxApplyService.trxHashErrorList[nextTrxErrorCodeIndex])
	trxApplyService.trxHashErrorMap[trxHash] = errCode

	trxApplyService.trxHashErrorList[nextTrxErrorCodeIndex] = trxHash	

	trxApplyService.curTrxErrorCodeIndex = nextTrxErrorCodeIndex
}


func (trxApplyService *TrxApplyService) GetTrxErrorCode(trxHash common.Hash) bottosErr.ErrCode {

	trxApplyService.mu.Lock()
	defer trxApplyService.mu.Unlock()

	errCode, ok := trxApplyService.trxHashErrorMap[trxHash]
	if ok { 
		return errCode
	}else {
		return bottosErr.ErrNoError
	}
}



func (trxApplyService *TrxApplyService) IsTrxInPendingPool(trxHash common.Hash) bool {
	TrxPoolInst.mu.Lock()
	defer TrxPoolInst.mu.Unlock()	

	_, ok := TrxPoolInst.pending[trxHash]
	if ok { 
		return true
	}else {
		return false
	}
}

//GetAvailableSpace
func (trxApplyService *TrxApplyService) GetAvailableSpace(acc string) (Limit, Limit, error) {
	/*cs, _ := trxApplyService.roleIntf.GetChainState()
	now := cs.LastBlockNum + 1

	var limit Limit
	ufsl, err := GetUserFreeSpaceLimit(trxApplyService.roleIntf, acc, now)
	if err != nil {
		log.Errorf("GetUserFreeSpaceLimit error:%v\n", err)
		return limit, limit, err
	}
	log.Infof("Account:%v, now:%v, userFreeSpaceLimit:%+v", acc, now, ufsl)

	usl, err := GetUserSpaceLimit(trxApplyService.roleIntf, acc, now)
	if err != nil {
		log.Errorf("GetUserFreeSpaceLimit error:%v\n", err)
		return limit, limit, err
	}
	log.Infof("Account:%v, now:%v, userSpaceLimit:%+v", acc, now, usl)

	return ufsl, usl, nil*/
	resService := CreateResProcessorService(trxApplyService.roleIntf)
	f, err := checkMinBalance(resService, acc)
	if err != nil {
		log.Warnf("RESOURCE:checkMinBalance failed:%v", err)
		//return limit, limit, err
	}
	return MaxAvailableSpace(CreateResProcessorService(trxApplyService.roleIntf), acc, f)
}

//GetAvailableTime
func (trxApplyService *TrxApplyService) GetAvailableTime(acc string) (Limit, Limit, error) {
	/*	cs, _ := trxApplyService.roleIntf.GetChainState()
	now := cs.LastBlockNum + 1

	var limit Limit
	ufsl, err := GetUserFreeTimeLimit(trxApplyService.roleIntf, acc, now)
	if err != nil {
		log.Errorf("GetUserFreeTimeLimit error:%v\n", err)
		return limit, limit, err
	}
	log.Infof("Account:%v, now:%v, userFreeTimeLimit:%+v", acc, now, ufsl)

	usl, err := GetUserTimeLimit(trxApplyService.roleIntf, acc, now)
	if err != nil {
		log.Errorf("GetUserTimeLimit error:%v\n", err)
		return limit, limit, err
	}
		log.Infof("Account:%v, now:%v, userTimeLimit:%+v", acc, now, usl)*/

	//return ufsl, usl, nil
	resService := CreateResProcessorService(trxApplyService.roleIntf)
	f, err := checkMinBalance(resService, acc)
	if err != nil {
		log.Warnf("RESOURCE:checkMinBalance failed:%v", err)
		//return limit, limit, err
	}
	return MaxAvailableTime(resService, acc, f)
}
