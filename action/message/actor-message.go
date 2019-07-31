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
 * file description:  actor entry
 * @Author:
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */

package message

import (
	"github.com/bottos-project/bottos/common"
	berr "github.com/bottos-project/bottos/common/errors"
	"github.com/bottos-project/bottos/common/types"
)

//PushTrxReq trx request info
type PushTrxReq struct {
	Trx *types.Transaction
}
//PushTrxForP2PReq trx request info
type PushTrxForP2PReq struct {
	P2PTrx *types.P2PTransaction
}
//QueryTrxReq the key of trx query
type QueryTrxReq struct {
	TrxHash common.Hash
}

//QueryTrxResp the response of trx query
type QueryTrxResp struct {
	Trx   *types.Transaction
	Error error
}

//QueryBlockTrxResp the response of trx query
type QueryBlockTrxResp struct {
	Trx   *types.BlockTransaction
	Error error
}

//QueryBlockReq the key of block query
type QueryBlockReq struct {
	BlockHash   common.Hash
	BlockNumber uint64
}

//QueryBlockResp the response of block query
type QueryBlockResp struct {
	Block *types.Block
	Error error
}

//QueryChainInfoReq the key of chain info query
type QueryChainInfoReq struct {
}

//QueryChainInfoResp the response of chain info query
type QueryChainInfoResp struct {
	HeadBlockVersion      uint32
	HeadBlockNum          uint64
	LastConsensusBlockNum uint64
	HeadBlockHash         common.Hash
	HeadBlockTime         uint64
	HeadBlockDelegate     string
	Error                 error
}

//QueryAccountReq the key of account query
type QueryAccountReq struct {
	AccountName string
}

//QueryAccountResp the response of account query
type QueryAccountResp struct {
	AccountName   string
	Balance       uint64
	StakedBalance uint64
	Error         error
}

//InsertBlockReq  block info
type InsertBlockReq struct {
	Block *types.Block
}

//InsertBlockRsp the response of insert block
type InsertBlockRsp struct {
	Hash  common.Hash
	Error error
}

//GetAllPendingTrxReq the key of pending trx query
type GetAllPendingTrxReq struct {
}

//GetAllPendingTrxRsp the response of pending trx query
type GetAllPendingTrxRsp struct {
	Trxs []*types.Transaction
}

//RemovePendingTrxsReq the key of remove trx
type RemovePendingBlockTrxsReq struct {
	Trxs []*types.BlockTransaction
}

type RemovePendingTrxsReq struct {
	Trxs []*types.Transaction
}

//NotifyTrx txactor->p2pactor
type NotifyTrx struct {
	P2PTrx *types.P2PTransaction
}

//NotifyBlock producer->p2pactor
type NotifyBlock struct {
	Block *types.Block
}

//ReceiveTrx p2pactor->trxpool
type ReceiveTrx struct {
	P2PTrx *types.P2PTransaction
}

//ReceiveBlock p2pactor->chainactor
type ReceiveBlock struct {
	Block *types.Block
}

//ReceiveBlockResp chainactor->p2pactor
type ReceiveBlockResp struct {
	BlockNum uint64
	ErrorNo  berr.ErrCode
}

//ProducedBlockReq produceractor->consensusactor
type ProducedBlockReq struct {
	Block *types.Block
}

//RcvPrevoteReq p2pactor->consensusactor
type RcvPrevoteReq struct {
	BlockState *types.ConsensusBlockState
}

//RcvPrecommitReq p2pactor->consensusactor
type RcvPrecommitReq struct {
	BlockState *types.ConsensusBlockState
}

//RcvPrecommitReq p2pactor->consensusactor
type RcvCommitReq struct {
	BftHeaderState *types.ConsensusHeaderState
}

//SendPrevote consensusactor->p2pactor
type SendPrevote struct {
	BlockState *types.ConsensusBlockState
}

//SendPrecommit consensusactor->p2pactor
type SendPrecommit struct {
	BlockState *types.ConsensusBlockState
}

//SendCommit consensusactor->p2pactor
type SendCommit struct {
	BftHeaderState *types.ConsensusHeaderState
}

//PrevoteReq consensusactor->chainactor
type PrevoteReq struct {
	Block        *types.Block
	IsMyProduced bool
}

//PrecommitReq consensusactor->chainactor
type PrecommitReq struct {
	Block *types.Block
}

//PrecommitReq chainactor->p2pactor
type PrecommitResp struct {
	BlockNum uint64
	ErrorNo  berr.ErrCode
	MainFork bool
}

//CommitReq consensusactor->chainactor
type CommitReq struct {
	Block *types.Block
}

//SyncCommitReq consensusactor->chainactor
type SyncCommitReq struct {
	BftHeaderState *types.ConsensusHeaderState
}
