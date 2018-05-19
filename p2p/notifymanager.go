package p2pserver

import (
	"fmt"
	"sync"
	"github.com/AsynkronIT/protoactor-go/actor"
)

//its function to sync the trx , blk and peer info with other p2p other
type notifyManager struct {
	//
	p2p      *p2pServer

	peerList []Peer

	stopSync chan bool
	pid      *actor.PID

	//for reading/writing peerlist
	sync.RWMutex
}

func (notify *notifyManager) Start() {
	fmt.Println("notifyManager::Start")
	for {
		//signal from actor
		go notify.BroadcastTrx()
		//signal from actor
		go notify.BroadcastBlk()

		go notify.SyncHash()
		go notify.SyncPeer()

		//receive
	}
}

func (notify *notifyManager) BroadCast(buf []byte, isSync bool) {
	notify.RLock()
	defer notify.RUnlock()

	for _ , node := range notify.peerList {
		fmt.Println("node: ",node)
	}

	return
}

//sync trx info with other peer
func (notify *notifyManager) BroadcastTrx() {
	fmt.Println("notifyManager::BroadcastTrx")
}

//sync blk info with other peer
func (notify *notifyManager) BroadcastBlk() {
	fmt.Println("notifyManager::BroadcastBlk")
}

//sync blk's hash info with other peer
func (notify *notifyManager) SyncHash() {
	fmt.Println("notifyManager::SyncHash")
}

//sync peer info with other peer
func (notify *notifyManager) SyncPeer() {
	fmt.Println("notifyManager::SyncPeer")
}

























