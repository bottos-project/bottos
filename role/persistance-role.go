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
 * file description:  persistance role
 * @Author:
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package role

import (
	"errors"
	"fmt"
	"time"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/config"
	"github.com/bottos-project/core/db"
	"gopkg.in/mgo.v2/bson"
)

type AccountInfo struct {
	ID               bson.ObjectId `bson:"_id"`
	AccountName      string        `bson:"account_name"`
	Balance          string        `bson:"bto_balance"`
	StakedBalance    string        `bson:"staked_balance"`
	UnstakingBalance string        `bson:"unstaking_balance"`
	PublicKey        []byte        `bson:"public_key"`
	VMType           byte          `bson:"vm_type"`
	VMVersion        byte          `bson:"vm_version"`
	CodeVersion      common.Hash   `bson:"code_version"`
	CreateTime       time.Time     `bson:"create_time"`
	ContractCode     []byte        `bson:"contract_code"`
	ContractAbi      []byte        `bson:"abi"`
	UpdatedTime      time.Time     `bson:"updated_time"`
}
type BlockInfo struct {
	ID              bson.ObjectId   `bson:"_id"`
	BlockHash       common.Hash     `bson:"block_hash"`
	PrevBlockHash   []byte          `bson:"prev_block_hash"`
	BlockNumber     uint32          `bson:"block_number"`
	Timestamp       uint64          `bson:"timestamp"`
	MerkleRoot      []byte          `bson:"merkle_root"`
	DelegateAccount string          `bson:"delegate"`
	Transactions    []bson.ObjectId `bson:"transactions"`
	CreateTime      time.Time       `bson:"create_time"`
}

type TxInfo struct {
	ID            bson.ObjectId `bson:"_id"`
	BlockNum      uint32        `bson:"block_number"`
	TransactionID common.Hash   `bson:"transaction_id"`
	SequenceNum   uint32        `bson:"sequence_num"`
	BlockHash     common.Hash   `bson:"block_hash"`
	CursorNum     uint32        `bson:"cursor_num"`
	CursorLabel   uint32        `bson:"cursor_label"`
	Lifetime      uint64        `bson:"lifetime"`
	Sender        string        `bson:"sender"`
	Contract      string        `bson:"contract"`
	Method        string        `bson:"method"`
	Param         []byte        `bson:"param"`
	SigAlg        uint32        `bson:"sig_alg"`
	Signature     []byte        `bson:"signature"`
	CreateTime    time.Time     `bson:"create_time"`
}

func insertAccountInfoRole(ldb *db.DBService, block *types.Block, trx *types.Transaction) error {

	return nil
}

func insertTxInfoRole(ldb *db.DBService, block *types.Block, oids []bson.ObjectId) error {
	if len(oids) != len(block.Transactions) {
		return errors.New("invalid param")
	}

	for i, trx := range block.Transactions {
		newtrx := &TxInfo{
			ID:            oids[i],
			BlockNum:      block.Header.Number,
			TransactionID: trx.Hash(),
			SequenceNum:   uint32(i),
			BlockHash:     block.Hash(),
			CursorNum:     trx.CursorNum,
			CursorLabel:   trx.CursorLabel,
			Lifetime:      trx.Lifetime,
			Sender:        trx.Sender,
			Contract:      trx.Contract,
			Method:        trx.Method,
			Param:         trx.Param,
			SigAlg:        trx.SigAlg,
			Signature:     trx.Signature,
			CreateTime:    time.Now(),
		}
		ldb.Insert("Transactions", newtrx)
		if trx.Contract == config.BOTTOS_CONTRACT_NAME {
			insertAccountInfoRole(ldb, block, trx)
		}
	}

	return nil
}
func insertBlockInfoRole(ldb *db.DBService, block *types.Block, oids []bson.ObjectId) error {

	newBlockInfo := &BlockInfo{
		bson.NewObjectId(),
		block.Hash(),
		block.Header.PrevBlockHash,
		block.Header.Number,
		block.Header.Timestamp,
		block.Header.MerkleRoot,
		string(block.Header.Delegate),
		oids,
		time.Now(),
	}
	ldb.Insert("Blocks", newBlockInfo)
	return nil
}

func ApplyPersistanceRole(ldb *db.DBService, block *types.Block) error {
	fmt.Println("beging............applyPersistance")
	oids := make([]bson.ObjectId, len(block.Transactions))
	for i := range block.Transactions {
		oids[i] = bson.NewObjectId()
	}
	insertBlockInfoRole(ldb, block, oids)
	insertTxInfoRole(ldb, block, oids)
	return nil

}
