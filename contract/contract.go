package contract

import (
	"time"
	"strconv"
	"fmt"

	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/common/types"
)

type Context struct {
	roleIntf role.RoleInterface
	trx *types.Transaction
}


func CreateNativeContractAccount(roleIntf role.RoleInterface) error {
	// account
	bto := &role.Account {
		AccountName: config.BTO_CONTRACT_NAME,
		CreateTime: uint64(time.Now().Unix()),
	}
	roleIntf.SetAccount(bto.AccountName, bto)

	// balance
	balance := &role.Balance{
		AccountName: bto.AccountName,
		Balance: config.BTO_INIT_SUPPLY,
	}
	roleIntf.SetBalance(bto.AccountName, balance)

	// staked_balance
	staked_balance := &role.StakedBalance{
		AccountName: bto.AccountName,
	}
	roleIntf.SetStakedBalance(bto.AccountName, staked_balance)

	return nil
}


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
		btoBalance, _ := roleIntf.GetBalance(config.BTO_CONTRACT_NAME)
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

func InitNativeContract(roleIntf role.RoleInterface) error {
	CreateNativeContractAccount(roleIntf)
	CreateInitialDelegates(roleIntf)

	return nil
}
