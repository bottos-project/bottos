package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"io/ioutil"
	"encoding/hex"
	"encoding/json"
	"bytes"

	"github.com/micro/go-micro"
	"golang.org/x/net/context"

	coreapi "github.com/bottos-project/bottos/api"
	"github.com/bottos-project/bottos/contract/msgpack"
)

// CLI responsible for processing command line arguments
type CLI struct {
	client coreapi.CoreApiClient
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  newaccount -name NAME -pubkey PUBKEY - Create a New account")
	fmt.Println("  transfer -from FROM -to TO -amount AMOUNT - transfer bottos from FROM account to TO")
	fmt.Println("  deploycode -contract NAME -wasm PATH - deploy contract NAME from .wasm file")
	fmt.Println("  deployabi -contract NAME -abi PATH - deploy contract NAME from .abi file")
	fmt.Println("")
}

func NewCLI() *CLI {
	cli := &CLI{}
	service := micro.NewService()
	service.Init()
	cli.client = coreapi.NewCoreApiClient("core", service.Client())

	return cli
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) queryChainInfo() (*coreapi.QueryChainInfoResponse_Result, error) {
	chainInfoRsp, err := cli.client.QueryChainInfo(context.TODO(), &coreapi.QueryChainInfoRequest{})
	if err != nil || chainInfoRsp == nil {
		fmt.Println(err)
		return nil, err
	}

	chainInfo := chainInfoRsp.GetResult()
	return chainInfo, nil
}

func (cli *CLI) transfer(from, to string, amount int) {
	chainInfo, err := cli.queryChainInfo()
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}

	type TransferParam struct {
		From string
		To string
		Value uint64
	}
	tp := &TransferParam{
		From: from,
		To: to,
		Value: uint64(amount),
	}
	param, _ := msgpack.Marshal(tp)

	trx := &coreapi.Transaction{
		Version:1,
		CursorNum: chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime: chainInfo.HeadBlockTime+1000,
		Sender:"bottos",
		Contract:"bottos",
		Method:"transfer",
		Param: BytesToHex(param),
		SigAlg:1,
		Signature:string(""),
	}

	newAccountRsp, err := cli.client.PushTrx(context.TODO(), trx)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(newAccountRsp)
	cli.jsonPrint(b)
}

func (cli *CLI) jsonPrint(data []byte) {
	var out bytes.Buffer  
	json.Indent(&out, data, "", "    ")
	
	fmt.Println(string(out.Bytes()))
}

func (cli *CLI) newaccount(name string, pubkey string) {
	chainInfo, err := cli.queryChainInfo()
	if err != nil {
		fmt.Println("QueryChainInfo error: ", err)
		return
	}

	// 1, new account trx
	type NewAccountParam struct {
		Name string
		Pubkey string
	}
	nps := &NewAccountParam{
		Name: name,
		Pubkey: pubkey,
	}
	param, _ := msgpack.Marshal(nps)

	trx := &coreapi.Transaction{
		Version:1,
		CursorNum: chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime: chainInfo.HeadBlockTime+1000,
		Sender:"bottos",
		Contract:"bottos",
		Method:"newaccount",
		Param: BytesToHex(param),
		SigAlg:1,
		Signature:string(""),
	}

	rsp, err := cli.client.PushTrx(context.TODO(), trx)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(rsp)
	cli.jsonPrint(b)
}


func (cli *CLI) deploycode(name string, path string) {
	chainInfo, err := cli.queryChainInfo()
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
		Name		 string		 `json:"name"`
		VMType       byte        `json:"vm_type"`
		VMVersion    byte        `json:"vm_version"`
		ContractCode []byte      `json:"contract_code"`
	}

	dcp := &DeployCodeParam{
		Name: name,
		VMType: 1,
		VMVersion: 1,
	}
	dcp.ContractCode = make([]byte, fi.Size())
	f.Read(dcp.ContractCode)
	//fmt.Printf("Code %x", dcp.ContractCode)
	param, _ := msgpack.Marshal(dcp)

	trx1 := &coreapi.Transaction{
		Version:1,
		CursorNum: chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime: chainInfo.HeadBlockTime+1000,
		Sender:"bottos",
		Contract:"bottos",
		Method:"deploycode",
		Param: BytesToHex(param),
		SigAlg:1,
		Signature:string(""),
	}
	deployCodeRsp, err := cli.client.PushTrx(context.TODO(), trx1)
	if err != nil {
		fmt.Println(err)
		return
	}

	b, _ := json.Marshal(deployCodeRsp)
	cli.jsonPrint(b)
}


func (cli *CLI) deployabi(name string, path string) {
	chainInfo, err := cli.queryChainInfo()
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
		Name		 string		 `json:"name"`
		ContractAbi  []byte      `json:"contract_abi"`
	}

	dcp := &DeployAbiParam{
		Name: name,
	}
	dcp.ContractAbi = make([]byte, fi.Size())
	f.Read(dcp.ContractAbi)
	fmt.Printf("Abi Hex: %x, Str: %v", dcp.ContractAbi, string(dcp.ContractAbi))
	param, _ := msgpack.Marshal(dcp)

	trx1 := &coreapi.Transaction{
		Version:1,
		CursorNum: chainInfo.HeadBlockNum,
		CursorLabel: chainInfo.CursorLabel,
		Lifetime: chainInfo.HeadBlockTime+1000,
		Sender:"bottos",
		Contract:"bottos",
		Method:"deployabi",
		Param: BytesToHex(param),
		SigAlg:1,
		Signature:string(""),
	}
	deployAbiRsp, err := cli.client.PushTrx(context.TODO(), trx1)
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

func BytesToHex(d []byte) string {
	return hex.EncodeToString(d)
}

func HexToBytes(str string) ([]byte, error) {
	h, err := hex.DecodeString(str)

	return h, err
}

