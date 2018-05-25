package role

import (
	"encoding/json"
	_ "fmt"
	//"time"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/db"
)

type ChainState struct {
	LastBlockNum          uint32      `json:"last_block_num"`
	LastBlockHash         common.Hash `json:"last_block_hash"`
	LastBlockTime         uint64      `json:"last_block_time"`
	LastConsensusBlockNum uint32	  `json:"last_consensus_block_num"`
	CurrentDelegate       string      `json:"current_delegate"`
	CurrentAbsoluteSlot   uint64      `json:"current_absolute_slot"`
	RecentSlotFilled      uint64      `json:"recent_slot_filled"`
}

const (
	ChainStateName       string = "chain_state"
	ChainStateDefaultKey string = "chain_state_defkey"
)

func getGenesisTime() uint64 {
	t := config.Genesis.GenesisTime
	//genesisTime := (uint32(t) / config.DEFAULT_BLOCK_INTERVAL) * config.DEFAULT_BLOCK_INTERVAL
	return uint64(t)
}

func CreateChainStateRole(ldb *db.DBService) error {
	object := &ChainState{
		LastBlockTime:    getGenesisTime(),
		RecentSlotFilled: ^uint64(0),
	}
	return SetChainStateRole(ldb, object)
}

func SetChainStateRole(ldb *db.DBService, value *ChainState) error {
	jsonvalue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return ldb.SetObject(ChainStateName, ChainStateDefaultKey, string(jsonvalue))
}

func GetChainStateRole(ldb *db.DBService) (*ChainState, error) {
	value, err := ldb.GetObject(ChainStateName, ChainStateDefaultKey)
	if err != nil {
		return nil, err
	}

	res := &ChainState{}
	json.Unmarshal([]byte(value), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
