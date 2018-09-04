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
		Name: "config",
		Usage: "Config file path",
	}
	GenesisFileFlag = cli.StringFlag{
		Name: "genesis",
		Usage: "Genesis config file path",
	}
	DataDirFlag = cli.StringFlag{
		Name: "datadir",
		Usage: "Data directory for the databases",
	}
	DisableRESTFlag = cli.BoolFlag{
		Name: "disable-api",
		Usage: "Disable RESTful API server",
	}
	RESTPortFlag = cli.IntFlag{
		Name: "restport",
		Usage: "RESTful API server listening port",
		Value: 8689,
	}
	RESTServerAddrFlag = cli.StringFlag{
		Name: "rest-servaddr",
		Usage: "RESTful API server address",
	}
	DisableRPCFlag = cli.BoolFlag{
		Name: "disable-rpc",
		Usage: "Disable RPC server",
	}
	RPCPortFlag = cli.IntFlag{
		Name: "rpcport",
		Usage: "RPC server listening port",
		Value: 8690,
	}
	P2PPortFlag = cli.IntFlag{
		Name: "p2pport",
		Usage: "P2P network listening port",
		Value: 9868,
	}
	P2PServerAddrFlag = cli.StringFlag{
		Name: "p2p-servaddr",
		Usage: "P2P network server address",
	}
	PeerListFlag = cli.StringFlag{
		Name: "peerlist",
		Usage: "Comma separated list of network peers (192.168.1.2:9868,192.168.1.3:9868,192.168.1.4:9868)",
	}
	DelegateSignkeyFlag = cli.StringFlag{
		Name: "delegate-signkey",
		Usage: "Sign key for delegate (<public key>,<private key>)",
	}
	DelegateFlag = cli.StringFlag{
		Name: "delegate",
		Usage: "Producer account name",
	}
	EnableStaleReportFlag = cli.BoolFlag{
		Name: "enable-stale-report",
		Usage: "Enable stale block production",
	}
	EnableMongoDBFlag = cli.BoolFlag{
		Name: "enable-mongodb",
		Usage: "Enable mongodb plugin",
	}
	MongoDBFlag = cli.StringFlag{
		Name: "mongodb",
		Usage: "MongoDB connection config",
	}
	LogConfigFlag = cli.StringFlag{
		Name: "logconfig",
		Usage: "Log config file path",
	}
	WalletDirFlag = cli.StringFlag{
		Name: "walletdir",
		Usage: "wallet directory",
	}
	EnableWalletFlag = cli.StringFlag{
		Name: "enable-wallet",
		Usage: "enable wallet",
	}
)
