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

// DO NOT EDIT. This file define op code
var (
	Unreachable = newOp(0x00, "unreachable", nil, noReturn)
	Nop         = newOp(0x01, "nop", nil, noReturn)
	Block       = newOp(0x02, "block", nil, noReturn)
	Loop        = newOp(0x03, "loop", nil, noReturn)
	If          = newOp(0x04, "if", []wasm.ValueType{wasm.ValueTypeI32}, noReturn)
	Else        = newOp(0x05, "else", nil, noReturn)
	End         = newOp(0x0b, "end", nil, noReturn)
	Br          = newPolymorphicOp(0x0c, "br")
	BrIf        = newOp(0x0d, "br_if", []wasm.ValueType{wasm.ValueTypeI32}, noReturn)
	BrTable     = newPolymorphicOp(0x0e, "br_table")
	Return      = newPolymorphicOp(0x0f, "return")
)
