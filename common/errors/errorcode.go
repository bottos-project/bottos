// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  provide common error type
 * @Author: Gong Zibin
 * @Date:   2017-12-05
 * @Last Modified by:
 * @Last Modified time:
 */

package errors



// ErrCode define the type of error
type ErrCode uint32

const (
	// ErrNoError no error
	ErrNoError ErrCode = 0

	// ErrTrxPendingNumLimit limit of pending trx number
	ErrTrxPendingNumLimit ErrCode = 10001
	// ErrTrxSignError invalid trx sign
	ErrTrxSignError ErrCode = 10002
	// ErrTrxAccountError invalid trx account
	ErrTrxAccountError ErrCode = 10003
	// ErrTrxLifeTimeError invalid trx life time
	ErrTrxLifeTimeError ErrCode = 10004
	// ErrTrxUniqueError trx does not unique
	ErrTrxUniqueError ErrCode = 10005
	// ErrTrxChainMathError invalid trx chain math
	ErrTrxChainMathError ErrCode = 10006
	// ErrTrxContractHanldeError handle trx contract error
	ErrTrxContractHanldeError ErrCode = 10007
	// ErrTrxContractDepthError handle trx contract depth error
	ErrTrxContractDepthError ErrCode = 10008
	// ErrTrxSubContractNumError handle trx sub contract num error
	ErrTrxSubContractNumError ErrCode = 10009


	// ErrContractAccountNameIllegal invalid contract account name
	ErrContractAccountNameIllegal ErrCode = 10101
	// ErrContractAccountNotFound contract account not found
	ErrContractAccountNotFound ErrCode = 10102
	// ErrContractAccountAlreadyExist contract account already exist
	ErrContractAccountAlreadyExist ErrCode = 10103
	// ErrContractParamParseError parse contract param error
	ErrContractParamParseError ErrCode = 10104
	// ErrContractInsufficientFunds insufficient fund
	ErrContractInsufficientFunds ErrCode = 10105
	// ErrContractInvalidContractCode invalid contract code
	ErrContractInvalidContractCode ErrCode = 10106
	// ErrContractInvalidContractAbi invalid abi
	ErrContractInvalidContractAbi ErrCode = 10107
	// ErrContractUnknownContract unknown contract
	ErrContractUnknownContract ErrCode = 10108
	// ErrContractUnknownMethod unknown method
	ErrContractUnknownMethod ErrCode = 10109
	// ErrContractTransferOverflow transfer overflow
	ErrContractTransferOverflow ErrCode = 10110
	// ErrContractAccountMismatch accoumnt mismatch
	ErrContractAccountMismatch ErrCode = 10111
	// ErrContractInsufficientCredits insufficient credit
	ErrContractInsufficientCredits ErrCode = 10112

	// ErrApiTrxNotFound api trx not found
	ErrApiTrxNotFound ErrCode = 10201
	// ErrApiBlockNotFound abi block not found
	ErrApiBlockNotFound ErrCode = 10202
	// ErrApiQueryChainInfoError query chain info error
	ErrApiQueryChainInfoError ErrCode = 10203
	// ErrApiAccountNotFound api account not found
	ErrApiAccountNotFound ErrCode = 10204
	// ErrApiObjectNotFound ap object not found
	ErrApiObjectNotFound ErrCode = 10205
	// ErrTransferCreditNotFound transfer crredit not found
	ErrTransferCreditNotFound ErrCode = 10206

	// ErrActorHandleError actor handle error
	ErrActorHandleError ErrCode = 10301

	// ErrInvalid max invalid enum
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
		ErrTrxContractDepthError:  "push trx: " + "contract depth error",
		ErrTrxSubContractNumError: "push trx: " + "sub contract num error",

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

// GetCodeString get code string
func GetCodeString(errorCode ErrCode) string {
	return aaa[errorCode]
}
