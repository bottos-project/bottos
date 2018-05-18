package p2pserver

import (
	"fmt"
	"net"
	"sync"
	"errors"
	//"log"
	//"encoding/binary"
	//"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"github.com/AsynkronIT/protoactor-go/actor"

)

const (
	CONF_FILE = "config.json"
)

type netServer struct {

	port            uint16
	addr            string

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

	p2pConfig  *P2PConfig

	p2pLock     sync.RWMutex
}

type P2PConfig struct {
	ServAddr    string
	ServPort    uint16
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

//parse json configuration
func ReadFile(filename string) *P2PConfig{

	if filename == "" {
		fmt.Println("*ERROR* parmeter is null")
		return &P2PConfig{}
	}
	var pc P2PConfig

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("*ERROR* Failed to read the config: ",filename)
		return  &P2PConfig{}
	}

	str:=string(bytes)

	if err := json.Unmarshal([]byte(str), &pc) ; err != nil{
		fmt.Println("Unmarshal: ", err.Error())
		return &P2PConfig{}
	}

	//fmt.Println("ReadFile pc = ",pc)

	return &pc
}

//
func NewServ() *p2pServer{
	fmt.Println("NewServ()")

	p2pconfig := ReadFile(CONF_FILE)

	return &p2pServer{
		serv:       NewNetServer(p2pconfig),
		p2pConfig:  p2pconfig,
	}
}


//it is the entry of p2p
func (p2p *p2pServer) Start() error {
	fmt.Println("p2pServer::Start()")

	if p2p.p2pConfig == nil {
		return errors.New("*ERROR* P2P Configuration hadn't been inited yet !!!")
	}

	if p2p.serv != nil {
		//wait for connection from others
		p2p.serv.Start()
	}

	//cconnect to other seed nodes
	go p2p.ConnectSeed()

	// ping/pong
	go p2p.RunHeartBeat()

	return nil
}

//
func (p2p *p2pServer) ConnectSeed() error {
	fmt.Println("p2pServer::ConnectSeed()")
	return nil
}


//run a heart beat to watch the network status
func  (p2p *p2pServer) RunHeartBeat() error {
	fmt.Println("p2pServer::RunHeartBeat()")
	return nil
}


func NewNetServer(config *P2PConfig) *netServer {
	if config == nil {
		fmt.Println("*ERROR* Parmeter is empty !!!")
		return nil
	}

	return &netServer{
		addr: config.ServAddr,
		port: config.ServPort,
	}
}

//start listener
func (serv *netServer) Start() error {
	fmt.Println("netServer::Start()")

	go serv.Listening()

	return nil
}

//run accept
func (serv *netServer) Listening() {

	fmt.Println("netServer::Listening()")

	userlist := make([]*net.UDPAddr, 0, 10)
	serv_addr := serv.addr+":"+fmt.Sprint(serv.port)

	serverAddr, err := net.ResolveUDPAddr("udp4", serv_addr)
	if err != nil {
		fmt.Println("*ERROR* Failed to Resolve UDP Address !!!")
		return
	}

	socket, err := net.ListenUDP("udp4", serverAddr) //(*UDPConn, error)
	if err != nil {
		fmt.Println("*ERROR* Failed to Listen !!! ", err)
		return
	}

	data := make([]byte, 4096)
	for {
		read, addr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("*ERROR* Failed to receive the data : ", err)
			continue
		}

		fmt.Println("data = ",data[0:read]," , addr = ",addr)
		userlist = append(userlist, addr)
		fmt.Println("userlist = ",userlist)
	}
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












