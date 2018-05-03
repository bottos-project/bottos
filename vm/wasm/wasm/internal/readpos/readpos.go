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

package readpos

import (
	"io"
)

// ReadPos implements io.Reader and stores the current number of bytes read from
// the reader
type ReadPos struct {
	R      io.Reader
	CurPos int64
}

// Read implements the io.Reader interface
func (r *ReadPos) Read(p []byte) (int, error) {
	n, err := r.R.Read(p)
	r.CurPos += int64(n)
	return n, err
}

// ReadByte implements the io.ByteReader interface
func (r *ReadPos) ReadByte() (byte, error) {
	p := make([]byte, 1)
	_, err := r.R.Read(p)
	return p[0], err
}
