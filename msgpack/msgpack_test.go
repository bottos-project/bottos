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
	"fmt"
	//"bytes"
	"encoding/hex"
	"testing"
)

func BytesToHex(d []byte) string {
	return hex.EncodeToString(d)
}

func HexToBytes(str string) ([]byte, error) {
	h, err := hex.DecodeString(str)

	return h, err
}

func TestMarshalStruct(t *testing.T) {
	type TestStruct struct {
		V1 string
		V2 uint8
		V3 uint16
		V4 uint32
		V5 uint64
		V6 []byte
	}

	ts := TestStruct{
		V1: "testuser",
		V2: 99,
		V3: 999,
		V4: 9999,
		V5: 99999,
		V6: []byte{0xac, 0xcd, 0xde},
	}

	b, err := Marshal(ts)

	fmt.Printf("%v\n", BytesToHex(b))
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
}
