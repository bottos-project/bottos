package nativecontract

import (
	"fmt"
	"time"
	"strconv"

	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/role"
)

func CreateNativeContractAccount(ldb *db.DBService) error {
	// account
	bto := &role.Account {
		AccountName: config.BTO_CONTRACT_NAME,
		CreateTime: uint64(time.Now().Unix()),
		// Abi
	}
	role.SetAccountRole(ldb, bto.AccountName, bto)

	// balance
	balance := &role.Balance{
		AccountName: bto.AccountName,
		Balance: config.BTO_INIT_SUPPLY,
	}
	role.SetBalanceRole(ldb, bto.AccountName, balance)

	// staked_balance
	staked_balance := &role.StakedBalance{
		AccountName: bto.AccountName,
	}
	role.SetStakedBalanceRole(ldb, bto.AccountName, staked_balance)

	fmt.Println(bto, balance, staked_balance)

	return nil
}


func CreateInitialDelegates(ldb *db.DBService) error {
	initAmount := 1

	var i uint32
	for i = 1; i <= config.INIT_DELEGATE_NUM; i++ {
		name := config.Genesis.InitDelegate.Name
		name = name + strconv.Itoa(int(i))

		// 1, create account
		delegate := &role.Account {
			AccountName: name,
			CreateTime: uint64(time.Now().Unix()),
		}
		role.SetAccountRole(ldb, delegate.AccountName, delegate)

		// 2, transfer
		btoBalance, _ := role.GetBalanceRoleByAccountName(ldb, config.BTO_CONTRACT_NAME)
		btoBalance.Balance -= uint64(config.Genesis.InitDelegate.Balance)
		btoBalance.Balance -= uint64(initAmount)
		role.SetBalanceRole(ldb, btoBalance.AccountName, btoBalance)

		// balance
		balance := &role.Balance{
			AccountName: delegate.AccountName,
			Balance: uint64(config.Genesis.InitDelegate.Balance),
		}
		role.SetBalanceRole(ldb, delegate.AccountName, balance)

		// staked_balance
		staked_balance := &role.StakedBalance{
			AccountName: delegate.AccountName,
			StakedBalance: uint64(initAmount),
		}
		role.SetStakedBalanceRole(ldb, delegate.AccountName, staked_balance)

		// 3, set delegate
		delegateObj := &role.Delegate{
			AccountName: delegate.AccountName,
			SigningKey: config.Genesis.InitDelegate.PublicKey,
		}
		role.SetDelegateRole(ldb, delegateObj.AccountName, delegateObj)

		// TODO votes object

		fmt.Println(delegate, balance, staked_balance, delegateObj)
	}
	
	return nil
}

func InitNativeContract(ldb *db.DBService) error {
	CreateNativeContractAccount(ldb)
	CreateInitialDelegates(ldb)

	return nil
}
