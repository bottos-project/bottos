// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

// This program is free software: you can distribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Bottos.  If not, see <http://www.gnu.org/licenses/>.

// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


/*
 * file description: the interface for WASM execution
 * @Author: Richard
 * @Date:   2018-02-10
 * @Last Modified by:
 * @Last Modified time:
 */

package discover

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/bottos-project/bottos/p2p/nat"

)

const (
	Version = 3
	macSize  = 256 / 8
	sigSize  = 520 / 8
	headSize = macSize + sigSize 
	respTimeout = 600 * time.Millisecond
	sendTimeout = 600 * time.Millisecond
	expiration  = 30 * time.Second
	refreshInterval = 1 * time.Hour
	pingPacket = iota + 1
	pongPacket
	findnodePacket
	neighborsPacket
)

var (
	headSpace = make([]byte, headSize)
	maxNeighbors int
	errTooSmall   = errors.New("udpPacket is too small")
	errBadHash          = errors.New("a bad hash")
	errExpired          = errors.New("expired")
	errBadVersion       = errors.New("version illegal")
	errUnsolicitedReply = errors.New("unsolicited udpReply")
	errUnknownNode      = errors.New("illegal node")
	errTimeout          = errors.New("timeout")
	errClosed           = errors.New("socket closed")
)


type (
	ping struct {
		Version    uint
		From, To   rpcEndpoint
		Expiration uint64
	}
	pong struct {
		To rpcEndpoint
		ReplyTok   []byte 
		Expiration uint64 
	}
	findnode struct {
		Target     NodeID 
		Expiration uint64
	}
	neighbors struct {
		Nodes      []rpcNode
		Expiration uint64
	}
	rpcNode struct {
		IP  net.IP 
		UDP uint64 
		TCP uint64 
		ID  NodeID
	}
	rpcEndpoint struct {
		IP  net.IP 
		UDP uint64 
		TCP uint64 
	}
)

type udp struct {
	udpConn        udpConn
	priv        *ecdsa.PrivateKey
	ourEndpoint rpcEndpoint
	addpending chan *pending
	gotreply   chan udpReply
	closing chan struct{}
	nat     string
	*Table
}
type udpPacket interface {
	handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error
}
type udpConn interface {
	ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error)
	WriteToUDP(b []byte, addr *net.UDPAddr) (n int, err error)
	Close() error
	LocalAddr() net.Addr
}
type udpReply struct {
	from  NodeID
	ptype byte
	data  interface{}
	matched chan<- bool
}



type pending struct {
	from  NodeID
	ptype byte
	deadline time.Time
	callback func(resp interface{}) (done bool)
	errc chan<- error
}


func makeUdpEndpoint(addr *net.UDPAddr, tcpPort uint64) rpcEndpoint {
	ip := addr.IP.To4()
	if ip == nil {
		ip = addr.IP.To16()
	}
	return rpcEndpoint{IP: ip, UDP: uint64(addr.Port), TCP: tcpPort}
}

func newNodeFromRPC(rn rpcNode) (n *Node, valid bool) {
	if rn.IP.IsMulticast() || rn.IP.IsUnspecified() || rn.UDP == 0 {
		return nil, false
	}
	return newNodeInfo(rn.ID, rn.IP, rn.UDP, rn.TCP), true
}

func RPCNode(n *Node) rpcNode {
	return rpcNode{ID: n.ID, IP: n.IP, UDP: n.UDP, TCP: n.TCP}
}



func ListenUDP(priv *ecdsa.PrivateKey, laddr string, natm nat.Interface, nodeDBPath string) (*Table, error) {
	addr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	tab, _ := newUDP(priv, udpConn, natm, nodeDBPath)
	return tab, nil
}

func newUDP(priv *ecdsa.PrivateKey, c udpConn, natm nat.Interface, nodeDBPath string) (*Table, *udp) {
	udp := &udp{
		udpConn:       c,
		priv:       priv,
		closing:    make(chan struct{}),
		gotreply:   make(chan udpReply),
		addpending: make(chan *pending),
	}
	realaddr := c.LocalAddr().(*net.UDPAddr)
	udp.ourEndpoint = makeUdpEndpoint(realaddr, uint64(realaddr.Port))
	udp.Table = newTableInfo(udp, PublickeyID(&priv.PublicKey), realaddr, nodeDBPath)
	go udp.loop()
	go udp.readLoop()
	return udp.Table, udp
}

func (t *udp) close() {
	close(t.closing)
	t.udpConn.Close()
}

func (t *udp) ping(toid NodeID, toaddr *net.UDPAddr) error {
	errc := t.pending(toid, pongPacket, func(interface{}) bool { return true })
	t.send(toaddr, pingPacket, ping{
		Version:    Version,
		From:       t.ourEndpoint,
		To:         makeUdpEndpoint(toaddr, 0), // TODO: maybe use known TCP port from DB
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
	return <-errc
}

func (t *udp) waitping(from NodeID) error {
	return <-t.pending(from, pingPacket, func(interface{}) bool { return true })
}


func (t *udp) loop() {
	var (
		pending      []*pending
		nextDeadline time.Time
		timeout      = time.NewTimer(0)
		refresh      = time.NewTicker(refreshInterval)
	)
	<-timeout.C 
	defer refresh.Stop()
	defer timeout.Stop()
	rearmTimeout := func() {
		now := time.Now()
		if len(pending) == 0 || now.Before(nextDeadline) {
			return
		}
		nextDeadline = pending[0].deadline
		timeout.Reset(nextDeadline.Sub(now))
	}
	for {
		select {
		case r := <-t.gotreply:
			var matched bool
			for i := 0; i < len(pending); i++ {
				if p := pending[i]; p.from == r.from && p.ptype == r.ptype {
					matched = true
					if p.callback(r.data) {	
						p.errc <- nil
						copy(pending[i:], pending[i+1:])
						pending = pending[:len(pending)-1]
						i--
					}
				}
			}
			r.matched <- matched
		case p := <-t.addpending:
			p.deadline = time.Now().Add(respTimeout)
			pending = append(pending, p)
			rearmTimeout()
		case <-refresh.C:
			go t.refresh()
		case now := <-timeout.C:
			i := 0
			for ; i < len(pending) && now.After(pending[i].deadline); i++ {
				pending[i].errc <- errTimeout
			}
			if i > 0 {
				copy(pending, pending[i:])
				pending = pending[:len(pending)-i]
			}
			rearmTimeout()
		case <-t.closing:
			for _, p := range pending {
				p.errc <- errClosed
			}
			pending = nil
			return
		}
	}
}


func init() {
	p := neighbors{Expiration: ^uint64(0)}
	maxSizeNode := rpcNode{IP: make(net.IP, 16), UDP: ^uint64(0), TCP: ^uint64(0)}
	for n := 0; ; n++ {
		p.Nodes = append(p.Nodes, maxSizeNode)
		if headSize+1 >= 1280 {
			maxNeighbors = n
			break
		}
	}
}
func (t *udp) handleReply(from NodeID, ptype byte, req udpPacket) bool {
	matched := make(chan bool)
	select {
	case t.gotreply <- udpReply{from, ptype, req, matched}:

		return <-matched
	case <-t.closing:
		return false
	}
}
func (req *findnode) handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error {
	if expired(req.Expiration) {
		return errExpired
	}
	if t.db.node(fromID) == nil {
		return errUnknownNode
	}
	var target  string 
	t.mutex.Lock()
	closest := t.closest(target, bucketSize).entries
	t.mutex.Unlock()
	p := neighbors{Expiration: uint64(time.Now().Add(expiration).Unix())}
	for i, n := range closest {
		p.Nodes = append(p.Nodes, RPCNode(n))
		if len(p.Nodes) == maxNeighbors || i == len(closest)-1 {
			t.send(from, neighborsPacket, p)
			p.Nodes = p.Nodes[:0]
		}
	}
	return nil
}
func (req *neighbors) handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error {
	if expired(req.Expiration) {
		return errExpired
	}
	if !t.handleReply(fromID, neighborsPacket, req) {
		return errUnsolicitedReply
	}
	return nil
}

func (req *pong) handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error {
	if expired(req.Expiration) {
		return errExpired
	}
	if !t.handleReply(fromID, pongPacket, req) {
		return errUnsolicitedReply
	}
	return nil
}
func (req *ping) handle(t *udp, from *net.UDPAddr, fromID NodeID, mac []byte) error {
	if req.Version != Version {
		return errBadVersion
	}
	if expired(req.Expiration) {
		return errExpired
	}
	t.send(from, pongPacket, pong{
		To:         makeUdpEndpoint(from, req.From.TCP),
		ReplyTok:   mac,
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
	if !t.handleReply(fromID, pingPacket, req) {
		go t.bond(true, fromID, from, req.From.TCP)
	}
	return nil
}




func (t *udp) findnode(toid NodeID, toaddr *net.UDPAddr, target NodeID) ([]*Node, error) {
	nrec := 0
	nodes := make([]*Node, 0, bucketSize)
	errc := t.pending(toid, neighborsPacket, func(r interface{}) bool {
		udpReply := r.(*neighbors)
		for _, rn := range udpReply.Nodes {
			nrec++
			if n, valid := newNodeFromRPC(rn); valid {
				nodes = append(nodes, n)
			}
		}
		return nrec >= bucketSize
	})
	t.send(toaddr, findnodePacket, findnode{
		Target:     target,
		Expiration: uint64(time.Now().Add(expiration).Unix()),
	})
	err := <-errc
	return nodes, err
}

func (t *udp) pending(id NodeID, ptype byte, callback func(interface{}) bool) <-chan error {
	ch := make(chan error, 1)
	p := &pending{from: id, ptype: ptype, callback: callback, errc: ch}
	select {
	case t.addpending <- p:
	case <-t.closing:
		ch <- errClosed
	}
	return ch
}



func expired(ts uint64) bool {
	return time.Unix(int64(ts), 0).Before(time.Now())
}

func (t *udp) send(toaddr *net.UDPAddr, ptype byte, req interface{}) error {
	udpPacket, err := encodePacket(t.priv, ptype, req)
	if err != nil {
		return err
	}
	
	if _, err = t.udpConn.WriteToUDP(udpPacket, toaddr); err != nil {
		
	}
	return err
}

func encodePacket(priv *ecdsa.PrivateKey, ptype byte, req interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	b.Write(headSpace)
	b.WriteByte(ptype)
	udpPacket := b.Bytes()
	sig := udpPacket
	copy(udpPacket[macSize:], sig)
	return udpPacket, nil
}
func (t *udp) readLoop() {
	defer t.udpConn.Close()
	buf := make([]byte, 1280)
	for {
		nbytes, from, err := t.udpConn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		t.handlePacket(from, buf[:nbytes])
	}
}

func (t *udp) handlePacket(from *net.UDPAddr, buf []byte) error {
	udpPacket, fromID, hash, err := decodePacket(buf)
	if err != nil {
		return err
	}
	if err = udpPacket.handle(t, from, fromID, hash); err != nil {
	}
	return err
}

func decodePacket(buf []byte) (udpPacket, NodeID, []byte, error) {
	if len(buf) < headSize+1 {
		return nil, NodeID{}, nil, errTooSmall
	}
	hash, sig, sigdata := buf[:macSize], buf[macSize:headSize], buf[headSize:]
	fmt.Println(sig)
	var  shouldhash []byte
	if !bytes.Equal(hash, shouldhash) {
		return nil, NodeID{}, nil, errBadHash
	}
	var fromID NodeID 
	var err error
	var req udpPacket
	switch ptype := sigdata[0]; ptype {
	case findnodePacket:
		req = new(findnode)
	case neighborsPacket:
		req = new(neighbors)
	case pingPacket:
		req = new(ping)
	case pongPacket:
		req = new(pong)
	default:
		return nil, fromID, hash, fmt.Errorf("unknown type: %d", ptype)
	}
	return req, fromID, hash, err
}
