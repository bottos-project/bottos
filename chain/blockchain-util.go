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
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package chain

import (
	//"fmt"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/db"

	"github.com/golang/protobuf/proto"
)

var (
	//BlockHashPrefix prefix of block hash
	BlockHashPrefix = []byte("bh-")
	//BlockNumberPrefix prefix of block number
	BlockNumberPrefix = []byte("bn-")
	//LastBlockKey prefix of block key
	LastBlockKey = []byte("lb")
)

//HasBlock check block in db
func HasBlock(db *db.DBService, hash common.Hash) bool {
	data, _ := db.Get(append(BlockHashPrefix, hash[:]...))
	if len(data) != 0 {
		return true
	}

	return false
}

//GetBlock get block from db by hash
func GetBlock(db *db.DBService, hash common.Hash) *types.Block {
	data, _ := db.Get(append(BlockHashPrefix, hash[:]...))
	if len(data) == 0 {
		return nil
	}

	block := types.Block{}
	if err := proto.Unmarshal(data, &block); err != nil {
		return nil
	}

	//fmt.Printf("GetBlock, hash: %x, data: %x\n", hash.Bytes(), data)

	return &block
}

//GetBlockHashByNumber get block from db by number
func GetBlockHashByNumber(db *db.DBService, number uint32) common.Hash {
	hash, _ := db.Get(append(BlockNumberPrefix, common.NumberToBytes(number, 32)...))
	if len(hash) == 0 {
		return common.Hash{}
	}
	return common.BytesToHash(hash)
}

//GetLastBlock get lastest block from db
func GetLastBlock(db *db.DBService) *types.Block {
	data, _ := db.Get(LastBlockKey)
	if len(data) == 0 {
		return nil
	}

	return GetBlock(db, common.BytesToHash(data))
}

//WriteGenesisBlock write the first block in db
func WriteGenesisBlock(db *db.DBService, block *types.Block) error {
	if err := WriteBlock(db, block); err != nil {
		return err
	}

	return nil
}

func writeHead(db *db.DBService, block *types.Block) error {
	key := append(BlockNumberPrefix, common.NumberToBytes(block.GetNumber(), 32)...)
	err := db.Put(key, block.Hash().Bytes())
	if err != nil {
		return err
	}

	err = db.Put(LastBlockKey, block.Hash().Bytes())
	if err != nil {
		return err
	}
	return nil
}

//WriteBlock write block in db
func WriteBlock(db *db.DBService, block *types.Block) error {
	key := append(BlockHashPrefix, block.Hash().Bytes()...)
	data, _ := proto.Marshal(block)

	err := db.Put(key, data)
	if err != nil {
		return err
	}

	//fmt.Printf("WriteBlock, hash: %x, key: %x, value: %x\n", block.Hash().Bytes(), key, data)

	return writeHead(db, block)
}
