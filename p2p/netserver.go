package p2pserver

import (
	"fmt"
	"net"
	"sync"
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
	peer_list  := make([]*net.UDPAddr, 0, 10)

	data := make([]byte, 4096)
	for {
		read, addr, err := serv.socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("*ERROR* Failed to receive the data : ", err)
			continue
		}

		fmt.Println("data = ",data[0:read]," , addr = ",addr)
		peer_list = append(peer_list, addr) //todo set a map[string]*Conn
		fmt.Println("peer_list = ",peer_list)
	}

	return
}
