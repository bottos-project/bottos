package role

import (
	"fmt"
	"encoding/json"
	"time"

	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/config"
)

type ChainStateObject struct {
	LastBlockNum				uint32				`json:"last_block_num"`
	LastBlockHash				common.Hash			`json:"last_block_hash"`
	LastBlockTime	       		uint64				`json:"last_block_time"`
	LastConfirmedBlockNum 		uint64				`json:"last_confirmed_block_num"`
	CurrentDelegate				string				`json:"current_delegate"`
	CurrentAbsoluteSlot			uint64				`json:"current_absolute_slot"`
	RecentSlotFilled			uint64				`json:"recent_slot_filled"`
}

const (
	ChainStateObjectName string = "chain_state_object"
	ChainStateObjectDefaultKey string = "chain_state_object_defkey"
)

func getGenesisTime() uint64 {
	t := time.Now().Unix()
	genesisTime := (uint32(t) / config.DEFAULT_BLOCK_INTERVAL) * config.DEFAULT_BLOCK_INTERVAL
	return uint64(genesisTime)
}

func CreateChainStateObjectRole(ldb *db.DBService) error {
	object := &ChainStateObject{
		LastBlockTime: getGenesisTime(),
	}
	return SetChainStateObjectRole(ldb, object)
}

func SetChainStateObjectRole(ldb *db.DBService, value *ChainStateObject) error {
	jsonvalue, _ := json.Marshal(value)
	return ldb.SetObject(ChainStateObjectName, ChainStateObjectDefaultKey, string(jsonvalue))
}

func GetChainStateObjectRole(ldb *db.DBService) (*ChainStateObject, error) {
	value, err := ldb.GetObject(ChainStateObjectName, ChainStateObjectDefaultKey)
	res := &ChainStateObject{}
	json.Unmarshal([]byte(value), res)
	fmt.Println("Get", ChainStateObjectDefaultKey, value)
	return res, err
}
