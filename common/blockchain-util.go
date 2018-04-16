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
 * @Date:   2017-12-01
 * @Last Modified by:
 * @Last Modified time:
 */

package common

import (
	"time"
	"fmt"
	"bytes"
	//"time"
	"github.com/bottos-project/core/common/types"
	proto "github.com/golang/protobuf/proto"
)

var (
	BlockHashPrefix  	= []byte("bh-")
	BlockNumberPrefix   = []byte("bn-")
	LastBlockKey		= []byte("lb")
)

func HasBlock(db Database, hash Hash) bool {
	data, _ := db.Get(append(BlockHashPrefix, hash[:]...))
	if len(data) != 0 {
		return true
	}

	return false
}

func GetBlock(db Database, hash Hash) *types.Block {
	data, _ := db.Get(append(BlockHashPrefix, hash[:]...))
	if len(data) == 0 {
		return nil
	}

	block := types.Block{}
	if err := proto.Unmarshal(data, &block); err != nil {
		return nil
	}

	fmt.Printf("GetBlock, hash: %x, data: %x\n", hash.Bytes(), data)

	return &block
}
func GetBlockHashByNumber(db Database, number uint32) Hash {
	hash, _ := db.Get(append(BlockNumberPrefix, NumberToBytes(number,32)...))
	if len(hash) == 0 {
		return Hash{}
	}
	return BytesToHash(hash)
}

func GetLastBlock(db Database) *types.Block {
	data, _ := db.Get(lastBlockPre)
	if len(data) == 0 {
		return nil
	}

	block := types.Block{}
	if err := proto.Unmarshal(data, &block); err != nil {
		return nil
	}

	return &block
}

func WriteGenesisBlock(blockDb Database) (*types.Block, error) {
	// TODO make block and write to db
	header := &types.Header {
		PrevBlockHash:	Hash{}[:],
		Number:			0,
		Timestamp:		uint64(time.Now().Unix()),
		MerkleRoot:     Hash{}[:],
		Producer:		[]byte{},
		ProducerChange:	[][]byte]{},
		ProducerSign:	Hash{}[:],
	}

	block := types.NewBlock(header, []*types.Transaction{})

	if err := WriteBlock(blockDb, block); err != nil {
		return nil, err
	}

	if err := WriteHead(blockDb, block); err != nil {
		return nil, err
	}
	
	return block, nil
}

func writeHead(db Database, block *types.Block) error {
	key := append(BlockNumberPrefix, NumberToBytes(block.GetNumber(),32)...)
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

func WriteBlock(db Database, block *types.Block) error {
	key := append(BlockHashPrefix, block.Hash().Bytes()...)
	data, _ := proto.Marshal(block)

	err := db.Put(key, data)
	if err != nil {
		return err
	}

	fmt.Printf("WriteBlock, hash: %x, key: %s, value: %x\n", block.Hash().Bytes(), string(key), data)

	return writeHead(block)
}
