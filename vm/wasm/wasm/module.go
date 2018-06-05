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

package wasm

import (
	"errors"
	"io"

	"github.com/bottos-project/bottos/vm/wasm/wasm/internal/readpos"
	log "github.com/cihub/seelog"
)

// ErrInvalidMagic invalid magic
var ErrInvalidMagic = errors.New("wasm: Invalid magic number")

const (
	// Magic magic number
	Magic uint32 = 0x6d736100
	// Version version number
	Version uint32 = 0x1
)

// Function represents an entry in the function index space of a module.
type Function struct {
	Sig     *FunctionSig
	Body    *FunctionBody
	EnvFunc bool
	Method  string
}

// Module represents a parsed WebAssembly module:
// http://webassembly.org/docs/modules/
type Module struct {
	Version uint32

	Types    *SectionTypes
	Import   *SectionImports
	Function *SectionFunctions
	Table    *SectionTables
	Memory   *SectionMemories
	Global   *SectionGlobals
	Export   *SectionExports
	Start    *SectionStartFunction
	Elements *SectionElements
	Code     *SectionCode
	Data     *SectionData

	// The function index space of the module
	FunctionIndexSpace []Function
	GlobalIndexSpace   []GlobalEntry
	// function indices into the global function space
	// the limit of each table is its capacity (cap)
	TableIndexSpace        [][]uint32
	LinearMemoryIndexSpace [][]byte

	Other []Section // Other holds the custom sections if any

	imports struct {
		Funcs    []uint32
		Globals  int
		Tables   int
		Memories int
	}
}

// EnvGlobal global environment
type EnvGlobal struct {
	Env bool
	Val uint64
}

// NewEnvGlobal new global environment
func NewEnvGlobal(env bool, val uint64) *EnvGlobal {
	e := &EnvGlobal{
		Env: env,
		Val: val,
	}

	return e
}

// ResolveFunc is a function that takes a module name and
// returns a valid resolved module.
type ResolveFunc func(name string) (*Module, error)

// ReadModule reads a module from the reader r. resolvePath must take a string
// and a return a reader to the module pointed to by the string.
func ReadModule(r io.Reader, resolvePath ResolveFunc) (*Module, error) {
	reader := &readpos.ReadPos{
		R:      r,
		CurPos: 0,
	}
	m := &Module{}
	magic, err := readU32(reader)
	if err != nil {
		return nil, err
	}
	if magic != Magic {
		return nil, ErrInvalidMagic
	}
	if m.Version, err = readU32(reader); err != nil {
		return nil, err
	}

	for {
		done, err := m.readSection(reader)
		if err != nil {
			return nil, err
		} else if done {
			break
		}
	}

	m.LinearMemoryIndexSpace = make([][]byte, 1)
	if m.Table != nil {
		m.TableIndexSpace = make([][]uint32, int(len(m.Table.Entries)))
	}

	if m.Import != nil && resolvePath != nil {
		err := m.resolveImports(resolvePath) //resolvePath is importer() function
		if err != nil {
			return nil, err
		}
	}

	for _, fn := range []func() error{
		m.populateGlobals,
		m.populateFunctions,
		m.populateTables,
		m.populateLinearMemory,
	} {
		if err := fn(); err != nil {
			return nil, err
		}

	}

	log.Trace("There are %d entries in the function index space.", len(m.FunctionIndexSpace))
	return m, nil
}
