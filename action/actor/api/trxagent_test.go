
package apiactor

import (
	"fmt"
	"time"
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/action/actor/transaction"
	"github.com/bottos-project/core/transaction"
	"github.com/bottos-project/core/action/message"
)

var trxActorPid *actor.PID



func TestPushTrxTest(t *testing.T) {

	trxActorPid = trxactor.NewTrxActor()

	InitTrxActorAgent()
	var trxPool = transaction.InitTrxPool()
	trxactor.SetTrxPool(trxPool)



	fmt.Println("Test PushTrxTest will called")
	
	trxTest := &types.Transaction{
		RefBlockNum: 11,
		Sender:      22,
		Action:      1,
	}
	
	reqMsg := &message.PushTrxReq{
		Trx: trxTest,
		TrxSender : message.TrxSenderTypeFront,
		
	}
	
	result, err := trxActorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result()

	if (nil == err) {
		fmt.Println("push trx req exec result:")
		fmt.Println("rusult is =======", result)
		fmt.Println("error  is =======", err)
	} else 	{ 
		t.Error("push trx failed, trx:", trxTest)
	}
	
}