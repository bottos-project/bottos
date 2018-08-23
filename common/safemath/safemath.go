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
	"math/big"

	"github.com/bottos-project/bottos/common"
)

//Uint64Add safe add
func Uint64Add(a uint64, b uint64) (uint64, error) {
	var c uint64
	c = a + b
	if c >= a {
		return c, nil
	}

	return 0, errors.New("uint64 overflow")
}

//Uint64Sub safe sub
func Uint64Sub(a uint64, b uint64) (uint64, error) {
	var c uint64
	if b <= a {
		c = a - b
		return c, nil
	}

	return 0, errors.New("uint64 overflow")
}

//Uint64Mul safe mul
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

//U256IntAdd safe add
func U256Add(result *big.Int, a *big.Int, b *big.Int) (*big.Int, error) {
	//var c uint64
	result = result.Add(a, b)
	if 1 != result.Cmp(common.MaxUint256()) {
		return result, nil
	}

	return result, errors.New("bigInt overflow1")
}

//U256Sub safe sub
func U256Sub(result *big.Int, a *big.Int, b *big.Int) (*big.Int, error) {
	//var c big.Int
	if 1 != b.Cmp(a) {
		result = result.Sub(a, b)
		return result, nil
	}

	return result, errors.New("bigInt overflow2")
}

//U256Mul safe mul
func U256Mul(result *big.Int, a *big.Int, b *big.Int) (*big.Int, error) {
	//var c uint64
	result = result.Mul(a, b)

	if 1 == result.Cmp(common.MaxUint256()) {
		return result, errors.New("bigInt overflow3")
	}

	if 0 != a.Cmp(big.NewInt(0)) {			
		d := big.NewInt(0)
		//d = c / a
		d = d.Div(result, a)
		if 0 != b.Cmp(d) {
			return result, errors.New("bigInt overflow4")
		}
	}

	return result, nil
}