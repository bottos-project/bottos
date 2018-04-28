package contract

import (
)

type NewAccountParam struct {
	Creator		string		`json:"creator"`
	Name		string		`json:"name"`
	Pubkey		string 		`json:"pubkey"`
	Deposit		uint64		`json:"deposit"`
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

type NativeContractInterface interface {
	IsNativeContract(contract string, method string) bool
	ExecuteNativeContract(*Context) error
}
