package contract

import (
	"encoding/json"
	"fmt"

	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/role"
	"github.com/bottos-project/core/common"
)

type NativeContractMethod func(*Context) error

type NativeContract struct {
	Handler map[string]NativeContractMethod
}

func NewNativeContractHandler() (NativeContractInterface, error) {
	nc := &NativeContract{
		Handler: make(map[string]NativeContractMethod),
	}

	nc.Handler["newaccount"] = newaccount
	nc.Handler["transfer"] = transfer
	nc.Handler["setdelegate"] = setdelegate
	nc.Handler["deploycode"] = deploycode
	nc.Handler["deployabi"] = deployabi

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
	contract := ctx.Trx.Contract
	method := ctx.Trx.Method
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

func check_account(RoleIntf role.RoleInterface, name string) error {
	if len(name) == 0 || len(name) > config.MAX_ACCOUNT_NAME_LENGTH {
		return fmt.Errorf("Invalid Account Name")
	}

	_, err := RoleIntf.GetAccount(name)
	if err != nil {
		return fmt.Errorf("Account Not Exist")
	}

	return nil
}

func newaccount(ctx *Context) error {
	// trx.param --> json
	newaccount := &NewAccountParam{}
	err := json.Unmarshal(ctx.Trx.Param, newaccount)
	if err != nil {
		return err
	}
	fmt.Println("new account param: ", newaccount)

	// TODO: check auth

	//check account
	err = check_account(ctx.RoleIntf, newaccount.Name)
	if err != nil {
		return err
	}

	chainState, _ := ctx.RoleIntf.GetChainState()

	// 1, create account
	account := &role.Account{
		AccountName: newaccount.Name,
		PublicKey:   []byte(newaccount.Pubkey),
		CreateTime:  chainState.LastBlockTime,
	}
	ctx.RoleIntf.SetAccount(account.AccountName, account)

	// 2, create balance
	balance := &role.Balance{
		AccountName: newaccount.Name,
		Balance:     0,
	}
	ctx.RoleIntf.SetBalance(newaccount.Name, balance)

	// 3, create staked_balance
	staked_balance := &role.StakedBalance{
		AccountName:   newaccount.Name,
		StakedBalance: 0,
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

	// TODO: check auth

	// check account
	err = check_account(ctx.RoleIntf, transfer.From)
	if err != nil {
		return err
	}

	err = check_account(ctx.RoleIntf, transfer.To)
	if err != nil {
		return err
	}

	// check funds
	// TODO safe math check
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

func setdelegate(ctx *Context) error {
	// trx.param --> json
	param := &SetDelegateParam{}
	err := json.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return err
	}

	fmt.Println("setdelegate param: ", param)

	// TODO: check auth

	// check account
	err = check_account(ctx.RoleIntf, param.Name)
	if err != nil {
		return err
	}

	_, err = ctx.RoleIntf.GetDelegateByAccountName(param.Name)

	if err != nil {
		// new delegate
		newdelegate := &role.Delegate{
			AccountName: param.Name,
			ReportKey:   param.Pubkey,
		}
		ctx.RoleIntf.SetDelegate(newdelegate.AccountName, newdelegate)
		fmt.Println(newdelegate)

		//create schedule delegate vote role
		scheduleDelegate, err := ctx.RoleIntf.GetScheduleDelegate()
		if err != nil {
			return fmt.Errorf("critical error schedule delegate is not exist")
		}
		//create delegate vote role
		ctx.RoleIntf.CreateDelegateVotes()

		newDelegateVotes := new(role.DelegateVotes).StartNewTerm(scheduleDelegate.CurrentTermTime)
		newDelegateVotes.OwnerAccount = newdelegate.AccountName
		err = ctx.RoleIntf.SetDelegateVotes(newdelegate.AccountName, newDelegateVotes)
		if err != nil {
			return fmt.Errorf("set Delegate vote failed")
		}
		fmt.Println("set delegate vote", newDelegateVotes)
	} else {
		return fmt.Errorf("Delegate Already Exist")
	}

	return nil
}

func check_code(code []byte) error {

	return nil
}

func deploycode(ctx *Context) error {
	// trx.param --> json
	param := &DeployCodeParam{}
	err := json.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return err
	}

	fmt.Println("deploycode param: ", param)

	// TODO: check auth

	// check account
	err = check_account(ctx.RoleIntf, param.Name)
	if err != nil {
		return err
	}

	// check code
	err = check_code(param.ContractCode)
	if err != nil {
		return err
	}

	codeHash := common.Sha256(param.ContractCode)

	account, _ := ctx.RoleIntf.GetAccount(param.Name)
	account.CodeVersion = codeHash
	account.ContractCode = make([]byte, len(param.ContractCode))
	copy(account.ContractCode, param.ContractCode)
	err = ctx.RoleIntf.SetAccount(account.AccountName, account)
	if err != nil {
		return fmt.Errorf("Set Code Fail")
	}

	return nil
}

func check_abi(abi []byte) error {

	return nil
}

func deployabi(ctx *Context) error {
	// trx.param --> json
	param := &DeployABIParam{}
	err := json.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return err
	}

	fmt.Println("deployabi param: ", param)

	// TODO: check auth

	// check account
	err = check_account(ctx.RoleIntf, param.Name)
	if err != nil {
		return err
	}

	// check code
	err = check_code(param.ContractAbi)
	if err != nil {
		return err
	}

	account, _ := ctx.RoleIntf.GetAccount(param.Name)
	account.ContractAbi = make([]byte, len(param.ContractAbi))
	copy(account.ContractAbi, param.ContractAbi)
	err = ctx.RoleIntf.SetAccount(account.AccountName, account)
	if err != nil {
		return fmt.Errorf("Set Abi Fail")
	}

	return nil
}