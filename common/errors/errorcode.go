package errors

// import "fmt"

type ErrCode uint32

const (
	ErrNoError ErrCode = 0

	ErrTrxPendingNumLimit     ErrCode = 10001
	ErrTrxSignError           ErrCode = 10002
	ErrTrxAccountError        ErrCode = 10003
	ErrTrxLifeTimeError       ErrCode = 10004
	ErrTrxUniqueError         ErrCode = 10005
	ErrTrxChainMathError      ErrCode = 10006
	ErrTrxContractHanldeError ErrCode = 10007

	ErrContractAccountNameIllegal  ErrCode = 10101
	ErrContractAccountNotFound     ErrCode = 10102
	ErrContractAccountAlreadyExist ErrCode = 10103
	ErrContractParamParseError     ErrCode = 10104
	ErrContractInsufficientFunds   ErrCode = 10105
	ErrContractInvalidContractCode ErrCode = 10106
	ErrContractInvalidContractAbi  ErrCode = 10107
	ErrContractUnknownContract     ErrCode = 10108
	ErrContractUnknownMethod       ErrCode = 10109
	ErrContractTransferOverflow    ErrCode = 10110
	ErrContractAccountMismatch     ErrCode = 10111
	ErrContractInsufficientCredits ErrCode = 10112

	ErrApiTrxNotFound         ErrCode = 10201
	ErrApiBlockNotFound       ErrCode = 10202
	ErrApiQueryChainInfoError ErrCode = 10203
	ErrApiAccountNotFound     ErrCode = 10204
	ErrApiObjectNotFound      ErrCode = 10205
	ErrTransferCreditNotFound ErrCode = 10206

	ErrActorHandleError ErrCode = 10301

	ErrInvalid ErrCode = 0xFFFFFFFF
)

var (
	aaa = map[ErrCode]string{
		ErrTrxPendingNumLimit:     "push trx: " + "check Pending pool max num error",
		ErrTrxSignError:           "push trx: " + "check signature error",
		ErrTrxAccountError:        "push trx: " + "check account valid error",
		ErrTrxLifeTimeError:       "push trx: " + "check life time error",
		ErrTrxUniqueError:         "push trx: " + "check trx unique error",
		ErrTrxChainMathError:      "push trx: " + "check match chain error",
		ErrTrxContractHanldeError: "push trx: " + "process contract error",

		ErrContractAccountNameIllegal:  "push trx: " + "illegal account name",
		ErrContractAccountNotFound:     "push trx: " + "account name not found",
		ErrContractAccountAlreadyExist: "push trx: " + "account name already exist",
		ErrContractParamParseError:     "push trx: " + "transaction param error",
		ErrContractInsufficientFunds:   "push trx: " + "transfer account insufficient funds",
		ErrContractInvalidContractCode: "push trx: " + "invalide contract code",
		ErrContractInvalidContractAbi:  "push trx: " + "invalide contract abi",
		ErrContractUnknownContract:     "push trx: " + "unknown contract",
		ErrContractUnknownMethod:       "push trx: " + "unknown contract method",
		ErrContractTransferOverflow:    "push trx: " + "transfer overflow",
		ErrContractAccountMismatch:     "push trx: " + "sender and param account mismatch",
		ErrContractInsufficientCredits: "push trx: " + "insufficient credits",

		ErrApiTrxNotFound:         "query trx: " + "trx not found",
		ErrApiBlockNotFound:       "query block: " + "block not found",
		ErrApiQueryChainInfoError: "query chain info: " + "error",
		ErrApiAccountNotFound:     "query account: " + "account not found",
		ErrApiObjectNotFound:      "query object: " + "object not found",
		ErrTransferCreditNotFound: "query credit: " + "credit not found",

		ErrActorHandleError: "actor: " + "process error",
	}
)

func GetCodeString(errorCode ErrCode) string {
	return aaa[errorCode]
}
