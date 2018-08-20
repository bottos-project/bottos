/*
@Time : 2018/8/13 10:15 
@Author : 星空之钥丶
@File : main
@Software: GoLand
*/
package main

import (
	cli "gopkg.in/urfave/cli.v1"
	"encoding/json"
	"os"
	"log"
	"fmt"
	"regexp"
	//"encoding/hex"
	//pack "github.com/bottos-project/msgpack-go"
	"github.com/bottos-project/magiccube/service/common/data"
	//user_proto "github.com/bottos-project/magiccube/service/user/proto"
	//push_sign "github.com/bottos-project/magiccube/service/common/signature/push"
	//"github.com/protobuf/proto"
	//"github.com/bottos-project/magiccube/service/common/bean"
	//"bytes"
	//"github.com/bottos-project/crypto-go/crypto"
	//"github.com/bottos-project/magiccube/config"
	//"github.com/bottos-project/magiccube/service/common/util"
)

func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.GlobalSet(name, ctx.String(name))
			}
		}
		return action(ctx)
	}
}

func (cli *CLI) BcliGetChainInfo(ctx *cli.Context) error {
	
	chainInfo, err := cli.getChainInfoOverHttp("http://"+CONFIG.ChainAddr+"/v1/block/height")
	if err != nil {
		fmt.Println("GetInfo error: ", err)
		return nil
	}
	fmt.Printf("\n==Chain Info==\n\n")
	
	b, _ := json.Marshal(chainInfo)
	cli.jsonPrint(b)
	
	return nil
}

func (cli *CLI) BcliGetBlockInfo(ctx *cli.Context) error {

	num := ctx.Uint64("num")
	hash := ctx.String("hash")

	blockInfo, err := cli.getBlockInfoOverHttp("http://"+CONFIG.ChainAddr+"/v1/block/detail", num, hash)
	if err != nil {
		return nil
	}
	fmt.Printf("\n==Block Info==\n\n")
	b, _ := json.Marshal(blockInfo)
	cli.jsonPrint(b)
	return nil
}

func (cli *CLI) BcliNewAccount(ctx *cli.Context) error {

	username := ctx.String("username")
	pubkey := ctx.String("pubkey")
	
	cli.newaccount(username, pubkey)
	
	return nil
}

func (cli *CLI) BcliGetAccount(ctx *cli.Context) error {

	username := ctx.String("username")
	
	cli.getaccount(username)
	
	return nil
}

func (Cli *CLI) RunNewCLI() {
	app := cli.NewApp()
	app.Name = "Bottos Cmd"
	app.Usage = "block chain bcli"
	app.Version = "0.0.1"

	app.Commands = []cli.Command {
		{
			Name: "getinfo",
			Usage: "Get chian info",
			Category: "general",
			Action: MigrateFlags(Cli.BcliGetChainInfo),
		},
		{
			Name: "getblock",
			Usage: "Geeter block info",
			Category: "general",
			Flags: []cli.Flag {
				cli.Uint64Flag{
					Name: "num",
					Value: 100,
					Usage: "get block by number",
				},
				cli.StringFlag{
					Name: "hash",
					Value: "",
					Usage: "get block by hash",
				},
			},
			Action: MigrateFlags(Cli.BcliGetBlockInfo),
		},
		{
			Name: "gettable",
			Usage: "",
			Category: "general",
			Flags: []cli.Flag {
				cli.StringFlag{
					Name: "contract",
					Value:"usermng",
					Usage: "contract name",
				},
				cli.StringFlag{
					Name: "table",
					Usage: "table name",
				},
				cli.StringFlag{
					Name: "key",
					Usage: "key value",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println(c.String("contract"))
				fmt.Println(c.String("table"))
				fmt.Println(c.String("key"))
				return nil
			},
		},
		{
			Name: "account",
			Usage: "Create or Get account",
			Category: "account",
			Subcommands: []cli.Command {
				{
					Name: "create",
					Usage: "Create account",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "username",
							Value:"",
							Usage: "acocunt name",
						},
						cli.StringFlag{
							Name: "pubkey",
							Value:"",
							Usage: "account public key",
						},
					},
					Action: MigrateFlags(Cli.BcliNewAccount),
				},
				{
					Name: "get",
					Usage: "Getter account info",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "username",
							Value:"",
							Usage: "acocunt name",
						},
					},
					Action: MigrateFlags(Cli.BcliGetAccount),
				},
			},
		},
		{
			Name: "transfer",
			Usage: "transfer",
			Category: "transfer",
			Flags:[]cli.Flag {
				cli.StringFlag{
					Name: "form",
					Usage: "",
				},
				cli.StringFlag{
					Name: "to",
					Usage: "",
				},
				cli.StringFlag{
					Name: "amount",
					Usage: "",
				},
				cli.StringFlag{
					Name: "sign",
					Usage: "",
				},
			},
			Action: func(c *cli.Context) error {
				// TODO
				fmt.Println(c.String("form"))
				fmt.Println(c.String("to"))
				fmt.Println(c.String("amount"))
				fmt.Println(c.String("sign"))

				return nil
			},
		},
		{
			Name: "transaction",
			Usage: "transaction lists",
			Category: "transaction",
			Subcommands: []cli.Command {
				{
					Name: "get",
					Usage: "Getter tx details",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "trxhash",
						},
					},
					Action: func(c *cli.Context) error {

						if !isNotEmpty(c.String("trxhash")){
							return fmt.Errorf("参数错误")
						}

						// TODO
						fmt.Println(c.String("trxhash"))

						return nil
					},
				},
				{
					Name: "push",
					Usage: "push transaction",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "sender",
							Usage: "acocunt name",
						},
						cli.StringFlag{
							Name: "contract",
							Usage: "contract name",
						},
						cli.StringFlag{
							Name: "method",
							Usage: "method name",
						},
						cli.StringFlag{
							Name: "param",
							Usage: "param value",
						},
						cli.StringFlag{
							Name: "sign",
							Usage: "sign value",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(data.AccountInfo(c.String("sender")))
						fmt.Println(data.AccountInfo(c.String("contract")))
						fmt.Println(data.AccountInfo(c.String("method")))
						fmt.Println(data.AccountInfo(c.String("param")))
						fmt.Println(data.AccountInfo(c.String("sign")))
						return nil
					},
				},
			},
		},
		{
			Name: "contract",
			Usage: "contract info",
			Category: "contract",
			Subcommands: []cli.Command {
				{
					Name: "deploy",
					Usage: "contract deploy",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "name",
						},
						cli.StringFlag{
							Name: "code",
							Usage:"",
						},
						cli.StringFlag{
							Name: "abi",
							Usage:"",
						},
						cli.StringFlag{
							Name: "sign",
							Usage:"",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("name"))
						fmt.Println(c.String("code"))
						fmt.Println(c.String("abi"))
						fmt.Println(c.String("sign"))

						return nil
					},
				},
				{
					Name: "deploycode",
					Usage: "contract  deploycode",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "name",
						},
						cli.StringFlag{
							Name: "code",
							Usage:"",
						},
						cli.StringFlag{
							Name: "sign",
							Usage:"",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("name"))
						fmt.Println(c.String("code"))
						fmt.Println(c.String("sign"))

						return nil
					},
				},
				{
					Name: "deployabi",
					Usage: "contract  deployabi",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "name",
						},
						cli.StringFlag{
							Name: "abi",
							Usage:"",
						},
						cli.StringFlag{
							Name: "sign",
							Usage:"",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("name"))
						fmt.Println(c.String("abi"))
						fmt.Println(c.String("sign"))

						return nil
					},
				},
				{
					Name: "get",
					Usage: "Getter contract",
					Flags:[]cli.Flag {
						cli.StringFlag{
							Name: "name",
						},
						cli.StringFlag{
							Name: "code",
							Usage:"",
						},
						cli.StringFlag{
							Name: "abi",
							Usage:"",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("name"))
						fmt.Println(c.String("code"))
						fmt.Println(c.String("abi"))

						return nil
					},
				},
			},
		},
		{
			Name:     "p2p",
			Category: "p2p",
			Subcommands: []cli.Command{
				{
					Name:  "connect",
					Usage: "connect address or port",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "peer",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("peer"))
						return nil
					},
				},
				{
					Name:  "disconnect",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "peer",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("peer"))
						return nil
					},
				},
				{
					Name:  "status",
					Usage: "p2p status",
					Action: func(c *cli.Context) error {
						// TODO

						return nil
					},
				},
				{
					Name:  "peers",
					Usage: "peers info",
					Action: func(c *cli.Context) error {
						// TODO

						return nil
					},
				},
			},
		},
		{
			Name: "delegate",
			Category: "delegate",
			Subcommands: []cli.Command{
				{
					Name:  "reg",
					Usage: "connect address or port",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "account",
							Usage:"account name",
						},
						cli.StringFlag{
							Name: "signkey",
							Usage:"sign key",
						},
						cli.StringFlag{
							Name: "url",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("account"))
						fmt.Println(c.String("signkey"))
						fmt.Println(c.String("url"))
						return nil
					},
				},
				{
					Name:  "unreg",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "account",
						},
					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("account"))
						return nil
					},
				},
				{
					Name:  "list",
					Flags: []cli.Flag{
						cli.Int64Flag{
							Name: "limit",
							Value:100,
						},
						cli.Int64Flag{
							Name: "start",
							Value:0,
						},

					},
					Action: func(c *cli.Context) error {
						// TODO
						fmt.Println(c.String("limit"))
						fmt.Println(c.String("start"))
						return nil
					},
				},
			},
		},

	}


	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}
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
		return false, fmt.Errorf("参数错误！")
	}

	return true, nil
}
//
//func registerAccount(username string, pubkey string)  {
//	account := &user_proto.AccountInfo{
//		Name: username,
//		Pubkey: pubkey,
//	}
//	accountBuf, _ := pack.Marshal(account)
//
//	block, _:= data.BlockHeader()
//
//	txAccountSign := &push_sign.TransactionSign{
//		Version: 1,
//		CursorNum: block.HeadBlockNum,
//		CursorLabel: block.CursorLabel,
//		Lifetime: block.HeadBlockTime + 20,
//		Sender: "delta",
//		Contract: "bottos",
//		Method: "newaccount",
//		Param: accountBuf,
//		SigAlg: 1,
//	}
//
//	msg, _ := proto.Marshal(txAccountSign)
//	seckey,_ := hex.DecodeString("b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45")
//
//	chainID,_:=hex.DecodeString(config.CHAIN_ID)
//	msg = bytes.Join([][]byte{msg, chainID}, []byte{})
//	sign, _ := crypto.Sign(util.Sha256(msg), seckey)
//
//
//	txAccount := &bean.TxBean{
//		Version:     1,
//		CursorNum:   block.HeadBlockNum,
//		CursorLabel: block.CursorLabel,
//		Lifetime:    block.HeadBlockTime + 20,
//		Sender:      "delta",
//		Contract:    "bottos",
//		Method:      "newaccount",
//		Param:       hex.EncodeToString(accountBuf),
//		SigAlg:      1,
//		Signature:   hex.EncodeToString(sign),
//	}
//	data.PushTransaction(txAccount)
//}

