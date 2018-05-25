package transaction

import (
	"github.com/bottos-project/bottos/common/types"
)

type trxApplyApi interface {
	ApplyTransaction(trx *types.Transaction)
}

func NewTrxApplyService() *TrxApplyService {	
	return GetTrxApplyService()
}


