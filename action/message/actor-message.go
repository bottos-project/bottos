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
	"github.com/bottos-project/bottos/common/types"
)

//PushTrxReq trx request info
type PushTrxReq struct {
	Trx *types.Transaction
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

//QueryBlockReq the key of block query
type QueryBlockReq struct {
	BlockHash   common.Hash
	BlockNumber uint32
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
	HeadBlockNum          uint32
	LastConsensusBlockNum uint32
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
type RemovePendingTrxsReq struct {
	Trxs []*types.Transaction
}

//NotifyTrx txactor->p2pactor
type NotifyTrx struct {
	Trx *types.Transaction
}

//NotifyBlock producer->p2pactor
type NotifyBlock struct {
	Block *types.Block
}

//ReceiveTrx p2pactor->trxpool
type ReceiveTrx struct {
	Trx *types.Transaction
}

//ReceiveBlock p2pactor->chainactor
type ReceiveBlock struct {
	Block *types.Block
}
