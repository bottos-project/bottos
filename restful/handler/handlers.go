package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"errors"

	"github.com/gorilla/mux"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/action/message"
	"github.com/AsynkronIT/protoactor-go/actor"


	"time"

	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/role"
	service "github.com/bottos-project/bottos/action/actor/api"
	log "github.com/cihub/seelog"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/config"
	"regexp"
	"runtime"
)

//ApiService is actor service
type ApiService struct {
	env *env.ActorEnv
}

var roleIntf role.RoleInterface
//SetChainActorPid set chain actor pid
func SetRoleIntf(tpid role.RoleInterface) {
	roleIntf = tpid
}

var chainActorPid *actor.PID

//SetChainActorPid set chain actor pid
func SetChainActorPid(tpid *actor.PID) {
	chainActorPid = tpid
}

var trxPreHandleActorPid *actor.PID

//SetTrxActorPid set trx actor pid
func SetTrxPreHandleActorPid(tpid *actor.PID) {
	trxPreHandleActorPid = tpid
}

/*
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := ResponseStructs{
		ResponseStruct{Msg: "Write presentation"},
		ResponseStruct{Msg: "Host meetup"},
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintf(w, "Todo show: %s\n", todoId)
}*/

//Node
func GetGenerateBlockTime(w http.ResponseWriter, r *http.Request) {
	/*	//func GetGenerateBlockTime(cmd map[string]interface{}) map[string]interface{} {
		resp := ResponsePack(error.SUCCESS)
		resp["Result"] = "aq"
		//fmt.Fprint(w, "Welcome!\n",resp	)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}*/
	//resp["Result"] = config.DEFAULT_GEN_BLOCK_TIME
	//return resp
}

//GetInfo query chain info
func GetInfo(w http.ResponseWriter, r *http.Request) {
	msgReq := &message.QueryChainInfoReq{}
	var resp comtool.ResponseStruct

	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		log.Errorf("REST:chain Actor Request failed,%v", err)
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(res, 1); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	response := res.(*message.QueryChainInfoResp)
	if response.Error != nil {
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
		encoderRestResponse(w, resp)
		return
	}

	result := &api.GetInfoResponse_Result{}
	result.HeadBlockVersion = response.HeadBlockVersion
	result.HeadBlockNum = response.HeadBlockNum
	result.LastConsensusBlockNum = response.LastConsensusBlockNum
	result.HeadBlockHash = response.HeadBlockHash.ToHexString()
	result.HeadBlockTime = response.HeadBlockTime
	result.HeadBlockDelegate = response.HeadBlockDelegate
	result.CursorLabel = response.HeadBlockHash.Label()
	result.ChainId = common.BytesToHex(config.GetChainID())

	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	resp.Result = result
	encoderRestResponse(w, resp)
}

//checkNil check param is or not Nil,flag 0:request; 1:response
func checkNil(req interface{}, flag int8) comtool.ResponseStruct {
	var resp comtool.ResponseStruct

	if req == nil {
		if flag == 0 {
			resp.Errcode = uint32(bottosErr.RestErrReqNil)
			resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrReqNil)
		} else {
			resp.Errcode = uint32(bottosErr.RestErrResultNil)
			resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrResultNil)
		}

		funcName, _, _, _ := runtime.Caller(1)
		log.Errorf("REST:check param is nil,%s errcode: %d, msg:%s", runtime.FuncForPC(funcName).Name(), resp.Errcode, resp.Msg)
		return resp
	}
	return resp
}

//GetBlock query block
func GetBlock(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	var msgReq *api.GetBlockRequest
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("REST:json Decoder failed:%v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err
		encoderRestResponse(w, resp)
		return
	}
	if resp := checkNil(msgReq, 0); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	msgReq2 := &message.QueryBlockReq{BlockHash: common.HexToHash(msgReq.BlockHash),
		BlockNumber: msgReq.BlockNum}

	res, err := chainActorPid.RequestFuture(msgReq2, 500*time.Millisecond).Result()
	if err != nil {
		log.Errorf("REST:chain Actor Request failed,%v", err)
		resp.Errcode = uint32(bottosErr.ErrApiBlockNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiBlockNotFound)
		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(res, 1); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	response := res.(*message.QueryBlockResp)
	if response.Block == nil {
		resp.Errcode = uint32(bottosErr.ErrApiBlockNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiBlockNotFound)
		encoderRestResponse(w, resp)
		return
	}

	//result := &api.GetBlockResponse_Result{}
	result := &types.BlockDetail{}
	hash := response.Block.Hash()
	result.BlockVersion = response.Block.GetVersion()
	result.PrevBlockHash = response.Block.GetPrevBlockHash().ToHexString()
	result.BlockNum = response.Block.GetNumber()
	result.BlockHash = hash.ToHexString()
	result.CursorBlockLabel = hash.Label()
	result.BlockTime = response.Block.GetTimestamp()
	result.TrxMerkleRoot = response.Block.ComputeMerkleRoot().ToHexString()
	result.Delegate = string(response.Block.GetDelegate())
	result.DelegateSign = common.BytesToHex(response.Block.GetDelegateSign())
	for _, v := range response.Block.BlockTransactions {
		tx := convertIntTrxToApiTrxInter(v, roleIntf)
		result.Trxs = append(result.Trxs, &tx)
	}

	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	resp.Result = result
	encoderRestResponse(w, resp)
}

//SendTransaction send transaction
func SendTransaction(w http.ResponseWriter, r *http.Request) {
	var trx *api.Transaction
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&trx)
	if err != nil {
		log.Errorf("REST:json Decoder failed:%v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(trx, 0); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	//verity Sender
	if !re.MatchString(trx.Sender) {
		log.Errorf("REST:match sender failed,sender:%v", trx.Sender)
		resp.Errcode = uint32(bottosErr.ErrTrxAccountError)
		resp.Msg = bottosErr.GetCodeString((bottosErr.ErrCode)(resp.Errcode))
		encoderRestResponse(w, resp)
		return
	}

	//verity Contract
	if !contractRe.MatchString(trx.Contract) && trx.Contract != config.BOTTOS_CONTRACT_NAME {
		log.Errorf("REST:match method failed,Contract:%v", trx.Contract)
		resp.Errcode = uint32(bottosErr.ErrTrxAccountError)
		resp.Msg = bottosErr.GetCodeString((bottosErr.ErrCode)(resp.Errcode))
		encoderRestResponse(w, resp)
		return
	}

	//verity Method
	if !re.MatchString(trx.Method) {
		log.Errorf("REST:match method failed,Method:%v", trx.Method)
		resp.Errcode = uint32(bottosErr.ErrTrxAccountError)
		resp.Msg = bottosErr.GetCodeString((bottosErr.ErrCode)(resp.Errcode))
		encoderRestResponse(w, resp)
		return
	}

	intTrx, err := service.ConvertApiTrxToIntTrx(trx)
	if err != nil {
		log.Errorf("REST:Convert ApiTrx to IntTrx failed:%v", err)
		resp.Errcode = uint32(bottosErr.RestErrInternal)
		resp.Msg = err.Error()
		encoderRestResponse(w, resp)
		return
	}

	reqMsg := &message.PushTrxReq{
		Trx: intTrx,
	}

	handlerErr, err := trxPreHandleActorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result() // await result

	if nil != err {
		log.Errorf("REST:trxPreHandle %x Actor process failed,", intTrx.Hash(), err)

		resp.Errcode = uint32(bottosErr.ErrActorHandleError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrActorHandleError)
		encoderRestResponse(w, resp)
		return
	}

	result := &api.SendTransactionResponse_Result{}
	if bottosErr.ErrNoError == handlerErr {
		resp.Errcode = uint32(bottosErr.ErrNoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)

		result.TrxHash = intTrx.Hash().ToHexString()
		result.Trx = service.ConvertIntTrxToApiTrx(intTrx)
		resp.Result = result
		encoderRestResponse(w, resp)
		return
	} else {
		result.TrxHash = intTrx.Hash().ToHexString()
		result.Trx = service.ConvertIntTrxToApiTrx(intTrx)
		resp.Result = result
		//resp.Msg = handlerErr.(string)GetCodeString
		//resp.Msg = "to be add detail error description"
		var tempErr bottosErr.ErrCode
		tempErr = handlerErr.(bottosErr.ErrCode)

		resp.Errcode = (uint32)(tempErr)
		resp.Msg = bottosErr.GetCodeString(tempErr)
		encoderRestResponse(w, resp)
		return
	}

	log.Infof("REST:trx hash:%v, response: %s", result.TrxHash, resp.Msg)

	encoderRestResponse(w, resp)
}

type reqStruct struct {
	TrxHash string `json:"trx_hash,omitemty"`
}

type BlockTransaction struct {
	Transaction     *Transaction
	ResourceReceipt *types.ResourceReceipt
	TrxHash         string
}

type Transaction struct {
	Version     uint32      `json:"version"`
	CursorNum   uint64      `json:"cursor_num"`
	CursorLabel uint32      `json:"cursor_label"`
	Lifetime    uint64      `json:"lifetime"`
	Sender      string      `json:"sender"`
	Contract    string      `json:"contract"`
	Method      string      `json:"method"`
	Param       interface{} `json:"param"`
	SigAlg      uint32      `json:"sig_alg"`
	Signature   string      `json:"signature"`
}

func getContractAbi(r role.RoleInterface, contractName string) (*abi.ABI, error) {
	contract, err := r.GetContract(contractName)
	if err != nil {
		log.Errorf("REST:GetContract failed,%v", err)
		return nil, errors.New("Get contract fail")
	}

	Abi, err := abi.ParseAbi(contract.ContractAbi)
	if err != nil {
		log.Errorf("REST:ParseAbi failed, %v", err)
		return nil, err
	}

	return Abi, nil
}

func ParseTransactionParam(r role.RoleInterface, Param []byte, Contract string, Method string) (interface{}, error) {
	var Abi *abi.ABI = nil
	if Contract != config.BOTTOS_CONTRACT_NAME {
		var err error
		Abi, err = getContractAbi(r, Contract)
		if err != nil {
			log.Errorf("REST:getContractAbi failed, %v", err)
			return nil, errors.New("External Abi is empty!")
		}
	} else {
		Abi = abi.GetAbi()
	}

	if Abi == nil {
		return nil, errors.New("Abi is empty!")
	}

	decodedParam, isOK := abi.UnmarshalAbiEx(Contract, Abi, Method, Param)
	if decodedParam == nil || len(decodedParam) <= 0 {
		if isOK {
			decodedParam = make(map[string]interface{})
		} else {
			return nil, errors.New("ParseTransactionParam: FAILED")
		}
	}
	return decodedParam, nil
}

//convertIntTrxToApiTrxInter convert IntTrx to Api TrxInter
func convertIntTrxToApiTrxInter(trxB *types.BlockTransaction, r role.RoleInterface) interface{} {
	trx := trxB.Transaction
	var apiTrx *BlockTransaction

	parmConvered, err := ParseTransactionParam(r, trx.Param, trx.Contract, trx.Method)
	if err != nil {
		log.Errorf("REST:ParseTransactionParam failed, %v", err)
		return apiTrx
	}

	apiTx := &Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       parmConvered,
		SigAlg:      trx.SigAlg,
		Signature:   common.BytesToHex(trx.Signature),
	}

	apiTrx = &BlockTransaction{
		Transaction:     apiTx,
		ResourceReceipt: trxB.ResourceReceipt,
		TrxHash:         trx.Hash().ToHexString(),
	}
	return apiTrx
}

//GetTransaction get transaction by Trx hash
func GetTransaction(w http.ResponseWriter, r *http.Request) {
	var req *reqStruct
	var resp comtool.ResponseStruct

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("REST:json Decoder failed: %v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(req, 0); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	msgReq := &message.QueryTrxReq{
		TrxHash: common.HexToHash(req.TrxHash),
	}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		log.Errorf("REST:chainActor process failed: %v", err)
		resp.Errcode = uint32(bottosErr.ErrActorHandleError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrActorHandleError)
		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(res, 1); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	response := res.(*message.QueryBlockTrxResp)
	if response.Trx == nil {
		resp.Errcode = uint32(bottosErr.ErrApiTrxNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiTrxNotFound)
		encoderRestResponse(w, resp)
		return
	}

	//resp.Result = service.ConvertIntTrxToApiTrx(response.Trx)
	resp.Result = convertIntTrxToApiTrxInter(response.Trx, roleIntf)

	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	encoderRestResponse(w, resp)
}

type TransactionStatus struct {
	Status string `json:"status"`
}

func GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
	var req *reqStruct
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("REST:json Decoder failed: %v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(req, 0); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	tx := actorenv.Chain.GetCommittedTransaction(common.HexToHash(req.TrxHash))
	if tx != nil {
		resp.Errcode = uint32(bottosErr.ErrNoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
		resp.Result = &TransactionStatus{Status: "committed"}
		encoderRestResponse(w, resp)
		return
	}

	tx = actorenv.Chain.GetTransaction(common.HexToHash(req.TrxHash))
	if tx != nil {
		resp.Errcode = uint32(bottosErr.RestErrTxPacked)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrTxPacked)
		resp.Result = &TransactionStatus{Status: "packed"}
		encoderRestResponse(w, resp)
		return
	}

	trxApply := transaction.NewTrxApplyService()
	errCode := trxApply.GetTrxErrorCode(common.HexToHash(req.TrxHash))
	if bottosErr.ErrNoError != errCode {
		log.Errorf("REST:get trx error code:%v",errCode)

		//resp.Errcode = uint32(errCode)
		resp.Errcode = uint32(bottosErr.ErrNoError)
		//resp.Msg = bottosErr.GetCodeString(errCode)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
		resp.Result = &TransactionStatus{Status: bottosErr.GetCodeString(errCode)}
		encoderRestResponse(w, resp)
		return
	}

	isInPool := trxApply.IsTrxInPendingPool(common.HexToHash(req.TrxHash))
	if true == isInPool {
		resp.Errcode = uint32(bottosErr.RestErrTxPending)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrTxPending)
		resp.Result = &TransactionStatus{Status: "pending"}
		encoderRestResponse(w, resp)
		return
	} else {
		resp.Errcode = uint32(bottosErr.RestErrTxNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrTxNotFound)
		resp.Result = &TransactionStatus{Status: "not found"}
		encoderRestResponse(w, resp)
		return
	}
}

//GetAccountBrief query account public key
func GetAccountBrief(w http.ResponseWriter, r *http.Request) {
	var msgReq api.GetAccountBriefRequest
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("REST:json Decoder failed: %v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(msgReq, 0); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}
	name := msgReq.AccountName

	result := &api.GetAccountBriefResponse_Result{}

	accountType, _ := common.AnalyzeName(name)
	if common.NameTypeUnknown == accountType {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNameIllegal)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNameIllegal)
		encoderRestResponse(w, resp)
		return
	} else if common.NameTypeAccount == accountType {
		account, err := roleIntf.GetAccount(name)
		if err != nil {
			resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
			resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)

			encoderRestResponse(w, resp)
			return
		}
		if resp := checkNil(account, 1); resp.Errcode != 0 {
			encoderRestResponse(w, resp)
			return
		}

		balance, err := roleIntf.GetBalance(name)
		if err != nil {
			log.Errorf("DB:GetBalance failed: %v", err)

			resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
			resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
			encoderRestResponse(w, resp)
			return
		}

		result.AccountName = name
		result.Pubkey = common.BytesToHex(account.PublicKey)
		result.Balance = balance.Balance.String()
	}

	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	resp.Result = result
	encoderRestResponse(w, resp)
}

//GetAccount query account info
func GetAccount(w http.ResponseWriter, r *http.Request) {
	var msgReq api.GetAccountRequest
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
		return
	}
	name := msgReq.AccountName

	account, err := roleIntf.GetAccount(name)
	var resp ResponseStruct
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	balance, err := roleIntf.GetBalance(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	stakedBalance, err := roleIntf.GetStakedBalance(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	result := &api.GetAccountResponse_Result{}
	result.AccountName = name
	result.Pubkey = common.BytesToHex(account.PublicKey)
	result.Balance = balance.Balance.String()
	result.StakedBalance = stakedBalance.StakedBalance.String()
	resp.Result=result
	resp.Errcode = 0

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func GetKeyValue(w http.ResponseWriter, r *http.Request) {
	var req *api.GetKeyValueRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
		return
	}

	contract := req.Contract
	object := req.Object
	key := req.Key
	value, err := roleIntf.GetBinValue(contract, object, key)
	var resp ResponseStruct
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiObjectNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiObjectNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	result := &api.GetKeyValueResponse_Result{}
	result.Contract = contract
	result.Object = object
	result.Key = key
	result.Value = common.BytesToHex(value)
	resp.Result = result
	resp.Errcode = 0

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func GetContractAbi(w http.ResponseWriter, r *http.Request) {
	var req *api.GetAbiRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
		return
	}
	//contract := req.Contract
	account, err := roleIntf.GetAccount(req.Contract)
	var resp ResponseStruct
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	if len(account.ContractAbi) > 0 {
		resp.Result = string(account.ContractAbi)
	} else {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func GetContractCode(w http.ResponseWriter, r *http.Request) {
	var req *api.GetAbiRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
		return
	}
	//contract := req.Contract
	account, err := roleIntf.GetAccount(req.Contract)
	var resp ResponseStruct
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	if len(account.ContractCode) > 0 {
		resp.Result = common.BytesToHex(account.ContractCode)
	} else {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func GetTransferCredit(w http.ResponseWriter, r *http.Request) {
	var req *api.GetTransferCreditRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
		return
	}
	name := req.Name
	spender := req.Spender
	credit, err := roleIntf.GetTransferCredit(name, spender)
	var resp ResponseStruct
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrTransferCreditNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrTransferCreditNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	result := &api.GetTransferCreditResponse_Result{}
	result.Name = credit.Name
	result.Spender = credit.Spender
	result.Limit = credit.Limit.String()
	resp.Result = result

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

//GetPeers get all peers
func ConnetPeerbyAddress(w http.ResponseWriter, r *http.Request) {
	var msgReq api.ConnectPeerByAddressRequest
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("REST:json Decoder failed: %v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		encoderRestResponse(w, resp)
		return
	}

	address := msgReq.Address
	isDone := actorenv.Protocol.UpdatePeerStateToActive(address)



	result := &api.ConnectPeerByAddressResponse_Result{
		Result: isDone,
	}
	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	resp.Result = result
	encoderRestResponse(w, resp)
}

//GetPeers get all peers
func DisConnectPeerbyAddress(w http.ResponseWriter, r *http.Request) {
	var msgReq api.DisconnectPeerByAddressRequest
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("REST:json Decoder failed: %v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		encoderRestResponse(w, resp)
		return
	}

	address := msgReq.Address
	isDone := actorenv.Protocol.UpdatePeerStateToUnActive(address)



	result := &api.DisconnectPeerByAddressResponse_Result{
		Result: isDone,
	}
	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	resp.Result = result
	encoderRestResponse(w, resp)
}


//GetPeers get all peers
func GetPeerStatebyAddress(w http.ResponseWriter, r *http.Request) {
	var msgReq api.GetPeerStateByAddressRequest
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("REST:json Decoder failed: %v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		encoderRestResponse(w, resp)
		return
	}

	address := msgReq.Address
	state := actorenv.Protocol.QueryPeerState(address)



	result := &api.GetPeerStateByAddressResponse_Result{
		IsActive: state,
	}
	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	resp.Result = result
	encoderRestResponse(w, resp)
}
func encoderRestResponse(w http.ResponseWriter, resp comtool.ResponseStruct) http.ResponseWriter {
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Expose-Headers", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	//w.Header().Set("Access-Control-Allow-Methods", "POST")

	//w.WriteHeader(404)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)

		funcName, _, _, _ := runtime.Caller(1)
		log.Errorf("REST:json encoder failed,%s errcode: %d json.NewEncoder(w).Encode(resp) error:%s", runtime.FuncForPC(funcName).Name(), resp.Errcode, err)
	}

	return w
}

func encoderRestRequest(r *http.Request, req interface{}) (interface{}, error) {
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("request error: %s", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		funcName, _, _, _ := runtime.Caller(1)
		log.Errorf("%s errcode: %d json.NewEncoder(w).Encode(resp) error:%s", runtime.FuncForPC(funcName).Name(), resp.Errcode, err)

		return resp, err
	}

	if req == nil {
		resp.Errcode = uint32(bottosErr.RestErrReqNil)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrReqNil)

		funcName, _, _, _ := runtime.Caller(1)
		log.Errorf("%s errcode: %d json.NewEncoder(w).Encode(resp) error:%s", runtime.FuncForPC(funcName).Name(), resp.Errcode, err)
		return resp, errors.New("request is nil")
	}

	return req, nil
}

func GetTrxHashForSign(sender, contract, method string, param []byte, h *api.GetInfoResponse) ([]byte, *types.Transaction, error) {
	var blockHeader *api.GetInfoResponse_Result
	if h == nil {
		var err error
		blockHeader, err = GetBlockHeader()
		if err != nil {
			return nil, nil, err
		}
	} else {
		blockHeader = h.Result
	}

	trx := &types.BasicTransaction{
		Version:     blockHeader.HeadBlockVersion,
		CursorNum:   blockHeader.HeadBlockNum,
		CursorLabel: blockHeader.CursorLabel,
		Lifetime:    blockHeader.HeadBlockTime + 100,
		Sender:      sender,
		Contract:    contract,
		Method:      method,
		Param:       param,
		SigAlg:      config.SIGN_ALG,
	}
	msg, err := bpl.Marshal(trx)
	if nil != err {
		log.Errorf("REST:bpl Marshal failed: %v", err)
		return nil, nil, err
	}

	//Add chainID Flag
	chainID, _ := hex.DecodeString(blockHeader.ChainId)
	msg = bytes.Join([][]byte{msg, chainID}, []byte{})

	intTrx := &types.Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       trx.Param,
		SigAlg:      config.SIGN_ALG,
	}
	return comtool.Sha256(msg), intTrx, err
}

func PushTrx(intTrx *types.Transaction) (comtool.ResponseStruct, error) {
	reqMsg := &message.PushTrxReq{
		Trx: intTrx,
	}

	handlerErr, err := trxPreHandleActorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result() // await result
	var resp comtool.ResponseStruct
	if nil != err {
		resp.Errcode = uint32(bottosErr.ErrActorHandleError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrActorHandleError)

		log.Errorf("REST:trx PreHandleActor: %x actor process failed,%v", reqMsg.Trx.Hash(), err)
		return resp, err
	}

	result := &api.SendTransactionResponse_Result{}
	if bottosErr.ErrNoError == handlerErr {
		result.TrxHash = reqMsg.Trx.Hash().ToHexString()
		result.Trx = service.ConvertIntTrxToApiTrx(reqMsg.Trx)
		resp.Result = result
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
		resp.Errcode = uint32(bottosErr.ErrNoError)
	} else {
		result.TrxHash = reqMsg.Trx.Hash().ToHexString()
		result.Trx = service.ConvertIntTrxToApiTrx(reqMsg.Trx)
		resp.Result = result
		//resp.Msg = handlerErr.(string)GetCodeString
		//resp.Msg = "to be add detail error description"
		var tempErr bottosErr.ErrCode
		tempErr = handlerErr.(bottosErr.ErrCode)

		resp.Errcode = (uint32)(tempErr)
		resp.Msg = bottosErr.GetCodeString((bottosErr.ErrCode)(resp.Errcode))
	}

	log.Infof("REST:PushTrx trx: %v %s", result.TrxHash, resp.Msg)

	return resp, nil
}

func JsonToBin(w http.ResponseWriter, r *http.Request) {

	//info := fmt.Sprintln(r.Header.Get("Content-Type"))
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	//fmt.Fprintln(w, info, string(body))

	/*	var data map[string]json.RawMessage
		err := json.Unmarshal([]byte(r.Body), &data)
		if err != nil {
			fmt.Println(err)
		}*/

	param, err := bpl.Marshal(body)

	var resp comtool.ResponseStruct
	if err != nil {
		resp.Errcode = uint32(bottosErr.RestErrBplMarshal)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrBplMarshal)
		resp.Result = err

		funcName, _, _, _ := runtime.Caller(1)
		log.Errorf("%s errcode: %d bpl.Marshal error:%s", runtime.FuncForPC(funcName).Name(), resp.Errcode, err)
		encoderRestResponse(w, resp)
		return
	}
	if resp := checkNil(body, 0); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	resp.Result = hex.EncodeToString(param)
	encoderRestResponse(w, resp)
	return
}

func GetContract(contractName string) (*role.Contract, error) {
	return roleIntf.GetContract(contractName)
}

func ReviewProposal(w http.ResponseWriter, r *http.Request) {
	var msgReq api.ReviewProposalRequest
	var resp comtool.ResponseStruct
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("REST:json Decoder failed: %v", err)
		resp.Errcode = uint32(bottosErr.RestErrJsonNewEncoder)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrJsonNewEncoder)
		resp.Result = err

		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(msgReq, 0); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}

	if !common.CheckAccountNameContent(msgReq.ProposalName) {
		resp.Errcode = uint32(bottosErr.ErrApiProposalNameIllegal)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiProposalNameIllegal)
		encoderRestResponse(w, resp)
		return
		if !common.CheckAccountNameContent(msgReq.Proposer) {
			resp.Errcode = uint32(bottosErr.ErrApiAccountNameIllegal)
			resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNameIllegal)
			encoderRestResponse(w, resp)
			return
		}
	}

	msignTransfer, err := roleIntf.GetMsignTransfer(msgReq.ProposalName)
	if err != nil {
		log.Errorf("REST:get Multi sign transfer failed: %v", err)

		resp.Errcode = uint32(bottosErr.RestErrGetMsignTransferError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrGetMsignTransferError)
		encoderRestResponse(w, resp)
		return
	}

	if resp := checkNil(msignTransfer, 1); resp.Errcode != 0 {
		encoderRestResponse(w, resp)
		return
	}
	if msignTransfer.ProposerName != msgReq.Proposer {
		log.Errorf("REST:mutli sign proposal:%+v", msignTransfer)
		resp.Errcode = uint32(bottosErr.ErrMsignProposalNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrMsignProposalNotFound)
		encoderRestResponse(w, resp)
		return
	}

	var authorList = []*api.AuthorList{}
	for _, v := range msignTransfer.RequestList {
		authorList = append(authorList, &api.AuthorList{
			AuthorAccount: v.AuthorAccount,
			IsApproved:    v.IsApproved,
		})
	}

	var param = &api.MsignTransferParam{}
	err = json.Unmarshal(msignTransfer.PackedTransaction, &param)
	if err != nil {
		log.Errorf("REST:Unmarshal failed: %v", err)

		resp.Errcode = uint32(bottosErr.RestErrBplMarshal)
		resp.Msg = bottosErr.GetCodeString(bottosErr.RestErrBplMarshal)
		encoderRestResponse(w, resp)
		return
	}

	result := &api.ReviewProposalResponse_Result{
		ProposalName:      msignTransfer.ProposalName,
		Proposer:          msignTransfer.ProposerName,
		MsignAccountName:  msignTransfer.MsignAccountName,
		AuthorList:        authorList,
		PackedTransaction: common.BytesToHex(msignTransfer.PackedTransaction),
		Transaction:       param,
		Available:         msignTransfer.Available,
		Time:              msignTransfer.Time,
	}

	resp.Errcode = uint32(bottosErr.ErrNoError)
	resp.Msg = bottosErr.GetCodeString(bottosErr.ErrNoError)
	resp.Result = result
	encoderRestResponse(w, resp)
}
