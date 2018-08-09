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
 * file description:  memory instructions
 * @Author: Stewart Li
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package exec

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

// Type define one variable type
type Type int

const (
	// Int8 type enum
	Int8 Type = iota
	// Int16 type enum
	Int16
	// Int32 type enum
	Int32
	// Int64 type enum
	Int64
	// Float32 type enum
	Float32
	// Float64 type enum
	Float64
	// String type enum
	String
	// Struct type enum
	Struct
	// Unknown type enum
	Unknown
)

type typeInfo struct {
	Type Type
	Len  uint64
}

// ErrOutOfBoundsMemoryAccess is the error value used while trapping the VM
// when it detects an out of bounds access to the linear memory.
var ErrOutOfBoundsMemoryAccess = errors.New("exec: out of bounds memory access")

func (vm *VM) fetchBaseAddr() int {
	return int(vm.fetchUint32() + uint32(vm.popInt32()))
}

// inBounds returns true when the next vm.fetchBaseAddr() + offset
// indices are in bounds accesses to the linear memory.
func (vm *VM) inBounds(offset int) bool {
	addr := endianess.Uint32(vm.ctx.code[vm.ctx.pc:]) + uint32(vm.ctx.stack[len(vm.ctx.stack)-1])
	return int(addr)+offset < len(vm.memory)
}

// curMem returns a slice to the memeory segment pointed to by
// the current base address on the bytecode stream.
func (vm *VM) curMem() []byte {
	return vm.memory[vm.fetchBaseAddr():]
}

func (vm *VM) i32Load() {
	if !vm.inBounds(3) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushUint32(endianess.Uint32(vm.curMem()))
}

func (vm *VM) i32Load8s() {
	if !vm.inBounds(0) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushInt32(int32(int8(vm.memory[vm.fetchBaseAddr()])))
}

func (vm *VM) i32Load8u() {
	if !vm.inBounds(0) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushUint32(uint32(uint8(vm.memory[vm.fetchBaseAddr()])))
}

func (vm *VM) i32Load16s() {
	if !vm.inBounds(1) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushInt32(int32(int16(endianess.Uint16(vm.curMem()))))
}

func (vm *VM) i32Load16u() {
	if !vm.inBounds(1) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushUint32(uint32(endianess.Uint16(vm.curMem())))
}

func (vm *VM) i64Load() {
	if !vm.inBounds(7) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushUint64(endianess.Uint64(vm.curMem()))
}

func (vm *VM) i64Load8s() {
	if !vm.inBounds(0) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushInt64(int64(int8(vm.memory[vm.fetchBaseAddr()])))
}

func (vm *VM) i64Load8u() {
	if !vm.inBounds(0) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushUint64(uint64(uint8(vm.memory[vm.fetchBaseAddr()])))
}

func (vm *VM) i64Load16s() {
	if !vm.inBounds(1) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushInt64(int64(int16(endianess.Uint16(vm.curMem()))))
}

func (vm *VM) i64Load16u() {
	if !vm.inBounds(1) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushUint64(uint64(endianess.Uint16(vm.curMem())))
}

func (vm *VM) i64Load32s() {
	if !vm.inBounds(3) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushInt64(int64(int32(endianess.Uint32(vm.curMem()))))
}

func (vm *VM) i64Load32u() {
	if !vm.inBounds(3) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushUint64(uint64(endianess.Uint32(vm.curMem())))
}

func (vm *VM) f32Store() {
	v := math.Float32bits(vm.popFloat32())
	if !vm.inBounds(3) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	endianess.PutUint32(vm.curMem(), v)
}

func (vm *VM) f32Load() {
	if !vm.inBounds(3) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushFloat32(math.Float32frombits(endianess.Uint32(vm.curMem())))
}

func (vm *VM) f64Store() {
	v := math.Float64bits(vm.popFloat64())
	if !vm.inBounds(7) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	endianess.PutUint64(vm.curMem(), v)
}

func (vm *VM) f64Load() {
	if !vm.inBounds(7) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.pushFloat64(math.Float64frombits(endianess.Uint64(vm.curMem())))
}

func (vm *VM) i32Store() {
	v := vm.popUint32()
	if !vm.inBounds(3) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	endianess.PutUint32(vm.curMem(), v)
}

func (vm *VM) i32Store8() {
	v := byte(uint8(vm.popUint32()))
	if !vm.inBounds(0) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.memory[vm.fetchBaseAddr()] = v
}

func (vm *VM) i32Store16() {
	v := uint16(vm.popUint32())
	if !vm.inBounds(1) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	endianess.PutUint16(vm.curMem(), v)
}

func (vm *VM) i64Store() {
	v := vm.popUint64()
	if !vm.inBounds(7) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	endianess.PutUint64(vm.curMem(), v)
}

func (vm *VM) i64Store8() {
	v := byte(uint8(vm.popUint64()))
	if !vm.inBounds(0) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	vm.memory[vm.fetchBaseAddr()] = v
}

func (vm *VM) i64Store16() {
	v := uint16(vm.popUint64())
	if !vm.inBounds(1) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	endianess.PutUint16(vm.curMem(), v)
}

func (vm *VM) i64Store32() {
	v := uint32(vm.popUint64())
	if !vm.inBounds(3) {
		panic(ErrOutOfBoundsMemoryAccess)
	}
	endianess.PutUint32(vm.curMem(), v)
}

func (vm *VM) currentMemory() {
	_ = vm.fetchInt8() // reserved (https://github.com/WebAssembly/design/blob/27ac254c854994103c24834a994be16f74f54186/BinaryEncoding.md#memory-related-operators-described-here)
	vm.pushInt32(int32(len(vm.memory) / wasmPageSize))
}

func (vm *VM) growMemory() {
	_ = vm.fetchInt8() // reserved (https://github.com/WebAssembly/design/blob/27ac254c854994103c24834a994be16f74f54186/BinaryEncoding.md#memory-related-operators-described-here)
	curLen := len(vm.memory) / wasmPageSize
	n := vm.popInt32()
	vm.memory = append(vm.memory, make([]byte, n*wasmPageSize)...) //auto extend range
	vm.pushInt32(int32(curLen))
}

// GetData retrieve data
func (vm *VM) GetData(pos uint64) ([]byte, error) {

	if pos < 0 || pos == uint64(math.MaxInt64) {
		return nil, ERR_OUT_BOUNDS
	}

	t, ok := vm.memType[pos] //map[uint64]*typeInfo
	if !ok {
		return nil, ERR_FINE_MAP
	}

	if pos + t.Len > uint64(len(vm.memory)) {
		return nil, ERR_OUT_BOUNDS
	}

	return vm.memory[int(pos) : pos + t.Len], nil
}

// StorageData store data
func (vm *VM) StorageData(data interface{}) (uint64, error) {

	if data == nil {
		return 0, ERR_EMPTY_INVALID_PARAM
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.String:
		return vm.storageMemory([]byte(data.(string)), String)
	case reflect.Slice:
		switch data.(type) {
		case []byte:
			return vm.storageMemory(data.([]byte), String)
		case []int:
			byteArray := make([]byte, len(data.([]int))*4)
			for i, v := range data.([]int) {
				array := make([]byte, 4)
				binary.LittleEndian.PutUint32(array, uint32(v))
				copy(byteArray[i*4:(i+1)*4], array)
			}
			return vm.storageMemory(byteArray, Int32)

		default:
			return 0, ERR_UNSUPPORT_TYPE
		}
	default:
		return 0, ERR_UNSUPPORT_TYPE
	}
}

func (vm *VM) getStoragePos(size uint64, t Type) (uint64, error) {

	if size <= 0 || vm.memory == nil {
		return 0, ERR_EMPTY_INVALID_PARAM
	}

	if vm.memPos + size > uint64(len(vm.memory)) {
		return 0, ERR_OUT_BOUNDS
	}

	newpos    := vm.memPos + 1
	vm.memPos += size

	vm.memType[uint64(newpos)] = &typeInfo{Type: t, Len: size}

	return newpos, nil
}

func (vm *VM) storageMemory(b []byte, t Type) (uint64 , error) {
	index, err := vm.getStoragePos(uint64(len(b)), t) //get new pos after storage new data
	if err != nil {
		return 0, ERR_GET_STORE_POS
	}
	copy(vm.memory[index : index + uint64(len(b))], b)
	vm.memory[index + uint64(len(b))] = 0

	return index, nil
}

func (vm *VM) registerMemory(pos int , size uint64, t Type) (int, error) {

	if pos <= 0 || size < 0 {
		return 0, ERR_EMPTY_INVALID_PARAM
	}

	vm.memType[uint64(pos)] = &typeInfo{Type: t, Len: size}
	return pos, nil
}

func (vm *VM) getDataLen(pos int) (uint64, error) {
	ti , ok := vm.memType[uint64(pos)]
	if !ok {
		return 0 , ERR_FINE_MAP
	}

	return ti.Len , nil
}