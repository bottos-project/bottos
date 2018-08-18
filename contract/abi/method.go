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
 * file description:  abi method
 * @Author: Gong Zibin
 * @Date:   2018-08-18
 * @Last Modified by:
 * @Last Modified time:
 */

package abi

import (
	"github.com/bottos-project/bottos/contract/abi/fieldmap"
)

type Method struct {
	Name    string
	Fields  *fieldmap.FeildMap
}

func (m Method) GetFieldNames() []string {
	return m.Fields.Keys()
}

func (m Method) GetFieldType(name string) (string, bool) {
	return m.Fields.GetStringVal(name)
}
