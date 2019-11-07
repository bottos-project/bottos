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
	"github.com/bottos-project/bottos/common/vm"
	"github.com/bottos-project/bottos/vm/duktape"
)

// TrxApplyService is to define a service for apply a transaction
type TrxApplyService struct {
	roleIntf role.RoleInterface
	ncIntf   contract.NativeContractInterface

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
		log.Errorf("TRX check unique error, trx %x, trx Expiration %x", trx.Hash(), transactionExpiration)

		return false
	}

	return true
}

// CheckTransactionMatchChain is to check whether it is the right chain
func (trxApplyService *TrxApplyService) CheckTransactionMatchChain(trx *types.Transaction) bool {

	blockHistory, err := trxApplyService.roleIntf.GetBlockHistory(trx.CursorNum)
	if nil != err || nil == blockHistory {
		log.Errorf("TRX check chain match error, trx %x, cursor %v", trx.Hash(), trx.CursorNum)
		return false
	}

	var chainCursorLabel uint32 = (uint32)(blockHistory.BlockHash[common.HashLength-1]) + (uint32)(blockHistory.BlockHash[common.HashLength-2])<<8 + (uint32)(blockHistory.BlockHash[common.HashLength-3])<<16 + (uint32)(blockHistory.BlockHash[common.HashLength-4])<<24

	if chainCursorLabel != trx.CursorLabel {
		log.Errorf("TRX check chain match error, trx %x, cursorlabel %v, chain cursolabel %v", trx.Hash(), trx.CursorLabel, chainCursorLabel)
		return false
	}

	return true
}

// SaveTransactionExpiration is to save the expiration of a transaction
func (trxApplyService *TrxApplyService) SaveTransactionExpiration(trx *types.Transaction) {

	var transactionExpiration = &role.TransactionExpiration{TrxHash: trx.Hash(), Expiration: trx.Lifetime}
	trxApplyService.roleIntf.SetTransactionExpiration(trx.Hash(), transactionExpiration)
}

func (trxApplyService *TrxApplyService) ApplyBlockTransaction(trx *types.BlockTransaction) (bool, bottosErr.ErrCode, *types.HandledTransaction, *types.ResourceReceipt) {

	log.Infof("RESOURCE: begin verify trx %x", trx.Transaction.Hash())
	if trx.Transaction.Sender == config.BOTTOS_CONTRACT_NAME {
		f, bErr, h, rr, _ := trxApplyService.ExecuteTransaction(trx.Transaction, false)
		log.Errorf("RESOURCE: execute native transaction SUCESS , trx %x, error %v", trx.Transaction.Hash(), bottosErr.GetCodeString(bErr))
		return f, bErr, h, rr
	}

	//todo
	var resouceReceipt *types.ResourceReceipt
	resService := CreateResProcessorService(trxApplyService.roleIntf)
	f, err := checkMinBalance(resService, trx.Transaction.Sender)
	if err != nil {
		log.Warnf("RESOURCE:checkMinBalance failed:%v", err)
		//return false, bottosErr.ErrTrxResourceCheckMinBalance, nil, resouceReceipt, resUsage
	}

	flag, bErr, _, resReceipt, resUsage := trxApplyService.ExecuteTransaction(trx.Transaction, false)
	if !flag {
		log.Errorf("RESOURCE: execute transaction failed, trx %x, error %v", trx.Transaction.Hash(), bottosErr.GetCodeString(bErr))
		return false, bErr, nil, nil
	}

	_, timeUsage, err, be := ProcessTimeResource(trxApplyService.roleIntf, trx.Transaction, trx.ResourceReceipt.TimeTokenCost, f)
	if err != nil {
		log.Errorf("RESOURCE: process time rsc failed, trx %x, error %v", trx.Transaction.Hash(), err)
		return false, bottosErr.ErrTrxCheckTimeInternalError, nil, resouceReceipt
	}
	if int(be) != 0 {
		log.Errorf("RESOURCE: process time rsc failed, trx %x, error %v", trx.Transaction.Hash(), bottosErr.GetCodeString(be))
		return false, be, nil, resouceReceipt
	}

	resUsage = generateNewUsage(resUsage, timeUsage)

	err = UpdateTimeUsage(trxApplyService.roleIntf, resUsage)
	if err != nil {
		log.Errorf("RESOURCE: update time usage failed, trx %x, error %v", trx.Transaction.Hash(), err)
		return false, bottosErr.ErrTrxCheckResourceInternalError, nil, nil
	}

	//verify space Cost
	if resReceipt.SpaceTokenCost != trx.ResourceReceipt.SpaceTokenCost {
		log.Errorf("RESOURCE: verify space failed, trx %x, nowCost:%v, trxSpaceCost:%v", trx.Transaction.Hash(), resReceipt.SpaceTokenCost, trx.ResourceReceipt.SpaceTokenCost)
		return false, bottosErr.ErrNoError, nil, nil
	}

	//TODO
	//if resReceipt.TimeTokenCost ==timeTokenCost  {
	//	return false, bottosErr.ErrNoError, nil, nil
	//}
	return true, bottosErr.ErrNoError, nil, nil
}

// ExecuteTransaction is to handle a transaction, include parameters checking
func (trxApplyService *TrxApplyService) ExecuteTransaction(trx *types.Transaction, verifyTimeFlag bool) (bool, bottosErr.ErrCode, *types.HandledTransaction, *types.ResourceReceipt, role.ResourceUsage) {
	start := common.MeasureStart()
	log.Infof("TRX begin exec trx, trx %x", trx.Hash())
	var resouceReceipt *types.ResourceReceipt
	var resUsage role.ResourceUsage

	if !trxApplyService.CheckTransactionUnique(trx) {
		return false, bottosErr.ErrTrxUniqueError, nil, resouceReceipt, resUsage
	}

	account, getAccountErr := trxApplyService.roleIntf.GetAccount(trx.Sender)
	if nil != getAccountErr || nil == account {
		log.Errorf("TRX exec trx get account error, trx %x", trx.Hash())
		return false, bottosErr.ErrTrxAccountError, nil, resouceReceipt, resUsage
	}

	if !trxApplyService.CheckTransactionLifeTime(trx) {
		return false, bottosErr.ErrTrxLifeTimeError, nil, resouceReceipt, resUsage
	}

	if !trxApplyService.CheckTransactionMatchChain(trx) {
		return false, bottosErr.ErrTrxChainMathError, nil, resouceReceipt, resUsage
	}

	trxApplyService.SaveTransactionExpiration(trx)

	resService := CreateResProcessorService(trxApplyService.roleIntf)

	f, err := checkMinBalance(resService, trx.Sender)
	if err != nil {
		log.Warnf("RESOURCE:checkMinBalance failed:%v", err)
		//return false, bottosErr.ErrTrxResourceCheckMinBalance, nil, resouceReceipt, resUsage
	}

	if trx.Sender != config.BOTTOS_CONTRACT_NAME {
		resConfig, err := trxApplyService.roleIntf.GetResourceConfig()
		if err != nil {
			log.Errorf("RESOURCE:get Resource Config failed,", err)
		}
		_, berr := GetTxSize(*resConfig, trx, 0)

		if int(berr) != 0 {
			log.Errorf("RESOURCE: check Process Space Resource failed:%v", bottosErr.GetCodeString(berr))
			return false, berr, nil, resouceReceipt, resUsage
		}
	}

	//max available time of user
	maxTime, err := MaxContractExecuteTime(trxApplyService.roleIntf, trx.Sender, f)
	if err != nil {
		log.Errorf("RESOURCE: get max time failed, trx: %x, err:%v", trx.Hash(), err)
		return false, bottosErr.ErrTrxCheckTimeInternalError, nil, resouceReceipt, resUsage
	}
	if maxTime < config.CONTRACT_EXEC_MIN_TIME {
		if (trx.Contract == config.BOTTOS_CONTRACT_NAME) && (trx.Method == "stake") {
		maxTime = config.CONTRACT_EXEC_MIN_TIME
		} else {
			log.Errorf("RESOURCE: max timeToken less than min required, maxTime: %v", maxTime)
			return false, bottosErr.ErrTrxCheckMinTimeError, nil, resouceReceipt, resUsage
		}
	}
	
	applyContext := &contract.Context{
		RoleIntf: trxApplyService.roleIntf, 
		Trx: trx, 
		CallContract: trx.Contract,
		CallMethod: trx.Method,
		DeepLimit:0, 
		MaxExecTime:maxTime}
	
	bottosError, derivedTrxList, space, execTime:= trxApplyService.ProcessTransaction(applyContext)

	if bottosErr.ErrNoError != bottosError {
		//log.Errorf("TRX process trx error, trx %x, error %v", trx.Hash(), bottosError)
		return false, bottosError, nil, resouceReceipt, resUsage
	}

	log.Infof("TRX contract exec succ, space is %v, time %v ", space, execTime)
	handleTrx := &types.HandledTransaction{
		Transaction: trx,
		DerivedTrx:  derivedTrxList,
	}

	log.Infof("TRX exec trx success, trx %x, elapsed time %v", trx.Hash(), common.Elapsed(start))
	return true, bottosErr.ErrNoError, handleTrx, resouceReceipt, resUsage
}

// ProcessTransaction is to handle a transaction without parameters checking
func (trxApplyService *TrxApplyService) ProcessTransaction(applyContext *contract.Context) (bottosErr.ErrCode, []*types.DerivedTransaction, uint64, uint64) {

	log.Debugf("TRX begein exec contract %v, deep %v", applyContext.CallContract,  applyContext.DeepLimit)

	if applyContext.DeepLimit >= config.DEFAUL_MAX_CONTRACT_DEPTH {
		log.Errorf("TRX exec trx failed, deep over %v",applyContext.DeepLimit)
		return bottosErr.ErrTrxContractDepthError, nil, 0, 0
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
			Transaction: applyContext.Trx,
			DerivedTrx:  subDerivedTrx,
		}

		derivedTrx = append(derivedTrx, handleTrx)
	}

	log.Debugf("TRX begin exec sub trx list:")	
	
	log.Debugf("TRX sub trx list:")
	for _, subTrx := range trxList {
		log.Debug(subTrx)
	}
	
	for _, subTrx := range trxList {
		subApplyContext := &contract.Context{
			RoleIntf: trxApplyService.roleIntf,
			Trx: subTrx, 
			CallContract: subTrx.Contract,
		    CallMethod: subTrx.Method,
			DeepLimit:applyContext.DeepLimit + 1,  
			MaxExecTime:applyContext.MaxExecTime - totalExecTime}

		subTrxErr, subDerivedTrx, subTrxDbDataSaveLen, subExecTime := trxApplyService.ProcessTransaction(subApplyContext)
		if bottosErr.ErrNoError != subTrxErr {
			log.Errorf("TRX exec sub trx failed, trx %x, error %v, sub trx %v, method %v,", applyContext.Trx.Hash(), subTrxErr, subTrx.Contract, subTrx.Method)
			return subTrxErr, nil, 0, 0
		}

		log.Infof("TRX exec sub trx, contract %v, db save len %v, exec time %v",subTrx.Contract, subTrxDbDataSaveLen, subExecTime)

		totalExecTime += subExecTime

		totalDbDataSaveLen += subTrxDbDataSaveLen

		if totalExecTime > applyContext.MaxExecTime {
			log.Errorf("TRX exec sub trx failed, trx %x, exec time over", applyContext.Trx.Hash())
			return bottosErr.ErrTrxExecTimeOver, nil, 0, 0
		}		

		handleTrx := &types.DerivedTransaction{
			Transaction: subTrx,
			DerivedTrx:  subDerivedTrx,
		}

		derivedTrx = append(derivedTrx, handleTrx)
	}

	log.Debugf("TRX exec contract %v done, deep %v", applyContext.CallContract, applyContext.DeepLimit)

	return bottosErr.ErrNoError, derivedTrx, totalDbDataSaveLen, totalExecTime
}

func (trxApplyService *TrxApplyService) AddTrxErrorCode(trxHash common.Hash, errCode bottosErr.ErrCode) {

	trxApplyService.mu.Lock()
	defer trxApplyService.mu.Unlock()

	var nextTrxErrorCodeIndex = (trxApplyService.curTrxErrorCodeIndex + 1) % config.DEFAULT_MAX_TRX_ERROR_CODE_NUM
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
	} else {
		return bottosErr.ErrNoError
	}
}

func (trxApplyService *TrxApplyService) IsTrxInPendingPool(trxHash common.Hash) bool {
	TrxPoolInst.mu.Lock()
	defer TrxPoolInst.mu.Unlock()

	_, ok := TrxPoolInst.pending[trxHash]
	if ok {
		return true
	} else {
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
