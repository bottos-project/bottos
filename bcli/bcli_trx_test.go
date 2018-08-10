package main

import (
	"fmt"
	"testing"

	//"github.com/bottos-project/bottos/contract/msgpack"
	//log "github.com/cihub/seelog"
)


func Test_PushTransaction(t *testing.T) {
	
	var pushtrx BcliPushTrxInfo
	cli := NewCLI()
	
	pushtrx.sender   = "bottos"
        pushtrx.contract = "nodeclustermng"
        pushtrx.method   = "reg"
        pushtrx.ParamMap = map[string]interface{}{"nodeIP":"0a0a0a0a", "clusterIP":"0b0b0b0b", "uuid":"33", "capacity":"2GB"}
	
	fmt.Printf("LYP: BcliPushTransaction %v!!!", pushtrx)
	
	cli.BcliPushTransaction(&pushtrx)
}

