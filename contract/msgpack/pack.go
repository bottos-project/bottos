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
 * file description:  contract db
 * @Author: Gong Zibin
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */

package msgpack

import (
	"io"
)

const (
	//BIN16 is byte array type identifier
	BIN16 = 0xc5
	//UINT8 is uint8
	UINT8 = 0xcc
	//UINT16 is uint16
	UINT16 = 0xcd
	//UINT32 is uint32
	UINT32 = 0xce
	//UINT64 is uint64
	UINT64 = 0xcf
	//STR16 is string type identifier
	STR16 = 0xda
	//ARRAY16 is array size type identifier
	ARRAY16 = 0xdc
	//LEN_INT32 value
	LEN_INT32 = 4
	//LEN_INT64 value
	LEN_INT64 = 8
	//MAX16BIT value
	MAX16BIT = 2 << (16 - 1)
	//REGULAR_UINT7_MAX value
	REGULAR_UINT7_MAX = 2 << (7 - 1)
	//REGULAR_UINT8_MAX value
	REGULAR_UINT8_MAX = 2 << (8 - 1)
	//REGULAR_UINT16_MAX value
	REGULAR_UINT16_MAX = 2 << (16 - 1)
	//REGULAR_UINT32_MAX value
	REGULAR_UINT32_MAX = 2 << (32 - 1)

	//SPECIAL_INT8 value
	SPECIAL_INT8 = 32
	//SPECIAL_INT16 value
	SPECIAL_INT16 = 2 << (8 - 2)
	//SPECIAL_INT32 value
	SPECIAL_INT32 = 2 << (16 - 2)
	//SPECIAL_INT64 value
	SPECIAL_INT64 = 2 << (32 - 2)
)

//Bytes is []byte type
type Bytes []byte

//PackUint8 is to pack a given value and writes it into the specified writer.
func PackUint8(writer io.Writer, value uint8) (n int, err error) {
	return writer.Write(Bytes{UINT8, value})
}

//PackUint16 is to pack a given value and writes it into the specified writer.
func PackUint16(writer io.Writer, value uint16) (n int, err error) {
	return writer.Write(Bytes{UINT16, byte(value >> 8), byte(value)})
}

//PackUint32 is to pack a given value and writes it into the specified writer.
func PackUint32(writer io.Writer, value uint32) (n int, err error) {
	return writer.Write(Bytes{UINT32, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

//PackUint64 is to pack a given value and writes it into the specified writer.
func PackUint64(writer io.Writer, value uint64) (n int, err error) {
	return writer.Write(Bytes{UINT64, byte(value >> 56), byte(value >> 48), byte(value >> 40), byte(value >> 32), byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

//PackBin16 is to pack a given value and writes it into the specified writer.
func PackBin16(writer io.Writer, value []byte) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{BIN16, byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write(value)
	return n1 + n2, err
}

//PackStr16 is to pack a given value and writes it into the specified writer.
func PackStr16(writer io.Writer, value string) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{STR16, byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write([]byte(value))
	return n1 + n2, err
}

//PackArraySize is to pack a given value and writes it into the specified writer.
func PackArraySize(writer io.Writer, length uint16) (n int, err error) {
	n, err = writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
	if err != nil {
		return n, err
	}
	return n, nil
}
