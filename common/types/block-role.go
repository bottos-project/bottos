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
	//"fmt"
	//"bytes"
	//"io"
	"crypto/sha256"
	"github.com/bottos-project/core/common"
	"github.com/golang/protobuf/proto"

)

func NewBlock(h *Header, txs []*Transaction) *Block {
	b := Block{Header: copyHeader(h)}

	if len(txs) == 0 {
	} else {
		b.Transactions = make([]*Transaction, len(txs))
		copy(b.Transactions, txs)
	}

	return &b
}

func (b *Block) Hash() common.Hash {
	return b.Header.Hash()
}

func (h *Header) Hash() common.Hash {
	data, _ := proto.Marshal(h)
	hash := sha256.Sum256(data)
	return hash
}

func copyHeader(h *Header) *Header {
	cpy := *h

	// TODO

	return &cpy
}

func (b *Block) GetPrevBlockHash() common.Hash {
	bh := b.GetHeader().GetPrevBlockHash()
	return common.BytesToHash(bh)
}

func (b *Block) GetNumber() uint32 { 
	return b.GetHeader().GetNumber()
}

func (b *Block) GetTimestamp() uint64 { 
	return b.GetHeader().GetTimestamp()
}

func (b *Block) GetMerkleRoot() common.Hash {
	bh := b.GetHeader().GetMerkleRoot()
	return common.BytesToHash(bh)
}

// TODO AccountName Type
func (b *Block) GetProducer() []byte {
	return b.GetHeader().GetProducer()
}

//func (b *Block) GetProducerChange() AccountName {
//	return b.header.Producer
//}

func (b *Block) GetProducerSign() common.Hash {
	bh := b.GetHeader().GetProducerSign()
	return common.BytesToHash(bh)
}

func (b *Block) ValidateSign() bool {
	return true
} 

func (b *Block) GetTransactionByHash(hash common.Hash) *Transaction {
	for _, transaction := range b.Transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}

	return nil
} 
