package types 

import (
	//"math/big"
	"fmt"
	"bytes"
	"io"
	"crypto/sha256"

	"github.com/bottos-project/bottos/core/library"
)

type Block struct {
	header			*Header
	transactions	[]*Transaction
}

type Header struct {
	PrevBlockHash	library.Hash		// Hash of Previos block
	Number      	uint32				// Block Number
	Timestamp       uint32         		// Creation time
	MerkleRoot		library.Hash
	Producer		library.AccountName
	ProducerChange	[]library.AccountName
	ProducerSign	library.Hash	// TODO ECSDA sign type
}

func NewBlock(h *Header, txs []*Transaction) *Block {
	b := Block{header: copyHeader(h)}

	if len(txs) == 0 {
	} else {
		// TODO Compute Hash
		b.transactions = make([]*Transaction, len(txs))
		copy(b.transactions, txs)
	}

	return &b
}

func (h *Header) Serialize(w io.Writer) error {

	return nil
}

func (h *Header) Deserialize(r io.Reader) error {

	return nil
}

func (b *Block) Serialize(w io.Writer) error {
	b.header.Serialize(w)

	return nil
}


func (b *Block) Deserialize(r io.Reader) error {
	var header Header
	err := header.Deserialize(r)
	if err != nil {
		return fmt.Errorf("Header Deserialize failed: %s", err)
	}

	// TODO

	b.header = &header
	return nil
}

func (h *Header) Hash() library.Hash {
	value := bytes.NewBuffer(nil)
	h.Serialize(value)
	temp := sha256.Sum256(value.Bytes())
	return temp
}

func (b *Block) Hash() library.Hash {
	return b.header.Hash()
}

func copyHeader(h *Header) *Header {
	cpy := *h

	// TODO

	return &cpy
}

func (b *Block) PrevBlockHash() library.Hash		{ return b.header.PrevBlockHash }
func (b *Block) Number() uint32     				{ return b.header.Number }
func (b *Block) Time() uint32						{ return b.header.Timestamp }
func (b *Block) MerkleRoot() library.Hash			{ return b.header.MerkleRoot }
func (b *Block) Producer() library.AccountName		{ return b.header.Producer }

func (b *Block) Header() *Header				{ return copyHeader(b.header) }

func (b *Block) Transactions() []*Transaction 	{ return b.transactions }