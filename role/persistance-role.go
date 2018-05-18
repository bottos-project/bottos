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
    "github.com/bottos-project/core/contract/msgpack"
)
type AccountInfo struct {
	ID               bson.ObjectId `bson:"_id"`
	AccountName      string        `bson:"account_name"`
	Balance          uint32        `bson:"bto_balance"`
	StakedBalance    uint64        `bson:"staked_balance"`
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
	BlockHash       string/*common.Hash*/     `bson:"block_hash"`
	PrevBlockHash   string/*[]byte*/   `bson:"prev_block_hash"`
	BlockNumber     uint32          `bson:"block_number"`
	Timestamp       uint64          `bson:"timestamp"`
	MerkleRoot      string/*[]byte*/   `bson:"merkle_root"`
	DelegateAccount string          `bson:"delegate"`
	Transactions    []bson.ObjectId `bson:"transactions"`
	CreateTime      time.Time       `bson:"create_time"`
}

type TxInfo struct {
	ID            bson.ObjectId `bson:"_id"`
	BlockNum      uint32        `bson:"block_number"`
	TransactionID string/*common.Hash*/   `bson:"transaction_id"`
	SequenceNum   uint32        `bson:"sequence_num"`
	BlockHash     string/*common.Hash*/   `bson:"block_hash"`
	CursorNum     uint32        `bson:"cursor_num"`
	CursorLabel   uint32        `bson:"cursor_label"`
	Lifetime      uint64        `bson:"lifetime"`
	Sender        string        `bson:"sender"`
	Contract      string        `bson:"contract"`
	Method        string        `bson:"method"`
	//Param         []byte        `bson:"param"`
	Param         TParam        `bson:"param"`
	SigAlg        uint32        `bson:"sig_alg"`
	Signature     string/*[]byte*/        `bson:"signature"`
	CreateTime    time.Time     `bson:"create_time"`
}

type transferparam struct {
    From        string  `json:"from"`
    To          string  `json:"to"`
    Value       uint64  `json: value`
}

type newaccountparam struct {
    Name        string  `json: name`
    Pubkey      string  `json: pubkey`
}

type reguser struct {
    Didid       string `json:"didid"`
    Didinfo     string `json:"didinfo"`
}

type TParam interface {
    //Accountparam
    //Transferparam transferpa
    //Reguser       reguser{}
    //DeployCodeParam
}

type DeployCodeParam struct {
    Name         string      `json:"name"`
    VMType       byte        `json:"vm_type"`
    VMVersion    byte        `json:"vm_version"`
    ContractCode []byte      `json:"contract_code"`
 }

type mgo_DeployCodeParam struct {
    Name         string      `json:"name"`
    VMType       byte        `json:"vm_type"`
    VMVersion    byte        `json:"vm_version"`
    ContractCode string      `json:"contract_code"`
 }

func findAcountInfo(ldb *db.DBService, accountName string) (*AccountInfo, error) {

	object, err := ldb.Find(config.DEFAULT_OPTIONDB_TABLE_TRX_NAME, "account_name", accountName)
	if err != nil {
		fmt.Println("find ", accountName, "failed ")
		return nil, errors.New("find " + accountName + "failed ")
	}
	return object.(*AccountInfo), nil
}

func ParseParam(Param []byte, Contract string, Method string) (interface{}, error) {
    var decodedParam interface{}
    if Contract == "bottos" {
        if Method == "newaccount" {
            decodedParam = &newaccountparam {}
        } else if Method == "transfer" {
            decodedParam = &transferparam {}
        } else if Method == "deploycode" {
            decodedParam = &DeployCodeParam {}
        } else {
            fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        } 
    } else if Contract == "usermsg" {
        if Method == "reguser" {
            decodedParam = &reguser{}
        }
    } else {
        fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
        return nil, errors.New("Not supported")
    }
    
    fmt.Println("insertTxInfoRole: done: Contract: ", Contract, ", Method: ", Method)
    err := msgpack.Unmarshal(Param, decodedParam)
    
    if Contract == "bottos" && Method == "deploycode" {
        p, ok := decodedParam.(DeployCodeParam)
        if ok {
            var mgo_param = mgo_DeployCodeParam {}
            mgo_param.Name      = p.Name
            mgo_param.VMType    = p.VMType
            mgo_param.VMVersion = p.VMVersion
            mgo_param.ContractCode = common.BytesToHex(p.ContractCode)
            fmt.Println("decodedParam OK!!!!")
            return mgo_param, nil
        } else {
            fmt.Println("decodedParam FAILED!!!: decodedParam: ", decodedParam)
            return nil, errors.New("Decode DeployCodeParam failed.")
        }
    }

    if err != nil {
        return nil, err
    }

    return decodedParam, nil
}

func insertTxInfoRole(r *Role, ldb *db.DBService, block *types.Block, oids []bson.ObjectId) error {
	
    fmt.Println("LYP: 11 START insertTxInfoRole!!!")
    
    if ldb == nil || block == nil {
		return errors.New("Error Invalid param")
	}
	if len(oids) != len(block.Transactions) {
		return errors.New("invalid param")
	}
    fmt.Println("LYP: 22 START  insertTxInfoRole!!! len(oids): ", len(oids), ", len(block.Transactions): ", len(block.Transactions))

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
			SigAlg:        trx.SigAlg,
			Signature:     common.BytesToHex(trx.Signature),
			CreateTime:    time.Now(),
		}
        
        fmt.Println("LYP: 33 START  insertTxInfoRole!!! Method: ", trx.Method, "|", newtrx.Method)
        decodedParam, err := ParseParam(trx.Param, newtrx.Contract, newtrx.Method)
        
        if err != nil {
            return err
        } else {
            newtrx.Param = decodedParam
        }

        fmt.Println("LYP: 33 START  insertTxInfoRole!!!")
		
        ldb.Insert(config.DEFAULT_OPTIONDB_TABLE_TRX_NAME, newtrx)
        fmt.Println("LYP: 44 START  insertTxInfoRole!!!")
		if trx.Contract == config.BOTTOS_CONTRACT_NAME {
            fmt.Println("LYP: 55 START  insertTxInfoRole!!!")
			insertAccountInfoRole(r, ldb, block, trx, oids[i])
		}
	}

	return nil
}

func insertBlockInfoRole(ldb *db.DBService, block *types.Block, oids []bson.ObjectId) error {
	if ldb == nil || block == nil {
		return errors.New("Error Invalid param")
	}
    fmt.Println("insertBlockInfoRole: len(oids):", len(oids), ", len(block.Transactions):", len(block.Transactions), ", block.Header.MerkleRoot: ", block.Header.MerkleRoot, " | ", common.BytesToHex(block.Header.MerkleRoot) )

	newBlockInfo := &BlockInfo{
		bson.NewObjectId(),
		block.Hash().ToHexString(),
		common.BytesToHex(block.Header.PrevBlockHash),
		block.Header.Number,
		block.Header.Timestamp,
		common.BytesToHex(block.Header.MerkleRoot),
		string(block.Header.Delegate),
		oids,
		time.Now(),
	}
	return ldb.Insert(config.DEFAULT_OPTIONDB_TABLE_BLOCK_NAME, newBlockInfo)
}

func insertAccountInfoRole(r *Role, ldb *db.DBService, block *types.Block, trx *types.Transaction, oid bson.ObjectId) error {
	if ldb == nil || block == nil {
		return errors.New("Error Invalid param")
	}
    
    if trx.Contract != config.BOTTOS_CONTRACT_NAME {
        return errors.New("Invalid contract param")
    }
    
    if trx.Method == "transfer" {
        data := &transferparam{}
        err :=  msgpack.Unmarshal(trx.Param, data)
        fmt.Printf("transfer struct: %v, msgpack: %x\n", trx.Param, data)
         
        FromAccountName := data.From
        ToAccountName   := data.To
        SrcBalanceInfo, err := r.GetBalance(FromAccountName)    //data.Value
         
        if(err != nil) {
            return err
        }

        DstBalanceInfo, err := r.GetBalance(ToAccountName)
         
        if(err != nil) {
            return err
        }

        if SrcBalanceInfo.Balance < data.Value {
            return err
        }

        SrcBalanceInfo.Balance -= data.Value
        DstBalanceInfo.Balance += data.Value
        
        err = r.SetBalance(FromAccountName, SrcBalanceInfo)
        if err != nil {
            return err
        }
        err = r.SetBalance(ToAccountName,   DstBalanceInfo)
        if err != nil {
            return err
        }
    } else if (trx.Method == "newaccount") {
        data := &newaccountparam{}
        err  :=  msgpack.Unmarshal(trx.Param, data)
        if err != nil {
            return err
        }
        fmt.Printf("transfer struct: %v, msgpack: %x\n", trx.Param, data)
            
        //accountInfos, err := GetAccount(data.Name)
            
        NewAccount := &AccountInfo {
            ID:               oid,
            AccountName:      data.Name,
            Balance:          0,//uint32        `bson:"bto_balance"`
            StakedBalance:    0,//uint64        `bson:"staked_balance"`
            UnstakingBalance: "",//             `bson:"unstaking_balance"`
            PublicKey:        []byte(data.Pubkey),
           // VMType:           0,// byte          `bson:"vm_type"`
           // VMVersion:        0, //byte          `bson:"vm_version"`
            //CodeVersion:      common.BytesToHash(123), //common.Hash   `bson:"code_version"`
            CreateTime:       time.Now(), //time.Time     `bson:"create_time"`
           // ContractCode:     "",  //[]byte   `bson:"contract_code"`
           // ContractAbi:      "",  //[]byte   `bson:"abi"`
            UpdatedTime:      time.Now(), //time.Time     `bson:"updated_time"`
       }
               
       return ldb.Insert(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, NewAccount)
    }

    return nil
}

func ApplyPersistanceRole(r *Role, ldb *db.DBService, block *types.Block) error {
	fmt.Println("LYP: len(block.Transactions): ", len(block.Transactions))
    oids := make([]bson.ObjectId, len(block.Transactions))
	for i := range block.Transactions {
		oids[i] = bson.NewObjectId()
	}
    
	insertBlockInfoRole(ldb, block, oids)
    fmt.Println("LYP: insertTxInfoRole!!!")
    insertTxInfoRole(r, ldb, block, oids)
    
    fmt.Printf("apply to mongodb block hash %x", block.Hash())
	return nil
}

//TODO start retro block when core start
func StartRetroBlock(ldb *db.DBService) {

}
