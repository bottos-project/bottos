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
	log "github.com/cihub/seelog"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/hashicorp/golang-lru"
)

const (
	//BlockChainCacheLimit cache limit
	BlockChainCacheLimit = 256
)

//BlockChainCache the cache of chain block
type BlockChainCache struct {
	headBlock  *types.Block
	startBlock *types.Block
	cache      *lru.Cache
}

//CreateBlockChainCache new a cache
func CreateBlockChainCache() (*BlockChainCache, error) {
	bcCache := BlockChainCache{
		headBlock:  nil,
		startBlock: nil,
	}
	bcCache.cache, _ = lru.New(BlockChainCacheLimit)

	return &bcCache, nil
}

//Reset purge cache
func (c *BlockChainCache) Reset() {
	c.cache.Purge()
}

//GetBlock get block by hash
func (c *BlockChainCache) GetBlock(hash common.Hash) *types.Block {
	if block, ok := c.cache.Get(hash); ok {
		return block.(*types.Block)
	}

	return nil
}

//GetBlockByNum get block by block number
func (c *BlockChainCache) GetBlockByNum(number uint64) *types.Block {
	if hash, ok := c.cache.Get(number); ok {
		return c.GetBlock(hash.(common.Hash))
	}

	return nil
}

//HasBlock check block
func (c *BlockChainCache) HasBlock(hash common.Hash) bool {
	if c.cache.Contains(hash) {
		return true
	}

	return false
}

//Trim remove block from cache
func (c *BlockChainCache) Trim(headBlockNum uint64, LIB uint64) error {
	if c.startBlock == nil {
		return nil
	}

	trimmed := false
	for start := c.startBlock.GetNumber(); start < LIB; start++ {
		if hash, ok := c.cache.Get(start); ok {
			c.cache.Remove(hash.(common.Hash))
			c.cache.Remove(start)
			trimmed = true
			log.Infof("remove block form block cache, num = %v, hash = %x\n", start, hash.(common.Hash))
		}
	}

	if trimmed {
		c.startBlock = c.GetBlockByNum(LIB)
	}

	log.Infof("BlockCache Trim, head block num = %v, LIB = %v, start = %v\n", headBlockNum, LIB, c.startBlock.GetNumber())
	return nil
}

func (c *BlockChainCache) add(block *types.Block) {
	hash := block.Hash()
	c.cache.Add(block.GetNumber(), hash)
	c.cache.Add(hash, block)
}

//Insert insert a block
func (c *BlockChainCache) Insert(block *types.Block) (*types.Block, error) {
	if c.headBlock == nil {
		c.add(block)

		c.headBlock = block
		c.startBlock = block

		return c.headBlock, nil
	}

	if block.GetPrevBlockHash() == c.headBlock.Hash() {
		c.add(block)

		c.headBlock = block

		return c.headBlock, nil
	}

	return nil, log.Errorf("BlockCache insert block error, block not link, block PrevBlockHash: %x, head block Hash: %x\n", block.GetPrevBlockHash(), c.headBlock.Hash())

}
