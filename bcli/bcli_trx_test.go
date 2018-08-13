package main

import (
	"fmt"
	"testing"
)


func Test_PushTransaction(t *testing.T) {
	cli := NewCLI()
	fmt.Println(cli)
	CONFIG = &CLIConfig{}
	CONFIG.KeyPairs = []KeyPair{{ PrivateKey: "b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45", PublicKey: "0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f"},}
	CONFIG.ChainId  = "00000000000000000000000000000000"
	var pushtrx BcliPushTrxInfo
	
	pushtrx.sender   = "tester01"
        pushtrx.contract = "nodeclustermng"
        pushtrx.method   = "reg"
        pushtrx.ParamMap = map[string]interface{}{"nodeIP":"0a0a0a0a", "clusterIP":"0b0b0b0b", "uuid":"33", "capacity":"2GB"}
	
	//fmt.Printf("BcliPushTransaction %v!!!", pushtrx)
	
	cli.BcliPushTransaction(&pushtrx)
}

