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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"

	chain "github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/crypto-go/crypto"
	proto "github.com/golang/protobuf/proto"
	"github.com/micro/go-micro"

	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"
)

// CLI responsible for processing command line arguments
type CLI struct {
	client chain.ChainService
}

//Transaction trx info
type Transaction struct {
	Version     uint32      `json:"version"`
	CursorNum   uint32      `json:"cursor_num"`
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

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  newaccount -name NAME -pubkey PUBKEY         - Create a New account")
	fmt.Println("  getaccount -name NAME                        - Get account balance")
	fmt.Println("  transfer -from FROM -to TO -amount AMOUNT    - Transfer BTO from FROM account to TO")
	fmt.Println("  deploycode -contract NAME -wasm PATH         - Deploy contract NAME from .wasm file")
	fmt.Println("  deployabi -contract NAME -abi PATH           - Deploy contract ABI from .abi file")
	fmt.Println("")
}

//NewCLI new console client
func NewCLI() *CLI {
	cli := &CLI{}
	service := micro.NewService()
	service.Init()
	cli.client = chain.NewChainService("bottos", service.Client())

	return cli
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
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

	data, err := proto.Marshal(ctrx)
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

func (cli *CLI) transfer(from, to string, amount int) {
	chainInfo, err := cli.getChainInfo()
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}

	type TransferParam struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount uint64 `json:"amount"`
	}
	var value uint64
	value = uint64(amount) * uint64(100000000)
	tp := &TransferParam{
		From:   from,
		To:     to,
		Amount: value,
	}

	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		return
	}

	param, _ := abi.MarshalAbi(tp, &Abi, "bottos", "transfer")

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

	newAccountRsp, err := cli.client.SendTransaction(context.TODO(), trx)
	if err != nil || newAccountRsp == nil {
		fmt.Println(err)
		return
	}

	if newAccountRsp.Errcode != 0 {
		fmt.Printf("Transfer error:\n")
		fmt.Printf("    %v\n", newAccountRsp.Msg)
		return
	}

	fmt.Printf("Transfer Succeed\n")
	fmt.Printf("    From: %v\n", from)
	fmt.Printf("    To: %v\n", to)
	fmt.Printf("    Amount: %v\n", amount)
	fmt.Printf("Trx: \n")

	tp.Amount = uint64(amount)
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
	NodeIp := "127.0.0.1"
	addr := "http://" + NodeIp + ":8080/rpc"
	params := `service=bottos&method=CoreApi.QueryAbi&request={
			"contract":"%s"}`
	s := fmt.Sprintf(params, contractname)
	respBody, err := http.Post(addr, "application/x-www-form-urlencoded",
		strings.NewReader(s))

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
		fmt.Println(err)
		return abi.ABI{}, err
	}

	Abi, err := abi.ParseAbi([]byte(abistring))
	if err != nil {
		fmt.Println("Parse abistring", abistring, " to abi failed!")
		return abi.ABI{}, err
	}

	return *Abi, nil
}

func (cli *CLI) newaccount(name string, pubkey string) {
	chainInfo, err := cli.getChainInfo()
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
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

	param, _ := abi.MarshalAbi(nps, &Abi, "bottos", "newaccount")

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

	rsp, err := cli.client.SendTransaction(context.TODO(), trx)
	if err != nil || rsp == nil {
		fmt.Println(err)
		return
	}

	if rsp.Errcode != 0 {
		fmt.Printf("Newaccount error:\n")
		fmt.Printf("    %v\n", rsp.Msg)
		return
	}

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
	accountRsp, err := cli.client.GetAccount(context.TODO(), &chain.GetAccountRequest{AccountName: name})
	if err != nil || accountRsp == nil {
		return
	}

	if accountRsp.Errcode == 10204 {
		fmt.Printf("Account: %s Not Exist\n", name)
		return
	}

	account := accountRsp.GetResult()
	fmt.Printf("    Account: %s\n", account.AccountName)
	fmt.Printf("    Balance: %d.%08d BTO\n", account.Balance/100000000, account.Balance%100000000)
}

func (cli *CLI) deploycode(name string, path string) {
	chainInfo, err := cli.getChainInfo()
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
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

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("Open wasm file error: ", err)
		return
	}

	type DeployCodeParam struct {
		Name         string `json:"contract"`
		VMType       byte   `json:"vm_type"`
		VMVersion    byte   `json:"vm_version"`
		ContractCode []byte `json:"contract_code"`
	}

	dcp := &DeployCodeParam{
		Name:      name,
		VMType:    1,
		VMVersion: 1,
	}
	dcp.ContractCode = make([]byte, fi.Size())
	f.Read(dcp.ContractCode)
	//fmt.Printf("Code %x", dcp.ContractCode)

	Abi, abierr := getAbibyContractName("bottos")
	if abierr != nil {
		return
	}

	param, _ := abi.MarshalAbi(dcp, &Abi, "bottos", "deploycode")

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

	deployCodeRsp, err := cli.client.SendTransaction(context.TODO(), trx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if deployCodeRsp.Errcode != 0 {
		fmt.Printf("Deploy contract error:\n")
		fmt.Printf("    %v\n", deployCodeRsp.Msg)
		return
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
	pdcp.Name = dcp.Name
	pdcp.VMType = dcp.VMType
	pdcp.VMVersion = dcp.VMVersion
	codeHex := BytesToHex(dcp.ContractCode[0:100])
	pdcp.ContractCode = codeHex + "..."
	//pdcp.ContractCode = BytesToHex(dcp.ContractCode)

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
	chainInfo, err := cli.getChainInfo()
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
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

	type DeployAbiParam struct {
		Name        string `json:"contract"`
		ContractAbi []byte `json:"contract_abi"`
	}

	dcp := &DeployAbiParam{
		Name: name,
	}
	tempAbi := make([]byte, fi.Size())
	f.Read(tempAbi)
	Abi, err := abi.ParseAbi(tempAbi)
	if err != nil {
		fmt.Printf("Abi Parse Error Hex: %x, Str: %v", tempAbi, string(tempAbi))
		return
	}

	dcp.ContractAbi, err = json.Marshal(Abi)
	if err != nil {
		fmt.Printf("Abi Reformat Error: %v", Abi)
		return
	}

	fmt.Printf("Abi Hex: %x, Str: %v", dcp.ContractAbi, string(dcp.ContractAbi))

	Abi2, abierr2 := getAbibyContractName("bottos")
	if abierr2 != nil {
		return
	}

	param, _ := abi.MarshalAbi(dcp, &Abi2, "bottos", "deployabi")

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

	trx1.Signature = sign

	deployAbiRsp, err := cli.client.SendTransaction(context.TODO(), trx1)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(deployAbiRsp)
	cli.jsonPrint(b)
}

// Run parses command line arguments and processes commands
func (cli *CLI) Run() {
	cli.validateArgs()

	transferCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	sendfrom := transferCmd.String("from", "", "transfer from")
	sendto := transferCmd.String("to", "", "transfer to")
	sendamount := transferCmd.Int("amount", 0, "transfer amount")

	NewAccountCmd := flag.NewFlagSet("newaccount", flag.ExitOnError)
	newAccountName := NewAccountCmd.String("name", "", "account name")
	newAccountPubkey := NewAccountCmd.String("pubkey", "", "pubkey")

	GetAccountCmd := flag.NewFlagSet("getaccount", flag.ExitOnError)
	getAccountName := GetAccountCmd.String("name", "", "account name")

	deployCodeCmd := flag.NewFlagSet("deploycode", flag.ExitOnError)
	deployCodeName := deployCodeCmd.String("contract", "", "contract name")
	deployCodePath := deployCodeCmd.String("wasm", "", ".wasm file path")

	deployAbiCmd := flag.NewFlagSet("deployabi", flag.ExitOnError)
	deployAbiName := deployAbiCmd.String("contract", "", "contract name")
	deployAbiPath := deployAbiCmd.String("abi", "", ".abi file path")

	switch os.Args[1] {
	case "transfer":
		err := transferCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "newaccount":
		err := NewAccountCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getaccount":
		err := GetAccountCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "deploycode":
		err := deployCodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "deployabi":
		err := deployAbiCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if transferCmd.Parsed() {
		if *sendfrom == "" || *sendto == "" || *sendamount <= 0 {
			transferCmd.Usage()
			os.Exit(1)
		}

		cli.transfer(*sendfrom, *sendto, *sendamount)
	}

	if NewAccountCmd.Parsed() {
		if *newAccountName == "" || *newAccountPubkey == "" {
			NewAccountCmd.Usage()
			os.Exit(1)
		}

		cli.newaccount(*newAccountName, *newAccountPubkey)
	}

	if GetAccountCmd.Parsed() {
		if *getAccountName == "" {
			GetAccountCmd.Usage()
			os.Exit(1)
		}

		cli.getaccount(*getAccountName)
	}

	if deployCodeCmd.Parsed() {
		if *deployCodeName == "" || *deployCodePath == "" {
			deployCodeCmd.Usage()
			os.Exit(1)
		}

		cli.deploycode(*deployCodeName, *deployCodePath)
	}

	if deployAbiCmd.Parsed() {
		if *deployAbiName == "" || *deployAbiPath == "" {
			deployAbiCmd.Usage()
			os.Exit(1)
		}

		cli.deployabi(*deployAbiName, *deployAbiPath)
	}
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
