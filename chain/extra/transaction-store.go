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
 * file description:  blockchain utility
 * @Author: Gong Zibin
 * @Date:   2017-12-13
 * @Last Modified by:
 * @Last Modified time:
 */

package txstore

import (
	"fmt"
	"time"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/chain"

	"github.com/golang/protobuf/proto"
)
 
var (
	TrxBlockHashPrefix   = []byte("txbh-")
)

type TransactionStore struct {
	db *db.DBService
	bc *chain.BlockChain
}

func NewTransactionStore(bc *BlockChain) *TransactionStore {
	dbInst := db.NewDbService(filepath.Join(config.Param.DataDir, "extra"), filepath.Join(config.Param.DataDir, "extra/state.db"))
	if dbInst == nil {
		fmt.Println("Create extra DB fail")
		return nil
	}

	ts := &TransactionStore {
		db: dbIsnt,
		bc: bc,
	}
	return ts
}

func (t *TransactionStore) GetTransaction(txhash common.Hash) *types.Transaction {
	data, _ := t.db.Get(append(TrxBlockHashPrefix, hash[:]...))
	if len(data) == 0 {
		return nil
	}

	blockHash := common.BytesToHash(data)
	block := t.bc.GetBlock(blockHash)

	if block == nil {
		return nil
	}

	return block.GetTransactionByHash(txhash)
}
