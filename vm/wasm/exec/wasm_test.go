// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

// This program is free software: you can distribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Bottos.  If not, see <http://www.gnu.org/licenses/>.

// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * file description:  wasm test suite
 * @Author: Stewart Li
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */
package exec

import (
	"fmt"
	"testing"

	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/contract/msgpack"
)

//the case is to test crx recursive call
func TestWasmRecursiveCall(t *testing.T) {

	type transferparam struct {
		To     string
		Amount uint32
	}

	param := transferparam{
		To:     "stewart",
		Amount: 1233,
	}

	bf, err := msgpack.Marshal(param)
	fmt.Println(" TestWasmRecursiveCall() bf = ", bf, " , err = ", err)

	trx := &types.Transaction{
		Version:     1,
		CursorNum:   1,
		CursorLabel: 1,
		Lifetime:    1,
		Sender:      "bottos",
		Contract:    "usermng",
		Method:      "reguser",
		Param:       bf,
		SigAlg:      1,
		Signature:   []byte{},
	}

	ctx := &contract.Context{Trx: trx}

	res, err := GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println("*ERROR* fail to execute start !!!")
		fmt.Println("err = ", err)
		return
	}

	fmt.Println("<================ *SUCCESS* res = ", res, " , err = ", err, " ================>")

	res, err = GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println("*ERROR* fail to execute start !!!")
		fmt.Println("err = ", err)
		return
	}

	fmt.Println("*SUCCESS* res = ", res, " , err = ", err)
}
