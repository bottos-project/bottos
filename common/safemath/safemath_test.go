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
 * file description:  safe math
 * @Author: Gong Zibin
 * @Date:   2017-01-20
 * @Last Modified by:
 * @Last Modified time:
 */

package safemath

import (
	"fmt"
	"math"
	"testing"
	//"fmt"

	"github.com/stretchr/testify/assert"
)

func TestOverFlow(t *testing.T) {
	var u64 uint64
	var u64max uint64
	var u64zero uint64

	u64max = math.MaxUint64
	u64zero = uint64(0)

	u64 = u64max
	assert.Equal(t, u64max, u64)

	// overflow
	u64 = u64 + 1
	assert.Equal(t, u64zero, u64)

	// overflow
	u64 = 0
	u64 = u64 - 1
	assert.Equal(t, u64max, u64)
}

func TestSafeMath(t *testing.T) {
	var u64 uint64
	var u64max uint64
	var u64zero uint64

	u64max = math.MaxUint64
	u64zero = uint64(0)

	u64 = u64max
	assert.Equal(t, u64max, u64)

	// overflow
	var a, b, c uint64
	a = u64max
	b = uint64(1)
	c, err := Uint64Add(a, b)
	assert.NotNil(t, err)

	// overflow
	u64 = 0
	u64 = u64 - 1
	assert.Equal(t, u64max, u64)
	a = u64zero
	b = uint64(1)
	c, err = Uint64Sub(a, b)
	assert.NotNil(t, err)

	// normal add
	a = 9999999
	b = 88888888
	c, err = Uint64Add(a, b)
	assert.Equal(t, a+b, c)

	// normal sub
	a = 99999999
	b = 88888888
	c, err = Uint64Sub(a, b)
	fmt.Printf(err)
	assert.Equal(t, a-b, c)
}
