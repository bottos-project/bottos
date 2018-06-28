package discover

import "github.com/bottos-project/bottos/p2p"

const (
	PEER_INFO_REQ = 1
	PEER_INFO_RSP = 2

	PEER_HANDSHAKE_REQ     = 3
	PEER_HANDSHAKE_RSP     = 4
	PEER_HANDSHAKE_RSP_ACK = 5

	PEER_NEIGHBOR_REQ = 7
	PEER_NEIGHBOR_RSP = 8

	PEER_PING = 9
	PEER_PONG = 10
)

type PeerInfoReq struct {
}

type PeerInfoRsp struct {
	Info p2p.PeerInfo
}

type PeerNeighborReq struct {
}

type PeerNeighborRsp struct {
	Neighbor []p2p.PeerInfo
}
