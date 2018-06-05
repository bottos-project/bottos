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
 * file description:  interpeters
 * @Author: Stewart Li
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package exec

import (
	"math"
)

// these operations are essentially no-ops.
// TODO(vibhavp): Add optimisations to package compiles that
// removes them from the original bytecode.

func (vm *VM) i32ReinterpretF32() {
	vm.pushUint32(math.Float32bits(vm.popFloat32()))
}

func (vm *VM) i64ReinterpretF64() {
	vm.pushUint64(math.Float64bits(vm.popFloat64()))
}

func (vm *VM) f32ReinterpretI32() {
	vm.pushFloat32(math.Float32frombits(vm.popUint32()))
}

func (vm *VM) f64ReinterpretI64() {
	vm.pushFloat64(math.Float64frombits(vm.popUint64()))
}
