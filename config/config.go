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
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/cihub/seelog"
)

const (
	// CONFIG_FILE_NAME is definition of config file name
	CONFIG_FILE_NAME = "./chainconfig.json"
)

// Param is var of Parameter type
var Param *Parameter

// Genesis is var of GenesisConfig type
var Genesis *GenesisConfig

// Parameter is definition of config param
type Parameter struct {
	GenesisJson       string    `json:"genesis_json"`
	DataDir           string    `json:"data_dir"`
	Consensus         string    `json:"consensus"`
	APIPort           int       `json:"api_port"`
	P2PPort           int       `json:"p2p_port"`
	ServAddr          string    `json:"serv_addr"`
	PeerList          []string  `json:"peer_list"`
	KeyPairs          []KeyPair `json:"key_pairs"`
	Delegates         []string  `json:"delegates"`
	ApiServiceEnable  bool      `json:"api_service_enable"`
	ApiServiceName    string    `json:"api_service_name"`
	ApiServiceVersion string    `json:"api_service_version"`
	EnableStaleReport bool      `json:"enable_stale_report"`
	OptionDb          string    `json:"option_db"`
	LogConfig         string    `json:"log_config"`
	ChainId           string    `json:"chain_id"`
}

// KeyPair is definition of key pair
type KeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

// GenesisConfig is definition of genesis config
type GenesisConfig struct {
	GenesisTime   uint64         `json:"genesis_time"`
	ChainId       string         `json:"chain_id"`
	InitDelegates []InitDelegate `json:"init_delegates"`
}

// InitDelegate is definition of init delegate
type InitDelegate struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	Balance   uint64 `json:"balance"`
}

// LoadConfig is to load config file
func LoadConfig() error {
	file, e := loadConfigJson(CONFIG_FILE_NAME)
	if e != nil {
		log.Error("Read config file error: ", e)
		return e
	}

	param := Parameter{}
	e = json.Unmarshal(file, &param)
	if e != nil {
		log.Error("Unmarshal config file error: ", e)
		return e
	}
	Param = &param

	file, e = loadConfigJson(param.GenesisJson)
	if e != nil {
		log.Error("Read genesis file error: ", e)
		return e
	}

	loadLogConfig(param.LogConfig)

	genesisConfig := GenesisConfig{}
	e = json.Unmarshal(file, &genesisConfig)
	if e != nil {
		log.Error("Unmarshal genesis file error: ", e)
		return e
	}
	Genesis = &genesisConfig

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

// loadLogConfig is to load log config file
func loadLogConfig(logConfigFile string) {
	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile(logConfigFile)
	if err != nil {
		log.Critical("*ERROR* Failed to parse config log file !!!", err)
		os.Exit(1)
		return
	}
	log.ReplaceLogger(logger)
}
