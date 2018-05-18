package p2pserver

import (
	"testing"
	"fmt"
)

func TestP2PServ(t *testing.T)  {
	fmt.Println("p2p_server::Test1")

	p2p := NewServ()
	p2p.Start()

	for{}

	return
}