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
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/chain"
	"github.com/bottos-project/core/role"
)
 
var (
	TrxBlockHashPrefix   = []byte("txbh-")
)

type TransactionStore struct {
	db *db.DBService
	bc chain.BlockChainInterface
}

func NewTransactionStore(bc chain.BlockChainInterface, db *db.DBService) *TransactionStore {
	ts := &TransactionStore {
		db: db,
		bc: bc,
	}
	bc.RegisterHandledBlockCallback(ts.ReceiveHandledBlock)
	return ts
}

func (t *TransactionStore) GetTransaction(txhash common.Hash) *types.Transaction {
	blockHash, err := role.GetBlockHashByTxHash(t.db, txhash)
	if err != nil {
		return nil
	}

	block := t.bc.GetBlockByHash(blockHash)
	if block == nil {
		return nil
	}

	return block.GetTransactionByHash(txhash)
}

func (t *TransactionStore) addTx(txhash common.Hash, blockhash common.Hash) error {
	return role.AddTransactionHistory(t.db, txhash, blockhash)
}

func (t *TransactionStore) delTx(txhash common.Hash) error {
	key := append(TrxBlockHashPrefix, txhash.Bytes()...)
	err := t.db.Delete(key)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionStore) ReceiveHandledBlock(block *types.Block) {
	blockHash := block.Hash()

	for _, tx := range block.Transactions {
		txHash := tx.Hash()
		t.addTx(txHash, blockHash)
	}
}

func (t *TransactionStore) RemoveBlock(block *types.Block) {
	for _, tx := range block.Transactions {
		txHash := tx.Hash()
		t.delTx(txHash)
	}
}
