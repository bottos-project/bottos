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
	"github.com/bottos-project/bottos/contract/msgpack"
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
	
	for i := range abi.Structs {
		var s ABIStruct
		for k, v := range abis.Structs[i].Fields {
			
			s = ABIStruct{Name: abis.Structs[i].Name, Fields: New()}
			s.Fields.Set(k, v)
		}
		abi.Structs = append(abi.Structs, s)
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

//getAbiFieldsByAbiEx function
func getAbiFieldsByAbiEx(contractname string, method string, abi ABI, subStructName string) *FeildMap {
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
		return fmt.Errorf("EncodeAbiEx: count is 0!", count)
	}

	msgpack.PackArraySize(w, uint16(count))

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
					msgpack.PackStr16(w, val.(string))
				case "uint8":
					msgpack.PackUint8(w, val.(uint8))
				case "uint16":
					msgpack.PackUint16(w, val.(uint16))
				case "uint32":
					msgpack.PackUint32(w, val.(uint32))
				case "uint64":
					msgpack.PackUint64(w, val.(uint64))
				case "bytes":
					msgpack.PackBin16(w, val.([]byte))
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
		_, errs = msgpack.UnpackArraySize(r)
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
					val, err := msgpack.UnpackStr16(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "uint8":
					val, err := msgpack.UnpackUint8(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "uint16":
					val, err := msgpack.UnpackUint16(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "uint32":
					val, err := msgpack.UnpackUint32(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "uint64":
					val, err := msgpack.UnpackUint64(r)
					if err != nil {
						return nil
					}
					Setmapval(mapResult, abiValKey, val)
					i++
				case "bytes":
					val, err := msgpack.UnpackBin16(r)
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

var a  *ABI

func GetAbi() *ABI {
	if a != nil {
		return a
	}
	
	a = CreateNativeContractABI()
	
	return a
}

func CreateNativeContractABI() *ABI {

	a = &ABI{}
	a.Actions = append(a.Actions, ABIAction{ActionName: "newaccount", Type: "NewAccount"})
	a.Actions = append(a.Actions, ABIAction{ActionName: "transfer", Type: "Transfer"})
	a.Actions = append(a.Actions, ABIAction{ActionName: "setdelegate", Type: "SetDelegate"})
	a.Actions = append(a.Actions, ABIAction{ActionName: "grantcredit", Type: "GrantCredit"})
	a.Actions = append(a.Actions, ABIAction{ActionName: "cancelcredit", Type: "CancelCredit"})
	a.Actions = append(a.Actions, ABIAction{ActionName: "transferfrom", Type: "TransferFrom"})
	a.Actions = append(a.Actions, ABIAction{ActionName: "deploycode", Type: "DeployCode"})
	a.Actions = append(a.Actions, ABIAction{ActionName: "deployabi", Type: "DeployABI"})

	s := ABIStruct{Name: "NewAccount", Fields: New()}
	s.Fields.Set("name", "string")
	s.Fields.Set("pubkey", "string")
	a.Structs = append(a.Structs, s)
	s = ABIStruct{Name: "Transfer", Fields: New()}
	s.Fields.Set("from", "string")
	s.Fields.Set("to", "string")
	s.Fields.Set("value", "uint64")
	a.Structs = append(a.Structs, s)
	s = ABIStruct{Name: "SetDelegate", Fields: New()}
	s.Fields.Set("name", "string")
	s.Fields.Set("pubkey", "string")
	a.Structs = append(a.Structs, s)
	s = ABIStruct{Name: "GrantCredit", Fields: New()}
	s.Fields.Set("name", "string")
	s.Fields.Set("spender", "string")
	s.Fields.Set("limit", "uint64")
	a.Structs = append(a.Structs, s)
	s = ABIStruct{Name: "CancelCredit", Fields: New()}
	s.Fields.Set("name", "string")
	s.Fields.Set("spender", "string")
	a.Structs = append(a.Structs, s)
	s = ABIStruct{Name: "TransferFrom", Fields: New()}
	s.Fields.Set("from", "string")
	s.Fields.Set("to", "string")
	s.Fields.Set("value", "uint64")
	a.Structs = append(a.Structs, s)
	s = ABIStruct{Name: "DeployCode", Fields: New()}
	s.Fields.Set("contract", "string")
	s.Fields.Set("vm_type", "uint8")
	s.Fields.Set("vm_version", "uint8")
	s.Fields.Set("contract_code", "bytes")
	a.Structs = append(a.Structs, s)
	s = ABIStruct{Name: "DeployABI", Fields: New()}
	s.Fields.Set("contract", "string")
	s.Fields.Set("contract_abi", "bytes")
	a.Structs = append(a.Structs, s)

	return a
}
