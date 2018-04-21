package transaction

import (
	"github.com/bottos-project/core/common/types"
)


func CheckTransactionLifeTime(trx *types.Transaction) { 

}

func CheckTransactionUnique(trx *types.Transaction) { 

}


func CheckTransactionMatchChain(trx *types.Transaction) { 

}

func ApplyTransaction(trx *types.Transaction) {

	CheckTransactionLifeTime(trx)
	CheckTransactionUnique(trx)
	CheckTransactionMatchChain(trx)
	/* save to db */

	/* call evm... */
}


