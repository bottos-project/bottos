package contract

import (
	bottosErr "github.com/bottos-project/bottos/common/errors"
)

//ContractError is to record error info in contracts
type ContractError uint32

const (
	//ERROR_NONE is no error
	ERROR_NONE ContractError = 0

	//ERROR_CONT_HANDLE_FAIL means no contract handle failed.
	ERROR_CONT_HANDLE_FAIL ContractError = 1
	//ERROR_CONT_ACCOUNT_NAME_NULL means account name is null.
	ERROR_CONT_ACCOUNT_NAME_NULL ContractError = 100
	//ERROR_CONT_ACCOUNT_NAME_TOO_LONG means account name is too long.
	ERROR_CONT_ACCOUNT_NAME_TOO_LONG ContractError = 101
	//ERROR_CONT_ACCOUNT_NAME_ILLEGAL means account name is illegal.
	ERROR_CONT_ACCOUNT_NAME_ILLEGAL ContractError = 102
	//ERROR_CONT_ACCOUNT_NOT_EXIST means account is not exist.
	ERROR_CONT_ACCOUNT_NOT_EXIST ContractError = 103
	//ERROR_CONT_ACCOUNT_ALREADY_EXIST means account is already exist.
	ERROR_CONT_ACCOUNT_ALREADY_EXIST ContractError = 104
	//ERROR_CONT_PARAM_TOO_LONG means param is too long.
	ERROR_CONT_PARAM_TOO_LONG ContractError = 105
	//ERROR_CONT_PARAM_PARSE_ERROR means param parse error.
	ERROR_CONT_PARAM_PARSE_ERROR ContractError = 106
	//ERROR_CONT_INSUFFICIENT_FUNDS means insufficient funds.
	ERROR_CONT_INSUFFICIENT_FUNDS ContractError = 107
	//ERROR_CONT_CODE_INVALID means code invalid.
	ERROR_CONT_CODE_INVALID ContractError = 108
	//ERROR_CONT_ABI_PARSE_FAIL means abi parse failed.
	ERROR_CONT_ABI_PARSE_FAIL ContractError = 109
	//ERROR_CONT_UNKNOWN_CONTARCT means unknown contract, please check.
	ERROR_CONT_UNKNOWN_CONTARCT ContractError = 110
	//ERROR_CONT_UNKNOWN_METHOD means unknown method.
	ERROR_CONT_UNKNOWN_METHOD ContractError = 111
	//ERROR_CONT_TRANSFER_OVERFLOW means transfer overflow.
	ERROR_CONT_TRANSFER_OVERFLOW ContractError = 112
	//ERROR_CONT_ACCOUNT_MISMATCH means account mismatch.
	ERROR_CONT_ACCOUNT_MISMATCH ContractError = 113
	//ERROR_CONT_INSUFFICIENT_CREDITS means insufficient credits.
	ERROR_CONT_INSUFFICIENT_CREDITS ContractError = 114
)

//ConvertErrorCode is corresponding to error info
func ConvertErrorCode(cerr ContractError) bottosErr.ErrCode {
	switch cerr {
	case ERROR_NONE:
		return bottosErr.ErrNoError
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
	case ERROR_CONT_TRANSFER_OVERFLOW:
		return bottosErr.ErrContractTransferOverflow
	case ERROR_CONT_ACCOUNT_MISMATCH:
		return bottosErr.ErrContractAccountMismatch
	case ERROR_CONT_INSUFFICIENT_CREDITS:
		return bottosErr.ErrContractInsufficientCredits
	}
	return bottosErr.ErrTrxContractHanldeError
}
