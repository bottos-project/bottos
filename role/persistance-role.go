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

package role

import (
	"errors"
	log "github.com/cihub/seelog"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/safemath"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	abi "github.com/bottos-project/bottos/contract/abi"
	"github.com/bottos-project/bottos/db"
	"gopkg.in/mgo.v2/bson"
)

// AccountInfo is definition of account
type AccountInfo struct {
	ID               bson.ObjectId `bson:"_id"`
	AccountName      string        `bson:"account_name"`
	Balance          uint64        `bson:"bto_balance"`
	StakedBalance    uint64        `bson:"staked_balance"`
	UnstakingBalance string        `bson:"unstaking_balance"`
	PublicKey        string        `bson:"public_key"`
	CreateTime       time.Time     `bson:"create_time"`
	UpdatedTime      time.Time     `bson:"updated_time"`
}

// BlockInfo is definition of block
type BlockInfo struct {
	ID              bson.ObjectId          `bson:"_id"`
	BlockHash       string /*common.Hash*/ `bson:"block_hash"`
	PrevBlockHash   string /*[]byte*/      `bson:"prev_block_hash"`
	BlockNumber     uint64                 `bson:"block_number"`
	Timestamp       uint64                 `bson:"timestamp"`
	MerkleRoot      string /*[]byte*/      `bson:"merkle_root"`
	DelegateAccount string                 `bson:"delegate"`
	DelegateSign    string                 `bson:"delegate_sign"`
	Transactions    []bson.ObjectId        `bson:"transactions"`
	CreateTime      time.Time              `bson:"create_time"`
}

// TxInfo is definition of tx
type TxInfo struct {
	ID            bson.ObjectId          `bson:"_id"`
	BlockNum      uint64                 `bson:"block_number"`
	TransactionID string /*common.Hash*/ `bson:"transaction_id"`
	SequenceNum   uint32                 `bson:"sequence_num"`
	BlockHash     string /*common.Hash*/ `bson:"block_hash"`
	CursorNum     uint64                 `bson:"cursor_num"`
	CursorLabel   uint32                 `bson:"cursor_label"`
	Lifetime      uint64                 `bson:"lifetime"`
	Sender        string                 `bson:"sender"`
	Contract      string                 `bson:"contract"`
	Method        string                 `bson:"method"`
	Param         TParam                 `bson:"param"`
	SigAlg        uint32                 `bson:"sig_alg"`
	Signature     string /*[]byte*/      `bson:"signature"`
	CreateTime    time.Time              `bson:"create_time"`
}

// TParam is interface definition
type TParam interface {
	//Accountparam
	//Transferparam transferpa
	//Reguser       reguser{}
	//DeployCodeParam
}

func findAcountInfo(ldb *db.DBService, accountName string) (interface{}, error) {
	return ldb.Find(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, "account_name", accountName)
}

//getMyPublicIPaddr function
func getMyPublicIPaddr() (string, error) {
	resp, err1 := http.Get("http://members.3322.org/dyndns/getip")
	if err1 != nil {
		return "", err1
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return "", err2
	}
	//log.Info("getMyPublicIPaddr"+ string(body))
	return string(body), nil
}

// ParseParam is to parase param by method
func ParseParam(r *Role, Param []byte, Contract string, Method string) (interface{}, error) {
	var Abi *abi.ABI = nil
	if Contract != "bottos" {
		var err error
		Abi, err = GetAbiForExternalContract(r, Contract)
		if  err != nil {
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

func insertTxInfoRole(r *Role, ldb *db.DBService, block *types.Block, oids []bson.ObjectId) error {

	if ldb == nil || block == nil {
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

		decodedParam, err := ParseParam(r, trx.Param, newtrx.Contract, newtrx.Method)

		if err != nil {
			continue
		}

		newtrx.Param = decodedParam

		ldb.Insert(config.DEFAULT_OPTIONDB_TABLE_TRX_NAME, newtrx)
		if trx.Contract == config.BOTTOS_CONTRACT_NAME {
			insertAccountInfoRole(r, ldb, block, trx, oids[i])
		}
	}

	return nil
}

func insertBlockInfoRole(ldb *db.DBService, block *types.Block, oids []bson.ObjectId) error {
	if ldb == nil || block == nil {
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
	return ldb.Insert(config.DEFAULT_OPTIONDB_TABLE_BLOCK_NAME, newBlockInfo)
}

// GetBalanceOp is to get account balance
func GetBalanceOp(ldb *db.DBService, accountName string) (*Balance, error) {
	var value2 AccountInfo
	value, err := findAcountInfo(ldb, accountName)

	if value == nil || err != nil {
		return nil, err
	}

	// convert bson.M to struct
	bsonBytes, _ := bson.Marshal(value)
	bson.Unmarshal(bsonBytes, &value2)

	res := &Balance{
		AccountName: accountName,
		Balance:     value2.Balance,
	}

	return res, nil
}

// SetBalanceOp is to save account balance
func SetBalanceOp(ldb *db.DBService, accountName string, balance uint64) error {
	return ldb.Update(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, "account_name", accountName, "bto_balance", balance)
}

func insertAccountInfoRole(r *Role, ldb *db.DBService, block *types.Block, trx *types.Transaction, oid bson.ObjectId) error {
	if ldb == nil || block == nil {
		return errors.New("Error Invalid param")
	}

	if trx.Contract != config.BOTTOS_CONTRACT_NAME {
		return errors.New("Invalid contract param")
	}

	var initSupply uint64
	var err error
	initSupply, err = safemath.Uint64Mul(config.BOTTOS_INIT_SUPPLY, config.BOTTOS_SUPPLY_MUL)
	if err != nil {
		return err
	}
	
	Abi := abi.GetAbi()

	_, err = findAcountInfo(ldb, config.BOTTOS_CONTRACT_NAME)
	if err != nil {

		bto := &AccountInfo{
			ID:               bson.NewObjectId(),
			AccountName:      config.BOTTOS_CONTRACT_NAME,
			Balance:          initSupply, //uint32        `bson:"bto_balance"`
			StakedBalance:    0,          //uint64        `bson:"staked_balance"`
			UnstakingBalance: "",         //             `bson:"unstaking_balance"`
			PublicKey:        config.Param.KeyPairs[0].PublicKey,
			CreateTime:       time.Unix(int64(config.Genesis.GenesisTime), 0), //time.Time     `bson:"create_time"`
			UpdatedTime:      time.Now(),                                      //time.Time     `bson:"updated_time"`
		}
		ldb.Insert(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, bto)
	}

	if trx.Method == "transfer" {

		data := abi.UnmarshalAbiEx(trx.Contract, Abi, trx.Method, trx.Param)
		if data == nil || len(data) <= 0 {
			log.Error("UnmarshalAbi for contract: ", trx.Contract, ", Method: ", trx.Method, " failed!")
		}

		FromAccountName := data["from"].(string)
		ToAccountName := data["to"].(string)
		DataVal := data["value"].(uint64)
		
		SrcBalanceInfo, err := GetBalanceOp(ldb, FromAccountName) //data.Value

		if err != nil {
			return err
		}

		DstBalanceInfo, err := GetBalanceOp(ldb, ToAccountName)

		if err != nil {
			return err
		}
		
		if SrcBalanceInfo.Balance < DataVal {
			return err
		}

		SrcBalanceInfo.Balance -= DataVal
		DstBalanceInfo.Balance += DataVal

		err = SetBalanceOp(ldb, FromAccountName, SrcBalanceInfo.Balance)
		if err != nil {
			return err
		}
		err = SetBalanceOp(ldb, ToAccountName, DstBalanceInfo.Balance)
		if err != nil {
			return err
		}
	} else if trx.Method == "newaccount" {

		data := abi.UnmarshalAbiEx(trx.Contract, Abi, trx.Method, trx.Param)
		if data == nil || len(data) <= 0 {
			log.Error("UnmarshalAbi for contract: ", trx.Contract, ", Method: ", trx.Method, " failed!")
			return err
		}
		DataName   := data["name"].(string)
		DataPubKey := data["pubkey"].(string)

		mesgs, err := findAcountInfo(ldb, DataName)
		if mesgs != nil {
			return nil /* Do not allow insert same account */
		}
		log.Info(err)

		NewAccount := &AccountInfo{
			ID:               oid,
			AccountName:      DataName,
			Balance:          0,  //uint32        `bson:"bto_balance"`
			StakedBalance:    0,  //uint64        `bson:"staked_balance"`
			UnstakingBalance: "", //             `bson:"unstaking_balance"`
			PublicKey:        DataPubKey,
			CreateTime:       time.Now(), //time.Time     `bson:"create_time"`
			UpdatedTime:      time.Now(), //time.Time     `bson:"updated_time"`
		}

		return ldb.Insert(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, NewAccount)
	}

	return nil
}

// ApplyPersistanceRole is to apply persistence
func ApplyPersistanceRole(r *Role, ldb *db.DBService, block *types.Block) error {
	if !ldb.IsOpDbConfigured() {
		return nil
	}
	oids := make([]bson.ObjectId, len(block.Transactions))
	for i := range block.Transactions {
		oids[i] = bson.NewObjectId()
	}
	insertBlockInfoRole(ldb, block, oids)
	insertTxInfoRole(r, ldb, block, oids)

	//fmt.Printf("apply to mongodb block hash %x, block number %d", block.Hash(), block.Header.Number)
	return nil
}

// StartRetroBlock is to do: start retro block when core start
func StartRetroBlock(ldb *db.DBService) {

}

//GetAbi function
var ExternalAbiMap map[string]interface{}

func GetAbiForExternalContract(r *Role, contract string) (*abi.ABI, error) {
	
	if len(ExternalAbiMap) <= 0 {
		ExternalAbiMap = make(map[string]interface{})
	}
	
	if _, ok := ExternalAbiMap[contract]; ok {
		return ExternalAbiMap[contract].(*abi.ABI), nil	
	}
	
	account, err := r.GetAccount(contract)
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

