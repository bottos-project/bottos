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
 * file description:  abi type definition
 * @Author: Gong Zibin
 * @Date:   2018-08-18
 * @Last Modified by:
 * @Last Modified time:
 */

package abi

import (
	"encoding/json"
	"github.com/bottos-project/bottos/contract/abi/fieldmap"
)

type ABIDefFeild struct {
	Name string
	Type  string
}

//ABIDefStruct
type ABIDefStruct struct {
	Name   string    `json:"name"`
	Base   string    `json:"base"`
	Fields *fieldmap.FeildMap `json:"fields"`
}

//ABIDefAction abi Method
type ABIDefAction struct {
	Name string `json:"action_name"`
	Type string `json:"type"`
}

//ABIDef for ABI definition
type ABIDef struct {
	Types   []interface{}  `json:"types"`
	Structs []ABIDefStruct `json:"structs"`
	Methods []ABIDefAction `json:"actions"`
	Tables  []interface{}  `json:"tables"`
}

func NewABIStruct(name, base string, fields ...ABIDefFeild) ABIDefStruct {
	s := ABIDefStruct{
		Name: name,
		Base: base,
		Fields: fieldmap.New(),
	}

	for _, f := range fields {
		s.Fields.Set(f.Name, f.Type)
	}

	return s
}

func NewABIMethod(name string, typ string) ABIDefAction {
	m := ABIDefAction{
		Name: name,
		Type: typ,
	}

	return m
}

func (def ABIDef) ToJson(beautify bool) string {
	data, _ := json.Marshal(def)
	if beautify {
		return jsonFormat(data)
	}
	return string(data)
}
