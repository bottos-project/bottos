package transaction

import (
	"github.com/bottos-project/core/common/types"
)

type trxApplyApi interface {
	ApplyTransaction(trx *types.Transaction)
}

func NewTrxApplyService() *TrxApplyService {	
	return GetTrxApplyService()
}


