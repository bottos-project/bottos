package p2pserver

import (
	//"io"
	"fmt"
	"net"
	"sync"
	"errors"
	//"crypto"
	//"crypto/rand"
	"crypto/rsa"
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

	/*
	prvKey, pubKey, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	fmt.Println("prvKey = ",prvKey," , pubKey = ",pubKey)
	*/

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


type RsaKeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// Key represents a crypto key that can be compared to another key
type Key interface {
	// Bytes returns a serialized, storeable representation of this key
	Bytes() ([]byte, error)

	// Equals checks whether two PubKeys are the same
	Equals(Key) bool
}

// PrivKey represents a private key that can be used to generate a public key,
// sign data, and decrypt data that was encrypted with a public key
type PrivKey interface {
	Key

	// Cryptographically sign the given bytes
	Sign([]byte) ([]byte, error)

	// Return a public key paired with this private key
	GetPublic() PubKey
}

type PubKey interface {
	Key

	// Verify that 'sig' is the signed hash of 'data'
	Verify(data []byte, sig []byte) (bool, error)
}


/*
// Generates a keypair
func GenerateKeyPairWithReader(typ, bits int, src io.Reader) (PrivKey, PubKey, error) {

	privateKey, err := rsa.GenerateKey(src, bits)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return &RsaKeyPair{privateKey:privateKey}, &RsaKeyPair{ publicKey:publicKey}, nil
}
*/







