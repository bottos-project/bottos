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
 * @Date:   2017-12-13
 * @Last Modified by:
 * @Last Modified time:
 */

package chain

import (
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/transaction"
	berr "github.com/bottos-project/bottos/common/errors"
	"gopkg.in/urfave/cli.v1"
)


//BlockCallback notify block
type BlockCallback func(*types.Block)

//BlockChainInterface the interface of chain
type BlockChainInterface interface {
	Init(ctx *cli.Context) error
	SetTrxPool(trxPool *transaction.TrxPool)
	InitOnRecover(ctx *cli.Context) error
	Close()

	IsMainForkHasBlock(block *types.Block) bool
	GetBlockByHash(hash common.Hash) *types.Block
	GetBlockByNumber(number uint64) *types.Block
	GetHeaderByNumber(number uint64) *types.Header
	GetTransaction(hash common.Hash) *types.BlockTransaction
	GetCommittedTransaction(hash common.Hash) *types.BlockTransaction
	GetLastBlockNumber() (uint64, error)

	HeadBlockTime() uint64
	HeadBlockNum() uint64
	HeadBlockHash() common.Hash
	HeadBlockDelegate() string
	LastConsensusBlockNum() uint64

	InsertBlock(block *types.Block) berr.ErrCode
	ImportBlock(block *types.Block) berr.ErrCode

	PreVoteBlock(block *types.Block) berr.ErrCode
	PreCommitBlock(block *types.Block) berr.ErrCode
	CommitBlock(block *types.Block) berr.ErrCode
	SyncCommitBlock(BftHeaderState *types.ConsensusHeaderState) berr.ErrCode

	RegisterHandledBlockCallback(cb BlockCallback)
	RegisterCommittedBlockCallback(cb BlockCallback)
}
