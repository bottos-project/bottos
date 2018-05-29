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
	"fmt"
	"sort"
	"sync"
	"unsafe"

	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"
	//trx "github.com/bottos-project/bottos/transaction"
)

type BlockChain struct {
	blockDb    *db.DBService
	roleIntf   role.RoleInterface
	blockCache *BlockChainCache
	nc         contract.NativeContractInterface

	handledBlockCB HandledBlockCallback

	genesisBlock *types.Block

	chainmu sync.RWMutex
}

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

func (bc *BlockChain) Close() {
	fmt.Println("BlockChain: Close")
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
	trxs, err := contract.NativeContractInitChain(bc.roleIntf, bc.nc)
	if err != nil {
		return err
	}
	block := types.NewBlock(header, trxs)

	// execute trxs
	for _, trx := range trxs {
		ctx := &contract.Context{RoleIntf: bc.roleIntf, Trx: trx}
		err := bc.nc.ExecuteNativeContract(ctx)
		if err != contract.ERROR_NONE {
			fmt.Println("NativeContractInitChain Error: ", trx, err)
			//return err
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

func (bc *BlockChain) RegisterHandledBlockCallback(cb HandledBlockCallback) {
	bc.handledBlockCB = cb
}

func (bc *BlockChain) GetGenesisBlock() *types.Block {
	return bc.genesisBlock
}

func (bc *BlockChain) HasBlock(hash common.Hash) bool {
	if bc.blockCache.HasBlock(hash) {
		return true
	}

	return HasBlock(bc.blockDb, hash)
}

func (bc *BlockChain) GetBlock(hash common.Hash) *types.Block {
	// cache
	block := bc.blockCache.GetBlock(hash)
	if block != nil {
		return block
	}

	return GetBlock(bc.blockDb, hash)
}

func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {
	return bc.GetBlock(hash)
}

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

func (bc *BlockChain) GetBlockHashByNumber(number uint32) common.Hash {
	return GetBlockHashByNumber(bc.blockDb, number)
}

func (bc *BlockChain) WriteBlock(block *types.Block) error {
	return WriteBlock(bc.blockDb, block)
}

func (bc *BlockChain) HeadBlockTime() uint64 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockTime
}

func (bc *BlockChain) HeadBlockNum() uint32 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockNum
}

func (bc *BlockChain) HeadBlockHash() common.Hash {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockHash
}

func (bc *BlockChain) HeadBlockDelegate() string {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.CurrentDelegate
}

func (bc *BlockChain) LastConsensusBlockNum() uint32 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastConsensusBlockNum
}

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

func (bc *BlockChain) LoadBlockDb() error {
	lastBlock := GetLastBlock(bc.blockDb)
	if lastBlock == nil {
		return fmt.Errorf("Loading block database fail, try recovering")
	}

	if lastBlock.GetNumber() == 0 {
		bc.updateChainState(lastBlock)
	}

	if bc.HeadBlockHash() != lastBlock.Hash() {
		return fmt.Errorf("Load block db fail, head block hash=%x, last block in blockdb hash=%x", bc.HeadBlockHash(), lastBlock.Hash())
	}

	fmt.Printf("Loading block database, Last block num = %v, hash = %x\n", lastBlock.GetNumber(), lastBlock.Hash())

	// TODO replay
	return nil
}

// TODO
func (bc *BlockChain) updateCoreState(block *types.Block) {

	if block.Header.Number%config.BLOCKS_PER_ROUND == 0 {
		schedule := bc.roleIntf.ElectNextTermDelegates()

		newCoreState, err := bc.roleIntf.GetCoreState()
		if err != nil {
			fmt.Errorf("Loading block database fail, try recovering")
			return
		}
		newCoreState.CurrentDelegates = schedule
		bc.roleIntf.SetCoreState(newCoreState)
		//TODO permission object
	}

}

// TODO
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
		fmt.Println("missBlocks", missBlocks)
		panic(1)
		return
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

	//oldNum := delegate.LastConfirmedBlockNum
	delegate.LastSlot = newSlot
	delegate.LastConfirmedBlockNum = block.GetNumber()
	bc.roleIntf.SetDelegate(delegate.AccountName, delegate)

	//fmt.Printf("delegate: %v, update last confirmed block num, old:%v, new:%v\n", delegate.AccountName, oldNum, delegate.LastConfirmedBlockNum)
}

// TODO
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
	//fmt.Println(lastConfirmedNums)

	consensusIndex := (100 - int(config.CONSENSUS_BLOCKS_PERCENT)) * len(delegates) / 100
	sort.Sort(lastConfirmedNums)
	//fmt.Println(lastConfirmedNums, consensusIndex)
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
	//fmt.Printf("lastBlockNum = %v, newLastConsensusBlockNum = %v\n", lastBlockNum, newLastConsensusBlockNum)

	if lastBlockNum < newLastConsensusBlockNum {
		for i := lastBlockNum + 1; i <= newLastConsensusBlockNum; i++ {
			block := bc.GetBlockByNumber(i)
			if block != nil {
				//bc.WriteBlock(block)
			} else {
				fmt.Printf("block num = %v not found\n", i)
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

func (bc *BlockChain) HandleBlock(block *types.Block) error {
	delegate, _ := bc.roleIntf.GetDelegateByAccountName(string(block.GetDelegate()))

	// TODO
	//for _, tx := range block.Transactions {
	//	trx.ApplyTransaction(tx)
	//	fmt.Println("BlockChain : Applying transactions")
	//}

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

func (bc *BlockChain) ValidateBlock(block *types.Block) error {
	prevBlockHash := block.GetPrevBlockHash()
	if prevBlockHash != bc.HeadBlockHash() {
		return fmt.Errorf("Block Prev Hash error, head block Hash = %x, block PrevBlockHash = %x", bc.HeadBlockHash(), prevBlockHash)
	}

	if block.GetNumber() != bc.HeadBlockNum()+1 {
		return fmt.Errorf("Block Number error, head block Number = %v, block Number = %v", bc.HeadBlockNum(), block.GetNumber())
	}

	// block timestamp check
	if block.GetTimestamp() <= bc.HeadBlockTime() && bc.HeadBlockNum() != 0 {
		return fmt.Errorf("Block Timestamp error, head block time=%v, block time=%v", bc.HeadBlockTime(), block.GetTimestamp())
	}

	if (block.GetTimestamp() > bc.HeadBlockTime()+uint64(config.DEFAULT_BLOCK_INTERVAL)) && bc.HeadBlockNum() != 0 {
		fmt.Printf("Block Timestamp, head block time=%v, block time=%v", bc.HeadBlockTime(), block.GetTimestamp())
		//return fmt.Errorf("Block Timestamp error, head block time=%v, block time=%v", bc.HeadBlockTime(), block.GetTimestamp())
	}

	//slot := bc.roleIntf.GetSlotAtTime(block.GetTimestamp())
	//scheduleDelegateName, _ := bc.roleIntf.GetScheduleDelegateRole(slot)
	//scheduleDelegate, _ := bc.roleIntf.GetDelegateByAccountName(scheduleDelegateName)
	// TODO delegate signature check
	if ok := block.ValidateSign( /*producer*/ ); !ok {
		return fmt.Errorf("Producer Sign Error")
	}

	// delegate schedule check
	/*
		blockDelegate := string(block.GetDelegate())
		if blockDelegate != scheduleDelegate.AccountName {
			return fmt.Errorf("Schedule Delegate Error: schedule delegate %v, block delegate %v", scheduleDelegate.AccountName, blockDelegate)
		}
	*/

	return nil
}

func (bc *BlockChain) InsertBlock(block *types.Block) error {
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	// TODO db lock

	err := bc.ValidateBlock(block)
	if err != nil {
		fmt.Println("Validate Block error: ", err)
		return err
	}

	// push to cache, block must link now, TODO: fork process
	_, err = bc.blockCache.Insert(block)
	if err != nil {
		fmt.Println("blockCache insert error: ", err)
		return err
	}

	// record stateDb revision
	//bc.stateDb.StartUndoSession(true)
	err = bc.HandleBlock(block)
	if err != nil {
		//bc.stateDb.Rollback()
		fmt.Println("InsertBlock error: ", err)
		return err
	}
	//bc.stateDb.Commit()

	fmt.Printf("Insert Block: number:%v, trxn:%v, delegate: %v\n\n", block.GetNumber(), len(block.Transactions), string(block.GetDelegate()))

	return nil
}
