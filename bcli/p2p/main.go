package main

import p2pserv "github.com/bottos-project/bottos/p2p"

func main() {

	p2p := p2pserv.NewServ()

	go p2p.Start()

	select {}
}
