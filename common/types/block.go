package types 

import (
	//"math/big"
)

type Block struct {
	header			Header
	trxs		 	[]Transaction
}

type Header struct {
	PrevBlockHash	[]byte
	Number      	int
}

func NewHeader(prevHash []byte, number int) *Header {
	h := Header{PrevBlockHash: prevHash, Number: number}
	return &h
}

func NewBlock(h *Header) *Block {
	b := Block{header: *h, trxs: []Transaction{}}
	return &b
}

func (b *Block) Decode() error {
	return nil
}

func (b *Block) Encode() error {
	return nil
}

func (b *Block) Hash() []byte {
	return []byte{}
}

func (b *Block) Number() int     { return b.header.Number }