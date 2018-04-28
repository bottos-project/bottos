package contract

import (
)


type NativeContractInterface interface {
	IsNativeContract(string, string) bool
	ExecuteNativeContract(*Context) error
}
