// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  producer actor
 * @Author: eripi
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package p2p

import (
	"errors"
	log "github.com/cihub/seelog"
	"sync"
	"encoding/hex"
)

type collection struct {
	peers map[uint16]*Peer

	lock sync.RWMutex
}

func createCollection() *collection {
	c := &collection{
		peers: make(map[uint16]*Peer),
	}

	return c
}

func (c *collection) addPeer(peer *Peer) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	log.Debugf("P2P collection add peer index: %d, id: %s , add: %s, port: %s， chainId: %s, signature: %s, version: %d",
		peer.Index, peer.Info.Id, peer.Info.Addr, peer.Info.Port, peer.Info.ChainId, hex.EncodeToString(peer.Info.Signature), peer.Info.Version)

	if peer.Info.IsIncomplete() {
		log.Info("P2P peer info error")
		return errors.New("peer info error")
	}

	for _, p := range c.peers {
		if p.Info.Equal(peer.Info) {
			if p.isconn {
				log.Info("P2P peer is already exist")
				return errors.New("peer is already exist")
			}
		}
	}

	c.peers[peer.Index] = peer
	return nil
}

func (c *collection) getPeer(index uint16) *PeerInfo {
	c.lock.Lock()
	defer c.lock.Unlock()

	var info PeerInfo
	peer, ok := c.peers[index]
	if ok {
		info.ChainId = peer.Info.ChainId
		info.Addr = peer.Info.Addr
		info.Port = peer.Info.Port
		info.Signature = peer.Info.Signature
		info.Version = peer.Info.Version
		return &info
	}

	return nil
}

func (c *collection) delPeer(index uint16) bool {
	c.lock.Lock()
	defer c.lock.Unlock()


	peer, ok := c.peers[index]
	if ok {
		if peer.isconn {
			return false
		}
		peer.Stop()
		delete(c.peers, index)
		return true
	}


	return false
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
func (c *collection) getPeerP2PInfo() []Peer {
	c.lock.Lock()
	defer c.lock.Unlock()
	var peers []Peer
	for _, p := range c.peers {
		peers = append(peers, *p)
	}

	return peers
}
func (c *collection) send(msg *UniMsgPacket) {
	c.lock.Lock()
	defer c.lock.Unlock()

	peer, ok := c.peers[msg.Index]
	if !ok {
		log.Errorf("P2P peer not exist %s", msg.Index)
		return
	}

	if !peer.isconn {
		return
	}

	go peer.Send(msg.P)

}

func (c *collection) sendBroadcast(msg *BcastMsgPacket) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for id, peer := range c.peers {
		if !peer.isconn {
			continue
		}

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
