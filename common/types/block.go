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
 * file description:  general Hash type
 * @Author: Gong Zibin
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package types 

import (
	//"math/big"
	"fmt"
	"bytes"
	"io"
	"crypto/sha256"

	proto "github.com/golang/protobuf/proto"
)

func NewBlock(h *Header, txs []*Transaction) *Block {
	b := Block{header: copyHeader(h)}

	if len(txs) == 0 {
	} else {
		b.transactions = make([]*Transaction, len(txs))
		copy(b.transactions, txs)
	}

	return &b
}

func (b *Block) Hash() Hash {
	return b.header.Hash()
}

func (h *Header) Hash() Hash {
	data, _ := proto.Marshal(h)
	h := sha256.Sum256(data)
	return h
}

func copyHeader(h *Header) *Header {
	cpy := *h

	// TODO

	return &cpy
}

func (b *Block) GetPrevBlockHash() Hash {
	bh := b.GetHeader().GetPrevBlockHash()
	return BytesToHash(bh)
}

func (b *Block) GetNumber() uint32 { 
	return b.GetHeader().GetNumber()
}

func (b *Block) GetTimestamp() uint64 { 
	return b.GetHeader().GetTimestamp()
}

func (b *Block) GetMerkleRoot() Hash {
	bh := b.GetHeader().GetMerkleRoot()
	return BytesToHash(bh)
}

// TODO AccountName Type
func (b *Block) GetProducer() []byte {
	return b.GetHeader().GetProducer()
}

//func (b *Block) GetProducerChange() AccountName {
//	return b.header.Producer
//}

func (b *Block) GetProducerSign() Hash {
	bh := b.GetHeader().GetProducerSign()
	return BytesToHash(bh)
}
