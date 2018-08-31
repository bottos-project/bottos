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
 * file description:  persistence role
 * @Author:
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package mongodb

import (
	"errors"
	log "github.com/cihub/seelog"
	"time"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"
	"gopkg.in/mgo.v2/bson"
	"math/big"
)

type MongoDBPlugin struct {
	Role role.RoleInterface
	Db   *db.OptionDBService
}

// AccountInfo is definition of account
type AccountInfo struct {
	ID               bson.ObjectId `bson:"_id"`
	AccountName      string        `bson:"account_name"`
	Balance          string        `bson:"bto_balance"`
	StakedBalance    string        `bson:"staked_balance"`
	UnstakingBalance string        `bson:"unstaking_balance"`
	PublicKey        string        `bson:"public_key"`
	CreateTime       time.Time     `bson:"create_time"`
	UpdatedTime      time.Time     `bson:"updated_time"`
}

// BlockInfo is definition of block
type BlockInfo struct {
	ID              bson.ObjectId `bson:"_id"`
	BlockHash       string/*common.Hash*/ `bson:"block_hash"`
	PrevBlockHash   string/*[]byte*/ `bson:"prev_block_hash"`
	BlockNumber     uint64 `bson:"block_number"`
	Timestamp       uint64 `bson:"timestamp"`
	MerkleRoot      string/*[]byte*/ `bson:"merkle_root"`
	DelegateAccount string          `bson:"delegate"`
	DelegateSign    string          `bson:"delegate_sign"`
	Transactions    []bson.ObjectId `bson:"transactions"`
	CreateTime      time.Time       `bson:"create_time"`
}

// TxInfo is definition of tx
type TxInfo struct {
	ID            bson.ObjectId `bson:"_id"`
	BlockNum      uint64        `bson:"block_number"`
	TransactionID string/*common.Hash*/ `bson:"transaction_id"`
	SequenceNum   uint32 `bson:"sequence_num"`
	BlockHash     string/*common.Hash*/ `bson:"block_hash"`
	CursorNum     uint64      `bson:"cursor_num"`
	CursorLabel   uint32      `bson:"cursor_label"`
	Lifetime      uint64      `bson:"lifetime"`
	Sender        string      `bson:"sender"`
	Contract      string      `bson:"contract"`
	Method        string      `bson:"method"`
	Param         interface{} `bson:"param"`
	SigAlg        uint32      `bson:"sig_alg"`
	Signature     string/*[]byte*/ `bson:"signature"`
	CreateTime    time.Time `bson:"create_time"`
}

func NewMongoDBPlugin(role role.RoleInterface, db *db.OptionDBService) *MongoDBPlugin {
	return &MongoDBPlugin{Role: role, Db: db}
}

func (mdb *MongoDBPlugin) ApplyBlock(block *types.Block) error {
	if !mdb.Db.IsOpDbConfigured() {
		return nil
	}
	oids := make([]bson.ObjectId, len(block.Transactions))
	for i := range block.Transactions {
		oids[i] = bson.NewObjectId()
	}
	mdb.insertBlockInfoRole(block, oids)
	mdb.insertTxInfoRole(block, oids)

	return nil
}

func (mdb *MongoDBPlugin) findAcountInfo(accountName string) (interface{}, error) {
	return mdb.Db.Find(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, "account_name", accountName)
}

// ParseParam is to parase param by method
func (mdb *MongoDBPlugin) ParseParam(Contract string, Method string, Param []byte) (map[string]interface{}, error) {
	var Abi *abi.ABI = nil
	if Contract != "bottos" {
		var err error
		Abi, err = mdb.getAbiForExternalContract(Contract)
		if err != nil {
			return nil, errors.New("External Abi is empty!")
		}
	} else {
		Abi = abi.GetAbi()
	}

	if Abi == nil {
		return nil, errors.New("Abi is empty!")
	}

	decodedParam := abi.UnmarshalAbiEx(Contract, Abi, Method, Param)
	if decodedParam == nil || len(decodedParam) <= 0 {
		log.Error("insertTxInfoRole: FAILED (decodedParam is nil!): Contract: ", Contract, ", Method: ", Method)
		return nil, errors.New("insertTxInfoRole: FAILED")
	}
	return decodedParam, nil
}

func (mdb *MongoDBPlugin) insertTxInfoRole(block *types.Block, oids []bson.ObjectId) error {
	if block == nil {
		return errors.New("Error Invalid param")
	}
	if len(oids) != len(block.Transactions) {
		return errors.New("invalid param")
	}
	for i, trx := range block.Transactions {
		newtrx := &TxInfo{
			ID:            oids[i],
			BlockNum:      block.Header.Number,
			TransactionID: trx.Hash().ToHexString(), //trx.Hash()
			SequenceNum:   uint32(i),
			BlockHash:     block.Hash().ToHexString(), //block.Hash(),
			CursorNum:     trx.CursorNum,
			CursorLabel:   trx.CursorLabel,
			Lifetime:      trx.Lifetime,
			Sender:        trx.Sender,
			Contract:      trx.Contract,
			Method:        trx.Method,
			//Param:         trx.Param,
			SigAlg:     trx.SigAlg,
			Signature:  common.BytesToHex(trx.Signature),
			CreateTime: time.Now(),
		}

		decodedParam, err := mdb.ParseParam(newtrx.Contract, newtrx.Method, trx.Param)
		if err != nil {
			continue
		}
		newtrx.Param = decodedParam

		mdb.Db.Insert(config.DEFAULT_OPTIONDB_TABLE_TRX_NAME, newtrx)
		if trx.Contract == config.BOTTOS_CONTRACT_NAME {
			mdb.insertAccountInfoRole(block, trx, oids[i])
		}
	}

	return nil
}

func (mdb *MongoDBPlugin) insertBlockInfoRole(block *types.Block, oids []bson.ObjectId) error {
	if block == nil {
		return errors.New("Error Invalid param")
	}

	newBlockInfo := &BlockInfo{
		bson.NewObjectId(),
		block.Hash().ToHexString(),
		common.BytesToHex(block.Header.PrevBlockHash),
		block.Header.Number,
		block.Header.Timestamp,
		common.BytesToHex(block.Header.MerkleRoot),
		string(block.Header.Delegate),
		common.BytesToHex(block.Header.DelegateSign),
		oids,
		time.Now(),
	}
	return mdb.Db.Insert(config.DEFAULT_OPTIONDB_TABLE_BLOCK_NAME, newBlockInfo)
}

// GetBalanceOp is to get account balance
func (mdb *MongoDBPlugin) GetBalanceOp(accountName string) (*role.Balance, error) {
	var value2 AccountInfo
	value, err := mdb.findAcountInfo(accountName)

	if value == nil || err != nil {
		return nil, err
	}

	// convert bson.M to struct
	bsonBytes, _ := bson.Marshal(value)
	bson.Unmarshal(bsonBytes, &value2)

	balance, errConvert := big.NewInt(0).SetString(value2.Balance, 10)

	if false == errConvert {
		return nil, errors.New("convert balance error")
	}

	res := &role.Balance{
		AccountName: accountName,
		Balance:     balance,
	}

	return res, nil
}

// SetBalanceOp is to save account balance
func (mdb *MongoDBPlugin) SetBalanceOp(accountName string, balance *big.Int) error {
	return mdb.Db.Update(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, "account_name", accountName, "bto_balance", balance)
}

func (mdb *MongoDBPlugin) insertAccountInfoRole(block *types.Block, trx *types.Transaction, oid bson.ObjectId) error {
	if block == nil {
		return errors.New("Error Invalid param")
	}

	if trx.Contract != config.BOTTOS_CONTRACT_NAME {
		return errors.New("Invalid contract param")
	}

	var initSupply *big.Int = big.NewInt(0)
	var err error
	initSupply, err = safemath.U256Mul(initSupply, new(big.Int).SetUint64(config.BOTTOS_INIT_SUPPLY), new(big.Int).SetUint64(config.BOTTOS_SUPPLY_MUL))

	if err != nil {
		return err
	}

	Abi := abi.GetAbi()

	_, err = mdb.findAcountInfo(config.BOTTOS_CONTRACT_NAME)
	if err != nil {

		bto := &AccountInfo{
			ID:               bson.NewObjectId(),
			AccountName:      config.BOTTOS_CONTRACT_NAME,
			Balance:          initSupply.String(),
			StakedBalance:    "0",
			UnstakingBalance: "0",
			PublicKey:        config.Param.KeyPairs[0].PublicKey,
			CreateTime:       time.Unix(int64(config.Genesis.GenesisTime), 0), //time.Time     `bson:"create_time"`
			UpdatedTime:      time.Now(),                                      //time.Time     `bson:"updated_time"`
		}
		mdb.Db.Insert(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, bto)
	}

	if trx.Method == "transfer" {
		data := abi.UnmarshalAbiEx(trx.Contract, Abi, trx.Method, trx.Param)
		if data == nil || len(data) <= 0 {
			log.Error("UnmarshalAbi for contract: ", trx.Contract, ", Method: ", trx.Method, " failed!")
		}

		FromAccountName := data["from"].(string)
		ToAccountName := data["to"].(string)
		DataVal := data["value"].(*big.Int)

		SrcBalanceInfo, err := mdb.GetBalanceOp(FromAccountName) //data.Value
		if err != nil {
			return err
		}

		DstBalanceInfo, err := mdb.GetBalanceOp(ToAccountName)
		if err != nil {
			return err
		}

		if -1 == SrcBalanceInfo.Balance.Cmp(DataVal) {
			return err
		}

		SrcBalanceInfo.SafeSub(DataVal)
		DstBalanceInfo.SafeAdd(DataVal)

		err = mdb.SetBalanceOp(FromAccountName, SrcBalanceInfo.Balance)
		if err != nil {
			return err
		}
		err = mdb.SetBalanceOp(ToAccountName, DstBalanceInfo.Balance)
		if err != nil {
			return err
		}
	} else if trx.Method == "newaccount" {
		data := abi.UnmarshalAbiEx(trx.Contract, Abi, trx.Method, trx.Param)
		if data == nil || len(data) <= 0 {
			log.Error("UnmarshalAbi for contract: ", trx.Contract, ", Method: ", trx.Method, " failed!")
			return err
		}
		DataName := data["name"].(string)
		DataPubKey := data["pubkey"].(string)

		mesgs, err := mdb.findAcountInfo(DataName)
		if mesgs != nil {
			return nil /* Do not allow insert same account */
		}
		log.Info(err)

		NewAccount := &AccountInfo{
			ID:               oid,
			AccountName:      DataName,
			Balance:          "0", //uint32        `bson:"bto_balance"`
			StakedBalance:    "0", //uint64        `bson:"staked_balance"`
			UnstakingBalance: "0", //             `bson:"unstaking_balance"`
			PublicKey:        DataPubKey,
			CreateTime:       time.Now(), //time.Time     `bson:"create_time"`
			UpdatedTime:      time.Now(), //time.Time     `bson:"updated_time"`
		}

		return mdb.Db.Insert(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, NewAccount)
	}

	return nil
}

//GetAbi function
var ExternalAbiMap map[string]interface{}

func (mdb *MongoDBPlugin) getAbiForExternalContract(contract string) (*abi.ABI, error) {

	if len(ExternalAbiMap) <= 0 {
		ExternalAbiMap = make(map[string]interface{})
	}

	if _, ok := ExternalAbiMap[contract]; ok {
		return ExternalAbiMap[contract].(*abi.ABI), nil
	}

	account, err := mdb.Role.GetAccount(contract)
	if err != nil {
		return nil, errors.New("Get account fail")
	}

	if len(account.ContractAbi) > 0 {

		Abi, err := abi.ParseAbi(account.ContractAbi)
		if err != nil {
			return nil, err
		}

		ExternalAbiMap[contract] = Abi

		return Abi, nil
	}

	// TODO
	return nil, errors.New("Get Abi failed!")
}
