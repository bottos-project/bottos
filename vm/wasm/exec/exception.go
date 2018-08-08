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
 * file description:  const push
 * @Author: Stewart Li
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package exec

import (
	"errors"
)

type TYPE uint32

const (
	VM_NOERROR TYPE = iota
	VM_ERROR_OUT_OF_MEMORY
	VM_ERROR_INVALID_PARAMETER_COUNT
	VM_ERROR_FAIL_EXECUTE_ENVFUNC
)

const (
	VM_NULL = iota
)

const (
	VM_FALSE = iota
	VM_TRUE
)

var ERR_EOF                      = errors.New("EOF")
var ERR_STORE_METHOD             = errors.New("*ERROR* failed to store the method name at the memory")
var ERR_STORE_PARAM              = errors.New("*ERROR* failed to store the method arguments at the memory")
var ERR_STORE_MEMORY             = errors.New("*ERROR* failed to storage data to memory")
var ERR_GET_STORE_POS            = errors.New("*ERROR* failed to get the position of data")
var ERR_CREATE_VM                = errors.New("*ERROR* failed to create a new VM instance")
var ERR_GET_VM                   = errors.New("*ERROR* failed to get a VM instance from memory")
var ERR_FIND_VM_METHOD           = errors.New("*ERROR* failed to find the method from the wasm module")
var ERR_PARAM_COUNT              = errors.New("*ERROR* parameters count is not right")
var ERR_UNSUPPORT_TYPE           = errors.New("*ERROR* unsupport type")
var ERR_OUT_BOUNDS               = errors.New("*ERROR* (array) index out of bounds")
var ERR_CALL_ENV_METHOD          = errors.New("*ERROR* failed to call the env method")
var ERR_EMPTY_INVALID_PARAM      = errors.New("*ERROR* empty parameter or invalid parameter")
var ERR_INVALID_WASM             = errors.New("*ERROR* invalid wasm module")
var ERR_DATA_INDEX               = errors.New("*ERROR* failed to get data index from memory")
var ERR_FINE_MAP                 = errors.New("*ERROR* the specified value can't be found by the key from the map")
// ErrMultipleLinearMemories is returned by (*VM).NewVM when the module
// has more then one entries in the linear memory space.
var ERR_MULTIPLE_LINEAR_MEMORIES = errors.New("*ERROR* more than one linear memories in module")
// ErrInvalidArgumentCount is returned by (*VM).ExecCode when an invalid
// number of arguments to the WebAssembly function are passed to it.
var ERR_INVALID_ARGUMENT_COUNT   = errors.New("*ERROR* invalid number of arguments to function")

