package transaction

import (
	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
)


func ApplyTransaction(trx *types.Transaction) {

	/* save to db */

	/* call evm... */
}


func GetAllPendingTransactions() ([]*types.Transaction){

	return nil;
}


func RemoveTransactions(trxs []*types.Transaction){

}


func RemoveSingleTransaction(*types.Transaction){

}


func GetPendingTransaction(trxHash common.Hash) *types.Transaction {	

	return nil;
}