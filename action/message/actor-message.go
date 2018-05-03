package message


import (
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
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
	Trx *types.Transaction
	Error error
}

type QueryBlockReq struct {
	BlockHash common.Hash
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
	Error error
}

type QueryAccountReq struct {
	AccountName	string
}

type QueryAccountResp struct {
	AccountName   string
	Balance       uint64
	StakedBalance uint64
	Error error
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

	