package message


import (
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
	