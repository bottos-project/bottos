package contract

import (
	"fmt"
	"math/big"
	//"strconv"
	//"time"

	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/contract/msgpack"
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
	var i int
	for i = 0; i < len(config.Genesis.InitDelegates); i++ {
		name := config.Genesis.InitDelegates[i].Name

		// 1, new account trx
		nps := &NewAccountParam{
			Name:   name,
			Pubkey: config.Genesis.InitDelegates[i].PublicKey,
		}
		nparam, _ := msgpack.Marshal(nps)
		trx := newTransaction(config.BOTTOS_CONTRACT_NAME, "newaccount", nparam)
		trxs = append(trxs, trx)

		// 2, transfer trx
		tps := &TransferParam{
			From:  config.BOTTOS_CONTRACT_NAME,
			To:    name,
			Value: uint64(config.Genesis.InitDelegates[i].Balance),
		}
		tparam, _ := msgpack.Marshal(tps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "transfer", tparam)
		trxs = append(trxs, trx)

		// 3, set delegate
		sps := &SetDelegateParam{
			Name:   name,
			Pubkey: config.Genesis.InitDelegates[i].PublicKey,
		}
		sparam, _ := msgpack.Marshal(sps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "setdelegate", sparam)
		trxs = append(trxs, trx)
	}

	// init CoreState delegates
	coreState, _ := roleIntf.GetCoreState()
	for i = 0; i < int(config.BLOCKS_PER_ROUND); i++ {
		name := config.Genesis.InitDelegates[i].Name

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
	var initSupply uint64
	initSupply, err = safemath.Uint64Mul(config.BOTTOS_INIT_SUPPLY, config.BOTTOS_SUPPLY_MUL)
	if err != nil {
		return err
	}

	balance := &role.Balance{
		AccountName: bto.AccountName,
		Balance:     initSupply,
	}
	roleIntf.SetBalance(bto.AccountName, balance)

	// staked_balance
	staked_balance := &role.StakedBalance{
		AccountName: bto.AccountName,
	}
	roleIntf.SetStakedBalance(bto.AccountName, staked_balance)

	return nil
}
