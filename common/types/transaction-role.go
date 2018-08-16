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
Package types is definition of common type
 * file description:  transaction
 * @Author: Gong Zibin
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */
package types

import (
	"crypto/sha256"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
	"github.com/bottos-project/bottos/bpl"
	"encoding/hex"
	"github.com/bottos-project/crypto-go/crypto"
	"github.com/ontio/ontology/common/log"
)

// Transaction define transaction struct for bottos protocol
type Transaction struct {
	Version     uint32
	CursorNum   uint32
	CursorLabel uint32
	Lifetime    uint64
	Sender      string // max length 21
	Contract    string // max length 21
	Method      string // max length 21
	Param       []byte
	SigAlg      uint32
	Signature   []byte
}

// HandledTransaction define transaction which is handled
type HandledTransaction struct {
	Transaction *Transaction
	DerivedTrx  []*DerivedTransaction
}

// DerivedTransaction define transaction which is derived from raw transaction
type DerivedTransaction struct {
	Transaction *Transaction
	DerivedTrx  []*DerivedTransaction
}

// Hash transaction hash
func (trx *Transaction) Hash() common.Hash {
	data, _ := bpl.Marshal(trx)
	temp := sha256.Sum256(data)
	hash := sha256.Sum256(temp[:])
	return hash
}

// BasicTransaction define transaction struct for transaction signature
type BasicTransaction struct {
	Version     uint32
	CursorNum   uint32
	CursorLabel uint32
	Lifetime    uint64
	Sender      string
	Contract    string
	Method      string
	Param       []byte
	SigAlg      uint32
}

// VerifySignature verify signature
func (trx *Transaction) VerifySignature(pubkey []byte) bool {
	data, err := bpl.Marshal(BasicTransaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       trx.Param,
		SigAlg:      trx.SigAlg,
	})

	if nil != err {
		return false
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(data)))
	h.Write([]byte(config.Param.ChainId))
	hash := h.Sum(nil)

	ok := crypto.VerifySign(pubkey, hash, trx.Signature)

	if false == ok {
		log.Errorf("trx %x verify signature failed, sender %s, pubkey %x", trx.Hash(), trx.Sender, pubkey)
	}

	return ok
}

// Sign sign a transaction with privkey
func (trx *Transaction) Sign(param []byte, privkey []byte) ([]byte, error) {
	data, err := bpl.Marshal(BasicTransaction{
		Version:     trx.Version,
		CursorNum:   trx.CursorNum,
		CursorLabel: trx.CursorLabel,
		Lifetime:    trx.Lifetime,
		Sender:      trx.Sender,
		Contract:    trx.Contract,
		Method:      trx.Method,
		Param:       param,
		SigAlg:      trx.SigAlg,
	})
	if nil != err {
		return []byte{}, err
	}

	h := sha256.New()
	h.Write([]byte(hex.EncodeToString(data)))
	h.Write([]byte(config.Param.ChainId))
	hash := h.Sum(nil)
	signdata, err := crypto.Sign(hash, privkey)

	return signdata, err
}
