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

import "strconv"

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

	// ErrTrxBlockSyncingError handle trx when block syncing error
	ErrTrxBlockSyncingError ErrCode = 10010

	// ErrTrxAlreadyInPoolError handle trx when trx already in pool
	ErrTrxAlreadyInPool ErrCode = 10011

	// ErrTrxContractError invalid trx contract
	ErrTrxContractError ErrCode = 10012
	// ErrTrxMethodError invalid trx method
	ErrTrxMethodError ErrCode = 10013

	// ErrTrxCheckSpaceError trx has not enough Space token
	ErrTrxCheckSpaceError ErrCode = 10014
	// ErrTrxCheckTimeError trx has not enough Time token
	ErrTrxCheckTimeError ErrCode = 10015
	// ErrTrxCheckMinSpaceError trx has not enough Space token
	ErrTrxCheckMinSpaceError ErrCode = 10016
	// ErrTrxCheckMinTimeError trx has not enough Time token
	ErrTrxCheckMinTimeError ErrCode = 10017

	// ErrTrxCheckSpaceInternalError trx check Space token error
	ErrTrxCheckSpaceInternalError ErrCode = 10018
	// ErrTrxCheckTimeInternalError trx check Time token error
	ErrTrxCheckTimeInternalError ErrCode = 10019
	// ErrTrxCheckResourceInternalError trx check Resource limit error
	ErrTrxCheckResourceInternalError ErrCode = 10020
	// ErrTrxResourceExceedMaxSpacePerTrx trx size Exceeding the maximum value
	ErrTrxResourceExceedMaxSpacePerTrx ErrCode = 10021
	// ErrTrxResourceCheckMinBalance min balance has not enough
	ErrTrxResourceCheckMinBalance ErrCode = 10022

	ErrTrxExecTimeOver           ErrCode = 10023
	ErrTrxVmTypeInvalid          ErrCode = 10024
	ErrTrxNoticeContractNumError ErrCode = 10025
	ErrTrxContractNotExist       ErrCode = 10026
	ErrTrxCacheNumLimit          ErrCode = 10027
	ErrTrxAlreadyInCache         ErrCode = 10028

	// ErrTrxVersionError invalid trx version
	ErrTrxVersionError ErrCode = 10051

	// ErrAccountNameIllegal invalid contract account name
	ErrAccountNameIllegal ErrCode = 10101
	// ErrAccountNotFound contract account not found
	ErrAccountNotFound ErrCode = 10102
	// ErrAccountAlreadyExist contract account already exist
	ErrAccountAlreadyExist ErrCode = 10103
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
	// ErrAccountMismatch accoumnt mismatch
	ErrAccountMismatch ErrCode = 10111
	// ErrContractInsufficientCredits insufficient credit
	ErrContractInsufficientCredits ErrCode = 10112
	// ErrAccountPubkeyLenIllegal invalid pubkey len
	ErrAccountPubkeyLenIllegal ErrCode = 10113
	// ErrContractGenesisPermissionError no permission
	ErrContractGenesisPermissionError ErrCode = 10114
	// ErrContractChainNotActivated chain not activated
	ErrContractChainNotActivated ErrCode = 10115
	// ErrContractTransferToSelf cannot transfer to self
	ErrContractTransferToSelf ErrCode = 10116
	// ErrContractGrantToSelf cannot grant to self
	ErrContractGrantToSelf                     ErrCode = 10117
	ErrContractNumReachMaxPerAccount           ErrCode = 10118
	ErrContractNotFound                        ErrCode = 10119
	ErrContractNameIllegal                     ErrCode = 10120
	ErrContractNoStakedVoteFunds               ErrCode = 10121
	ErrContractMustVoteToValidDelegate         ErrCode = 10122
	ErrContractDelegateVoteAlreadyImported     ErrCode = 10123
	ErrContractBlockProducingAlreadyTransfered ErrCode = 10124

	ErrContractJSNotSupport                 ErrCode = 10125
	ErrContractAlreadyExist                 ErrCode = 10126
	ErrContractInsufficientTransferValue    ErrCode = 10127
	ErrContractAlreadyClaimedReward         ErrCode = 10128
	ErrContractUnstakeReleaseTimeNotReached ErrCode = 10129
	ErrContractInsufficientUnstakeFunds     ErrCode = 10130

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
	ErrApiAccountNameIllegal  ErrCode = 10207
	ErrApiProposalNameIllegal ErrCode = 10210

	// ErrActorHandleError actor handle error
	ErrActorHandleError ErrCode = 10301

	ErrContractInsufficientStakeBalanceToRegDelegate ErrCode = 10302

	// ErrBlockInsertErrorGeneral general error
	ErrBlockInsertErrorGeneral ErrCode = 10401
	// ErrBlockInsertErrorNotLinked the block not linked to the chain
	ErrBlockInsertErrorNotLinked ErrCode = 10402
	// ErrBlockInsertErrorGeneral block validate fail
	ErrBlockInsertErrorValidateFail ErrCode = 10403
	// ErrBlockInsertErrorGeneral different lib block but linked
	ErrBlockInsertErrorDiffLibLinked ErrCode = 10404
	// ErrBlockInsertErrorGeneral different lib block and not linked in this chain
	ErrBlockInsertErrorDiffLibNotLinked ErrCode = 10405
	// ErrBlockVersionError invalid block version
	ErrBlockVersionError ErrCode = 10406

	// ErrConsensusReceiveMessageTimeout general error
	ErrConsensusReceiveMessageTimeout ErrCode = 10501

	//ErrWalletPasswdError input password error
	ErrWalletPasswdError ErrCode = 10601

	// ErrMsignProposalNameIllegal invalid msign prposal name
	ErrMsignProposalNameIllegal ErrCode = 10701
	// ErrPrposalNotFound msign prposal not found
	ErrMsignProposalNotFound ErrCode = 10702
	// ErrMsignProposalAlreadyExist msign prposal already exist
	ErrMsignProposalAlreadyExist ErrCode = 10703
	// ErrProposalNotFound msign prposal approved not found
	ErrMsignProposalApprovedNotFound ErrCode = 10704
	// ErrProposalNotFound msign prposal transfer illegal
	ErrMsignProposalTransferIllegal ErrCode = 10705
	// ErrProposalNotFound msign prposal authority weight not enough
	ErrMsignProposalWeightNotEnough ErrCode = 10706
	// ErrMsignProposalNoAuthority msign prposal no authority data
	ErrMsignProposalNoAuthorityData ErrCode = 10707
	// ErrMsignProposalNoAuthority msign prposal threshold illegal
	ErrMsignProposalThresholdIllegal ErrCode = 10708
	// ErrMsignProposalNoAuthority msign prposal authority account illegal
	ErrMsignProposalAuthorityAccountIllegal ErrCode = 10709

	// RestErrJsonNewEncoder rest Json NewEncoder
	RestErrInternal          ErrCode = 20000
	RestErrGenerateParm      ErrCode = 20001
	RestErrJsonNewEncoder    ErrCode = 20002
	RestErrBplMarshal        ErrCode = 20003
	RestErrDecodeStringError ErrCode = 20004
	RestErrStringToBig       ErrCode = 20005

	RestErrReqNil      ErrCode = 20100
	RestErrResultNil   ErrCode = 20101
	RestErrPriKeyError ErrCode = 20102
	RestErrPubKeyError ErrCode = 20103
	RestErrHashError   ErrCode = 20104

	RestErrTrxSignError   ErrCode = 20202
	RestErrUnlockError    ErrCode = 20211
	RestErrUnkownAccType  ErrCode = 20221
	RestErrWalletLocked   ErrCode = 20222
	RestErrWalletExist    ErrCode = 20223
	RestErrWalletNotExist ErrCode = 20224

	RestErrTxPending  ErrCode = 20300
	RestErrTxNotFound ErrCode = 20301
	RestErrTxPacked   ErrCode = 20302
	RestErrTxSending  ErrCode = 20303

	RestErrGetResSpaceError ErrCode = 20320
	RestErrGetResTimeError  ErrCode = 20321

	RestErrGetMsignTransferError ErrCode = 20330

	RestErrLogItemInvalid ErrCode = 20331

	// contract exec code: 0x30000:

	ContractExecStart ErrCode = 0x30000
	/* wasm js program exec error : 0x30001~ 0x30fff */

	/* wasm vm exe error : 0x31000~ 0x31fff */
	WASMExecErrorStart ErrCode = 0x31000

	WASMEXecError_VM_ERR_CREATE_VM      ErrCode = 0x31001
	WASMEXecError_VM_ERR_GET_VM         ErrCode = 0x31002
	WASMEXecError_VM_ERR_FIND_VM_METHOD ErrCode = 0x31003
	WASMEXecError_VM_ERR_PARAM_COUNT    ErrCode = 0x31004
	WASMEXecError_VM_ERR_UNSUPPORT_TYPE ErrCode = 0x31005
	WASMEXecError_VM_ERR_EXEC_FAILED    ErrCode = 0x31006

	WASMEXecError_VM_ERR_OUT_OF_MEMORY           ErrCode = 0x31007
	WASMEXecError_VM_ERR_INVALID_PARAMETER_COUNT ErrCode = 0x31008
	WASMEXecError_VM_ERR_FAIL_EXECUTE_ENVFUNC    ErrCode = 0x31009
	WASMEXecError_VM_ERR_FAIL_STORAGE_MEMORY     ErrCode = 0x3100a
	WASMEXecError_VM_ERR_EXEC_TIME_OVER          ErrCode = 0x3100b
	WASMEXecError_VM_ERR_EXEC_PANIC              ErrCode = 0x3100c
	WASMEXecError_VM_ERR_EXEC_DEFER              ErrCode = 0x3100d

	/* js vm exe error : 0x32000~ 0x32fff */
	JSExecErrorStart ErrCode = 0x32000

	JSExecError_VM_CTX_INIT_FAIL    ErrCode = 0x32001
	JSExecError_VM_LOAD_CODE_FAIL   ErrCode = 0x32002
	JSExecError_VM_JS_CALL_FAIL     ErrCode = 0x32003
	JSExecError_VM_JS_RTN_FAIL      ErrCode = 0x32004
	JSExecError_VM_ADD_SUB_TRX_FAIL ErrCode = 0x32005
)

var (
	aaa = map[ErrCode]string{
		ErrNoError:                "success",
		ErrTrxPendingNumLimit:     "push trx: " + "trx pool busy",
		ErrTrxCacheNumLimit:       "push trx: " + "trx cache busy",
		ErrTrxSignError:           "push trx: " + "check signature error",
		ErrTrxAccountError:        "push trx: " + "check account valid error",
		ErrTrxLifeTimeError:       "push trx: " + "check life time error",
		ErrTrxUniqueError:         "push trx: " + "check trx unique error",
		ErrTrxChainMathError:      "push trx: " + "check match chain error",
		ErrTrxContractHanldeError: "push trx: " + "process contract error",
		ErrTrxContractDepthError:  "push trx: " + "contract depth error",
		ErrTrxSubContractNumError: "push trx: " + "sub contract num error",
		ErrTrxBlockSyncingError:   "push trx: " + "block syncing error",
		ErrTrxAlreadyInPool:       "push trx: " + "already in pool",
		ErrTrxAlreadyInCache:      "push trx: " + "already in cache",
		ErrTrxContractError:       "push trx: " + "check contract valid error",
		ErrTrxMethodError:         "push trx: " + "check method valid error",
		ErrTrxCheckSpaceError:     "push trx: " + "space token is not enough",
		ErrTrxCheckTimeError:      "push trx: " + "time token is not enough",
		ErrTrxCheckMinSpaceError:  "push trx: " + "check minimum space token failed",
		ErrTrxCheckMinTimeError:   "push trx: " + "check minimum time token failed",

		ErrTrxCheckSpaceInternalError:      "push trx: " + "check space token fee error",
		ErrTrxCheckTimeInternalError:       "push trx: " + "check time token fee error",
		ErrTrxCheckResourceInternalError:   "push trx: " + "check resource limit error",
		ErrTrxResourceExceedMaxSpacePerTrx: "push trx: " + "trx size Exceeding the maximum value",
		ErrTrxResourceCheckMinBalance:      "push trx: " + "check minimum balance failed",
		ErrTrxExecTimeOver:                 "push trx: " + "trx exec time over",
		ErrTrxVmTypeInvalid:                "push trx: " + "contract vm type error",
		ErrTrxNoticeContractNumError:       "push trx: " + "notice contract num error",
		ErrTrxContractNotExist:             "push trx: " + "contract not exist",
		ErrTrxCacheNumLimit:                "push trx: " + "trx cache busy",
		ErrTrxAlreadyInCache:               "push trx: " + "already in cache",
		ErrTrxVersionError:                 "push trx: " + "trx version not match",

		ErrAccountNameIllegal:             "push trx: " + "illegal account name",
		ErrAccountNotFound:                "push trx: " + "account name not found",
		ErrAccountAlreadyExist:            "push trx: " + "account name already exist",
		ErrContractParamParseError:        "push trx: " + "transaction param error",
		ErrContractInsufficientFunds:      "push trx: " + "account insufficient funds",
		ErrContractInvalidContractCode:    "push trx: " + "invalide contract code",
		ErrContractInvalidContractAbi:     "push trx: " + "invalide contract abi",
		ErrContractUnknownContract:        "push trx: " + "unknown contract",
		ErrContractUnknownMethod:          "push trx: " + "unknown contract method",
		ErrContractTransferOverflow:       "push trx: " + "transfer overflow",
		ErrAccountMismatch:                "push trx: " + "sender and param account mismatch",
		ErrContractInsufficientCredits:    "push trx: " + "insufficient credits",
		ErrAccountPubkeyLenIllegal:        "push trx: " + "pubkey len error",
		ErrContractGenesisPermissionError: "push trx: " + "no genesis node permission",
		ErrContractChainNotActivated:      "push trx: " + "chain does not activated",
		ErrContractTransferToSelf:         "push trx: " + "cannot transfer to self",
		ErrContractGrantToSelf:            "push trx: " + "cannot grant to self",

		ErrContractNumReachMaxPerAccount:           "push trx: " + "contract num reach max",
		ErrContractNotFound:                        "push trx: " + "contract not found",
		ErrContractNameIllegal:                     "push trx: " + "illegal contract name",
		ErrContractNoStakedVoteFunds:               "push trx: " + "must stake first before voting",
		ErrContractMustVoteToValidDelegate:         "push trx: " + "must vote to a valid delegate account",
		ErrContractDelegateVoteAlreadyImported:     "push trx: " + "transit delegate vote already imported",
		ErrContractBlockProducingAlreadyTransfered: "push trx: " + "block producing already transfered",
		ErrContractJSNotSupport:                    "push trx: " + "js contract not support for current chain",
		ErrContractAlreadyExist:                    "push trx: " + "contract name already exist",

		ErrContractInsufficientTransferValue:    "push trx: " + "insufficient transfer value",
		ErrContractAlreadyClaimedReward:         "push trx: " + "already claimed within past day",
		ErrContractUnstakeReleaseTimeNotReached: "push trx: " + "not reach unstake release time",
		ErrContractInsufficientUnstakeFunds:     "push trx: " + "insufficient unstake funds",

		ErrApiTrxNotFound:         "query trx: " + "trx not found",
		ErrApiBlockNotFound:       "query block: " + "block not found",
		ErrApiQueryChainInfoError: "query chain info: " + "error",
		ErrApiAccountNotFound:     "query account: " + "account not found",
		ErrApiObjectNotFound:      "query object: " + "object not found",
		ErrTransferCreditNotFound: "query credit: " + "credit not found",
		ErrApiAccountNameIllegal:  "query account: " + "account name illegal",

		ErrApiProposalNameIllegal: "query proposal: " + "proposal name illegal",

		ErrActorHandleError: "actor: " + "process error",

		ErrContractInsufficientStakeBalanceToRegDelegate: "push trx: " + "please stake at least 490,000 BTO to get qualification of delegate",

		ErrBlockInsertErrorGeneral:          "block insert general error",
		ErrBlockInsertErrorNotLinked:        "block not linked to previous block",
		ErrBlockInsertErrorValidateFail:     "block validate fail",
		ErrBlockInsertErrorDiffLibLinked:    "receive a lib block on another fork",
		ErrBlockInsertErrorDiffLibNotLinked: "receive a lib block on unknown fork",
		ErrBlockVersionError:                "block version not match",
		ErrConsensusReceiveMessageTimeout:   "consensus rsp timeout",

		ErrWalletPasswdError: "password of wallet is error",

		ErrMsignProposalNameIllegal:             "push trx: " + "illegal msign proposal",
		ErrMsignProposalNotFound:                "push trx: " + "msign proposal not found or excuted",
		ErrMsignProposalAlreadyExist:            "push trx: " + "msign proposal already exist",
		ErrMsignProposalApprovedNotFound:        "push trx: " + "msign proposal approved not found",
		ErrMsignProposalTransferIllegal:         "push trx: " + "illegal msign proposal transfer",
		ErrMsignProposalWeightNotEnough:         "push trx: " + "there is not enough authority to execute msign proposal",
		ErrMsignProposalNoAuthorityData:         "push trx: " + "msign proposal no authority data",
		ErrMsignProposalThresholdIllegal:        "push trx: " + "msign proposal threshold illegal",
		ErrMsignProposalAuthorityAccountIllegal: "push trx: " + "msign proposal authority account illegal",

		RestErrInternal:          "internal error",
		RestErrGenerateParm:      "generate parameter error",
		RestErrJsonNewEncoder:    "json NewEncoder or Encode error",
		RestErrBplMarshal:        "BPL Marshal Data error",
		RestErrDecodeStringError: "data decode string error",
		RestErrStringToBig:       "input is not valid data",

		RestErrReqNil:      "request Body is null",
		RestErrResultNil:   "result is null",
		RestErrPriKeyError: "check private key valid error",
		RestErrPubKeyError: "check public key valid error",
		RestErrHashError:   "check hash value invalid",

		RestErrTrxSignError:   "push trx: " + "signature Param error",
		RestErrUnlockError:    "unlock account error",
		RestErrUnkownAccType:  "unkown account type",
		RestErrWalletLocked:   "account is locked",
		RestErrWalletExist:    "account already exists",
		RestErrWalletNotExist: "account is not exists",

		RestErrTxPending:  "Trx is pending",
		RestErrTxNotFound: "Trx execute failed",
		RestErrTxPacked:   "Trx is packed",
		RestErrTxSending:  "Trx is sending",

		RestErrGetResSpaceError: "get space resource failed",
		RestErrGetResTimeError:  "get time resource failed",

		RestErrGetMsignTransferError: "Get multi sign transfer failed",

		RestErrLogItemInvalid: "log item invalid",

		ContractExecStart: "contract exec failed, error code: ",

		WASMEXecError_VM_ERR_CREATE_VM:      "failed to create a new VM instance",
		WASMEXecError_VM_ERR_GET_VM:         "failed to get a VM instance from memory",
		WASMEXecError_VM_ERR_FIND_VM_METHOD: "failed to find the method from the wasm module",
		WASMEXecError_VM_ERR_PARAM_COUNT:    "parameters count is not right",
		WASMEXecError_VM_ERR_UNSUPPORT_TYPE: "contract return type not support",
		WASMEXecError_VM_ERR_EXEC_FAILED:    "failed to call contract method",

		WASMEXecError_VM_ERR_OUT_OF_MEMORY:           "out of memory",
		WASMEXecError_VM_ERR_INVALID_PARAMETER_COUNT: "invalid parameter count",
		WASMEXecError_VM_ERR_FAIL_EXECUTE_ENVFUNC:    "failed to exec env func",
		WASMEXecError_VM_ERR_FAIL_STORAGE_MEMORY:     "failed to storeage memory",
		WASMEXecError_VM_ERR_EXEC_TIME_OVER:          "exec time over",
		WASMEXecError_VM_ERR_EXEC_PANIC:              "exec enter panic",
		WASMEXecError_VM_ERR_EXEC_DEFER:              "exec enter defer",

		JSExecError_VM_CTX_INIT_FAIL:    "failed to init context",
		JSExecError_VM_LOAD_CODE_FAIL:   "faild to load js code",
		JSExecError_VM_JS_CALL_FAIL:     "faild to call js contract method",
		JSExecError_VM_JS_RTN_FAIL:      "js exec failed",
		JSExecError_VM_ADD_SUB_TRX_FAIL: "failed to add sub transaction to list",
	}
)

// GetCodeString get code string
func GetCodeString(errorCode ErrCode) string {
	if ContractExecStart == errorCode&0xFF0000 {
		if 0 == errorCode&0xf000 {
			return aaa[ContractExecStart] + strconv.Itoa(int(errorCode&0xfff))
		} else {
			return aaa[errorCode]
		}
	} else {
		return aaa[errorCode]
	}
}
