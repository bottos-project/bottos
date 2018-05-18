package p2pserver

import (
	"fmt"
	"net"
	"sync"
	//"log"
	//"encoding/binary"
	//"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"github.com/AsynkronIT/protoactor-go/actor"

)

const (
	SERV_ADDR = "169.254.46.185:8080"
)

type netServer struct {

	port            uint16

	listener        net.Listener

	seed_peer       []string
	neighborNode    []*net.UDPAddr
	serverAddr      *net.UDPAddr
	nodeMap         map[uint64]*peer

	publicKey       string           //todo

	netLock         sync.RWMutex
}

//
type p2pServer struct{
	serv       *netServer
	notify     *notifyManager

	p2pLock     sync.RWMutex
}

type P2PInfo struct {
	ServAddr    string
	ServPort    string
	PeerLst     []string
}

type peer struct {
	peerName     string

	publicKey    string

	syncState    uint32
	neighborNode []peer
}

//its function to sync the trx , blk and peer info with other p2p other
type notifyManager struct {
	//
	p2p      *p2pServer

	peerList []peer

	stopSync chan bool
	pid      *actor.PID

	//for reading/writing peerlist
	syncLock sync.RWMutex
}

//
func ReadFile(filename string) P2PInfo{

	if filename == "" {
		fmt.Println("Error ! parmeter is null")
		return P2PInfo{}
	}
	var pi P2PInfo

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return  P2PInfo{}
	}

	str:=string(bytes)

	if err := json.Unmarshal([]byte(str), &pi) ; err != nil{
		fmt.Println("Unmarshal: ", err.Error())
		return P2PInfo{}
	}

	return pi
}

//
func NewServ() *p2pServer{
	fmt.Println("NewServ()")

	return nil
}

//it is the entry of p2p
func (p2p *p2pServer) Start() error {
	fmt.Println("p2p_server::Start()")

	if p2p.serv != nil {
		//wait for connection from others
		p2p.serv.Start()
	}

	//cconnect to other nodes
	go p2p.ConnectSeed()

	// ping/pong
	go p2p.RunHeartBeat()

	return nil
}

//
func (p2p *p2pServer) ConnectSeed() error {
	fmt.Println("p2p_server::ConnectSeed()")
	return nil
}


//run a heart beat to watch the network status
func  (p2p *p2pServer) RunHeartBeat() error {
	fmt.Println("p2p_server::RunHeartBeat()")
	return nil
}


//start listener
func (serv *netServer) Start() error {
	fmt.Println("net_server::Start()")

	go serv.Listening()

	return nil
}

//run accept
func (serv *netServer) Listening() {
	fmt.Println("net_server::Listening()")

}


func (notify *notifyManager) Start() {
	fmt.Println("block_sync_manager::Start")
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

//sync trx info with other peer
func (notify *notifyManager) BroadcastTrx() {
	fmt.Println("block_sync_manager::BroadcastTrx")
}

//sync blk info with other peer
func (notify *notifyManager) BroadcastBlk() {
	fmt.Println("block_sync_manager::BroadcastBlk")
}

//sync blk's hash info with other peer
func (notify *notifyManager) SyncHash() {
	fmt.Println("block_sync_manager::SyncHash")
}

//sync peer info with other peer
func (notify *notifyManager) SyncPeer() {
	fmt.Println("block_sync_manager::SyncPeer")
}












