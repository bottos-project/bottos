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
