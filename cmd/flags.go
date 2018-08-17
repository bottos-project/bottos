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
	cli "gopkg.in/urfave/cli.v1"
)

var (
	ConfigFileFlag = cli.StringFlag{
		Name: "config",
		Usage: "config file path the greeting,If without this path, the bottos process will boot up with default config in hardcode",
	}
	GenesisFileFlag = cli.StringFlag{
		Name: "genesis",
		Usage: "genesis config file path the greeting",
	}
	DataDirFlag = cli.StringFlag{
		Name: "datadir",
		Usage: "datadir's path",
	}
	DisableAPIFlag = cli.BoolFlag{
		Name: "disable-api",
		Usage: "disable restful api's requests",
	}
	APIPortFlag = cli.IntFlag{
		Name: "apiport",
		Usage: "api service port for the greeting",
	}
	DisableRPCFlag = cli.BoolFlag{
		Name: "disable-rpc",
		Usage: "disable rpc requests",
	}
	RPCPortFlag = cli.IntFlag{
		Name: "rpcport",
		Usage: "json-rpc port for the greeting",
	}
	P2PPortFlag = cli.IntFlag{
		Name: "p2pport",
		Usage: "local listen on this p2p port to receive remote p2p messages",
	}
	ServerAddrFlag = cli.StringFlag{
		Name: "servaddr",
		Usage: "for p2p sync / reply local server ip& port info",
	}
	PeerListFlag = cli.StringFlag{
		Name: "peerlist",
		Usage: "for p2p add pne / add neighbour. Example: 192.168.1.2:9868, 192.168.1.3:9868, 192.168.1.4:9868",
	}
	DelegateSignkeyFlag = cli.StringFlag{
		Name: "delegate-signkey",
		Usage: "--delegate-signkey=<pubkey>,<private key>.Param struct needs be modified ,public and private key for native contract, external contracts' accounts",
	}
	DelegateFlag = cli.StringFlag{
		Name: "delegate",
		Usage: "Assign one producer. Later this section will no more be used.\n Only one delegate is allowed in one node(other than bottos account).",
	}
	EnableStaleReportFlag = cli.BoolFlag{
		Name: "enable-stale-report",
		Usage: "",
	}
	EnableMongoDBFlag = cli.BoolFlag{
		Name: "enable-mongodb",
		Usage: "",
	}
	MongoDBFlag = cli.StringFlag{
		Name: "mongodb",
		Usage: "db inst for load mongodb",
	}
	LogConfigFlag = cli.StringFlag{
		Name: "logconfig",
		Usage: "for seelog config",
	}
)
