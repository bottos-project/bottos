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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bottos-project/core/vm/wasm/exec"
	"github.com/bottos-project/core/vm/wasm/validate"
	"github.com/bottos-project/core/vm/wasm/wasm"
)

func main() {
	log.SetPrefix("wasm-run: ")
	log.SetFlags(0)

	verbose := flag.Bool("v", false, "enable/disable verbose mode")
	verify := flag.Bool("verify-module", false, "run module verification")

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}

	wasm.SetDebugMode(*verbose)

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	m, err := wasm.ReadModule(f, importer)
	if err != nil {
		log.Fatalf("could not read module: %v", err)
	}

	if *verify {
		err = validate.VerifyModule(m)
		if err != nil {
			log.Fatalf("could not verify module: %v", err)
		}
	}

	if m.Export == nil {
		log.Fatalf("module has no export section")
	}

	vm, err := exec.NewVM(m)
	if err != nil {
		log.Fatalf("could not create VM: %v", err)
	}

	for name, e := range m.Export.Entries {
		i := int64(e.Index)
		fidx := m.Function.Types[int(i)]
		ftype := m.Types.Entries[int(fidx)]
		switch len(ftype.ReturnTypes) {
		case 1:
			fmt.Printf("%s() %s => ", name, ftype.ReturnTypes[0])
		case 0:
			fmt.Printf("%s() => ", name)
		default:
			log.Printf("running exported functions with more than one return value is not supported")
			continue
		}
		if len(ftype.ParamTypes) > 0 {
			log.Printf("running exported functions with input parameters is not supported")
			continue
		}
		o, err := vm.ExecCode(i)
		if err != nil {
			fmt.Printf("\n")
			log.Printf("err=%v", err)
		}
		if len(ftype.ReturnTypes) == 0 {
			fmt.Printf("\n")
			continue
		}
		fmt.Printf("%[1]v (%[1]T)\n", o)
	}
}

func importer(name string) (*wasm.Module, error) {
	f, err := os.Open(name + ".wasm")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, err := wasm.ReadModule(f, nil)
	if err != nil {
		return nil, err
	}
	err = validate.VerifyModule(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
