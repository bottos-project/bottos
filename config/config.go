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
	CONFIG_FILE_NAME = "./config/config.json"
)

var Param *Parameter
var Genesis *GenesisConfig

type Parameter struct {
	GenesisJson				string				`json:"GenesisJson"`
	DataDir					string				`json:"DataDir"`
	Consensus       		string				`json:"Consensus"`
	APIPort					int					`json:"APIPort"`
	P2PPort					int					`json:"P2PPort"`
	PeerList				[]string			`json:"PeerList"`
}

type GenesisConfig struct {
	InitialTimestamp		string				`json:"InitialTimestamp"`
	InitAccounts			[]InitAccount		`json:"InitAccounts"`
	InitDelegates			[]InitDelegate		`json:"InitDelegates"`
	InitialChainId			string				`json:"InitialChainId"`
}

type InitAccount struct {
	Name					string				`json:"Name"`
	OwnerKey				string				`json:"OwnerKey"`
	ActiveKey				string				`json:"ActiveKey"`
	InitialBalance			uint32				`json:"InitialBalance"`
}

type InitDelegate struct {
	OwnerName 		    	string				`json:"OwnerName"`
	BlockSigningKey			string				`json:"BlockSigningKey"`
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
