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

//UniMsgPacket it is a unicast packet , Index is peer id to send to
type UniMsgPacket struct {
	Index uint16
	P     Packet
}

//BcastMsgPacket it is a multicast packet , Indexs is filter peers index which not send to
type BcastMsgPacket struct {
	Indexs []uint16
	P      Packet
}

//DO NOT EDIT
const (
	MAX_PACKET_LEN = 10000000
)

//Head packet head
type Head struct {
	ProtocolType uint16
	PacketType   uint16
	Pad          uint16
}

//Packet packet
type Packet struct {
	H    Head
	Data []byte
}
