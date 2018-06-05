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
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/role"
)

var (
	TrxBlockHashPrefix = []byte("txbh-")
)

type TransactionStore struct {
	roleIntf role.RoleInterface
	bc       chain.BlockChainInterface
}

func NewTransactionStore(bc chain.BlockChainInterface, roleIntf role.RoleInterface) *TransactionStore {
	ts := &TransactionStore{
		roleIntf: roleIntf,
		bc:       bc,
	}
	bc.RegisterHandledBlockCallback(ts.ReceiveHandledBlock)
	return ts
}

func (t *TransactionStore) GetTransaction(txhash common.Hash) *types.Transaction {
	blockHash, err := t.roleIntf.GetTransactionHistory(txhash)
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
	return t.roleIntf.AddTransactionHistory(txhash, blockhash)
}

func (t *TransactionStore) delTx(txhash common.Hash) error {
	// TODO

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
