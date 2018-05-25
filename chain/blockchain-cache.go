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
 * file description:  general Hash type
 * @Author: Gong Zibin
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

 package chain

 import (
	"fmt"
	//"sync"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"

	"github.com/hashicorp/golang-lru"
 )
 
 const (
	 BlockChainCacheLimit = 256
 )
 
 type BlockChainCache struct {
	 //memdb		*db.MemDb
	 headBlock		 *types.Block
	 startBlock		 *types.Block
	 cache           *lru.Cache
 }
 
 func CreateBlockChainCache() (*BlockChainCache, error) {
	 /*
	 memdb, err := db.NewMemDatabase()
	 if err != nil {
		 return nil, err
	 }
	 */
 
	 bcCache := BlockChainCache {
		 //memdb: memdb,
		 headBlock: nil,
		 startBlock: nil,
	 }
	 bcCache.cache, _ = lru.New(BlockChainCacheLimit)
 
	 return &bcCache, nil
 }
 
 func (self *BlockChainCache) GetBlock(hash common.Hash) *types.Block {
	 if block, ok := self.cache.Get(hash); ok {
		 return block.(*types.Block)
	 } else {
		 return nil
	 }
 }
 
 func (self *BlockChainCache) GetBlockByNum(number uint32) *types.Block {
	 if hash, ok := self.cache.Get(number); ok {
		 return self.GetBlock(hash.(common.Hash))
	 } else {
		 return nil
	 }
 }
 
 func (self *BlockChainCache) HasBlock(hash common.Hash) bool {
	 if self.cache.Contains(hash) {
		 return true
	 }
	 return false
 }
 
 func (self *BlockChainCache) Trim(headBlockNum uint32, LIB uint32) error {
	 if self.startBlock == nil {
		 return nil
	 }

	 trimmed := false
	 for start := self.startBlock.GetNumber(); start < LIB; start++ {
		 if hash, ok := self.cache.Get(start); ok {
			 self.cache.Remove(hash.(common.Hash))
			 self.cache.Remove(start)
			 trimmed = true
			 fmt.Printf("remove block form block cache, num = %v, hash = %x\n", start, hash.(common.Hash))
		 }
	 }

	 if (trimmed) {
		self.startBlock = self.GetBlockByNum(LIB)
	 }
 
	 fmt.Printf("BlockCache Trim, head block num = %v, LIB = %v, start = %v\n", headBlockNum, LIB, self.startBlock.GetNumber())
	 return nil
 }
 
 func (self *BlockChainCache) add(block *types.Block) {
	 hash := block.Hash()
	 self.cache.Add(block.GetNumber(), hash)
	 self.cache.Add(hash, block)
 }
 
 func (self *BlockChainCache) Insert(block *types.Block) (*types.Block, error) {
	 if self.headBlock == nil {
		 self.add(block)
 
		 self.headBlock = block
		 self.startBlock = block
		 
		 return self.headBlock, nil
	 } else {
		 if block.GetPrevBlockHash() == self.headBlock.Hash() {
			 self.add(block)
 
			 self.headBlock = block
 
			 return self.headBlock, nil
		 } else {
			 return nil, fmt.Errorf("BlockCache insert block error, block not link, block PrevBlockHash: %x, head block Hash: %x\n", block.GetPrevBlockHash(), self.headBlock.Hash())
		 }
	 }
 }
 