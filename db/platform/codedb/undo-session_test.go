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
 * file description: database interface
 * @Author: May Luo
 * @Date:   2017-12-04
 * @Last Modified by:
 * @Last Modified time:
 */
package codedb

import (
	"fmt"
	"testing"

	"github.com/bottos-project/bottos/config"
)

func TestUndoSession_Rollback(t *testing.T) {
	fmt.Println("abcdo")
	ins, _ := NewMultindexDB("./file3")
	ins.CallAddObject("abcundo")
	ins.CallAddObject("elfundo")
	ins.CallAddObject("lllundo")
	ins.CallAddObject("eeeundo")
	ins.CallSetRevision(1)
	ins.CallSetObject("abcundo", "111", "222")
	ins.CallSetObject("elfundo", "112", "elf222")
	rl := ins.CallGetRevision()
	fmt.Println(rl)
	session := ins.CallBeginUndo(config.PRIMARY_TRX_SESSION)
	fmt.Println(session)
	fmt.Println("111111111")
	ins.CallSetObject("lllundo", "lbc", "222")
	ins.CallSetObject("eeeundo", "cbl", "elf222")
	fmt.Println("2222222222")
	session.rollback()
	val, err := ins.CallGetObject("lllundo", "lbc")
	fmt.Println("value", val, err)
	val, err = ins.CallGetObject("eeeundo", "cbl")
	fmt.Println("value", val, err)

	val, err = ins.CallGetObject("abcundo", "111")
	fmt.Println("value", val, err)
}

func TestUndoSession_Reset(t *testing.T) {
	fmt.Println("abcdo")
	ins, _ := NewMultindexDB("./file3")
	ins.CallAddObject("abcundo")
	ins.CallAddObject("elfundo")
	ins.CallAddObject("lllundo")
	ins.CallAddObject("eeeundo")
	ins.CallSetRevision(1)
	ins.CallSetObject("abcundo", "111", "222")
	ins.CallSetObject("elfundo", "112", "elf222")
	rl := ins.CallGetRevision()
	fmt.Println(rl)
	session := ins.CallBeginUndo(config.PRIMARY_TRX_SESSION)
	fmt.Println(session)
	fmt.Println("111111111")
	ins.CallSetObject("lllundo", "lbc", "222")
	ins.CallSetObject("eeeundo", "cbl", "elf222")
	fmt.Println("2222222222")
	ins.CallResetSession()
	ses := ins.CallGetSession()
	fmt.Println("reset", ses)
	val, err := ins.CallGetObject("lllundo", "lbc")
	fmt.Println("value", val, err)
	val, err = ins.CallGetObject("eeeundo", "cbl")
	fmt.Println("value", val, err)

	val, err = ins.CallGetObject("abcundo", "111")
	fmt.Println("value", val, err)
}
