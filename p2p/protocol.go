package p2p

import (
	"fmt"
	//"time"
	//"sync"

	"github.com/bottos-project/bottos/core/event"
	"github.com/bottos-project/bottos/core/common"
)

type Protocol struct {
	em *event.TypeMux
	bc *common.BlockChain
}

func NewProtocol(em *event.TypeMux, bc *common.BlockChain) *Protocol {
	proto := Protocol{em, bc}
	return &proto
}

func (p *Protocol) ProtocolLoop() {
	fmt.Println("P2P : start")

	for {

	}
}