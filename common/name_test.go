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
		"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameEncoding(t *testing.T) {
	name, err := NewName("bottos")
	assert.Nil(t, err)
	assert.Equal(t, name.Bytes(), fromHex("000000000000000000000002d875d61c"))
	raw := name.ToString()
	assert.Equal(t, raw, "bottos")
}

func TestErrorEncoding(t *testing.T) {
	_, err := NewName("O")
	assert.NotNil(t, err)
}

func TestErrorEncoding1(t *testing.T) {
	_, err := NewName("L")
	assert.NotNil(t, err)
}

func TestLongName(t *testing.T) {
	_, err := NewName("bottos-1234-bottos-2345-bottos")
	assert.NotNil(t, err)
}

func fromHex(str string) []byte {
	b, err := hex.DecodeString(strings.Replace(str, " ", "", -1))
	if err != nil {
		panic(fmt.Sprintf("invalid hex string: %q", str))
	}
	return b
}
