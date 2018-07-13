// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  blockchain general interface and logic
 * @Author: Gong Zibin
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package chain

import (
	"sort"
	"sync"
	"unsafe"

	log "github.com/cihub/seelog"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"
)

//BlockChain the chain info
type BlockChain struct {
	blockDb    *db.DBService
	roleIntf   role.RoleInterface
	blockCache *BlockChainCache
	nc         contract.NativeContractInterface

	handledBlockCB HandledBlockCallback

	genesisBlock *types.Block

	chainmu sync.RWMutex
}

//CreateBlockChain create a chain
func CreateBlockChain(dbInstance *db.DBService, roleIntf role.RoleInterface, nc contract.NativeContractInterface) (BlockChainInterface, error) {
	blockCache, err := CreateBlockChainCache()
	if err != nil {
		return nil, err
	}

	bc := &BlockChain{
		blockDb:    dbInstance,
		blockCache: blockCache,
		roleIntf:   roleIntf,
		nc:         nc,
	}

	err = bc.initChain()
	if err != nil {
		return nil, err
	}

	err = bc.LoadBlockDb()
	if err != nil {
		return nil, err
	}

	// init block cache
	bc.initBlockCache()

	return bc, nil
}

//Close close chain cache
func (bc *BlockChain) Close() {
	log.Info("BlockChain: Close")
	bc.blockCache.Reset()
}

func (bc *BlockChain) initChain() error {
	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock != nil {
		return nil
	}

	header := &types.Header{
		Version:   1,
		Number:    0,
		Timestamp: config.Genesis.GenesisTime,
		Delegate:  []byte(config.BOTTOS_CONTRACT_NAME),
	}
	trxs, err := contract.NativeContractInitChain(bc.blockDb, bc.roleIntf, bc.nc)
	if err != nil {
		return err
	}
	block := types.NewBlock(header, trxs)

	// execute trxs
	for _, trx := range trxs {
		ctx := &contract.Context{RoleIntf: bc.roleIntf, Trx: trx}
		err := bc.nc.ExecuteNativeContract(ctx)
		if err != contract.ERROR_NONE {
			log.Infof("NativeContractInitChain Error: ", trx, err)
			break
		}
	}

	err = WriteGenesisBlock(bc.blockDb, block)
	if err != nil {
		return err
	}

	bc.genesisBlock = block
	bc.roleIntf.ApplyPersistance(block)

	return nil
}

//RegisterHandledBlockCallback call back register
func (bc *BlockChain) RegisterHandledBlockCallback(cb HandledBlockCallback) {
	bc.handledBlockCB = cb
}

//GetGenesisBlock get the first block
func (bc *BlockChain) GetGenesisBlock() *types.Block {
	return bc.genesisBlock
}

//HasBlock check block
func (bc *BlockChain) HasBlock(hash common.Hash) bool {
	if bc.blockCache.HasBlock(hash) {
		return true
	}

	return HasBlock(bc.blockDb, hash)
}

//GetBlock get block from cache and chain by hash
func (bc *BlockChain) GetBlock(hash common.Hash) *types.Block {
	// cache
	block := bc.blockCache.GetBlock(hash)
	if block != nil {
		return block
	}

	return GetBlock(bc.blockDb, hash)
}

//GetBlockByHash get block from chain by hash
func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {
	return bc.GetBlock(hash)
}

//GetBlockByNumber get block from cache and chain by number
func (bc *BlockChain) GetBlockByNumber(number uint32) *types.Block {
	block := bc.blockCache.GetBlockByNum(number)
	if block != nil {
		return block
	}

	hash := GetBlockHashByNumber(bc.blockDb, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return bc.GetBlock(hash)
}

func (bc *BlockChain) GetHeaderByNumber(number uint32) *types.Header {
	block := bc.GetBlockByNumber(number)
	if block != nil {
		return block.Header
	}

	return nil
}

//GetBlockHashByNumber get block hash from chain by number
func (bc *BlockChain) GetBlockHashByNumber(number uint32) common.Hash {
	return GetBlockHashByNumber(bc.blockDb, number)
}

//WriteBlock write block to chain
func (bc *BlockChain) WriteBlock(block *types.Block) error {
	return WriteBlock(bc.blockDb, block)
}

//HeadBlockTime get lastest block time
func (bc *BlockChain) HeadBlockTime() uint64 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockTime
}

//HeadBlockNum get lastest block number
func (bc *BlockChain) HeadBlockNum() uint32 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockNum
}

//HeadBlockHash get lastest block hash
func (bc *BlockChain) HeadBlockHash() common.Hash {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockHash
}

//HeadBlockDelegate get current delegator
func (bc *BlockChain) HeadBlockDelegate() string {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.CurrentDelegate
}

//LastConsensusBlockNum get lastest consensus block numnber
func (bc *BlockChain) LastConsensusBlockNum() uint32 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastConsensusBlockNum
}

//GenesisTimestamp get genesis time
func (bc *BlockChain) GenesisTimestamp() uint64 {
	return config.Genesis.GenesisTime
}

// internal
func (bc *BlockChain) getBlockDbLastBlock() *types.Block {
	return GetLastBlock(bc.blockDb)
}

func (bc *BlockChain) initBlockCache() error {
	block := bc.getBlockDbLastBlock()
	if block != nil {
		_, err := bc.blockCache.Insert(block)
		return err
	}

	return nil
}

//LoadBlockDb load block from db
func (bc *BlockChain) LoadBlockDb() error {
	lastBlock := GetLastBlock(bc.blockDb)
	if lastBlock == nil {
		return log.Error("Loading block database fail, try recovering")
	}

	if lastBlock.GetNumber() == 0 {
		bc.updateChainState(lastBlock)
	}

	if bc.HeadBlockHash() != lastBlock.Hash() {
		return log.Errorf("Load block db fail, head block hash=%x, last block in blockdb hash=%x", bc.HeadBlockHash(), lastBlock.Hash())
	}

	log.Infof("Loading block database, Last block num = %v, hash = %x\n", lastBlock.GetNumber(), lastBlock.Hash())

	return nil
}

func (bc *BlockChain) updateCoreState(block *types.Block) {
	if block.Header.Number%config.BLOCKS_PER_ROUND == 0 {
		schedule, err := bc.roleIntf.ShuffleEelectCandidateList(block)
		if err != nil {
			return
		}
		newCoreState, err := bc.roleIntf.GetCoreState()
		if err != nil {
			return
		}
		if uint32(len(schedule)) > config.BLOCKS_PER_ROUND {
			log.Error("invalid schedule length which is greater than BLOCKS_PER_ROUND")
			return
		}

		copy(newCoreState.CurrentDelegates, schedule)
		//log.Info("CurrentDelegates", newCoreState.CurrentDelegates)
		bc.roleIntf.SetCoreState(newCoreState)
	}
}

func (bc *BlockChain) updateChainState(block *types.Block) {
	var missBlocks uint64
	var i uint64
	chainSate, _ := bc.roleIntf.GetChainState()

	chainSate.CurrentDelegate = string(block.GetDelegate())

	if chainSate.LastBlockNum == 0 {
		missBlocks = 1
	} else {
		slot := bc.roleIntf.GetSlotAtTime(block.GetTimestamp())
		missBlocks = slot
	}
	if missBlocks == 0 {
		log.Infof("missBlocks", missBlocks)
		panic(1)
	}
	missBlocks--

	for i = 0; i < missBlocks; i++ {
		name, err := bc.roleIntf.GetCandidateBySlot(i + 1)

		delegateLeave, err := bc.roleIntf.GetDelegateByAccountName(name)
		if err != nil {
			continue
		}
		if delegateLeave.AccountName != chainSate.CurrentDelegate {
			delegateLeave.TotalMissed++
			bc.roleIntf.SetDelegate(delegateLeave.AccountName, delegateLeave)
		}
	}

	//update chain state
	chainSate.CurrentAbsoluteSlot += missBlocks + 1

	size := uint64(unsafe.Sizeof(chainSate.RecentSlotFilled))
	if missBlocks < size*8 {
		chainSate.RecentSlotFilled <<= 1
		chainSate.RecentSlotFilled += 1
		chainSate.RecentSlotFilled <<= missBlocks
	} else {
		coreSate, _ := bc.roleIntf.GetCoreState()

		if uint64(uint32(len(coreSate.CurrentDelegates))/config.BLOCKS_PER_ROUND) > config.DELEGATE_PATICIPATION {

			chainSate.RecentSlotFilled = ^uint64(0)
		} else {
			chainSate.RecentSlotFilled = 0
		}
	}

	chainSate.LastBlockNum = block.GetNumber()
	chainSate.LastBlockHash = block.Hash()
	chainSate.LastBlockTime = block.GetTimestamp()

	bc.roleIntf.SetChainState(chainSate)

}

func (bc *BlockChain) updateDelegate(delegate *role.Delegate, block *types.Block) {
	chainSate, _ := bc.roleIntf.GetChainState()
	blockTime := block.GetTimestamp()
	newSlot := chainSate.CurrentAbsoluteSlot + uint64(bc.roleIntf.GetSlotAtTime(blockTime))
	delegate.LastSlot = newSlot
	delegate.LastConfirmedBlockNum = block.GetNumber()
	bc.roleIntf.SetDelegate(delegate.AccountName, delegate)
}

func (bc *BlockChain) updateConsensusBlock(block *types.Block) {
	chainSate, _ := bc.roleIntf.GetChainState()
	coreState, _ := bc.roleIntf.GetCoreState()

	delegates := make([]*role.Delegate, len(coreState.CurrentDelegates))
	lastConfirmedNums := make(ConfirmedNum, len(coreState.CurrentDelegates))
	for i, name := range coreState.CurrentDelegates {
		delegate, _ := bc.roleIntf.GetDelegateByAccountName(name)
		delegates[i] = delegate
		lastConfirmedNums[i] = delegates[i].LastConfirmedBlockNum
	}

	consensusIndex := (100 - int(config.CONSENSUS_BLOCKS_PERCENT)) * len(delegates) / 100
	sort.Sort(lastConfirmedNums)
	if consensusIndex >= len(lastConfirmedNums) {
		log.Errorf("out of range: index=%v, len=%v", consensusIndex, len(lastConfirmedNums))
		return
	}
	newLastConsensusBlockNum := lastConfirmedNums[consensusIndex]
	if newLastConsensusBlockNum > chainSate.LastConsensusBlockNum {
		chainSate.LastConsensusBlockNum = newLastConsensusBlockNum
	}
	bc.roleIntf.SetChainState(chainSate)

	// write LCB to blockDb
	lastBlockNum := uint32(0)
	lastBlock := bc.getBlockDbLastBlock()
	if lastBlock != nil {
		lastBlockNum = lastBlock.GetNumber()
	}

	if lastBlockNum < newLastConsensusBlockNum {
		for i := lastBlockNum + 1; i <= newLastConsensusBlockNum; i++ {
			block := bc.GetBlockByNumber(i)
			if block != nil {
				bc.WriteBlock(block)
			} else {
				log.Errorf("block num = %v not found\n", i)
			}
		}

		// trim blockCache
		bc.blockCache.Trim(chainSate.LastBlockNum, newLastConsensusBlockNum)
	}
}

func (bc *BlockChain) clearTransactionExpiration(block *types.Block) error {
	// TODO
	return nil
}

func (bc *BlockChain) addBlockHistory(block *types.Block) {
	bc.roleIntf.SetBlockHistory(block.GetNumber(), block.Hash())
}

//HandleBlock update state when handle a new block
func (bc *BlockChain) HandleBlock(block *types.Block) error {
	delegate, _ := bc.roleIntf.GetDelegateByAccountName(string(block.GetDelegate()))

	// update consensus
	bc.updateCoreState(block)
	bc.updateChainState(block)
	bc.updateDelegate(delegate, block)
	bc.updateConsensusBlock(block)

	// clear transaction expiration
	bc.clearTransactionExpiration(block)

	bc.addBlockHistory(block)

	// block handled callback
	if bc.handledBlockCB != nil {
		bc.handledBlockCB(block)
	}
	bc.WriteBlock(block)
	bc.roleIntf.ApplyPersistance(block)
	return nil
}

//ValidateBlock verify a block
func (bc *BlockChain) ValidateBlock(block *types.Block) uint32 {
	prevBlockHash := block.GetPrevBlockHash()
	if prevBlockHash != bc.HeadBlockHash() {
		log.Errorf("Block Prev Hash error, head block Hash = %x, block PrevBlockHash = %x", bc.HeadBlockHash(), prevBlockHash)
		return InsertBlockErrorValidateFail
	}

	if block.GetNumber() != bc.HeadBlockNum()+1 {
		log.Errorf("Block Number error, head block Number = %v, block Number = %v", bc.HeadBlockNum(), block.GetNumber())
		return InsertBlockErrorValidateFail
	}

	// block timestamp check
	if block.GetTimestamp() <= bc.HeadBlockTime() && bc.HeadBlockNum() != 0 {
		log.Errorf("Block Timestamp error, head block time=%v, block time=%v", bc.HeadBlockTime(), block.GetTimestamp())
		return InsertBlockErrorValidateFail
	}

	return InsertBlockSuccess
}

func (bc *BlockChain) checkConsensusedBlock(block *types.Block) uint32 {
	num := block.GetNumber()
	localBlock := bc.GetBlockByNumber(num)
	if localBlock.Hash() == block.Hash() {
		return InsertBlockSuccess
	} else {
		if localBlock.GetPrevBlockHash() == block.GetPrevBlockHash() {
			return InsertBlockErrorDiffLibLinked
		} else {
			return InsertBlockErrorDiffLibNotLinked
		}
	}
}

//InsertBlock write a new block
func (bc *BlockChain) InsertBlock(block *types.Block) uint32 {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	if block.GetNumber() <= bc.LastConsensusBlockNum() {
		return bc.checkConsensusedBlock(block)
	}

	errcode := bc.ValidateBlock(block)
	if errcode != InsertBlockSuccess {
		return errcode
	}

	// push to cache, block must link now
	_, err := bc.blockCache.Insert(block)
	if err != nil {
		log.Infof("blockCache insert error: ", err)
		return InsertBlockErrorNotLinked
	}

	err = bc.HandleBlock(block)
	if err != nil {
		log.Infof("InsertBlock error: ", err)
		return InsertBlockErrorGeneral
	}

	log.Infof("Insert block: block num:%v, trxn:%v, delegate: %v, hash:%x\n\n", block.GetNumber(), len(block.Transactions), string(block.GetDelegate()), block.Hash())

	return InsertBlockSuccess
}
