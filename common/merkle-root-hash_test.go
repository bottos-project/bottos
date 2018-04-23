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
 * file description:  general Hash type
 * @Author: Gong Zibin
 * @Date:   2017-12-12
 * @Last Modified by:
 * @Last Modified time:
 */

package common

import (
	"fmt"
	"testing"
)

func TestMerkleRootHash_Odd(t *testing.T) {
	var hs []Hash
	hs = append(hs, Sha256([]byte("1")))
	hs = append(hs, Sha256([]byte("2")))
	hs = append(hs, Sha256([]byte("3")))
	hs = append(hs, Sha256([]byte("4")))
	hs = append(hs, Sha256([]byte("5")))

	root := ComputeMerkleRootHash(hs)
	fmt.Printf("root hash: %x\n", root)
}

func TestMerkleRootHash_Even(t *testing.T) {
	var hs []Hash
	hs = append(hs, Sha256([]byte("1")))
	hs = append(hs, Sha256([]byte("2")))
	hs = append(hs, Sha256([]byte("3")))
	hs = append(hs, Sha256([]byte("4")))

	root := ComputeMerkleRootHash(hs)
	fmt.Printf("root hash: %x\n", root)
}


func TestMerkleRootHash_NULL(t *testing.T) {
	var hs []Hash

	root := ComputeMerkleRootHash(hs)
	fmt.Printf("root hash: %x\n", root)
}