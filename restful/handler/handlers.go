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
		Todo{Name: "Write presentation"},
		Todo{Name: "Host meetup"},
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
	_ = json.NewDecoder(r.Body).Decode(&msgReq)

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

}
func GetTransaction(w http.ResponseWriter, r *http.Request) {

}
//GetAccount query account info
func GetAccount(w http.ResponseWriter, r *http.Request) {
	var msgReq api.GetAccountRequest
	_ = json.NewDecoder(r.Body).Decode(&msgReq)
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

}

func GetAbi(w http.ResponseWriter, r *http.Request) {

}

func GetTransferCredit(w http.ResponseWriter, r *http.Request) {

}
