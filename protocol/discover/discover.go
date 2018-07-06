package discover

import (
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/p2p"
	log "github.com/cihub/seelog"
	"net"
	"time"
)

type Discover struct {
	c      *candidates
	p      *pne
	k      *keeplive
	sendup p2p.SendupCb
}

const (
	//TIME_DISCOVER connect to unknow peer, second
	TIME_DISCOVER = 5
	//NEIGHBOR_DISCOVER_COUNT
	NEIGHBOR_DISCOVER_COUNT = 10
)

func MakeDiscover(config *config.Parameter) *Discover {
	d := &Discover{}
	d.p = makePne(config)
	d.c = makeCandidates(d.p)
	d.k = makeKeeplive(d.c, d.p)

	d.c.setKeeplive(d.k)

	return d
}

func (d *Discover) Start() {
	d.c.start()
	d.p.start()
	d.k.start()

	go d.discoverTimer()
}

func (d *Discover) SetSendupCallback(cb p2p.SendupCb) {
	d.sendup = cb
}

func (d *Discover) Dispatch(index uint16, p *p2p.Packet) {
	switch p.H.PacketType {
	case PEER_INFO_REQ:
		d.c.processPeerInfoReq(index, p.Data)
	case PEER_INFO_RSP:
		d.c.processPeerInfoRsp(index, p.Data)
	case PEER_HANDSHAKE_REQ:
		d.c.processHandshakeReq(index, p.Data)
	case PEER_HANDSHAKE_RSP:
		d.c.processHandshakeRsp(index, p.Data)
	case PEER_HANDSHAKE_RSP_ACK:
		d.c.processHandshakeRspAck(index, p.Data)
	case PEER_NEIGHBOR_REQ:
		d.p.processPneNeighborReq(index, p.Data)
	case PEER_NEIGHBOR_RSP:
		d.p.processPneNeighborRsp(index, p.Data)
	case PEER_PING:
		d.k.processPing(index, p.Data)
	case PEER_PONG:
		d.k.processPong(index, p.Data)
	default:
		log.Errorf("discover Dispatch packet type:%d error", p.H.PacketType)
	}

}

func (d *Discover) NewConnCb(conn net.Conn, sendup p2p.SendupCb) {
	//new candidate peer
	info := p2p.PeerInfo{}
	p := p2p.CreatePeer(info, conn, true, sendup)

	error := d.c.addCandidate(p)
	if error != nil {
		p.Stop()
		return
	}

	p.Start()
}

func (d *Discover) newConn(peer p2p.PeerInfo) error {
	addrPort := peer.Addr + ":" + peer.Port
	conn, err := net.DialTimeout("tcp", addrPort, 2*time.Second)
	if err != nil {
		log.Debugf("failed to connect to peer：%s:%s", peer.Addr, peer.Port)
		return err
	}

	p := p2p.CreatePeer(peer, conn, false, d.sendup)

	err = d.c.addCandidate(p)
	if err != nil {
		p.Stop()
		return err
	}

	p.Start()
	return nil
}

func (d *Discover) discoverTimer() {
	log.Debug("discoverTimer")

	dicover := time.NewTimer(TIME_DISCOVER * time.Second)

	defer func() {
		log.Debug("discoverTimer stop")
		dicover.Stop()
	}()

	for {
		select {
		case <-dicover.C:
			if d.c.isCandidateFull() {
				dicover.Reset(TIME_DISCOVER * time.Second)
				continue
			}

			neighbors := d.p.n.nextPneNeighbors()
			if neighbors == nil {
				dicover.Reset(TIME_DISCOVER * time.Second)
				continue
			}

			/*try to connect peer*/
			for _, peer := range neighbors {
				err := d.newConn(peer)
				if err != nil {
					continue
				}

				d.p.n.delNeighbor(peer)
			}

			dicover.Reset(TIME_DISCOVER * time.Second)
		}
	}
}
