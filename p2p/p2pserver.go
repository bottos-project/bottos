package p2pserver

import (
	"fmt"
	"net"
	"sync"
	"errors"
	"time"
	//"log"
	//"encoding/binary"
	//"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"github.com/AsynkronIT/protoactor-go/actor"

)

type netServer struct {

	port            uint16
	addr            string

	listener        net.Listener

	seed_peer       []string

	neighborList    []*net.UDPAddr
	serverAddr      *net.UDPAddr

	socket          *net.UDPConn

	nodeMap         map[uint64]*peer

	publicKey       string           //todo

	netLock         sync.RWMutex
}

//
type p2pServer struct{
	serv          *netServer
	notify        *notifyManager

	p2pConfig     *P2PConfig

	time_interval *time.Timer

	p2pLock        sync.RWMutex
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
	sync.RWMutex
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
		serv:          NewNetServer(p2pconfig),
		p2pConfig:     p2pconfig,
		time_interval: time.NewTimer(TIME_INTERVAL * time.Second),
	}
}

func (p2p *p2pServer) Init() error {

	fmt.Println("p2pServer::Init()")

	serv_addr := p2p.serv.addr + ":" + fmt.Sprint(p2p.serv.port)
	var err error

	p2p.serv.serverAddr, err = net.ResolveUDPAddr("udp4", serv_addr)
	if err != nil {
		fmt.Println("*ERROR* Failed to Resolve UDP Address !!!")
		return errors.New("*ERROR* Failed to Resolve UDP Address !!!")
	}

	//return (*UDPConn, error)
	p2p.serv.socket, err = net.ListenUDP("udp4", p2p.serv.serverAddr)
	if err != nil {
		fmt.Println("*ERROR* Failed to Listen !!! ", err)
		return errors.New("*ERROR* Failed to Listen !!! ")
	}


	return nil
}


//it is the entry of p2p
func (p2p *p2pServer) Start() error {
	fmt.Println("p2pServer::Start()")

	if p2p.p2pConfig == nil {
		return errors.New("*ERROR* P2P Configuration hadn't been inited yet !!!")
	}

	p2p.Init()

	if p2p.serv != nil {
		//wait for connection from others
		p2p.serv.Start()
	}

	//Todo connect to other seed nodes
	go p2p.ActiveSeeds()

	// Todo ping/pong
	go p2p.RunHeartBeat()

	return nil
}

func (p2p *p2pServer) ActiveSeeds() error {
	var i = 0
	for {
		select {
		case <- p2p.time_interval.C:
			fmt.Println("i = ",i)
			p2p.ResetTimer()
		}
	}
}

func  (p2p *p2pServer) ResetTimer ()  {
	p2p.time_interval.Stop()
	p2p.time_interval.Reset(time.Second * TIME_INTERVAL)
}

//connect seed during start p2p server
func (p2p *p2pServer) ConnectSeeds() error {
	fmt.Println("p2pServer::ConnectSeed()")

	for _ , peer := range p2p.p2pConfig.PeerLst {
		//fmt.Println(" peer = ",peer)
		p2p.Connect(peer)  //todo connect remote seed peer , if it's successful , add it into remote_list
	}

	return nil
}

//
func (p2p *p2pServer) Connect(addr string) error {
	//fmt.Println("p2pServer::ConnectSeed()")
	remoteAddr, err := net.ResolveUDPAddr("udp4", addr+":"+fmt.Sprint(p2p.serv.port))
	if err != nil {
		fmt.Println("*ERROR* Failed to create a remote server addr !!!")
		return errors.New("*ERROR* Failed to create a remote server addr !!!")
	}



	fmt.Println("remoteAddr = ",remoteAddr)

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

	//userlist should be packaged as a "peer" struct
	userlist  := make([]*net.UDPAddr, 0, 10)


	data := make([]byte, 4096)
	for {
		read, addr, err := serv.socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("*ERROR* Failed to receive the data : ", err)
			continue
		}

		fmt.Println("data = ",data[0:read]," , addr = ",addr)
		userlist = append(userlist, addr) //todo set a map[string]*Conn
		fmt.Println("userlist = ",userlist)
	}

	return
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












