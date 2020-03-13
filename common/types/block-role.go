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
 * file description:  block
 * @Author: Gong Zibin
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package types

import (
	"crypto/sha256"

	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/signature"
	log "github.com/cihub/seelog"
)

type BlockV0 struct {
	Header       *Header
	Transactions []*Transaction
	ValidatorSet []*Validator
}

//how to add field when upgrade version: NewField   uint64 `version:"x.x.x"`，x.x.x is new version
type Block struct {
	Header       *Header
	BlockTransactions []*BlockTransaction
	//Transactions []*Transaction
	ValidatorSet []*Validator
}

type Header struct {
	Version         uint32
	PrevBlockHash   []byte
	Number          uint64
	Timestamp       uint64
	MerkleRoot      []byte
	Delegate        []byte
	DelegateSign    []byte
	DelegateChanges []string
}

type BlockDetail struct {
	BlockVersion      uint32        `json:"block_version"`
	PrevBlockHash    string         `json:"prev_block_hash"`
	BlockNum         uint64         `json:"block_num"`
	BlockHash        string         `json:"block_hash"`
	CursorBlockLabel uint32         `json:"cursor_block_label"`
	BlockTime        uint64         `json:"block_time"`
	TrxMerkleRoot    string         `json:"trx_merkle_root"`
	Delegate         string         `json:"delegate"`
	DelegateSign     string         `json:"delegate_sign"`
	Trxs             []*interface{} `json:"trxs"`
}

func NewBlock(h *Header, txs []*BlockTransaction) *Block {
	b := Block{Header: copyHeader(h)}

	if len(txs) > 0 {
		b.BlockTransactions = make([]*BlockTransaction, len(txs))
		copy(b.BlockTransactions, txs)
	}

	b.Header.MerkleRoot = b.ComputeMerkleRoot().Bytes()

	return &b
}

func NewHeader(version uint32) *Header {
	h := &Header{Version: version}

	return h
}

func (b *Block) Hash() common.Hash {
	return b.Header.Hash()
}

func (h *Header) Hash() common.Hash {
	nh := copyHeader(h)
	nh.DelegateSign = []byte{}
	data, _ := bpl.Marshal(nh)
	temp := sha256.Sum256(data)
	hash := sha256.Sum256(temp[:])
	return hash
}

func copyHeader(h *Header) *Header {
	cpy := *h
	return &cpy
}

func (b *Block) Copy() *Block {
	if b == nil {
		return nil
	}
	block := Block{Header: copyHeader(b.Header)}
	if len(b.BlockTransactions) > 0 {
		block.BlockTransactions = make([]*BlockTransaction, len(b.BlockTransactions))
		copy(block.BlockTransactions, b.BlockTransactions)

	}
	//map copy
	if len(b.ValidatorSet) > 0 {
		block.ValidatorSet = make([]*Validator, len(b.ValidatorSet))
		copy(block.ValidatorSet, b.ValidatorSet)
	}
	return &block
}

func (b *Block) InitFromV0(oldblock *BlockV0) {
	b.Header = copyHeader(oldblock.Header)
	if len(oldblock.Transactions) > 0 {
		b.BlockTransactions = make([]*BlockTransaction, len(oldblock.Transactions))
		for i := 0; i < len(oldblock.Transactions); i++ {
			b.BlockTransactions[i] = &BlockTransaction{}
			b.BlockTransactions[i].Transaction = oldblock.Transactions[i]
		}
	}
	//map copy
	if len(oldblock.ValidatorSet) > 0 {
		b.ValidatorSet = make([]*Validator, len(oldblock.ValidatorSet))
		copy(b.ValidatorSet, oldblock.ValidatorSet)
	}
}

func (b *Block) CovertToV0(oldblock *BlockV0) {
	oldblock.Header = copyHeader(b.Header)
	if len(b.BlockTransactions) > 0 {
		oldblock.Transactions = make([]*Transaction, len(b.BlockTransactions))
		for i := 0; i < len(b.BlockTransactions); i++ {
			oldblock.Transactions[i] = b.BlockTransactions[i].Transaction
		}
	}
	//map copy
	if len(b.ValidatorSet) > 0 {
		oldblock.ValidatorSet = make([]*Validator, len(b.ValidatorSet))
		copy(oldblock.ValidatorSet, b.ValidatorSet)
	}
}

func (b *Block) GetPrevBlockHash() common.Hash {
	bh := b.Header.GetPrevBlockHash()
	return common.BytesToHash(bh)
}

func (b *Block) GetNumber() uint64 {
	return b.Header.GetNumber()
}

func (b *Block) GetTimestamp() uint64 {
	return b.Header.GetTimestamp()
}

func (b *Block) GetMerkleRoot() common.Hash {
	bh := b.Header.GetMerkleRoot()
	return common.BytesToHash(bh)
}

func (b *Block) ComputeMerkleRoot() common.Hash {
	if len(b.BlockTransactions) > 0 {
		var hs []common.Hash
		for _, tx := range b.BlockTransactions {
			hs = append(hs, tx.Hash())
		}
		return common.ComputeMerkleRootHash(hs)
	}
	return common.Hash{}
}

func (b *Block) GetDelegate() []byte {
	return b.Header.GetDelegate()
}

func (b *Block) GetDelegateSign() []byte {
	sign := b.Header.GetDelegateSign()
	return sign
}

func (b *Block) GetTransactionByHash(hash common.Hash) *BlockTransaction {
	for _, transaction := range b.BlockTransactions {
		if transaction.Transaction.Hash() == hash {
			return transaction
		}
	}

	return nil
}

func (b *Block) GetVersion() uint32 {
	return b.Header.GetVersion()
}

func (m *Header) GetPrevBlockHash() []byte {
	if m != nil {
		return m.PrevBlockHash
	}
	return nil
}

func (m *Header) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

func (m *Header) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Header) GetMerkleRoot() []byte {
	if m != nil {
		return m.MerkleRoot
	}
	return nil
}

func (m *Header) GetDelegate() []byte {
	if m != nil {
		return m.Delegate
	}
	return nil
}

func (m *Header) GetDelegateSign() []byte {
	if m != nil {
		return m.DelegateSign
	}
	return nil
}

func (m *Header) GetDelegateChanges() []string {
	if m != nil {
		return m.DelegateChanges
	}
	return nil
}

func (m *Header) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (b *Block) SignVote(account string, vote *Validator) (*Validator, error) {
	myVote := vote.Copy()
	digest := signDigest(b, myVote)
	signature, err := signValidator(account, digest)
	if err != nil {
		log.Errorf("COMMON SignVote delegate %s, voteInfo.Height %v, voteInfo.Round %v,voteInfo.Step %v, voteInfo.VoteResult %v,DelegateSignature %x",
			myVote.Delegate, myVote.VoteInfo.Height, myVote.VoteInfo.Round, myVote.VoteInfo.Step, myVote.VoteInfo.VoteResult, signature)

		log.Errorf("COMMON SignVote failed: signdata %x, hash %x, account=%s", signature, digest.Bytes(), account)
		return nil, err
	}

	myVote.DelegateSignature = signature

	return myVote, nil
}

func (b *Block) VerifyVote(pubkey []byte, vote *Validator) bool {
	myVote := vote.Copy()
	digest := signDigest(b, myVote)
	signdata := myVote.DelegateSignature
	result := signature.VerifySign(pubkey, digest.Bytes(), signdata)

	if result == false {
		log.Errorf("COMMON VerifyVote delegate %s, voteInfo.Height %v, voteInfo.Round %v,voteInfo.Step %v, voteInfo.VoteResult %v,DelegateSignature %x",
			myVote.Delegate, myVote.VoteInfo.Height, myVote.VoteInfo.Round, myVote.VoteInfo.Step, myVote.VoteInfo.VoteResult, myVote.DelegateSignature)

		log.Errorf("COMMON VerifyVote failed: signdata %x, pubkey %x, hash %x, result=%v", signdata, pubkey, digest.Bytes(), result)
		return result
	}
	return result
}

func (b *Block) MarshalCompatibly() ([]byte, error) {
	if b.GetVersion() == 0 {
		oldblock := &BlockV0{}
		b.CovertToV0(oldblock)
		return bpl.Marshal(oldblock)
	} else {
		return bpl.Marshal(b)
	}
}


func (b *Block) Marshal() ([]byte, error) {
	return bpl.Marshal(b)
}

func (b *Block) UnmarshalCompatibly(data []byte) error {
	if err := bpl.UnmarshalUntilField(data, b, "Version"); err != nil { //unmarshal version field first
		return err
	}
	version := b.GetVersion()
	if version == 0 {
		oldblock := &BlockV0{}
		if err := bpl.Unmarshal(data, oldblock); err != nil {
			return err
		}
		b.InitFromV0(oldblock)
	} else {
		if err := bpl.Unmarshal(data, b); err != nil {
			return err
		}
	}
	return nil
}

func (b *Block) Unmarshal(data []byte) error {
	return bpl.Unmarshal(data, b)
}
