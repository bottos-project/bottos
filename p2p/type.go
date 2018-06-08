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
	"unsafe"

	log "github.com/cihub/seelog"
)

//message type
const (
	REQUEST = iota //0
	RESPONSE
	CRX_BROADCAST
	BLK_BROADCAST
	PEERNEIGHBOR_REQ
	PEERNEIGHBOR_RSP
	DEFAULT
)

//connection state
const (
	ESTABLISH  = 10 //receive peer`s verack
	INACTIVITY = 11 //link broken
)

//p2p call type
const (
	TRANSACTION = 100
	BLOCK       = 101
	BLOCK_INFO  = 102
	BLOCK_REQ   = 103
	BLOCK_RES   = 104
)

const (
	BLOCK_PRINT        = 30
	RED_PRINT          = 31
	GREEN_PRINT        = 32
	YELLO_PRINT        = 33
	BLUE_PRINT         = 34
	PURPLISH_RED_PRINT = 35
	AUQA_PRINT         = 36
	WHITE_PRINT        = 37
)

//CommonMessage message struct
type CommonMessage struct {
	Src     string
	Dst     string
	MsgType uint8
	Content []byte
}

//BlockInfo block info
type BlockInfo struct {
	BlockNum  uint32
	HeaderNum uint32
}

//BlockReq block request message
type BlockReq struct {
	BlockNum uint32
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//Hash hash calc
func Hash(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

//SuperPrint log with color
func SuperPrint(color uint8, args ...interface{}) {
	for _, v := range args {
		log.Infof("%c[%d;%d;%dm%v%c[0m", 0x1B, 123, 40, color, v, 0x1B)
	}
	log.Info("\n")
}

