package p2pserver

type Peer struct {
	peerAddr     string

	publicKey    string

	syncState    uint32
	neighborNode []*Peer
}

func NewPeer(addr string) *Peer {
	return &Peer{
		peerAddr:   addr,
		syncState:  0,
	}
}

func (p *Peer) GetPeerAddr() string {
	return p.peerAddr
}

func (p *Peer) SetPeerAddr(addr string) {
	p.peerAddr = addr
}



