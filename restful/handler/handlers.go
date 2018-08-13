package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/bottos-project/bottos/action/env"
	//"../error"
	//"github.com/bottos-project/bottos/restful/error"
	"github.com/bottos-project/bottos/action/message"
	"github.com/AsynkronIT/protoactor-go/actor"

	"time"

	bottosErr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/role"
	service "github.com/bottos-project/bottos/action/actor/api"
	log "github.com/cihub/seelog"
	"github.com/bottos-project/bottos/contract/contractdb"
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

var contractDbIns *contractdb.ContractDB
//SetChainActorPid set chain actor pid
func SetContractDbIns(tpid *contractdb.ContractDB) {
	contractDbIns = tpid
}

var chainActorPid *actor.PID

//SetChainActorPid set chain actor pid
func SetChainActorPid(tpid *actor.PID) {
	chainActorPid = tpid
}

var trxactorPid *actor.PID

//SetTrxActorPid set trx actor pid
func SetTrxActorPid(tpid *actor.PID) {
	trxactorPid = tpid
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Msg: "Write presentation"},
		Todo{Msg: "Host meetup"},
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintf(w, "Todo show: %s\n", todoId)
}

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
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()

	var resp Todo
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	response := res.(*message.QueryChainInfoResp)
	if response.Error != nil {
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	result := &api.GetInfoResponse_Result{}
	result.HeadBlockNum = response.HeadBlockNum
	result.LastConsensusBlockNum = response.LastConsensusBlockNum
	result.HeadBlockHash = response.HeadBlockHash.ToHexString()
	result.HeadBlockTime = response.HeadBlockTime
	result.HeadBlockDelegate = response.HeadBlockDelegate
	result.CursorLabel = response.HeadBlockHash.Label()
	resp.Result = result

	resp.Errcode = 0
	json.NewEncoder(w).Encode(resp)
}

//GetBlock query block
func GetBlock(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	var msgReq *message.QueryBlockReq
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
	}

	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	var resp Todo
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiBlockNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiBlockNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	response := res.(*message.QueryBlockResp)
	if response.Block == nil {
		resp.Errcode = uint32(bottosErr.ErrApiBlockNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiBlockNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	result := &api.GetBlockResponse_Result{}
	hash := response.Block.Hash()
	result.PrevBlockHash = response.Block.GetPrevBlockHash().ToHexString()
	result.BlockNum = response.Block.GetNumber()
	result.BlockHash = hash.ToHexString()
	result.CursorBlockLabel = hash.Label()
	result.BlockTime = response.Block.GetTimestamp()
	result.TrxMerkleRoot = response.Block.ComputeMerkleRoot().ToHexString()
	result.Delegate = string(response.Block.GetDelegate())
	result.DelegateSign = response.Block.GetDelegateSign().ToHexString()
	resp.Result = result

	resp.Errcode = 0
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func SendTransaction(w http.ResponseWriter, r *http.Request) {
	var trx *api.Transaction
	err := json.NewDecoder(r.Body).Decode(&trx)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
	}

	if trx == nil {
		//rsp.retCode = ??
		if err := json.NewEncoder(w).Encode(trx); err != nil {
			panic(err)
		}
	}

	intTrx, err := service.ConvertApiTrxToIntTrx(trx)
	if err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	reqMsg := &message.PushTrxReq{
		Trx: intTrx,
	}

	handlerErr, err := trxactorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result() // await result

	var resp Todo
	if nil != err {
		resp.Errcode = uint32(bottosErr.ErrActorHandleError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrActorHandleError)

		log.Errorf("trx: %x actor process failed", intTrx.Hash(), )

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	result := &api.SendTransactionResponse_Result{}
	if bottosErr.ErrNoError == handlerErr {
		result.TrxHash = intTrx.Hash().ToHexString()
		result.Trx = service.ConvertIntTrxToApiTrx(intTrx)
		resp.Result = result
		resp.Msg = "trx receive succ"
		resp.Errcode = 0
	} else {
		result.TrxHash = intTrx.Hash().ToHexString()
		result.Trx = service.ConvertIntTrxToApiTrx(intTrx)
		resp.Result = result
		//resp.Msg = handlerErr.(string)GetCodeString
		//resp.Msg = "to be add detail error description"
		var tempErr bottosErr.ErrCode
		tempErr = handlerErr.(bottosErr.ErrCode)

		resp.Errcode = (uint32)(tempErr)
		resp.Msg = bottosErr.GetCodeString((bottosErr.ErrCode)(resp.Errcode))
	}

	log.Infof("trx: %v %s", result.TrxHash, resp.Msg)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

type reqStruct struct {
	TrxHash string `json:"trx_hash,omitemty"`
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	var req *reqStruct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
	}

	msgReq := &message.QueryTrxReq{
		TrxHash: common.HexToHash(req.TrxHash),
}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	var resp Todo
	if err != nil {
		resp.Errcode = 1
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	response := res.(*message.QueryTrxResp)
	if response.Trx == nil {
		resp.Errcode = uint32(bottosErr.ErrApiTrxNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiTrxNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	//resp.Result = service.ConvertIntTrxToApiTrx(response.Trx)
	role := &role.Role{Db: contractDbIns.Db}
	resp.Result = service.ConvertIntTrxToApiTrxInter(response.Trx, role)

	resp.Errcode = uint32(bottosErr.ErrNoError)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

//GetAccount query account info
func GetAccount(w http.ResponseWriter, r *http.Request) {
	var msgReq api.GetAccountRequest
	err := json.NewDecoder(r.Body).Decode(&msgReq)
	if err != nil {
		log.Errorf("request error: %s", err)
		panic(err)
	}
	name := msgReq.AccountName

	account, err := roleIntf.GetAccount(name)
	var resp Todo
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	balance, err := roleIntf.GetBalance(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	stakedBalance, err := roleIntf.GetStakedBalance(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	result := &api.GetAccountResponse_Result{}
	result.AccountName = name
	result.Pubkey = common.BytesToHex(account.PublicKey)
	result.Balance = balance.Balance
	result.StakedBalance = stakedBalance.StakedBalance
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
	}

	contract := req.Contract
	object := req.Object
	key := req.Key
	value, err := contractDbIns.GetBinValue(contract, object, key)
	var resp Todo
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiObjectNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiObjectNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
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
	}
	//contract := req.Contract
	account, err := roleIntf.GetAccount(req.Contract)
	var resp Todo
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
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
	}
	//contract := req.Contract
	account, err := roleIntf.GetAccount(req.Contract)
	var resp Todo
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
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
	}
	name := req.Name
	spender := req.Spender
	credit, err := roleIntf.GetTransferCredit(name, spender)
	var resp Todo
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrTransferCreditNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrTransferCreditNotFound)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}

	result := &api.GetTransferCreditResponse_Result{}
	result.Name = credit.Name
	result.Spender = credit.Spender
	result.Limit = credit.Limit
	resp.Result = result

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}
