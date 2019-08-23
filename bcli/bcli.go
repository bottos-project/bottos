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

func (cli *CLI) getBlockInfoOverHttp(http_url string, block_num uint64, block_hash string) (*chain.GetBlockResponse_Result, error) {
		var getinfo *chain.GetBlockRequest
		if block_num >= 0 {
			getinfo = &chain.GetBlockRequest{BlockNum:block_num}
		} else if len(block_hash) > 0 {
			getinfo = &chain.GetBlockRequest{BlockHash:block_hash}
		} else {
			getinfo = &chain.GetBlockRequest{BlockNum:0}
		}

		req, _ := json.Marshal(getinfo)
		req_new := bytes.NewBuffer([]byte(req))
		httpRspBody, err := send_httpreq("POST", http_url, req_new)
		if err != nil || httpRspBody == nil {
			fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
			return nil, errors.New("Error!")
		}
		
		var trxrespbody  TODO.ResponseStruct
		
		err = json.Unmarshal(httpRspBody, &trxrespbody)
		
		if err != nil {
		    fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		    return nil, errors.New("Error!")
		} else if trxrespbody.Errcode != 0 {
		    fmt.Println("Error! ",trxrespbody.Errcode, ":", trxrespbody.Msg)
		    return nil, errors.New("Error!")
			
		}
		
		b, _ := json.Marshal(trxrespbody.Result)
		//cli.jsonPrint(b)
		var blockInfo chain.GetBlockResponse_Result
		json.Unmarshal(b, &blockInfo)
		
	return &blockInfo, nil
}

func (cli *CLI) getAccountInfoOverHttp(name string, http_url string) (*chain.GetAccountResponse_Result, error) {
		
		getinfo := &chain.GetAccountRequest{AccountName:name}
		req, _ := json.Marshal(getinfo)
		req_new := bytes.NewBuffer([]byte(req))
		httpRspBody, err := send_httpreq("POST", http_url, req_new)
		if err != nil || httpRspBody == nil {
			fmt.Println("Error. httpRspBody: ", httpRspBody, ", err: ", err)
			return nil, errors.New("Error!")
		}
		
		var respbody  TODO.ResponseStruct
		
		err = json.Unmarshal(httpRspBody, &respbody)
		
		if err != nil {
		    fmt.Println("Error! Unmarshal to trx failed: ", err, "| body is: ", string(httpRspBody), ". trxrsp:")
		    return nil, errors.New("Error!")
		} else if respbody.Errcode != 0 {
		    fmt.Println("Error! ", respbody.Errcode, ":", respbody.Msg)
		    return nil, errors.New("Error!") 
		}
		
		b, _ := json.Marshal(respbody.Result)
		//cli.jsonPrint(b)
		var accountInfo chain.GetAccountResponse_Result
		json.Unmarshal(b, &accountInfo)
		
	return &accountInfo, nil
}

func (cli *CLI) signTrx(trx *chain.Transaction, param []byte) (string, error) {
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
	seckey, err := GetDefaultKey()
	signdata, err := crypto.Sign(hashData, seckey)

	return BytesToHex(signdata), err
}

func (cli *CLI) transfer(from string, to string, amount big.Int) {
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("GetInfo error: ", err)
		return
	}

	type TransferParam struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount big.Int `json:"value"`
	}

	value := big.NewInt(100000000)
	value2 := big.NewInt(0)

	value2, _ = safemath.U256Mul(value2, &amount, value)
	tp := &TransferParam{
		From:   from,
		To:     to,
		Amount: *value2,
	}

	Abi, abierr := getAbibyContractName("bottos")
        if abierr != nil {
           return
        }

	mapstruct := make(map[string]interface{})
	abi.Setmapval(mapstruct, "from", from)
	abi.Setmapval(mapstruct, "to", to)
	abi.Setmapval(mapstruct, "value", *value2)
	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "transfer")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      from,
		Contract:    "bottos",
		Method:      "transfer",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}

	sign, err := cli.signTrx(trx, param)
	if err != nil {
		return
	}

	trx.Signature = sign

	/*newAccountRsp, err := cli.client.SendTransaction(context.TODO(), trx)
	if err != nil || newAccountRsp == nil {
		fmt.Println(err)
		return
	}

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %v\n", newAccountRsp.Msg)
		return
	}*/
	
	req, _ := json.Marshal(trx)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", "http://" + ChainAddr + "/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("Error! ", respbody.Errcode, ":", respbody.Msg)
		return 
	}
	newAccountRsp := &respbody


	fmt.Printf("Transfer Succeed\n")
	fmt.Printf("    From: %v\n", from)
	fmt.Printf("    To: %v\n", to)
	fmt.Println("    Amount:", value2)
	fmt.Printf("Trx: \n")

	tp.Amount = amount
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
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", newAccountRsp.Result.TrxHash)
}

func (cli *CLI) jsonPrint(data []byte) {
	var out bytes.Buffer
	json.Indent(&out, data, "", "    ")

	fmt.Println(string(out.Bytes()))
}

//getAbibyContractName function
func getAbibyContractName(contractname string) (abi.ABI, error) {
	var abistring string
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
	abistring, err := GetAbiOverHttp(http_url, contractname)
	if len(abistring) <= 0 || err != nil {
		return abi.ABI{}, errors.New("len(abistring) <= 0")
	}

	Abi, err := abi.ParseAbi([]byte(abistring))
	if err != nil {
		fmt.Println("Parse abistring", abistring, " to abi failed!")
		return abi.ABI{}, err
	}

	return *Abi, nil
}

func (cli *CLI) newaccount(name string, pubkey string) {

	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
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

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      "delta",
		Contract:    "bottos",
		Method:      "newaccount",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}

	sign, err := cli.signTrx(trx, param)
	if err != nil {
		return
	}

	trx.Signature = sign

	/*rsp, err := cli.client.SendTransaction(context.TODO(), trx)
	if err != nil || rsp == nil {
		fmt.Println(err)
		return
	}

	if rsp.Errcode != 0 {
		fmt.Printf("Newaccount error:\n")
		fmt.Printf("    %v\n", rsp.Msg)
		return
	} */
	
	req, _ := json.Marshal(trx)
	req_new := bytes.NewBuffer([]byte(req))
	httpRspBody, err := send_httpreq("POST", "http://" + ChainAddr + "/v1/transaction/send", req_new)
	if err != nil || httpRspBody == nil {
		fmt.Println("BcliPushTransaction Error:", err, ", httpRspBody: ", httpRspBody)
		return
	}
	
	var respbody chain.SendTransactionResponse
	json.Unmarshal(httpRspBody, &respbody)
	if respbody.Errcode != 0 {
		fmt.Println("Error! ", respbody.Errcode, ":", respbody.Msg)
		return 
	}
	rsp := &respbody
	

	fmt.Printf("Create account: %v Succeed\n", name)
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
		ParamBin:    trx.Param,
		SigAlg:      trx.SigAlg,
		Signature:   trx.Signature,
	}

	b, _ := json.Marshal(printTrx)
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", rsp.Result.TrxHash)
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
	mulval  := big.NewInt(100000000)

	balanceResult, result := balance.SetString(account.Balance, 10)
	if false == result {
		fmt.Println("Error: balance.SetString failed. account: ", account)
		return
	}
	
	var mulrestlt *big.Int = big.NewInt(0)
	var modrestlt *big.Int = big.NewInt(0)

	mulrestlt, err = safemath.U256Div(mulrestlt, balanceResult, mulval) 

	if err != nil {
		return
	}
	
	mulval2  := big.NewInt(100000000)
	modrestlt, err = safemath.U256Mod(modrestlt, balanceResult, mulval2)
	if err != nil {
		return
	}

	fmt.Printf("    Account: %s\n", account.AccountName)
	fmt.Printf("    Balance: %d.%08d BTO\n", mulrestlt, modrestlt)
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

func (cli *CLI) deploycode(name string, path string, fileTypeInput string) {
	
	//chainInfo, err := cli.getChainInfo()
	infourl := "http://" + ChainAddr + "/v1/block/height"
	chainInfo, err := cli.GetChainInfoOverHttp(infourl)
	
	if err != nil {
		fmt.Println("GetInfo error: ", err)
		return
	}
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
        abi.Setmapval(mapstruct, "vm_version", uint8(0))
	
	abi.Setmapval(mapstruct, "contract_code", ContractCodeVal)
	//fmt.Printf("contract_code: %x", ContractCodeVal)
	param, _ := abi.MarshalAbiEx(mapstruct, &Abi, "bottos", "deploycode")

	trx := &chain.Transaction{
		Version:     1,
		CursorNum:   chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime:    chainInfo.HeadBlockTime + 100,
		Sender:      name,
		Contract:    "bottos",
		Method:      "deploycode",
		Param:       BytesToHex(param),
		SigAlg:      1,
	}

	sign, err := cli.signTrx(trx, param)
	if err != nil {
		return
	}

	trx.Signature = sign
	var deployCodeRsp *chain.SendTransactionResponse
	
	http_method := "restful"

	if http_method == "grpc" {
		deployCodeRsp, err = cli.client.SendTransaction(context.TODO(), trx)
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
		
		if respbody.Errcode != 0 {
		    fmt.Println("Error! ", respbody.Errcode, ":", respbody.Msg)
		    return
			
		}

		deployCodeRsp = &respbody
	}
	
	fmt.Printf("Deploy contract: %v Succeed\n", name)
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
	cli.jsonPrint(b)
	fmt.Printf("TrxHash: %v\n", deployCodeRsp.Result.TrxHash)
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
