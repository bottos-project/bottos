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
 * file description:  core state role
 * @Author: Gong Zibin
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"

	"github.com/bottos-project/bottos/db"
)

type CoreState struct {
	Config           ChainConfig `json:"chain_config"`
	CurrentDelegates []string    `json:"current_delegates"`
}

type ChainConfig struct {
	MaxBlockSize   uint32 `json:"max_block_size"`
	MaxTrxLifetime uint32 `json:"max_trx_lifetime"`
	MaxTrxRuntime  uint32 `json:"max_trx_runtime"`
	InDepthLeimit  uint32 `json:"in_depth_limit"`
}

const (
	CoreStateName       string = "core_state"
	CoreStateDefaultKey string = "core_state_defkey"
)

func CreateCoreStateRole(ldb *db.DBService) error {
	dgp := &CoreState{
		Config: ChainConfig{
			MaxBlockSize:   5242880,
			MaxTrxLifetime: 3600,
			MaxTrxRuntime:  10000,
			InDepthLeimit:  4,
		},
		CurrentDelegates: []string{},
	}
	return SetCoreStateRole(ldb, dgp)
}

func SetCoreStateRole(ldb *db.DBService, value *CoreState) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(CoreStateName, CoreStateDefaultKey, string(jsonvalue))
}

func GetCoreStateRole(ldb *db.DBService) (*CoreState, error) {
	value, err := ldb.GetObject(CoreStateName, CoreStateDefaultKey)
	if err != nil {
		return nil, err
	}

	res := &CoreState{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
