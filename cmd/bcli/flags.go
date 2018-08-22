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
	"gopkg.in/urfave/cli.v1"
)

var (

	globalServAddrFlag = cli.StringFlag{
		Name:  "servaddr",
		Value: "localhost:8689",
	}

	globalSignerFlag = cli.StringFlag{
		Name:  "signer",
		Value: "",
		Usage: "sign account name",
	}

	blockNumberFlag = cli.Uint64Flag{
		Name: "number",
		Value: 100,
		Usage: "get block by number",
	}

	blockHashFlag = cli.StringFlag{
		Name: "hash",
		Value: "",
		Usage: "get block by hash",
	}

	contractNameFlag = cli.StringFlag{
		Name: "contract",
		Value:"",
		Usage: "contract name",
	}

	tableNameFlag = cli.StringFlag{
		Name: "table",
		Usage: "table name",
	}

	tableKeyNameFlag = cli.StringFlag{
		Name: "key",
		Usage: "key value",
	}

	accountNameFlag = cli.StringFlag{
		Name: "name",
		Value: "",
		Usage: "acocunt name",
	}

	publicKeyFlag = cli.StringFlag{
		Name: "pubkey",
		Value: "",
		Usage: "account public key",
	}

	transferFromFlag = cli.StringFlag{
		Name: "from",
		Usage: "transfer from account name",
	}

	transferToFlag = cli.StringFlag{
		Name: "to",
		Usage: "transfer to account name",
	}

	transferAmountFlag = cli.StringFlag{
		Name: "amount",
		Usage: "transfer amount",
	}

	transactionHashFlag = cli.StringFlag{
		Name: "trxhash",
		Usage: "transaction hash",
	}

	transactionSenderFlag = cli.StringFlag{
		Name: "sender",
		Usage: "acocunt name",
	}

	transactionMethodFlag = cli.StringFlag{
		Name: "method",
		Usage: "method name",
	}

	transactionParamFlag = cli.StringFlag{
		Name: "param",
		Usage: "parameter hex string",
	}

	contractCodeFlag = cli.StringFlag{
		Name: "code",
		Usage: "",
	}

	contractAbiFlag = cli.StringFlag{
		Name: "abi",
		Usage: "",
	}
)
