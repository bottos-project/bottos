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
	"fmt"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"os"
)

const (
	CONFIG_FILE_NAME = "./config.json"
)

var Param *Parameter
var Genesis *GenesisConfig

type Parameter struct {
	GenesisJson				string				`json:"genesis_json"`
	DataDir					string				`json:"data_dir"`
	Consensus       		string				`json:"consensus"`
	APIPort					int					`json:"api_port"`
	P2PPort					int					`json:"p2p_port"`
	PeerList				[]string			`json:"peer_list"`
	KeyPairs				[]KeyPair			`json:"key_pairs"`
}

type KeyPair struct {
	PrivateKey				string				`json:"private_key"`
	PublicKey				string				`json:"public_key"`
}

type GenesisConfig struct {
	GenesisTime				string				`json:"genesis_time"`
	ChainId					string				`json:"chain_id"`
	InitDelegate			InitDelegate		`json:"init_delegate"`
	
}

type InitDelegate struct {
	Name					string				`json:"name"`
	PublicKey				string				`json:"public_key"`
	Balance					uint32				`json:"balance"`
}

func init() {
	file, e := loadConfigJson(CONFIG_FILE_NAME)
	if e != nil {
		fmt.Println("Read config file error: ", e)
		os.Exit(1)
	}

	param := Parameter{}
	e = json.Unmarshal(file, &param)
	if e != nil {
		fmt.Println("Unmarshal config file error: ", e)
		os.Exit(1)
	}
	Param = &param

	file, e = loadConfigJson(param.GenesisJson)
	if e != nil {
		fmt.Println("Read genesis file error: ", e)
		os.Exit(1)
	}

	genesisConfig := GenesisConfig{}
	e = json.Unmarshal(file, &genesisConfig)
	if e != nil {
		fmt.Println("Unmarshal genesis file error: ", e)
		os.Exit(1)
	}
	Genesis = &genesisConfig

	fmt.Println(Param, Genesis)
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
