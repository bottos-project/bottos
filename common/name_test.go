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
 * file description:  Name test
 * @Author: Gong Zibin
 * @Date:   2018-08-08
 * @Last Modified by:
 * @Last Modified time:
 */

package common

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameEncoding(t *testing.T) {
	decoded, err := NewName("bottos")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	assert.Equal(t, decoded.Bytes(), fromHex("000000000000000000000008aecf3bb2"))
	encoded := decoded.ToString()
	assert.Equal(t, encoded, "bottos")
}

func TestErrorEncoding(t *testing.T) {
	_, err := NewName("O")
	assert.NotNil(t, err)
}

func TestErrorEncoding1(t *testing.T) {
	_, err := NewName("l")
	assert.NotNil(t, err)
}

func fromHex(str string) []byte {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		panic(fmt.Sprintf("invalid hex string: %q", str))
	}
	return b
}
