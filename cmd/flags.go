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

/*
 * file description:  bottos command line
 * @Author: Zhang Lei
 * @Date:   2018-07-27
 * @Last Modified by:
 * @Last Modified time:
 */
package cmd

import (
	"gopkg.in/urfave/cli.v1"
	)

var (
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Config file path",
	}
	GenesisFileFlag = cli.StringFlag{
		Name:  "genesis",
		Usage: "Genesis config file path",
	}
	DataDirFlag = cli.StringFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases",
	}
	DisableRESTFlag = cli.BoolFlag{
		Name:  "disable-api",
		Usage: "Disable RESTful API server",
	}
	RESTPortFlag = cli.IntFlag{
		Name:  "restport",
		Usage: "RESTful API server listening port",
		Value: 8689,
	}
	RESTServerAddrFlag = cli.StringFlag{
		Name:  "rest-servaddr",
		Usage: "RESTful API server address",
	}
	RestMaxLimit = cli.IntFlag{
		Name:  "rest_max_limit",
		Usage: "tps of rest",
		Value: 1000,
	}
	WalletRestMaxLimit = cli.IntFlag{
		Name:  "wallet_rest_max_limit",
		Usage: "tps of wallet rest",
		Value: 10,
	}

	EnableRPCFlag = cli.BoolFlag{
		Name:  "enable-rpc",
		Usage: "Enable RPC server",
	}
	RPCPortFlag = cli.IntFlag{
		Name:  "rpcport",
		Usage: "RPC server listening port",
		Value: 8690,
	}
	P2PPortFlag = cli.IntFlag{
		Name:  "p2pport",
		Usage: "P2P network listening port",
		Value: 9868,
	}
	P2PServerAddrFlag = cli.StringFlag{
		Name:  "p2p-servaddr",
		Usage: "P2P network server address",
	}
	PeerListFlag = cli.StringFlag{
		Name:  "peerlist",
		Usage: "Comma separated list of network peers (192.168.1.2:9868,192.168.1.3:9868,192.168.1.4:9868)",
	}
	DelegateSignkeyFlag = cli.StringFlag{
		Name:  "delegate-signkey",
		Usage: "Sign key for delegate ('key:<public key>,<private key>' or 'wallet:<wallet url>')",
	}
	DelegateFlag = cli.StringFlag{
		Name:  "delegate",
		Usage: "Producer account name",
	}
	DelegatePrateFlag = cli.IntFlag{
		Name:  "delegate-prate",
		Usage: "Config delegate participate threshold",
		Value: 33,
	}
	EnableMongoDBFlag = cli.BoolFlag{
		Name:  "enable-mongodb",
		Usage: "Enable mongodb plugin",
	}
	MongoDBFlag = cli.StringFlag{
		Name:  "mongodb",
		Usage: "MongoDB connection config",
	}
	LogConfigFlag = cli.StringFlag{
		Name:  "logconfig",
		Usage: "Log config file path",
	}
	WalletDirFlag = cli.StringFlag{
		Name:  "walletdir",
		Usage: "wallet directory",
	}
	EnableWalletFlag = cli.BoolFlag{
		Name:  "enable-wallet",
		Usage: "enable wallet",
	}
	WalletRESTPortFlag = cli.IntFlag{
		Name:  "wallet-rest-port",
		Usage: "Wallet RESTful API server listening port",
		Value: 6869,
	}
	WalletRESTServerAddrFlag = cli.StringFlag{
		Name:  "wallet-rest-servaddr",
		Usage: "Wallet RESTful API server address",
		Value: "localhost",
	}
	DebugFlag = cli.BoolFlag{
		Name:  "debug",
		Usage: "Enable debug mode",
	}
	LogMinLevelFlag = cli.StringFlag{
		Name:  "log-minlevel",
		Usage: "log minlevel",
		Value: "error",
	}
	LogMaxLevelFlag = cli.StringFlag{
		Name:  "log-maxlevel",
		Usage: "log maxlevel",
		Value: "critical",
	}
	LogLevelsFlag = cli.StringFlag{
		Name:  "log-levels",
		Usage: "log levels",
		Value: "debug,info,warn,error,critical",
	}
	LogMaxrollsFlag = cli.StringFlag{
		Name:  "log-maxrolls",
		Usage: "log maxrolls",
		Value: "999",
	}
	RecoverAtBlockNumFlag = cli.IntFlag{
		Name:  "recover_at_blocknum",
		Usage: "recover at blocknum",
		Value: 0,
	}
)
