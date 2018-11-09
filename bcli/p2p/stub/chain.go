package stub

import (
	"sync"

	"github.com/bottos-project/bottos/action/message"
	"github.com/bottos-project/bottos/chain"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"gopkg.in/urfave/cli.v1"
)

//type HandledBlockCallback func(*types.Block)

type BlockChainStub struct {
	blocks []types.Block

	beginNumber uint32
	libNumber   uint32
	l           sync.Mutex
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
func (b *BlockChainStub) GetBlockByNumber(number uint64) *types.Block {
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
func (b *BlockChainStub) HeadBlockNum() uint64 {
	b.l.Lock()
	defer b.l.Unlock()

	return uint64(len(b.blocks))

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

func (b *BlockChainStub) LastConsensusBlockNum() uint64 {
	return uint64(b.libNumber)
}

func (b *BlockChainStub) GenesisTimestamp() uint64 {
	return 0
}

func (b *BlockChainStub) InsertBlock(new *types.Block) errors.ErrCode {
	b.l.Lock()
	defer b.l.Unlock()

	for _, block := range b.blocks {
		if block.Header.Number == new.Header.Number {
			return 0
		}
	}

	b.blocks = append(b.blocks, *new)

	if len(b.blocks) > 100 {
		b.libNumber = 100
	}

	return 0
}

func (b *BlockChainStub) RegisterHandledBlockCallback(cb chain.BlockCallback) {
	return
}
func (b *BlockChainStub) RegisterCommittedBlockCallback(cb chain.BlockCallback) {
	return
}
func (b *BlockChainStub) GetHeaderByNumber(number uint64) *types.Header {
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

	b.beginNumber = uint64(len(b.blocks))
}

func (b *BlockChainStub) SetLibNumber(number uint64) {
	b.libNumber = number
}

func (b *BlockChainStub) Tell(message interface{}) {

}

func (b *BlockChainStub) NewBlockMsg() *message.NotifyBlock {
	last := b.blocks[len(b.blocks)-1]
	new := types.NewBlock(last.Header, last.Transactions)

	new.Header.Number++
	new.Header.PrevBlockHash = last.Header.Hash().Bytes()

	b.InsertBlock(new)

	msg := &message.NotifyBlock{Block: new}

	return msg
}

func (b *BlockChainStub) ValidateBlock(block *types.Block) uint32 {
	return 0
}
