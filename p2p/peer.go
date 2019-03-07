// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description:  producer actor
 * @Author: eripi
 * @Date:   2017-12-06
 * @Last Modified by:
 * @Last Modified time:
 */

package p2p

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	log "github.com/cihub/seelog"
	"io"
	"net"
	"strings"
)

//PeerInfo peer's info
type PeerInfo struct {
	//Id peer id
	Id string
	//Addr peer address
	Addr string
	//Port peer port
	Port string
	//ChainId peer work chain id
	ChainId string
	//Signature peer auth
	Signature []byte
	//Version
	Version uint32
}

// BasicPeerInfo define struct for PeerInfo signature
type BasicPeerInfo struct {
	ChainId string
}

// Hash BasicPeerInfo hash
func (pi *BasicPeerInfo) Hash() common.Hash {
	data, _ := bpl.Marshal(pi)
	temp := sha256.Sum256(data)
	hash := sha256.Sum256(temp[:])
	return hash
}

// Hash PeerInfo hash
func (pi *PeerInfo) Hash() common.Hash {
	data, _ := bpl.Marshal(pi)
	temp := sha256.Sum256(data)
	hash := sha256.Sum256(temp[:])
	return hash
}

// Sign sign a PeerInfo with privkey
func (pi *PeerInfo) Sign(bp BasicPeerInfo, privkey []byte) ([]byte, error) {
	data, err := bpl.Marshal(bp)

	if nil != err {
		return []byte{}, err
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(data)))
	h.Write([]byte(hex.EncodeToString(config.GetChainID())))
	hash := h.Sum(nil)
	signdata, err := crypto.Sign(hash, privkey)

	return signdata, err
}

//Equal peer's info compare
func (a *PeerInfo) Equal(b PeerInfo) bool {
	return (a.Id == b.Id && a.Id != "" && b.Id != "") || (a.Addr == b.Addr && a.Port == b.Port)
}

//IsIncomplete judege peer's info is complete or not
func (a *PeerInfo) IsIncomplete() bool {
	return a.Id == "" || a.Addr == "" || a.Port == "" || a.ChainId == ""
}

//Bigger 1 a > b;  0 a = b ; -1 a < b
func (a *PeerInfo) Bigger(b PeerInfo) int {
	return strings.Compare(a.Id, b.Id)
}

//PeerData peer's key info
type PeerData struct {
	Id    string
	Index uint16
}

//PeerDataSet peer's key info slice
type PeerDataSet []PeerData

//Len length
func (s PeerDataSet) Len() int {
	return len(s)
}

//Less small or not
func (s PeerDataSet) Less(i, j int) bool {
	return s[i].Id > s[j].Id
}

//Swap swap
func (s PeerDataSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//Peer peer...
type Peer struct {
	Info  PeerInfo
	Index uint16

	/*peer state*/
	State int

	conn   net.Conn
	isconn bool
	reader *bufio.Reader

	In bool

	sendup SendupCb
}

//CreatePeer create a instance
func CreatePeer(info PeerInfo, conn net.Conn, in bool, sendup SendupCb) *Peer {
	return &Peer{
		Info:   info,
		State:  PEER_STATE_INIT,
		conn:   conn,
		isconn: true,
		reader: bufio.NewReader(conn),
		In:     in,
		sendup: sendup,
	}
}

//Start start peer routine
func (p *Peer) Start() {
	go p.recvRoutine()
}

//Stop stop peer net conn
func (p *Peer) Stop() {
	p.conn.Close()
}

//Send send a packet
func (p *Peer) Send(packet Packet) error {
	var length uint32
	var head Head
	headsize := uint32(binary.Size(head))

	if packet.Data == nil {
		length = headsize
	} else {
		length = headsize + uint32(len(packet.Data))
	}

	if length > MAX_PACKET_LEN {
		log.Errorf("p2p Send packet length large than max packet length")
		return errors.New("large than max packet length")
	}

	buf := &bytes.Buffer{}
	err := binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		log.Error("p2p send write packet length error")
		return err
	}

	err = binary.Write(buf, binary.BigEndian, packet.H)
	if err != nil {
		log.Error("p2p send write packet protocolType error")
		return err
	}

	_, err = buf.Write(packet.Data)
	if err != nil {
		log.Error("p2p send write packet Data error")
		return err
	}

	if !p.isconn {
		return errors.New("peer disconnected")
	}

	//log.Debugf("p2p peer index: %d send packet %d %d", p.Index, packet.H.ProtocolType, packet.H.PacketType)

	_, err = p.conn.Write(buf.Bytes())
	return err
}

func (p *Peer) recvRoutine() {
	defer p.conn.Close()

	bl := make([]byte, 4)
	var packetLen uint32
	var len int
	var head Head
	readerr := false
	headsize := uint32(binary.Size(head))

	for {
		_, err := io.ReadFull(p.reader, bl)
		if err != nil {
			log.Errorf("p2p recvRoutine read head error:%s,  peer index: %d, %s:%s", err, p.Index, p.Info.Addr, p.Info.Port)
			p.isconn = false
			return
		}

		packetLen = binary.BigEndian.Uint32(bl)
		if packetLen < headsize || packetLen > MAX_PACKET_LEN {
			log.Errorf("p2p recvRoutine drop packet wrong packet lenght %d", packetLen)
			continue
		}

		buf := make([]byte, packetLen)
		len, err = io.ReadFull(p.reader, buf)
		if err != nil {
			log.Errorf("p2p recvRoutine read data error:%s,  peer index: %d, %s:%s", err, p.Index, p.Info.Addr, p.Info.Port)
			p.isconn = false
			return
		}

		if uint32(len) < packetLen {
			for {
				length, err := io.ReadFull(p.reader, buf[len:])
				if err != nil {
					log.Errorf("p2p recvRoutine continue read data error:%s,  peer index: %d, %s:%s", err, p.Index, p.Info.Addr, p.Info.Port)
					p.isconn = false
					return
				}

				len += length

				if uint32(len) < packetLen {
					continue
				} else if uint32(len) == packetLen {
					break
				} else {
					log.Errorf("p2p recvRoutine continue read data length wrong packet length:%d, read:%d", packetLen, len)
					readerr = true
					break
				}
			}
		}

		if readerr {
			readerr = false
			continue
		}

		var packet Packet

		packet.H.ProtocolType = uint16(binary.BigEndian.Uint16(buf))
		packet.H.PacketType = uint16(binary.BigEndian.Uint16(buf[2:]))

		if packetLen > headsize {
			packet.Data = buf[headsize:packetLen]
		}

		p.sendup(p.Index, &packet)
	}

	p.isconn = false
}
