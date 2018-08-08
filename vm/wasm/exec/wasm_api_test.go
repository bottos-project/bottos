package exec

import (
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/contract/msgpack"
	"testing"
	"fmt"
	"time"
	"reflect"
)

/*
func TestCallSubTrx(t *testing.T) {
	type transferparam struct {
		To     string
		Amount uint32
	}

	param := transferparam{
		To:     "Clinton",
		Amount: 1233,
	}

	bf, err := msgpack.Marshal(param)
	log.Infof(" TestCallSubTrx bf = ", bf, " , err = ", err)

	trx := &types.Transaction{
		Version:     1,
		CursorNum:   1,
		CursorLabel: 1,
		Lifetime:    1,
		Sender:      "bottos",
		Contract:    "usermng",
		Method:      "r",
		Param:       bf,
		SigAlg:      1,
		Signature:   []byte{},
	}

	ctx := &contract.Context{Trx: trx}

	res, err := GetInstance().Start(ctx, 1, false)
	if err != nil {
		log.Infof("*ERROR* fail to execute start !!! ", err.Error())
		return
	}

	//check sub trx
	var tf transferparam
	for _, sub_trx := range res {
		msgpack.Unmarshal(sub_trx.Param, &tf)
		log.Infof("TestCallSubTrx sub_trx = ", sub_trx.Param, " , tf = ", tf)
	}
}
*/

func TestSafeMath(t *testing.T) {
	type transferparam struct {
		To			string
		Amount		uint32
	}

	param := transferparam {
		To     : "stewart",
		Amount : 1233,
	}

	_ , err :=  msgpack.Marshal(param)

	var p string    = "dc0004da00087465737466726f6dda000b646174616465616c6d6e67da000344544fcf0000000000000064"
	var data []byte = []byte(p)

	fmt.Println("data = ",data)

	trx := &types.Transaction{
		Version:     1,
		CursorNum:   1,
		CursorLabel: 1,
		Lifetime:    1,
		Sender:      "bottos",
		Contract:    "usermng",
		Method:      "start",
		Param:       data,
		SigAlg:      1,
		Signature:   []byte{},
	}

	ctx := &contract.Context{Trx:trx}

	res , err := GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	var tf transferparam
	for _ , sub_trx := range res {
		//var tf transferparam
		msgpack.Unmarshal(sub_trx.Param , &tf)
	}
	fmt.Println("end of testcase")


	time.Sleep(time.Second * 3)

	vmi := GetInstance().GetWasteVM()
	if vmi == nil {
		return
	}

	fmt.Println("After GetInstance().Start(): ",reflect.TypeOf(vmi),",vmi: ",vmi)
}