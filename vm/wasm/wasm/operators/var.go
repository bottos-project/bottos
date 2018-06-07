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

// DO NOT EDIT. follow define op code
var (
	GetLocal  = newPolymorphicOp(0x20, "get_local")
	SetLocal  = newPolymorphicOp(0x21, "set_local")
	TeeLocal  = newPolymorphicOp(0x22, "tee_local")
	GetGlobal = newPolymorphicOp(0x23, "get_global")
	SetGlobal = newPolymorphicOp(0x24, "set_global")
)
