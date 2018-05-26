package contract

import (
	"github.com/bottos-project/bottos/config"
)

type NewAccountParam struct {
	Name		string		`json:"name"`
	Pubkey		string 		`json:"pubkey"`
}

type TransferParam struct {
	From		string		`json:"from"`
	To			string		`json:"to"`
	Value		uint64		`json:"value"`
}

type SetDelegateParam struct {
	Name		string		`json:"name"`
	Pubkey		string 		`json:"pubkey"`
	// TODO CONFIG
}

type DeployCodeParam struct {
	Name		 string		 `json:"name"`
	VMType       byte        `json:"vm_type"`
	VMVersion    byte        `json:"vm_version"`
	ContractCode []byte      `json:"contract_code"`
}

type DeployABIParam struct {
	Name		 string		 `json:"name"`
	ContractAbi	 []byte      `json:"contract_abi"`
}

type NativeContractInterface interface {
	IsNativeContract(contract string, method string) bool
	ExecuteNativeContract(*Context) ContractError
}

type NativeContractMethod func(*Context) ContractError

type NativeContract struct {
	Handler map[string]NativeContractMethod
}

func NewNativeContractHandler() (NativeContractInterface, error) {
	nc := &NativeContract{
		Handler: make(map[string]NativeContractMethod),
	}

	nc.Handler["newaccount"] = newAccount
	nc.Handler["transfer"] = transfer
	nc.Handler["setdelegate"] = setDelegate
	nc.Handler["deploycode"] = deployCode
	nc.Handler["deployabi"] = deployAbi

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

func (nc *NativeContract) ExecuteNativeContract(ctx *Context) ContractError {
	contract := ctx.Trx.Contract
	method := ctx.Trx.Method
	if nc.IsNativeContract(contract, method) {
		if handler, ok := nc.Handler[method]; ok {
			contErr := handler(ctx)
			return contErr
		} else {
			return ERROR_CONT_UNKNOWN_METHOD
		}
	} else {
		return ERROR_CONT_UNKNOWN_CONTARCT
	}
}
