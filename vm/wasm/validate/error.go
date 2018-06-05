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
 * file description: the interface for WASM execution
 * @Author: Stewart Li
 * @Date:   2017-12-04
 * @Last Modified by:   Stewart Li
 * @Last Modified time: 2017-05-15
 */

package validate

import (
	"errors"
	"fmt"

	"github.com/bottos-project/bottos/vm/wasm/wasm"
	ops "github.com/bottos-project/bottos/vm/wasm/wasm/operators"
)

// Error define the error struct
type Error struct {
	Offset   int // Byte offset in the bytecode vector where the error occurs.
	Function int // Index into the function index space for the offending function.
	Err      error
}

// Error define the error func
func (e Error) Error() string {
	return fmt.Sprintf("error while validating function %d at offset %d: %v", e.Function, e.Offset, e.Err)
}

// ErrStackUnderflow define the error message
var ErrStackUnderflow = errors.New("validate: stack underflow")

// InvalidImmediateError define invalid immediate error
type InvalidImmediateError struct {
	ImmType string
	OpName  string
}

// Error define invalid immediate error message
func (e InvalidImmediateError) Error() string {
	return fmt.Sprintf("invalid immediate for op %s at (should be %s)", e.OpName, e.ImmType)
}

// UnmatchedOpError define byte type
type UnmatchedOpError byte

// Error define unmathed operator error
func (e UnmatchedOpError) Error() string {
	n1, _ := ops.New(byte(e))
	return fmt.Sprintf("encountered unmatched %s", n1.Name)
}

// InvalidLabelError define uint32 type
type InvalidLabelError uint32

// Error define invalid label error
func (e InvalidLabelError) Error() string {
	return fmt.Sprintf("invalid nesting depth %d", uint32(e))
}

// InvalidLocalIndexError define uint32 type
type InvalidLocalIndexError uint32

// Error define invalid local index error
func (e InvalidLocalIndexError) Error() string {
	return fmt.Sprintf("invalid index for local variable %d", uint32(e))
}

// InvalidTypeError define invalid type error struct
type InvalidTypeError struct {
	Wanted wasm.ValueType
	Got    wasm.ValueType
}

// Error define invalid type error
func (e InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type, got: %v, wanted: %v", e.Got, e.Wanted)
}

// InvalidElementIndexError define invalid element index error
type InvalidElementIndexError uint32

// Error
func (e InvalidElementIndexError) Error() string {
	return fmt.Sprintf("invalid element index %d", uint32(e))
}

// NoSectionError define section id
type NoSectionError wasm.SectionID

// Error define no section error
func (e NoSectionError) Error() string {
	return fmt.Sprintf("reference to non exist section (id %d) in module", wasm.SectionID(e))
}
