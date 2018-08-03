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
func Decode(v interface{}, r io.Reader) error {
	if v == nil {
		return errors.New("msgpack: nil pointer")
	}

	rv := reflect.ValueOf(v)
	t := rv.Type()
	if t.Kind() != reflect.Ptr {
		return errors.New("msgpack: need a pointer")
	}
	if rv.IsNil() {
		return errors.New("msgpack: need a pointer")
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
	case kind == reflect.Struct:
		return makeStructDecoder(t, r)
	case kind == reflect.Ptr:
		return makePtrDecoder(t, r)
	default:
		return nil, fmt.Errorf("msgpack, type %v, kind %v not support", t, kind)
	}
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
	fmt.Println(etype)
	decoder, err := getDecoder(etype, r)
	if err != nil {
		return nil, err
	}
	dec := func(val reflect.Value, r io.Reader) (err error) {
		newval := val
		if val.IsNil() {
			newval = reflect.New(etype)
		}
		if err = decoder(newval.Elem(), r); err == nil {
			fmt.Println(newval.Type(), val.Type())
			val.Set(newval)
		}
		return err
	}
	return dec, nil
}
