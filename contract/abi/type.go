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
 * file description:  abi type
 * @Author: Gong Zibin
 * @Date:   2018-08-27
 * @Last Modified by:
 * @Last Modified time:
 */

package abi

import (
	"io"
	"reflect"
	"math/big"
	"github.com/bottos-project/bottos/bpl"
)

const (
	ABITypeUINT8 byte = iota
	ABITypeUINT16
	ABITypeUINT32
	ABITypeUINT64
	ABITypeUINT256
	ABITypeString
	ABITypeBytes
)

var (
	RefTypeUINT8    = reflect.TypeOf(uint8(0))
	RefTypeUINT16    = reflect.TypeOf(uint16(0))
	RefTypeUINT32    = reflect.TypeOf(uint32(0))
	RefTypeUINT64    = reflect.TypeOf(uint64(0))
	RefTypeUINT256   = reflect.TypeOf(&big.Int{})
	RefTypeString    = reflect.TypeOf(string(""))
	//RefTypeBytes = reflect.TypeOf([]byte{})
)

type TypeWriter func(io.Writer, reflect.Value) error

type Type struct {
	StringType string
	AbiType    byte
	RefType    reflect.Type
	RefKind    reflect.Kind
	Writer     TypeWriter
}

type ABIType struct {
	TypeMap map[string]Type
}

func NewABIType() *ABIType {
	at := &ABIType{}
	at.TypeMap = make(map[string]Type)

	at.TypeMap["uint8"] = Type{
		StringType: "uint8",
		AbiType: ABITypeUINT8,
		RefType: RefTypeUINT8,
		RefKind: reflect.Uint8,
		Writer: writeUint8,
	}
	at.TypeMap["uint16"] = Type{
		StringType: "uint16",
		AbiType: ABITypeUINT16,
		RefType: RefTypeUINT16,
		RefKind: reflect.Uint16,
		Writer: writeUint16,
	}
	at.TypeMap["uint32"] = Type{
		StringType: "uint32",
		AbiType: ABITypeUINT32,
		RefType: RefTypeUINT32,
		RefKind: reflect.Uint32,
		Writer: writeUint32,
	}
	at.TypeMap["uint64"] = Type{
		StringType: "uint64",
		AbiType: ABITypeUINT64,
		RefType: RefTypeUINT64,
		RefKind: reflect.Uint64,
		Writer: writeUint64,
	}
	at.TypeMap["uint256"] = Type{
		StringType: "uint256",
		AbiType: ABITypeUINT256,
		RefType: RefTypeUINT256,
		RefKind: reflect.Ptr,
		Writer: writeUint256,
	}
	at.TypeMap["string"] = Type{
		StringType: "string",
		AbiType: ABITypeString,
		RefType: RefTypeString,
		RefKind: reflect.String,
		Writer: writeString,
	}
	at.TypeMap["bytes"] = Type{
		StringType: "bytes",
		AbiType:    ABITypeBytes,
		RefType:    RefTypeUINT8,
		RefKind:    reflect.Slice,
		Writer:     writeBytes,
	}

	return at
}

func writeUint8(w io.Writer, val reflect.Value) error {
	bpl.PackUint8(w, uint8(val.Uint()))
	return nil
}

func writeUint16(w io.Writer, val reflect.Value) error {
	bpl.PackUint16(w, uint16(val.Uint()))
	return nil
}

func writeUint32(w io.Writer, val reflect.Value) error {
	bpl.PackUint32(w, uint32(val.Uint()))
	return nil
}

func writeUint64(w io.Writer, val reflect.Value) error {
	bpl.PackUint64(w, uint64(val.Uint()))
	return nil
}

func writeUint256(w io.Writer, val reflect.Value) error {
	return nil
}

func writeString(w io.Writer, val reflect.Value) error {
	bpl.PackStr16(w, val.String())
	return nil
}

func writeBytes(w io.Writer, val reflect.Value) error {
	bpl.PackBin16(w, val.Bytes())
	return nil
}
