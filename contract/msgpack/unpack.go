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
 * @Date:   2017-12-20
 * @Last Modified by:
 * @Last Modified time:
 */

package msgpack

import (
	"fmt"
	"io"
)

type (
	//Bytes1 is first
	Bytes1 [1]byte
	//Bytes2 is second
	Bytes2 [2]byte
	//Bytes4 is third
	Bytes4 [4]byte
	//Bytes8 is forth
	Bytes8 [8]byte
)

const (
	//NEGFIXNUM is negfix maxnum
	NEGFIXNUM = 0xe0
	//FIXMAPMAX is fixmap maxnum
	FIXMAPMAX = 0x8f
	//FIXARRAYMAX is fixarray maxnum
	FIXARRAYMAX = 0x9f
	//FIXRAWMAX is fix raw max
	FIXRAWMAX = 0xbf
	//FIRSTBYTEMASK is first byte mask
	FIRSTBYTEMASK = 0xf
)

func readByte(reader io.Reader) (v uint8, err error) {
	var data Bytes1
	_, e := reader.Read(data[0:])
	if e != nil {
		return 0, e
	}
	return data[0], nil
}

//UnpackUint8 is to unpack message
func UnpackUint8(reader io.Reader) (v uint8, err error) {
	c, e := readByte(reader)
	if e == nil && c == UINT8 {
		v, err = readByte(reader)
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func readUint16(reader io.Reader) (v uint16, n int, err error) {
	var data Bytes2
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint16(data[0]) << 8) | uint16(data[1]), n, nil
}

//UnpackUint16 is to unpack message
func UnpackUint16(reader io.Reader) (v uint16, err error) {
	c, e := readByte(reader)
	if e == nil && c == UINT16 {
		v, _, err = readUint16(reader)
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func readUint32(reader io.Reader) (v uint32, n int, err error) {
	var data Bytes4
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3]), n, nil
}

//UnpackUint32 is to unpack message
func UnpackUint32(reader io.Reader) (v uint32, err error) {
	c, e := readByte(reader)
	if e == nil && c == UINT32 {
		v, _, err = readUint32(reader)
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func readUint64(reader io.Reader) (v uint64, n int, err error) {
	var data Bytes8
	n, e := reader.Read(data[0:])
	if e != nil {
		return 0, n, e
	}
	return (uint64(data[0]) << 56) | (uint64(data[1]) << 48) | (uint64(data[2]) << 40) | (uint64(data[3]) << 32) | (uint64(data[4]) << 24) | (uint64(data[5]) << 16) | (uint64(data[6]) << 8) | uint64(data[7]), n, nil
}

//UnpackUint64 is to unpack message
func UnpackUint64(reader io.Reader) (v uint64, err error) {
	c, e := readByte(reader)
	if e == nil && c == UINT64 {
		v, _, err = readUint64(reader)
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

//UnpackArraySize is to unpack message
func UnpackArraySize(reader io.Reader) (size uint16, err error) {
	c, e := readByte(reader)
	if e != nil {
		return 0, e
	}

	header := uint16(c)
	if header != ARRAY16 {
		return 0, fmt.Errorf("Not Array 16")
	}

	size, _, e = readUint16(reader)
	if e != nil {
		return 0, e
	}

	return size, nil
}

//UnpackStr16 is to unpack message
func UnpackStr16(reader io.Reader) (string, error) {
	c, e := readByte(reader)
	if e == nil && c == STR16 {
		size, _, e := readUint16(reader)
		if e == nil {
			value := make([]byte, size)
			n, e := reader.Read(value)
			if e == nil && uint16(n) == size {
				return string(value), nil
			}
		}
	}

	return "", e
}

//UnpackBin16 is to unpack message
func UnpackBin16(reader io.Reader) ([]byte, error) {
	c, e := readByte(reader)
	if e == nil && c == BIN16 {
		size, _, e := readUint16(reader)
		if e == nil {
			value := make([]byte, size)
			n, e := reader.Read(value)
			if e == nil && uint16(n) == size {
				return value, nil
			}
		}
	}

	return []byte{}, e
}
