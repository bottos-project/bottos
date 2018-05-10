package contract

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/role"
)

func NewNativeContract(roleIntf role.RoleInterface) (NativeContractInterface, error) {
	intf, err := NewNativeContractHandler()
	if err != nil {
		return nil, err
	}
	roleIntf.SetScheduleDelegate(&role.ScheduleDelegate{big.NewInt(2)})
	CreateNativeContractAccount(roleIntf)
	NativeContractInitChain(roleIntf, intf)

	return intf, nil
}

func newTransaction(contract string, method string, param []byte) *types.Transaction {
	trx := &types.Transaction{
		Sender:   contract,
		Contract: contract,
		Method:   method,
		Param:    param,
	}

	return trx
}

func NativeContractInitChain(roleIntf role.RoleInterface, ncIntf NativeContractInterface) error {
	var trxs []*types.Transaction

	// construct trxs
	var i uint32
	for i = 1; i <= config.INIT_DELEGATE_NUM; i++ {
		name := config.Genesis.InitDelegate.Name
		name = name + strconv.Itoa(int(i))

		// 1, new account trx
		nps := &NewAccountParam{
			Name:   name,
			Pubkey: config.Genesis.InitDelegate.PublicKey,
		}
		nparam, _ := json.Marshal(nps)
		trx := newTransaction(config.BOTTOS_CONTRACT_NAME, "newaccount", nparam)
		trxs = append(trxs, trx)

		// 2, transfer trx
		tps := &TransferParam{
			From:  config.BOTTOS_CONTRACT_NAME,
			To:    name,
			Value: uint64(config.Genesis.InitDelegate.Balance),
		}
		tparam, _ := json.Marshal(tps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "transfer", tparam)
		trxs = append(trxs, trx)

		// 3, set delegate
		sps := &SetDelegateParam{
			Name:   name,
			Pubkey: config.Genesis.InitDelegate.PublicKey,
		}
		sparam, _ := json.Marshal(sps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "setdelegate", sparam)
		trxs = append(trxs, trx)
	}

	// execute trxs
	for _, trx := range trxs {
		ctx := &Context{RoleIntf: roleIntf, Trx: trx}
		err := ncIntf.ExecuteNativeContract(ctx)
		if err != nil {
			fmt.Println("NativeContractInitChain Error: ", trx, err)
			//return err
			break
		}
	}

	// init CoreState delegates
	coreState, _ := roleIntf.GetCoreState()
	for i = 1; i <= config.INIT_DELEGATE_NUM; i++ {
		name := config.Genesis.InitDelegate.Name
		name = name + strconv.Itoa(int(i))

		coreState.CurrentDelegates = append(coreState.CurrentDelegates, name)
	}
	roleIntf.SetCoreState(coreState)

	return nil
}

func CreateNativeContractAccount(roleIntf role.RoleInterface) error {
	// account
	_, err := roleIntf.GetAccount(config.BOTTOS_CONTRACT_NAME)
	if err == nil {
		return nil
	}

	bto := &role.Account{
		AccountName: config.BOTTOS_CONTRACT_NAME,
		CreateTime:  uint64(time.Now().Unix()),
	}
	roleIntf.SetAccount(bto.AccountName, bto)

	// balance
	balance := &role.Balance{
		AccountName: bto.AccountName,
		Balance:     config.BOTTOS_INIT_SUPPLY,
	}
	roleIntf.SetBalance(bto.AccountName, balance)

	// staked_balance
	staked_balance := &role.StakedBalance{
		AccountName: bto.AccountName,
	}
	roleIntf.SetStakedBalance(bto.AccountName, staked_balance)

	return nil
}
