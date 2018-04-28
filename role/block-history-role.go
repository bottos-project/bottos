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
 * file description:  block history role
 * @Author: Gong Zibin
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"encoding/json"
	_"fmt"
	"strconv"

	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/common"
)

const BlockHistoryObjectName string = "block_history"

type BlockHistory struct {
	BlockNumber			uint32			`json:"block_number"`
	BlockHash			common.Hash		`json:"block_hash"`
}

func blockNumberToKey(blockNumber uint32) string {
	id := blockNumber & 0xFFFF
	key := strconv.Itoa(int(id))
	return key
}

func CreateBlockHistoryRole(ldb *db.DBService) error {
	for i := 0; i < 65536; i++ {
		value := &BlockHistory{}
		jsonvalue, err := json.Marshal(value)
		if err != nil {
			return err
		}

		err = ldb.SetObject(BlockHistoryObjectName, blockNumberToKey(uint32(i)), string(jsonvalue))
		if err != nil {
			return err
		}
	}

	return nil
}

func SetBlockHistoryRole(ldb *db.DBService, blockNumber uint32, blockHash common.Hash) error {
	key := blockNumberToKey(blockNumber)
	value := &BlockHistory {
		BlockNumber: blockNumber,
		BlockHash: blockHash,
	}
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ldb.SetObject(BlockHistoryObjectName, key, string(jsonvalue))
}

func GetBlockHistoryRole(ldb *db.DBService, blockNumber uint32) (*BlockHistory, error) {
	key := blockNumberToKey(blockNumber)
	value, err := ldb.GetObject(BlockHistoryObjectName, key)
	if err != nil {
		return nil, err
	}
	res := &BlockHistory{}
	err = json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
  