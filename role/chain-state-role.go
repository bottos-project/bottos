package role

import (
	_"fmt"
	"encoding/json"
	//"time"

	"github.com/bottos-project/core/db"
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/config"
)

type ChainState struct {
	LastBlockNum				uint32				`json:"last_block_num"`
	LastBlockHash				common.Hash			`json:"last_block_hash"`
	LastBlockTime	       		uint64				`json:"last_block_time"`
	LastConfirmedBlockNum 		uint32				`json:"last_confirmed_block_num"`
	CurrentDelegate				string				`json:"current_delegate"`
	CurrentAbsoluteSlot			uint64				`json:"current_absolute_slot"`
	RecentSlotFilled			uint64				`json:"recent_slot_filled"`
}

const (
	ChainStateName string = "chain_state"
	ChainStateDefaultKey string = "chain_state_defkey"
)

func getGenesisTime() uint64 {
	t := config.Genesis.GenesisTime
	//genesisTime := (uint32(t) / config.DEFAULT_BLOCK_INTERVAL) * config.DEFAULT_BLOCK_INTERVAL
	return uint64(t)
}

func CreateChainStateRole(ldb *db.DBService) error {
	object := &ChainState{
		LastBlockTime: getGenesisTime(),
	}
	return SetChainStateRole(ldb, object)
}

func SetChainStateRole(ldb *db.DBService, value *ChainState) error {
	jsonvalue, _ := json.Marshal(value)
	//mt.Println("Set", ChainStateObjectDefaultKey, value)
	return ldb.SetObject(ChainStateName, ChainStateDefaultKey, string(jsonvalue))
}

func GetChainStateRole(ldb *db.DBService) (*ChainState, error) {
	value, err := ldb.GetObject(ChainStateName, ChainStateDefaultKey)
	res := &ChainState{}
	json.Unmarshal([]byte(value), res)
	//fmt.Println("Get", ChainStateObjectDefaultKey, value)
	return res, err
}
