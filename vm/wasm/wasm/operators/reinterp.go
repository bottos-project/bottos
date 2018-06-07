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
	I32ReinterpretF32 = newOp(0xbc, "i32.reinterpret/f32", []wasm.ValueType{wasm.ValueTypeF32}, wasm.ValueTypeI32)
	I64ReinterpretF64 = newOp(0xbd, "i64.reinterpret/f64", []wasm.ValueType{wasm.ValueTypeF64}, wasm.ValueTypeI64)
	F32ReinterpretI32 = newOp(0xbe, "f32.reinterpret/i32", []wasm.ValueType{wasm.ValueTypeI32}, wasm.ValueTypeF32)
	F64ReinterpretI64 = newOp(0xbf, "f64.reinterpret/i64", []wasm.ValueType{wasm.ValueTypeI64}, wasm.ValueTypeF64)
)
