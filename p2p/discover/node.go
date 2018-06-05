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
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const nodeIDBitLength = 512


type Node struct {
	IP       net.IP 
	UDP, TCP uint64 
	ID       NodeID
	hash      string
}

type NodeID [nodeIDBitLength / 8]byte

var ncount = [256]int{
	8, 7, 6, 6, 5, 5, 5, 5,
	4, 4, 4, 4, 4, 4, 4, 4,
	3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

func newNodeInfo(id NodeID, ip net.IP, udpPort, tcpPort uint64) *Node {
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}
	return &Node{
		IP:  ip,
		UDP: udpPort,
		TCP: tcpPort,
		ID:  id,
		
	}
}

func (n *Node) addr() *net.UDPAddr {
	return &net.UDPAddr{IP: n.IP, Port: int(n.UDP)}
}


func (n *Node) String() string {
	addr := net.TCPAddr{IP: n.IP, Port: int(n.TCP)}
	u := url.URL{
		Scheme: "enode",
		User:   url.User(fmt.Sprintf("%x", n.ID[:])),
		Host:   addr.String(),
	}
	if n.UDP != n.TCP {
		u.RawQuery = "discport=" + strconv.Itoa(int(n.UDP))
	}
	return u.String()
}


func ParseNodeInfo(rawurl string) (*Node, error) {
	var (
		id               NodeID
		ip               net.IP
		tcpPort, udpPort uint64
	)
	u, err := url.Parse(rawurl)
	if u.Scheme != "enode" {
		return nil, errors.New("invalid URL, set like \"enode\"")
	}
	
	if u.User == nil {
		return nil, errors.New("missing nodeID")
	}
	if id, err = HexNodeID(u.User.String()); err != nil {
		return nil, fmt.Errorf("invalid nodeID: %v", err)
	}
	
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, fmt.Errorf("invalid host: %v", err)
	}
	if ip = net.ParseIP(host); ip == nil {
		return nil, errors.New("invalid IP")
	}
	
	if ipv4 := ip.To4(); ipv4 != nil {
		ip = ipv4
	}
	
	if tcpPort, err = strconv.ParseUint(port, 10, 16); err != nil {
		return nil, errors.New("invalid port")
	}
	udpPort = tcpPort
	qv := u.Query()
	if qv.Get("discport") != "" {
		udpPort, err = strconv.ParseUint(qv.Get("discport"), 10, 16)
		if err != nil {
			return nil, errors.New("invalid discport")
		}
	}
	return newNodeInfo(id, ip, uint64(udpPort), uint64(tcpPort)), nil
}


func MustParseNode(rawurl string) *Node {
	n, err := ParseNodeInfo(rawurl)
	if err != nil {
		return nil
	}
	return n
}





func (n NodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}


func (n NodeID) GoString() string {
	return fmt.Sprintf("discover.HexNodeID(\"%x\")", n[:])
}


func HexNodeID(in string) (NodeID, error) {
	var id NodeID
	if strings.HasPrefix(in, "0x") {
		in = in[2:]
	}

	b, err := hex.DecodeString(in)
	if err != nil {
		return id, err
	} else if len(b) != len(id) {
		return id, fmt.Errorf("invalid length, need %d hex bytes", len(id))
	}
	copy(id[:], b)
	return id, nil
}


func MustHexNodeID(in string) NodeID {
	id, err := HexNodeID(in)
	if err != nil {
		panic(err)
	}
	return id
}


func PublickeyID(pub *ecdsa.PublicKey) NodeID {
	var id NodeID
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != len(id) {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}


func (id NodeID) Publickey() (*ecdsa.PublicKey, error) {
	p := &ecdsa.PublicKey{Curve: nil, X: new(big.Int), Y: new(big.Int)}
	half := len(id) / 2
	p.X.SetBytes(id[:half])
	p.Y.SetBytes(id[half:])
	if !p.Curve.IsOnCurve(p.X, p.Y) {
		return nil, errors.New("not on the curve")
	}
	return p, nil
}


func recoverNode(hash, sig []byte, pubkey string) (id NodeID, err error) {
	
	if len(pubkey)-1 != len(id) {
		return id, fmt.Errorf("invalid length, has %d bits, want %d bits", len(pubkey)*8, (len(id)+1)*8)
	}
	for i := range id {
		id[i] = pubkey[i+1]
	}
	return id, nil
}

func distancecmp(target, x, y string) int {
	for i := range target {
		dx := x[i] ^ target[i]
		dy := y[i] ^ target[i]
		if dx < dy {
			return -1
		} else if dx > dy {
			return 1
		}
	}
	return 0
}



func logdistance(x, y string) int {
	lz := 0
	for i := range x {
		p := x[i] ^ y[i]
		if p != 0 {
			lz += ncount[p]
			break
		} else {
			lz += 8
		}
	}
	return len(x)*8 - lz
}


func nodeDistance(x string, n int) (y string) {
	if n == 0 {
		return x
	}
	y = x
	pos := len(x) - n/8 - 1
	bit := byte(0x01) << (byte(n%8) - 1)
	if bit == 0 {
		pos++
		bit = 0x80
	}
	return y
}
