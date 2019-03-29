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
 * file description:  bpl test
 * @Author: Gong Zibin
 * @Date:   2018-08-02
 * @Last Modified by:
 * @Last Modified time:
 */
package bpl

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/bottos-project/bottos/common"
	"github.com/stretchr/testify/assert"
	)

func BytesToHex(d []byte) string {
	return hex.EncodeToString(d)
}

func HexToBytes(str string) []byte {
	h, _ := hex.DecodeString(str)

	return h
}

func TestMarshalStruct(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 uint8
		V3 uint16
		V4 uint32
		V5 uint64
		V6 []byte
		V7 bool
	}

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: 999,
		V4: 9999,
		V5: 99999,
		V6: []byte{0xac, 0xcd, 0xde},
		V7: true,
	}

	b, err := Marshal(ts)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc0007da00087465737475736572cc63cd03e7ce0000270fcf000000000001869fc50003accddec3"))

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	assert.Equal(t, ts, ts1)
}


// The value of an uninitialized slice is nil.
// ref. to https://golang.org/ref/spec#Slice_types
func TestMarshalNilSlice(t *testing.T) {

	type TestStruct struct {
		V1 uint32
		V2 []byte
	}

	ts := TestStruct{V1:1}
	assert.Nil(t, ts.V2)

	ts = TestStruct{V1:1, V2:[]byte{}}
	assert.NotNil(t, ts.V2)
	b, err := Marshal(ts)
	assert.Nil(t, err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	assert.NotNil(t, ts1.V2)
	assert.Equal(t, ts, ts1)
}

func TestMarshalNestStruct(t *testing.T) {
	type TestSubStruct struct {
		V1 string
		V2 uint32
	}

	type TestStruct struct {
		V1 string
		V2 uint32
		V3 TestSubStruct
	}

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: TestSubStruct{V1: "123", V2: 3},
	}
	b, err := Marshal(ts)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc0003da00087465737475736572ce00000063dc0002da0003313233ce00000003"))

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	assert.Equal(t, ts, ts1)
}

func TestMarshalNestStructPtr(t *testing.T) {
	type TestSubStruct struct {
		V1 string
		V2 uint32
	}

	type TestStruct struct {
		V1 string
		V2 uint32
		V3 *TestSubStruct
	}

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: &TestSubStruct{V1: "123", V2: 3},
	}
	b, err := Marshal(ts)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc0003da00087465737475736572ce00000063dc0002da0003313233ce00000003"))

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	assert.Equal(t, ts, ts1)
}

func TestMarshalNilPtr(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 *uint32
		V3 uint64
	}

	ts := TestStruct{
		V1: "testuser",
		V2: nil,
		V3: 999999,
	}
	b, err := Marshal(ts)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc0003da00087465737475736572c0cf00000000000f423f"))

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	assert.Equal(t, ts, ts1)
}

func TestMarshalCustomHashType(t *testing.T) {
	type Hash [32]byte

	type TestStruct struct {
		AccountName string
		CodeVersion Hash
	}

	ts := TestStruct{
		AccountName: "testuser",
		CodeVersion: sha256.Sum256([]byte("testuser")),
	}
	b, err := Marshal(ts)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc0002da00087465737475736572c50020ae5deb822e0d71992900471a7199d0d95b8e7c9d05c40a8245a281fd2c1d6684"))

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	assert.Equal(t, ts, ts1)
}

func TestMarshalBigInt(t *testing.T) {
	type TestStruct struct {
		V1 uint32
		V2 *big.Int
	}

	ts := TestStruct{
		V1: 9999,
		V2: new(big.Int).SetUint64(uint64(999999)),
	}
	b, err := Marshal(ts)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc0002ce0000270fc80003010f423f"))
	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	assert.Equal(t, ts, ts1)
}

func TestMarshalTransaction(t *testing.T) {
	type Transaction struct {
		Version     uint32
		CursorNum   uint32
		CursorLabel uint32
		Lifetime    uint64
		Sender      string
		Contract    string
		Method      string
		Param       []byte
		SigAlg      uint32
		Signature   []byte
	}

	ts := Transaction{
		Version:     1,
		CursorNum:   999,
		CursorLabel: 86868797,
		Lifetime:    0x5b691451,
		Sender:      "alice",
		Contract:    "bottos",
		Method:      "transfer",
		Param:       HexToBytes("dc000212345678"),
		SigAlg:      1,
		Signature:   []byte{},
	}
	b, err := Marshal(ts)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc000ace00000001ce000003e7ce052d833dcf000000005b691451da0005616c696365da0006626f74746f73da00087472616e73666572c50007dc000212345678ce00000001c50000"))
	ts1 := Transaction{}
	err = Unmarshal(b, &ts1)
	assert.Equal(t, ts, ts1)
}

func TestMarshalBalance(t *testing.T) {
	type Balance struct {
		AccountName string
		Balance     *big.Int
	}

	balance, _ := new(big.Int).SetString("100000000001000000000", 10)
	ts := Balance{
		AccountName: "alice",
		Balance:     balance,
	}
	b, err := Marshal(ts)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc0002da0005616c696365c8000901056bc75e2d9eaaca00"))

	ts1 := Balance{}
	err = Unmarshal(b, &ts1)
	assert.Equal(t, ts, ts1)
}

type Name [16]byte // uint128
func StringToName(s string) Name {
	var name Name
	name.SetBytes([]byte(s))
	return name
}
func (h *Name) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-16:]
	}
	copy(h[16-len(b):], b)
}

func NewName(s string) common.Name {
	name, err := common.NewName(s)
	if err != nil {
		panic(err)
	}
	return name
}

func TestMarshalBlock(t *testing.T) {

	type Header struct {
		Version         uint32
		PrevBlockHash   common.Hash
		Number          uint32
		Timestamp       uint64
		MerkleRoot      common.Hash
		Delegate        common.Name
		DelegateSign    []byte
		DelegateChanges []common.Name
	}

	type Transaction struct {
		Version     uint32
		CursorNum   uint32
		CursorLabel uint32
		Lifetime    uint64
		Sender      common.Name
		Contract    common.Name
		Method      common.Name
		Param       []byte
		SigAlg      uint32
		Signature   []byte
	}

	type Block struct {
		Header       *Header
		Transactions []*Transaction
	}

	header := Header{
		Version:         1,
		PrevBlockHash:   common.Sha256([]byte("123")),
		Number:          123,
		Timestamp:       0x5b691451,
		MerkleRoot:      common.Sha256([]byte("234")),
		Delegate:        NewName("sirus"),
		DelegateSign:    []byte{},
		DelegateChanges: []common.Name{NewName("sirus"), NewName("ran")},
	}

	tx1 := Transaction{
		Version:     1,
		CursorNum:   999,
		CursorLabel: 86868797,
		Lifetime:    0x5b691451,
		Sender:      NewName("bob"),
		Contract:    NewName("bottos"),
		Method:      NewName("transfer"),
		Param:       HexToBytes("dc000212345678"),
		SigAlg:      1,
		Signature:   []byte{},
	}
	tx2 := Transaction{
		Version:     1,
		CursorNum:   999,
		CursorLabel: 1412312421,
		Lifetime:    0x5b691451,
		Sender:      NewName("bob"),
		Contract:    NewName("bottos"),
		Method:      NewName("transfer"),
		Param:       HexToBytes("dc000212345678"),
		SigAlg:      1,
		Signature:   []byte{},
	}

	block := Block{Header: &header}
	block.Transactions = make([]*Transaction, 2)
	block.Transactions = append(block.Transactions, &tx1)
	block.Transactions = append(block.Transactions, &tx2)

	b, err := Marshal(block)
	assert.Nil(t, err)
	assert.Equal(t, b, fromHex("dc0002dc0008ce00000001c50020a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3ce0000007bcf000000005b691451c50020114bd151f8fb0c58642d2170da4ae7d7c57977260ac2cc8905306cab6b2acabcc500100000000000000000000000001c49b79cc50000dc0002c500100000000000000000000000001c49b79cc500100000000000000000000000000001b297dc0004c0c0dc000ace00000001ce000003e7ce052d833dcf000000005b691451c500100000000000000000000000000000b60bc50010000000000000000000000002d875d61cc500100000000000000000000075b29770f39bc50007dc000212345678ce00000001c50000dc000ace00000001ce000003e7ce542e2d65cf000000005b691451c500100000000000000000000000000000b60bc50010000000000000000000000002d875d61cc500100000000000000000000075b29770f39bc50007dc000212345678ce00000001c50000"))
	block1 := Block{}
	err = Unmarshal(b, &block1)
	assert.Equal(t, block, block1)
}

type bTestStruct struct {
	V1 uint32
	V2 string
	V3 *big.Int
}

var val bTestStruct = bTestStruct{
	V1: uint32(999),
	V2: "bottos",
	V3: big.NewInt(0x7FFFFFFFFFFFFFFF),
}

func bplEncode() ([]byte, error) {
	b := new(bytes.Buffer)
	err := Encode(val, b)
	return b.Bytes(), err
}

func jsonEncode() ([]byte, error) {
	b, err := json.Marshal(val)
	return b, err
}

func BenchmarkBpl(b *testing.B) {
	for n := 0; n < b.N; n++ {
		bplEncode()
	}
}

func BenchmarkJson(b *testing.B) {
	for n := 0; n < b.N; n++ {
		jsonEncode()
	}
}
