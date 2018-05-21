package contract

import (
	"fmt"
	"math/big"
	"strconv"
	//"time"

	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/contract/msgpack"
)

func NewNativeContract(roleIntf role.RoleInterface) (NativeContractInterface, error) {
	intf, err := NewNativeContractHandler()
	if err != nil {
		return nil, err
	}
	roleIntf.SetScheduleDelegate(&role.ScheduleDelegate{big.NewInt(2)})

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

func NativeContractInitChain(roleIntf role.RoleInterface, ncIntf NativeContractInterface) ([]*types.Transaction, error) {
	err := CreateNativeContractAccount(roleIntf)
	if err != nil {
		return nil, err
	}

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
		nparam, _ := msgpack.Marshal(nps)
		trx := newTransaction(config.BOTTOS_CONTRACT_NAME, "newaccount", nparam)
		trxs = append(trxs, trx)

		// 2, transfer trx
		tps := &TransferParam{
			From:  config.BOTTOS_CONTRACT_NAME,
			To:    name,
			Value: uint64(config.Genesis.InitDelegate.Balance),
		}
		tparam, _ := msgpack.Marshal(tps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "transfer", tparam)
		trxs = append(trxs, trx)

		// 3, set delegate
		sps := &SetDelegateParam{
			Name:   name,
			Pubkey: config.Genesis.InitDelegate.PublicKey,
		}
		sparam, _ := msgpack.Marshal(sps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "setdelegate", sparam)
		trxs = append(trxs, trx)
	}

	// init CoreState delegates
	coreState, _ := roleIntf.GetCoreState()
	i = 1
	for i = 1; i <= config.BLOCKS_PER_ROUND; i++ {
		name := config.Genesis.InitDelegate.Name
		name = name + strconv.Itoa(int(i))

		coreState.CurrentDelegates = append(coreState.CurrentDelegates, name)
	}
	roleIntf.SetCoreState(coreState)

	fmt.Println("NativeContractInitChain: ", coreState)

	return trxs, nil
}

func CreateNativeContractAccount(roleIntf role.RoleInterface) error {
	// account
	_, err := roleIntf.GetAccount(config.BOTTOS_CONTRACT_NAME)
	if err == nil {
		return nil
	}

	bto := &role.Account{
		AccountName: config.BOTTOS_CONTRACT_NAME,
		CreateTime:  config.Genesis.GenesisTime,
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
