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
	"errors"
)

func Uint64Add(a uint64, b uint64) (uint64, error) {
	var c uint64
	c = a + b
	if c >= a {
		return c, nil
	} else {
		return 0, errors.New("uint64 overflow")
	}
}

func Uint64Sub(a uint64, b uint64) (uint64, error) {
	var c uint64
	if b <= a {
		c = a - b
		return c, nil
	} else {
		return 0, errors.New("uint64 overflow")
	}
}

func Uint64Mul(a uint64, b uint64) (uint64, error) {
	var c uint64
	c = a * b
	if a != 0 {
		var d uint64
		d = c / a
		if b != d {
			return 0, errors.New("uint64 overflow")
		}
	}

	return c, nil
}
