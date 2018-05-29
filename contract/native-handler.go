package contract

import (
	"fmt"
	"regexp"

	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/contract/msgpack"
)


func newAccount(ctx *Context) ContractError {
	newaccount := &NewAccountParam{}
	err := msgpack.Unmarshal(ctx.Trx.Param, newaccount)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}
	//fmt.Println("new account param: ", newaccount)

	// TODO: check auth

	//check account
	cerr := checkAccountName(newaccount.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	if isAccountNameExist(ctx.RoleIntf, newaccount.Name) {
		return ERROR_CONT_ACCOUNT_ALREADY_EXIST
	}

	// TODO: check pubkey

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

	//fmt.Println(account, balance, staked_balance)

	return ERROR_NONE
}

func transfer(ctx *Context) ContractError {
	transfer := &TransferParam{}
	err := msgpack.Unmarshal(ctx.Trx.Param, transfer)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	//fmt.Println("transfer param: ", transfer)

	// TODO: check auth

	// check account
	cerr := checkAccount(ctx.RoleIntf, transfer.From)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, transfer.To)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check Sender

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(transfer.From)
	if from.Balance < transfer.Value {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}
	to, _ := ctx.RoleIntf.GetBalance(transfer.To)

	err = from.SafeSub(transfer.Value)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = to.SafeAdd(transfer.Value)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}

	err = ctx.RoleIntf.SetBalance(from.AccountName, from)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}
	err = ctx.RoleIntf.SetBalance(to.AccountName, to)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	//fmt.Println(from, to)

	return ERROR_NONE
}

func setDelegate(ctx *Context) ContractError {
	param := &SetDelegateParam{}
	err := msgpack.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	//fmt.Println("setDelegate param: ", param)

	// TODO: check auth

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	// TODO check pubkey

	_, err = ctx.RoleIntf.GetDelegateByAccountName(param.Name)
	if err != nil {
		// new delegate
		newdelegate := &role.Delegate{
			AccountName: param.Name,
			ReportKey:   param.Pubkey,
		}
		ctx.RoleIntf.SetDelegate(newdelegate.AccountName, newdelegate)
		//fmt.Println(newdelegate)

		//create schedule delegate vote role
		scheduleDelegate, err := ctx.RoleIntf.GetScheduleDelegate()
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}
		//create delegate vote role
		ctx.RoleIntf.CreateDelegateVotes()

		newDelegateVotes := new(role.DelegateVotes).StartNewTerm(scheduleDelegate.CurrentTermTime)
		newDelegateVotes.OwnerAccount = newdelegate.AccountName
		err = ctx.RoleIntf.SetDelegateVotes(newdelegate.AccountName, newDelegateVotes)
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}
		fmt.Println("set delegate vote", newDelegateVotes)
	} else {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func grantCredit(ctx *Context) ContractError {
	param := &GrantCreditParam{}
	err := msgpack.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	//fmt.Println("grantcredit param: ", param)

	// TODO: check auth

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, param.Spender)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ctx.Trx.Sender)
	if cerr != ERROR_NONE {
		return cerr
	}

	// sender must be from
	if ctx.Trx.Sender != param.Name {
		return ERROR_CONT_ACCOUNT_MISMATCH
	}

	// check limit
	balance, err := ctx.RoleIntf.GetBalance(param.Name)
	if balance.Balance < param.Limit {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}

	credit := &role.TransferCredit{
		Name: param.Name,
		Spender: param.Spender,
		Limit: param.Limit,
	}
	err = ctx.RoleIntf.SetTransferCredit(credit.Name, credit)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	//fmt.Println(credit)

	return ERROR_NONE
}

func cancelCredit(ctx *Context) ContractError {
	param := &CancelCreditParam{}
	err := msgpack.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	//fmt.Println("cancelcredit param: ", param)

	// TODO: check auth

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, param.Spender)
	if cerr != ERROR_NONE {
		return cerr
	}

	_, err = ctx.RoleIntf.GetTransferCredit(param.Name, param.Spender)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	err = ctx.RoleIntf.DeleteTransferCredit(param.Name, param.Spender)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func transferFrom(ctx *Context) ContractError {
	transfer := &TransferFromParam{}
	err := msgpack.Unmarshal(ctx.Trx.Param, transfer)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	//fmt.Println("transferfrom param: ", transfer)

	// TODO: check auth

	// check account
	cerr := checkAccount(ctx.RoleIntf, transfer.From)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, transfer.To)
	if cerr != ERROR_NONE {
		return cerr
	}

	cerr = checkAccount(ctx.RoleIntf, ctx.Trx.Sender)
	if cerr != ERROR_NONE {
		return cerr
	}

	// Note: sender is the spender
	// check limit
	credit, err := ctx.RoleIntf.GetTransferCredit(transfer.From, ctx.Trx.Sender)
	if err != nil {
		return ERROR_CONT_INSUFFICIENT_CREDITS
	}
	if transfer.Value > credit.Limit {
		return ERROR_CONT_INSUFFICIENT_CREDITS
	}

	// check funds
	from, _ := ctx.RoleIntf.GetBalance(transfer.From)
	if from.Balance < transfer.Value {
		return ERROR_CONT_INSUFFICIENT_FUNDS
	}
	to, _ := ctx.RoleIntf.GetBalance(transfer.To)

	err = from.SafeSub(transfer.Value)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = credit.SafeSub(transfer.Value)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}
	err = to.SafeAdd(transfer.Value)
	if err != nil {
		return ERROR_CONT_TRANSFER_OVERFLOW
	}

	err = ctx.RoleIntf.SetBalance(from.AccountName, from)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}
	err = ctx.RoleIntf.SetBalance(to.AccountName, to)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	if credit.Limit > 0 {
		err = ctx.RoleIntf.SetTransferCredit(credit.Name, credit)
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}
	} else {
		err = ctx.RoleIntf.DeleteTransferCredit(credit.Name, ctx.Trx.Sender)
		if err != nil {
			return ERROR_CONT_HANDLE_FAIL
		}
	}

	//fmt.Println(from, to, credit)

	return ERROR_NONE
}

func checkCode(code []byte) error {
	// TODO 
	return nil
}

func deployCode(ctx *Context) ContractError {
	param := &DeployCodeParam{}
	err := msgpack.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	//fmt.Println("deployCode param: ", param)

	// TODO: check auth

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check code
	err = checkCode(param.ContractCode)
	if err != nil {
		return ERROR_CONT_CODE_INVALID
	}

	codeHash := common.Sha256(param.ContractCode)

	account, _ := ctx.RoleIntf.GetAccount(param.Name)
	account.CodeVersion = codeHash
	account.ContractCode = make([]byte, len(param.ContractCode))
	copy(account.ContractCode, param.ContractCode)
	err = ctx.RoleIntf.SetAccount(account.AccountName, account)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func checkAbi(abiRaw []byte) error {
	_, err := ParseAbi(abiRaw)
	if err != nil {
		return fmt.Errorf("ABI Parse error: %v", err) 
	}
	return nil
}

func deployAbi(ctx *Context) ContractError {
	param := &DeployABIParam{}
	err := msgpack.Unmarshal(ctx.Trx.Param, param)
	if err != nil {
		return ERROR_CONT_PARAM_PARSE_ERROR
	}

	//fmt.Println("deployAbi param: ", param)

	// TODO: check auth

	// check account
	cerr := checkAccount(ctx.RoleIntf, param.Name)
	if cerr != ERROR_NONE {
		return cerr
	}

	// check abi
	err = checkAbi(param.ContractAbi)
	if err != nil {
		return ERROR_CONT_ABI_PARSE_FAIL
	}

	account, _ := ctx.RoleIntf.GetAccount(param.Name)
	account.ContractAbi = make([]byte, len(param.ContractAbi))
	copy(account.ContractAbi, param.ContractAbi)
	err = ctx.RoleIntf.SetAccount(account.AccountName, account)
	if err != nil {
		return ERROR_CONT_HANDLE_FAIL
	}

	return ERROR_NONE
}

func checkAccountName(name string) ContractError {
	if len(name) == 0 {
		return ERROR_CONT_ACCOUNT_NAME_NULL
	}

	if len(name) > config.MAX_ACCOUNT_NAME_LENGTH {
		return ERROR_CONT_ACCOUNT_NAME_TOO_LONG
	}

	if !checkAccountNameContent(name) {
		return ERROR_CONT_ACCOUNT_NAME_ILLEGAL
	}

	return ERROR_NONE
}

func checkAccountNameContent(name string) bool {
	match, err := regexp.MatchString(config.ACCOUNT_NAME_REGEXP, name)
	if err != nil {
		return false
	}
	if !match {
		return false
	}

	return true
}

func isAccountNameExist(RoleIntf role.RoleInterface, name string) bool {
	account, err := RoleIntf.GetAccount(name)
	if err == nil {
		if account != nil && account.AccountName == name {
			return true
		}
	}
	return false
}

func checkAccount(RoleIntf role.RoleInterface, name string) ContractError {
	cerr := checkAccountName(name)
	if cerr != ERROR_NONE {
		return cerr
	}

	if !isAccountNameExist(RoleIntf, name) {
		return ERROR_CONT_ACCOUNT_NOT_EXIST
	}

	return ERROR_NONE
}
