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
	//"fmt"
	"time"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/db"
	"gopkg.in/mgo.v2/bson"
    "github.com/bottos-project/bottos/contract/msgpack"
	"github.com/bottos-project/bottos/common/safemath"
)
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

type SetDelegateParam struct {
    Name        string  `json: name`
    Pubkey      string  `json: pubkey`
}

type GrantCreditParam struct {
	Name		string		`json:"name"`
	Spender		string 		`json:"spender"`
	Limit		uint64		`json:"limit"`
}

type CancelCreditParam struct {
	Name		string		`json:"name"`
	Spender		string 		`json:"spender"`
}

type TransferFromParam struct {
	From		string		`json:"from"`
	To			string		`json:"to"`
	Value		uint64		`json:"value"`
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
    AssetType           uint64 `json:"assettype"`
    FeatureTag     string `json:"featuretag"`
    SampleHash     string `json:"samplehash"`
    StorageHash    string `json:"storagehash"`
    ExpireTime          uint32 `json:"expiretime"`
    OpType              uint32 `json:"optype"`
    Price               uint64 `json:"price"`
    Description         string `json:"description"`
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
    AssetId      string `json:"assetid"`
}

type DataDealReq struct {
    DataExchangeId string  `json:"dataexchangeid"`
    Info DataDealnfo
}

type PresaleInfo struct {
    UserName string    `json:"username"`
    AssetId string     `json:"assetid"`
    DataReqId string   `json:"datareqid"`
    Consumer string    `json:"consumer"`
    OpType   uint32    `json:"optype"` 
}


type PresaleReq struct {
    DataPresaleId string   `json:"datapresaleid"`
    Info PresaleInfo
}

type DataFileInfo struct {
    UserName    string  `json:"username"` 
    FileSize    uint64 `json:"filesize"`
    FileName    string `json:"filename"`
    FilePolicy  string `json:"filepolicy"`
    FileNumber  uint64 `json:"filenumber"`
    Simorass    uint32 `json:"simorass"`
    OpType      uint32 `json:"optype"`
    StoreAddr   string `json:"storeaddr"`
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

func findAcountInfo(ldb *db.DBService, accountName string) (interface{}, error) {
    return ldb.Find(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, "account_name", accountName)
}

func ParseParam(Param []byte, Contract string, Method string) (interface{}, error) {
    var decodedParam interface{}
    if Contract == "bottos" {
        if Method == "newaccount" {
            decodedParam = &newaccountparam {}
        } else if Method == "setdelegate" {
            decodedParam = &SetDelegateParam {}
        } else if Method == "transfer" {
            decodedParam = &transferparam {}
        } else if Method == "deploycode" {
            decodedParam = &DeployCodeParam {}
        } else if Method == "grantcredit" {
            decodedParam = &GrantCreditParam {}
        } else if Method == "cancelcredit" {
            decodedParam = &CancelCreditParam {}
        } else if Method == "transferfrom" {
            decodedParam = &TransferFromParam {}
        } else {
            //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract, ", Method: ", Method)
            return nil, errors.New("Not supported")
        } 
    } else if Contract == "usermng" {
        if Method == "reguser" {
            decodedParam = &reguser{}
        } else if Method == "userlogin" {
            decodedParam = &UserLogin{}
        } else {
            //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        } 
    } else if Contract == "assetmng" {
         if Method == "assetreg" {
            decodedParam = &RegAssetReq {}
        } else {
            //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        } 
    } else if Contract == "datadealmng" {
        if Method == "buydata" {
            decodedParam = &DataDealReq {}
        } else if Method == "presale" {
            decodedParam = &PresaleReq {}
        } else {
            //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        }
    } else if Contract == "datafilemng" {
        if Method == "datafilereg" {
           decodedParam = &DataFileRegReq {}
        } else if Method == "fileauthreg" {
            decodedParam = &DataFileAuthReq {}
        } else {
            //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        }
    } else if Contract == "datareqmng" {
       if Method == "datareqreg" {
           decodedParam = &RegDataReqReq {}
       } else {
           //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
           return nil, errors.New("Not supported")
       }
    } else if Contract == "favoritemng" {
        if Method == "favoritepro" {
            decodedParam = &GoodsProReq {}
        } else {
           //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
           return nil, errors.New("Not supported")
       }
    } else if Contract == "nodeclustermng" {
        if Method == "reg" {
            decodedParam = &NodeClusterReg {}
        } else {
            //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        }
    } else if Contract == "nodemng" {
        if Method == "nodeinforeg" {
            decodedParam = &NodeInfoReq {}       
        } else {
            //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
            return nil, errors.New("Not supported")
        }
    } else {
        //fmt.Println("insertTxInfoRole:Not supported: Contract: ", Contract)
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
        //fmt.Println("insertTxInfoRole: FAILED: Contract: ", Contract, ", Method: ", Method)
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
    //fmt.Println("insertBlockInfoRole: len(oids):", len(oids), ", len(block.Transactions):", len(block.Transactions), ", block.Header.MerkleRoot: ", block.Header.MerkleRoot, " | ", common.BytesToHex(block.Header.MerkleRoot) )

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

func GetBalanceOp(ldb *db.DBService, accountName string) (*Balance, error) {
    var value2 AccountInfo
    value, err := findAcountInfo(ldb, accountName)

    if value == nil || err != nil {
        return nil, err
    }
    
    // convert bson.M to struct
    bsonBytes, _ := bson.Marshal(value)
    bson.Unmarshal(bsonBytes, &value2)

    //fmt.Println("GetBalanceOp: value2 is: ", value2, ", value2.Balance is: ", value2.Balance)
    
    res := &Balance{
                       AccountName: accountName,
                       Balance:     value2.Balance,
                   }
    
    return res, nil
}

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

    _, err = findAcountInfo(ldb, config.BOTTOS_CONTRACT_NAME)
    if err != nil {

        bto := &AccountInfo{
            ID:               bson.NewObjectId(),
            AccountName:      config.BOTTOS_CONTRACT_NAME,
            Balance:          initSupply,//uint32        `bson:"bto_balance"`
            StakedBalance:    0,//uint64        `bson:"staked_balance"`
            UnstakingBalance: "",//             `bson:"unstaking_balance"`
            PublicKey:        config.Param.KeyPairs[0].PublicKey,
            CreateTime:       time.Unix(int64(config.Genesis.GenesisTime), 0), //time.Time     `bson:"create_time"`
            UpdatedTime:      time.Now(), //time.Time     `bson:"updated_time"`
        }
        ldb.Insert(config.DEFAULT_OPTIONDB_TABLE_ACCOUNT_NAME, bto)
    }

    if trx.Method == "transfer" {
        
        data := &transferparam{}
        err :=  msgpack.Unmarshal(trx.Param, data)
        //fmt.Printf("transfer struct: %v, msgpack: %x\n", trx.Param, data)
         
        FromAccountName := data.From
        ToAccountName   := data.To
        SrcBalanceInfo, err := GetBalanceOp(ldb, FromAccountName)    //data.Value
         
        if(err != nil) {
            return err
        }

        DstBalanceInfo, err := GetBalanceOp(ldb, ToAccountName)
         
        if(err != nil) {
            return err
        }

        if SrcBalanceInfo.Balance < data.Value {
            return err
        }

        SrcBalanceInfo.Balance -= data.Value
        DstBalanceInfo.Balance += data.Value
        
        err = SetBalanceOp(ldb, FromAccountName, SrcBalanceInfo.Balance)
        if err != nil {
            return err
        }
        err = SetBalanceOp(ldb, ToAccountName, DstBalanceInfo.Balance)
        if err != nil {
            return err
        }
    } else if (trx.Method == "newaccount") {
        
        data := &newaccountparam{}
        err  :=  msgpack.Unmarshal(trx.Param, data)
        if err != nil {
            return err
        }
        
        mesgs, err := findAcountInfo(ldb, data.Name)
        if mesgs != nil {
           return nil /* Do not allow insert same account */ 
        }

        NewAccount := &AccountInfo {
            ID:               oid,
            AccountName:      data.Name,
            Balance:          0,//uint32        `bson:"bto_balance"`
            StakedBalance:    0,//uint64        `bson:"staked_balance"`
            UnstakingBalance: "",//             `bson:"unstaking_balance"`
            PublicKey:        data.Pubkey,
            CreateTime:       time.Now(), //time.Time     `bson:"create_time"`
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
    
    //fmt.Printf("apply to mongodb block hash %x, block number %d", block.Hash(), block.Header.Number)
	return nil
}

//TODO start retro block when core start
func StartRetroBlock(ldb *db.DBService) {

}
