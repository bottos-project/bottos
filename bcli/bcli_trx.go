package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/context"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/common"
	chain "github.com/bottos-project/bottos/api"
)

type BcliPushTrxInfo struct {
	sender string
	contract string
	method string
	ParamMap map[string]interface{}
}

func (cli *CLI) BcliGetTransaction (trxhash common.Hash) {
	
}

func (cli *CLI) BcliPushTransaction (pushtrxinfo *BcliPushTrxInfo) {
	
	Abi, abierr := getAbibyContractName(pushtrxinfo.contract)
        if abierr != nil {
           return
        }
	
	chainInfo, err := cli.getChainInfo()
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
		Contract:    "nodeclustermng",
		Method:      "reg",
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
