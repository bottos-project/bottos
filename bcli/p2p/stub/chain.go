package stub

import (
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
)

//type HandledBlockCallback func(*types.Block)

type BlockChainStub struct {
	blocknumber uint32
	blocks      []types.Block
}

//func MakeBlockChainStub() chain.BlockChainInterface {
//	return &BlockChainStub{}
//}

func MakeBlockChainStub() *BlockChainStub {
	return &BlockChainStub{}
}

func (b *BlockChainStub) Close() {

}

func (b *BlockChainStub) HasBlock(hash common.Hash) bool {
	return true
}

func (b *BlockChainStub) GetBlockByHash(hash common.Hash) *types.Block {
	return nil

}
func (b *BlockChainStub) GetBlockByNumber(number uint32) *types.Block {
	for _, block := range b.blocks {
		if block.Header.Number == number {
			return &block
		}
	}

	return nil
}

func (b *BlockChainStub) HeadBlockTime() uint64 {
	return 0
}
func (b *BlockChainStub) HeadBlockNum() uint32 {
	return b.blocknumber
}
func (b *BlockChainStub) HeadBlockHash() common.Hash {
	return common.Hash{}
}
func (b *BlockChainStub) HeadBlockDelegate() string {
	return ""
}
func (b *BlockChainStub) LastConsensusBlockNum() uint32 {
	return 0
}
func (b *BlockChainStub) GenesisTimestamp() uint64 {
	return 0
}

func (b *BlockChainStub) InsertBlock(block *types.Block) error {
	return nil
}

func (b *BlockChainStub) RegisterHandledBlockCallback(cb chain.HandledBlockCallback) {
	return
}

func (b *BlockChainStub) GetHeaderByNumber(number uint32) *types.Header {
	for _, block := range b.blocks {
		if block.Header.Number == number {
			return block.Header
		}
	}

	return nil
}

func (b *BlockChainStub) SetBlockNumber(blocknumber uint32) {
	b.blocknumber = blocknumber
}

func (b *BlockChainStub) SetBlocks(blocks []types.Block) {
	b.blocks = blocks
}

func (b *BlockChainStub) Tell(message interface{}) {

}
