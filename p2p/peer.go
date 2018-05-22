package p2pserver

import  (
	"net"
)

type Peer struct {
	peerAddr     string
	publicKey    string

	peer_sock    *net.UDPAddr
	conn         *net.Conn

	syncState    uint32
	neighborNode []*Peer
}

func NewPeer(addr_name string , conn *net.Conn) *Peer {
	return &Peer{
		peerAddr:   addr_name,
		conn:       conn,
		syncState:  0,
	}
}

func (p *Peer) GetPeerAddr() string {
	return p.peerAddr
}

func (p *Peer) SetPeerAddr(addr string) {
	p.peerAddr = addr
}



