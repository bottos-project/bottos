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
    DelegateSign    string          `bson:"delegate_sign"`
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

/**======Internal Contract struct definition====*/

type transferparam struct {
    From        string  `json:"from"`
    To          string  `json:"to"`
    Value       uint64  `json: value`
}

type newaccountparam struct {
    Name        string  `json: name`
    Pubkey      string  `json: pubkey`
}

type TParam interface {
    //Accountparam
    //Transferparam transferpa
    //Reguser       reguser{}
    //DeployCodeParam
}

type DeployCodeParam struct {
    Name         string     `json:"name"`
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


/**======External Contract struct definition====*/

type AssetInfo struct {
    UserName            string `json:“username”`
    AssetName           string `json:"assetname"`
    AssetType           string `json:"assettype"`
    FeatureTag     string `json:"featuretag"`
    SamplePath     string `json:"samplepath"`
    SampleHash     string `json:"samplehash"`
    StoragePath    string `json:"storagepath"`      // not enough for multislices of big files
    StorageHash    string `json:"storagehash"`
    ExpireTime          uint32 `json:"expiretime"`
    Price               uint64 `json:"price"`
    Description         string `json:"description"`
    UploadDate          uint32 `json:"uploaddate"`
    Signature           string `json:"signature"`
}

type RegAssetReq struct {
    AssetId string `json:"assedid"`
    Info AssetInfo 
}

type reguser struct {
    Didid        string `json:"didid"`
    Didinfo      string `json:"didinfo"`
}

type UserLogin struct {
    UserName string    `json:"username"`
    RandomNum uint32 `json:"randomnum"`
}

type DataDealnfo struct {
    UserName     string `json:"username"`
    SessionId    string `json:”sessionid“`
    AssetId      string `json:"assetid"`
    Random_num   uint64 `json:"random_num"`
    Signature    uint64 `json:"signature"`
}

type DataDealReq struct {
    DataExchangeId string  `json:"dataexchangeid"`
    Info DataDealnfo
}

type PresaleInfo struct {
    UserName string    `json:"username"`
    SessionId string   `json:"sessionid"`
    AssetId string     `json:"assetid"`
    AssetName string   `json:"assetname"`
    DataReqId string   `json:"datareqid"`
    DataReqName string `json:"datareqname"`
    Consumer string    `json:"consumer"`
    Random_num uint64  `json:"randomnum"`
    Signature string   `json:"signature"`
}

type PresaleReq struct {
    DataPresaleId string   `json:"datapresaleid"`
    Info PresaleInfo
}

type DataFileInfo struct {
    UserName    string  `json:"username"` 
    SessionId   string `json:"sessonid"`
    FileSize    uint64 `json:"filesize"`
    FileName    string `json:"filename"`
    FilePolicy  string `json:"filepolicy"`
    AuthPath    string `json:"authpath"`
    FileNumber  uint64 `json:"filenumber"`
    Signature   string `json:"signature"`
}

type DataFileRegReq struct {
    FileHash   string  `json:"filehash"`
    Info    DataFileInfo
}

type AuthBasicInfo struct {
    AuthType   string  `json:"authType"`
    AuthPath   string  `json:"authpath"`
}

type DataFileAuthInfo struct {
    HashUserName    string `json:"hashusername"`
    Info            AuthBasicInfo
}

type DataFileAuthReq struct {
    StorgeHash  string `json:"storagehash"`
    UserName    string `json:"username"`
}

type DataReqInfo struct {

    UserName    string `json:"username"`
    ReqName     string `json:"reqname"`
    ReqType     uint64 `json:"reqtype"`
    FeatureTag  uint64 `json:"featuretag"`
    SampleHash  string `json:"samplehash"`
    ExpireTime  uint64 `json:"expiretime"`
    OpType      uint32 `json:"optype"`
    Price       uint64 `json:"price"`
    FavoriFlag  uint32 `json:"favoriflag"`
    Description string `description`
}

type RegDataReqReq struct{
    DataReqId   string `json:"datareqid"`
    Info        DataReqInfo
}

type GoodsProReq struct {
    UserName    string `json:"username"`
    OpType      uint32 `json:"optype"`
    GoodsType   string `json:"goodstype"`
    GoodsId     string `json:"goodsid"`
}

type NodeClusterReg struct {
    NodeIP    string   `json:"nodeip"`
    ClusterIP string   `json:"clusterip"`
}

type NodeBaseInfo struct {
    NodeIp     string  `json:"nodeip"`
    NodePort   string  `json:"nodeport"`
    NodeAddress string `json:"nodeaddress"`
}

type NodeInfoReq struct {
    NodeId  string `json:"nodeid"`
    Info    NodeBaseInfo
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
    } else if Contract == "usermng" {
        if Method == "reguser" {
            decodedParam = &reguser{}
        } else if Method == "userlogin" {
            decodedParam = &UserLogin{}
        } else {
            fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        } 
    } else if Contract == "assetmng" {
         if Method == "assetreg" {
            decodedParam = &RegAssetReq {}
        } else {
            fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        } 
    } else if Contract == "datadealmng" {
        if Method == "buydata" {
            decodedParam = &DataDealReq {}
        } else if Method == "presale" {
            decodedParam = &PresaleReq {}
        } else {
            fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        }
    } else if Contract == "datafilemng" {
        if Method == "datafilereg" {
           decodedParam = &DataFileRegReq {}
        } else if Method == "fileauthreg" {
            decodedParam = &DataFileAuthReq {}
        } else {
            fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        }
    } else if Contract == "datareqmng" {
       if Method == "datareqreg" {
           decodedParam = &RegDataReqReq {}
       } else {
           fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
           return nil, errors.New("Not supported")
       }
    } else if Contract == "favoritemng" {
        if Method == "favoritepro" {
            decodedParam = &GoodsProReq {}
        } else {
           fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
           return nil, errors.New("Not supported")
       }
    } else if Contract == "nodeclustermng" {
        if Method == "reg" {
            decodedParam = &NodeClusterReg {}
        } else {
            fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        }
    } else if Contract == "nodemng" {
        if Method == "nodeinforeg" {
            decodedParam = &NodeInfoReq {}       
        } else {
            fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        }
    } else {
        fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
        return nil, errors.New("Not supported")
    }
    
    err := msgpack.Unmarshal(Param, decodedParam)
    
    if Contract == "bottos" && Method == "deploycode" {
        //p, ok := decodedParam.(DeployCodeParam)
        
        var tmpval = &DeployCodeParam {}
        err = msgpack.Unmarshal(Param, tmpval)
        //if ok {
        if err == nil {
            var mgo_param = mgo_DeployCodeParam {}
            mgo_param.Name      = tmpval.Name
            mgo_param.VMType    = tmpval.VMType
            mgo_param.VMVersion = tmpval.VMVersion
            mgo_param.ContractCode = common.BytesToHex(tmpval.ContractCode)
            return mgo_param, nil
        } else {
            return nil, errors.New("Decode DeployCodeParam failed.")
        }
    }
    
    if err != nil {
        fmt.Println("insertTxInfoRole: FAILED: Contract: ", Contract, ", Method: ", Method)
        return nil, err
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
			SigAlg:        trx.SigAlg,
			Signature:     common.BytesToHex(trx.Signature),
			CreateTime:    time.Now(),
		}
        
        decodedParam, err := ParseParam(trx.Param, newtrx.Contract, newtrx.Method)
        
        if err != nil {
            return err
        } else {
            newtrx.Param = decodedParam
        }
		
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
    fmt.Println("insertBlockInfoRole: len(oids):", len(oids), ", len(block.Transactions):", len(block.Transactions), ", block.Header.MerkleRoot: ", block.Header.MerkleRoot, " | ", common.BytesToHex(block.Header.MerkleRoot) )

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
    oids := make([]bson.ObjectId, len(block.Transactions))
	for i := range block.Transactions {
		oids[i] = bson.NewObjectId()
	}
    
	insertBlockInfoRole(ldb, block, oids)
    insertTxInfoRole(r, ldb, block, oids)
    
    fmt.Printf("apply to mongodb block hash %x", block.Hash())
	return nil
}

//TODO start retro block when core start
func StartRetroBlock(ldb *db.DBService) {

}
