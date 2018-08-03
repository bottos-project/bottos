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
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type EncodeWriter func(reflect.Value, io.Writer) error

//Marshal is to serialize the message
func Marshal(v interface{}) ([]byte, error) {
	writer := &bytes.Buffer{}
	err := Encode(v, writer)
	if err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}

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
	case kind == reflect.Struct:
		return makeStructEncoder(t, w)
	case kind == reflect.Ptr:
		return makePtrEncoder(t, w)
	default:
		return nil, fmt.Errorf("msgpack, type %v not support", t)
	}
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
			return fmt.Errorf("msgpack: Ptr is Nil")
		} else {
			return encodeWriter(val.Elem(), w)
		}
	}

	return encoder, nil
}

/*
//Encode is to encode message
func Encode(w io.Writer, structs interface{}) error {
	v := reflect.ValueOf(structs)

	if !v.IsValid() {
		log.Errorf("Not Valid %T\n", structs)
		return fmt.Errorf("Not Valid %T\n", structs)
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if !v.IsValid() {
			log.Errorf("Nil Ptr: %T\n", structs)
			return fmt.Errorf("Nil Ptr: %T\n", structs)
		}
	}

	values := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	PackArraySize(w, uint16(len(values)))
	for i := 0; i < len(values); i++ {
		t := reflect.TypeOf(values[i])
		val := reflect.ValueOf(values[i])

		kind := t.Kind()
		switch kind {
		case reflect.String:
			PackStr16(w, val.String())
		case reflect.Uint8:
			PackUint8(w, uint8(val.Uint()))
		case reflect.Uint16:
			PackUint16(w, uint16(val.Uint()))
		case reflect.Uint32:
			PackUint32(w, uint32(val.Uint()))
		case reflect.Uint64:
			PackUint64(w, uint64(val.Uint()))
		case reflect.Slice: // []byte
			if t.Elem().Kind() == reflect.Uint8 {
				PackBin16(w, val.Bytes())
			} else {
				return fmt.Errorf("Unsupported Slice Type")
			}
		case reflect.Struct:
			Encode(w, val.Interface())
		case reflect.Ptr:
			vvt := reflect.TypeOf(val.Elem().Interface())
			if vvt.Kind() == reflect.Struct {
				Encode(w, val.Elem().Interface())
			} else {
				return fmt.Errorf("Unsupported Type: %T", val)
			}
		default:
			return fmt.Errorf("Unsupported Type: %v", kind)
		}
	}
	return nil
}
*/
