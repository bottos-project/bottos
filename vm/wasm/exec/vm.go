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

// Package exec provides functions for executing WebAssembly bytecode.
package exec

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"sync"

	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/vm/wasm/disasm"
	"github.com/bottos-project/bottos/vm/wasm/exec/internal/compile"
	"github.com/bottos-project/bottos/vm/wasm/wasm"
	ops "github.com/bottos-project/bottos/vm/wasm/wasm/operators"
)

type SrcFileType int
const (
	CPP SrcFileType = iota
	JS
	PY
)

// InvalidReturnTypeError is returned by (*VM).ExecCode when the module
// specifies an invalid return type value for the executed function.
type InvalidReturnTypeError int8

func (e InvalidReturnTypeError) Error() string {
	return fmt.Sprintf("Function has invalid return value_type: %d", int8(e))
}

// InvalidFunctionIndexError is returned by (*VM).ExecCode when the function
// index provided is invalid.
type InvalidFunctionIndexError int64

func (e InvalidFunctionIndexError) Error() string {
	return fmt.Sprintf("Invalid index to function index space: %d", int64(e))
}

type context struct {
	stack   []uint64
	locals  []uint64
	code    []byte
	pc      int64
	curFunc int64
}

// VM is the execution context for executing WebAssembly bytecode.
type VM struct {
	ctx context

	module        *wasm.Module
	globals       []uint64
	memory        []byte
	compiledFuncs []compiledFunction

	funcTable     [256]func()

	memPos        uint64
	//To avoid the too much the number of recursion execution(dep) in contract
	callDep       int
	//To limit the too much the number of new contract execution(wid) in contract
	callWid       int
	// define a map relationship between memory address and data's type
	memType       map[uint64]*typeInfo
	//define env function
	envFunc      *EnvFunc
	funcInfo      FuncInfo

	contract     *contract.Context

	vmLock       *sync.Mutex
	//the channel be used to communcate with vm_engine
	vmChannel     chan []byte

	//record sub-trx for recursive call[wid]
	subTrxLst     []*types.Transaction
	subCtnLst     []*contract.Context

	codeVersion   uint32

	//to identify the type of source file
	sourceFile    SrcFileType
}

// As per the WebAssembly spec: https://github.com/WebAssembly/design/blob/27ac254c854994103c24834a994be16f74f54186/Semantics.md#linear-memory
const wasmPageSize = 65536 // (64 KB)

var endianess = binary.LittleEndian

// NewVM creates a new VM from a given module. If the module defines a
// start function, it will be executed.
func NewVM(module *wasm.Module) (*VM, error) {

	var value interface{}
	var err   error

	var vm = &VM{
		envFunc:  NewEnvFunc(),
		memPos:   0,
		memType:  make(map[uint64]*typeInfo),
		memory:   make([]byte, wasmPageSize),
		contract: nil,
		vmLock:   new(sync.Mutex),
		callDep:  0,
		callWid:  0,
	}

	if len(module.LinearMemoryIndexSpace) <= 0 {
		return nil, ERR_INVALID_WASM
	}

	if module.Memory != nil && len(module.Memory.Entries) != 0 {
		if len(module.Memory.Entries) > 1 {
			return nil, ERR_MULTIPLE_LINEAR_MEMORIES
		}
		vm.memory = make([]byte, uint(module.Memory.Entries[0].Limits.Initial)*wasmPageSize)
	}

	indexSpaceLen := len(module.LinearMemoryIndexSpace[0])
	if copy(vm.memory, module.LinearMemoryIndexSpace[0]) == indexSpaceLen {
		vm.memPos += uint64(len(module.LinearMemoryIndexSpace[0]))
	}else{
		return nil , ERR_CREATE_VM
	}

	//it need modify if adding python or compiler change
	if module.Other == nil {
		vm.sourceFile = CPP
	} else {
		vm.sourceFile = JS
	}

	if module.Data != nil {
		for _, funcList := range module.Data.Entries {
			if value, err = module.ExecInitExpr(funcList.Offset); err != nil {
				return nil, err
			}

			index, ok := value.(int32)
			if !ok {
				return nil, ERR_DATA_INDEX
			}

			// if it contains multi-function(splited by '0')
			if bytes.Contains(funcList.Data, []byte{byte(0)}) != true {
				vm.memType[uint64(index)] = &typeInfo{Type: String, Len: uint64(len(funcList.Data))}
			} else {
				var idx = int(index)
				funcArray := bytes.Split(funcList.Data, []byte{byte(0)})
				for _, function := range funcArray {
					vm.memType[uint64(idx)] = &typeInfo{Type: String, Len: uint64(len(function) + 1)}
					idx += len(function) + 1
				}
			}
		}
	} else {
		vm.memPos = uint64(len(vm.memory) / 2)
	}

	vm.compiledFuncs = make([]compiledFunction, len(module.FunctionIndexSpace))
	vm.globals       = make([]uint64, len(module.GlobalIndexSpace))
	vm.newFuncTable()
	vm.module = module

	for i, fn := range module.FunctionIndexSpace {
		disassembly, err := disasm.Disassemble(fn, module)
		if err != nil {
			return nil, err
		}

		totalLocalVars := 0
		totalLocalVars += len(fn.Sig.ParamTypes)
		for _, entry := range fn.Body.Locals {
			totalLocalVars += int(entry.Count)
		}

		code, table := compile.Compile(disassembly.Code)
		vm.compiledFuncs[i] = compiledFunction{
			code:           code,
			branchTables:   table,
			maxDepth:       disassembly.MaxDepth,
			totalLocalVars: totalLocalVars,
			args:           len(fn.Sig.ParamTypes),
			returns:        len(fn.Sig.ReturnTypes) != 0,
			funcProp:       fn,
		}
	}

	for i, global := range module.GlobalIndexSpace {
		val, err := module.ExecInitExpr(global.Init)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case int32:
			vm.globals[i] = uint64(v)
		case int64:
			vm.globals[i] = uint64(v)
		case float32:
			vm.globals[i] = uint64(math.Float32bits(v))
		case float64:
			vm.globals[i] = uint64(math.Float64bits(v))
		}
	}

	if module.Start != nil {
		_, err := vm.ExecCode(int64(module.Start.Index))
		if err != nil {
			return nil, err
		}
	}

	return vm, nil
}

// Memory returns the linear memory space for the VM.
func (vm *VM) Memory() []byte {
	return vm.memory
}

func (vm *VM) pushBool(v bool) {
	if v {
		vm.pushUint64(1)
	} else {
		vm.pushUint64(0)
	}
}

func (vm *VM) fetchBool() bool {
	return vm.fetchInt8() != 0
}

func (vm *VM) fetchInt8() int8 {
	i := int8(vm.ctx.code[vm.ctx.pc])
	vm.ctx.pc++
	return i
}

func (vm *VM) fetchUint32() uint32 {
	v := endianess.Uint32(vm.ctx.code[vm.ctx.pc:])
	vm.ctx.pc += 4
	return v
}

func (vm *VM) fetchInt32() int32 {
	return int32(vm.fetchUint32())
}

func (vm *VM) fetchFloat32() float32 {
	return math.Float32frombits(vm.fetchUint32())
}

func (vm *VM) fetchUint64() uint64 {
	v := endianess.Uint64(vm.ctx.code[vm.ctx.pc:])
	vm.ctx.pc += 8
	return v
}

func (vm *VM) fetchInt64() int64 {
	return int64(vm.fetchUint64())
}

func (vm *VM) fetchFloat64() float64 {
	return math.Float64frombits(vm.fetchUint64())
}

func (vm *VM) popUint64() uint64 {
	i := vm.ctx.stack[len(vm.ctx.stack)-1]
	vm.ctx.stack = vm.ctx.stack[:len(vm.ctx.stack)-1]
	return i
}

func (vm *VM) popInt64() int64 {
	return int64(vm.popUint64())
}

func (vm *VM) popFloat64() float64 {
	return math.Float64frombits(vm.popUint64())
}

func (vm *VM) popUint32() uint32 {
	return uint32(vm.popUint64())
}

func (vm *VM) popInt32() int32 {
	return int32(vm.popUint32())
}

func (vm *VM) popFloat32() float32 {
	return math.Float32frombits(vm.popUint32())
}

func (vm *VM) pushUint64(i uint64) {
	vm.ctx.stack = append(vm.ctx.stack, i)
}

func (vm *VM) pushInt64(i int64) {
	vm.pushUint64(uint64(i))
}

func (vm *VM) pushFloat64(f float64) {
	vm.pushUint64(math.Float64bits(f))
}

func (vm *VM) pushUint32(i uint32) {
	vm.pushUint64(uint64(i))
}

func (vm *VM) pushInt32(i int32) {
	vm.pushUint64(uint64(i))
}

func (vm *VM) pushFloat32(f float32) {
	vm.pushUint32(math.Float32bits(f))
}

// ExecCode calls the function with the given index and arguments.
// fnIndex should be a valid index into the function index space of
// the VM's module.
func (vm *VM) ExecCode(fnIndex int64, args ...uint64) (interface{}, error) {
	if int(fnIndex) > len(vm.compiledFuncs) {
		return nil, InvalidFunctionIndexError(fnIndex)
	}

	if len(vm.module.GetFunction(int(fnIndex)).Sig.ParamTypes) != len(args) {
		return nil, ERR_INVALID_ARGUMENT_COUNT
	}

	compiled      := vm.compiledFuncs[fnIndex]
	vm.ctx.stack   = make([]uint64, 0, compiled.maxDepth)
	vm.ctx.locals  = make([]uint64, compiled.totalLocalVars) // number of local variables used by the function
	vm.ctx.pc      = 0
	vm.ctx.code    = compiled.code
	vm.ctx.curFunc = fnIndex

	for i, arg := range args {
		vm.ctx.locals[i] = arg
	}

	var rtrn interface{}

	res := vm.execCode(compiled)
	if compiled.returns {
		rtrnType := vm.module.GetFunction(int(fnIndex)).Sig.ReturnTypes[0]
		switch rtrnType {
		case wasm.ValueTypeI32:
			rtrn = uint32(res)
		case wasm.ValueTypeI64:
			rtrn = uint64(res)
		case wasm.ValueTypeF32:
			rtrn = math.Float32frombits(uint32(res))
		case wasm.ValueTypeF64:
			rtrn = math.Float64frombits(res)
		default:
			return nil, InvalidReturnTypeError(rtrnType)
		}
	}

	return rtrn, nil
}

func (vm *VM) execCode(compiled compiledFunction) uint64 {
	if compiled.funcProp.EnvFunc == true {
		err := vm.ExecEnvFunc(compiled)
		if err != nil {
			return uint64(VM_ERROR_FAIL_EXECUTE_ENVFUNC)
		}
	}
outer:
	for int(vm.ctx.pc) < len(vm.ctx.code) {
		op := vm.ctx.code[vm.ctx.pc]
		vm.ctx.pc++
		switch op {
		case ops.Return:
			break outer
		case compile.OpJmp:
			vm.ctx.pc = vm.fetchInt64()
			continue
		case compile.OpJmpZ:
			target := vm.fetchInt64()
			if vm.popUint32() == 0 {
				vm.ctx.pc = target
				continue
			}
		case compile.OpJmpNz:
			target      := vm.fetchInt64()
			preserveTop := vm.fetchBool()
			discard := vm.fetchInt64()
			if vm.popUint32() != 0 {
				vm.ctx.pc = target
				var top uint64
				if preserveTop {
					top = vm.ctx.stack[len(vm.ctx.stack)-1]
				}
				vm.ctx.stack = vm.ctx.stack[:len(vm.ctx.stack)-int(discard)]
				if preserveTop {
					vm.pushUint64(top)
				}
				continue
			}
		case ops.BrTable:
			index := vm.fetchInt64()
			label := vm.popInt32()
			table := vm.compiledFuncs[vm.ctx.curFunc].branchTables[index]
			var target compile.Target
			if label >= 0 && label < int32(len(table.Targets)) {
				target = table.Targets[int32(label)]
			} else {
				target = table.DefaultTarget
			}

			if target.Return {
				break outer
			}
			vm.ctx.pc = target.Addr
			var top uint64
			if target.PreserveTop {
				top = vm.ctx.stack[len(vm.ctx.stack)-1]
			}
			vm.ctx.stack = vm.ctx.stack[:len(vm.ctx.stack)-int(target.Discard)]
			if target.PreserveTop {
				vm.pushUint64(top)
			}
			continue
		case compile.OpDiscard:
			place := vm.fetchInt64()
			if len(vm.ctx.stack)-int(place) > 0 {
				vm.ctx.stack = vm.ctx.stack[:len(vm.ctx.stack)-int(place)]
			}

		case compile.OpDiscardPreserveTop:
			top   := vm.ctx.stack[len(vm.ctx.stack)-1]
			place := vm.fetchInt64()
			vm.ctx.stack = vm.ctx.stack[:len(vm.ctx.stack)-int(place)]
			vm.pushUint64(top)
		default:
			vm.funcTable[op]()
		}
	}

	if compiled.returns && len(vm.ctx.stack) >= 1 {
		return vm.ctx.stack[len(vm.ctx.stack)-1]
	}

	return uint64(VM_NOERROR)
}

// GetMemory get memory
func (vm *VM) GetMemory() []byte {
	return vm.memory
}

// GetFuncParams get param
func (vm *VM) GetFuncParams() []uint64 {
	envFunc := vm.envFunc
	params  := envFunc.envFuncParam

	return params
}

// ExecEnvFunc exec function
func (vm *VM) ExecEnvFunc(compiled compiledFunction) error {

	vm.envFunc.envFuncParam = vm.ctx.locals
	vm.envFunc.envFuncCtx   = vm.ctx
	oldCtx                 := vm.ctx

	if compiled.returns {
		vm.envFunc.envFuncRtn = true
	} else {
		vm.envFunc.envFuncRtn = false
	}

	fc, ok := vm.envFunc.envFuncMap[compiled.funcProp.Method] //get env function
	if !ok {
		fmt.Println("*ERROR* Failed to search the method: " + compiled.funcProp.Method)
		return ERR_FIND_VM_METHOD
	}

	_, err := fc(vm)
	if err != nil {
		vm.ctx = oldCtx
		if compiled.returns {
			vm.pushUint64(0)
		}

		return ERR_CALL_ENV_METHOD
	}

	return nil
}

// GetMsgBytes get message bytes
func (vm *VM) GetMsgBytes() ([]byte, error) {

	bytesbuf := bytes.NewBuffer(nil)
	return bytesbuf.Bytes(), nil
}

// SetContract set contract param
func (vm *VM) SetContract(contract *contract.Context) error {

	if contract == nil {
		return ERR_EMPTY_INVALID_PARAM
	}

	vm.contract = contract
	return nil
}

// GetContract get contract
func (vm *VM) GetContract() *contract.Context {
	return vm.contract
}

// SetChannel set channel
func (vm *VM) SetChannel(channel chan []byte) error {
	vm.vmChannel = channel
	return nil
}

func (vm *VM) RecoverContext() bool {
	if vm.envFunc != nil {
		vm.ctx = vm.envFunc.envFuncCtx
	}

	return true
}