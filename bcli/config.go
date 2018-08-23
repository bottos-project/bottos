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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	
	"github.com/bottos-project/bottos/common"
)

const (
	//CONFIG_FILE_NAME configure file
	CONFIG_FILE_NAME = "./cliconfig.json"
)

//CONFIG configure pointer
var CONFIG *CLIConfig

//CLIConfig configure key pairs
type CLIConfig struct {
	KeyPairs []KeyPair `json:"key_pairs"`
}

//KeyPair key pair
type KeyPair struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

//LoadConfig read configure
func LoadConfig() error {
	file, e := loadConfigJson(CONFIG_FILE_NAME)
	if e != nil {
		fmt.Println("Read config file error: ", e)
		return e
	}

	config := CLIConfig{}
	e = json.Unmarshal(file, &config)
	if e != nil {
		fmt.Println("Unmarshal config file error: ", e)
		return e
	}
	CONFIG = &config
	return nil
}

//GetPrivateKey get private key
func GetPrivateKey(pubkey string) ([]byte, error) {
	if CONFIG != nil {
		for _, keypair := range CONFIG.KeyPairs {
			if pubkey == keypair.PublicKey {
				return common.HexStringToBytes(keypair.PrivateKey), nil
			}
		}
	}

	return []byte{}, fmt.Errorf("Key Not Found")
}

//GetDefaultKey get default private key
func GetDefaultKey() ([]byte, error) {
	if CONFIG != nil {
		return common.HexStringToBytes(CONFIG.KeyPairs[0].PrivateKey), nil
	}

	return []byte{}, fmt.Errorf("Key Not Found")
}


//GetChainId get chain id
func GetChainId() ([]byte, error) {
	if CONFIG != nil {
		infourl := "http://" + ChainAddr + "/v1/block/height"
                cli := &CLI{}
		chainInfo, err := cli.GetChainInfoOverHttp(infourl)
		if err != nil {
			return []byte("00000000000000000000000000000000"), nil
		}
		return common.HexStringToBytes(chainInfo.ChainId), nil
	}

	return []byte{}, fmt.Errorf("Chain Id Not Found")
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
