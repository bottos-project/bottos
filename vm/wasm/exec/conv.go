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
 * file description:  convert variable
 * @Author: Stewart Li
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package exec

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"fmt"
)

func (vm *VM) i32Wrapi64() {
	vm.pushUint32(uint32(vm.popUint64()))
}

func (vm *VM) i32TruncSF32() {
	vm.pushInt32(int32(math.Trunc(float64(vm.popFloat32()))))
}

func (vm *VM) i32TruncUF32() {
	vm.pushUint32(uint32(math.Trunc(float64(vm.popFloat32()))))
}

func (vm *VM) i32TruncSF64() {
	vm.pushInt32(int32(math.Trunc(vm.popFloat64())))
}

func (vm *VM) i32TruncUF64() {
	vm.pushUint32(uint32(math.Trunc(vm.popFloat64())))
}

func (vm *VM) i64ExtendSI32() {
	vm.pushInt64(int64(vm.popInt32()))
}

func (vm *VM) i64ExtendUI32() {
	vm.pushUint64(uint64(vm.popUint32()))
}

func (vm *VM) i64TruncSF32() {
	vm.pushInt64(int64(math.Trunc(float64(vm.popFloat32()))))
}

func (vm *VM) i64TruncUF32() {
	vm.pushUint64(uint64(math.Trunc(float64(vm.popFloat32()))))
}

func (vm *VM) i64TruncSF64() {
	vm.pushInt64(int64(math.Trunc(vm.popFloat64())))
}

func (vm *VM) i64TruncUF64() {
	vm.pushUint64(uint64(math.Trunc(vm.popFloat64())))
}

func (vm *VM) f32ConvertSI32() {
	vm.pushFloat32(float32(vm.popInt32()))
}

func (vm *VM) f32ConvertUI32() {
	vm.pushFloat32(float32(vm.popUint32()))
}

func (vm *VM) f32ConvertSI64() {
	vm.pushFloat32(float32(vm.popInt64()))
}

func (vm *VM) f32ConvertUI64() {
	vm.pushFloat32(float32(vm.popUint64()))
}

func (vm *VM) f32DemoteF64() {
	vm.pushFloat32(float32(vm.popFloat64()))
}

func (vm *VM) f64ConvertSI32() {
	vm.pushFloat64(float64(vm.popInt32()))
}

func (vm *VM) f64ConvertUI32() {
	vm.pushFloat64(float64(vm.popUint32()))
}

func (vm *VM) f64ConvertSI64() {
	vm.pushFloat64(float64(vm.popInt64()))
}

func (vm *VM) f64ConvertUI64() {
	vm.pushFloat64(float64(vm.popUint64()))
}

func (vm *VM) f64PromoteF32() {
	vm.pushFloat64(float64(vm.popFloat32()))
}

// BytesToString convert bytes to string
func BytesToString(bytes []byte) string {

	for i, b := range bytes {
		if b == 0 {
			return string(bytes[:i])
		}
	}
	return string(bytes)
}

// F32ToBytes convert float32 to bytes
func F32ToBytes(f32 float32) []byte {
	bytes := make([]byte, 4)
	bits := math.Float32bits(f32)

	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

// BytesToF32 convert bytes to float32
func BytesToF32(b []byte) float32 {
	f32 := math.Float32frombits(binary.LittleEndian.Uint32(b))
	return f32
}

// F64ToBytes convert float64 to bytes
func F64ToBytes(f64 float64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(f64))
	return bytes
}

// BytesToF64 convert bytes to float64
func BytesToF64(b []byte) float64 {
	f64 := math.Float64frombits(binary.LittleEndian.Uint64(b))
	return f64
}

// I32ToBytes convert int32 to bytes
func I32ToBytes(i32 uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, i32)
	return bytesBuffer.Bytes()
}

// I64ToBytes convert int64 to bytes
func I64ToBytes(i64 uint64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, i64)
	return bytesBuffer.Bytes()
}

// ByteToFloat64 convert byte to float64
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

func (vm *VM) StrLen(pos uint64) uint64 {
	if pos <= 0 {
		return 0
	}

	if vm.sourceFile == JS {
		//if source file is js file , the byte array's first element will be array's length
		return uint64(vm.memory[pos])
	}

	var i uint64 = pos
	for {
		if vm.memory[i] == 0 {
			break
		}
		i++
	}

	return (i - pos)
}

//for js
func ConvertStr (param []byte , pos uint64 , length uint64) ([]byte , error) {
	var i   uint64 = pos + 4
	var cnt uint64 = 0
	var res []byte
	if length == 0 {
		return nil , ERR_EMPTY_INVALID_PARAM
	}

	for {
		//fmt.Println("param[i]: ",param[i]," , i: ",i)
		res = append(res , param[i])
		i   += 2
		cnt += 1
		if cnt == length {
			break
		}
	}

	return res , nil
}

func UnConvertStr(param []byte , pos uint64 , len uint64) ([]byte , error) {



	return nil,nil
}

//for common
func Convert(vm *VM , pos uint64 , length uint64) ([]byte , error) {
	if vm.memory == nil || pos < 0 || length < 0 {
		return nil , ERR_EMPTY_INVALID_PARAM
	}

	vmLen := uint64(len(vm.memory))
	if pos >= vmLen || pos + length >= vmLen {
		return nil , ERR_EMPTY_INVALID_PARAM
	}

	var value []byte
	var err   error

	if vm.sourceFile == CPP {
		value = vm.memory[pos : pos + length]
	} else if vm.sourceFile == JS {
		value , err = ConvertStr(vm.memory , pos , length)
		if err != nil {
			return nil , err
		}
	} else if vm.sourceFile == PY {
		//it will be implemented in the feature
	}

	return value , nil
}

func RemoveElement(slice []byte, start, end uint64) []byte {
	return append(slice[:start], slice[end:]...)
}

func ByteArrToNum(numArray []byte) (byte , error) {
	num ,err := strconv.Atoi(string(numArray))
	if err != nil{
		fmt.Println("*ERROR* Failed to convert number from string !!!")
		return 0 , err
	}

	return byte(num) , nil
}

//for js
func PackStrToByteArray(vm *VM, pos uint64 , length uint64) ([]byte , error) {
	if length == 0 {
		return nil , ERR_EMPTY_INVALID_PARAM
	}

	arrByte , err := ConvertStr(vm.memory , pos , length)
	if err != nil {
		return nil , err
	}

	var i     uint64 = 0
	var idx   uint64 = 0
	var tmp []byte
	var res []byte
	var num   byte   = 0

	for {
		if arrByte[i] == 44 { // 44 = ','
			if len(tmp) == 0 {
				continue
			}

			if num ,err = ByteArrToNum(tmp) ; err != nil {
				return nil , err
			}
			res = append(res , num)

			tmp = RemoveElement(tmp , 0 , idx)
			idx = 0

		} else {
			tmp = append(tmp , arrByte[i])
			idx += 1
		}

		i += 1

		if i == length {
			if num ,err = ByteArrToNum(tmp) ; err != nil {
				return nil , err
			}
			res = append(res , num)
			break
		}

	}

	return res , nil
}

