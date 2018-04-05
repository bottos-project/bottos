package common

import (
	"fmt"
	"sync"
	"github.com/bottos-project/bottos/core/common/types"
	//"github.com/bottos-project/bottos/core/define"
)

/*
ChainController---eos
ChainManager---eth1.0.3
*/

type BlockChain struct {
	//blockDB			define.Database
	//stateDB			def.Database

	chainmu			sync.RWMutex

	CurrentBlock	*types.Block
}


func CreateBlockChain(/*blockDB define.Database*/) (*BlockChain, error) {
	bc := &BlockChain{}
	//bc.blockDB = blockDB

	return bc, nil
}

/*InsertChain---eth1.0.3*/
func (self *BlockChain) PushBlock(block *types.Block) error {
	self.chainmu.Lock()
	defer self.chainmu.Unlock()

	// validate and execute
	// TODO self.blockDB.Put()
	self.CurrentBlock = block
	fmt.Println("BlockChain : push a block")
	return nil
}

func (self *BlockChain) GetBlock(hash []byte) *types.Block {
	return &types.Block{}
}

func (self *BlockChain) GetBlockByNumber(Number uint64) *types.Block {
	return &types.Block{}
}
