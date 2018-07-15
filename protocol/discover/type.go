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

package discover

import "github.com/bottos-project/bottos/p2p"

//DO NOT EDIT
const (
	PEER_INFO_REQ = 1
	PEER_INFO_RSP = 2

	PEER_HANDSHAKE_REQ     = 3
	PEER_HANDSHAKE_RSP     = 4
	PEER_HANDSHAKE_RSP_ACK = 5

	PEER_NEIGHBOR_REQ = 7
	PEER_NEIGHBOR_RSP = 8

	PEER_PING = 9
	PEER_PONG = 10
)

//PeerInfoReq peer info
type PeerInfoReq struct {
}

//PeerInfoRsp peer info response
type PeerInfoRsp struct {
	Info p2p.PeerInfo
}

//PeerNeighborReq pne exchange request
type PeerNeighborReq struct {
}

//PeerNeighborRsp pne exchange response
type PeerNeighborRsp struct {
	Neighbor []p2p.PeerInfo
}
