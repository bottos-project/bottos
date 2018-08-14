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
 * file description:  bpl pack
 * @Author: Gong Zibin
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */

package bpl

import (
	"io"
)

const (
	// NIL nil or null type identifier
	NIL = 0xc0
	// FALSE is a Bool type identifier
	FALSE = 0xc2
	// TRUE is a Bool type identifier
	TRUE = 0xc3
	//BIN16 is byte array type identifier
	BIN16 = 0xc5
	//EXT16 is extension type identifier
	EXT16 = 0xc8
	//UINT8 is uint8 type identifier
	UINT8 = 0xcc
	//UINT16 is uint16 type identifier
	UINT16 = 0xcd
	//UINT32 is uint32 type identifier
	UINT32 = 0xce
	//UINT64 is uint64 type identifier
	UINT64 = 0xcf
	//STR16 is string type identifier
	STR16 = 0xda
	//ARRAY16 is array size type identifier
	ARRAY16 = 0xdc
)

//Bytes is []byte type
type Bytes []byte

// PackNil pack a nil type
func PackNil(writer io.Writer) (n int, err error) {
	return writer.Write(Bytes{NIL})
}

// PackBool pack a bool value.
func PackBool(writer io.Writer, value bool) (n int, err error) {
	if value {
		return writer.Write(Bytes{TRUE})
	}
	return writer.Write(Bytes{FALSE})
}

//PackUint8 is to pack a uint8.
func PackUint8(writer io.Writer, value uint8) (n int, err error) {
	return writer.Write(Bytes{UINT8, value})
}

//PackUint16 is to pack a uint16
func PackUint16(writer io.Writer, value uint16) (n int, err error) {
	return writer.Write(Bytes{UINT16, byte(value >> 8), byte(value)})
}

//PackUint32 is to pack a uint32.
func PackUint32(writer io.Writer, value uint32) (n int, err error) {
	return writer.Write(Bytes{UINT32, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

//PackUint64 is to pack a uint64.
func PackUint64(writer io.Writer, value uint64) (n int, err error) {
	return writer.Write(Bytes{UINT64, byte(value >> 56), byte(value >> 48), byte(value >> 40), byte(value >> 32), byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
}

//PackBin16 is to pack a byte array.
func PackBin16(writer io.Writer, value []byte) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{BIN16, byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write(value)
	return n1 + n2, err
}

//PackStr16 is to pack a string.
func PackStr16(writer io.Writer, value string) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{STR16, byte(length >> 8), byte(length)})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write([]byte(value))
	return n1 + n2, err
}

//PackExt16 is to pack a extension type.
func PackExt16(writer io.Writer, t byte, value []byte) (n int, err error) {
	length := len(value)
	n1, err := writer.Write(Bytes{EXT16, byte(length >> 8), byte(length), t})
	if err != nil {
		return n1, err
	}
	n2, err := writer.Write([]byte(value))
	return n1 + n2, err
}

//PackArraySize is to pack a array size header.
func PackArraySize(writer io.Writer, length uint16) (n int, err error) {
	n, err = writer.Write(Bytes{ARRAY16, byte(length >> 8), byte(length)})
	if err != nil {
		return n, err
	}
	return n, nil
}
