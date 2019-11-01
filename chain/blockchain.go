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

	"github.com/bottos-project/bottos/cmd"
	"github.com/bottos-project/bottos/common"
	berr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/db"
	"github.com/bottos-project/bottos/role"
)

//BlockChain the chain info
type BlockChain struct {
	dbInst   *db.DBService
	blockDb  *BlockDB
	roleIntf role.RoleInterface
	forkdb   *ForkDB
	nc       contract.NativeContractInterface
	trxPool  *transaction.TrxPool

	handledBlockCB   []BlockCallback
	committedBlockCB []BlockCallback

	chainmu sync.RWMutex
}

//CreateBlockChain create a chain
func CreateBlockChain(datadir string, db *db.DBService, roleIntf role.RoleInterface, nc contract.NativeContractInterface) (BlockChainInterface, error) {
	forkDBPath := filepath.Join(datadir, "data/forkdb")
	forkdb, err := CreateForkDB(forkDBPath)
	if err != nil {
		return nil, err
	}

	bc := &BlockChain{
		dbInst:   db,
		blockDb:  NewBlockDB(db),
		forkdb:   forkdb,
		roleIntf: roleIntf,
		nc:       nc,
		trxPool:  nil,
	}

	return bc, nil
}

func (bc *BlockChain) Init(ctx *cli.Context) error {
	bc.dbInst.Lock()
	defer bc.dbInst.UnLock()

	gsb := bc.blockDb.GetBlockByNumber(0)
	if gsb == nil {
		err := bc.initNewChain(ctx)
		return err
	}

	if err := bc.loadDB(ctx); err != nil {
		return err
	}

	return nil
}

func (bc *BlockChain) SetTrxPool(trxPool *transaction.TrxPool) {
	bc.dbInst.Lock()
	defer bc.dbInst.UnLock()

	bc.trxPool = trxPool
}


//Close close chain cache
func (bc *BlockChain) Close() {
	bc.dbInst.Lock()
	defer bc.dbInst.UnLock()

	bc.forkdb.Close()
	log.Info("CHAIN closed")
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
		if err != berr.ErrNoError {
			log.Infof("NativeContractInit Error: ", trx, err)
			return fmt.Errorf("genesis block execute fail")
		}
	}

	err = WriteGenesisBlock(bc.blockDb, block)
	if err != nil {
		return err
	}
	bc.genesisBlock = block

	bc.handledBlockCallback(block)
	return nil
}


func (bc *BlockChain) InitOnRecover(ctx *cli.Context) error {
	gsconfig, err := bc.blockDb.GetGenesisConfig()
	if err != nil {
		return fmt.Errorf("CHAIN InitOnRecover fail, genesis configuration not found")
	}
	config.SetGenesisConfig(gsconfig)

	err = contract.NativeContractInitChain(bc.dbInst, bc.roleIntf, bc.nc)
	if err != nil {
		return err
	}

	b := bc.GetBlockByNumber(uint64(0))

	bc.updateChainState(b)

	return nil
}

func (bc *BlockChain) loadDB(ctx *cli.Context) error {
	gsconfig, err := bc.blockDb.GetGenesisConfig()
	if err != nil {
		return fmt.Errorf("CHAIN Loading block database fail, genesis configuration not found")
	}
	config.SetGenesisConfig(gsconfig)

	if ctx.GlobalIsSet(cmd.GenesisFileFlag.Name) {
		log.Info("CHAIN Ignoring genesis configuration path, already have genesis configuration within block database")
	}

	lastb := bc.blockDb.GetLastBlock()
	if lastb == nil {
		return fmt.Errorf("CHAIN Loading block database fail, last block not found")
	}

	err = bc.dbInst.RollbackAll()
	if err != nil {
		return err
	}

	headBlockNum := bc.HeadBlockNum()
	libNum := bc.LastConsensusBlockNum()
	lastBlockNum := lastb.GetNumber()
	log.Errorf("CHAIN Loading block database, headBlockNum %v, libNum %v, lastBlockNum %v", headBlockNum, libNum, lastBlockNum)

	if bc.HeadBlockNum() != lastBlockNum || bc.LastConsensusBlockNum() != lastBlockNum {
		return fmt.Errorf("CHAIN Head block number not match")
	}

	_, errcode := bc.forkdb.Insert(lastb)
	if errcode != berr.ErrNoError {
		return fmt.Errorf("CHAIN initialize forkdb fail")
	}

	if lastb.GetNumber() == 0 {
		bc.updateChainState(lastb)
	}
	//release undo information after rollback success
	bc.dbInst.ReleaseUndoInfo()
	return nil
}

func (bc *BlockChain) makeGenesisBlock() *types.Block {
	header := &types.Header{
		Version: 1,
		Number:  0,
		//PrevBlockHash: config.GetChainID(),
		Timestamp: config.Genesis.GenesisTime,
		Delegate:  []byte(config.BOTTOS_CONTRACT_NAME),
	}

	var trxs []*types.BlockTransaction
	block := types.NewBlock(header, trxs)

	return block
}

//RegisterHandledBlockCallback call back register
func (bc *BlockChain) RegisterHandledBlockCallback(cb BlockCallback) {
	bc.handledBlockCB = append(bc.handledBlockCB, cb)
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

//GetBlock get block from forkdb main fork and blockdb by hash
func (bc *BlockChain) GetBlock(hash common.Hash) *types.Block {
	// cache
	block := bc.forkdb.GetMainForkBlock(hash)
	if block != nil {
		return block
	}

	return bc.blockDb.GetBlock(hash)
}


//GetBlockByHash get block from chain by hash
func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {
	return bc.GetBlock(hash)
}

//GetCommittedTransaction get committed transaction from blockdb by hash
func (bc *BlockChain) GetCommittedTransaction(hash common.Hash) *types.BlockTransaction {
	return bc.blockDb.GetTransaction(hash)
}


//GetBlockByNumber get block from forkdb main fork and blockdb by number
func (bc *BlockChain) GetBlockByNumber(number uint64) *types.Block {
	block := bc.forkdb.GetMainForkBlockByNum(number)
	if block != nil {
		return block
	}

	hash := bc.blockDb.GetBlockHashByNumber(number)
	if hash == (common.Hash{}) {
		log.Errorf("CHAIN GetBlockHashByNumber fail, number = %v", number)
		return nil
	}
	return bc.blockDb.GetBlock(hash)
}

func (bc *BlockChain) GetHeaderByNumber(number uint64) *types.Header {
	block := bc.GetBlockByNumber(number)
	if block != nil {
		return block.Header
	}

	return nil
}

//GetBlockByHash get block in forkdb(no matter which forks, linked or unlinked) and blockdb
func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {
	block := bc.forkdb.GetBlock(hash)
	if block != nil {
		return block
	}

	return bc.blockDb.GetBlock(hash)
}

//GetBlockHashByNumber get block hash from chain by number
func (bc *BlockChain) GetBlockHashByNumber(number uint64) common.Hash {
	return bc.blockDb.GetBlockHashByNumber(number)
}

//WriteBlock write block to blockdb
func (bc *BlockChain) WriteBlock(block *types.Block) error {
	err := bc.blockDb.WriteBlock(block)
	return err
}

func (bc *BlockChain) handledBlockCallback(block *types.Block) {
	for _, cb := range bc.handledBlockCB {
		if cb != nil {
			cb(block)
		}
	}
}

func (bc *BlockChain) committedBlockCallback(block *types.Block) {
	for _, cb := range bc.committedBlockCB {
		if cb != nil {
			cb(block)
		}
	}
}

//HeadBlockTime get head block time
func (bc *BlockChain) HeadBlockTime() uint64 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockTime
}

//HeadBlockNum get head block number
func (bc *BlockChain) HeadBlockNum() uint64 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockNum
}

//HeadBlockHash get head block hash
func (bc *BlockChain) HeadBlockHash() common.Hash {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastBlockHash
}

//HeadBlockDelegate get delegate of head block
func (bc *BlockChain) HeadBlockDelegate() string {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.CurrentDelegate
}

//LastConsensusBlockNum get lib block num
func (bc *BlockChain) LastConsensusBlockNum() uint64 {
	coreState, _ := bc.roleIntf.GetChainState()
	return coreState.LastConsensusBlockNum
}

func (bc *BlockChain) checkConsensusedBlock(block *types.Block) berr.ErrCode {
	num := block.GetNumber()
	localBlock := bc.GetBlockByNumber(num)
	if localBlock.Hash() == block.Hash() {
		return berr.ErrNoError
	} else {
		if localBlock.GetPrevBlockHash() == block.GetPrevBlockHash() {
			return berr.ErrBlockInsertErrorDiffLibLinked
		} else {
			return berr.ErrBlockInsertErrorDiffLibNotLinked
		}
	}
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

//InsertBlock insert the latest block to block chain.
func (bc *BlockChain) InsertBlock(block *types.Block) berr.ErrCode {
	bc.dbInst.Lock()
	defer bc.dbInst.UnLock()

	if block.GetNumber() <= bc.forkdb.GetStartBlockNum() {
		return bc.checkConsensusedBlock(block)
	}

	if err := version.CheckBlock(block, "InsertBlock"); err != berr.ErrNoError {
		return err
	}

	newHeadBlock, errcode := bc.forkdb.Insert(block)
	if errcode != berr.ErrNoError {
		log.Errorf("ForkDB insert error: %v", errcode)
		return errcode
	}

	session := bc.dbInst.GetSession()
	if session != nil {
		bc.dbInst.ResetSession()
	}

	if newHeadBlock.GetPrevBlockHash() != bc.HeadBlockHash() {
		return bc.processSwitchFork(newHeadBlock)
	}

	session = bc.dbInst.BeginUndo(config.PRIMARY_TRX_SESSION)
	errcode = bc.handleBlock(block, false)
	if errcode != berr.ErrNoError {
		bc.dbInst.ResetSession()
		bc.forkdb.Remove(block)
		return errcode
	}

	chainSate, _ := bc.roleIntf.GetChainState()
	lastLib := chainSate.LastConsensusBlockNum
	validated := bc.CheckAcceptedValidators(block)
	if validated {
		bc.updateLib(block)
		bc.updateProduceTransfering(block)
	}
	bc.dbInst.Push(session)
	if validated {
		bc.commitLib(block, lastLib)
	}
	bc.handledBlockCallback(block)

	fmt.Printf("InsertBlock, number:%v, time:%v, delegate:%v, trxn:%v, hash:%x, prevHash:%x, version:%v\n",
		block.GetNumber(), common.TimeFormat(block.GetTimestamp()), string(block.GetDelegate()), len(block.BlockTransactions), block.Hash(), block.GetPrevBlockHash(), version.GetStringVersion(block.GetVersion()))

	return berr.ErrNoError
}
func (bc *BlockChain) ImportBlock(block *types.Block) berr.ErrCode {
	bc.dbInst.Lock()
	defer bc.dbInst.UnLock()

	errcode := bc.handleBlock(block, true)
	if errcode != berr.ErrNoError {
		return errcode
	}

	bc.updateImportLib(block)
	bc.handledBlockCallback(block)

	fmt.Printf("ImportBlock, number:%v, time:%v, delegate:%v, trxn:%v, hash:%x, prevHash:%x, version:%v\n",
		block.GetNumber(), common.TimeFormat(block.GetTimestamp()), string(block.GetDelegate()), len(block.BlockTransactions), block.Hash(), block.GetPrevBlockHash(), version.GetStringVersion(block.GetVersion()))

	return berr.ErrNoError
}
func (bc *BlockChain) popBlock() error {
	session := bc.dbInst.GetSession()
	if session != nil {
		bc.dbInst.ResetSession()
	}
	headBlockHash := bc.HeadBlockHash()

	b := bc.GetBlockByHash(headBlockHash)
	if b == nil {
		return fmt.Errorf("CHAIN block not found, block hash %x", headBlockHash)
	}

	bc.dbInst.Rollback()
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
	if block.Header.Number%uint64(config.BLOCKS_PER_ROUND) == 0 {
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
	lastBlockNum := uint64(0)
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

	bc.WriteBlock(block)
	bc.handledBlockCallback(block)
	return nil
}

func (bc *BlockChain) handledBlockCallback(block *types.Block) {
	for _, cb := range bc.handledBlockCB {
		if cb != nil {
			cb(block)
		}
	}
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
	start := common.MeasureStart()
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
	span := common.Elapsed(start)
	log.Infof("Insert block: block num:%v, trxn:%v, delegate: %v, hash:%x, span:%v\n\n", block.GetNumber(), len(block.Transactions), string(block.GetDelegate()), block.Hash(), span)

	return InsertBlockSuccess
}


func (b *BlockChain) GetLastBlockNumber() (uint64, error) {
	lastb := b.blockDb.GetLastBlock()
	if lastb == nil {
		return 0, fmt.Errorf("CHAIN last block not found")
	}

	lastBlockNum := lastb.GetNumber()
	return lastBlockNum, nil
}