﻿package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"errors"
	"io"
	"math/big"
	"golang.org/x/net/context"
	"github.com/bottos-project/bottos/contract/abi"
	chain "github.com/bottos-project/bottos/api"
	TODO "github.com/bottos-project/bottos/restful/handler"
)

type BcliPushTrxInfo struct {
	sender string
	contract string
	method string
	ParamMap map[string]interface{}
}

func send_httpreq (get_or_post string, ReqUrl string, ReqMsg io.Reader) ([]byte, error) {
    var err error	
    
    client := &http.Client{}
    req, _ := http.NewRequest(get_or_post, ReqUrl, ReqMsg)
    //req.Header.Set("Connection", "keep-alive")
    resp, err := client.Do(req)
    
    if err != nil || resp == nil {
	fmt.Println("Error: get_or_post:", get_or_post, ", resp:", resp, ",err: ", err)
	return nil, errors.New("Error: send http failed")
    }	   
	
    if resp.StatusCode != 200 {
	fmt.Println("Error: get_or_post:", get_or_post, ", ReqUrl:",ReqUrl, ",resp: ", resp)
	return nil, errors.New("Error: resp.StatusCode is not 200")
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
		return nil, err
    }

    if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
    }
    
    return body, nil
}

type GetTransactionResponse struct {
	Errcode uint32       `protobuf:"varint,1,opt,name=errcode" json:"errcode,omitempty"`
	Msg     string       `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	Result  interface{} `protobuf:"bytes,3,opt,name=result" json:"result,omitempty"`
}

type ResponseStruct struct {
	Errcode uint32      `json:"errcode"`
	Msg     string      `json:"msg"`
	Result  interface{} `json:"result"`
}

func Sha256(msg []byte) []byte {
	sha := sha256.New()
	sha.Write([]byte(hex.EncodeToString(msg)))
	return sha.Sum(nil)
}

// call wallet's v1/wallet/signhash
type SignDataResponse_Result struct {
	SignValue string `protobuf:"bytes,1,opt,name=sign_value,json=signValue" json:"sign_value"`
}

type SignDataResponse struct {
	Errcode uint32      `protobuf:"varint,1,opt,name=errcode" json:"errcode"`
	Msg     string      `protobuf:"bytes,2,opt,name=msg" json:"msg"`
	Result  interface{} `protobuf:"bytes,3,opt,name=result" json:"result"`
}

func SignHash(digest []byte, account string, url string) ([]byte, error, berr.ErrCode) {
	values := map[string]interface{}{
		"account_name": account,
		"type":         "normal",
		"hash":         common.BytesToHex(digest),
	}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("wallet signature failed1: %s, %v", err, resp), berr.RestErrInternal
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("wallet signature failed2: %v", resp), berr.RestErrInternal
	}

	var respStruct SignDataResponse
	err = json.Unmarshal(body, &respStruct)
	if err != nil || respStruct.Errcode != uint32(berr.ErrNoError) {
		return nil, fmt.Errorf("wallet signature failed3!"), berr.ErrCode(respStruct.Errcode)
	} else if respStruct.Result == nil {
		return nil, fmt.Errorf("respStruct.Result is empty!"), berr.ErrCode(respStruct.Errcode)
	}

	var respStruct2 SignDataResponse_Result
	b, _ := json.Marshal(respStruct.Result)
	err = json.Unmarshal(b, &respStruct2)

	if err != nil {
		return nil, fmt.Errorf("wallet signature failed4!"), berr.ErrCode(respStruct.Errcode)
	}

	signdata, err := common.HexToBytes(respStruct2.SignValue)
	return signdata, err, 0
}

func (cli *CLI) BcliSignTrxOverHttp(http_url string, sender string, contract string, method string, param string) (*chain.Transaction, error) {

	chaininfo, err := cli.GetChainInfoOverHttp("http://" + ChainAddr + "/v1/block/height")
	if err != nil {
		fmt.Println("cli.GetChainInfoOverHttp failed!")
		return nil, err
	}

	//fmt.Printf("Current block num: %d, version: %d [%s] ", num, blockInfo.VersionNum, version.ParseStringVersion(blockInfo.VersionNum))
	param_bin, _ := common.HexToBytes(param)
	basictrx := &types.BasicTransaction{
		Version:     chaininfo.HeadBlockVersion,
		CursorNum:   chaininfo.HeadBlockNum,
		CursorLabel: chaininfo.CursorLabel,
		Lifetime:    chaininfo.HeadBlockTime + 100,
		Sender:      sender,
		Contract:    contract,
		Method:      method,
		Param:       param_bin,
		SigAlg:      1,
	}

	msg, err := bpl.Marshal(basictrx)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	//Add chainID Flag
	chainID, _ := hex.DecodeString(chaininfo.ChainId)
	msg = bytes.Join([][]byte{msg, chainID}, []byte{})
	hash := Sha256(msg)

	intTrx := &chain.Transaction{
		Version:     chaininfo.HeadBlockVersion,
		CursorNum:   chaininfo.HeadBlockNum,
		CursorLabel: chaininfo.CursorLabel,
		Lifetime:    chaininfo.HeadBlockTime + 100,
		Sender:      sender,
		Contract:    contract,
		Method:      method,
		Param:       param,
		SigAlg:      1,
		//Signature:   signature,
	}

	http_url_sign := "http://" + ChainAddrWallet + "/v1/wallet/signhash"
	tmpval, err, errcode := SignHash(hash, sender, http_url_sign)
	intTrx.Signature = common.BytesToHex(tmpval)

	if err != nil {
		if errcode == berr.RestErrWalletLocked {
			fmt.Printf("\nYour wallet of account [%s] is locked. Please unlock it first.\n\n", sender)
		} else {
			fmt.Println("BcliSignTrxOverHttp failed: ", err)
		}
		return nil, err
	}

	return intTrx, nil
}

func (cli *CLI) BcliGetTransaction (trxhash string) {
	
	var newAccountRsp *chain.GetTransactionResponse
	var err error
	
	http_method := "restful"

	if http_method == "grpc" {
		gettrx := &chain.GetTransactionRequest{trxhash}
		newAccountRsp, err = cli.client.GetTransaction(context.TODO(), gettrx)
		if err != nil || newAccountRsp == nil {
			fmt.Println(err)
			return
		}

		if newAccountRsp.Errcode != 0 {
			fmt.Printf("Transfer error:\n")
			fmt.Printf("    %v\n", newAccountRsp.Msg)
			return
		}

		fmt.Printf("GetTransaction Succeed\n")

		b, _ := json.Marshal(newAccountRsp.Result)
		cli.jsonPrint(b)

		fmt.Printf("Trx: %v\n", newAccountRsp.Result)

	} else {
		http_url := "http://"+ChainAddr+ "/v1/transaction/get"
		gettrx := &chain.GetTransactionRequest{trxhash}
		req, _ := json.Marshal(gettrx)
		req_new := bytes.NewBuffer([]byte(req))
		httpRspBody, err := send_httpreq("POST", http_url, req_new)
		if err != nil || httpRspBody == nil {
			fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
			return
		}
		
		var trxrespbody TODO.ResponseStruct
		
		err = json.Unmarshal(httpRspBody, &trxrespbody)
		
		if err != nil {
		    fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		    return
		}
		
		b, _ := json.Marshal(trxrespbody.Result)
		cli.jsonPrint(b)
	}
}

func (cli *CLI) BcliPushTransaction (pushtrxinfo *BcliPushTrxInfo) {

	Abi, abierr := getAbibyContractName(pushtrxinfo.contract)
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed:", pushtrxinfo.contract)
           return
        }
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}
	
	mapstruct := make(map[string]interface{})
	
	for key, value := range(pushtrxinfo.ParamMap) {
        	abi.Setmapval(mapstruct, key, value)
	}

	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, pushtrxinfo.contract, pushtrxinfo.method)

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      pushtrxinfo.sender,
		Contract:    pushtrxinfo.contract,
		Method:      pushtrxinfo.method,
		Param:       BytesToHex(param),
		SigAlg:      1,
	}
	
	sign, err := cli.signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	http_method := "restful"
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	if http_method == "grpc" {
		newAccountRsp, err = cli.client.SendTransaction(context.TODO(), trx)
		if err != nil || newAccountRsp == nil {
			fmt.Println(err)
	   		fmt.Println("Push Transaction fail due to get grpc response failed.")
			return
		}
	} else {
		http_url := "http://"+ChainAddr+ "/v1/transaction/send"
		req, _ := json.Marshal(trx)
    		req_new := bytes.NewBuffer([]byte(req))
		httpRspBody, err := send_httpreq("POST", http_url, req_new)
		if err != nil || httpRspBody == nil {
			fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
			return
		}
		var respbody chain.SendTransactionResponse
		json.Unmarshal(httpRspBody, &respbody)
		newAccountRsp = &respbody
	}

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}

	fmt.Printf("Transfer Succeed:\n")
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       BytesToHex(param),
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

type GetContractCodeAbi struct {
	Contract string `json:"contract"`
}

func (cli *CLI) BcliGetContractCode (contract string, save_to_wasm_path string, save_to_abi_path string) (string, string) {
	
	var err error
	
	httpurl_contractcode := "http://"+ChainAddr+"/v1/contract/code"
	getcontract := &GetContractCodeAbi{Contract: contract}
	req, _ := json.Marshal(getcontract)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", httpurl_contractcode, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
		return "", ""
	}
	
	var trxrespbody TODO.ResponseStruct
	var contractcode string
	var abivalue     string
	
	err = json.Unmarshal(httpRspBody, &trxrespbody)
	contractcode = trxrespbody.Result.(string)
	if err != nil {
	    fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
	    return "", ""
	}
	
	fmt.Println("\n============CONTRACTCODE===============\n", req_new)
	b, _ := json.Marshal(trxrespbody.Result)
	cli.jsonPrint(b)

	http_urlAbi := "http://"+ChainAddr+ "/v1/contract/abi"

	getcontract = &GetContractCodeAbi{Contract: contract}
	req, _ = json.Marshal(getcontract)
	req_new = bytes.NewBuffer([]byte(req))
	
	httpRspBody, err = send_httpreq("POST", http_urlAbi, req_new)
	
	if err != nil || httpRspBody == nil {
		fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
		return "", ""
	}
	
	err = json.Unmarshal(httpRspBody, &trxrespbody)
	abivalue = trxrespbody.Result.(string)

	if err != nil {
	    fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
	    return "", ""
	}
	
	fmt.Println("\n============ABI===============\n", req_new)
	fmt.Println(trxrespbody.Result)
	
	writeFileToBinary(contractcode, save_to_wasm_path)
        ioutil.WriteFile(save_to_abi_path, []byte(abivalue), 0644)

	return contractcode, abivalue
}


func (cli *CLI) BCliGetTableInfo (contract string, table string, key string) {

		http_url := "http://"+ChainAddr+ "/v1/common/query"
		GetKeyReq := chain.GetKeyValueRequest{ Contract: contract, Object:table, Key: key }

		req, _ := json.Marshal(GetKeyReq)
    		req_new := bytes.NewBuffer([]byte(req))
		httpRspBody, err := send_httpreq("POST", http_url, req_new)
		if err != nil || httpRspBody == nil {
			fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
			return
		}
		var respbody chain.GetKeyValueResponse
		json.Unmarshal(httpRspBody, &respbody)
		newAccountRsp := &respbody

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("BCliGetTableInfo error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}
}

func (cli *CLI) BCliAccountStakeInfo(account string, amount big.Int) {
	
	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed:", "bottos")
           return
        }
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "amount", amount)

	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "stake")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      "delta",
		Contract:    "bottos",
		Method:      "stake",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}
	
	sign, err := cli.signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	http_url := "http://"+ChainAddr+ "/v1/transaction/send"
	req, _ := json.Marshal(trx)
    	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	newAccountRsp = &respbody

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}

	fmt.Printf("Transfer Succeed:\n")
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       BytesToHex(param),
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

func (cli *CLI) BCliAccountUnStakeInfo(account string, amount big.Int) {
	
	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed:", "bottos")
           return
        }
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "amount", amount)

	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "unstake")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      "delta",
		Contract:    "bottos",
		Method:      "unstake",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}
	
	sign, err := cli.signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	http_url := "http://"+ChainAddr+ "/v1/transaction/send"
	req, _ := json.Marshal(trx)
    	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	newAccountRsp = &respbody

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}

	fmt.Printf("Transfer Succeed:\n")
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       BytesToHex(param),
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

func (cli *CLI) BCliAccountClaimInfo(account string, amount big.Int) {
	
	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed:", "bottos")
           return
        }
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "amount", amount)

	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "claim")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      "delta",
		Contract:    "bottos",
		Method:      "claim",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}
	
	sign, err := cli.signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	http_url := "http://"+ChainAddr+ "/v1/transaction/send"
	req, _ := json.Marshal(trx)
    	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	newAccountRsp = &respbody

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}

	fmt.Printf("Transfer Succeed:\n")
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       BytesToHex(param),
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

func (cli *CLI) BCliVoteInfo(vouter string, delegate string) {
	
	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed:", "bottos")
           return
        }
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "voteop", 1)
        abi.Setmapval(mapstruct, "vouter", vouter)
        abi.Setmapval(mapstruct, "delegate", delegate)
	

	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "votedelegate")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      "delta",
		Contract:    "bottos",
		Method:      "votedelegate",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}
	
	sign, err := cli.signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	http_url := "http://"+ChainAddr+ "/v1/transaction/send"
	req, _ := json.Marshal(trx)
    	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	newAccountRsp = &respbody

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}

	fmt.Printf("Transfer Succeed:\n")
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       BytesToHex(param),
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

func (cli *CLI) BCliCancelVoteInfo(vouter string, delegate string) {
	
	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed:", "bottos")
           return
        }
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "voteop", 0)
        abi.Setmapval(mapstruct, "vouter", vouter)
        abi.Setmapval(mapstruct, "delegate", delegate)

	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "votedelegate")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      "delta",
		Contract:    "bottos",
		Method:      "votedelegate",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}
	
	sign, err := cli.signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	http_url := "http://"+ChainAddr+ "/v1/transaction/send"
	req, _ := json.Marshal(trx)
    	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	newAccountRsp = &respbody

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}

	fmt.Printf("Transfer Succeed:\n")
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       BytesToHex(param),
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

func (cli *CLI) BCliDelegateRegInfo(account string, signkey string, location string, description string) {
	
	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed:", "bottos")
           return
        }
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "name", account)
        abi.Setmapval(mapstruct, "pubkey", signkey)
        abi.Setmapval(mapstruct, "location", location)
        abi.Setmapval(mapstruct, "description", description)

	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "regdelegate")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      "delta",
		Contract:    "bottos",
		Method:      "regdelegate",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}
	
	sign, err := cli.signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	http_url := "http://"+ChainAddr+ "/v1/transaction/send"
	req, _ := json.Marshal(trx)
    	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	newAccountRsp = &respbody

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}

	fmt.Printf("Transfer Succeed:\n")
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       BytesToHex(param),
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

func (cli *CLI) BCliDelegateUnRegInfo(account string) {
	
	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed:", "bottos")
           return
        }
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}
	
	mapstruct := make(map[string]interface{})
	
        abi.Setmapval(mapstruct, "name", account)

	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "unregdelegate")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      "delta",
		Contract:    "bottos",
		Method:      "unregdelegate",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}
	
	sign, err := cli.signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	http_url := "http://"+ChainAddr+ "/v1/transaction/send"
	req, _ := json.Marshal(trx)
    	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", http_url, req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	newAccountRsp = &respbody

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %s\n", newAccountRsp)
		return
	}

	fmt.Printf("Transfer Succeed:\n")
	fmt.Printf("Trx: \n")

	printTrx := Transaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       BytesToHex(param),
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

