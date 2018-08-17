package main

import (
	"testing"
	//"time"
	"io/ioutil"
)

func Test_PushTransaction(t *testing.T) {
	cli := NewCLI()
	CONFIG = &CLIConfig{}
	CONFIG.KeyPairs = []KeyPair{{ PrivateKey: "b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45", PublicKey: "0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f"},}
	CONFIG.ChainId  = "00000000000000000000000000000000"
	var pushtrx BcliPushTrxInfo
	
	pushtrx.sender   = "bottos"
        pushtrx.contract = "nodeclustermng"
        pushtrx.method   = "reg"
        pushtrx.ParamMap = map[string]interface{}{"nodeIP":"0a0a0a0a", "clusterIP":"0b0b0b0b", "uuid":"33", "capacity":"2GB"}

	cli.BcliPushTransaction("restful", "http://192.168.2.189:8689/v1/transaction/send", &pushtrx)

}

func Test_GetTransaction(t *testing.T) {
	cli := NewCLI()
	CONFIG = &CLIConfig{}
	CONFIG.KeyPairs = []KeyPair{{ PrivateKey: "b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45", PublicKey: "0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f"},}
	CONFIG.ChainId  = "00000000000000000000000000000000"
	CONFIG.ChainAddr = "192.168.2.189:8689"
	
	trxhash := "3711490b5cbd86d82f98d8d635ed80685460577f05eaf05a698d74b3f161d60b"
	
	cli.BcliGetTransaction("restful", "http://192.168.2.189:8689/v1/transaction/get", trxhash)
}

func Test_DeployCode(t *testing.T) {
	cli := NewCLI()
	CONFIG = &CLIConfig{}
	CONFIG.KeyPairs = []KeyPair{{ PrivateKey: "b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45", PublicKey: "0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f"},}
	CONFIG.ChainId  = "00000000000000000000000000000000"
	CONFIG.ChainAddr = "192.168.2.189:8689"
	cli.deploycode("restful", "http://192.168.2.189:8689/v1/transaction/send", "nodeclustermng", "./nodeclustermng.wasm")
}

func Test_DeployAbi(t *testing.T) {
	cli := NewCLI()
	CONFIG = &CLIConfig{}
	CONFIG.KeyPairs = []KeyPair{{ PrivateKey: "b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45", PublicKey: "0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f"},}
	CONFIG.ChainId  = "00000000000000000000000000000000"
	CONFIG.ChainAddr = "192.168.2.189:8689"
	cli.deployabi("grpc", "http://192.168.2.189:8689/v1/transaction/send", "nodeclustermng", "./nodeclustermng.abi")
}

func Test_DeployCodeAndAbi(t *testing.T) {
	cli := NewCLI()
	CONFIG = &CLIConfig{}
	CONFIG.KeyPairs = []KeyPair{{ PrivateKey: "b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45", PublicKey: "0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f"},}
	CONFIG.ChainId  = "00000000000000000000000000000000"
	CONFIG.ChainAddr = "192.168.2.189:8689"
	cli.deploycode("restful", "http://192.168.2.189:8689/v1/transaction/send", "nodeclustermng", "./nodeclustermng.wasm")
	time.Sleep(time.Duration(1) * time.Second)
	cli.deployabi("grpc", "http://192.168.2.189:8689/v1/transaction/send", "nodeclustermng", "./nodeclustermng.abi")
}

func Test_GetContractCode(t *testing.T) {
	cli := NewCLI()
	CONFIG = &CLIConfig{}
	CONFIG.KeyPairs = []KeyPair{{ PrivateKey: "b799ef616830cd7b8599ae7958fbee56d4c8168ffd5421a16025a398b8a4be45", PublicKey: "0454f1c2223d553aa6ee53ea1ccea8b7bf78b8ca99f3ff622a3bb3e62dedc712089033d6091d77296547bc071022ca2838c9e86dec29667cf740e5c9e654b6127f"},}
	CONFIG.ChainId  = "00000000000000000000000000000000"
	CONFIG.ChainAddr = "192.168.2.189:8689"
	contractcode, abi := cli.BcliGetContractCode("nodeclustermng", "./nodeclustermng_saved.wasm", "./nodeclustermng_saved.abi")
	writeFileToBinary(contractcode, "lypcontract.bin")
        ioutil.WriteFile("lypabi.abi", []byte(abi), 0644)
	
}
