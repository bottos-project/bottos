package exec

import (
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/contract/msgpack"
	log "github.com/cihub/seelog"
	"testing"
)

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
