package p2pserver

import (
	"testing"
	"fmt"
	"os"
	"github.com/bottos-project/core/config"
)

func TestP2PServ(t *testing.T)  {
	fmt.Println("p2p_server::Test1")
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Load config fail")
		os.Exit(1)
	}

	p2p := NewServ()
	p2p.Start()

	for{}

	return
}