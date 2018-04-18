

package transaction

import (
	"time"

	"github.com/bottos-project/core/common"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/action/message"
	
)


type TxPool struct {
	pending     map[common.Hash]*types.Transaction       
	expiration  map[common.Hash]*time.Time       
}


func CheckTransactionBaseConditionFromFront(){

	/* check max pending trx num */
	/* check account validate */
	/* check signature */

}


func CheckTransactionBaseConditionFromP2P(){	

}



// HandlTransactionFromFront handles a transaction from front
func HandleTransactionFromFront(trx *types.Transaction) {
	
    CheckTransactionBaseConditionFromFront()
	//start db session
	ApplyTransaction(trx)

	//revert db session

	//tell P2P actor to notify trx	
}



// HandlTransactionFromP2P handles a transaction from P2P
func HandleTransactionFromP2P(trx *types.Transaction) {

	CheckTransactionBaseConditionFromP2P()

	// start db session
	ApplyTransaction(trx)
	//revert db session	
}



func HandlePushTransactionReq(TrxSender message.TrxSenderType, trx *types.Transaction){

	if (message.TrxSenderTypeFront == TrxSender){ 
		HandleTransactionFromFront(trx)
	} else if (message.TrxSenderTypeP2P == TrxSender) {
		HandleTransactionFromP2P(trx)
	}	
}
