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

func (vm *VM) i32Const() {
	vm.pushUint32(vm.fetchUint32())
}

func (vm *VM) i64Const() {
	vm.pushUint64(vm.fetchUint64())
}

func (vm *VM) f32Const() {
	vm.pushFloat32(vm.fetchFloat32())
}

func (vm *VM) f64Const() {
	vm.pushFloat64(vm.fetchFloat64())
}
