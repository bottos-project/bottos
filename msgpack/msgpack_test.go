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
 * file description:  msgpack go
 * @Author: Gong Zibin
 * @Date:   2018-08-02
 * @Last Modified by:
 * @Last Modified time:
 */
package msgpack

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/bottos-project/bottos/common"
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

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Printf("ts1: %#v \n", ts1)
	fmt.Println(err)
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
	fmt.Println("TestMarshalNestStruct...")

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: TestSubStruct{V1: "123", V2: 3},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Printf("ts1: %#v \n", ts1)
	fmt.Println(err)
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
	fmt.Println("TestMarshalNestStructPtr...")

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: &TestSubStruct{V1: "123", V2: 3},
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Printf("ts1: %#v, %#v\n", ts1, *ts1.V3)
	fmt.Println(err)
}

func TestMarshalNilPtr(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 *uint32
		V3 uint64
	}

	fmt.Println("TestMarshalNilPtr...")

	ts := TestStruct{
		V1: "testuser",
		V2: nil,
		V3: 999999,
	}
	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Printf("ts1: %#v\n", ts1)
	fmt.Println(err)
}

func TestMarshalCustomHashType(t *testing.T) {
	type Hash [32]byte

	type Account struct {
		AccountName string
		CodeVersion Hash
	}

	fmt.Println("TestMarshalCustomHashType...")

	ts := Account{
		AccountName: "testuser",
		CodeVersion: sha256.Sum256([]byte("testuser")),
	}
	b, err := Marshal(ts)

	fmt.Printf("%x\n", b)
	fmt.Println(err)

	ts1 := Account{}
	err = Unmarshal(b, &ts1)
	fmt.Printf("ts1: %#x\n", ts1)
	fmt.Println(err)
}

func TestMarshalBigInt(t *testing.T) {
	type TestStruct struct {
		V1 uint32
		V2 *big.Int
	}

	fmt.Println("TestMarshalBigInt...")

	ts := TestStruct{
		V1: 9999,
		V2: new(big.Int).SetUint64(uint64(999999)),
	}
	b, err := Marshal(ts)

	fmt.Printf("%x\n", b)
	fmt.Println(err)

	ts1 := TestStruct{}
	err = Unmarshal(b, &ts1)
	fmt.Printf("ts1: %#v\n", ts1)
	fmt.Println(err)
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

	fmt.Println("TestMarshalTransaction...")

	ts := Transaction{
		Version:     1,
		CursorNum:   999,
		CursorLabel: 86868797,
		Lifetime:    uint64(time.Now().Unix()),
		Sender:      "alice",
		Contract:    "bottos",
		Method:      "transfer",
		Param:       HexToBytes("dc000212345678"),
		SigAlg:      1,
		Signature:   []byte{},
	}
	b, err := Marshal(ts)

	fmt.Printf("%x\n", b)
	fmt.Println(err)

	ts1 := Transaction{}
	err = Unmarshal(b, &ts1)
	fmt.Printf("ts1: %#x\n", ts)
	fmt.Printf("ts1: %#x\n", ts1)
	fmt.Println(err)
}

func TestMarshalBalance(t *testing.T) {
	type Balance struct {
		AccountName string
		Balance     *big.Int
	}

	fmt.Println("TestMarshalTransaction...")

	balance, _ := new(big.Int).SetString("100000000001000000000", 10)
	ts := Balance{
		AccountName: "alice",
		Balance:     balance,
	}
	b, err := Marshal(ts)

	fmt.Printf("%x\n", b)
	fmt.Println(err)

	ts1 := Balance{}
	err = Unmarshal(b, &ts1)
	fmt.Printf("ts1: %#v\n", ts)
	fmt.Printf("ts1: %#v\n", ts1)
	fmt.Println(err)
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

func TestMarshalBlock(t *testing.T) {

	type Header struct {
		Version         uint32
		PrevBlockHash   common.Hash
		Number          uint32
		Timestamp       uint64
		MerkleRoot      common.Hash
		Delegate        Name
		DelegateSign    []byte
		DelegateChanges []Name
	}

	type Transaction struct {
		Version     uint32
		CursorNum   uint32
		CursorLabel uint32
		Lifetime    uint64
		Sender      Name
		Contract    Name
		Method      Name
		Param       []byte
		SigAlg      uint32
		Signature   []byte
	}

	type Block struct {
		Header       *Header
		Transactions []*Transaction
	}

	fmt.Println("TestMarshalBlock...")

	header := Header{
		Version:         1,
		PrevBlockHash:   common.Sha256([]byte("123")),
		Number:          123,
		Timestamp:       uint64(time.Now().Unix()),
		MerkleRoot:      common.Sha256([]byte("234")),
		Delegate:        StringToName("toliman"),
		DelegateSign:    []byte{},
		DelegateChanges: []Name{StringToName("toliman"), StringToName("ran")},
	}
	tx1 := Transaction{
		Version:     1,
		CursorNum:   999,
		CursorLabel: 86868797,
		Lifetime:    uint64(time.Now().Unix()),
		Sender:      StringToName("alice"),
		Contract:    StringToName("bottos"),
		Method:      StringToName("transfer"),
		Param:       HexToBytes("dc000212345678"),
		SigAlg:      1,
		Signature:   []byte{},
	}
	tx2 := Transaction{
		Version:     1,
		CursorNum:   999,
		CursorLabel: 1412312421,
		Lifetime:    uint64(time.Now().Unix()),
		Sender:      StringToName("alice"),
		Contract:    StringToName("bottos"),
		Method:      StringToName("transfer"),
		Param:       HexToBytes("dc000212345678"),
		SigAlg:      1,
		Signature:   []byte{},
	}

	block := Block{Header: &header}
	block.Transactions = append(block.Transactions, &tx1)
	block.Transactions = append(block.Transactions, &tx2)

	fmt.Printf("block: %#v\n", block)

	b, err := Marshal(block)

	fmt.Printf("%x\n", b)
	fmt.Println(err)

	block1 := Block{}
	err = Unmarshal(b, &block1)
	fmt.Printf("block: %#x\n", block.Header)
	fmt.Printf("block1: %#x\n", block1.Header)
	fmt.Println(err)
}
