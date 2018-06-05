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

// Package exec provides functions for executing WebAssembly bytecode.

/*
 * file description: the interface for WASM execution
 * @Author: Stewart Li
 * @Date:   2018-02-08
 * @Last Modified by:
 * @Last Modified time:
 */

package p2pserver

import (
	"hash/fnv"
	"crypto/rsa"
	"unsafe"
)

//message type
const (
	REQUEST = iota //0
	RESPONSE
	CRX_BROADCAST
	BLK_BROADCAST
	OTHER
)

//connection state
const (
	ESTABLISH  = iota //receive peer`s verack
	INACTIVITY        //link broken
)

//p2p call type
const (
	TRANSACTION = iota
	BLOCK
)

type message struct {
	Src     string
	Dst     string
	MsgType uint8
	Content []byte
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Hash(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

type RsaKeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type Key interface {
	Bytes() ([]byte, error)
	Equals(Key) bool
}

type PrivKey interface {
	Key
	Sign([]byte) ([]byte, error)
	GetPublicKey() PubKey
}

type PubKey interface {
	Key
	VerifyKey(data []byte, sig []byte) (bool, error)
}

