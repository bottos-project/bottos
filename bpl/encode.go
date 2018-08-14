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
 * file description:  bpl encode
 * @Author: Gong Zibin
 * @Date:   2018-08-02
 * @Last Modified by:
 * @Last Modified time:
 */

package bpl

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"sync"
)

//EncodeWriter function type of the encoder
type EncodeWriter func(reflect.Value, io.Writer) error

var encoderCache sync.Map // map[reflect.Type]EncodeWriter

//Encoder interface for customization
type Encoder interface {
	EncodeMsgpack(io.Writer) error
}

const (
	EXT_BIGINT = 1
)

var (
	encoderInterface = reflect.TypeOf(new(Encoder)).Elem()
	bigInt           = reflect.TypeOf(big.Int{})
)

//Encode encodes struct, ptr, slice/array and basic types to byte stream
func Encode(v interface{}, w io.Writer) error {
	rv := reflect.ValueOf(v)
	encoder, err := getEncoder(rv.Type())
	if err != nil {
		return err
	}
	return encoder(rv, w)
}

func newEncoder(t reflect.Type) (EncodeWriter, error) {
	kind := t.Kind()
	switch {
	case t.Implements(encoderInterface):
		return encodeCustom, nil
	case kind != reflect.Ptr && reflect.PtrTo(t).Implements(encoderInterface):
		return encodeCustomNoPtr, nil
	case t.AssignableTo(reflect.PtrTo(bigInt)):
		return encodeBigIntPtr, nil
	case t.AssignableTo(bigInt):
		return encodeBigIntNoPtr, nil
	case kind == reflect.Bool:
		return encodeBool, nil
	case kind == reflect.Uint8:
		return encodeUint8, nil
	case kind == reflect.Uint16:
		return encodeUint16, nil
	case kind == reflect.Uint32:
		return encodeUint32, nil
	case kind == reflect.Uint64:
		return encodeUint64, nil
	case kind == reflect.String:
		return encodeString, nil
	case kind == reflect.Slice && t.Elem().Kind() == reflect.Uint8:
		return encodeBytes, nil
	case kind == reflect.Array && t.Elem().Kind() == reflect.Uint8:
		return encodeByteArray, nil
	case kind == reflect.Slice || kind == reflect.Array:
		return makeSliceEncoder(t)
	case kind == reflect.Struct:
		return makeStructEncoder(t)
	case kind == reflect.Ptr:
		return makePtrEncoder(t)
	default:
		return unsupportedTypeEncoder, nil
	}
}

func getEncoder(t reflect.Type) (EncodeWriter, error) {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(EncodeWriter), nil
	}

	var (
		wg sync.WaitGroup
		f  EncodeWriter
	)
	wg.Add(1)
	fi, loaded := encoderCache.LoadOrStore(t, EncodeWriter(func(val reflect.Value, w io.Writer) error {
		wg.Wait()
		return f(val, w)
	}))
	if loaded {
		return fi.(EncodeWriter), nil
	}

	encoder, err := newEncoder(t)
	if err == nil {
		encoderCache.Store(t, encoder)
	}
	return encoder, err
}

func unsupportedTypeEncoder(val reflect.Value, w io.Writer) error {
	return fmt.Errorf("msgpack: supported type %v", val.Type())
}

func encodeBool(val reflect.Value, w io.Writer) error {
	PackBool(w, val.Bool())
	return nil
}

func encodeUint8(val reflect.Value, w io.Writer) error {
	PackUint8(w, uint8(val.Uint()))
	return nil
}

func encodeUint16(val reflect.Value, w io.Writer) error {
	PackUint16(w, uint16(val.Uint()))
	return nil
}

func encodeUint32(val reflect.Value, w io.Writer) error {
	PackUint32(w, uint32(val.Uint()))
	return nil
}

func encodeUint64(val reflect.Value, w io.Writer) error {
	PackUint64(w, uint64(val.Uint()))
	return nil
}

func encodeString(val reflect.Value, w io.Writer) error {
	PackStr16(w, val.String())
	return nil
}

func encodeBytes(val reflect.Value, w io.Writer) error {
	PackBin16(w, val.Bytes())
	return nil
}

func encodeByteArray(val reflect.Value, w io.Writer) error {
	if !val.CanAddr() {
		copy := reflect.New(val.Type()).Elem()
		copy.Set(val)
		val = copy
	}
	size := val.Len()
	slice := val.Slice(0, size).Bytes()
	PackBin16(w, slice)
	return nil
}

func makeSliceEncoder(t reflect.Type) (EncodeWriter, error) {
	elemEncoder, err := getEncoder(t.Elem())
	if err != nil {
		return nil, err
	}

	encoder := func(val reflect.Value, w io.Writer) error {
		vlen := val.Len()
		PackArraySize(w, uint16(vlen))
		for i := 0; i < vlen; i++ {
			if err := elemEncoder(val.Index(i), w); err != nil {
				return err
			}
		}
		return nil
	}
	return encoder, nil
}

type structField struct {
	t     reflect.Type
	index int
}

func makeStructEncoder(t reflect.Type) (EncodeWriter, error) {
	fields := []structField{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fields = append(fields, structField{f.Type, i})
	}

	encoder := func(val reflect.Value, w io.Writer) error {
		PackArraySize(w, uint16(len(fields)))
		for _, f := range fields {
			encoder, err := getEncoder(f.t)
			if err != nil {
				return err
			}
			if err := encoder(val.Field(f.index), w); err != nil {
				return err
			}
		}
		return nil
	}
	return encoder, nil
}

func makePtrEncoder(t reflect.Type) (EncodeWriter, error) {
	encodeWriter, err := getEncoder(t.Elem())
	if err != nil {
		return nil, err
	}

	encoder := func(val reflect.Value, w io.Writer) error {
		if val.IsNil() {
			_, err := PackNil(w)
			return err
		}
		return encodeWriter(val.Elem(), w)
	}

	return encoder, nil
}

func encodeBigIntPtr(val reflect.Value, w io.Writer) error {
	ptr := val.Interface().(*big.Int)
	if ptr == nil {
		return errors.New("nil ptr")
	}
	return encodeBigInt(ptr, w)
}

func encodeBigIntNoPtr(val reflect.Value, w io.Writer) error {
	i := val.Interface().(big.Int)
	return encodeBigInt(&i, w)
}

func encodeBigInt(i *big.Int, w io.Writer) error {
	_, err := PackExt16(w, EXT_BIGINT, i.Bytes())
	return err
}

func encodeCustom(val reflect.Value, w io.Writer) error {
	return val.Interface().(Encoder).EncodeMsgpack(w)
}

func encodeCustomNoPtr(val reflect.Value, w io.Writer) error {
	if !val.CanAddr() {
		return fmt.Errorf("msgpack encode: unadressable value of type %v", val.Type())
	}
	return val.Addr().Interface().(Encoder).EncodeMsgpack(w)
}
