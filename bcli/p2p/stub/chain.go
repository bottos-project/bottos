package stub

import (
	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"sync"
)

//type HandledBlockCallback func(*types.Block)

type BlockChainStub struct {
	blocks []types.Block

	l sync.Mutex
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
	b.l.Lock()
	defer b.l.Unlock()

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
	b.l.Lock()
	defer b.l.Unlock()

	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1].Header.Number
	} else {
		return 0
	}
}
func (b *BlockChainStub) HeadBlockHash() common.Hash {
	b.l.Lock()
	defer b.l.Unlock()

	if len(b.blocks) > 0 {
		return b.blocks[len(b.blocks)-1].Header.Hash()
	} else {
		return common.Hash{}
	}

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

func (b *BlockChainStub) InsertBlock(block *types.Block) uint32 {
	b.l.Lock()
	defer b.l.Unlock()

	b.blocks = append(b.blocks, *block)
	return 0
}

func (b *BlockChainStub) RegisterHandledBlockCallback(cb chain.HandledBlockCallback) {
	return
}

func (b *BlockChainStub) GetHeaderByNumber(number uint32) *types.Header {
	b.l.Lock()
	defer b.l.Unlock()

	for _, block := range b.blocks {
		if block.Header.Number == number {
			return block.Header
		}
	}

	return nil
}

func (b *BlockChainStub) SetBlocks(blocks []types.Block) {
	b.blocks = blocks
}

func (b *BlockChainStub) Tell(message interface{}) {

}

func (b *BlockChainStub) NewBlockMsg() *message.NotifyBlock {
	last := b.blocks[len(b.blocks)-1]
	new := types.NewBlock(last.Header, last.Transactions)

	new.Header.Number++
	new.Header.PrevBlockHash = last.Header.Hash().Bytes()

	b.blocks = append(b.blocks, *new)

	msg := &message.NotifyBlock{Block: new}

	return msg
}

func (b *BlockChainStub) ValidateBlock(block *types.Block) uint32 {
	return 0
}
