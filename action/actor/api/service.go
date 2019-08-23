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
 * file description:  trx agent
 * @Author:
 * @Date:   2017-12-13
 * @Last Modified by:
 * @Last Modified time:
 */

package apiactor

import (
	"context"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/bottos/action/env"
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"

	bottosErr "github.com/bottos-project/bottos/common/errors"
	log "github.com/cihub/seelog"
	"encoding/hex"
		"github.com/bottos-project/crypto-go/crypto"
	)

//ApiService is actor service
type ApiService struct {
	env *env.ActorEnv
}

//NewApiService new api service
func NewApiService(env *env.ActorEnv) api.ChainHandler {
	apiService := &ApiService{env: env}
	return apiService
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

func ConvertApiTrxToIntTrx(trx *api.Transaction) (*types.Transaction, error) {
	param, err := common.HexToBytes(trx.Param)
	if err != nil {
		return nil, err
	}

	signature, err := common.HexToBytes(trx.Signature)
	if err != nil {
		return nil, err
	}

	intTrx := &types.Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       param,
		SigAlg:      trx.SigAlg,
		Signature:   signature,
	}

	return intTrx, nil
}

func ConvertIntTrxToApiTrx(trx *types.Transaction) *api.Transaction {
	apiTrx := &api.Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       common.BytesToHex(trx.Param),
		SigAlg:      trx.SigAlg,
		Signature:   common.BytesToHex(trx.Signature),
	}

	return apiTrx
}

//SendTransaction push trx
func (a *ApiService) SendTransaction(ctx context.Context, trx *api.Transaction, resp *api.SendTransactionResponse) error {
	if trx == nil {
		//rsp.retCode = ??
		return nil
	}

	log.Info("rcv trx, detail: ", trx)

	intTrx, err := ConvertApiTrxToIntTrx(trx)
	if err != nil {
		return nil
	}

	reqMsg := &message.PushTrxReq{
		Trx: intTrx,
	}

	start := common.MeasureStart()

	handlerErr, err := trxPreHandleActorPid.RequestFuture(reqMsg, 1000*time.Millisecond).Result() // await result

	if nil != err {
		log.Error("fc ", common.Elapsed(start))
		resp.Errcode = uint32(bottosErr.ErrActorHandleError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrActorHandleError)

		log.Error("trx: %x actor process failed", intTrx.Hash())

		return nil
	}

	log.Error("succ, elapsed time ", common.Elapsed(start))

	if bottosErr.ErrNoError == handlerErr {
		resp.Result = &api.SendTransactionResponse_Result{}
		resp.Result.TrxHash = intTrx.Hash().ToHexString()
		resp.Result.Trx = ConvertIntTrxToApiTrx(intTrx)
		resp.Msg = "trx receive succ"
		resp.Errcode = 0
	} else {
		resp.Result = &api.SendTransactionResponse_Result{}
		resp.Result.TrxHash = intTrx.Hash().ToHexString()
		resp.Result.Trx = ConvertIntTrxToApiTrx(intTrx)
		//resp.Msg = handlerErr.(string)GetCodeString
		//resp.Msg = "to be add detail error description"
		var tempErr bottosErr.ErrCode
		tempErr = handlerErr.(bottosErr.ErrCode)

		resp.Errcode = (uint32)(tempErr)
		resp.Msg = bottosErr.GetCodeString((bottosErr.ErrCode)(resp.Errcode))
	}

	log.Infof("trx: %v %s", resp.Result.TrxHash, resp.Msg)

	return nil
}

//GetTransaction query trx
func (a *ApiService) GetTransaction(ctx context.Context, req *api.GetTransactionRequest, resp *api.GetTransactionResponse) error {
	msgReq := &message.QueryTrxReq{
		TrxHash: common.HexToHash(req.TrxHash),
	}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Errcode = 1
		return nil
	}

	response := res.(*message.QueryTrxResp)
	if response.Trx == nil {
		resp.Errcode = uint32(bottosErr.ErrApiTrxNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiTrxNotFound)
		return nil
	}

	resp.Result = ConvertIntTrxToApiTrx(response.Trx)
	resp.Errcode = uint32(bottosErr.ErrNoError)
	return nil
}

//GetBlock query block
func (a *ApiService) GetBlock(ctx context.Context, req *api.GetBlockRequest, resp *api.GetBlockResponse) error {
	msgReq := &message.QueryBlockReq{
		BlockHash:   common.HexToHash(req.BlockHash),
		BlockNumber: req.BlockNum,
	}
	res, err := chainActorPid.RequestFuture(msgReq, 500*time.Millisecond).Result()
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiBlockNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiBlockNotFound)
		return nil
	}

	response := res.(*message.QueryBlockResp)
	if response.Block == nil {
		resp.Errcode = uint32(bottosErr.ErrApiBlockNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiBlockNotFound)
		return nil
	}

	resp.Result = &api.GetBlockResponse_Result{}
	hash := response.Block.Hash()
	resp.Result.PrevBlockHash = response.Block.GetPrevBlockHash().ToHexString()
	resp.Result.BlockNum = response.Block.GetNumber()
	resp.Result.BlockHash = hash.ToHexString()
	resp.Result.CursorBlockLabel = hash.Label()
	resp.Result.BlockTime = response.Block.GetTimestamp()
	resp.Result.TrxMerkleRoot = response.Block.ComputeMerkleRoot().ToHexString()
	resp.Result.Delegate = string(response.Block.GetDelegate())
	resp.Result.DelegateSign = common.BytesToHex(response.Block.GetDelegateSign())

	resp.Errcode = 0
	return nil
}

var QuerhChainInfoCntReq uint32 = 0
var QuerhChainInfoCntSuc uint32 = 0

//GetInfo query chain info
func (a *ApiService) GetInfo(ctx context.Context, req *api.GetInfoRequest, resp *api.GetInfoResponse) error {

	QuerhChainInfoCntReq++

	log.Error("api actor rcv QueryChainInfo Req, cnt ", QuerhChainInfoCntReq)

	/*

	msgReq := &message.QueryChainInfoReq{}
	start := common.MeasureStart()
	res, err := chainActorPid.RequestFuture(msgReq, 1000*time.Millisecond).Result()
	if err != nil {
		log.Error("failed, elapsed time ",common.Elapsed(start) )
		QuerhChainInfoCntTimeout++
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
		return nil
	}

	log.Error("succed, elapsed time ",common.Elapsed(start) )

	response := res.(*message.QueryChainInfoResp)
	if response.Error != nil {
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
		return nil
	}
	*/

	coreState, err := a.env.RoleIntf.GetChainState()
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiQueryChainInfoError)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiQueryChainInfoError)
		return nil
	}

	resp.Result = &api.GetInfoResponse_Result{}
	resp.Result.HeadBlockNum = coreState.LastBlockNum
	resp.Result.LastConsensusBlockNum = coreState.LastConsensusBlockNum
	resp.Result.HeadBlockHash = coreState.LastBlockHash.ToHexString()
	resp.Result.HeadBlockTime = coreState.LastBlockTime
	resp.Result.HeadBlockDelegate = coreState.CurrentDelegate
	resp.Result.CursorLabel = coreState.LastBlockHash.Label()
	resp.Errcode = 0

	QuerhChainInfoCntSuc++
	log.Error("api actor rcv QueryChainInfo suc, cnt,direct", QuerhChainInfoCntSuc)

	return nil
}

//GetAccount query account info
func (a *ApiService) GetAccount(ctx context.Context, req *api.GetAccountRequest, resp *api.GetAccountResponse) error {
	name := req.AccountName
	account, err := a.env.RoleIntf.GetAccount(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound) + "_" + name + "_1"
		return nil
	}

	balance, err := a.env.RoleIntf.GetBalance(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound) + "_" + name + "_1"
		return nil
	}

	stakedBalance, err := a.env.RoleIntf.GetStakedBalance(name)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound) + "_" + name + "_1"
		return nil
	}

	resp.Result = &api.GetAccountResponse_Result{}
	resp.Result.AccountName = name
	resp.Result.Pubkey = common.BytesToHex(account.PublicKey)
	resp.Result.Balance = balance.Balance.String()
	resp.Result.StakedBalance = stakedBalance.StakedBalance.String()
	resp.Result.UnStakingBalance = stakedBalance.UnstakingBalance.String()
	resp.Result.UnStakingTimestamp = stakedBalance.LastUnstakingTime
	resp.Errcode = 0

	return nil
}

//QueryDBValue query contract object
func (a *ApiService) QueryDBValue(ctx context.Context, req *api.QueryDBValueRequest, resp *api.QueryDBValueResponse) error {
	contract := req.Contract
	object := req.TableName
	key := req.Key
	value, err := a.env.RoleIntf.GetStrValue(contract, object, key)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiObjectNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiObjectNotFound)
		return nil
	}

	resp.Result = &api.QueryDBValueResponse_Result{}
	resp.Result.Contract = contract
	resp.Result.TableName = object
	resp.Result.Key = key
	resp.Result.Value = value
	resp.Errcode = 0

	return nil
}

//GetAbi query contract abi info
func (a *ApiService) GetAbi(ctx context.Context, req *api.GetAbiRequest, resp *api.GetAbiResponse) error {
	contract := req.Contract
	contractInfo, err := a.env.RoleIntf.GetContract(contract)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrApiAccountNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrApiAccountNotFound)
		return nil
	}

	if len(contractInfo.ContractAbi) > 0 {
		resp.Result = string(contractInfo.ContractAbi)
	} else {
		// TODO
		return nil
	}

	return nil
}

//GetTransferCredit query trx credit info
func (a *ApiService) GetTransferCredit(ctx context.Context, req *api.GetTransferCreditRequest, resp *api.GetTransferCreditResponse) error {
	name := req.Name
	spender := req.Spender
	credit, err := a.env.RoleIntf.GetTransferCredit(name, spender)
	if err != nil {
		resp.Errcode = uint32(bottosErr.ErrTransferCreditNotFound)
		resp.Msg = bottosErr.GetCodeString(bottosErr.ErrTransferCreditNotFound)
		return nil
	}

	resp.Result = &api.GetTransferCreditResponse_Result{}
	resp.Result.Name = credit.Name
	resp.Result.Spender = credit.Spender
	resp.Result.Limit = credit.Limit.String()

	return nil
}

func newTransaction(contract string, method string, param []byte) *types.Transaction {
	trx := &types.Transaction{
		Sender:   contract,
		Contract: contract,
		Method:   method,
		Param:    param,
	}

	return trx
}

//GetAllDelegates get all delegates
func (a *ApiService) GetAllDelegates(ctx context.Context, req *api.GetAllDelegatesRequest, resp *api.GetAllDelegatesResponse) error {
	//pubkey, seckey := crypto.GenerateKey()
	//
	//resp.Result.PublicKey = hex.EncodeToString(pubkey)
	//resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//GenerateKeyPair query chain info
func (a *ApiService) GenerateKeyPair(ctx context.Context, req *api.GenerateKeyPairRequest, resp *api.GenerateKeyPairResponse) error {
	pubkey, seckey := crypto.GenerateKey()

	resp.Result.PublicKey = hex.EncodeToString(pubkey)
	resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//CreateAccount get all delegates
func (a *ApiService) CreateAccount(ctx context.Context, req *api.CreateAccountRequest, resp *api.CreateAccountResponse) error {
	//pubkey, seckey := crypto.GenerateKey()
	//
	//resp.Result.PublicKey = hex.EncodeToString(pubkey)
	//resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//CreateWallet get all delegates
func (a *ApiService) CreateWallet(ctx context.Context, req *api.CreateWalletRequest, resp *api.CreateWalletResponse) error {
	//pubkey, seckey := crypto.GenerateKey()
	//
	//resp.Result.PublicKey = hex.EncodeToString(pubkey)
	//resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}
//CreateWalletManual get all delegates
func (a *ApiService) CreateWalletManual(ctx context.Context, req *api.CreateWalletManualRequest, resp *api.CreateWalletManualResponse) error {
	//pubkey, seckey := crypto.GenerateKey()
	//
	//resp.Result.PublicKey = hex.EncodeToString(pubkey)
	//resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//GetKeyPair query chain info
func (a *ApiService) UnlockAccount(ctx context.Context, req *api.UnlockAccountRequest, resp *api.UnlockAccountResponse) error {
	//pubkey, seckey := crypto.GenerateKey()

	//resp.Result.PublicKey = hex.EncodeToString(pubkey)
	//resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//GetKeyPair query chain info
func (a *ApiService) LockAccount(ctx context.Context, req *api.LockAccountRequest, resp *api.LockAccountResponse) error {
	//pubkey, seckey := crypto.GenerateKey()
	//
	//resp.Result.PublicKey = hex.EncodeToString(pubkey)
	//resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//ListWallet get Wallet List
func (a *ApiService) ListWallet(ctx context.Context, req *api.ListWalletRequest, resp *api.ListWalletResponse) error {
	//pubkey, _ := crypto.GenerateKey()

	//resp.Result.PublicKey = hex.EncodeToString(pubkey)
	//resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//GetKeyPair query chain info
func (a *ApiService) GetKeyPair(ctx context.Context, req *api.GetKeyPairRequest, resp *api.GetKeyPairResponse) error {
	pubkey, seckey := crypto.GenerateKey()

	resp.Result.PublicKey = hex.EncodeToString(pubkey)
	resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//SignTransaction query chain info
func (a *ApiService) SignTransaction(ctx context.Context, req *api.SignTransactionRequest, resp *api.SignTransactionResponse) error {
	//pubkey, seckey := crypto.GenerateKey()
	//
	//resp.Result.PublicKey = hex.EncodeToString(pubkey)
	//resp.Result.PrivateKey = hex.EncodeToString(seckey)
	resp.Errcode = 0

	return nil
}

//SignData query chain info
func (a *ApiService) SignData(ctx context.Context, req *api.SignDataRequest, resp *api.SignDataResponse) error {
	resp.Errcode = 0

	return nil
}

//SignHash query chain info
func (a *ApiService) SignHash(ctx context.Context, req *api.SignHashRequest, resp *api.SignHashResponse) error {
	resp.Errcode = 0

	return nil
}

//GetDelegate query chain info
func (a *ApiService) GetDelegate(ctx context.Context, req *api.GetDelegateRequest, resp *api.GetDelegateResponse) error {
	resp.Errcode = 0

	return nil
}

//GetDelegate query chain info
func (a *ApiService) GetPubKey(ctx context.Context, req *api.GetPubKeyRequest, resp *api.GetPubKeyResponse) error {
	resp.Errcode = 0

	return nil
}

//GetDelegate query chain info
func (a *ApiService) GetForecastResBalance(ctx context.Context, req *api.GetForecastResBalanceRequest, resp *api.GetForecastResBalanceResponse) error {
	resp.Errcode = 0

	return nil
}
//GetPeers query peers
func (a *ApiService) GetPeers(ctx context.Context, req *api.GetPeersRequest, resp *api.GetPeersResponse) error {
	resp.Errcode = 0

	return nil
}

//GetAccountBrief query account Brief
func (a *ApiService) GetAccountBrief(ctx context.Context, req *api.GetAccountBriefRequest, resp *api.GetAccountBriefResponse) error {
	resp.Errcode = 0

	return nil
}

//ReviewProposal query account Brief
func (a *ApiService) ReviewProposal(ctx context.Context, req *api.ReviewProposalRequest, resp *api.ReviewProposalResponse) error {
	resp.Errcode = 0

	return nil
}