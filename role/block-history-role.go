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
	"fmt"
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
		jsonvalue, _ := json.Marshal(value)
		ldb.SetObject(BlockHistoryObjectName, blockNumberToKey(uint32(i)), string(jsonvalue))
	}

	return nil
}

func SetBlockHistoryRole(ldb *db.DBService, blockNumber uint32, blockHash common.Hash) error {
	key := blockNumberToKey(blockNumber)
	value := &BlockHistory {
		BlockNumber: blockNumber,
		BlockHash: blockHash,
	}
	jsonvalue, _ := json.Marshal(value)
	return ldb.SetObject(BlockHistoryObjectName, key, string(jsonvalue))
}

func GetBlockHistoryByNumber(ldb *db.DBService, blockNumber uint32) (*BlockHistory, error) {
	key := blockNumberToKey(blockNumber)
	value, err := ldb.GetObject(BlockHistoryObjectName, key)
	res := &BlockHistory{}
	json.Unmarshal([]byte(value), res)
	fmt.Println("Get", key, value)
	return res, err
}
  