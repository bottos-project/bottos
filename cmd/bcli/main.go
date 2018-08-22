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
	"fmt"
	"os"
	"gopkg.in/urfave/cli.v1"
	"github.com/micro/go-micro"
	chain "github.com/bottos-project/bottos/api"
)

var (
	app = cli.NewApp()
	client chain.ChainService

	getInfoCommand = cli.Command {
		Name: "getinfo",
		Usage: "Get chian info",
		Category: "general",
		Action: MigrateFlags(BcliGetChainInfo),
	}

	getBlockCommand = cli.Command {
		Name: "getblock",
		Usage: "Get block info",
		Category: "general",
		Flags: []cli.Flag {
			blockNumberFlag,
			blockHashFlag,
		},
		Action: MigrateFlags(BcliGetBlockInfo),
	}

	getTableCommand = cli.Command {
		Name: "gettable",
		Usage: "get table info",
		Category: "general",
		Flags: []cli.Flag {
			contractNameFlag,
			tableNameFlag,
			tableKeyNameFlag,
		},
		Action: MigrateFlags(BCLIGetTableInfo),
	}

	accountCommand = cli.Command {
		Name: "account",
		Usage: "Create or Get account",
		Category: "account",
		Subcommands: []cli.Command {
			{
				Name: "create",
				Usage: "Create account",
				Flags:[]cli.Flag {
					accountNameFlag,
					publicKeyFlag,
				},
				Action: MigrateFlags(BcliNewAccount),
			},
			{
				Name: "get",
				Usage: "Getter account info",
				Flags:[]cli.Flag {
					accountNameFlag,
				},
				Action: MigrateFlags(BcliGetAccount),
			},
		},
	}

	transferCommand = cli.Command {
		Name: "transfer",
		Usage: "transfer",
		Category: "transfer",
		Flags:[]cli.Flag {
			transferFromFlag,
			transferToFlag,
			transferAmountFlag,
		},
		Action: MigrateFlags(BcliTransfer),
	}

	transactionCommand = cli.Command {
		Name: "transaction",
		Usage: "transaction lists",
		Category: "transaction",
		Subcommands: []cli.Command {
			{
				Name: "get",
				Usage: "Getter tx details",
				Flags:[]cli.Flag {
					transactionHashFlag,
				},
				Action: MigrateFlags(BCLIGetTransaction),
			},
			{
				Name: "push",
				Usage: "push transaction",
				Flags:[]cli.Flag {
					transactionSenderFlag,
					contractNameFlag,
					transactionMethodFlag,
					transactionParamFlag,
				},
				Action: MigrateFlags(BCLIPushTransaction),
			},
		},
	}

	contractCommand = cli.Command {
		Name: "contract",
		Usage: "contract info",
		Category: "contract",
		Subcommands: []cli.Command {
			{
				Name: "deploy",
				Usage: "contract deploy",
				Flags:[]cli.Flag {
					contractNameFlag,
					contractCodeFlag,
					contractAbiFlag,
				},
				Action: MigrateFlags(BCLIDeployBoth),
			},
			{
				Name: "deploycode",
				Usage: "contract  deploycode",
				Flags:[]cli.Flag {
					contractNameFlag,
					contractCodeFlag,
				},
				Action: MigrateFlags(BCLIDeployCode),
			},
			{
				Name: "deployabi",
				Usage: "contract  deployabi",
				Flags:[]cli.Flag {
					contractNameFlag,
					contractAbiFlag,
				},
				Action: MigrateFlags(BCLIDeployAbi),
			},
			{
				Name: "get",
				Usage: "Getter contract",
				Flags:[]cli.Flag {
					contractNameFlag,
					contractCodeFlag,
					contractAbiFlag,
				},
				Action: MigrateFlags(BCLIGetContractCode),
			},
		},
	}

	p2pCommand = cli.Command {
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
	}

	delegateCommand = cli.Command {
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
	}
)

func init() {
	app.Name = "Bottos Cmd"
	app.Usage = "block chain bcli"
	app.Version = "0.0.1"
	app.Commands = []cli.Command {
		getInfoCommand,
		getBlockCommand,
		getTableCommand,
		accountCommand,
		transferCommand,
		transactionCommand,
		contractCommand,
		p2pCommand,
		delegateCommand,
	}
	app.Flags = []cli.Flag {
		globalServAddrFlag,
		globalSignerFlag,
	}
	app.Before = func(ctx *cli.Context) error {
		service := micro.NewService()
		client = chain.NewChainService("bottos", service.Client())
		return nil
	}
}

func main() {
	if err := LoadConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
