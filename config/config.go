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
 * file description:  config load
 * @Author: Gong Zibin
 * @Date:   2017-12-11
 * @Last Modified by:
 * @Last Modified time:
 */

package config

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"strings"
	"github.com/bottos-project/bottos/cmd"
	cli "gopkg.in/urfave/cli.v1"
	log "github.com/cihub/seelog"
	"time"
)

const (
	// DEFAULT_CONFIG_FILENAME is definition of config file name
	DEFAULT_CONFIG_FILENAME = "./chainconfig.json"
)

// Param is var of Parameter type
var Param Parameter

// Genesis is var of GenesisConfig type
var Genesis GenesisConfig

// Parameter is definition of config param
type Parameter struct {
	GenesisJson       string    `json:"genesis_json"`
	DataDir           string    `json:"data_dir"`
	Consensus         string    `json:"consensus"`
	APIPort           int       `json:"api_port"`
	P2PPort           int       `json:"p2p_port"`
	ServAddr          string    `json:"serv_addr"`
	ServInterAddr	  string    `json:"serv_inter_addr"`
	PeerList          []string  `json:"peer_list"`
	KeyPairs          []KeyPair `json:"key_pairs"`
	Delegates         []string  `json:"delegates"`
	RpcServiceEnable  bool      `json:"rpc_service_enable"`
	RpcServiceName    string    `json:"rpc_service_name"`
	RpcServiceVersion string    `json:"rpc_service_version"`
	RestFulApiServiceEnable  bool      `json:"restful_api_service_enable"`
	EnableStaleReport bool      `json:"enable_stale_report"`
	OptionDb          string    `json:"option_db"`
	LogConfig         string    `json:"log_config"`
	ChainId           string    `json:"chain_id"`
	DelegateSignKey   KeyPair   `json:"delegate_signkey_pair"`
}

// KeyPair is definition of key pair
type KeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

// GenesisConfig is definition of genesis config
type GenesisConfig struct {
	GenesisTime   uint64         `json:"genesis_time"`
	InitDelegates []InitDelegate `json:"init_delegates"`
}

// InitDelegate is definition of init delegate
type InitDelegate struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	Balance   uint64 `json:"balance"`
}

func InitConfig() {
	Param.GenesisJson = "./genesis.json"
	Param.DataDir     = "./datadir/"
	Param.Consensus   = "dpos"
	Param.APIPort     = 8689
	Param.P2PPort     = 9868
	Param.ServAddr    = "192.168.1.1"
	Param.ServInterAddr = "127.0.0.1"
	Param.PeerList    = []string{}
	Param.KeyPairs    = []KeyPair {
		{
			PrivateKey: "b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45",
			PublicKey: "0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f",
		},
	}
	Param.Delegates   = []string{}
	Param.RpcServiceEnable = true
	Param.RpcServiceName   = "bottos"
	Param.RpcServiceVersion = "3.0.0"
	Param.EnableStaleReport = true
	Param.OptionDb          = ""
	Param.LogConfig         = "./corelog.xml"
	Param.ChainId           = "00000000000000000000000000000000"

	Genesis.GenesisTime    = 1524801531
	Genesis.InitDelegates  = []InitDelegate{}
}

func loadConfigFile(fn string) error {
	file, e := loadConfigJson(fn)
	if e != nil {
		return fmt.Errorf("Load config file error: ", e)
	}

	e = json.Unmarshal(file, &Param)
	if e != nil {
		return fmt.Errorf("Parse config file error: %v", e)
	}

	return nil
}

func loadGenesisFile(fn string) error {
	file, e := loadConfigJson(fn)
	if e != nil {
		return fmt.Errorf("Load genesis file error: ", e)
	}

	type GenesisStruct struct {
		GenesisTime   string         `json:"genesis_time"`
		InitDelegates []InitDelegate `json:"init_delegates"`
	}
	gs := GenesisStruct{}
	e = json.Unmarshal(file, &gs)
	if e != nil {
		return fmt.Errorf("Parse genesis file error: %v", e)
	}

	gtstr := gs.GenesisTime
	if !strings.HasSuffix(gtstr, "Z") {
		gtstr += "Z"
	}
	gt, e := time.Parse(time.RFC3339, gtstr)
	if e != nil {
		return fmt.Errorf("Parse genesis time error: %v", e)
	}
	Genesis.GenesisTime = uint64(gt.Unix())
	Genesis.InitDelegates = make([]InitDelegate, len(gs.InitDelegates))
	copy(Genesis.InitDelegates, gs.InitDelegates)

	return nil
}


// LoadConfig is to load config file
func LoadConfig(ctx *cli.Context) error {
	configFn := DEFAULT_CONFIG_FILENAME
	if ctx.GlobalIsSet(cmd.ConfigFileFlag.Name) {
		configFn = ctx.GlobalString(cmd.ConfigFileFlag.Name)
	}
	if err := loadConfigFile(configFn); err != nil {
		return err
	}

	genesisFn := Param.GenesisJson
	if ctx.GlobalIsSet(cmd.GenesisFileFlag.Name) {
		genesisFn = ctx.GlobalString(cmd.GenesisFileFlag.Name)
	}
	if err := loadGenesisFile(genesisFn); err != nil {
		return err
	}
	Param.GenesisJson = genesisFn

	if ctx.GlobalIsSet(cmd.DataDirFlag.Name) {
		Param.DataDir = ctx.GlobalString(cmd.DataDirFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.LogConfigFlag.Name) {
		Param.LogConfig = ctx.GlobalString(cmd.LogConfigFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.ServerAddrFlag.Name) {
		Param.ServAddr = ctx.GlobalString(cmd.ServerAddrFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.ServerAddrFlag.Name) {
		Param.ServAddr = ctx.GlobalString(cmd.ServerAddrFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.RestPortFlag.Name) {
		Param.APIPort = ctx.GlobalInt(cmd.RestPortFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.RPCPortFlag.Name) {
		//Param.RPCPort = ctx.GlobalInt(cmd.RPCPortFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.P2PPortFlag.Name) {
		Param.P2PPort = ctx.GlobalInt(cmd.P2PPortFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.DelegateFlag.Name) {
		d := ctx.GlobalString(cmd.DelegateFlag.Name)
		Param.Delegates = []string{d}
	}

	if ctx.GlobalIsSet(cmd.MongoDBFlag.Name) {
		Param.OptionDb = ctx.GlobalString(cmd.MongoDBFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.EnableStaleReportFlag.Name) {
		Param.EnableStaleReport = ctx.GlobalBool(cmd.EnableStaleReportFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.DisableRPCFlag.Name) {
		Param.RpcServiceEnable = ctx.GlobalBool(cmd.DisableRPCFlag.Name)
	}

	if ctx.GlobalIsSet(cmd.PeerListFlag.Name) {
		raw := ctx.GlobalString(cmd.PeerListFlag.Name)
		peerList, err := parsePeerListCLI(raw)
		if err != nil {
			return fmt.Errorf("parse peerlist error")
		}
		Param.PeerList = make([]string, len(peerList))
		copy(Param.PeerList, peerList)
	}

	if ctx.GlobalIsSet(cmd.DelegateSignkeyFlag.Name) {
		raw := ctx.GlobalString(cmd.DelegateSignkeyFlag.Name)
		keypair, err := parseKeyPairCLI(raw)
		if err != nil {
			return err
		}
		Param.DelegateSignKey = keypair
	}

	return nil
}

func parsePeerListCLI(raw string) ([]string, error) {
	var peerList []string
	val := strings.Replace(raw, " ", "", -1)
	peerList = strings.Split(val, ",")
	// check peers
	return peerList, nil
}

func parseKeyPairCLI(raw string) (KeyPair, error) {
	val := strings.Replace(raw, " ", "", -1)
	keys := strings.Split(val, ",")
	if len(keys) < 2 {
		return KeyPair{}, fmt.Errorf("parse delegate sign key error")
	}
	keypair := KeyPair{}
	keypair.PrivateKey = keys[0]
	keypair.PublicKey = keys[1]
	// check key
	return keypair, nil
}

// InitLogConfig initialize log config
func InitLogConfig(ctx *cli.Context) error {
	if ctx.GlobalIsSet(cmd.LogConfigFlag.Name) {
		Param.LogConfig = ctx.GlobalString(cmd.LogConfigFlag.Name)
	}

	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile(Param.LogConfig)
	if err != nil {
		return fmt.Errorf("parse log config file error: ", err)
	}
	log.ReplaceLogger(logger)
	return nil
}

func loadConfigJson(fn string) ([]byte, error) {
	file, e := ioutil.ReadFile(fn)
	if e != nil {
		return nil, e
	}

	// Remove the UTF-8 Byte Order Mark
	file = bytes.TrimPrefix(file, []byte("\xef\xbb\xbf"))
	return file, nil
}
