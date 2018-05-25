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
 * file description: database interface
 * @Author: May Luo
 * @Date:   2017-12-04
 * @Last Modified by:
 * @Last Modified time:
 */
package txstore

import (
	"fmt"
	"testing"
	"os"
	"io"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/db"
)

type MockBlockChain struct {
	cb chain.HandledBlockCallback
	block *types.Block
}

func NewMockBlockChain() chain.BlockChainInterface {
	bc := &MockBlockChain{}
	return bc
}

func (bc *MockBlockChain) InsertBlock(block *types.Block) error {
	bc.block = types.NewBlock(block.Header, block.Transactions)
	if bc.cb != nil {
		bc.cb(block)
	}
	return nil
}

func (bc *MockBlockChain) GetBlockByHash(hash common.Hash) *types.Block {
	if bc.block.Hash() == hash {
		return bc.block
	}

	return nil
}

func (bc *MockBlockChain) RegisterHandledBlockCallback(cb chain.HandledBlockCallback) {
	bc.cb = cb
}

func (bc *MockBlockChain) HasBlock(hash common.Hash) bool {return false}
func (bc *MockBlockChain) GetBlockByNumber(number uint32) *types.Block  {return nil}
func (bc *MockBlockChain) HeadBlockTime() uint64  {return 0}
func (bc *MockBlockChain) HeadBlockNum() uint32  {return 0}
func (bc *MockBlockChain) HeadBlockHash() common.Hash  {return common.Hash{}}
func (bc *MockBlockChain) HeadBlockDelegate() string  {return ""}
func (bc *MockBlockChain) GenesisTimestamp() uint64  {return 0}

func TestTxStore(t *testing.T) {
	bc := NewMockBlockChain()
	dbInst := db.NewDbService("./datadir", "./datadir/db.db")
	if dbInst == nil {
		fmt.Println("Create DB service fail")
		os.Exit(1)
	}
	txStore := NewTransactionStore(bc, dbInst)

	var txs []*types.Transaction
	tx1 := &types.Transaction{RefBlockNum:1}
	tx2 := &types.Transaction{RefBlockNum:2}
	tx3 := &types.Transaction{RefBlockNum:3}
	txs = append(txs, tx1)
	txs = append(txs, tx2)
	txs = append(txs, tx3)
	header := types.NewHeader()
	block := types.NewBlock(header, txs)
	bc.InsertBlock(block)

	// check
	expTx1 := txStore.GetTransaction(tx1.Hash())
	fmt.Printf("tx1.hash=%x, expTx1.hash=%x\n", tx1.Hash(), expTx1.Hash())
}

func CopyFile(dstName, srcName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        return
    }
    defer src.Close()

    dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}