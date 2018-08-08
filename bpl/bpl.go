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
 * file description:  BPL, Bottos Pack Library
 * @Author: Gong Zibin
 * @Date:   2018-08-02
 * @Last Modified by:
 * @Last Modified time:
 */

package bpl

import (
	"bytes"
)

//Marshal is to serialize the message
func Marshal(v interface{}) ([]byte, error) {
	writer := &bytes.Buffer{}
	err := Encode(v, writer)
	if err != nil {
		return []byte{}, err
	}
	return writer.Bytes(), nil
}

//Unmarshal is to unserialize the message
func Unmarshal(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	err := Decode(r, v)
	return err
}
