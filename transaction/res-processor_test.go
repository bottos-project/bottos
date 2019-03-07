package transaction

import (
	"testing"
	"math"
	"github.com/bottos-project/bottos/role"
	"os"
	"github.com/bottos-project/bottos/db"
	"github.com/ontio/ontology/common/log"
	"math/big"
	"github.com/stretchr/testify/assert"
	"encoding/hex"
	"github.com/bottos-project/magiccube/config"
	"bytes"
	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/magiccube/service/common/util"
	"github.com/bottos-project/crypto-go/crypto"
	bottosErr "github.com/bottos-project/bottos/common/errors"
)

var rl role.ResourceLimit
var ru role.ResourceUsage
var b role.Balance
var sb role.StakedBalance
var chainState role.ChainState

var roleIntf role.RoleInterface
var dbInst *db.DBService
var acc = "bob"

func Test_Division_2(t *testing.T) {
	var a, b uint64
	a = 100
	b = 39
	c := float64(a) / float64(b)
	d := a / b
	t.Log("c:", c)
	t.Log("d:", d)
	t.Log("c:", math.Floor(c+0.5))

}

func Test_TxProcessSpaceResouce(t *testing.T) {
	//var now uint64
	//now = 10000
	sb, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", sb)

	spaceTokenCost, spaceUsage, err, be := ProcessSpaceResource(roleIntf, initTx())
	if err != nil {
		t.Logf(err.Error())
	}
	if int(be) != 0 {
		t.Logf(bottosErr.GetCodeString(be))
		//return false, be, nil
	}

	sb2, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", sb2)
	t.Logf("spaceToken:%v", spaceTokenCost)

	timeTokenCost, timeUsage, err, be := ProcessTimeResouce(roleIntf, initTx(), 10000000)
	if err != nil {
		t.Logf(err.Error())
	}
	if int(be) != 0 {
		t.Logf(bottosErr.GetCodeString(be))
		//return false, be, nil
	}

	UpdateResourceUsage(roleIntf, spaceUsage, timeUsage)

	sb2, _ = role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", sb2)
	t.Logf("%+v:", sb2)

	f := assert.ObjectsAreEqual(sb, sb2)
	t.Logf("%v:", f)
	t.Logf("timeToken:%v", timeTokenCost)

	AddResourceReceipt(timeUsage.AccountName, spaceTokenCost, timeTokenCost)

}

func Test_ProcessSpaceResouce(t *testing.T) {
	//var now uint64
	//now = 10000
	sb, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", sb)

	spaceToken, _, err, be := ProcessSpaceResource(roleIntf, initTx())
	if err != nil {
		t.Logf(err.Error())
	}
	if int(be) != 0 {
		t.Logf(bottosErr.GetCodeString(be))
		//return false, be, nil
	}
	t.Logf("spaceToken:%v", spaceToken)

	sb2, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", sb2)
	t.Logf("%+v:", sb2)

	f := assert.ObjectsAreEqual(sb, sb2)
	t.Logf("%v:", f)
	//t.Logf("Para lastUsage: %v, lastTime:%v, now:%v ", lastUsage, lastTime, now)
}

func Test_GetUserSpaceLimit(t *testing.T) {
	var now uint64
	now = 10000

	usl, _ := GetUserSpaceLimit(roleIntf, acc, now)

	//t.Logf("Para lastUsage: %v, lastTime:%v, now:%v ", lastUsage, lastTime, now)

	t.Log("GetUserSpaceLimit result, usl.Max:", usl.Max)
	t.Log("GetUserSpaceLimit result, usl.Available:", usl.Available)
	t.Log("GetUserSpaceLimit result, usl.Used:", usl.Used)
}

func Test_GetUserFreeSpaceLimit(t *testing.T) {
	var now uint64
	now = 10000

	usl, _ := GetUserFreeSpaceLimit(roleIntf, acc, now)

	//t.Logf("Para lastUsage: %v, lastTime:%v, now:%v ", lastUsage, lastTime, now)

	t.Log("GetUserFreeSpaceLimit result, usl.Max:", usl.Max)
	t.Log("GetUserFreeSpaceLimit result, usl.Available:", usl.Available)
	t.Log("GetUserFreeSpaceLimit result, usl.Used:", usl.Used)
}

func Test_ProcessTimeResouce(t *testing.T) {
	//var now uint64
	//now = 10000
	sb, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", sb)

	timeToken, _, err, be := ProcessTimeResouce(roleIntf, initTx(), 10000000)
	if err != nil {
		t.Logf(err.Error())
	}
	if int(be) != 0 {
		t.Logf(bottosErr.GetCodeString(be))
		//return false, be, nil
	}

	sb2, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", sb2)
	t.Logf("%+v:", sb2)

	f := assert.ObjectsAreEqual(sb, sb2)
	t.Logf("%v:", f)
	t.Logf("timeToken:%v", timeToken)
}

func Test_GetUserTimeLimit(t *testing.T) {
	var now uint64
	now = 10000

	usl, _ := GetUserTimeLimit(roleIntf, acc, now)

	//t.Logf("Para lastUsage: %v, lastTime:%v, now:%v ", lastUsage, lastTime, now)

	t.Log("GetUserSpaceLimit result, usl.Max:", usl.Max)
	t.Log("GetUserSpaceLimit result, usl.Available:", usl.Available)
	t.Log("GetUserSpaceLimit result, usl.Used:", usl.Used)
}

func Test_GetUserFreeTimeLimit(t *testing.T) {
	var now uint64
	now = 10000

	usl, _ := GetUserFreeTimeLimit(roleIntf, acc, now)

	//t.Logf("Para lastUsage: %v, lastTime:%v, now:%v ", lastUsage, lastTime, now)

	t.Log("GetUserFreeSpaceLimit result, usl.Max:", usl.Max)
	t.Log("GetUserFreeSpaceLimit result, usl.Available:", usl.Available)
	t.Log("GetUserFreeSpaceLimit result, usl.Used:", usl.Used)
}

func Test_lastUsageNow(t *testing.T) {
	var lastUsage, lastTime, now uint64

	rur, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	lastUsage = rur.PledgedSpaceTokenUsedInWin

	lastTime = rur.LastSpaceCursorBlock
	now = 10000

	t.Logf("Para lastUsage: %v, lastTime:%v, now:%v ", lastUsage, lastTime, now)

	res, _ := lastUsageNow(lastUsage, lastTime, now)
	t.Log("lastUsageNow result:", res)
}

func Test_Add(t *testing.T) {
	var lastUsage, usage, lastTime, now uint64
	lastUsage = 220
	usage = 2
	lastTime = 10
	now = 10000

	res := add(lastUsage, usage, lastTime, now)
	t.Logf("Para lastUsage: %v, lastTime:%v, usage:%v, now:%v ", lastUsage, lastTime, usage, now)

	now += 18700
	res = add(lastUsage, usage, lastTime, now)
	t.Logf("Add Para lastUsage: %v, lastTime:%v, usage:%v, now:%v ", lastUsage, lastTime, usage, now)
	t.Log("UpdateMargin result:", res)
}

func Test_DB(t *testing.T) {
	sb, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", sb)
	usage := role.ResourceUsage{
		AccountName: acc,
		//PledgedSpaceTokenUsedInWin: txSize - ufsl.Available,
		//PledgedTimeTokenUsedInWin:,
		//FreeTimeTokenUsedInWin:,
		//FreeSpaceTokenUsedInWin:    ufsl.Available,
		//LastSpaceCursorBlock:,
		LastTimeCursorBlock: 100,
	}
	role.SetResourceUsageRole(dbInst, acc, &usage)

	ru, _ := role.GetResourceUsageRoleByName(dbInst, acc)
	t.Logf("%v:", ru)
	t.Logf("%#v:", ru)
}

func init() {
	rl = role.ResourceLimit{
		AccountName:                acc,
		PledgedSpaceLimitInWin:     100,
		PledgedTimeTokenLimitInWin: 100,
		FreeTimeTokenLimitInWin:    100,
		FreeSpaceTokenLimitInWin:   100,
	}

	ru = role.ResourceUsage{
		AccountName:                acc,
		PledgedSpaceTokenUsedInWin: 90,
		PledgedTimeTokenUsedInWin:  100,
		FreeTimeTokenUsedInWin:     90,
		FreeSpaceTokenUsedInWin:    90,
		LastSpaceCursorBlock:       10,
		LastTimeCursorBlock:        10,
	}

	aa := new(big.Int).SetUint64(uint64(100000000000))
	allSb := new(big.Int).SetUint64(uint64(1000 * 1000 * 100000000))
	balance := new(big.Int).SetUint64(uint64(1000000000))
	b = role.Balance{
		AccountName: acc,
		Balance:     balance,
	}

	sb = role.StakedBalance{
		AccountName:           acc,
		StakedBalance:         aa,
		StakedTimeBalance:     aa,
		StakedSpaceBalance:    aa,
		UnstakingBalance:      aa,
		LastUnstakingTime:     1542353706,
	}
	chainState = role.ChainState{
		AllStakedSpaceBalance: allSb,
		AllStakedTimeBalance:  allSb,
	}

	initDb()
	role.SetResourceUsageRole(dbInst, acc, &ru)
	role.SetBalanceRole(dbInst, acc, &b)
	role.SetStakedBalanceRole(dbInst, acc, &sb)
	role.SetChainStateRole(dbInst,&chainState)
}

func initDb() {
	blockDBPath := "./data/block/"
	stateDBPath := "./data/state.db"
	dbInst = db.NewDbService(blockDBPath, stateDBPath)
	if dbInst == nil {
		log.Error("Create DB service fail")

		os.Exit(1)
	}

	dbInst.LoadStateDB()
	roleIntf = role.NewRole(dbInst)
}

func Test_TX(t *testing.T) {
	initTx()

}

func initTx() (*types.Transaction) {
	type transfer struct {
		From  string
		To    string
		Value uint64
		Memo  string
	}
	account := &transfer{
		From:  "fdsfgfdsfdsf123fdasf",
		To:    "fdsfgfdsfdsf123fdasf",
		Value: 18446744073709551615,
		Memo:  "fdsfgfdsfdsf123fdfdsas4f123fdasf",
	}
	accountBuf, _ := bpl.Marshal(account)

	txAccountSign := &types.BasicTransaction{
		Version:     1,
		CursorNum:   10000,
		CursorLabel: 100,
		Lifetime:    1220,
		Sender:      "bob",
		Contract:    "bottos",
		Method:      "transfer",
		Param:       accountBuf,
		SigAlg:      1,
	}

	msg, _ := bpl.Marshal(txAccountSign)
	seckey, _ := hex.DecodeString("b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45")

	chainID, _ := hex.DecodeString(config.CHAIN_ID)
	msg = bytes.Join([][]byte{msg, chainID}, []byte{})
	sign, _ := crypto.Sign(util.Sha256(msg), seckey)

	trx := &types.Transaction{
		Version:     txAccountSign.Version,
		CursorNum:   txAccountSign.CursorNum,
		CursorLabel: txAccountSign.CursorLabel,
		Lifetime:    txAccountSign.Lifetime,
		Sender:      txAccountSign.Sender,
		Contract:    txAccountSign.Contract,
		Method:      txAccountSign.Method,
		Param:       accountBuf,
		SigAlg:      txAccountSign.SigAlg,
		Signature:   sign,
	}
	b, _ := bpl.Marshal(trx)
	log.Infof("length Param:%v, sign:%v, tx:%v", len(accountBuf), len(sign), len(b))

	return trx

}
