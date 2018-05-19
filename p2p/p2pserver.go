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


)

//
type p2pServer struct{
	serv          *netServer
	notify        *notifyManager

	p2pConfig     *P2PConfig

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

	if err := p2p.Init(); err != nil {
		return err
	}

	if p2p.serv != nil {
		//wait for connection from others
		p2p.serv.Start()
	}

	//Todo connect to other seed nodes
	go p2p.serv.ActiveSeeds()

	// Todo ping/pong
	go p2p.RunHeartBeat()

	return nil
}

//run a heart beat to watch the network status
func  (p2p *p2pServer) RunHeartBeat() error {
	fmt.Println("p2pServer::RunHeartBeat()")
	return nil
}













