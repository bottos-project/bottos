package p2pserver

import (
	//"io"
	"fmt"
	"net"
	"sync"
	"time"
	"errors"
	"strings"
	//"reflect"
	"unsafe"
	//"crypto/sha1"
	"hash/fnv"
	"encoding/json"
	//"github.com/bottos-project/core/contract/msgpack"
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

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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
		peerMap:       make(map[uint64]*Peer),
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
	//listener, err := net.Listen("tcp", serv.addr+":"+fmt.Sprint(serv.port))
	listener, err := net.Listen("tcp", ":"+fmt.Sprint(serv.port))
	if err != nil {
		fmt.Println("*ERROR* Failed to listen at port: "+fmt.Sprint(serv.port))
		return
	}

	defer listener.Close()

	data := make([]byte, 4096)
	var msg message

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("netServer::Listening() Failed to accept")
			continue
		}

		len , err := conn.Read(data)
		if err != nil {
			fmt.Println("*WRAN* Can't read data from remote peer !!!")
			continue
		}

		err = json.Unmarshal(data[0:len] , &msg)
		if err != nil {
			fmt.Println("*WRAN* Can't unmarshal data from remote peer !!!")
			continue
		}

		go serv.HandleMsg(conn , &msg)
	}

}

//run accept
func (serv *netServer) HandleMsg(conn net.Conn , msg *message) {
	//defer conn.Close()
	fmt.Println("netServer::HandleMsg() ")

	switch msg.MsgType {
	case request:

		rsp := message {
			Src:      serv.addr,
			Dst:      msg.Src,
			MsgType:  response,
		}

		data , err := json.Marshal(rsp)
		if err != nil{
			fmt.Println("*WRAN* Failed to package the response message : ", err)
		}

		//fmt.Println("netServer::HandleMsg() request len(data) = ",len(data) , ", rsp = ",rsp," , conn.RemoteAddr = ",conn.RemoteAddr()," , conn.LocalAddr = ",conn.LocalAddr())

		//create a new conn to response the remote peer
		newconn , err := net.Dial("tcp", msg.Src+":"+fmt.Sprint(serv.port))
		if err != nil {
			fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ",err)
			return
		}

		len , err := newconn.Write(data)
		if err != nil {
			fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
			return
		}

		if len < 0 {
			fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
			return
		}

		//package remote peer info as "peer" struct and add it into peer list
		addr_port := msg.Dst + ":" + fmt.Sprint(serv.port)
		peer := NewPeer(msg.Dst , &conn)
		peer_identify := Hash(addr_port)
		serv.peerMap[uint64(peer_identify)] = peer
	case response:

		fmt.Println("response - netServer::Listening() serv.peerMap = ",serv.peerMap)
		/*
		for key , peer := range serv.peerMap {
			fmt.Println("msg.Dst = ",msg.Dst," , key = ",key," , peer = ",peer)
		}
		*/
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
		//check if the new peer is in peer list
		if serv.IsExist(peer , false) {
			fmt.Println("the Peer had existed !")
			continue
		}

		var msg = message {
			Src:      serv.addr,
			Dst:      peer,
			MsgType:  request,
		}

		req , err := json.Marshal(msg)
		if err != nil {
			return err
		}

		serv.Connect(peer , req , false)  //todo connect remote seed peer , if it's successful , add it into remote_list
	}

	return nil
}

//to connect specified peer
func (serv *netServer) ConnectTo (conn net.Conn , msg []byte , isExist bool) error {
	fmt.Println("p2pServer::ConnectTo()")
	if conn == nil {
		return errors.New("*ERROR* Invalid parameter !!!")
	}

	_  , err := conn.Write(msg)
	if err != nil {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
		return err
	}

	return nil
}

//to connect to certain peer
func (serv *netServer) Connect(addr string , msg []byte , isExist bool) error {
	addr_port := addr+":"+fmt.Sprint(serv.port)
	//fmt.Println("p2pServer::Connect() addr_port = ",addr_port)

	conn , err := net.Dial("tcp", addr_port)
	if err != nil {
		//fmt.Println("*ERROR* Failed to create a connection for remote server !!! err: ",err)
		return err
	}

	//len , err := conn.Write(msg)
	_ , err = conn.Write(msg)
	if err != nil {
		fmt.Println("*ERROR* Failed to send data to the remote server addr !!! err: ",err)
		return err
	}

	//fmt.Println("p2pServer::Connect() len = ",len)

	return nil
}


//to connect certain peer with udp
func (serv *netServer) ConnectUDP(addr string , msg []byte , isExist bool) error {
	fmt.Println("p2pServer::ConnectSeed() addr = ",addr)

	addr_port := addr+":"+fmt.Sprint(serv.port)
	remoteAddr, err := net.ResolveUDPAddr("udp4", addr_port)
	if err != nil {
		fmt.Println("*ERROR* Failed to create a remote server addr !!!")
		return errors.New("*ERROR* Failed to create a remote server addr !!!")
	}

	/*
	//test connection with remote peer
	var msg = message {
		src:      serv.addr,
		dst:      addr,
		msg_type: request,
	}

	req , err := json.Marshal(msg)
	if err != nil {
		return err
	}
	*/

	_ , err = serv.socket.WriteToUDP(msg , remoteAddr)
	if err != nil { //todo check len
		fmt.Println("*ERROR* Failed to send Test message to remote peer !!! ",err)
		return errors.New("*ERROR* Failed to send Test message to remote peer !!!")
	}

	/*
	//package remote peer info as "peer" struct and add it into peer list
	peer := NewPeer(addr)
	peer_identify := Hash(addr_port)
	serv.peerMap[uint64(peer_identify)] = peer
	*/

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

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
