package p2p

import (
	"fmt"
	//"time"
	//"sync"

	"github.com/bottos-project/core/common"
)

type Protocol struct {
	bc *common.BlockChain
}

func NewProtocol(bc *common.BlockChain) *Protocol {
	proto := Protocol{bc}
	return &proto
}

func (p *Protocol) ProtocolLoop() {
	fmt.Println("P2P : start")

	for {

	}
}