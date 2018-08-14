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
 * file description:  bpl encode test
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
	"strings"
	"testing"
)

type testStruct struct {
	V1 uint16
	V2 string
}

type nestStruct struct {
	V1 uint16
	V2 *nestStruct
}

type encTest struct {
	val           interface{}
	output, error string
}

var tests string = string("If you shed tears when you miss the sun, you also miss the stars.")

type NamedByteArray [4]byte

var encTests = []encTest{
	// booleans
	{val: true, output: "C3"},
	{val: false, output: "C2"},

	// integers
	{val: uint8(0), output: "CC00"},
	{val: uint8(127), output: "CC7F"},
	{val: uint8(0xFF), output: "CCFF"},
	{val: uint16(0), output: "CD0000"},
	{val: uint16(0x7FFF), output: "CD7FFF"},
	{val: uint16(0xFFFF), output: "CDFFFF"},
	{val: uint32(0), output: "CE00000000"},
	{val: uint32(0x7FFF), output: "CE00007FFF"},
	{val: uint32(0x7FFFFFFF), output: "CE7FFFFFFF"},
	{val: uint32(0xFFFFFFFF), output: "CEFFFFFFFF"},
	{val: uint64(0), output: "CF0000000000000000"},
	{val: uint64(0x7FFF), output: "CF0000000000007FFF"},
	{val: uint64(0x7FFFFFFF), output: "CF000000007FFFFFFF"},
	{val: uint64(0x7FFFFFFFFFFFFFFF), output: "CF7FFFFFFFFFFFFFFF"},
	{val: uint64(0xFFFFFFFFFFFFFFFF), output: "CFFFFFFFFFFFFFFFFF"},

	// string
	{
		val:    string("If you shed tears when you miss the sun, you also miss the stars."),
		output: "DA0041496620796F752073686564207465617273207768656E20796F75206D697373207468652073756E2C20796F7520616C736F206D697373207468652073746172732E",
	},
	// big.int
	{val: big.NewInt(0), output: "C8000001"},
	{val: big.NewInt(1), output: "C800010101"},
	{val: big.NewInt(0x7F), output: "C80001017F"},
	{val: big.NewInt(0x7FFFFFFFFFFFFFFF), output: "C80008017FFFFFFFFFFFFFFF"},
	{
		val:    big.NewInt(0).SetBytes(fromHex("0102030405060708090A0B0C0D0E0F")),
		output: "C8000F010102030405060708090A0B0C0D0E0F",
	},
	{val: *big.NewInt(0), output: "C8000001"},
	{val: *big.NewInt(0x7F), output: "C80001017F"},

	// bytes
	{val: []byte{}, output: "C50000"},
	{val: []byte{1, 2, 3}, output: "C50003010203"},

	// slices
	{val: []uint16{1, 2, 3}, output: "DC0003CD0001CD0002CD0003"},
	{
		val:    []string{"alice", "bob", "cindy"},
		output: "DC0003DA0005616C696365DA0003626F62DA000563696E6479",
	},

	// struct
	{val: testStruct{}, output: "DC0002CD0000DA0000"},
	{val: testStruct{V1: 0x7F, V2: "bottos"}, output: "DC0002CD007FDA0006626F74746F73"},
	// struct with nil
	{val: &nestStruct{V1: 5, V2: nil}, output: "DC0002CD0005C0"},
	{val: &nestStruct{1, &nestStruct{2, &nestStruct{3, nil}}}, output: "DC0002CD0001DC0002CD0002DC0002CD0003C0"},

	// named byte array type
	{val: NamedByteArray{0x1, 0x2, 0x3, 0x4}, output: "C5000401020304"},
}

func runEncTests(t *testing.T, f func(val interface{}) ([]byte, error)) {
	for i, test := range encTests {
		output, err := f(test.val)
		if err != nil && test.error == "" {
			t.Errorf("test %d: unexpected error: %v\nvalue %#v\ntype %T",
				i, err, test.val, test.val)
			continue
		}
		if test.error != "" && fmt.Sprint(err) != test.error {
			t.Errorf("test %d: error mismatch\ngot   %v\nwant  %v\nvalue %#v\ntype  %T",
				i, err, test.error, test.val, test.val)
			continue
		}
		if err == nil && !bytes.Equal(output, fromHex(test.output)) {
			t.Errorf("test %d: output mismatch:\ngot   %X\nwant  %s\nvalue %#v\ntype  %T",
				i, output, test.output, test.val, test.val)
		}
	}
}

func TestEncode(t *testing.T) {
	runEncTests(t, func(val interface{}) ([]byte, error) {
		b := new(bytes.Buffer)
		err := Encode(val, b)
		return b.Bytes(), err
	})
}

func fromHex(str string) []byte {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		panic(fmt.Sprintf("invalid hex string: %q", str))
	}
	return b
}

func toHex(b []byte) string {
	return hex.EncodeToString(b)
}
