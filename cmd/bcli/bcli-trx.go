package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"errors"
	"io"
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

func BcliGetTransaction (trxhash string) {
	
	var newAccountRsp *chain.GetTransactionResponse
	var err error
	
	http_method := "restful"

	if http_method == "grpc" {
		gettrx := &chain.GetTransactionRequest{trxhash}
		newAccountRsp, err = client.GetTransaction(context.TODO(), gettrx)
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
		jsonPrint(b)

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
		
		var trxrespbody TODO.Todo
		
		err = json.Unmarshal(httpRspBody, &trxrespbody)
		
		if err != nil {
		    fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		    return
		}
		
		b, _ := json.Marshal(trxrespbody.Result)
		jsonPrint(b)
	}
}

func BcliPushTransaction (pushtrxinfo *BcliPushTrxInfo) {

	Abi, abierr := getAbibyContractName(pushtrxinfo.contract)
        if abierr != nil {
	   fmt.Println("Push Transaction fail due to get Abi failed.")
           return
        }
	
	//chainInfo, err := getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := getChainInfoOverHttp(infourl)
	
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
	
	sign, err := signTrx(trx, param)
	if err != nil {
	   	fmt.Println("Push Transaction fail due to sign Trx failed.")
		return
	}
	
	http_method := "restful"
	trx.Signature = sign
	var newAccountRsp *chain.SendTransactionResponse
	
	if http_method == "grpc" {
		newAccountRsp, err = client.SendTransaction(context.TODO(), trx)
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
	jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

type GetContractCodeAbi struct {
	Contract string `json:"contract"`
}

func BcliGetContractCode (contract string, save_to_wasm_path string, save_to_abi_path string) (string, string) {
	
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
	
	var trxrespbody TODO.Todo
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
	jsonPrint(b)

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


func BCliGetTableInfo (contract string, table string, key string) {

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
