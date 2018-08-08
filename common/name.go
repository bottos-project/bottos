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
 * file description:  Base58 encoded Name
 * @Author: Gong Zibin
 * @Date:   2018-08-08
 * @Last Modified by:
 * @Last Modified time:
 */

package common

import (
	"fmt"
	"math/big"
)

const (
	// MAX_NAME_BYTE is max byte length of Name type
	MAX_NAME_BYTE = 16
	// ENCODE_RADIX base58
	ENCODE_RADIX = 58
	// ENCODE_BIT_LEN base58
	ENCODE_BIT_LEN = 6
)

// Name basic type for account name, method and contract
type Name [MAX_NAME_BYTE]byte

var defaultEncoding = encoding([]byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"))

// NewName encode a string name to Name type
func NewName(s string) (Name, error) {
	encoded, err := defaultEncoding.encode([]byte(s))
	if err != nil {
		return Name{}, err
	}
	return encoded, nil
}

// ToString decode Name type to string type
func (n Name) ToString() string {
	decoded, err := defaultEncoding.decode(n)
	if err != nil {
		return ""
	}
	return string(decoded)
}

func (n Name) toBig() *big.Int {
	return big.NewInt(0).SetBytes(n[:])
}

// Bytes get bytes of the name
func (n Name) Bytes() []byte {
	return n[:]
}

// EncodingStruct is a radix 58 encoding/decoding scheme.
type EncodingStruct struct {
	alphabet  [ENCODE_RADIX]byte
	decodeMap map[byte]int64
}

func encoding(alphabet []byte) *EncodingStruct {
	enc := &EncodingStruct{}
	copy(enc.alphabet[:], alphabet[:])
	for i := range enc.decodeMap {
		enc.decodeMap[i] = -1
	}
	enc.decodeMap = make(map[byte]int64)
	for i, b := range enc.alphabet {
		enc.decodeMap[b] = int64(i)
	}
	return enc
}

// string name -> Name
func (encoding *EncodingStruct) encode(src []byte) (Name, error) {
	if len(src) == 0 {
		return Name{}, nil
	}

	bigname := big.NewInt(0)
	for _, c := range src {
		if idx, ok := encoding.decodeMap[c]; ok {
			bigname.Lsh(bigname, ENCODE_BIT_LEN)
			bigname.Add(bigname, big.NewInt(idx))
		} else {
			return Name{}, fmt.Errorf("invalid character '%c' in decoding a base58 string \"%s\"", c, src)
		}
	}

	name := Name{}
	name.setBytes(bigname.Bytes())
	return name, nil
}

// Name -> string name
func (encoding *EncodingStruct) decode(name Name) ([]byte, error) {
	bigname := name.toBig()

	var decoded []byte
	zero := big.NewInt(0)
	for {
		switch bigname.Cmp(zero) {
		case 1:
			val := bigname.Int64() & 0x3F
			if val >= ENCODE_RADIX {
				return []byte{}, fmt.Errorf("invalid encoded value %v", val)
			}
			decoded = append(decoded, encoding.alphabet[val])
			bigname.Rsh(bigname, ENCODE_BIT_LEN)
		case 0:
			reverse(decoded)
			return decoded, nil
		default:
			return nil, fmt.Errorf("expecting a positive number in base58 encoding but got %q", bigname)
		}
	}
}

func (n *Name) setBytes(b []byte) {
	if len(b) > len(n) {
		b = b[len(b)-MAX_NAME_BYTE:]
	}
	copy(n[MAX_NAME_BYTE-len(b):], b)
}
