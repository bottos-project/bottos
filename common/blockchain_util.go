package common

import (
	"time"
	"fmt"
	"bytes"
	//"time"
	"github.com/bottos-project/bottos/core/common/types"
	"github.com/bottos-project/bottos/core/library"
	"github.com/bottos-project/bottos/core/library/rlp"
)

var (
	blockHashPre  = []byte("block-hash-")
	blockNumPre   = []byte("block-num-")
)

// GetBlockByHash returns the block corresponding to the hash or nil if not found
func GetBlock(db library.Database, hash library.Hash) *types.Block {
	data, _ := db.Get(append(blockHashPre, hash[:]...))
	if len(data) == 0 {
		return nil
	}

	fmt.Printf("GetBlockByHash, hash: %s\n", library.ToHex(hash.Bytes()))
	fmt.Printf("GetBlockByHash, data: %s\n", library.ToHex(data))

	var block types.Block
	if err := rlp.Decode(bytes.NewReader(data), &block); err != nil {
		//glog.V(logger.Error).Infof("invalid block RLP for hash %x: %v", hash, err)
		return nil
	}
	return &block
}

func GetBlockHashByNumber(db library.Database, number uint32) library.Hash {
	hash, _ := db.Get(append(blockNumPre, library.NumberToBytes(number,32)...))
	if len(hash) == 0 {
		return library.Hash{}
	}
	return library.BytesToHash(hash)
}


// WriteGenesisBlock writes the genesis block to the database as block number 0
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

// WriteCanonNumber writes the canonical hash for the given block
func WriteCanonNumber(db library.Database, block *types.Block) error {
	key := append(blockNumPre, library.NumberToBytes(block.Number(),32)...)
	err := db.Put(key, block.Hash().Bytes())
	if err != nil {
		return err
	}
	return nil
}

// WriteHead force writes the current head
func WriteHead(db library.Database, block *types.Block) error {
	err := WriteCanonNumber(db, block)
	if err != nil {
		return err
	}
	err = db.Put([]byte("LastBlock"), block.Hash().Bytes())
	if err != nil {
		return err
	}
	return nil
}

// WriteBlock writes a block to the database
func WriteBlock(db library.Database, block *types.Block) error {
	//tstart := time.Now()

	enc, _ := rlp.EncodeToBytes(block)
	key := append(blockHashPre, block.Hash().Bytes()...)

	fmt.Printf("hash: %s\n", library.ToHex(block.Hash().Bytes()))
	fmt.Printf("data: %s\n", library.ToHex(enc))

	err := db.Put(key, enc)
	if err != nil {
		//glog.Fatal("db write fail:", err)
		return err
	}

	//if glog.V(logger.Debug) {
	//	glog.Infof("wrote block #%v. Took %v\n", block.Number(), time.Since(tstart))
	//}

	return nil
}
