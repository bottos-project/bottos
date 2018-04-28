package contract

import (
	"encoding/json"
	"fmt"

	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/config"
	//"github.com/bottos-project/core/common/types"
)

type NativeContractMethod func(*Context) error

type NativeContract struct {
	Handler map[string]NativeContractMethod
}

func NewNativeContractHandler() (NativeContractInterface, error) {
	nc := &NativeContract {
		Handler: make(map[string]NativeContractMethod),
	}

	nc.Handler["newaccount"] = newaccount
	nc.Handler["transfer"] = transfer
	nc.Handler["setcode"] = setcode
	nc.Handler["setdelegate"] = setdelegate

	return nc, nil
}

func (nc *NativeContract) IsNativeContract(contract string, method string) bool {
	if contract == config.BOTTOS_CONTRACT_NAME {
		if _, ok := nc.Handler[method]; ok {
			return true
		}
	}
	return false
}

func (nc *NativeContract) ExecuteNativeContract(ctx *Context) error {
	contract := ctx.Trx.Contract.Name
	method := ctx.Trx.Method.Name
	if !nc.IsNativeContract(contract, method) {
		return fmt.Errorf("No Native Contract Method")
	}

	if handler, ok := nc.Handler[method]; ok {
		err := handler(ctx)
		return err
	}

	// TODO
	return fmt.Errorf("No Native Contract Method")
}


func check_account(RoleIntf role.RoleInterface, name string) bool {
	ac, _ := RoleIntf.GetAccount(name)
	if ac == nil {
		return false
	}

	balance, _ := RoleIntf.GetBalance(name)
	if balance == nil {
		return false
	}

	return true
}

func newaccount(ctx *Context) error {
	// trx.param --> json
	newaccount := &NewAccountParam{}
	err := json.Unmarshal(ctx.Trx.Param, newaccount)
	if err != nil {
		return err
	}
	fmt.Println("new account param: ", newaccount)

	// TODO: check from auth

	// check creator
	if !check_account(ctx.RoleIntf, newaccount.Creator) {
		return fmt.Errorf("Creator Account Not Exist")
	}

	//check name
	if check_account(ctx.RoleIntf, newaccount.Name) {
		return fmt.Errorf("Account Exist")
	}

	chainState, _ := ctx.RoleIntf.GetChainState()

	// 1, create account
	account := &role.Account {
		AccountName: newaccount.Name,
		PublicKey: []byte(newaccount.Pubkey),
		CreateTime: chainState.LastBlockTime,
	}
	ctx.RoleIntf.SetAccount(account.AccountName, account)

	// 2, transfer
	creatorBalance, _ := ctx.RoleIntf.GetBalance(newaccount.Creator)
	creatorBalance.Balance -= uint64(newaccount.Deposit)
	ctx.RoleIntf.SetBalance(newaccount.Creator, creatorBalance)

	// balance
	balance := &role.Balance{
		AccountName: newaccount.Name,
		Balance: 0,
	}
	ctx.RoleIntf.SetBalance(newaccount.Name, balance)

	// staked_balance
	staked_balance := &role.StakedBalance{
		AccountName: newaccount.Name,
		StakedBalance: uint64(newaccount.Deposit),
	}
	ctx.RoleIntf.SetStakedBalance(newaccount.Name, staked_balance)

	fmt.Println(account, balance, staked_balance)

	return nil
}

func transfer(ctx *Context) error {
	// trx.param --> json
	transfer := &TransferParam{}
	err := json.Unmarshal(ctx.Trx.Param, transfer)
	if err != nil {
		return err
	}

	fmt.Println("transfer param: ", transfer)

	// check account name
	if !check_account(ctx.RoleIntf, transfer.From) || !check_account(ctx.RoleIntf, transfer.To) {
		return fmt.Errorf("Account Not Exist")
	}

	// TODO: check from auth

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(transfer.From)
	if from.Balance < transfer.Value {
		return fmt.Errorf("Insufficient Funds")
	}
	to, _ := ctx.RoleIntf.GetBalance(transfer.To)
	
	from.Balance -= transfer.Value
	to.Balance += transfer.Value

	err = ctx.RoleIntf.SetBalance(from.AccountName, from)
	if err != nil {
		return fmt.Errorf("Transfer Error")
	}
	err = ctx.RoleIntf.SetBalance(to.AccountName, to)
	if err != nil {
		return fmt.Errorf("Transfer Error")
	}

	fmt.Println(from, to)

	return nil
}

func setcode(ctx *Context) error {
	return nil
}

func setdelegate(ctx *Context) error {
	// trx.param --> json
	param := &SetDelegateParam{}
	err := json.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return err
	}

	fmt.Println("setdelegate param: ", param)

	// TODO: check from auth

	// check account name
	if !check_account(ctx.RoleIntf, param.Name) {
		return fmt.Errorf("Account Not Exist")
	}

	_, err = ctx.RoleIntf.GetDelegateByAccountName(param.Name)

	if err != nil {
		// new delegate
		newdelegate := &role.Delegate{
			AccountName: param.Name,
			SigningKey: param.Pubkey,
		}
		ctx.RoleIntf.SetDelegate(newdelegate.AccountName, newdelegate)
		fmt.Println(newdelegate)
		// TODO votes object
	} else {
		return fmt.Errorf("Delegate Already Exist")
	}

	return nil
}
