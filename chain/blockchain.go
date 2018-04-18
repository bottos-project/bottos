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
 * file description:  blockchain general interface and logic
 * @Author: Gong Zibin
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package chain

import (
	"fmt"
	"sync"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/config"
)

var chainInstance *BlockChain

type BlockChain struct {
	blockDb    *db.DBService
	blockCache *BlockChainCache
	//stateDb		Database

	genesisBlock *types.Block

	chainmu sync.RWMutex
}

func GetChain() *BlockChain {
	return chainInstance
}

func CreateBlockChain(blockDb *db.DBService) (*BlockChain, error) {
	blockCache, err := CreateBlockChainCache()
	if err != nil {
		return nil, err
	}

	bc := &BlockChain{
		blockDb:    blockDb,
		blockCache: blockCache,
		//stateDb:  stateDb,
	}

	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock == nil {
		var err error
		bc.genesisBlock, err = WriteGenesisBlock(blockDb)
		if err != nil {
			return nil, err
		}
	}

	err = bc.LoadBlockDb()
	if err != nil {
		return nil, err
	}

	// init block cache
	bc.initBlockCache()

	chainInstance = bc

	return bc, nil
}

func (bc *BlockChain) GetGenesisBlock() *types.Block {
	return bc.genesisBlock
}

func (bc *BlockChain) HasBlock(hash common.Hash) bool {
	if bc.blockCache.HasBlock(hash) {
		return true
	}

	return HasBlock(bc.blockDb, hash)
}

func (bc *BlockChain) GetBlock(hash common.Hash) *types.Block {
	// cache
	block := bc.blockCache.GetBlock(hash)
	if block != nil {
		return block
	}

	return GetBlock(bc.blockDb, hash)
}

func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {
	return bc.GetBlock(hash)
}

func (bc *BlockChain) GetBlockByNumber(number uint32) *types.Block {
	block := bc.blockCache.GetBlockByNum(number)
	if block != nil {
		return block
	}

	hash := GetBlockHashByNumber(bc.blockDb, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return bc.GetBlock(hash)
}

func (bc *BlockChain) GetBlockHashByNumber(number uint32) common.Hash {
	return GetBlockHashByNumber(bc.blockDb, number)
}

func (bc *BlockChain) WriteBlock(block *types.Block) error {
	return WriteBlock(bc.blockDb, block)
}

func (bc *BlockChain) HeadBlockTime() uint64 {
	//dpo := bc.stateDb.GetDynamicGlobalPropertyObject()
	//return dpo.Time
	return 1
}

func (bc *BlockChain) HeadBlockNum() uint32 {
	//dpo := bc.stateDb.GetDynamicGlobalPropertyObject()
	//return dpo.HeadBlockNum
	return 1
}

func (bc *BlockChain) HeadBlockHash() common.Hash {
	//dpo := bc.stateDb.GetDynamicGlobalPropertyObject()
	//return dpo.HeadBlockHash
	return common.Hash{}
}

func (bc *BlockChain) HeaderBlockProducer() []byte {
	//dpo := bc.stateDb.GetDynamicGlobalPropertyObject()
	//return dpo.CurrentProducer
	return []byte{}
}

// internal
func (bc *BlockChain) getBlockDbLastBlock() *types.Block {
	return GetLastBlock(bc.blockDb)
}

func (bc *BlockChain) initBlockCache() error {

	// TODO 鐤戦棶
	block := bc.getBlockDbLastBlock()
	if block != nil {
		_, err := bc.blockCache.Insert(block)
		return err
	}

	return nil
}

func (bc *BlockChain) LoadBlockDb() error {
	lastBlock := GetLastBlock(bc.blockDb)
	if lastBlock == nil {
		// TODO blockDb Recover()
		return fmt.Errorf("Loading block database fail, try recovering")
	}

	bc.blockCache.Insert(lastBlock)

	// TODO  for testing
	//bc.updateGlobalDynamicProperty(lastBlock)
	//dpo := bc.stateDb.GetDynamicGlobalPropertyObject()
	//dpo.LastIrreversibleBlockNum = lastBlock.Number()
	//bc.stateDb.SetDynamicGlobalPropertyObject(dpo)
	fmt.Printf("current block num = %v, hash = %x\n", lastBlock.GetNumber(), lastBlock.Hash())

	// TODO replay
	if bc.HeadBlockNum() < lastBlock.GetNumber() {
		// LoadAndExcuteBlocks()
	}

	return nil
}

// TODO
func (bc *BlockChain) updateGlobalProperty(block *types.Block) {
}

// TODO
func (bc *BlockChain) updateGlobalDynamicProperty(block *types.Block) {
	//dpo := bc.stateDb.GetDynamicGlobalPropertyObject()
	//dpo.HeadBlockNum = block.Number()
	//dpo.HeadBlockHash = block.Hash()
	//dpo.Time = block.Time()
	//dpo.CurrentProducer = block.Producer()
	//bc.stateDb.SetDynamicGlobalPropertyObject(dpo)
}

// TODO
func (bc *BlockChain) updateLastIrreversibleBlock(block *types.Block) {
	/*
		 // TODO  compute new LIB
		 dpo := bc.stateDb.GetDynamicGlobalPropertyObject()
		 // for test
		 if dpo.HeadBlockNum > dpo.LastIrreversibleBlockNum + 7 {
			 dpo.LastIrreversibleBlockNum = dpo.HeadBlockNum - 7
		 }
		 bc.stateDb.SetDynamicGlobalPropertyObject(dpo)

		 // write LIB to blockDb
		 newLIB := dpo.LastIrreversibleBlockNum
		 lastBlockNum := uint32(0)
		 lastBlock := bc.getBlockDbLastBlock()
		 if lastBlock != nil {
			 lastBlockNum = lastBlock.Number()
		 }

		 fmt.Printf("lastBlockNum = %v, newLIB = %v\n", lastBlockNum, newLIB)

		 if lastBlockNum < newLIB {
			 for i := lastBlockNum + 1; i <= newLIB; i++ {
				 block := bc.GetBlockByNumber(i)
				 if block != nil {
					 bc.WriteBlock(block)
				 } else {
					 fmt.Printf("block num = %v not found\n", i)
				 }
			 }
		 }

		 // trim blockCache
		 bc.blockCache.Trim(dpo.HeadBlockNum, newLIB)
	*/
}


func (bc *BlockChain) HandleBlock(block *types.Block) error {

	// TODO excute block
	fmt.Println("BlockChain : Handling block")

	// TODO
	//for _, tx := range block.Transactions {
		// ValidateTransaction(tx)
		// HandleTransaction(tx)
	//}

	// update global property
	bc.updateGlobalProperty(block)
	bc.updateGlobalDynamicProperty(block)
	bc.updateLastIrreversibleBlock(block)

	return nil
}

func (bc *BlockChain) ValidateBlock(block *types.Block) error {
	prevBlockHash := block.GetPrevBlockHash()
	if prevBlockHash != bc.HeadBlockHash() {
		return fmt.Errorf("Block Prev Hash error, head block Hash = %x, block PrevBlockHash = %x", bc.HeadBlockHash(), prevBlockHash)
	}

	if block.GetNumber() != bc.HeadBlockNum() {
		return fmt.Errorf("Block Number error, head block Number = %v, block Number = %v", bc.HeadBlockNum(), block.GetNumber())
	}

	// block timestamp check
	if block.GetTimestamp() <= bc.HeadBlockTime() {
		return fmt.Errorf("Block Timestamp error, head block time=%v, block time=%v", bc.HeadBlockTime(), block.GetTimestamp())
	}

	if block.GetTimestamp() > bc.HeadBlockTime() + uint64(config.DEFAULT_BLOCK_INTERVAL) {
		return fmt.Errorf("Block Timestamp error, head block time=%v, block time=%v", bc.HeadBlockTime(), block.GetTimestamp())
	}

	// TODO producer_change check
	// ...

	// TODO producer signature check
	//slot := store.GetSlotAtTime(block.Time())
	//producerName := store.GetScheduledProducer(slot)
	//scheduleProducerObj := bc.stateDb.GetProducerObject(producerName)
	if ok := block.ValidateSign(/*producer*/); !ok {
		return fmt.Errorf("Producer Sign Error")
	}

	// producer schedule check
	//blockProducer := block.Producer()
	//if string(blockProducer[:]) != scheduleProducerObj.Owner {
	//	return fmt.Errorf("Producer Producer Error")
	//}

	return nil
} 

func (bc *BlockChain) InsertBlock(block *types.Block) error {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	err := bc.ValidateBlock(block)
	if err != nil {
		fmt.Printf("Validate Block error %v\n", err)
		return err
	}

	// push to cache, block must link now, TODO: fork process
	_, err = bc.blockCache.Insert(block)
	if err != nil {
		fmt.Printf("blockCache insert error %v", err)
		return err
	}

	// TODO record stateDb revision
	err = bc.HandleBlock(block)
	if err != nil {
		// TODO restore stateDb revision
		fmt.Printf("InsertBlock error %v", err)
		return err
	}
	// TODO commit stateDb revision

	return nil
}

