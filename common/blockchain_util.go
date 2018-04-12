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

func GetBlock(db library.Database, hash library.Hash) *types.Block {
	return nil
}

func GetBlockHashByNumber(db library.Database, number uint32) library.Hash {
	return library.Hash{}
}


func WriteGenesisBlock(blockDb library.Database) (*types.Block, error) {
	// TODO process init account
	
	// TODO make block and write to db
	header := &types.Header {
		PrevBlockHash:	library.Hash{},
		Number:			0,
		Timestamp:		uint32(time.Now().Unix()),
		MerkleRoot:     library.Hash{}, // TODO
		Producer:		library.AccountName{},
		ProducerChange:	[]library.AccountName{}, //  TODO
		ProducerSign:	library.Hash{}, // TODO
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

func WriteChainNumber(db library.Database, block *types.Block) error {
	return nil
}

func WriteHead(db library.Database, block *types.Block) error {
	return nil
}

func WriteBlock(db library.Database, block *types.Block) error {
	return nil
}
