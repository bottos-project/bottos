package p2pserver

type peer struct {
	peerName     string

	publicKey    string

	syncState    uint32
	neighborNode []peer
}


