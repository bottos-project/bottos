package p2p

import (
	"errors"
	log "github.com/cihub/seelog"
	"sync"
)

type collection struct {
	peers map[uint16]*Peer

	lock sync.RWMutex
}

const (
	MAX_PEER_NUM = 100
)

func createCollection() *collection {
	c := &collection{
		peers: make(map[uint16]*Peer),
	}

	return c
}

func (c *collection) addPeer(peer *Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debugf("collection add peer index: %d, id: %s , add: %s, port: %s",
		peer.Index, peer.Info.Id, peer.Info.Addr, peer.Info.Port)

	if peer.Info.IsIncomplete() {
		log.Info("peer info error")
		return errors.New("peer info error")
	}

	for _, p := range c.peers {
		if p.Info.Equal(peer.Info) {
			if p.isconn {
				log.Info("peer is already exist")
				return errors.New("peer is already exist")
			}
		}
	}

	c.peers[peer.Index] = peer
	return nil
}

func (c *collection) delPeer(index uint16) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debugf("collection delete peer index: %d", index)
	peer, ok := c.peers[index]
	if ok {
		log.Debug("delete peer")
		if peer.isconn {
			log.Error("peer is connected , don't delete")
			return false
		}

		peer.Stop()
		delete(c.peers, index)
		return true
	} else {
		log.Error("pee not exist")
		return false
	}
}

func (c *collection) isPeerExist(index uint16) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.peers[index]
	return ok
}

func (c *collection) isPeerInfoExist(info PeerInfo) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, value := range c.peers {
		if value.Info.Equal(info) {
			return true
		}
	}

	return false
}

func (c *collection) getPeers() []PeerInfo {
	c.lock.Lock()
	defer c.lock.Unlock()

	var peers []PeerInfo
	for _, p := range c.peers {
		peers = append(peers, p.Info)
	}

	return peers
}

func (c *collection) getPeersData() PeerDataSet {
	c.lock.Lock()
	defer c.lock.Unlock()

	var peers PeerDataSet
	for _, p := range c.peers {
		peers = append(peers, PeerData{Id: p.Info.Id, Index: p.Index})
	}

	return peers
}

func (c *collection) send(msg *UniMsgPacket) {
	c.lock.Lock()
	defer c.lock.Unlock()

	peer, ok := c.peers[msg.Index]
	if !ok {
		log.Errorf("peer not exist %s", msg.Index)
		return
	}

	go peer.Send(msg.P)

}

func (c *collection) sendBroadcast(msg *BcastMsgPacket) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for id, peer := range c.peers {
		if len(msg.Indexs) == 0 {
			go peer.Send(msg.P)
			continue
		}

		//filter index by msg index set
		bsend := true
		for _, eid := range msg.Indexs {
			if id == eid {
				bsend = false
				break
			}
		}

		if bsend {
			go peer.Send(msg.P)
		}

	}
}
