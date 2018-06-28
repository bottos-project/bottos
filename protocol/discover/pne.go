package discover

import (
	"container/list"
	"encoding/json"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/p2p"
	pcommon "github.com/bottos-project/bottos/protocol/common"
	log "github.com/cihub/seelog"
	"strings"
	"sync"
	"time"
)

type pne struct {
	qPeers *common.Queue
	lock   sync.RWMutex

	n *neighbors

	seeds []p2p.PeerInfo
}

const (
	//TIME_PNE_EXCHANGE time to exchange peer neighbor info, minute
	TIME_PNE_EXCHANGE = 1
)

func makePne(config *config.Parameter) *pne {

	pne := &pne{
		qPeers: common.NewQueue(),
		n:      makeNeighbors(),
	}

	pne.parseSeeds(config)

	return pne
}

func (p *pne) start() {
	go p.pneTimer()
}

func (p *pne) pushPeerIndex(index uint16) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.qPeers.Push(index)
}

func (p *pne) parseSeeds(config *config.Parameter) {

	var peers []p2p.PeerInfo

	for _, element := range config.PeerList {
		peerCfg := strings.Split(element, ":")

		addr := peerCfg[0]
		port := peerCfg[1]

		var peer p2p.PeerInfo
		peer.Addr = addr
		peer.Port = port
		peers = append(peers, peer)
		p.n.addNeighbor(peers)
	}
}

func (p *pne) pneTimer() {
	log.Debug("pneTimer")

	exchange := time.NewTimer(TIME_PNE_EXCHANGE * time.Minute)

	defer func() {
		log.Debug("pneTimer stop")
		exchange.Stop()
	}()

	for {
		select {
		case <-exchange.C:
			index := p.nextPeer()
			if index != 0 {
				log.Debugf("pneTimer peer index: %d", index)
				p.sendPneRequest(index)
			}

			exchange.Reset(TIME_PNE_EXCHANGE * time.Minute)
		}
	}
}

//NextPnePeer pop a peer id
func (p *pne) nextPeer() uint16 {
	p.lock.Lock()
	defer p.lock.Unlock()

	value := p.qPeers.Pop()
	if value == nil {
		return 0
	}

	return value.(uint16)
}

func (p *pne) sendPneRequest(index uint16) {
	//check peer is exist or not
	ok := p2p.Runner.IsPeerExist(index)
	if !ok {
		return
	}

	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_NEIGHBOR_REQ,
	}

	packet := p2p.Packet{H: head}

	send := p2p.MsgPacket{
		Index: []uint16{index},
		P:     packet,
	}

	p2p.Runner.SendUnicast(send)

	// add back to queue
	p.pushPeerIndex(index)
}

func (p *pne) sendPneResponse(index uint16) {
	peers := p2p.Runner.GetPeers()
	if len(peers) == 0 {
		return
	}

	resp := PeerNeighborRsp{
		Neighbor: peers,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("addrs Marshal error:%s", err)
		return
	}

	head := p2p.Head{ProtocolType: pcommon.P2P_PACKET,
		PacketType: PEER_NEIGHBOR_RSP,
	}

	packet := p2p.Packet{H: head,
		Data: data,
	}

	send := p2p.MsgPacket{
		Index: []uint16{index},
		P:     packet,
	}

	p2p.Runner.SendUnicast(send)
}

func (p *pne) processPneNeighborReq(index uint16, date []byte) {
	//check peer is exist or not
	ok := p2p.Runner.IsPeerExist(index)
	if !ok {
		return
	}

	p.sendPneResponse(index)
}

func (p *pne) processPneNeighborRsp(index uint16, date []byte) {
	//check peer is exist or not
	ok := p2p.Runner.IsPeerExist(index)
	if !ok {
		return
	}

	var rsp PeerNeighborRsp
	err := json.Unmarshal(date, &rsp)
	if err != nil {
		log.Errorf("ProcessPneNeighborRsp Unmarshal error")
		return
	}

	//add neighbor
	p.n.addNeighbor(rsp.Neighbor)
}

type neighbors struct {
	neighbor *list.List
	lock     sync.RWMutex
	nIndex   uint16
}

func makeNeighbors() *neighbors {
	return &neighbors{neighbor: list.New()}
}

//AddNeighbor add new neighbor
func (n *neighbors) addNeighbor(peers []p2p.PeerInfo) {
	n.lock.Lock()
	n.lock.Unlock()

	//filter neighbor which is exist
	for j := range peers {
		//filter neighbor of ourself
		if p2p.LocalPeerInfo.Equal(peers[j]) {
			continue
		}

		//filter neighbor of peer
		ok := p2p.Runner.IsPeerInfoExist(peers[j])
		if ok {
			return
		}

		//filter neighbor which is exist
		find := false
		for e := n.neighbor.Front(); e != nil; e = e.Next() {
			peerinfo := e.Value.(p2p.PeerInfo)
			if peerinfo.Equal(peers[j]) {
				find = true
				break
			}
		}

		if !find {
			n.neighbor.PushBack(peers[j])
		}
	}
}

//DelNeighbor delete neighbor
func (n *neighbors) delNeighbor(peer p2p.PeerInfo) {
	n.lock.Lock()
	n.lock.Unlock()

	//find addr and remove
	var count uint16
	for e := n.neighbor.Front(); e != nil; e = e.Next() {
		peerinfo := e.Value.(p2p.PeerInfo)
		if peerinfo.Equal(peer) {
			n.neighbor.Remove(e)

			if n.nIndex > count {
				n.nIndex--
			}
			return
		}

		count++
	}
}

//NextPneNeighbors get neighbors to discover
func (n *neighbors) nextPneNeighbors() []p2p.PeerInfo {
	n.lock.Lock()
	n.lock.Unlock()

	len := n.neighbor.Len()
	if len <= 0 {
		return nil
	}

	var start uint16
	var end uint16
	if n.nIndex+NEIGHBOR_DISCOVER_COUNT > uint16(len) {
		start = n.nIndex
		end = uint16(len)

		n.nIndex = 0
	} else {
		start = n.nIndex
		end = start + NEIGHBOR_DISCOVER_COUNT

		n.nIndex = end + 1
	}

	var count uint16
	var peers []p2p.PeerInfo
	for e := n.neighbor.Front(); e != nil; e = e.Next() {
		if count >= start && count <= end {
			peers = append(peers, e.Value.(p2p.PeerInfo))
		}
		count++
	}

	return peers
}
