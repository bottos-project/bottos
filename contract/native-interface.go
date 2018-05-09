package contract

import (
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
	ExecuteNativeContract(*Context) error
}
