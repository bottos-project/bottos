package p2pserver

import (
	"io"
	"fmt"
	"net"
	"sync"
	"time"
	"errors"
	"strings"
	"crypto/sha1"
)

type netServer struct {
	config          *P2PConfig
	port            uint16
	addr            string

	listener        net.Listener

	seed_peer       []string

	neighborList    []*net.UDPAddr
	serverAddr      *net.UDPAddr
	socket          *net.UDPConn

	peerMap         map[uint64]*Peer

	publicKey       string           //todo

	time_interval   *time.Timer

	netLock         sync.RWMutex
}

func NewNetServer(config *P2PConfig) *netServer {
	if config == nil {
		fmt.Println("*ERROR* Parmeter is empty !!!")
		return nil
	}

	return &netServer{
		config:        config,
		addr:          config.ServAddr,
		port:          config.ServPort,
		time_interval: time.NewTimer(TIME_INTERVAL * time.Second),
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
	peer_list  := make([]*net.UDPAddr, 0, 10)

	data := make([]byte, 4096)
	for {
		read, addr, err := serv.socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("*ERROR* Failed to receive the data : ", err)
			continue
		}

		//In here . it need use different handler function according to requirement, eg. handle login/income blk.ctx and so on

		fmt.Println("data = ",data[0:read]," , addr = ",addr)
		peer_list = append(peer_list, addr) //todo set a map[string]*Conn
		fmt.Println("peer_list = ",peer_list)
	}

	return
}

func (serv *netServer) ActiveSeeds() error {
	fmt.Println("p2pServer::ActiveSeeds()")
	for {
		select {
		case <- serv.time_interval.C:
			serv.ConnectSeeds()
			serv.ResetTimer()
		}
	}
}

//reset timer
func  (serv *netServer) ResetTimer ()  {
	serv.time_interval.Stop()
	serv.time_interval.Reset(time.Second * TIME_INTERVAL)
}

//connect seed during start p2p server
func (serv *netServer) ConnectSeeds() error {
	fmt.Println("p2pServer::ConnectSeed()")

	for _ , peer := range serv.config.PeerLst {
		fmt.Println("try to connect peer: ",peer)
		serv.Connect(peer , false)  //todo connect remote seed peer , if it's successful , add it into remote_list
	}

	return nil
}

//to connect certain peer
func (serv *netServer) Connect(addr string , isExist bool) error {
	fmt.Println("p2pServer::ConnectSeed()")

	//check if the new peer is in peer list
	if serv.IsExist(addr , isExist) {
		return nil
	}

	addr_port := addr+":"+fmt.Sprint(serv.port)
	remoteAddr, err := net.ResolveUDPAddr("udp4", addr_port)
	if err != nil {
		fmt.Println("*ERROR* Failed to create a remote server addr !!!")
		return errors.New("*ERROR* Failed to create a remote server addr !!!")
	}

	//test connection with remote peer
	var TST_DAT []byte
	_ , err = serv.socket.WriteToUDP(TST_DAT , remoteAddr)
	if err != nil { //todo check len
		fmt.Println("*ERROR* Failed to send Test message to remote peer !!!")
		return errors.New("*ERROR* Failed to send Test message to remote peer !!!")
	}

	//todo package remote peer info as "peer" struct and add it into peer list
	peer := NewPeer(addr)

	sha_handler := sha1.New()
	io.WriteString(sha_handler , addr_port)

	sha_handler.Sum(nil)
	serv.peerMap[1] = peer

	return nil
}

func (serv *netServer) IsExist(addr string , isExist bool) bool {

	for _ , peer := range serv.peerMap {
		if res := strings.Compare(peer.peerAddr , addr); res == 0 {
			return true
		}
	}

	return false
}

