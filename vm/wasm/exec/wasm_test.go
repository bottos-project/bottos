package exec

import (
	"fmt"
	"testing"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract/msgpack"
)


//the case is to test crx recursive call
func TestWasmRecursiveCall (t *testing.T) {

	type transferparam struct {
		To			string
		Amount		uint32
	}

	param := transferparam {
		To     : "stewart",
		Amount : 1233,
	}

	bf , err :=  msgpack.Marshal(param)
	fmt.Println(" TestWasmRecursiveCall() bf = ",bf," , err = ",err)


	trx := &types.Transaction{
		Version:1,
		CursorNum:1,
		CursorLabel:1,
		Lifetime:1,
		Sender:"bottos",
		Contract: "usermng",
		Method:  "reguser",
		Param: bf,
		SigAlg:1,
		Signature:[]byte{},
	}

	ctx := &contract.Context{ Trx:trx}

	res , err := GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println("*ERROR* fail to execute start !!!")
		fmt.Println("err = ",err)
		return
	}

	fmt.Println("*SUCCESS* res = ",res, " , err = ",err)


	res , err = GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println("*ERROR* fail to execute start !!!")
		fmt.Println("err = ",err)
		return
	}

	fmt.Println("*SUCCESS* res = ",res, " , err = ",err)
}






