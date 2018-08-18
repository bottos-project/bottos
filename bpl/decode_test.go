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
 * file description:  bpl decode test
 * @Author: Gong Zibin
 * @Date:   2018-08-06
 * @Last Modified by:
 * @Last Modified time:
 */
package bpl

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

type decTest struct {
	input string
	ptr   interface{}
	value interface{}
	error string
}

var decTests = []decTest{
	// booleans
	{input: "C3", ptr: new(bool), value: true},
	{input: "C2", ptr: new(bool), value: false},
	{input: "C4", ptr: new(bool), error: "bpl decode: unknown type identifier C4"},

	// integers
	{input: "CC00", ptr: new(uint8), value: uint8(0)},
	{input: "CC7F", ptr: new(uint8), value: uint8(0x7F)},
	{input: "CCFF", ptr: new(uint8), value: uint8(0xFF)},
	{input: "CD0000", ptr: new(uint16), value: uint16(0)},
	{input: "CD7FFF", ptr: new(uint16), value: uint16(0x7FFF)},
	{input: "CDFFFF", ptr: new(uint16), value: uint16(0xFFFF)},
	{input: "CE00000000", ptr: new(uint32), value: uint32(0)},
	{input: "CE00007FFF", ptr: new(uint32), value: uint32(0x7FFF)},
	{input: "CE7FFFFFFF", ptr: new(uint32), value: uint32(0x7FFFFFFF)},
	{input: "CEFFFFFFFF", ptr: new(uint32), value: uint32(0xFFFFFFFF)},
	{input: "CF0000000000000000", ptr: new(uint64), value: uint64(0)},
	{input: "CF0000000000007FFF", ptr: new(uint64), value: uint64(0x7FFF)},
	{input: "CF000000007FFFFFFF", ptr: new(uint64), value: uint64(0x7FFFFFFF)},
	{input: "CF7FFFFFFFFFFFFFFF", ptr: new(uint64), value: uint64(0x7FFFFFFFFFFFFFFF)},
	{input: "CFFFFFFFFFFFFFFFFF", ptr: new(uint64), value: uint64(0xFFFFFFFFFFFFFFFF)},

	// string
	{
		input: "DA0041496620796F752073686564207465617273207768656E20796F75206D697373207468652073756E2C20796F7520616C736F206D697373207468652073746172732E",
		ptr:   new(string),
		value: string("If you shed tears when you miss the sun, you also miss the stars."),
	},

	// big.int
	{input: "C8000001", ptr: new(*big.Int), value: big.NewInt(0)},
	{input: "C800010101", ptr: new(*big.Int), value: big.NewInt(1)},
	{input: "C80001017F", ptr: new(*big.Int), value: big.NewInt(0x7F)},
	{input: "C80008017FFFFFFFFFFFFFFF", ptr: new(*big.Int), value: big.NewInt(0x7FFFFFFFFFFFFFFF)},
	{
		input: "C8000F010102030405060708090A0B0C0D0E0F",
		ptr:   new(*big.Int),
		value: big.NewInt(0).SetBytes(fromHex("0102030405060708090A0B0C0D0E0F")),
	},

	// byte slice
	{input: "C50000", ptr: new([]byte), value: []byte{}},
	{input: "C50003010203", ptr: new([]byte), value: []byte{1, 2, 3}},

	// byte array
	{input: "C50000", ptr: new([0]byte), value: [0]byte{}},
	{input: "C50003010203", ptr: new([3]byte), value: [3]byte{1, 2, 3}},

	// slices
	{input: "DC0003CD0001CD0002CD0003", ptr: new([]uint16), value: []uint16{1, 2, 3}},
	{
		input: "DC0003DA0005616C696365DA0003626F62DA000563696E6479",
		ptr:   new([]string),
		value: []string{"alice", "bob", "cindy"},
	},

	// struct
	{
		input: "DC0002CD0000DA0000",
		ptr:   new(testStruct),
		value: testStruct{},
	},
	{
		input: "DC0002CD007FDA0006626F74746F73",
		ptr:   new(testStruct),
		value: testStruct{V1: 0x7F, V2: "bottos"},
	},
	{
		input: "DC0002CD0005C0",
		ptr:   new(nestStruct),
		value: nestStruct{V1: 5, V2: nil},
	},
	{
		input: "DC0002CD0001DC0002CD0002DC0002CD0003C0",
		ptr:   new(nestStruct),
		value: nestStruct{1, &nestStruct{2, &nestStruct{3, nil}}},
	},
}

func runTests(t *testing.T, decode func([]byte, interface{}) error) {
	for i, test := range decTests {
		input, err := hex.DecodeString(test.input)
		if err != nil {
			t.Errorf("test %d: invalid hex input %q", i, test.input)
			continue
		}
		err = decode(input, test.ptr)
		if err != nil && test.error == "" {
			t.Errorf("test %d: unexpected Decode error: %v\ndecoding into %T\ninput %q",
				i, err, test.ptr, test.input)
			continue
		}
		if test.error != "" && fmt.Sprint(err) != test.error {
			t.Errorf("test %d: Decode error mismatch\ngot  %v\nwant %v\ndecoding into %T\ninput %q",
				i, err, test.error, test.ptr, test.input)
			continue
		}
		deref := reflect.ValueOf(test.ptr).Elem().Interface()
		if err == nil && !reflect.DeepEqual(deref, test.value) {
			t.Errorf("test %d: value mismatch\ngot  %#v\nwant %#v\ndecoding into %T\ninput %q",
				i, deref, test.value, test.ptr, test.input)
		}
	}
}

func TestDecodeWithByteReader(t *testing.T) {
	runTests(t, func(input []byte, into interface{}) error {
		return Decode(bytes.NewReader(input), into)
	})
}
