
package apiactor

import (
	"fmt"
	"time"
	//"context"
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/bottos-project/core/common/types"
	//"github.com/bottos-project/core/api"
	//caapi "github.com/bottos-project/core/action/actor/api"
	"github.com/bottos-project/core/action/actor/transaction"
	"github.com/bottos-project/core/transaction"
	"github.com/bottos-project/core/action/message"
)

var trxActorPid *actor.PID



func TestPushTrxTest(t *testing.T) {

	trxActorPid = trxactor.NewTrxActor()
	
	fmt.Println("*****trxactorPid is " ,trxactorPid)


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

	
	result, err := trxactorPid.RequestFuture(reqMsg, 500*time.Millisecond).Result() // await result

	//_, err := trxactorPid.RequestFuture(req, 500*time.Millisecond).Result() // await result

	if (nil != err) {
		
	}
	
	fmt.Println("rusult is =======", result, err)

	fmt.Println("exec push transaction done !!!!!!,trx:", trxTest)
}


