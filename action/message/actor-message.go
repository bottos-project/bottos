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

type TrxSenderType uint8

const (
	InvalidSenderType TrxSenderType = iota
	TrxSenderTypeFront
	TrxSenderTypeP2P

	MaxTrxSenderType
)

type PushTrxReq struct {
	Trx *types.Transaction

	TrxSender TrxSenderType
}

type QueryTrxReq struct {
	TrxHash common.Hash
}

type QueryTrxResp struct {
	Trx   *types.Transaction
	Error error
}

type QueryBlockReq struct {
	BlockHash   common.Hash
	BlockNumber uint32
}

type QueryBlockResp struct {
	Block *types.Block
	Error error
}

type QueryChainInfoReq struct {
}

type QueryChainInfoResp struct {
	HeadBlockNum          uint32
	LastConsensusBlockNum uint32
	HeadBlockHash         common.Hash
	HeadBlockTime         uint64
	HeadBlockDelegate     string
	Error                 error
}

type QueryAccountReq struct {
	AccountName string
}

type QueryAccountResp struct {
	AccountName   string
	Balance       uint64
	StakedBalance uint64
	Error         error
}

type InsertBlockReq struct {
	Block *types.Block
}

type InsertBlockRsp struct {
	Hash  common.Hash
	Error error
}

type GetAllPendingTrxReq struct {
}

type GetAllPendingTrxRsp struct {
	Trxs []*types.Transaction
}

type RemovePendingTrxsReq struct {
	Trxs []*types.Transaction
}

type RemovePendingTrxsRsp struct {
}

// txactor->p2pactor
type NotifyTrx struct {
	Trx *types.Transaction
}

// producer->p2pactor
type NotifyBlock struct {
	Block *types.Block
}

// p2pactor->trxpool
type ReceiveTrx struct {
	Trx *types.Transaction
}

// p2pactor->chainactor
type ReceiveBlock struct {
	Block *types.Block
}
