package main

import (
	"encoding/json"
	"github.com/bottos-project/bottos/bcli/p2p/stub"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/protocol"
	log "github.com/cihub/seelog"
	"io/ioutil"
)

type p2pConfig struct {
	ServAddr string
	ServPort string
	PeerLst  []string
}

type chainConfig struct {
	BlockNumber uint32
	Blocks      []types.Block
}

func main() {

	pc := readP2PConfig("p2pconfig.json")
	if pc == nil {
		return
	}

	param := config.Parameter{ServAddr: pc.ServAddr,
		P2PPort:  pc.ServPort,
		PeerList: pc.PeerLst,
	}

	bc := stub.MakeBlockChainStub()

	chain := readChainConfig("chainconfig.json")
	if chain != nil {
		bc.SetBlockNumber(chain.BlockNumber)
		bc.SetBlocks(chain.Blocks)
	}

	log.Info("blocknumber:", chain.BlockNumber)
	p := protocol.MakeProtocol(&param, bc)

	actor := stub.NewDumActor()

	p.SetChainActor(actor)

	go p.Start()

	select {}
}

//ReadFile parse json configuration
func readP2PConfig(filename string) *p2pConfig {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	str := string(bytes)

	var pc p2pConfig
	if err := json.Unmarshal([]byte(str), &pc); err != nil {
		return nil
	}

	return &pc
}

//ReadFile parse json configuration
func readChainConfig(filename string) *chainConfig {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	str := string(bytes)

	var pc chainConfig
	if err := json.Unmarshal([]byte(str), &pc); err != nil {
		return nil
	}

	return &pc
}
