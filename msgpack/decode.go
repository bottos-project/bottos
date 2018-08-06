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
 * file description:  msgpack decode
 * @Author: Gong Zibin
 * @Date:   2018-08-02
 * @Last Modified by:
 * @Last Modified time:
 */

package msgpack

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
)

type DecoderReader func(reflect.Value, io.Reader) error

type Decoder interface {
	DecodeMsgpack(io.Reader) error
}

var (
	decoderInterface = reflect.TypeOf(new(Decoder)).Elem()
)

//Decode is to decode message
func Decode(r io.Reader, v interface{}) error {
	if v == nil {
		return errors.New("msgpack decode: nil pointer")
	}

	rv := reflect.ValueOf(v)
	t := rv.Type()
	if t.Kind() != reflect.Ptr {
		return errors.New("msgpack decode: need a pointer")
	}
	if rv.IsNil() {
		return errors.New("msgpack decode: nil pointer")
	}

	decoder, err := getDecoder(t.Elem(), r)
	if err != nil {
		return err
	}

	return decoder(rv.Elem(), r)
}

func getDecoder(t reflect.Type, r io.Reader) (DecoderReader, error) {
	kind := t.Kind()
	switch {
	case t.AssignableTo(reflect.PtrTo(bigInt)):
		return decodeBigInt, nil
	case t.AssignableTo(bigInt):
		return decodeBigIntNoPtr, nil
	case kind == reflect.Bool:
		return decodeBool, nil
	case kind == reflect.Uint8:
		return decodeUint8, nil
	case kind == reflect.Uint16:
		return decodeUint16, nil
	case kind == reflect.Uint32:
		return decodeUint32, nil
	case kind == reflect.Uint64:
		return decodeUint64, nil
	case kind == reflect.String:
		return decodeString, nil
	case kind == reflect.Slice && t.Elem().Kind() == reflect.Uint8:
		return decodeBytes, nil
	case kind == reflect.Array && t.Elem().Kind() == reflect.Uint8:
		return decodeByteArray, nil
	case kind == reflect.Slice || kind == reflect.Array:
		return makeArrayDecoder(t, r)
	case kind == reflect.Struct:
		return makeStructDecoder(t, r)
	case kind == reflect.Ptr:
		return makePtrDecoder(t, r)
	default:
		return nil, fmt.Errorf("msgpack, type %v, kind %v not support", t, kind)
	}
}

func decodeBigIntNoPtr(v reflect.Value, r io.Reader) error {
	return decodeBigInt(v.Addr(), r)
}

func decodeBigInt(v reflect.Value, r io.Reader) error {
	val, t, err := UnpackExt16(r)
	if err != nil {
		return errors.New("msgpack: unpack ext fail")
	}
	if t != EXT_BIGINT {
		return errors.New("msgpack: unpack ext type error")
	}
	i := v.Interface().(*big.Int)
	if i == nil {
		i = new(big.Int)
		v.Set(reflect.ValueOf(i))
	}
	i.SetBytes(val)
	return nil
}

func decodeBytes(v reflect.Value, r io.Reader) error {
	val, err := UnpackBin16(r)
	if err != nil {
		return err
	}
	v.SetBytes(val)
	return nil
}

func decodeByteArray(v reflect.Value, r io.Reader) error {
	vlen := v.Len()
	slice := v.Slice(0, vlen).Interface().([]byte)
	val, err := UnpackBin16(r)
	if err != nil {
		return err
	}

	if len(val) != vlen {
		return errors.New("msgpack: wrong array size")
	}

	copy(slice, val)

	return nil
}

func decodeBool(v reflect.Value, r io.Reader) error {
	val, err := UnpackBool(r)
	if err != nil {
		return err
	}
	v.SetBool(val)
	return nil
}

func decodeUint8(v reflect.Value, r io.Reader) error {
	val, err := UnpackUint8(r)
	if err != nil {
		return err
	}
	v.SetUint(uint64(val))
	return nil
}

func decodeUint16(v reflect.Value, r io.Reader) error {
	val, err := UnpackUint16(r)
	if err != nil {
		return err
	}
	v.SetUint(uint64(val))
	return nil
}

func decodeUint32(v reflect.Value, r io.Reader) error {
	val, err := UnpackUint32(r)
	if err != nil {
		return err
	}
	v.SetUint(uint64(val))
	return nil
}

func decodeUint64(v reflect.Value, r io.Reader) error {
	val, err := UnpackUint64(r)
	if err != nil {
		return err
	}
	v.SetUint(uint64(val))
	return nil
}

func decodeString(v reflect.Value, r io.Reader) error {
	val, err := UnpackStr16(r)
	if err != nil {
		return err
	}
	v.SetString(string(val))
	return nil
}

func decodeArray(v reflect.Value, r io.Reader, elemDecoder DecoderReader) error {
	vlen := v.Len()
	size, err := UnpackArraySize(r)
	if err != nil {
		return err
	}

	if vlen != int(size) {
		return fmt.Errorf("msgpack decoder: wrong array size %v, expected %v", int(size), vlen)
	}
	i := 0
	for ; i < vlen; i++ {
		if err := elemDecoder(v.Index(i), r); err != nil {
			return fmt.Errorf("msgpack decoder: array decode error")
		}
	}
	if i < vlen {
		return fmt.Errorf("msgpack decoder: array size array")
	}
	return nil
}

func decodeSlice(v reflect.Value, r io.Reader, elemDecoder DecoderReader) error {
	size, err := UnpackArraySize(r)
	if err != nil {
		return err
	}
	if size == 0 {
		v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		return nil
	}

	i := 0
	for ; i < int(size); i++ {
		if i >= v.Cap() {
			newcap := v.Cap() + v.Cap()/2
			if newcap < 4 {
				newcap = 4
			}
			newv := reflect.MakeSlice(v.Type(), v.Len(), newcap)
			reflect.Copy(newv, v)
			v.Set(newv)
		}
		if i >= v.Len() {
			v.SetLen(i + 1)
		}
		if err := elemDecoder(v.Index(i), r); err != nil {
			return fmt.Errorf("msgpack decoder: slice decode error")
		}
	}
	return nil
}

func makeArrayDecoder(t reflect.Type, r io.Reader) (DecoderReader, error) {
	etype := t.Elem()
	elemDecoder, err := getDecoder(etype, r)
	if err != nil {
		return nil, err
	}

	var dec DecoderReader
	switch {
	case t.Kind() == reflect.Array:
		dec = func(val reflect.Value, r io.Reader) error {
			return decodeArray(val, r, elemDecoder)
		}
	default: // t.Kind() == reflect.Slice
		dec = func(val reflect.Value, r io.Reader) error {
			return decodeSlice(val, r, elemDecoder)
		}
	}
	return dec, nil
}

type DecField struct {
	decoder DecoderReader
	index   int
}

func makeStructDecoder(t reflect.Type, r io.Reader) (DecoderReader, error) {
	fields := []DecField{}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		decoder, e := getDecoder(f.Type, r)
		if e != nil {
			return nil, e
		}
		fields = append(fields, DecField{decoder, i})
	}

	dec := func(val reflect.Value, r io.Reader) (err error) {
		_, err = UnpackArraySize(r)
		if err != nil {
			return err
		}
		for _, f := range fields {
			err = f.decoder(val.Field(f.index), r)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return dec, nil
}

func makePtrDecoder(t reflect.Type, r io.Reader) (DecoderReader, error) {
	etype := t.Elem()
	decoder, err := getDecoder(etype, r)
	if err != nil {
		return nil, err
	}
	dec := func(val reflect.Value, r io.Reader) (err error) {
		newval := val
		if val.IsNil() {
			newval = reflect.New(etype)
		}
		/*
			suc, err := TryUnpackNil(r)
			fmt.Println(suc, err)
			if err != nil {
				return err
			}

			if suc {
				val.Set(reflect.Zero(t))
				fmt.Println("no")
			}
		*/

		if err = decoder(newval.Elem(), r); err == nil {
			//fmt.Println(newval.Type(), val.Type())
			val.Set(newval)
		}
		return err
	}
	return dec, nil
}
