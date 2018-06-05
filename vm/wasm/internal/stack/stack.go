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

// Package stack implements a growable uint64 stack
package stack

// Stack define stack type
type Stack struct {
	slice []uint64
}

// Push define push instruction
func (s *Stack) Push(b uint64) {
	s.slice = append(s.slice, b)
}

// Pop define pop instruction
func (s *Stack) Pop() uint64 {
	v := s.Top()
	s.slice = s.slice[:len(s.slice)-1]
	return v
}

// SetTop define settop instruction
func (s *Stack) SetTop(v uint64) {
	s.slice[len(s.slice)-1] = v
}

// Top define top instruction
func (s *Stack) Top() uint64 {
	return s.slice[len(s.slice)-1]
}

// Get define get instruction
func (s *Stack) Get(i int) uint64 {
	return s.slice[i]
}

// Set define set instruction
func (s *Stack) Set(i int, v uint64) {
	s.slice[i] = v
}

// Len define len instruction
func (s *Stack) Len() int {
	return len(s.slice)
}
