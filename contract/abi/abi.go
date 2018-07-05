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
 * file description:  abi
 * @Author: Gong Zibin
 * @Date:   2017-01-20
 * @Last Modified by:
 * @Last Modified time:
 */

package abi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bottos-project/bottos/contract/msgpack"
	log "github.com/cihub/seelog"
	"io"
	"reflect"
)

//ABIAction abi Action(Method)
type ABIAction struct {
	ActionName string `json:"action_name"`
	Type       string `json:"type"`
}

//ABIStruct parameter struct for abi Action(Method)
type ABIStruct struct {
	Name   string    `json:"name"`
	Base   string    `json:"base"`
	Fields *FeildMap `json:"fields"`
}

//ABI struct for abi
type ABI struct {
	Types   []interface{} `json:"types"`
	Structs []ABIStruct   `json:"structs"`
	Actions []ABIAction   `json:"actions"`
	Tables  []interface{} `json:"tables"`
}

//ABIStructs structs for ABI
type ABIStructs struct {
	Structs []struct {
		Name   string            `json:"name"`
		Base   string            `json:"base"`
		Fields map[string]string `json:"fields"`
	} `json:"structs"`
}

//ParseAbi parse abiraw to struct for contracts
func ParseAbi(abiRaw []byte) (*ABI, error) {
	abis := &ABIStructs{}
	err := json.Unmarshal(abiRaw, abis)
	if err != nil {
		return &ABI{}, err
	}

	abi := &ABI{}
	abi.Structs = make([]ABIStruct, len(abis.Structs))
	for i := range abi.Structs {
		abi.Structs[i].Fields = New()
	}
	err = json.Unmarshal(abiRaw, abi)
	if err != nil {
		return &ABI{}, err
	}

	return abi, nil
}

//AbiToJson parse abi to json for contracts
func AbiToJson(abi *ABI) (string, error) {
	data, err := json.Marshal(abi)
	if err != nil {
		return "", err
	}
	return jsonFormat(data), nil
}

func jsonFormat(data []byte) string {
	var out bytes.Buffer
	json.Indent(&out, data, "", "    ")

	return string(out.Bytes())
}

//MarshalAbi is to serialize the message
func MarshalAbi(v interface{}, Abi *ABI, contractName string, method string) ([]byte, error) {
	var err error
	var abi ABI

	if Abi == nil {
		return []byte{}, err
	}
	
	abi = *Abi
	

	writer := &bytes.Buffer{}
	err = EncodeAbi(contractName, method, writer, v, abi, "")
	if err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}


func getAbiFieldsByAbi(contractname string, method string, abi ABI, subStructName string) map[string]interface{} {
	for _, subaction := range abi.Actions {
		if subaction.ActionName != method {
			continue
		}

		structname := subaction.Type

		for _, substruct := range abi.Structs {
			if subStructName != "" {
				if substruct.Name != subStructName {
					continue
				}
			} else if structname != substruct.Name {
				continue
			}

			return substruct.Fields.values
		}
	}

	return nil
}

//DecodeAbi is to encode message
func DecodeAbi(contractName string, method string, r io.Reader, dst interface{}, abi ABI, subStructName string) error {
	abiFields := getAbiFieldsByAbi(contractName, method, abi, subStructName)
	if abiFields == nil {
		return fmt.Errorf("DecodeAbi: getAbiFieldsByAbi failed: %s", abi)
	}
	v := reflect.ValueOf(dst)

	if !v.IsValid() {
		log.Errorf("Not Valid %T\n", dst)
		return fmt.Errorf("Not Valid %T\n", dst)
	}

	if v.Kind() != reflect.Ptr {
		log.Errorf("dst Not Settable %T\n", dst)
		return fmt.Errorf("dst Not Settable %T)", dst)
	}

	if !v.Elem().IsValid() {
		log.Errorf("Nil Ptr: %T\n", dst)
		return fmt.Errorf("Nil Ptr: %T\n", dst)
	}

	if v.Elem().NumField() > 0 {
		msgpack.UnpackArraySize(r)
	}

	v = v.Elem()

	vt := reflect.TypeOf(dst)
	vt = vt.Elem()

	count := v.NumField()
	
	for i := 0; i < count; i++ {
		
		field := v.Field(i)
		feildAddr := field.Addr().Interface()

		fieldname := vt.Field(i).Tag.Get("json")
		if _, ok := abiFields[fieldname]; !ok {
			return fmt.Errorf("DecodeAbi: getAbiFieldsByAbi failed: %s, %s", abi, !ok)
		}

		switch abiFields[fieldname] {
		case "string":
			val, err := msgpack.UnpackStr16(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*string)
			*ptr = val
		case "uint8":
			val, err := msgpack.UnpackUint8(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*uint8)
			*ptr = val
		case "uint16":
			val, err := msgpack.UnpackUint16(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*uint16)
			*ptr = val

		case "uint32":
			val, err := msgpack.UnpackUint32(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*uint32)
			*ptr = val
		case "uint64":
			val, err := msgpack.UnpackUint64(r)
			if err != nil {
				return err
			}
			ptr := feildAddr.(*uint64)
			*ptr = val
		case "bytes":
			t := reflect.TypeOf(v.Field(i).Interface())
			if t.Elem().Kind() == reflect.Uint8 {
				val, err := msgpack.UnpackBin16(r)
				if err != nil {
					return err
				}
				ptr := feildAddr.(*[]byte)
				*ptr = val
			} else {
				return fmt.Errorf("Unsupported Slice Type")
			}
		default:
			vt = reflect.TypeOf(v.Field(i).Interface())
			if vt.Kind() == reflect.Struct {
				DecodeAbi(contractName, method, r, v.Field(i).Interface(), abi, fieldname)
			} else if vt.Kind() == reflect.Ptr {
				DecodeAbi(contractName, method, r, v.Elem().Field(i).Interface(), abi, fieldname)
			} else {
				return fmt.Errorf("Unsupported Type: %v", vt.Kind())
			}
		}

	}

	return nil
}

//EncodeAbi is to encode message
func EncodeAbi(contractName string, method string, w io.Writer, value interface{}, abi ABI, subStructName string) error {
	abiFields := getAbiFieldsByAbi(contractName, method, abi, subStructName)
	if abiFields == nil {
		return fmt.Errorf("EncodeAbi: getAbiFieldsByAbi failed: %s", abi)
	}

	v := reflect.ValueOf(value)
	vt := reflect.TypeOf(value)

	if !v.IsValid() {
		log.Errorf("Not Valid %T\n", value)
		return fmt.Errorf("Not Valid %T\n", value)
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		vt = vt.Elem()
		if !v.IsValid() {
			log.Errorf("Nil Ptr: %T\n", value)
			return fmt.Errorf("Nil Ptr: %T\n", value)
		}
	}

	count := v.NumField()
	msgpack.PackArraySize(w, uint16(count))

	for i := 0; i < count; i++ {
		fieldname := vt.Field(i).Tag.Get("json")
		vals := v.Field(i).Interface()

		types := reflect.TypeOf(vals)
		val := reflect.ValueOf(vals)

		if _, ok := abiFields[fieldname]; !ok {
			return fmt.Errorf("%s is not in abiFields [%s]!", fieldname, abiFields)
		}

		switch abiFields[fieldname] {
		case "string":
			msgpack.PackStr16(w, val.String())
		case "uint8":
			msgpack.PackUint8(w, uint8(val.Uint()))
		case "uint16":
			msgpack.PackUint16(w, uint16(val.Uint()))
		case "uint32":
			msgpack.PackUint32(w, uint32(val.Uint()))
		case "uint64":
			msgpack.PackUint64(w, uint64(val.Uint()))
		case "bytes":
			t := reflect.TypeOf(v.Field(i).Interface())
			if t.Elem().Kind() == reflect.Uint8 {
				msgpack.PackBin16(w, val.Bytes())
			} else {
				return fmt.Errorf("Unsupported Slice Type")
			}
		default:
			t := reflect.TypeOf(v.Field(i).Interface())
			if t.Kind() == reflect.Struct || t.Kind() == reflect.Ptr {
				EncodeAbi(contractName, method, w, v.Field(i).Interface(), abi, fieldname)
			} else {
				return fmt.Errorf("Unsupported Type: %v", types)
			}
		}
	}

	return nil
}

//UnmarshalAbi is to unserialize the message
func UnmarshalAbi(contractName string, Abi *ABI, method string, data []byte, dst interface{}) error {
	var err error
	var abi ABI
	if Abi == nil {
	    return err
	}
	
	abi = *Abi

	r := bytes.NewReader(data)
	err = DecodeAbi(contractName, method, r, dst, abi, "")
	return err
}
