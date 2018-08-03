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
 * file description:  msgpack encode
 * @Author: Gong Zibin
 * @Date:   2018-08-02
 * @Last Modified by:
 * @Last Modified time:
 */

package msgpack

import (
	"fmt"
	"io"
	"reflect"
)

type EncodeWriter func(reflect.Value, io.Writer) error

type Encoder interface {
	EncodeMsgpack(io.Writer) error
}

var (
	encoderInterface = reflect.TypeOf(new(Encoder)).Elem()
)

//Encode is to encode message
func Encode(v interface{}, w io.Writer) error {
	rv := reflect.ValueOf(v)
	encoder, err := getEncoder(rv.Type(), w)
	if err != nil {
		return err
	}

	return encoder(rv, w)
}

func getEncoder(t reflect.Type, w io.Writer) (EncodeWriter, error) {
	kind := t.Kind()
	switch {
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
	case kind == reflect.Struct:
		return makeStructEncoder(t, w)
	case kind == reflect.Ptr:
		return makePtrEncoder(t, w)
	case t.Implements(encoderInterface):
		return writeEncoder, nil
	case kind != reflect.Ptr && reflect.PtrTo(t).Implements(encoderInterface):
		return writeEncoderNoPtr, nil
	default:
		return nil, fmt.Errorf("msgpack, type %v not support", t)
	}
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

type Field struct {
	encoder EncodeWriter
	index   int
}

func makeStructEncoder(t reflect.Type, w io.Writer) (EncodeWriter, error) {
	fields := []Field{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		encoder, err := getEncoder(f.Type, w)
		if err != nil {
			return nil, err
		}
		fields = append(fields, Field{encoder, i})
	}

	encoder := func(val reflect.Value, w io.Writer) error {
		PackArraySize(w, uint16(len(fields)))
		for _, f := range fields {
			if err := f.encoder(val.Field(f.index), w); err != nil {
				return err
			}
		}
		return nil
	}
	return encoder, nil
}

func makePtrEncoder(t reflect.Type, w io.Writer) (EncodeWriter, error) {
	encodeWriter, err := getEncoder(t.Elem(), w)
	if err != nil {
		return nil, err
	}

	encoder := func(val reflect.Value, w io.Writer) error {
		if val.IsNil() {
			_, err := PackNil(w)
			return err
		} else {
			return encodeWriter(val.Elem(), w)
		}
	}

	return encoder, nil
}

func writeEncoder(val reflect.Value, w io.Writer) error {
	fmt.Println(val)
	return val.Interface().(Encoder).EncodeMsgpack(w)
}

func writeEncoderNoPtr(val reflect.Value, w io.Writer) error {
	fmt.Println(val)
	if !val.CanAddr() {
		return fmt.Errorf("rlp: game over: unadressable value of type %v, EncodeRLP is pointer method", val.Type())
	}
	return val.Addr().Interface().(Encoder).EncodeMsgpack(w)
}
