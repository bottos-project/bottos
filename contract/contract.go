package contract

import (
	"time"
	"strconv"
	"fmt"
	"encoding/json"

	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/role"
)


func NewNativeContract(roleIntf role.RoleInterface) (NativeContractInterface, error) {
	intf, err := NewNativeContractHandler()
	if err != nil {
		return nil, err
	}

	CreateNativeContractAccount(roleIntf)
	NativeContractInitChain(roleIntf, intf)

	return intf, nil
}

func newTransaction(contract string, method string, param []byte) *types.Transaction {
	trx := &types.Transaction {
		Sender: &types.AccountName{Name:contract},
		Contract: &types.ContractName{Name:contract},
		Method: &types.MethodName{Name:method},
		Param: param,
	}

	return trx
}

func NativeContractInitChain(roleIntf role.RoleInterface, ncIntf NativeContractInterface) error {
	var trxs []*types.Transaction
	initAmount := uint64(1)

	// construct trxs
	var i uint32
	for i = 1; i <= config.INIT_DELEGATE_NUM; i++ {
		name := config.Genesis.InitDelegate.Name
		name = name + strconv.Itoa(int(i))

		// 1, new account trx
		nps := &NewAccountParam{
			Creator: config.BOTTOS_CONTRACT_NAME, 
			Name: name, 
			Pubkey: config.Genesis.InitDelegate.PublicKey, 
			Deposit: initAmount,
		}
		nparam, _ := json.Marshal(nps)
		trx := newTransaction(config.BOTTOS_CONTRACT_NAME, "newaccount", nparam)
		trxs = append(trxs, trx)

		// 2, transfer trx
		tps := &TransferParam{
			From: config.BOTTOS_CONTRACT_NAME, 
			To: name, 
			Value: uint64(config.Genesis.InitDelegate.Balance),
		}
		tparam, _ := json.Marshal(tps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "transfer", tparam)
		trxs = append(trxs, trx)

		// 3, set delegate
		sps := &SetDelegateParam{
			Name: name, 
			Pubkey: config.Genesis.InitDelegate.PublicKey, 
		}
		sparam, _ := json.Marshal(sps)
		trx = newTransaction(config.BOTTOS_CONTRACT_NAME, "setdelegate", sparam)
		trxs = append(trxs, trx)
	}

	// execute trxs
	for _, trx := range trxs {
		ctx := &Context{roleIntf: roleIntf, Trx: trx}
		err := ncIntf.ExecuteNativeContract(ctx)
		if err != nil {
			fmt.Printf("NativeContractInitChain Error: ", trx, err)
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
	bto := &role.Account {
		AccountName: config.BOTTOS_CONTRACT_NAME,
		CreateTime: uint64(time.Now().Unix()),
	}
	roleIntf.SetAccount(bto.AccountName, bto)

	// balance
	balance := &role.Balance{
		AccountName: bto.AccountName,
		Balance: config.BOTTOS_INIT_SUPPLY,
	}
	roleIntf.SetBalance(bto.AccountName, balance)

	// staked_balance
	staked_balance := &role.StakedBalance{
		AccountName: bto.AccountName,
	}
	roleIntf.SetStakedBalance(bto.AccountName, staked_balance)

	return nil
}

/*
func CreateInitialDelegates(roleIntf role.RoleInterface) error {
	initAmount := 1

	coreState, _ := roleIntf.GetCoreState()

	var i uint32
	for i = 1; i <= config.INIT_DELEGATE_NUM; i++ {
		name := config.Genesis.InitDelegate.Name
		name = name + strconv.Itoa(int(i))

		// 1, create account
		delegate := &role.Account {
			AccountName: name,
			CreateTime: uint64(time.Now().Unix()),
		}
		roleIntf.SetAccount(delegate.AccountName, delegate)

		// 2, transfer
		btoBalance, _ := roleIntf.GetBalance(config.BOTTOS_CONTRACT_NAME)
		btoBalance.Balance -= uint64(config.Genesis.InitDelegate.Balance)
		btoBalance.Balance -= uint64(initAmount)
		roleIntf.SetBalance(btoBalance.AccountName, btoBalance)

		// balance
		balance := &role.Balance{
			AccountName: delegate.AccountName,
			Balance: uint64(config.Genesis.InitDelegate.Balance),
		}
		roleIntf.SetBalance(delegate.AccountName, balance)

		// staked_balance
		staked_balance := &role.StakedBalance{
			AccountName: delegate.AccountName,
			StakedBalance: uint64(initAmount),
		}
		roleIntf.SetStakedBalance(delegate.AccountName, staked_balance)

		// 3, set delegate
		delegateObj := &role.Delegate{
			AccountName: delegate.AccountName,
			SigningKey: config.Genesis.InitDelegate.PublicKey,
		}
		roleIntf.SetDelegate(delegateObj.AccountName, delegateObj)

		// TODO votes object

		coreState.CurrentDelegates = append(coreState.CurrentDelegates, delegate.AccountName)
	}

	roleIntf.SetCoreState(coreState)

	fmt.Println(coreState)
	
	return nil
}
*/
