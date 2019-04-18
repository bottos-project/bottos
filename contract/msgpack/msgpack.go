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

	log "github.com/cihub/seelog"
)

//Marshal is to serialize the message
func Marshal(v interface{}) ([]byte, error) {
	writer := &bytes.Buffer{}
	err := Encode(writer, v)
	if err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}

//Unmarshal is to unserialize the message
func Unmarshal(data []byte, dst interface{}) error {
	r := bytes.NewReader(data)
	err := Decode(r, dst)
	return err
}

//Encode is to encode message
func Encode(w io.Writer, structs interface{}) error {
	v := reflect.ValueOf(structs)

	if !v.IsValid() {
		log.Errorf("MSGPACK encode error, Not Valid %T\n", structs)
		return fmt.Errorf("Not Valid %T\n", structs)
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if !v.IsValid() {
			log.Errorf("MSGPACK encode error, Nil Ptr: %T\n", structs)
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

//Decode is to encode message
func Decode(r io.Reader, dst interface{}) error {
	v := reflect.ValueOf(dst)

	if !v.IsValid() {
		log.Errorf("MSGPACK decode error, Not Valid %T\n", dst)
		return fmt.Errorf("Not Valid %T\n", dst)
	}

	if v.Kind() != reflect.Ptr {
		log.Errorf("MSGPACK decode error, dst Not Settable %T\n", dst)
		return fmt.Errorf("dst Not Settable %T)", dst)
	}

	if !v.Elem().IsValid() {
		log.Errorf("MSGPACK decode error, Nil Ptr: %T\n", dst)
		return fmt.Errorf("Nil Ptr: %T\n", dst)
	}

	if v.Elem().NumField() > 0 {
		UnpackArraySize(r)
	}

	for i := 0; i < v.Elem().NumField(); i++ {
		feild := v.Elem().Field(i)
		feildAddr := feild.Addr().Interface()
		switch feild.Kind() {
		case reflect.String:
			val, err := UnpackStr16(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*string)
			*ptr = val
		case reflect.Uint8:
			val, err := UnpackUint8(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*uint8)
			*ptr = val
		case reflect.Uint16:
			val, err := UnpackUint16(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*uint16)
			*ptr = val
		case reflect.Uint32:
			val, err := UnpackUint32(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*uint32)
			*ptr = val
		case reflect.Uint64:
			val, err := UnpackUint64(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*uint64)
			*ptr = val
		case reflect.Slice:
			if feild.Type().Elem().Kind() == reflect.Uint8 {
				val, err := UnpackBin16(r)
				if err != nil {
					return err
				}
				ptr := feildAddr.(*[]byte)
				*ptr = val
			} else {
				return fmt.Errorf("Unsupported Slice Type")
			}
		case reflect.Struct:
			vv := feild.Interface()
			if reflect.ValueOf(vv).Kind() != reflect.Ptr {
				vv = feildAddr
			}
			Decode(r, vv)
		case reflect.Ptr:
			vv := feild.Interface()
			Decode(r, vv)
		default:
			return fmt.Errorf("Unsupported Type")
		}
	}

	return nil
}
