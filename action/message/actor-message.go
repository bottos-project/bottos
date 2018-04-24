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
	TxHash common.Hash
}

type QueryTrxResp struct {
	Tx *types.Transaction
	Error error
}

type QueryBlockReq struct {
	BlockHash common.Hash
}

type QueryBlockResp struct {
	Block *types.Block
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

	