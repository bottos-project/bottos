package transaction

import (
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/role"
	"math/big"
	"math"
	log "github.com/cihub/seelog"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/bpl"
	bottosErr "github.com/bottos-project/bottos/common/errors"
	"errors"
	"fmt"
	"github.com/bottos-project/bottos/common/safemath"
	"strings"
)

// TrxApplyService is to define a service for apply a transaction
type ResProcessorService struct {
	roleIntf  role.RoleInterface
	resConfig role.ResourceConfig
}

type Limit struct {
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
	Max       uint64 `json:"max"`
}

var resProcessorServiceInst *ResProcessorService

//var resConfig *role.ResourceConfig
var blockSizePerWindow uint64

// CreateResProcessorService is to new a Resource ProcessorService
func CreateResProcessorService(roleIntf role.RoleInterface) *ResProcessorService {
	resConfig, err := roleIntf.GetResourceConfig()
	if err != nil {
		log.Errorf("RESOURCE:get Resource Config failed,", err)
	}
	resProcessorServiceInst = &ResProcessorService{roleIntf: roleIntf, resConfig: *resConfig}

	blockSizePerWindow = resConfig.BlockSizeAverageWindows

	log.Infof("RESOURCE:get Resource Config:%+v\n", resConfig)

	return resProcessorServiceInst
}

//Space Module
//ProcessSpaceResource Get current time free space resource usage
func ProcessSpaceResource(roleIntf role.RoleInterface, trx *types.Transaction, dbUseSize uint64, checkMinBala bool) (uint64, role.ResourceUsage, error, bottosErr.ErrCode) {
	resService := CreateResProcessorService(roleIntf)

	cs, _ := resService.roleIntf.GetChainState()
	now := cs.LastBlockNum + 1

	acc := trx.Sender
	var usage role.ResourceUsage
	usagee, _ := roleIntf.GetResourceUsage(trx.Sender)
	usage = *usagee

	txSize, bErr := GetTxSize(resService.resConfig, trx, dbUseSize)
	if bErr != bottosErr.ErrNoError {
		log.Errorf("RESOURCE:GetTxSize error:%v", bottosErr.GetCodeString(bErr))
		return 0, usage, nil, bErr
	}

	ufsl, usl, err := MaxAvailableSpace(resService, acc, checkMinBala)
	if err != nil {
		log.Errorf("RESOURCE:GetMaxAvailableSpace failed:%v", err)
		return 0, usage, err, bottosErr.ErrNoError
	}

	if ufsl.Available+usl.Available < resService.resConfig.MinAccountSpaceBalance {
		log.Errorf("RESOURCE:check min SpaceBalance failed,Account:%v,SpaceBalance:%v,min SpaceBlance:%v", acc, ufsl.Available+usl.Available, resService.resConfig.MinAccountSpaceBalance)
		return 0, usage, nil, bottosErr.ErrTrxCheckMinSpaceError
	}

	if ufsl.Available+usl.Available < txSize {
		log.Errorf("RESOURCE:CheckSpace failed,Account:%v, space size:%v, userSpaceAllAvailable:%v, error:%v", acc, txSize, ufsl.Available+usl.Available, bottosErr.GetCodeString(bottosErr.ErrTrxCheckSpaceError))
		return 0, usage, nil, bottosErr.ErrTrxCheckSpaceError
	}

	log.Infof("RESOURCE:GetUserSpaceLimit,trx size:%v, userSpaceLimit:%v", txSize, usl)

	if txSize <= usl.Available {
		usage.AccountName = acc
		//usage.FreeSpaceTokenUsedInWin = 0
		usage.PledgedSpaceTokenUsedInWin = txSize + usl.Used
		usage.LastSpaceCursorBlock = now

		return txSize, usage, err, bottosErr.ErrNoError
	} else {
		log.Infof("RESOURCE:getUserFreeSpaceLimit,Account:%v, now:%v, space size:%v, Limit:%v", acc, now, txSize, ufsl)

		usage.AccountName = acc
		usage.FreeSpaceTokenUsedInWin = txSize - usl.Available + ufsl.Used
		usage.PledgedSpaceTokenUsedInWin = usl.Max
		usage.LastSpaceCursorBlock = now

		return txSize, usage, nil, bottosErr.ErrNoError
	}
}

//GetUserSpaceLimit Get current time resource usage
func GetUserSpaceLimit(resService *ResProcessorService, sender string, now uint64) (Limit, error) {
	var limit Limit
	sb, err := resService.roleIntf.GetStakedBalance(sender)
	if err != nil {
		log.Errorf("RESOURCE:GetStakedBalance failed:%v", err)
		return limit, err
	}

	//Get all user use in window
	allUserUseInWindow := resService.resConfig.MaxSpacePerWindow
	userStakedSpaceBalance := sb.StakedSpaceBalance
	allStakedSpaceBalance, err := getallStakedSpaceBalance(resService.roleIntf)
	if err != nil {
		log.Errorf("RESOURCE:getAllStakedSpaceBalance failed:%v", err)
		return limit, err
	}

	//Get the user's Max Space size
	maxUserUseInWindow := big.NewInt(0)

	if i := allStakedSpaceBalance.Cmp(maxUserUseInWindow); i > 0 {
		result := new(big.Int).Mul(new(big.Int).SetUint64(allUserUseInWindow), userStakedSpaceBalance)
		maxUserUseInWindow, err = sb.SafeDivide(result, allStakedSpaceBalance)
		if err != nil {
			log.Errorf("RESOURCE:getallStakedSpaceBalance SafeDivide failed:%v", err)
			return limit, err
		}
	}

	//Get user Usage
	usage, _ := resService.roleIntf.GetResourceUsage(sender)

	usageNow, err := lastUsageNow(usage.PledgedSpaceTokenUsedInWin, usage.LastSpaceCursorBlock, now)
	if err != nil {
		log.Errorf("RESOURCE:lastUsageNow failed:%v", err)
		return limit, err
	}

	limit.Max = maxUserUseInWindow.Uint64()
	limit.Used = usageNow
	if maxUserUseInWindow.Uint64() == 0 {
		limit.Available = 0
	} else {
		limit.Available = maxUserUseInWindow.Uint64() - usageNow
	}
	log.Infof("RESOURCE:account:%v, space limit:%+v", sender, limit)

	return limit, nil
}

//GetUserFreeSpaceLimit Get current time free resource usage
func GetUserFreeSpaceLimit(resService *ResProcessorService, sender string, now uint64) (Limit, error) {
	var limit Limit

	//Get user Usage
	usage, err := resService.roleIntf.GetResourceUsage(sender)
	if err != nil {
		log.Errorf("DB:GetResourceUsage failed:%v", err)
		return limit, err
	}

	usageNow, err := lastUsageNow(usage.FreeSpaceTokenUsedInWin, usage.LastSpaceCursorBlock, now)
	if err != nil {
		log.Errorf("RESOURCE:lastUsageNow failed:%v", err)
		return limit, err
	}

	limit.Available = resService.resConfig.FreeSpaceTokenPerWindow - usageNow
	limit.Max = resService.resConfig.FreeSpaceTokenPerWindow
	limit.Used = usageNow
	log.Infof("RESOURCE: account:%v, free space limit:%+v", sender, limit)

	return limit, nil
}

//GetTxSize Get current trx size
func GetTxSize(resConfig role.ResourceConfig, trx *types.Transaction, dbUseSize uint64) (uint64, bottosErr.ErrCode) {
	tb, err := bpl.Marshal(trx)
	if err != nil {
		log.Errorf("RESOURCE: marshal failed:%v", err)
		return 0, bottosErr.RestErrBplMarshal
	}
	log.Infof("RESOURCE:Account is:%s, Trx Size:%v,dbUseSize:%v", trx.Sender, len(tb), dbUseSize)

	l := uint64(len(tb)) + dbUseSize
	if strings.Contains(trx.Method, "deploy") {
		if l > resConfig.MaxSpacePerDeployTrx {
			log.Errorf("RESOURCE: Check max space per Deploy trx failed,Account:%s,Method:%v, Trx Size:%v, maxSizePerDeployTrx:%v", trx.Sender, trx.Method, l, resConfig.MaxSpacePerDeployTrx)
			return l, bottosErr.ErrTrxResourceExceedMaxSpacePerTrx
		}
		return l, bottosErr.ErrNoError
	}

	if l > resConfig.MaxSpacePerTrx {
		log.Errorf("RESOURCE: Check max space per trx failed,Account:%s, Trx Size:%v,maxSizePerTrx:%v", trx.Sender, l, resConfig.MaxSpacePerTrx)
		return l, bottosErr.ErrTrxResourceExceedMaxSpacePerTrx
	}

	return l, bottosErr.ErrNoError
}

//UpdateResourceUsage
func UpdateResourceUsage(roleIntf role.RoleInterface, timeUsage role.ResourceUsage) error {
	err := roleIntf.SetResourceUsage(timeUsage.AccountName, &timeUsage)
	if err != nil {
		log.Errorf("DB: set Resource Usage failed:%v", err)
		return err
	}
	return nil
}

//generateNewUsage
func generateNewUsage(spaceUsage, timeUsage role.ResourceUsage) role.ResourceUsage {
	log.Debugf("RESOURCE:generateNewUsage,spaceUsage:%+v,timeUsage:%+v", spaceUsage, timeUsage)
	timeUsage.FreeSpaceTokenUsedInWin = spaceUsage.FreeSpaceTokenUsedInWin
	timeUsage.PledgedSpaceTokenUsedInWin = spaceUsage.PledgedSpaceTokenUsedInWin
	timeUsage.LastSpaceCursorBlock = spaceUsage.LastSpaceCursorBlock

	return timeUsage
}

//getallStakedSpaceBalance
func getallStakedSpaceBalance(roleIntf role.RoleInterface) (*big.Int, error) {
	sb, err := roleIntf.GetChainState()
	if err != nil {
		log.Errorf("DB: get ChainState failed:%v", err)
		return nil, err
	}

	if sb.AllStakedSpaceBalance != nil {
		return sb.AllStakedSpaceBalance, nil
	} else {
		sb.AllStakedSpaceBalance = big.NewInt(0)
		return sb.AllStakedSpaceBalance, nil
	}
}

//MaxAvailableSpace
func MaxAvailableSpace(resService *ResProcessorService, acc string, checkMinFlag bool) (Limit, Limit, error) {
	var limit Limit

	cs, err := resService.roleIntf.GetChainState()
	if err != nil {
		log.Errorf("RESOURCE: GetChainState failed:%v", err)
		return limit, limit, err
	}
	now := cs.LastBlockNum + 1

	ufsl, err := GetUserFreeSpaceLimit(resService, acc, now)
	if err != nil {
		log.Errorf("RESOURCE: GetUserFreeSpaceLimit failed:%v", err)
		return limit, limit, err
	}
	log.Debugf("RESOURCE:Account:%v, now:%v, userFreeSpaceLimit:%+v", acc, now, ufsl)

	if !checkMinFlag {
		ufsl.Available = 0
	}

	usl, err := GetUserSpaceLimit(resService, acc, now)
	if err != nil {
		log.Errorf("RESOURCE: GetUserSpaceLimit failed:%v", err)
		return limit, limit, err
	}

	log.Debugf("RESOURCE:Account:%v, now:%v, userSpaceLimit:%+v", acc, now, usl)
	return ufsl, usl, nil
}

//Time Module
//ProcessTimeResource Get current time free Time resource usage
func ProcessTimeResource(roleIntf role.RoleInterface, trx *types.Transaction, txTime uint64, checkMinBala bool) (uint64, role.ResourceUsage, error, bottosErr.ErrCode) {
	log.Debugf("RESOURCE:Account is:%s,  Trx Exec time:%v", trx.Sender, txTime)
	var usage role.ResourceUsage
	usagee, _ := roleIntf.GetResourceUsage(trx.Sender)
	usage = *usagee
	resService := CreateResProcessorService(roleIntf)

	if txTime < resService.resConfig.ContractExecMinTime {
		txTime = resService.resConfig.ContractExecMinTime
	}

	cs, err := resService.roleIntf.GetChainState()
	if err != nil {
		log.Errorf("DB: get ChainState failed:%v", err)
		return 0, usage, err, bottosErr.ErrNoError
	}

	now := cs.LastBlockNum + 1
	acc := trx.Sender

	ufsl, usl, err := MaxAvailableTime(resService, acc, checkMinBala)
	if err != nil {
		log.Errorf("RESOURCE:GetMaxAvailableTime failed:%v", err)
		return 0, usage, err, bottosErr.ErrNoError
	}

	if ufsl.Available+usl.Available <= resService.resConfig.MinAccountTimeBalance {
		log.Errorf("RESOURCE:check min TimeBalance failed,Account:%v,TimeBalance:%v,min TimeBlance:%v", acc, ufsl.Available+usl.Available, resService.resConfig.MinAccountTimeBalance)
		return 0, usage, nil, bottosErr.ErrTrxCheckTimeError
	}

	if ufsl.Available+usl.Available < txTime {
		log.Errorf("RESOURCE:CheckTime failed, Account:%v, time size:%v, TimeAllAvailable:%v, error:%v", acc, txTime, ufsl.Available+usl.Available, bottosErr.GetCodeString(bottosErr.ErrTrxCheckTimeError))
		return 0, usage, nil, bottosErr.ErrTrxCheckTimeError
	}

	//ufsl, err := GetUserFreeTimeLimit(roleIntf, acc, now)
	//if err != nil {
	//	log.Errorf("RESOURCE: GetUserFreeTimeLimit failed:%v", err)
	//	return 0, usage, err, bottosErr.ErrNoError
	//}
	log.Infof("RESOURCE:GetUserTimeLimit Account:%v, duration time:%v, limit:%+v", acc, txTime, usl)

	if txTime <= usl.Available {
		usage.AccountName = acc
		//usage.FreeTimeTokenUsedInWin = 0
		usage.PledgedTimeTokenUsedInWin = txTime + usl.Used
		usage.LastTimeCursorBlock = now
		return txTime, usage, err, bottosErr.ErrNoError
	} else {
		log.Infof("RESOURCE:getUserFreeTimeLimit Account:%v, now:%v, duration time:%v, limit:%+v", acc, now, txTime, ufsl)

		usage.AccountName = acc
		usage.FreeTimeTokenUsedInWin = txTime - usl.Available + ufsl.Used
		usage.PledgedTimeTokenUsedInWin = usl.Max
		usage.LastTimeCursorBlock = now

		return txTime, usage, nil, bottosErr.ErrNoError
	}
}

//MaxAvailableTime
func MaxAvailableTime(resService *ResProcessorService, acc string, checkMinFlag bool) (Limit, Limit, error) {
	var limit Limit

	if acc == config.BOTTOS_CONTRACT_NAME {
		return limit, limit, nil
	}

	cs, err := resService.roleIntf.GetChainState()
	if err != nil {
		log.Errorf("RESOURCE: GetChainState failed:%v", err)
		return limit, limit, err
	}

	now := cs.LastBlockNum + 1

	ufsl, err := GetUserFreeTimeLimit(resService, acc, now)
	if err != nil {
		log.Errorf("RESOURCE: GetUserFreeTimeLimit failed:%v", err)
		return limit, limit, err
	}

	if !checkMinFlag {
		ufsl.Available = 0
	}
	//f, err := checkMinBalance(resService, acc)
	//if err != nil {
	//	log.Errorf("RESOURCE:checkMinBalance failed:%v", err)
	//	return limit, limit, err
	//}
	//if !f {
	//	ufsl.Available = 0
	//}

	usl, err := GetUserTimeLimit(resService, acc, now)
	if err != nil {
		log.Errorf("RESOURCE: GetUserTimeLimit failed:%v", err)
		return limit, limit, err
	}

	return ufsl, usl, nil
}

//GetUserTimeLimit Get current time resource usage
func GetUserTimeLimit(resService *ResProcessorService, sender string, now uint64) (Limit, error) {
	var limit Limit
	sb, err := resService.roleIntf.GetStakedBalance(sender)
	if err != nil {
		log.Errorf("DB: GetUserTimeLimit failed:%v", err)
		return limit, err
	}
	//Get all user use in window
	allUserUseInWindow := resService.resConfig.MaxTimePerWindow
	userStakedTimeBalance := sb.StakedTimeBalance
	allStakedTimeBalance, err := getallStakedTimeBalance(resService.roleIntf)
	if err != nil {
		log.Errorf("RESOURCE: get allStakedTimeBalance failed:%v", err)
		return limit, err
	}

	//Get the user's Max Time size
	maxUserUseInWindow := big.NewInt(0)
	if i := userStakedTimeBalance.Cmp(maxUserUseInWindow); i > 0 {
		result := new(big.Int).Mul(new(big.Int).SetUint64(allUserUseInWindow), userStakedTimeBalance)
		maxUserUseInWindow, err = sb.SafeDivide(result, allStakedTimeBalance)
		if err != nil {
			log.Errorf("RESOURCE: Divide failed:%v", err)
			return limit, err
		}
	}

	//Get user Usage
	usage, _ := resService.roleIntf.GetResourceUsage(sender)

	usageNow, err := lastUsageNow(usage.PledgedTimeTokenUsedInWin, usage.LastTimeCursorBlock, now)
	if err != nil {
		log.Errorf("RESOURCE:lastUsageNow failed:%v", err)
		return limit, err
	}

	limit.Max = maxUserUseInWindow.Uint64()
	limit.Used = usageNow
	if maxUserUseInWindow.Uint64() == 0 {
		limit.Available = 0
	} else {
		limit.Available = maxUserUseInWindow.Uint64() - usageNow
	}

	log.Infof("RESOURCE:GetUserTimeLimit account:%v, limit:%+v", sender, limit)
	return limit, nil
}

//GetUserFreeTimeLimit Get current time free resource usage
func GetUserFreeTimeLimit(resService *ResProcessorService, sender string, now uint64) (Limit, error) {
	var limit Limit

	usage, err := resService.roleIntf.GetResourceUsage(sender)
	if err != nil {
		log.Errorf("DB:GetResourceUsage failed:%v", err)
		return limit, err
	}

	usageNow, err := lastUsageNow(usage.FreeTimeTokenUsedInWin, usage.LastTimeCursorBlock, now)
	if err != nil {
		log.Errorf("RESOURCE:lastUsageNow failed:%v", err)
		return limit, err
	}

	limit.Available = resService.resConfig.FreeTimeTokenPerWindow - usageNow
	limit.Max = resService.resConfig.FreeTimeTokenPerWindow
	limit.Used = usageNow

	log.Infof("RESOURCE:GetUserFreeTime account:%v, limit:%+v", sender, limit)
	return limit, nil
}

//checkMinBalance
func checkMinBalance(resService *ResProcessorService, sender string) (bool, error) {
	tb, err := GetTotalBalance(resService.roleIntf, sender)
	if err != nil {
		log.Errorf("RESOURCE:GetTotalBalance failed:%v", err)
		return false, err
	}

	i := new(big.Int).SetUint64(resService.resConfig.MinAccountBalance)
	if tb.Cmp(i) < 0 {
		log.Warnf("RESOURCE:checkMinBalance failed,Account:%v,balance:%v,minBlance:%v", sender, tb, resService.resConfig.MinAccountBalance)
		return false, nil
	}
	return true, nil
}

//GetTotalBalance Get Total Balance
func GetTotalBalance(role role.RoleInterface, sender string) (*big.Int, error) {
	b, err := role.GetBalance(sender)
	if err != nil {
		log.Errorf("DB:GetBalance failed:%v", err)
		return nil, err
	}

	sb, err := role.GetStakedBalance(sender)
	if err != nil {
		log.Errorf("DB:GetStakedBalance failed:%v", err)
		return nil, err
	}

	result := big.NewInt(0)
	if sender == config.BOTTOS_CONTRACT_NAME {
		return b.Balance, nil
	}

	result, err = safemath.U256Add(result, sb.StakedBalance, sb.StakedSpaceBalance)
	result, err = safemath.U256Add(result, result, sb.StakedTimeBalance)
	result, err = safemath.U256Add(result, result, sb.UnstakingBalance)
	result, err = safemath.U256Add(result, result, b.Balance)

	if err != nil {
		log.Errorf("RESOURCE:getTotalBalance U256Add failed:%v", err)
		return result, err
	}

	log.Infof("RESOURCE:TotalBalance:%v", result)
	return result, nil
}

//UpdateTimeUsage
func UpdateTimeUsage(roleIntf role.RoleInterface, usage role.ResourceUsage) (error) {
	//ru, err := roleIntf.GetResourceUsage(usage.AccountName)
	//if err != nil {
	//	return err
	//}

	/*	usage = role.ResourceUsage{
			AccountName:                usage.AccountName,
			PledgedSpaceTokenUsedInWin: usage.PledgedSpaceTokenUsedInWin,
			PledgedTimeTokenUsedInWin:  usage.PledgedTimeTokenUsedInWin,
			FreeTimeTokenUsedInWin:     usage.FreeTimeTokenUsedInWin,
			FreeSpaceTokenUsedInWin:    usage.FreeSpaceTokenUsedInWin,
			LastSpaceCursorBlock:       ru.LastSpaceCursorBlock,
			LastTimeCursorBlock:        ru.LastTimeCursorBlock,
		}*/

	err := roleIntf.SetResourceUsage(usage.AccountName, &usage)
	if err != nil {
		log.Errorf("DB:SetResourceUsage failed:%v", err)
		return err
	}
	log.Debugf("RESOURCE:UpdateTime usage:%+v", usage)
	return nil
}

//Common Module
//AddResourceReceipt
func AddResourceReceipt(account string, spaceTokenCost, timeTokenCost uint64) *types.ResourceReceipt {
	resourceReceipt := &types.ResourceReceipt{
		AccountName:    account,
		SpaceTokenCost: spaceTokenCost,
		TimeTokenCost:  timeTokenCost,
	}
	log.Infof("RESOURCE:resource receipt::%v", resourceReceipt)
	return resourceReceipt
}

//getallStakedTimeBalance
func getallStakedTimeBalance(roleIntf role.RoleInterface) (*big.Int, error) {
	sb, err := roleIntf.GetChainState()
	if err != nil {
		log.Errorf("DB:GetChainState failed:%v", err)
		return nil, err
	}

	if sb.AllStakedTimeBalance != nil {
		return sb.AllStakedTimeBalance, nil
	} else {
		sb.AllStakedTimeBalance = big.NewInt(0)
		return sb.AllStakedTimeBalance, nil
	}
}

//lastUsageNow
func lastUsageNow(lastUsage, lastTime, now uint64) (uint64, error) {
	//averageLastUsage := divideCeil(lastUsage*uint64(config.RATE_LIMITING_PRECISION), blockSizePerWindow)
	averageLastUsage := lastUsage

	if lastTime > now {
		err := fmt.Sprintf("The last time lags behind the current time. lastTime:%v, now:%v", lastTime, now)
		log.Errorf("RESOURCE:check Last time failed :%v", err)
		return 0, errors.New(err)
	}
	if (lastTime + blockSizePerWindow) > now {
		delta := now - lastTime
		decay := float64(blockSizePerWindow-delta) / float64(blockSizePerWindow)
		//decay := fmt.Sprintf("%0.6f", (float64(blockSizePerWindow - delta) / float64(blockSizePerWindow)))

		f := float64(averageLastUsage) * decay
		n10 := math.Pow10(0)
		a := math.Trunc((f+0.5/n10)*n10) / n10

		averageLastUsage = uint64(a)
	} else {
		averageLastUsage = 0
	}
	return averageLastUsage, nil
}

//add
func add(lastUsage, usage, lastTime, now uint64) uint64 {
	//TODO
	//averageLastUsage := divideCeil(lastUsage*uint64(config.RATE_LIMITING_PRECISION), blockSizePerWindow)
	//averageUsage := divideCeil(usage*uint64(config.RATE_LIMITING_PRECISION), blockSizePerWindow)
	//averageLastUsage:=lastUsage
	//averageUsage:=usage

	lastUsageNow, err := lastUsageNow(lastUsage, lastTime, now)
	if err != nil {
		log.Errorf("RESOURCE:lastUsageNow failed :%v", err)
	}

	//return getUsage(averageLastUsage)
	return lastUsageNow + usage
}

//divideCeil
func divideCeil(numerator, denominator uint64) uint64 {
	if (numerator % denominator) > 0 {
		return (numerator / denominator) + 1
	} else {
		return numerator / denominator
	}
}

//getUsage
func getUsage(usage uint64) uint64 {
	//return usage * blockSizePerWindow / resConfig.RateLimitingPrecision
	return 0
}

//MaxContractExecuteTime
func MaxContractExecuteTime(roleIntf role.RoleInterface, acc string, checkMinBala bool) (uint64, error) {
	resService := CreateResProcessorService(roleIntf)
	if acc == config.BOTTOS_CONTRACT_NAME {
		return resService.resConfig.ContractExecMaxTime, nil
	}
	ufsl, usl, err := MaxAvailableTime(resService, acc, checkMinBala)

	if err != nil {
		log.Errorf("RESOURCE:MaxAvailableTime failed,account:%s, %v", acc, err)
		return 0, err
	}

	if ufsl.Available+usl.Available > resService.resConfig.ContractExecMaxTime {
		return resService.resConfig.ContractExecMaxTime, nil
	}

	return ufsl.Available + usl.Available, nil
}

func (resService *ResProcessorService) ForeCastRes(userStakedSpaceBalance, userStakedTimeBalance *big.Int) (*big.Int, *big.Int, error) {
	calcUserSpace := big.NewInt(0)
	calcUserTime := big.NewInt(0)
	allStakeSpaceTmp := big.NewInt(0)
	allStakeTimeTmp := big.NewInt(0)

	allUserUseSpaceInWindow := resService.resConfig.MaxSpacePerWindow
	allStakedSpaceBalance, err := getallStakedSpaceBalance(resService.roleIntf)
	if err != nil {
		log.Errorf("RESOURCE: get allStakedTimeBalance failed:%v", err)
		return calcUserSpace, calcUserTime, err
	}
	allStakeSpaceTmp, err = safemath.U256Add(allStakeSpaceTmp, userStakedSpaceBalance, allStakedSpaceBalance)
	if err != nil {
		log.Errorf("RESOURCE: U256Add failed:%v", err)
		return calcUserSpace, calcUserTime, err
	}

	calcUserSpace, err = calcStakeResBala(userStakedSpaceBalance, allStakeSpaceTmp, allUserUseSpaceInWindow)
	if err != nil {
		log.Errorf("RESOURCE: calc Stake Resource Balance failed:%v", err)
		return calcUserSpace, calcUserTime, err
	}

	//calculate stake ==> Time resource
	allUserUseTimeInWindow := resService.resConfig.MaxTimePerWindow
	allStakedTimeBalance, err := getallStakedTimeBalance(resService.roleIntf)
	if err != nil {
		log.Errorf("RESOURCE: get allStakedTimeBalance failed:%v", err)
		return calcUserSpace, calcUserTime, err
	}
	allStakeTimeTmp, err = safemath.U256Add(allStakeTimeTmp, userStakedTimeBalance, allStakedTimeBalance)
	if err != nil {
		log.Errorf("RESOURCE: U256Add failed:%v", err)
		return calcUserSpace, calcUserTime, err
	}

	calcUserTime, err = calcStakeResBala(userStakedTimeBalance, allStakeTimeTmp, allUserUseTimeInWindow)
	return calcUserSpace, calcUserTime, err
}

func calcStakeResBala(userStakedBalance, globalStakeBala *big.Int, globalUseInWindow uint64) (*big.Int, error) {
	calcRes := big.NewInt(0)
	if i := userStakedBalance.Cmp(calcRes); i > 0 {
		result := new(big.Int).Mul(new(big.Int).SetUint64(globalUseInWindow), userStakedBalance)
		calcRes, err := role.SafeDivide(result, globalStakeBala)
		if err != nil {
			log.Errorf("RESOURCE: Divide failed:%v", err)
			return calcRes, err
		}
		log.Infof("user stake:%v,globalStake:%v,globalUseIn Window:%v", userStakedBalance, globalStakeBala, globalUseInWindow)
		return calcRes, nil
	}
	return calcRes, nil
}
