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
	nodeApi "github.com/bottos-project/magiccube/service/node/api"
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
	BlockNumber     uint32                 `bson:"block_number"`
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
	BlockNum      uint32                 `bson:"block_number"`
	TransactionID string /*common.Hash*/ `bson:"transaction_id"`
	SequenceNum   uint32                 `bson:"sequence_num"`
	BlockHash     string /*common.Hash*/ `bson:"block_hash"`
	CursorNum     uint32                 `bson:"cursor_num"`
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

/**======Internal Contract struct definition====*/

// transferparam is interface definition of transfer method
type transferparam struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint64 `json:"value"`
}

// newaccountparam is interface definition of new account method
type newaccountparam struct {
	Name   string `json:"name"`
	Pubkey string `json:"pubkey"`
}

// SetDelegateParam is interface definition of new set delegate method
type SetDelegateParam struct {
	Name   string `json:"name"`
	Pubkey string `json:"pubkey"`
}

// GrantCreditParam is interface definition of grant credit method
type GrantCreditParam struct {
	Name    string `json:"name"`
	Spender string `json:"spender"`
	Limit   uint64 `json:"limit"`
}

// CancelCreditParam is interface definition of cancel credit method
type CancelCreditParam struct {
	Name    string `json:"name"`
	Spender string `json:"spender"`
}

// TransferFromParam is interface definition of transfer from method
type TransferFromParam struct {
	From      string `json:"from"`
	To        string `json:"to"`
	TokenType string `json:"tokenType"`
	Value     uint64 `json:"value"`
}

// TParam is interface definition
type TParam interface {
	//Accountparam
	//Transferparam transferpa
	//Reguser       reguser{}
	//DeployCodeParam
}

// DeployCodeParam is interface definition of deploy code method
type DeployCodeParam struct {
	Name         string `json:"name"`
	VMType       byte   `json:"vm_type"`
	VMVersion    byte   `json:"vm_version"`
	ContractCode []byte `json:"contract_code"`
}

// DeployAbiParam is interface definition of deploy code method
type DeployAbiParam struct {
	Name        string `json:"contract"`
	ContractAbi []byte `json:"contract_abi"`
}

// mgo_DeployCodeParam is interface definition of mgo deploy code
type mgoDeployCodeParam struct {
	Name         string `json:"name"`
	VMType       byte   `json:"vm_type"`
	VMVersion    byte   `json:"vm_version"`
	ContractCode string `json:"contract_code"`
}

type mgoDeployAbiParam struct {
	Name        string `json:"contract"`
	ContractAbi string `json:"contract_abi"`
}

/**======External Contract struct definition====*/

// AssetInfo is definition of asset info
type AssetInfo struct {
	UserName    string `json:"username"`
	AssetName   string `json:"assetname"`
	AssetType   uint64 `json:"assettype"`
	FeatureTag  string `json:"featuretag"`
	SampleHash  string `json:"samplehash"`
	StorageHash string `json:"storagehash"`
	ExpireTime  uint32 `json:"expiretime"`
	OpType      uint32 `json:"optype"`
	TokenType   string `json:"tokenType"`
	Price       uint64 `json:"price"`
	Description string `json:"description"`
}

// RegAssetReq is definition of
type RegAssetReq struct {
	AssetId string `json:"assedid"`
	Info    AssetInfo
}

// reguser is definition of reg user
type reguser struct {
	Didid   string `json:"didid"`
	Didinfo string `json:"didinfo"`
}

// UserLogin is definition of user login
type UserLogin struct {
	UserName  string `json:"username"`
	RandomNum uint32 `json:"randomnum"`
}

// DataDealnfo is definition of data deal info
type DataDealnfo struct {
	UserName string `json:"username"`
	AssetId  string `json:"assetid"`
}

// DataDealReq is definition of
type DataDealReq struct {
	DataExchangeId string `json:"dataexchangeid"`
	Info           DataDealnfo
}

// PresaleInfo is definition of pre sale info
type PresaleInfo struct {
	UserName  string `json:"username"`
	AssetId   string `json:"assetid"`
	DataReqId string `json:"datareqid"`
	Consumer  string `json:"consumer"`
	OpType    uint32 `json:"optype"`
}

// PresaleReq is definition of pre sale req
type PresaleReq struct {
	DataPresaleId string `json:"datapresaleid"`
	Info          PresaleInfo
}

// DataFileInfo is definition of data file info
type DataFileInfo struct {
	UserName   string `json:"username"`
	FileSize   uint64 `json:"filesize"`
	FileName   string `json:"filename"`
	FilePolicy string `json:"filepolicy"`
	FileNumber uint64 `json:"filenumber"`
	Simorass   uint32 `json:"simorass"`
	OpType     uint32 `json:"optype"`
	StoreAddr  string `json:"storeaddr"`
}

// DataFileRegReq is definition of data file req reg
type DataFileRegReq struct {
	FileHash string `json:"filehash"`
	Info     DataFileInfo
}

// AuthBasicInfo is definition of auth basic info
type AuthBasicInfo struct {
	AuthType string `json:"authType"`
	AuthPath string `json:"authpath"`
}

// DataFileAuthInfo is definition of data file auth info
type DataFileAuthInfo struct {
	HashUserName string `json:"hashusername"`
	Info         AuthBasicInfo
}

// DataFileAuthReq is definition of data file auth req
type DataFileAuthReq struct {
	StorgeHash string `json:"storagehash"`
	UserName   string `json:"username"`
}

// DataReqInfo is definition of data req info
type DataReqInfo struct {
	UserName    string `json:"username"`
	ReqName     string `json:"reqname"`
	ReqType     uint64 `json:"reqtype"`
	FeatureTag  uint64 `json:"featuretag"`
	SampleHash  string `json:"samplehash"`
	ExpireTime  uint64 `json:"expiretime"`
	OpType      uint32 `json:"optype"`
	TokenType   string `json:"tokenType"`
	Price       uint64 `json:"price"`
	FavoriFlag  uint32 `json:"favoriflag"`
	Description string `json:"description"`
}

// RegDataReqReq is definition of reg data req reg
type RegDataReqReq struct {
	DataReqId string `json:"datareqid"`
	Info      DataReqInfo
}

// GoodsProReq is definition of goods req
type GoodsProReq struct {
	UserName  string `json:"username"`
	OpType    uint32 `json:"optype"`
	GoodsType string `json:"goodstype"`
	GoodsId   string `json:"goodsid"`
}

// NodeClusterReg is definition for node cluster reg
type NodeClusterReg struct {
	NodeIP    string `bson:"seedip" json:"nodeIP"`
	ClusterIP string `bson:"slaveiplist" json:"clusterIP"`
	NodeUUID  string `bson:"nodeuuid" json:"uuid"`
	StorageCapacity string `bson:"capacity" json:"capacity"`
}

// NodeBaseInfo is definition for node base info
type NodeBaseInfo struct {
	NodeIp      string `json:"nodeip"`
	NodePort    string `json:"nodeport"`
	NodeAddress string `json:"nodeaddress"`
}

// NodeInfoReq is definition for node info req
type NodeInfoReq struct {
	NodeId string `json:"nodeid"`
	Info   NodeBaseInfo
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
func ParseParam(ldb *db.DBService, Param []byte, Contract string, Method string) (interface{}, error) {
	var decodedParam interface{}
	if Contract == "bottos" {
		if Method == "newaccount" {
			decodedParam = &newaccountparam{}
		} else if Method == "setdelegate" {
			decodedParam = &SetDelegateParam{}
		} else if Method == "transfer" {
			decodedParam = &transferparam{}
		} else if Method == "deploycode" {
			decodedParam = &DeployCodeParam{}
		} else if Method == "grantcredit" {
			decodedParam = &GrantCreditParam{}
		} else if Method == "cancelcredit" {
			decodedParam = &CancelCreditParam{}
		} else if Method == "transferfrom" {
			decodedParam = &TransferFromParam{}
		} else if Method == "deployabi" {
			decodedParam = &DeployAbiParam{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract, ", Method: ", Method)
			return nil, errors.New("Not supported")
		}
	} else if Contract == "usermng" {
		if Method == "reguser" {
			decodedParam = &reguser{}
		} else if Method == "userlogin" {
			decodedParam = &UserLogin{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
			return nil, errors.New("Not supported")
		}
	} else if Contract == "assetmng" {
		if Method == "assetreg" {
			decodedParam = &RegAssetReq{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
			return nil, errors.New("Not supported")
		}
	} else if Contract == "datadealmng" {
		if Method == "buydata" {
			decodedParam = &DataDealReq{}
		} else if Method == "presale" {
			decodedParam = &PresaleReq{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
			return nil, errors.New("Not supported")
		}
	} else if Contract == "datafilemng" {
		if Method == "datafilereg" {
			decodedParam = &DataFileRegReq{}
		} else if Method == "fileauthreg" {
			decodedParam = &DataFileAuthReq{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
			return nil, errors.New("Not supported")
		}
	} else if Contract == "datareqmng" {
		if Method == "datareqreg" {
			decodedParam = &RegDataReqReq{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
			return nil, errors.New("Not supported")
		}
	} else if Contract == "favoritemng" {
		if Method == "favoritepro" {
			decodedParam = &GoodsProReq{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
			return nil, errors.New("Not supported")
		}
	} else if Contract == "nodeclustermng" {
		if Method == "reg" {
			myPublicIP, err := getMyPublicIPaddr()
			if err == nil && len(myPublicIP) > 0 {
				nodeApi.SaveIpPonixToBlockchain(myPublicIP)
			}
			decodedParam = &NodeClusterReg{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
			return nil, errors.New("Not supported")
		}
	} else if Contract == "nodemng" {
		if Method == "nodeinforeg" {
			decodedParam = &NodeInfoReq{}
		} else {
			//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
			return nil, errors.New("Not supported")
		}
	} else {
		//log.Info("insertTxInfoRole:Not supported: Contract: ", Contract)
		return nil, errors.New("Not supported")
	}
	
	Abi, errval := GetAbi(ldb, Contract)
	if errval != nil {
		log.Error("Get Abi failed for Contract: ", Contract, ", Method: ", Method)
		return nil, errors.New("ParseParam: Get Abi failed")
	}
	
	if Contract == "bottos" && Method == "deploycode" {
		//p, ok := decodedParam.(DeployCodeParam)

		var tmpval = &DeployCodeParam{}
		err := abi.UnmarshalAbi(Contract, &Abi, Method, Param, tmpval)
		if err != nil {
			return nil, errors.New("ParseParam: UnmarshalAbi failed")
		}
		//if ok {
		var mgoParam = mgoDeployCodeParam{}
		mgoParam.Name = tmpval.Name
		mgoParam.VMType = tmpval.VMType
		mgoParam.VMVersion = tmpval.VMVersion
		mgoParam.ContractCode = common.BytesToHex(tmpval.ContractCode)
		return mgoParam, nil

		return nil, errors.New("Decode DeployCodeParam failed.")
	} else if Contract == "bottos" && Method == "DeployAbiParam" {
		var tmpval = &DeployAbiParam{}
		err := abi.UnmarshalAbi(Contract, &Abi, Method, Param, tmpval)
		if err != nil {
			return nil, errors.New("ParseParam: UnmarshalAbi failed")
		}
		
		var mgoParam = mgoDeployAbiParam{}
		mgoParam.Name = tmpval.Name
		mgoParam.ContractAbi = common.BytesToHex(tmpval.ContractAbi)
		return mgoParam, nil
	}
	
	err := abi.UnmarshalAbi(Contract, &Abi, Method, Param, decodedParam)

	/*if Contract == "nodeclustermng" {
		decodedParam := &NodeClusterReg{}
		Abi, errval := GetAbi(ldb, "nodeclustermng")
		if errval != nil {
			return nil, errors.New("nodeclustermng: Get Abi failed!!")
		}
		err = abi.UnmarshalAbi(Contract, &Abi, Method, Param, decodedParam)
		if err != nil {
			return nil, errors.New("UnmarshalAbi failed for contract nodeclustermng")
		}
	}*/

 
	if err != nil {
		log.Error("insertTxInfoRole: FAILED: Contract: ", Contract, ", Method: ", Method)
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
			SigAlg:     trx.SigAlg,
			Signature:  common.BytesToHex(trx.Signature),
			CreateTime: time.Now(),
		}

		decodedParam, err := ParseParam(ldb, trx.Param, newtrx.Contract, newtrx.Method)

		if err != nil {
			return err
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
	//log.Info("insertBlockInfoRole: len(oids):", len(oids), ", len(block.Transactions):", len(block.Transactions), ", block.Header.MerkleRoot: ", block.Header.MerkleRoot, " | ", common.BytesToHex(block.Header.MerkleRoot) )

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

	//log.Info("GetBalanceOp: value2 is: ", value2, ", value2.Balance is: ", value2.Balance)

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
	
	Abi, errval := GetAbi(ldb, trx.Contract)
	if errval != nil {
		log.Error("Get Abi failed for Contract: ", trx.Contract, ", Method: ", trx.Method)
		return errval
	}

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

		data := &transferparam{}
		err := abi.UnmarshalAbi(trx.Contract, &Abi, trx.Method, trx.Param, data)
		if err != nil{
			log.Error("UnmarshalAbi for contract: ", trx.Contract, ", Method: ", trx.Method, " failed!")
		}

		FromAccountName := data.From
		ToAccountName := data.To
		SrcBalanceInfo, err := GetBalanceOp(ldb, FromAccountName) //data.Value

		if err != nil {
			return err
		}

		DstBalanceInfo, err := GetBalanceOp(ldb, ToAccountName)

		if err != nil {
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
	} else if trx.Method == "newaccount" {

		data := &newaccountparam{}
		err := abi.UnmarshalAbi(trx.Contract, &Abi, trx.Method, trx.Param, data)
		if err != nil{
			log.Error("UnmarshalAbi for contract: ", trx.Contract, ", Method: ", trx.Method, " failed!")
			return err
		}

		mesgs, err := findAcountInfo(ldb, data.Name)
		if mesgs != nil {
			return nil /* Do not allow insert same account */
		}
		log.Info(err)

		NewAccount := &AccountInfo{
			ID:               oid,
			AccountName:      data.Name,
			Balance:          0,  //uint32        `bson:"bto_balance"`
			StakedBalance:    0,  //uint64        `bson:"staked_balance"`
			UnstakingBalance: "", //             `bson:"unstaking_balance"`
			PublicKey:        data.Pubkey,
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
