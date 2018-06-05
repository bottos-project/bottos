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
 * file description:  main functions
 * @Author: Stewart Li
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package exec

import (
	"fmt"
	"math"
	"reflect"

	"github.com/bottos-project/bottos/vm/wasm/exec/internal/compile"
	"github.com/bottos-project/bottos/vm/wasm/wasm"
)

type function interface {
	call(vm *VM, index int64)
}

type compiledFunction struct {
	code           []byte //it means the internal call order for a method
	branchTables   []*compile.BranchTable
	maxDepth       int           // maximum stack depth reached while executing the function body
	totalLocalVars int           // number of local variables used by the function
	args           int           // number of arguments the function accepts
	returns        bool          // whether the function returns a value
	funcProp       wasm.Function //record function's properties
}

type goFunction struct {
	val reflect.Value
	typ reflect.Type
}

func (fn goFunction) call(vm *VM, index int64) {
	numIn := fn.typ.NumIn()
	args := make([]reflect.Value, numIn)

	for i := numIn - 1; i >= 0; i-- {
		val := reflect.New(fn.typ.In(i)).Elem()
		raw := vm.popUint64()
		kind := fn.typ.In(i).Kind()

		switch kind {
		case reflect.Float64, reflect.Float32:
			val.SetFloat(math.Float64frombits(raw))
		case reflect.Uint32, reflect.Uint64:
			val.SetUint(raw)
		case reflect.Int32, reflect.Int64:
			val.SetUint(raw)
		default:
			panic(fmt.Sprintf("exec: args %d invalid kind=%v", i, kind))
		}

		args[i] = val
	}

	rtrns := fn.val.Call(args)
	for i, out := range rtrns {
		kind := out.Kind()
		switch kind {
		case reflect.Float64, reflect.Float32:
			vm.pushFloat64(out.Float())
		case reflect.Uint32, reflect.Uint64:
			vm.pushUint64(out.Uint())
		case reflect.Int32, reflect.Int64:
			vm.pushInt64(out.Int())
		default:
			panic(fmt.Sprintf("exec: return value %d invalid kind=%v", i, kind))
		}
	}
}

func (compiled compiledFunction) call(vm *VM, index int64) {
	newStack := make([]uint64, compiled.maxDepth)
	locals := make([]uint64, compiled.totalLocalVars)

	for i := compiled.args - 1; i >= 0; i-- {
		locals[i] = vm.popUint64()
	}

	//save execution context
	prevCtxt := vm.ctx

	vm.ctx = context{
		stack:   newStack,
		locals:  locals,
		code:    compiled.code,
		pc:      0,
		curFunc: index,
	}

	rtrn := vm.execCode(compiled)

	//restore execution context
	vm.ctx = prevCtxt

	if compiled.returns {
		vm.pushUint64(rtrn)
	}
}
