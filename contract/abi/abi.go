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
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/bpl"
	"github.com/bottos-project/bottos/contract/abi/fieldmap"
	//log "github.com/cihub/seelog"
	"io"
	"reflect"
	"strings"
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
	Fields *fieldmap.FeildMap `json:"fields"`
}

//ABI struct for abi
type ABI struct {
	Types   []interface{} `json:"types"`
	Structs []ABIStruct   `json:"structs"`
	Actions []ABIAction   `json:"actions"`
	Tables  []interface{} `json:"tables"`
}

type ABI1 struct {
	AbiDef *ABIDef
	Methods map[string]Method
}

//NewABIFromJSON parse abi json definition
func NewABIFromJSON(abiJson []byte) (*ABI1, error) {
	def := &ABIDef{}
	err := json.Unmarshal(abiJson, def)
	if err != nil {
		return nil, err
	}
	return NewABIFromDef(def)
}

func NewABIFromDef(abiDef *ABIDef) (*ABI1, error) {
	fmt.Println("abiDef: ", abiDef)
	a := ABI1{
		AbiDef: abiDef,
		Methods: make(map[string]Method),
	}
	for _, method := range abiDef.Methods {
		m := Method{Name: method.Name}
		for _, defStruct := range abiDef.Structs {
			if method.Type == defStruct.Name {
				m.Fields = defStruct.Fields
				a.Methods[method.Name] = m
				break
			}
		}
	}

	fmt.Println("abi: ",a)

	return &a, nil
}

func (abi *ABI1) Pack(methodName string, args ...interface{}) ([]byte, error)  {
	method, exist := abi.Methods[methodName]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", methodName)
	}

	argNames := method.GetFieldNames()
	if len(args) != len(argNames) {
		return nil, fmt.Errorf("method '%s' arguments mismatch", methodName)
	}

	w := &bytes.Buffer{}
	bpl.PackArraySize(w, uint16(len(args)))
	for i, a := range args {
		name := argNames[i]
		typ, exist := method.GetFieldType(name)
		if !exist {
			return nil, fmt.Errorf("field '%s' not found", name)
		}
		switch typ {
		case "string":
			bpl.PackStr16(w, a.(string))
		case "uint64":
			bpl.PackUint64(w, a.(uint64))
		}
	}

	return w.Bytes(), nil
}

func ParseAbi(abiRaw []byte) (*ABI, error) {
	a := &ABI{}
	err := json.Unmarshal(abiRaw, a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

//AbiToJson parse abi to json for contracts
func (abi *ABI1) ToJson(beautify bool) string {
	data, _ := json.Marshal(abi.AbiDef)
	if beautify {
		return jsonFormat(data)
	}
	return string(data)
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

//getAbiFieldsByAbiEx function
func getAbiFieldsByAbiEx(contractname string, method string, abi ABI, subStructName string) *fieldmap.FeildMap {
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

			return substruct.Fields
		}
	}

	return nil
}

//EncodeAbiEx is to encode message
func EncodeAbiEx(contractName string, method string, w io.Writer, value map[string]interface{}, abi ABI, subStructName string) error {
        abiFieldsAttr := getAbiFieldsByAbiEx(contractName, method, abi, subStructName)
	if abiFieldsAttr == nil {
		return fmt.Errorf("EncodeAbiEx: getAbiFieldsByAbi failed: %s", abi)

	}

	abiFields := abiFieldsAttr.GetStringPair()
	
	count  := len(abiFields)
	count2 := len(value)
	
	if count != count2 {
		return fmt.Errorf("EncodeAbiEx: fields number mismatch! count: %d, count2: %d", count, count2)
	}
	
	if (count <= 0) {
		return fmt.Errorf("EncodeAbiEx: count is 0!")
	}

	bpl.PackArraySize(w, uint16(count))

		for _, abiValTypeAttr := range abiFields {
			abiValKey   := abiValTypeAttr.Key
			abiValType := abiValTypeAttr.Value

			val, ok := value[abiValKey]
			
			if !ok {
				return fmt.Errorf("EncodeAbiEx: value abiValKey %s not found in map", abiValKey)
			}
			
			valType := reflect.TypeOf(val).Name()
			
			if reflect.ValueOf(val).Kind() == reflect.Slice {
				valType = reflect.TypeOf(val).Elem().Name()
				if valType == "uint8"	{
					valType = "bytes"
				}
			}
			
			if valType != abiValType {
				return fmt.Errorf("EncodeAbiEx: abiValType %s mismatch to valType %s", abiValType, valType)
			}

			switch abiValType {
				case "string":
					bpl.PackStr16(w, val.(string))
				case "uint8":
					bpl.PackUint8(w, val.(uint8))
				case "uint16":
					bpl.PackUint16(w, val.(uint16))
				case "uint32":
					bpl.PackUint32(w, val.(uint32))
				case "uint64":
					bpl.PackUint64(w, val.(uint64))
				case "bytes":
					bpl.PackBin16(w, val.([]byte))
				default:
					if reflect.ValueOf(value[abiValKey]).Kind() == reflect.Struct {
						EncodeAbiEx(contractName, method, w, value, abi, abiValKey)
					} else {
						return fmt.Errorf("Unsupported Type: %v | %v", valType, abiValType)
					}
				}
		}
		
	return nil
}

func Setmapval(structmap map[string]interface{}, key string, val interface{}) {
        structmap[key] = val
}

//MarshalAbiEx is to serialize the message
func MarshalAbiEx(v map[string]interface{}, Abi *ABI, contractName string, method string) ([]byte, error) {
	var err error
	var abi ABI
	
	
	if Abi == nil {
		return []byte{}, err
	}
	
	abi = *Abi

	writer := &bytes.Buffer{}
	err = EncodeAbiEx(contractName, method, writer, v, abi, "")
	if err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}

//DecodeAbiEx is to encode message
func DecodeAbiEx(contractName string, method string, r io.Reader, abi ABI, subStructName string, subStructValueName string, mapResultIn *map[string]interface{}) (map[string]interface{}) {
	var errs error
	mapResult := make(map[string]interface{})
	
	if(mapResultIn != nil && len(subStructName) > 0) {
		mapResult = *mapResultIn 
		mapResult[subStructValueName] = make(map[string]interface{})
		mapResult = mapResult[subStructValueName].(map[string]interface{})
	}
	
	abiFieldsAttr := getAbiFieldsByAbiEx(contractName, method, abi, subStructName)
	if abiFieldsAttr == nil {
		return nil
	}
	
	abiFields := abiFieldsAttr.GetStringPair()
	
	count  := len(abiFields)
	if (count <= 0) {
		return nil
	}
	
	if len(abiFields) > 0 {
		_, errs = bpl.UnpackArraySize(r)
		if errs != nil {
			return nil
		}
	} else {
		return nil
	}
	var i uint64 = 0
	for _, abiValTypeAttr := range abiFields {
			abiValKey   := strings.ToLower(abiValTypeAttr.Key)
			abiValType := abiValTypeAttr.Value

			switch abiValType {
				case "string":
					val, err := bpl.UnpackStr16(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "uint8":
					val, err := bpl.UnpackUint8(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "uint16":
					val, err := bpl.UnpackUint16(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "uint32":
					val, err := bpl.UnpackUint32(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "uint64":
					val, err := bpl.UnpackUint64(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "bytes":
					val, err := bpl.UnpackBin16(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, common.BytesToHex(val))
					i++
				default:
					DecodeAbiEx(contractName, method, r, abi, abiValType, abiValKey, &mapResult)
				}
			i += 1
		}
	
	return mapResult
}

//UnmarshalAbiEx is to unserialize the message
func UnmarshalAbiEx(contractName string, Abi *ABI, method string, data []byte) (map[string]interface{}) {
	var abi ABI
	
	if Abi == nil {
	    return nil
	}
	
	abi = *Abi

	r := bytes.NewReader(data)
	mapResult := DecodeAbiEx(contractName, method, r, abi, "", "", nil)
	if mapResult == nil {
               return nil
        }

	return mapResult
}
