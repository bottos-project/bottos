package common

import (
	"time"
	"fmt"
	"bytes"
	//"time"
	"github.com/bottos-project/core/common/types"
)

var (
	blockHashPre  = []byte("bh-")
	blockNumPre   = []byte("bn-")
)

func GetBlock(db Database, hash Hash) *types.Block {
	return nil
}

func GetBlockHashByNumber(db Database, number uint32) Hash {
	return Hash{}
}


func WriteGenesisBlock(blockDb Database) (*types.Block, error) {
	// TODO process init account
	
	// TODO make block and write to db
	header := &types.Header {
		PrevBlockHash:	Hash{},
		Number:			0,
		Timestamp:		uint32(time.Now().Unix()),
		MerkleRoot:     Hash{}, // TODO
		Producer:		AccountName{},
		ProducerChange:	[]AccountName{}, //  TODO
		ProducerSign:	Hash{}, // TODO
	}

	block := types.NewBlock(header, []*types.Transaction{})

	err := WriteBlock(blockDb, block)
	if err != nil {
		return nil, err
	}
	err = WriteHead(blockDb, block)
	if err != nil {
		return nil, err
	}
	
	return block, nil
}

func WriteChainNumber(db Database, block *types.Block) error {
	return nil
}

func WriteHead(db Database, block *types.Block) error {
	return nil
}

func WriteBlock(db Database, block *types.Block) error {
	return nil
}
