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


)

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

	fmt.Println("p2pServer::Start() p2p.serv.socket = ",p2p.serv.socket)

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
	fmt.Println("p2pServer::ActiveSeeds()")
	for {
		select {
		case <- p2p.time_interval.C:
			p2p.ConnectSeeds()
			p2p.ResetTimer()
		}
	}
}

//reset timer
func  (p2p *p2pServer) ResetTimer ()  {
	p2p.time_interval.Stop()
	p2p.time_interval.Reset(time.Second * TIME_INTERVAL)
}

//connect seed during start p2p server
func (p2p *p2pServer) ConnectSeeds() error {
	fmt.Println("p2pServer::ConnectSeed()")

	for _ , peer := range p2p.p2pConfig.PeerLst {
		//fmt.Println(" peer = ",peer)
		p2p.Connect(peer , false)  //todo connect remote seed peer , if it's successful , add it into remote_list
	}

	return nil
}

//
func (p2p *p2pServer) Connect(addr string , isExist bool) error {
	fmt.Println("p2pServer::ConnectSeed()")
	//todo check if the new peer is in peer list

	remoteAddr, err := net.ResolveUDPAddr("udp4", addr+":"+fmt.Sprint(p2p.serv.port))
	if err != nil {
		fmt.Println("*ERROR* Failed to create a remote server addr !!!")
		return errors.New("*ERROR* Failed to create a remote server addr !!!")
	}

	fmt.Println("remoteAddr = ",remoteAddr)
	//todo test connection with remote peer


	//todo package remote peer info as "peer" struct

	return nil
}

//run a heart beat to watch the network status
func  (p2p *p2pServer) RunHeartBeat() error {
	fmt.Println("p2pServer::RunHeartBeat()")
	return nil
}













