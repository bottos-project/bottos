package common

import (
	"fmt"
	"sync"
	"github.com/bottos-project/core/common/types"

	"github.com/hashicorp/golang-lru"
)

const (
	blockCacheLimit     = 256
	checkpointLimit     = 200
)

type BlockChain struct {
	blockDb			library.Database
	//stateDb		library.Database
	//extraDb		library.Database

	genesisBlock 	*types.Block

	mu				sync.RWMutex
	chainmu			sync.RWMutex

	checkpoint      int // TODO
	currentBlock	*types.Block
	lastBlockHash   library.Hash

	cache           *lru.Cache
}

func CreateBlockChain(blockDb library.Database) (*BlockChain, error) {
	cache, _ := lru.New(blockCacheLimit)
	bc := &BlockChain{
		blockDb:  blockDb,
		//stateDb:  stateDb,
		//extraDb:  extraDb,
		cache:    cache,
	}

	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock == nil {
		var err error
		bc.genesisBlock, err = WriteGenesisBlock(blockDb)
		if err != nil {
			return nil, err
		}
	}

	err := bc.loadLastState()
	if err != nil {
		return nil, err
	}

	bc.makeCache()

	return bc, nil
}

func (bc *BlockChain) GetGenesisBlock() *types.Block {
	return bc.genesisBlock
}

func (bc *BlockChain) HasBlock(hash library.Hash) bool {
	if bc.cache.Contains(hash) {
		return true
	}

	data, _ := bc.blockDb.Get(append(blockHashPre, hash[:]...))
	return len(data) != 0
}

func (bc *BlockChain) GetBlock(hash library.Hash) *types.Block {
	if block, ok := bc.cache.Get(hash); ok {
		return block.(*types.Block)
	}
	block := GetBlock(bc.blockDb, hash)
	if block == nil {
		return nil
	}
	bc.cache.Add(block.Hash(), block)
	return block
}

func (bc *BlockChain) GetBlockByHash(hash library.Hash) *types.Block {
	return bc.GetBlock(hash)
}

func (bc *BlockChain) GetBlockByNumber(number uint32) *types.Block {
	hash := GetBlockHashByNumber(bc.blockDb, number)
	if hash == (library.Hash{}) {
		return nil
	}
	return bc.GetBlock(hash)
}

func (bc *BlockChain) GetBlockHashByNumber(number uint32) library.Hash {
	return GetBlockHashByNumber(bc.blockDb, number)
}


// internal
func (bc *BlockChain) loadLastState() error {
	data, _ := bc.blockDb.Get([]byte("LastBlock"))
	if len(data) != 0 {
		block := bc.GetBlockByHash(library.BytesToHash(data))
		if block != nil {
			bc.currentBlock = block
			bc.lastBlockHash = block.Hash()

			fmt.Printf("current block num = %v\n", bc.currentBlock.Number())
		} else {
			if bc.recover() {
			} else {
			}
		}
	} else {
		bc.Reset()
	}

	return nil
}

func (bc *BlockChain) recover() bool{
	fmt.Println("TODO: BlockChain recover")
	return true
}

func (bc *BlockChain) Reset() {
	fmt.Println("TODO: BlockChain reset")
}

func (bc *BlockChain) insert(block *types.Block) {
	err := WriteHead(bc.blockDb, block)
	if err != nil {
	}

	bc.currentBlock = block
	bc.lastBlockHash = block.Hash()
}


func (bc *BlockChain) makeCache() {
	bc.cache, _ = lru.New(blockCacheLimit)
	bc.cache.Add(bc.genesisBlock.Hash(), bc.genesisBlock)
	for _, block := range bc.GetBlocksFromHash(bc.currentBlock.Hash(), blockCacheLimit) {
		bc.cache.Add(block.Hash(), block)
	}
}

func (bc *BlockChain) GetBlocksFromHash(hash library.Hash, n int) (blocks []*types.Block) {
	for i := 0; i < n; i++ {
		block := bc.GetBlockByHash(hash)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
		hash = block.PrevBlockHash()
	}
	return
}


func (bc *BlockChain) ApplyBlock(block *types.Block) error {
	fmt.Println("Apply block")

	return nil
}

func (bc *BlockChain) PushBlock(block *types.Block) error {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	fmt.Println("BlockChain : push block")

	// TODO: validate

	// TODO record stateDb revision
	err := bc.ApplyBlock(block)
	if (err != nil) {
		// TODO restore stateDb revision
	}

	// TODO commit stateDb revision

	// TODO self.blockDB.Put()
	bc.currentBlock = block
	
	return nil
}


