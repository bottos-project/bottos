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

package operators

import (
	"github.com/bottos-project/bottos/vm/wasm/wasm"
)

// DO NOT EDIT. follow define op code
var (
	I32Load    = newOp(0x28, "i32.load", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI32)
	I64Load    = newOp(0x29, "i64.load", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI64)
	F32Load    = newOp(0x2a, "f32.load", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeF32)
	F64Load    = newOp(0x2b, "f64.load", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeF64)
	I32Load8s  = newOp(0x2c, "i32.load8_s", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI32)
	I32Load8u  = newOp(0x2d, "i32.load8_u", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI32)
	I32Load16s = newOp(0x2e, "i32.load16_s", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI32)
	I32Load16u = newOp(0x2f, "i32.load16_u", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI32)
	I64Load8s  = newOp(0x30, "i64.load8_s", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI64)
	I64Load8u  = newOp(0x31, "i64.load8_u", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI64)
	I64Load16s = newOp(0x32, "i64.load16_s", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI64)
	I64Load16u = newOp(0x33, "i64.load16_u", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI64)
	I64Load32s = newOp(0x34, "i64.load32_s", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI64)
	I64Load32u = newOp(0x35, "i64.load32_u", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI64)

	I32Store   = newOp(0x36, "i32.store", []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, noReturn)
	I64Store   = newOp(0x37, "i64.store", []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI32}, noReturn)
	F32Store   = newOp(0x38, "f32.store", []wasm.ValueType{wasm.ValueTypeF32, wasm.ValueTypeI32}, noReturn)
	F64Store   = newOp(0x39, "f64.store", []wasm.ValueType{wasm.ValueTypeF64, wasm.ValueTypeI32}, noReturn)
	I32Store8  = newOp(0x3a, "i32.store8", []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, noReturn)
	I32Store16 = newOp(0x3b, "i32.store16", []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32}, noReturn)
	I64Store8  = newOp(0x3c, "i64.store8", []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI32}, noReturn)
	I64Store16 = newOp(0x3d, "i64.store16", []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI32}, noReturn)
	I64Store32 = newOp(0x3e, "i64.store32", []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI32}, noReturn)

	CurrentMemory = newOp(0x3f, "current_memory", nil, wasm.ValueTypeI32)
	GrowMemory    = newOp(0x40, "grow_memory", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeI32)
)
