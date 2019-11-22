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

package main

import (
	"bytes"
	"errors"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
		"fmt"
	"io/ioutil"
		"os"

	"golang.org/x/net/context"

	chain "github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/crypto-go/crypto"
	//"net/http"
	//"strings"
	"math/big"
	"github.com/bottos-project/bottos/common/safemath"
	//"github.com/bitly/go-simplejson"
	TODO "github.com/bottos-project/bottos/restful/handler"
	"github.com/bottos-project/bottos/common/vm"
)



//Transaction trx info
type Transaction struct {
	Version     uint32      `json:"version"`
	CursorNum   uint64      `json:"cursor_num"`
	CursorLabel uint32      `json:"cursor_label"`
	Lifetime    uint64      `json:"lifetime"`
	Sender      string      `json:"sender"`
	Contract    string      `json:"contract"`
	Method      string      `json:"method"`
	Param       interface{} `json:"param"`
	ParamBin    string      `json:"param_bin"`
	SigAlg      uint32      `json:"sig_alg"`
	Signature   string      `json:"signature"`
}

func getBottosAbi() (abi.ABI, error) {
	abistring := `{
	    "types": null,
	    "structs": [
		{
		    "name": "NewAccount",
		    "base": "",
		    "fields": {
			"name": "string",
			"pubkey": "string"
		    }
		},
		{
		    "name": "Transfer",
		    "base": "",
		    "fields": {
			"from": "string",
			"to": "string",
			"value": "uint256",
			"memo" : "string"
		    }
		},
		{
		    "name": "SetDelegate",
		    "base": "",
		    "fields": {
			"name": "string",
			"pubkey": "string",
			"location": "string",
			"description": "string"
		    }
		},
		{
		    "name": "UnsetDelegate",
		    "base": "",
		    "fields": {
			"name": "string"
		    }
		},
		{
		    "name": "GrantCredit",
		    "base": "",
		    "fields": {
			"name": "string",
			"spender": "string",
			"limit": "uint256"
		    }
		},
		{
		    "name": "CancelCredit",
		    "base": "",
		    "fields": {
			"name": "string",
			"spender": "string"
		    }
		},
		{
		    "name": "TransferFrom",
		    "base": "",
		    "fields": {
			"from": "string",
			"to": "string",
			"value": "uint256"
		    }
		},
		{
		    "name": "DeployCode",
		    "base": "",
		    "fields": {
			"contract": "string",
			"vm_type": "uint8",
			"vm_version": "uint8",
			"contract_code": "bytes"
		    }
		},
		{
		    "name": "DeployABI",
		    "base": "",
		    "fields": {
			"contract": "string",
			"contract_abi": "bytes",
			"filetype":"string"
		    }
		},
		{
		    "name": "RegDelegate",
		    "base": "",
		    "fields": {
			"name": "string",
			"pubkey": "string",
			"location": "string",
			"description": "string"
		    }
		},
		{
		    "name": "UnregDelegate",
		    "base": "",
		    "fields": {
			"name": "string"
		    }
		},
		{
		    "name": "VoteDelegate",
		    "base": "",
		    "fields": {
			"voteop": "uint8",
			"voter": "string",
			"delegate": "string"
		    }
		},
		{
		    "name": "Stake",
		    "base": "",
		    "fields": {
			"amount": "uint256",
			"target": "string"
		    }
		},
		{
		    "name": "Unstake",
		    "base": "",
		    "fields": {
			"amount": "uint256",
			"source": "string"
		    }
		},
		{
		    "name": "Claim",
		    "base": "",
		    "fields": {
			"amount": "uint256"
		    }
		},
		{
		    "name": "BlkProdTrans",
		    "base": "",
		    "fields": {
			"actblknum": "uint64"
		    }
		},
		{
		    "name": "SetTransitVote",
		    "base": "",
		    "fields": {
			"name": "string",
			"vote": "uint64"
		    }
		}
	    ],
	    "actions": [
		{
		    "action_name": "newaccount",
		    "type": "NewAccount"
		},
		{
		    "action_name": "transfer",
		    "type": "Transfer"
		},
		{
		    "action_name": "grantcredit",
		    "type": "GrantCredit"
		},
		{
		    "action_name": "cancelcredit",
		    "type": "CancelCredit"
		},
		{
		    "action_name": "transferfrom",
		    "type": "TransferFrom"
		},
		{
		    "action_name": "deploycode",
		    "type": "DeployCode"
		},
		{
		    "action_name": "deployabi",
		    "type": "DeployABI"
		},
		{
		    "action_name": "regdelegate",
		    "type": "RegDelegate"
		},
		{
		    "action_name": "unregdelegate",
		    "type": "UnregDelegate"
		},
		{
		    "action_name": "votedelegate",
		    "type": "VoteDelegate"
		},
		{
		    "action_name": "stake",
		    "type": "Stake"
		},
		{
		    "action_name": "unstake",
		    "type": "Unstake"
		},
		{
		    "action_name": "claim",
		    "type": "Claim"
		},
		{
		    "action_name": "setdelegate",
		    "type": "SetDelegate"
		},
		{
		    "action_name": "settransitvote",
		    "type": "SetTransitVote"
		},
		{
		    "action_name": "unsetdelegate",
		    "type": "UnsetDelegate"
		},
		{
		    "action_name": "blkprodtrans",
		    "type": "BlkProdTrans"
		}
	    ],
	    "tables": null
	}
	`
	Abi, err := abi.ParseAbi([]byte(abistring))
	if err != nil {
		fmt.Println("Parse abistring", abistring, " to abi failed!")
		return abi.ABI{}, err
	}

	return *Abi, nil
}

func (cli *CLI) UnlockWalletOverHttp(http_url string, account string, password string, storepath string) (*chain.UnlockAccountResponse, error) {
	var getinfo *chain.UnlockAccountRequest
	getinfo = &chain.UnlockAccountRequest{AccountName: account, Passwd: password}

	req, _ := json.Marshal(getinfo)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
		return nil, errors.New("Error!")
	}

	var trxrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &trxrespbody)

	if err != nil {
		fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		return nil, errors.New("Error!")
	} else if trxrespbody.Errcode != 0 {
		fmt.Println("Error! ", trxrespbody.Errcode, ":", trxrespbody.Msg)
		return nil, errors.New("Error!")

	} else if trxrespbody.Result == nil {
		fmt.Println("Error! trxrespbody.Result is empty!")
		return nil, errors.New("Error!")
	}

	b, _ := json.Marshal(trxrespbody.Result)
	//cli.jsonPrint(b)
	var RspInfo chain.UnlockAccountResponse
	json.Unmarshal(b, &RspInfo)

	return &RspInfo, nil
}

func (cli *CLI) GetPrivateKeyOverHttp(http_url string, account string) (string, error) {
	var getinfo *chain.GetKeyPairRequest
	getinfo = &chain.GetKeyPairRequest{AccountName: account}

	req, _ := json.Marshal(getinfo)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
		return "", errors.New("Error!")
	}

	var trxrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &trxrespbody)

	if err != nil {
		fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		return "", errors.New("Error!")
	} else if trxrespbody.Errcode != 0 {
		fmt.Println("Error! ", trxrespbody.Errcode, ":", trxrespbody.Msg)
		return "", errors.New("Error!")

	} else if trxrespbody.Result == nil {
		fmt.Println("Error! trxrespbody.Result is empty!")
		return "", errors.New("Error!")
	}

	b, _ := json.Marshal(trxrespbody.Result)
	//cli.jsonPrint(b)
	var RspInfo chain.GetKeyPairResponse
	json.Unmarshal(b, &RspInfo)

	if RspInfo.Result != nil {
		return RspInfo.Result.PrivateKey, nil
	}

	return "", errors.New("Error!")
}

func (cli *CLI) getChainInfo() (*chain.GetInfoResponse_Result, error) {
	chainInfoRsp, err := cli.client.GetInfo(context.TODO(), &chain.GetInfoRequest{})
	if err != nil || chainInfoRsp == nil {
		fmt.Println(err)
		return nil, err
	}

	chainInfo := chainInfoRsp.GetResult()
	return chainInfo, nil
}

func (cli *CLI) GetChainInfoOverHttp(http_url string) (*chain.GetInfoResponse_Result, error) {
	getinfo := &chain.GetInfoRequest{}
	req, _ := json.Marshal(getinfo)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("GET", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
		return nil, errors.New("Error!")
	}

	var trxrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &trxrespbody)

	if err != nil {
		fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		return nil, errors.New("Error!")
	} else if trxrespbody.Errcode != 0 {
		fmt.Println("Error! ", trxrespbody.Errcode, ":", trxrespbody.Msg)
		return nil, errors.New("Error!")

	} else if trxrespbody.Result == nil {
		fmt.Println("Error! trxrespbody.Result is empty!")
		return nil, errors.New("Error!")
	}

	b, _ := json.Marshal(trxrespbody.Result)
	//cli.jsonPrint(b)
	var chainInfo chain.GetInfoResponse_Result
	json.Unmarshal(b, &chainInfo)

	return &chainInfo, nil
}

func (cli *CLI) getBlockInfoOverHttp(http_url string, block_num uint64, block_hash string, choice uint64) (*types.BlockDetail, error) {
	var getinfo *chain.GetBlockRequest
	if choice == 0 {
		getinfo = &chain.GetBlockRequest{BlockNum: block_num}
	} else if choice == 1 {
		getinfo = &chain.GetBlockRequest{BlockHash: block_hash}
	} else {
		getinfo = &chain.GetBlockRequest{BlockNum: 0}
	}

	req, _ := json.Marshal(getinfo)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
		return nil, errors.New("Error!")
	}

	var trxrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &trxrespbody)

	if err != nil {
		fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		return nil, errors.New("Error!")
	} else if trxrespbody.Errcode != 0 {
		fmt.Println("Error! ", trxrespbody.Errcode, ":", trxrespbody.Msg)
		return nil, errors.New("Error!")

	} else if trxrespbody.Result == nil {
		fmt.Println("Error! trxrespbody.Result is empty!")
		return nil, errors.New("Error!")
	}

	b, _ := json.Marshal(trxrespbody.Result)
	var blockInfo types.BlockDetail
	json.Unmarshal(b, &blockInfo)
	//cli.jsonPrint(b)

	return &blockInfo, nil
}

func (cli *CLI) getAccountInfoOverHttp(name string, http_url string, silent ...bool) (*chain.GetAccountResponse_Result, error) {

	getinfo := &chain.GetAccountRequest{AccountName: name}
	req, _ := json.Marshal(getinfo)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		if len(silent) <= 0 {
			fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
		}
		return nil, errors.New("Error!")
	}

	var respbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &respbody)

	if err != nil {
		if len(silent) <= 0 {
			fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		}
		return nil, errors.New("Error!")
	} else if respbody.Errcode != 0 {
		if len(silent) <= 0 {
			fmt.Println("Error! ", respbody.Errcode, ":", respbody.Msg)
		}
		return nil, errors.New("Error!")
	} else if respbody.Result == nil {
		fmt.Println("Error! trxrespbody.Result is empty!")
		return nil, errors.New("Error!")
	}

	b, _ := json.Marshal(respbody.Result)
	//cli.jsonPrint(b)
	var accountInfo chain.GetAccountResponse_Result
	json.Unmarshal(b, &accountInfo)

	return &accountInfo, nil
}

func (cli *CLI) signTrx(trx *chain.Transaction, param []byte, seckey string) (string, error) {
	ctrx := &types.BasicTransaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       param,
		SigAlg:      trx.SigAlg,
	}

	data, err := bpl.Marshal(ctrx)
	if nil != err {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(data)))
	chainId, err := GetChainId()
	h.Write([]byte(hex.EncodeToString(chainId)))
	hashData := h.Sum(nil)
	//seckey, err := GetDefaultKey()
	seckey2, _ := hex.DecodeString(seckey)
	//do not use []byte(seckey) here.
	signdata, err := crypto.Sign(hashData, seckey2)

	return BytesToHex(signdata), err
}

func (cli *CLI) transfer(from, to string, amount big.Int, memo string) {

	infourl := "http://" + ChainAddr + "/v1/account/info"
	account, err := cli.getAccountInfoOverHttp(from, infourl)

	if err != nil || account == nil {
		fmt.Println("Account 'from' does not exist!")
		return
	}

	balance := big.NewInt(0)
	mulval := big.NewInt(100000000)

	balanceResult1, result := balance.SetString(account.Balance, 10)
	if false == result {
		fmt.Println("Error: balance.SetString failed. account: ", account)
		return
	}

	infourl = "http://" + ChainAddr + "/v1/account/brief"
	account, err = cli.getAccountInfoOverHttp(to, infourl)

	if err != nil || account == nil {
		fmt.Println("Account 'to' does not exist!")
		return
	}

	if balanceResult1.Cmp(&amount) < 0 /* < amount */ {
		var mulrestlt *big.Int = big.NewInt(0)
		var modrestlt *big.Int = big.NewInt(0)

		mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult1, mulval)

		if err != nil {
			return
		}

		mulval2 := big.NewInt(100000000)
		modrestlt, err = safemath.U256Mod(modrestlt, balanceResult1, mulval2)
		if err != nil {
			return
		}

		fmt.Printf("Error: User %s has %d.%08d BTO, it is less than your transfer amount!\n", from, mulrestlt, modrestlt)
		return
	}

	type TransferParam struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount string `json:"value"`
		Memo   string `json:"memo"`
	}

	value := big.NewInt(1)
	value2 := big.NewInt(0)

	value2, _ = safemath.U256Mul(value2, &amount, value)

	value2str := value2.String()
	if len(value2str) <= 8 {
		idx := 0

		add_zero_cnt := 9 - len(value2str)
		for idx < add_zero_cnt {
			idx++
			value2str = "0" + value2str //ensure 0.*** can be guaranteed
		}
	}

	/*tp := &TransferParam{
		From:   from,
		To:     to,
		Amount: value2str[0:len(value2str)-8] + "." + value2str[len(value2str)-8:],
		Memo:   memo,
	}*/

	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		return
	}

	mapstruct := make(map[string]interface{})
	abi.Setmapval(mapstruct, "from", from)
	abi.Setmapval(mapstruct, "to", to)
	abi.Setmapval(mapstruct, "value", *value2)
	abi.Setmapval(mapstruct, "memo", memo)

	param, msgPackErr := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "transfer")

	if nil != msgPackErr {
		fmt.Println("msg pack err: ", msgPackErr)
		return
	}

	var realSender string
	if accountType, accountName := common.AnalyzeName(from); accountType == common.NameTypeExContract {
		realSender = accountName
	}
	if realSender == "" {
		realSender = from
	}
	http_url := "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, realSender, "bottos", "transfer", BytesToHex(param))

	if err != nil || ptrx == nil {
		return
	}

	trx := *ptrx

	req, _ := json.Marshal(trx)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", "http://"+ChainAddr+"/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}

	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("Error! ", respbody.Errcode, ":", respbody.Msg)
		return
	}
	newAccountRsp := &respbody

	/*fmt.Printf("\nPush transaction done.\n")
	fmt.Printf("    From: %v\n", from)
	fmt.Printf("    To: %v\n", to)
	fmt.Println("    Amount:", value2str[0:len(value2str)-8]+"."+value2str[len(value2str)-8:])
	fmt.Printf("    Memo: %v\n", memo)
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       tp,
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)*/
	fmt.Printf("\nTrxHash: %v\n", newAccountRsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}

func (cli *CLI) newmsignaccount(name string, authority string, threshold uint32, referrer string) {


    //0check-name
	infourl := "http://" + ChainAddr + "/v1/account/brief"
	account, _ := cli.getAccountInfoOverHttp(name, infourl, true)

	if account != nil {
		fmt.Println("\nError: The multisign account has been already registered.\n")
		return
	}

	//0check-authority
	beginIndex := strings.Index(authority, "[")
	endIndex := strings.LastIndex(authority, "]")
	if beginIndex < 0 || endIndex < 0 {
		fmt.Println("\nError: Invalid multisign account authority. Ensure your multisign account authority is valid.\n")
		return
	}
	var msignAccountAuthority []role.MsignAccountAuthority
	if err := json.Unmarshal([]byte(authority[beginIndex:endIndex+1]), &msignAccountAuthority); err != nil {
		fmt.Println("\nError: Invalid multisign account authority. Ensure your multisign account authority is valid.\n")
		return
	}
	for _, mauthority := range msignAccountAuthority {
		//0check-authority-account
		infourl := "http://" + ChainAddr + "/v1/account/brief"
		authorityAccount, err := cli.getAccountInfoOverHttp(mauthority.AuthorAccount, infourl)
		if err != nil || authorityAccount == nil {
			fmt.Println("\nError: Account", mauthority.AuthorAccount,"in 'authority' has not been found!\n")
			return
		}
		//0check-authority-weight
		if mauthority.Weight > threshold {
			fmt.Println("\nError:  Forbid multisign account authority's weight be larger than threshold\n")
			return
		}
	}

	// 1, new account trx
	type NewMultiAccountParam struct {
		Account   string `json:"account"`
		Authority string `json:"authority"`
		Threshold uint32 `json:"threshold"`
	}


	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		fmt.Println("\nError: get abi err: ", abierr)
		return
	}

	mapstruct := make(map[string]interface{})

	abi.Setmapval(mapstruct, "account", name)
	abi.Setmapval(mapstruct, "authority", authority)
	abi.Setmapval(mapstruct, "threshold", threshold)

	param, err := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "newmsignaccount")

	Sender := "bottos"
	if len(referrer) > 0 {
		Sender = referrer
	}

	http_url := "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, Sender, "bottos", "newmsignaccount", BytesToHex(param))

	if err != nil || ptrx == nil {
		return
	}

	trx := *ptrx

	req, _ := json.Marshal(trx)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", "http://"+ChainAddr+"/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError: BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}

	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("Error! ", respbody.Errcode, ":", respbody.Msg)
		return
	}
	rsp := &respbody


	fmt.Printf("\nTrxHash: %v\n", rsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}

func (cli *CLI) pushmsignproposal(proposal, account, transfer, proposer string) {



	type PushmsignproposalParam struct {
		Proposal   string `json:"proposal"`
		Account     string `json:"account"`
		Transfer string `json:"transfer"`
		Proposer   string `json:"proposer"`
	}

	//0check-account/proposer
	infourl := "http://" + ChainAddr + "/v1/account/brief"
	accountName, _ := cli.getAccountInfoOverHttp(account, infourl, true)
	proposerName, _ := cli.getAccountInfoOverHttp(proposer, infourl, true)

	if accountName == nil || proposerName == nil{
		fmt.Println("\nError: The multisign account or proposer has not been found.\n")
		return
	}

	//0check-transfer
	transferParam := &role.MsignTransferParam{}
	if err := json.Unmarshal([]byte(transfer), transferParam); err != nil {
		fmt.Println("\nError: Wrong transfer infomation.\n")
		return
	}
	//0check-transfer-from/to
	fromName, _ := cli.getAccountInfoOverHttp(transferParam.From, infourl, true)
	toName, _ := cli.getAccountInfoOverHttp(transferParam.To, infourl, true)

	if fromName == nil || toName == nil{
		fmt.Println("\nError: The multisign proposal transfer from or to has not been found.\n ")
		return
	}

	//0check-transfer-amount
	_, transfer_value_str, err := transferStrinToBigInt(transferParam.Amount)
	if err != nil {
		fmt.Println("\nError: Invalid multisign proposal transfer amount. Ensure your multisign proposal transfer amount is valid.\n")
		return
	}

	transferParam.Amount = transfer_value_str
	transfer_bytes, err := json.Marshal(transferParam)
	if err != nil {
		fmt.Println("\nError: Wrong transfer infomation.\n")
		return
	}
	transfer = string(transfer_bytes)


	//1ABI
	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		fmt.Println("\nError: get abi err: ", abierr)
		return
	}

	mapstruct := make(map[string]interface{})
	abi.Setmapval(mapstruct, "proposal", proposal)
	abi.Setmapval(mapstruct, "account", account)
	abi.Setmapval(mapstruct, "transfer", transfer)
	abi.Setmapval(mapstruct, "proposer", proposer)

	param, msgPackErr := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "pushmsignproposal")

	if nil != msgPackErr {
		fmt.Println("\nError: msg pack err: ", msgPackErr)
		return
	}

	var realSender string
	if accountType, accountName := common.AnalyzeName(proposer); accountType == common.NameTypeExContract {
		realSender = accountName
	}
	if realSender == "" {
		realSender = proposer
	}
	http_url := "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, realSender, "bottos", "pushmsignproposal", BytesToHex(param))

	if err != nil || ptrx == nil {
		return
	}

	trx := *ptrx

	req, _ := json.Marshal(trx)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", "http://"+ChainAddr+"/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError:  BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}

	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("\nError! ", respbody.Errcode, ":", respbody.Msg)
		return
	}
	newAccountRsp := &respbody


	fmt.Printf("\nTrxHash: %v\n", newAccountRsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}

func (cli *CLI) reviewmsignproposal(proposer,proposal string) {


	//0check-proposer
	infourl := "http://" + ChainAddr + "/v1/account/brief"
	proposerName, _ := cli.getAccountInfoOverHttp(proposer, infourl, true)

	if proposerName == nil{
		fmt.Println("\nError: The multisign account or proposer has not been found.")
		return
	}

	//0check-proposal
	http_url := "http://" + ChainAddr + "/v1/proposal/review"
	getproposal := &chain.ReviewProposalRequest{ProposalName: proposal,Proposer:proposer}
	req, _ := json.Marshal(getproposal)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError: . httpRspBody: ", httpRspBody, ", err: ", err)
		return
	}

	var proposalrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &proposalrespbody)

	if err != nil {
		fmt.Println("\nError:  Unmarshal to proposal failed: ", err, "| body is: ", string(httpRspBody), ". proposalresp:")
		return
	}


	if proposalrespbody.Result == nil {
		fmt.Printf("\nError: The multisign proposal has not been found.\n\n")
		return
	}



	b, _ := json.Marshal(proposalrespbody.Result)
	var reviewProposal chain.ReviewProposalResponse_Result
	err = json.Unmarshal(b, &reviewProposal)

	if err != nil {
		fmt.Println("Error! Unmarshal to proposal failed: ", err, "| body is: ", string(httpRspBody))
		return
	}

	fmt.Printf("\n    ProposalName: %s\n\n", reviewProposal.ProposalName)
	fmt.Printf("    MsignAccountName: %s\n\n", reviewProposal.MsignAccountName)
	fmt.Printf("    AuthorList: %s\n\n", reviewProposal.AuthorList)
	fmt.Printf("    Available: %t\n\n", reviewProposal.Available)
	fmt.Printf("    PackedTransaction: %s\n\n", reviewProposal.PackedTransaction)
	fmt.Printf("    Transaction: [%s]\n\n", reviewProposal.Transaction)
	fmt.Printf("    Time: %d\n\n", reviewProposal.Time)

}
func (cli *CLI) approvemsignproposal(proposal, account, proposer string) {



	type ApproveMsignProposalParam struct {
		Proposal   string `json:"proposal"`
		Account     string `json:"account"`
		Proposer   string `json:"proposer"`
	}

	//0check-account/proposer
	infourl := "http://" + ChainAddr + "/v1/account/brief"
	accountName, _ := cli.getAccountInfoOverHttp(account, infourl, true)
	proposerName, _ := cli.getAccountInfoOverHttp(proposer, infourl, true)

	if accountName == nil || proposerName == nil{
		fmt.Println("\nError: The multisign account or proposer has not been found.\n")
		return
	}

	//0check-proposal
	http_url := "http://" + ChainAddr + "/v1/proposal/review"
	getproposal := &chain.ReviewProposalRequest{ProposalName: proposal,Proposer:proposer}
	req, _ := json.Marshal(getproposal)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError:  httpRspBody: ", httpRspBody, ", err: ", err,"\n")
		return
	}

	var proposalrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &proposalrespbody)

	if err != nil {
		fmt.Println("\nError:  Unmarshal to proposal failed: ", err, "| body is: ", string(httpRspBody), ". proposalresp:\n")
		return
	}


	if proposalrespbody.Result == nil {
		fmt.Printf("\nError: The multisign proposal has not been found.\n\n")
		return
	}

	//1ABI
	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		fmt.Println("\nError: get abi err: ", abierr)
		return
	}

	mapstruct := make(map[string]interface{})
	abi.Setmapval(mapstruct, "proposal", proposal)
	abi.Setmapval(mapstruct, "account", account)
	abi.Setmapval(mapstruct, "proposer", proposer)

	param, msgPackErr := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "approvemsignproposal")

	if nil != msgPackErr {
		fmt.Println("msg pack err: ", msgPackErr)
		return
	}

	var realSender string
	if accountType, accountName := common.AnalyzeName(account); accountType == common.NameTypeExContract {
		realSender = accountName
	}
	if realSender == "" {
		realSender = account
	}
	http_url = "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, realSender, "bottos", "approvemsignproposal", BytesToHex(param))

	if err != nil || ptrx == nil {
		return
	}

	trx := *ptrx
	req, _ = json.Marshal(trx)
	req_new = bytes.NewBuffer([]byte(req))
	httpRspBody, err = send_httpreq("POST", "http://"+ChainAddr+"/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError: BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}

	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("\nError:  ", respbody.Errcode, ":", respbody.Msg)
		return
	}
	newAccountRsp := &respbody


	fmt.Printf("\nTrxHash: %v\n", newAccountRsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}
func (cli *CLI) unapprovemsignproposal(proposal, account, proposer string) {



	type UnapproveMsignProposalParam struct {
		Proposal   string `json:"proposal"`
		Account     string `json:"account"`
		Proposer   string `json:"proposer"`
	}


	//0check-account/proposer
	infourl := "http://" + ChainAddr + "/v1/account/brief"
	accountName, _ := cli.getAccountInfoOverHttp(account, infourl, true)
	proposerName, _ := cli.getAccountInfoOverHttp(proposer, infourl, true)

	if accountName == nil || proposerName == nil{
		fmt.Println("\nError: The multisign account or proposer has not been found.\n")
		return
	}

	//0check-proposal
	http_url := "http://" + ChainAddr + "/v1/proposal/review"
	getproposal := &chain.ReviewProposalRequest{ProposalName: proposal,Proposer:proposer}
	req, _ := json.Marshal(getproposal)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError. httpRspBody: ", httpRspBody, ", err: ", err,"\n")
		return
	}

	var proposalrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &proposalrespbody)

	if err != nil {
		fmt.Println("\nError: Unmarshal to proposal failed: ", err, "| body is: ", string(httpRspBody), ". proposalresp:\n")
		return
	}


	if proposalrespbody.Result == nil {
		fmt.Printf("\nError: The multisign proposal has not been found.\n\n")
		return
	}

	//1ABI
	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		fmt.Println("\nError: get abi err: ", abierr)
		return
	}

	mapstruct := make(map[string]interface{})
	abi.Setmapval(mapstruct, "proposal", proposal)
	abi.Setmapval(mapstruct, "account", account)
	abi.Setmapval(mapstruct, "proposer", proposer)

	param, msgPackErr := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "unapprovemsign")

	if nil != msgPackErr {
		fmt.Println("\nError: msg pack err: ", msgPackErr)
		return
	}

	var realSender string
	if accountType, accountName := common.AnalyzeName(proposer); accountType == common.NameTypeExContract {
		realSender = accountName
	}
	if realSender == "" {
		realSender = proposer
	}
	http_url = "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, realSender, "bottos", "unapprovemsign", BytesToHex(param))

	if err != nil || ptrx == nil {
		return
	}

	trx := *ptrx
	req, _ = json.Marshal(trx)
	req_new = bytes.NewBuffer([]byte(req))
	httpRspBody, err = send_httpreq("POST", "http://"+ChainAddr+"/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError: BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}

	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("\nError:  ", respbody.Errcode, ":", respbody.Msg)
		return
	}
	newAccountRsp := &respbody


	fmt.Printf("\nTrxHash: %v\n", newAccountRsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}
func (cli *CLI) execmsignproposal(proposal, proposer string) {



	type ExecMsignProposalParam struct {
		Proposal   string `json:"proposal"`
		Proposer   string `json:"proposer"`
	}
	//0check-proposer
	infourl := "http://" + ChainAddr + "/v1/account/brief"
	proposerName, _ := cli.getAccountInfoOverHttp(proposer, infourl, true)

	if proposerName == nil{
		fmt.Println("\nError: The multisign account or proposer has not been found.")
		return
	}

	//0check-proposal
	http_url := "http://" + ChainAddr + "/v1/proposal/review"
	getproposal := &chain.ReviewProposalRequest{ProposalName: proposal,Proposer:proposer}
	req, _ := json.Marshal(getproposal)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError:  httpRspBody: ", httpRspBody, ", err: ", err)
		return
	}

	var proposalrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &proposalrespbody)

	if err != nil {
		fmt.Println("\nError:  Unmarshal to proposal failed: ", err, "| body is: ", string(httpRspBody), ". proposalresp:")
		return
	}


	if proposalrespbody.Result == nil {
		fmt.Printf("\nError: The multisign proposal has not been found.\n\n")
		return
	}

	//1-ABI
	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		fmt.Println("\nError: get abi err: ", abierr)
		return
	}

	mapstruct := make(map[string]interface{})
	abi.Setmapval(mapstruct, "proposal", proposal)
	abi.Setmapval(mapstruct, "proposer", proposer)

	param, msgPackErr := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "execmsignproposal")

	if nil != msgPackErr {
		fmt.Println("\nError: msg pack err: ", msgPackErr)
		return
	}

	var realSender string
	if accountType, accountName := common.AnalyzeName(proposer); accountType == common.NameTypeExContract {
		realSender = accountName
	}
	if realSender == "" {
		realSender = proposer
	}
	http_url = "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, realSender, "bottos", "execmsignproposal", BytesToHex(param))

	if err != nil || ptrx == nil {
		return
	}

	trx := *ptrx

	req, _ = json.Marshal(trx)
	req_new = bytes.NewBuffer([]byte(req))
	httpRspBody, err = send_httpreq("POST", "http://"+ChainAddr+"/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError: BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}

	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("\nError:  ", respbody.Errcode, ":", respbody.Msg)
		return
	}
	newAccountRsp := &respbody


	fmt.Printf("\nTrxHash: %v\n", newAccountRsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}
func (cli *CLI) cancelmsignproposal(proposal, proposer string) {



	type CancelMsignProposalParam struct {
		Proposal   string `json:"proposal"`
		Proposer   string `json:"proposer"`
	}


	//0check-proposer
	infourl := "http://" + ChainAddr + "/v1/account/brief"
	proposerName, _ := cli.getAccountInfoOverHttp(proposer, infourl, true)

	if proposerName == nil{
		fmt.Println("\nError: The multisign account or proposer has not been found.")
		return
	}

	//0check-proposal
	http_url := "http://" + ChainAddr + "/v1/proposal/review"
	getproposal := &chain.ReviewProposalRequest{ProposalName: proposal,Proposer:proposer}
	req, _ := json.Marshal(getproposal)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError:  httpRspBody: ", httpRspBody, ", err: ", err)
		return
	}

	var proposalrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &proposalrespbody)

	if err != nil {
		fmt.Println("\nError:  Unmarshal to proposal failed: ", err, "| body is: ", string(httpRspBody), ". proposalresp:")
		return
	}


	if proposalrespbody.Result == nil {
		fmt.Printf("\nError: The multisign proposal has not been found.\n\n")
		return
	}

	//1-ABI
	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		fmt.Println("\nError: get abi err: ", abierr)
		return
	}

	mapstruct := make(map[string]interface{})
	abi.Setmapval(mapstruct, "proposal", proposal)
	abi.Setmapval(mapstruct, "proposer", proposer)

	param, msgPackErr := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "cancelmsignproposal")

	if nil != msgPackErr {
		fmt.Println("\nError: msg pack err: ", msgPackErr)
		return
	}

	var realSender string
	if accountType, accountName := common.AnalyzeName(proposer); accountType == common.NameTypeExContract {
		realSender = accountName
	}
	if realSender == "" {
		realSender = proposer
	}
	http_url = "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, realSender, "bottos", "cancelmsignproposal", BytesToHex(param))

	if err != nil || ptrx == nil {
		return
	}

	trx := *ptrx

	req, _ = json.Marshal(trx)
	req_new = bytes.NewBuffer([]byte(req))
	httpRspBody, err = send_httpreq("POST", "http://"+ChainAddr+"/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("\nError: BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}

	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("\nError: ", respbody.Errcode, ":", respbody.Msg)
		return
	}
	newAccountRsp := &respbody

	fmt.Printf("\nTrxHash: %v\n", newAccountRsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}

func (cli *CLI) jsonPrint(data []byte) {
	var out bytes.Buffer
	json.Indent(&out, data, "", "    ")

	fmt.Println(string(out.Bytes()))
}

func IsContractExist(contractname string) bool {
	var err error
	if contractname == config.BOTTOS_CONTRACT_NAME {
		return true
	}

	httpurl_contractcode := "http://" + ChainAddr + "/v1/contract/code"
	getcontract := &GetContractCodeAbi{Contract: contractname}
	req, _ := json.Marshal(getcontract)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", httpurl_contractcode, req_new)

	if err != nil || httpRspBody == nil {
		return false
	}
	var trxrespbody TODO.ResponseStruct

	err = json.Unmarshal(httpRspBody, &trxrespbody)
	if err != nil || trxrespbody.Result == nil {
		return false
	}

	return true
}

//getAbibyContractName function
func getAbibyContractName(contractname string) (abi.ABI, error) {
	/*NodeIp := "127.0.0.1"
	addr := "http://" + NodeIp + ":8080/rpc"
	params := `service=bottos&method=Chain.GetAbi&request={
			"contract":"%s"}`
	s := fmt.Sprintf(params, contractname)
	respBody, err := http.Post(addr, "application/x-www-form-urlencoded", strings.NewReader(s))
	
	if err != nil {
		fmt.Println(err)
		return abi.ABI{}, err
	}
	

	defer respBody.Body.Close()
	body, err := ioutil.ReadAll(respBody.Body)
	if err != nil {
		fmt.Println(err)
		return abi.ABI{}, err
	}

	jss, _ := simplejson.NewJson([]byte(body))
	abistring = jss.Get("result").MustString()
	if len(abistring) <= 0 {
		return abi.ABI{}, errors.New("len(abistring) <= 0")
	}
	
	*/

	http_url := "http://" + ChainAddr + "/v1/contract/abi"
	Abi, err := GetAbiOverHttp(http_url, contractname)
	if err != nil {
		return abi.ABI{}, errors.New("GetAbiOverHttp failed")
	}

	return Abi, nil
}

//getAbiFieldsByAbiEx function
func getAbiFieldsByAbiEx(contractname string, method string, abi abi.ABI, subStructName string) *abi.FeildMap {
	for _, subaction := range abi.Actions {
		if subaction.ActionName != method {
			continue
		}
		structname := subaction.Type

		for _, substruct := range abi.Structs {
			if subStructName != "" {
				if substruct.Name != subStructName {
					continue
				}
			} else if structname != substruct.Name {
				continue
	}

	Abi, err := abi.ParseAbi([]byte(abistring))
	if err != nil {
		fmt.Println("Parse abistring", abistring, " to abi failed!")
		return abi.ABI{}, err
	}

	return *Abi, nil
}

func (cli *CLI) newaccount(name string, pubkey string, referrer string) {

	var err error

	infourl := "http://" + ChainAddr + "/v1/account/brief"
	account, _ := cli.getAccountInfoOverHttp(name, infourl, true)

	if account != nil {
		fmt.Println("The account has been already registered.")
		return
	}

	if len(pubkey) != PUBKEY_LEN {
		fmt.Println("\nNewaccount error: public key len is invalid! Public key sample: 0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f")
		return
	}
	
	if err != nil {
		fmt.Println("GetInfo error: ", err)
		return
	}

	// 1, new account trx
	type NewAccountParam struct {
		Name   string `json:"name"`
		Pubkey string `json:"pubkey"`
	}
	nps := &NewAccountParam{
		Name:   name,
		Pubkey: pubkey,
	}
	
        Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
           return
        }
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "name", name)
        abi.Setmapval(mapstruct, "pubkey", pubkey)
        
	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "newaccount")
	Sender := "bottos"
	if len(referrer) > 0 {
		Sender = referrer
	}

	http_url := "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, Sender, "bottos", "newaccount", BytesToHex(param))

	if err != nil || ptrx == nil {
		return
	}

	trx := *ptrx
	
	req, _ := json.Marshal(trx)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", "http://" + ChainAddr + "/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("Error! ", respbody.Errcode, ":", respbody.Msg)
		return 
	}
	rsp := &respbody
	
	/*fmt.Printf("\nPush transaction done for creating account %v.\n", name)
	fmt.Printf("Trx: \n")
	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       nps,
		ParamBin:    common.BytesToHex(param),
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)*/
	fmt.Printf("\nTrxHash: %v\n", rsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
	fmt.Printf("Please create wallet for your new account.\n")
}

func (cli *CLI) getaccount(name string) {
	//accountRsp, err := cli.client.GetAccount(context.TODO(), &chain.GetAccountRequest{AccountName: name})

	infourl := "http://" + ChainAddr + "/v1/account/info"
	account, err := cli.getAccountInfoOverHttp(name, infourl)

	if err != nil || account == nil {
		return
	}

	/*if accountRsp.Errcode == 10204 {
		fmt.Printf("Account: %s Not Exist\n", name)
		return
	}

	account := accountRsp.GetResult()
	*/
	balance := big.NewInt(0)
	mulval := big.NewInt(100000000)

	balanceResult, result := balance.SetString(account.Balance, 10)
	if false == result {
		fmt.Println("Error: balance.SetString Balance failed. account: ", account)
		return
	}

	var mulrestlt *big.Int = big.NewInt(0)
	var modrestlt *big.Int = big.NewInt(0)

	mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult, mulval)

	if err != nil {
		return
	}

	mulval2 := big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, balanceResult, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("\n    Account: %s\n", account.AccountName)
	fmt.Printf("    Authority: %v\n", account.Authority)
	fmt.Printf("    Threshold: %d\n", account.Threshold)
	fmt.Printf("    Balance: %d.%08d BTO\n", mulrestlt, modrestlt)
	fmt.Printf("    Pubkey: %s\n\n", account.Pubkey)

	balanceResult, result = balance.SetString(account.StakedBalance, 10)
	if false == result {
		fmt.Println("Error: balance.SetString StakedBalance failed. account: ", account)
		return
	}

	mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult, mulval)

	if err != nil {
		return
	}

	mulval2 = big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, balanceResult, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("    StakedBalance: %d.%08d BTO\n", mulrestlt, modrestlt)

	balanceResult, result = balance.SetString(account.UnStakingBalance, 10)
	if false == result {
		fmt.Println("Error: balance.SetString UnStakingBalance failed. account: ", account)
		return
	}

	mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult, mulval)

	if err != nil {
		return
	}

	mulval2 = big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, balanceResult, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("    UnStakingBalance: %d.%08d BTO\n", mulrestlt, modrestlt)

	balanceResult, result = balance.SetString(account.StakedSpaceBalance, 10)
	if false == result {
		fmt.Println("Error: balance.SetString StakedSpaceBalance failed. account: ", account)
		return
	}

	mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult, mulval)

	if err != nil {
		return
	}

	mulval2 = big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, balanceResult, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("    StakedSpaceBalance: %d.%08d BTO\n", mulrestlt, modrestlt)

	balanceResult, result = balance.SetString(account.StakedTimeBalance, 10)
	if false == result {
		fmt.Println("Error: balance.SetString StakedTimeBalance failed. account: ", account)
		return
	}

	mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult, mulval)

	if err != nil {
		return
	}

	mulval2 = big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, balanceResult, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("    StakedTimeBalance: %d.%08d BTO\n", mulrestlt, modrestlt)

	fmt.Printf("    UnStakingTimestamp: %d\n\n", account.UnStakingTimestamp)

	if account.Resource == nil {
		fmt.Printf("    Resource: N/A\n\n")
	} else {
		a, _ := json.MarshalIndent(&account.Resource, "     ", "\t")

		fmt.Printf("    Resource: %v\n\n", string(a))
	}

	balanceResult, result = balance.SetString(account.UnClaimedBlockReward, 10)
	if false == result {
		fmt.Println("Error: balance.SetString UnClaimedBlockReward failed. account: ", account)
		return
	}
	UnClaimedBlockReward := big.NewInt(0).Add(balanceResult, big.NewInt(0))

	mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult, mulval)

	if err != nil {
		return
	}

	mulval2 = big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, balanceResult, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("    UnClaimedBlockReward: %d.%08d BTO\n", mulrestlt, modrestlt)

	balanceResult, result = balance.SetString(account.UnClaimedVoteReward, 10)
	UnClaimedVoteReward := big.NewInt(0).Add(balanceResult, big.NewInt(0))
	if false == result {
		fmt.Println("Error: balance.SetString UnClaimedVoteReward failed. account: ", account)
		return
	}

	mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult, mulval)

	if err != nil {
		return
	}

	mulval2 = big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, balanceResult, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("    UnClaimedVoteReward: %d.%08d BTO\n", mulrestlt, modrestlt)

	UnClaimedTotalReward := big.NewInt(0).Add(UnClaimedBlockReward, UnClaimedVoteReward)
	mulrestlt, err = safemath.U256Div(mulrestlt, UnClaimedTotalReward, mulval)

	if err != nil {
		return
	}

	mulval2 = big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, UnClaimedTotalReward, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("    UnClaimedTotalReward: %d.%08d BTO\n\n", mulrestlt, modrestlt)

	if account.Vote == nil {
		fmt.Printf("    Vote: N/A\n\n")
	} else {
		fmt.Printf("    Vote: %v\n\n", account.Vote)
	}

	if len(account.DeployContractList) <= 0 {
		fmt.Printf("    Contracts: N/A\n\n")
	} else {
		fmt.Printf("    Contracts: %s\n\n", account.DeployContractList)
	}

}

func sendTransaction(trx chain.Transaction) ([]byte, error) {
	http_url := "http://" + ChainAddr + "/v1/transaction/send"
	req, _ := json.Marshal(trx)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil {
		fmt.Println("BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return nil, err
	}
	if httpRspBody == nil {
		fmt.Println("BcliSendTransaction Error:httpRspBody is null")
		return nil, errors.New("deploy send contract resp is null")
	}
	return httpRspBody, nil
}

func getTransactionResp(httpRspBody []byte) (*chain.SendTransactionResponse, error) {
	var respbody *chain.SendTransactionResponse
	if err := json.Unmarshal(httpRspBody, &respbody); err != nil {
		return nil, err
	}

	if respbody.Errcode != 0 {
		fmt.Println("Deploy code error! ", respbody.Errcode, ":", respbody.Msg)
		return respbody, errors.New(respbody.Msg)
	}
	return respbody, nil
}

func readContractFile(contractPath string) ([]byte, error) {
	_, err := ioutil.ReadFile(contractPath)

	if err != nil {
		fmt.Printf("Open %s error: %v", contractPath, err)
		return nil, err
	}

	contractFile, err := os.Open(contractPath)
	defer contractFile.Close()

	if err != nil {
		fmt.Printf("Open %s error: %v", contractPath, err)
		return nil, err
	}

	contractFileInfo, err := contractFile.Stat()
	if err != nil {
		fmt.Printf("Open %s error: %v", contractPath, err)
		return nil, err
	}

	contractCode := make([]byte, contractFileInfo.Size())
	if _, err := contractFile.Read(contractCode); err != nil {
		fmt.Printf("Read %s error: %v", contractPath, err)
		return nil, err
	}

	return contractCode, nil
}

func (cli *CLI) deploycontract(name string, codePath, abiPath string, user string, fileTypeInput string) {
	contractCode, err := readContractFile(codePath)
	if err != nil {
		return
	}

	contractAbi, err := readContractFile(abiPath)
	if err != nil {
		return
	}

	//get file type(wasm or js)
	var fileType vm.VmType
	vmType, err := getCodeFileType(fileTypeInput)
	if err != nil {
		return
	}
	fileType = vmType

	Abi, err := getAbibyContractName("bottos")
	if err != nil {
		return
	}

	//Marshal contract
	mapStruct := buildContractMapStruct(contractCode, contractAbi, name, fileType)
	param, _ := abi.MarshalAbiEx(mapStruct, &Abi, "bottos", "deploycontract")

	//sign transaction
	ptrx, err := signTransaction(cli, user, param)
	if err != nil {
		return
	}

	//send transaction
	var deployContractRsp *chain.SendTransactionResponse

	respBodyByte, err := sendTransaction(*ptrx)
	if err != nil {
		return
	}
	respBody, err := getTransactionResp(respBodyByte)
	if err != nil || respBody == nil {
		return
	}

	deployContractRsp = respBody

	//show resp
	/*fmt.Printf("\nPush transaction done for deploying contract %v.\n", name)
	fmt.Printf("Trx: \n")

	deployContractInfo := showContractInfo(name, fileType, contractCode, *ptrx)
	cli.jsonPrint(deployContractInfo)*/
	fmt.Printf("\nTrxHash: %v\n", deployContractRsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}

func signTransaction(cli *CLI, user string, param []byte) (*chain.Transaction, error) {
	http_url := "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, user, "bottos", "deploycontract", BytesToHex(param))
	if err != nil || ptrx == nil {
		fmt.Println("Deploy contract error! May be your wallet has not been created ok unlocked?")
		return nil, err
	}

	return ptrx, nil
}

func showContractInfo(name string, fileType vm.VmType, ContractCodeVal []byte, trx chain.Transaction) []byte {
	type PrintDeployCodeParam struct {
		Name         string `json:"name"`
		VMType       byte   `json:"vm_type"`
		VMVersion    byte   `json:"vm_version"`
		ContractCode string `json:"contract_code"`
	}
	pdcp := &PrintDeployCodeParam{}
	pdcp.Name = name
	pdcp.VMType = byte(fileType)
	pdcp.VMVersion = 1

	//decide the length of show val
	codeLength := len(ContractCodeVal)
	paramLength := len([]byte(trx.Param))
	if codeLength > 100 {
		codeLength = 100
	}
	if paramLength > 200 {
		paramLength = 200
	}
	codeHex := BytesToHex(ContractCodeVal[0:codeLength])
	pdcp.ContractCode = codeHex + "..."
	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       pdcp,
		ParamBin:    string([]byte(trx.Param)[0:paramLength]) + "...",
		//ParamBin: trx.Param,
		SigAlg:    trx.SigAlg,
		Signature: trx.Signature,
	}
	b, _ := json.Marshal(printTrx)
	return b
}

func buildContractMapStruct(contractCodeVal, contractAbiVal []byte, name string, fileType vm.VmType) map[string]interface{} {
	mapstruct := make(map[string]interface{})

	abi.Setmapval(mapstruct, "contract", name)
	abi.Setmapval(mapstruct, "vm_type", uint8(fileType))
	abi.Setmapval(mapstruct, "vm_version", uint8(0))
	abi.Setmapval(mapstruct, "contract_code", contractCodeVal)
	abi.Setmapval(mapstruct, "contract_abi", contractAbiVal)

	return mapstruct
}

func getCodeFileType(fileTypeInput string) (vm.VmType, error) {
	if fileTypeInput == "wasm" {
		return vm.VmTypeWasm, nil
	} else if fileTypeInput == "js" {
		return vm.VmTypeJS, nil
	} else {
		fmt.Println("file type should be wasm or js.")
		return vm.VmTypeUnkonw, errors.New("file type should be wasm or js")
	}
}

func (cli *CLI) deploycode(name string, path string, user string, fileTypeInput string) {
	var err error
	_, err = ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Open wasm file error: ", err)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Open wasm file error: ", err)
		return
	}

	var fileType vm.VmType

	if fileTypeInput == "wasm" {
		fileType = vm.VmTypeWasm
	} else if fileTypeInput == "js" {
		fileType = vm.VmTypeJS
	} else {
		fmt.Println("file type should be wasm or js.")
		return
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("Open wasm file error: ", err)
		return
	}

	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		return
	}

	var ContractCodeVal []byte
	ContractCodeVal = make([]byte, fi.Size())
	f.Read(ContractCodeVal)
	mapstruct := make(map[string]interface{})

	abi.Setmapval(mapstruct, "contract", name)
	abi.Setmapval(mapstruct, "vm_type", uint8(fileType))
	abi.Setmapval(mapstruct, "vm_version", uint8(1))

	abi.Setmapval(mapstruct, "contract_code", ContractCodeVal)
	//fmt.Printf("contract_code: %x", ContractCodeVal)
	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "deploycode")

	http_url := "http://" + ChainAddrWallet + "/v1/wallet/signtransaction"
	ptrx, err := cli.BcliSignTrxOverHttp(http_url, user, "bottos", "deploycode", BytesToHex(param))

	if err != nil || ptrx == nil {
		fmt.Println("Deploy code error! May be your wallet has not been created ok unlocked?")
		return
	}

	trx := *ptrx

	var deployCodeRsp *chain.SendTransactionResponse

	http_method := "restful"

	if http_method == "grpc" {
		deployCodeRsp, err = cli.client.SendTransaction(context.TODO(), ptrx)
		if err != nil {
			fmt.Println(err)
			return
		}

		if deployCodeRsp.Errcode != 0 {
			fmt.Printf("Deploy contract error:\n")
			fmt.Printf("    %v\n", deployCodeRsp.Msg)
			return
		}
	} else {
		http_url := "http://" + ChainAddr + "/v1/transaction/send"
		req, _ := json.Marshal(trx)
		req_new := bytes.NewBuffer([]byte(req))
		httpRspBody, err := send_httpreq("POST", http_url, req_new)
		if err != nil || httpRspBody == nil {
			fmt.Println("BcliSendTransaction Error:", err, ", httpRspBody: ", httpRspBody)
			return
		}
		var respbody chain.SendTransactionResponse
		json.Unmarshal(httpRspBody, &respbody)

		if respbody.Errcode != 0 {
			fmt.Println("Deploy code error! ", respbody.Errcode, ":", respbody.Msg)
			return

		}

		deployCodeRsp = &respbody
	}

	/*fmt.Printf("\nPush transaction done for deploying contract %v.\n", name)
	fmt.Printf("Trx: \n")

	type PrintDeployCodeParam struct {
		Name         string `json:"name"`
		VMType       byte   `json:"vm_type"`
		VMVersion    byte   `json:"vm_version"`
		ContractCode string `json:"contract_code"`
	}

	pdcp := &PrintDeployCodeParam{}
	pdcp.Name = name
	pdcp.VMType = byte(fileType)
	pdcp.VMVersion = 1
	codeHex := BytesToHex(ContractCodeVal[0:100])
	pdcp.ContractCode = codeHex + "..."

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       pdcp,
		ParamBin:    string([]byte(trx.Param)[0:200]) + "...",
		//ParamBin: trx.Param,
		SigAlg:    trx.SigAlg,
		Signature: trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)*/
	fmt.Printf("\nTrxHash: %v\n", deployCodeRsp.Result.TrxHash)
	fmt.Printf("\nThis transaction is sent. Please check its result by command : bcli transaction get --trxhash  <hash>\n\n")
}

func checkAbi(abiRaw []byte) error {
	_, err := abi.ParseAbi(abiRaw)
	if err != nil {
		return fmt.Errorf("ABI Parse error: %v", err)
	}
	return nil
}

func (cli *CLI) deployabi(name string, path string) {
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("GetInfo error: ", err)
		return
	}

	_, err = ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Open abi file error: ", err)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Open abi file error: ", err)
		return
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("Open abi file error: ", err)
		return
	}

	tempAbi := make([]byte, fi.Size())
	f.Read(tempAbi)

	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
	   fmt.Println("getAbibyContractName of bottos failed!")
           return
        }
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "contract", name)
	abi.Setmapval(mapstruct, "contract_abi", tempAbi)
	
	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "deployabi")

	trx1 := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      name,
		Contract:    "bottos",
		Method:      "deployabi",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}

	sign, err := cli.signTrx(trx1, param)
	if err != nil {
		return
	}
	
	http_method := "restful"
	trx1.Signature = sign
	
	var deployAbiRsp *chain.SendTransactionResponse
	
	if http_method == "grpc" {
		deployAbiRsp, err = cli.client.SendTransaction(context.TODO(), trx1)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		http_url := "http://"+ChainAddr+ "/v1/transaction/send"
		req, _ := json.Marshal(trx1)
    		req_new := bytes.NewBuffer([]byte(req))
		httpRspBody, err := send_httpreq("POST", http_url, req_new)
		if err != nil || httpRspBody == nil {
			fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
			return
		}
		var respbody chain.SendTransactionResponse
		json.Unmarshal(httpRspBody, &respbody)
		if respbody.Errcode != 0 {
		    fmt.Println("Error! ",respbody.Errcode, ":", respbody.Msg)
		    return
		}
		deployAbiRsp = &respbody
	}

	b, _ := json.Marshal(deployAbiRsp)
	cli.jsonPrint(b)
}

func GetAbiOverHttp(http_url string, contract string) (/**chain.GetInfoResponse_Result*/string, error) {
		getinfo := &chain.GetAbiRequest{Contract: contract}
		req, _ := json.Marshal(getinfo)
		req_new := bytes.NewBuffer([]byte(req))
		httpRspBody, err := send_httpreq("POST", http_url, req_new)
		if err != nil || httpRspBody == nil {
			fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
			return "", errors.New("Error!")
		}
		
		var trxrespbody  TODO.ResponseStruct
		
		err = json.Unmarshal(httpRspBody, &trxrespbody)
		
		if err != nil {
		    fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		    return "", errors.New("Error!")
		} else if trxrespbody.Errcode != 0 {
		    fmt.Println("Error! ",trxrespbody.Errcode, ":", trxrespbody.Msg)
		    return "", errors.New("Error!")
			
		}
		
		//b, _ := json.Marshal(trxrespbody.Result)
		//cli.jsonPrint(b)
		//var abiInfo chain.GetAbiResponse
		//json.Unmarshal(b, &abiInfo)
		
	return trxrespbody.Result.(string), nil 
}

//BytesToHex hex encode
func BytesToHex(d []byte) string {
	return hex.EncodeToString(d)
}

//HexToBytes hex decode
func HexToBytes(str string) ([]byte, error) {
	h, err := hex.DecodeString(str)

	return h, err
}
