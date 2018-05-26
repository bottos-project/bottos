package contract

import (
	bottosErr "github.com/bottos-project/bottos/common/errors"
)

type ContractError uint32

const (
	ERROR_NONE 								ContractError 				= 0
	
	ERROR_CONT_HANDLE_FAIL					ContractError				= 1
	ERROR_CONT_ACCOUNT_NAME_NULL			ContractError				= 100
	ERROR_CONT_ACCOUNT_NAME_TOO_LONG		ContractError				= 101
	ERROR_CONT_ACCOUNT_NAME_ILLEGAL			ContractError				= 102
	ERROR_CONT_ACCOUNT_NOT_EXIST			ContractError				= 103
	ERROR_CONT_ACCOUNT_ALREADY_EXIST		ContractError				= 104
	ERROR_CONT_PARAM_TOO_LONG				ContractError				= 105
	ERROR_CONT_PARAM_PARSE_ERROR			ContractError				= 106
	ERROR_CONT_INSUFFICIENT_FUNDS			ContractError				= 107
	ERROR_CONT_CODE_INVALID					ContractError				= 108
	ERROR_CONT_ABI_PARSE_FAIL				ContractError				= 109
	ERROR_CONT_UNKNOWN_CONTARCT				ContractError				= 110
	ERROR_CONT_UNKNOWN_METHOD				ContractError				= 111
	
)

func ConvertErrorCode(cerr ContractError) bottosErr.ErrCode {
	switch cerr {
	case ERROR_CONT_HANDLE_FAIL:
		return bottosErr.ErrTrxContractHanldeError 
    case ERROR_CONT_ACCOUNT_NAME_NULL:
        return bottosErr.ErrContractAccountNameIllegal
    case ERROR_CONT_ACCOUNT_NAME_TOO_LONG:
        return bottosErr.ErrContractAccountNameIllegal
    case ERROR_CONT_ACCOUNT_NAME_ILLEGAL:
		return bottosErr.ErrContractAccountNameIllegal
	case ERROR_CONT_ACCOUNT_NOT_EXIST:
		return bottosErr.ErrContractAccountNotFound
	case ERROR_CONT_ACCOUNT_ALREADY_EXIST:
		return bottosErr.ErrContractAccountAlreadyExist
	case ERROR_CONT_PARAM_TOO_LONG:
		return bottosErr.ErrContractParamParseError
	case ERROR_CONT_PARAM_PARSE_ERROR:
		return bottosErr.ErrContractParamParseError
	case ERROR_CONT_INSUFFICIENT_FUNDS:
		return bottosErr.ErrContractInsufficientFunds
	case ERROR_CONT_CODE_INVALID:
		return bottosErr.ErrContractInvalidContractCode
	case ERROR_CONT_ABI_PARSE_FAIL:
		return bottosErr.ErrContractInvalidContractAbi
	case ERROR_CONT_UNKNOWN_CONTARCT:
		return bottosErr.ErrContractUnknownContract
	case ERROR_CONT_UNKNOWN_METHOD:
        return bottosErr.ErrContractUnknownMethod
    }
	return bottosErr.ErrNoError
}

