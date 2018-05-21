package p2pserver

import  (
	"net"
)

type Peer struct {
	peerAddr     string
	publicKey    string

	peer_sock    *net.UDPAddr

	syncState    uint32
	neighborNode []*Peer
}

func NewPeer(addr_name string , addr_sock *net.UDPAddr) *Peer {
	return &Peer{
		peerAddr:   addr_name,
		peer_sock:  addr_sock,
		syncState:  0,
	}
}

func (p *Peer) GetPeerAddr() string {
	return p.peerAddr
}

func (p *Peer) SetPeerAddr(addr string) {
	p.peerAddr = addr
}



