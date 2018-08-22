/*
@Time : 2018/8/13 10:15 
@Author : 星空之钥丶
@File : main
@Software: GoLand
*/
package main

import (
	"gopkg.in/urfave/cli.v1"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var ChainAddr string = "127.0.0.1:8689"


func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.GlobalSet(name, ctx.String(name))
			}
		}
		if ctx.GlobalIsSet("servaddr") {
			ChainAddr = ctx.GlobalString("servaddr")
		}
		return action(ctx)
	}
}

func BcliGetChainInfo(ctx *cli.Context) error {
	chainInfo, err := getChainInfoOverHttp("http://"+ChainAddr+"/v1/block/height")
	if err != nil {
		fmt.Println("GetInfo error: ", err)
		return nil
	}
	fmt.Printf("\n==Chain Info==\n\n")
	
	b, _ := json.Marshal(chainInfo)
	jsonPrint(b)
	
	return nil
}

func BcliGetBlockInfo(ctx *cli.Context) error {

	num := ctx.Uint64("num")
	hash := ctx.String("hash")

	blockInfo, err := getBlockInfoOverHttp("http://"+ChainAddr+"/v1/block/detail", num, hash)
	if err != nil {
		return nil
	}
	fmt.Printf("\n==Block Info==\n\n")
	b, _ := json.Marshal(blockInfo)
	jsonPrint(b)
	return nil
}

func BcliNewAccount(ctx *cli.Context) error {

	username := ctx.String("username")
	pubkey := ctx.String("pubkey")
	
	newaccount(username, pubkey)
	
	return nil
}

func BcliGetAccount(ctx *cli.Context) error {

	username := ctx.String("username")
	
	getaccount(username)
	
	return nil
}

func BcliTransfer(ctx *cli.Context) error {

	from := ctx.String("from")
	to   := ctx.String("to")
	amount := ctx.Int("amount")

	transfer(from, to, amount)
	
	return nil
}

func BCLIGetTransaction(ctx *cli.Context) error {

	trxhash := ctx.String("trxhash")

	BcliGetTransaction(trxhash)
	
	return nil
}

func BCLIPushTransaction(ctx *cli.Context) error {
	
	var pushtrx BcliPushTrxInfo
        
        pushtrx.sender   = ctx.String("sender")
        pushtrx.contract = ctx.String("contract")
        pushtrx.method   = ctx.String("method")
	pushtrx.ParamMap = make(map[string]interface{})

	param1 := ctx.String("param")
	param1 = strings.Replace(param1, " ", "", -1)
	param2 := strings.Split(param1, ",")
	for _, item := range(param2) {
		param3 := strings.Split(item, ":")
		pushtrx.ParamMap[param3[0]] = param3[1]
	}

	BcliPushTransaction(&pushtrx)
	
	return nil
}

func BCLIDeployCode(ctx *cli.Context) error {
	name := ctx.String("name")
	code := ctx.String("code")

	deploycode(name, code)

	return nil
}

func BCLIDeployAbi(ctx *cli.Context) error {
	name := ctx.String("name")
	Abi := ctx.String("abi")

	deployabi(name, Abi)

	return nil
}

func BCLIDeployBoth(ctx *cli.Context) error {
	name := ctx.String("name")
	Abi  := ctx.String("abi")
	code := ctx.String("code")
	
	deploycode(name, code)

	deployabi(name, Abi)

	return nil
}

func BCLIGetContractCode(ctx *cli.Context) error {
	name := ctx.String("name")
	SaveToAbiPath  := ctx.String("abi")
	SaveTocodePath := ctx.String("code")
	
	BcliGetContractCode(name, SaveTocodePath, SaveToAbiPath)

	return nil
}

func BCLIGetTableInfo(ctx *cli.Context) error {
	contract := ctx.String("contract")
	table := ctx.String("table")
	key  := ctx.String("key")
	
	BCliGetTableInfo(contract, table, key)

	return nil
}



func isNotEmpty(str string) bool {
	if len(str) > 0 {
		return true
	}
	return false
}

func validatorUsername(str string) (bool,error) {
	if !isNotEmpty(str) {
		return false,  fmt.Errorf("Parameter anomaly！")
	}

	match, err := regexp.MatchString("^[a-z][a-z1-9]{2,15}$", str);
	if err != nil {
		return false, err
	}

	if !match {
		return false, fmt.Errorf("Error parameter!")
	}

	return true, nil
}

