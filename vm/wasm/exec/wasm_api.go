// Copyright 2017~2022 The Bottos Authors
// This file is part of the Bottos Chain library.
// Created by Rocket Core Team of Bottos.

//This program is free software: you can distribute it and/or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.

//This program is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//GNU General Public License for more details.

//You should have received a copy of the GNU General Public License
// along with bottos.  If not, see <http://www.gnu.org/licenses/>.

/*
 * file description: the interface for WASM execution
 * @Author: Stewart Li
 * @Date:   2017-12-04
 * @Last Modified by:   Stewart Li
 * @Last Modified time: 2017-05-15
 */

package exec

import (
	"bytes"
	"encoding/binary"
	"errors"
	log "github.com/cihub/seelog"
	"os"
	"sync"
	"time"
	"fmt"
	"io/ioutil"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/vm/wasm/validate"
	"github.com/bottos-project/bottos/vm/wasm/wasm"
)

var accountName uint64

const (
	// ENTRY_FUNCTION config the entry name
	ENTRY_FUNCTION = "start"

	// VM_PERIOD_OF_VALIDITY config the VM period of validity
	VM_PERIOD_OF_VALIDITY = "1h"

	// WAIT_TIME config the wait times
	WAIT_TIME = 4

	// BOTTOS_INVALID_CODE config the status of code
	BOTTOS_INVALID_CODE = 1

	// CALL_DEP_LIMIT config the max depth of call
	CALL_DEP_LIMIT = 5
	// CALL_WID_LIMIT config the max width of call
	CALL_WID_LIMIT = 10

	CALL_MAX_VM_LIMIT = 10
	// CTX_WASM_FILE config ctx wasm file
	CTX_WASM_FILE = "/opt/bin/go/usermng.wasm"
	// SUB_WASM_FILE config sub wasm file
	SUB_WASM_FILE = "/opt/bin/go/sub.wasm"
	// TST Test status
	TST = true
)

// ParamList define param array
type ParamList struct {
	Params []ParamInfo
}

// ParamInfo define param info
type ParamInfo struct {
	Type string
	Val  string
}

// Rtn define the return type
type Rtn struct {
	Type string
	Val  string
}

// ApplyContext define the apply context
type ApplyContext struct {
	Msg Message
}

// Authorization define the struct of authorization
type Authorization struct {
	Accout      string
	CodeVersion common.Hash
}

// Message define the message info
type Message struct {
	WasmName    string //crx name
	MethodName  string //method name
	Auth        Authorization
	MethodParam []byte //parameter
}

//FuncInfo is function information
type FuncInfo struct {
	funcIndex int64
	actIndex  uint64
	argIndex  uint64

	funcEntry wasm.ExportEntry
	funcType  wasm.FunctionSig
}

type subCrxMsg struct {
	ctx     *contract.Context
	callDep  int
}

var wasmEng *wasmEngine

// vmInstance it means a VM instance , include its created time , end time and status
type vmInstance struct {
	vm         *VM        //it means a vm , it is a WASM module/file
	createTime  time.Time //vm instance's created time
	updateTime  time.Time //the time vm last used
}

// wasmEngine struct wasm is a executable environment for other caller
type wasmEngine struct {
	//the string type need be modified
	vmMap        map[string]*vmInstance
	vmEngineLock *sync.Mutex

	//the channel is to communicate with each vm
	vmChannel    chan []byte
}

type wasmInterface interface {
	Init() error
	//a wrap for vmCall
	Start(ctx *contract.Context, executionTime uint32, receivedBlock bool) (uint32, error)
	Process(ctx *contract.Context, depth uint8, executionTime uint32, receivedBlock bool) (uint32, error)
	GetFuncInfo(module wasm.Module, entry wasm.ExportEntry) error
}

type vmRuntime struct {
	vmList []vmInstance
}

//GetInstance is to get instance of wasm engine
func GetInstance() *wasmEngine {

	if wasmEng == nil {
		wasmEng = &wasmEngine{
			vmMap:        make(map[string]*vmInstance),
			vmEngineLock: new(sync.Mutex),
			vmChannel:    make(chan []byte, 10),
		}
		wasmEng.Init()
	}

	return wasmEng
}

//GetFuncInfo is to get function information
func (vm *VM) GetFuncInfo(method string, param []byte) error {

	index     := vm.funcInfo.funcEntry.Index
	typeIndex := vm.module.Function.Types[int(index)]

	vm.funcInfo.funcType  = vm.module.Types.Entries[int(typeIndex)]
	vm.funcInfo.funcIndex = int64(index)

	var err error
	var idx int

	idx, err = vm.StorageData(method)
	if err != nil {
		return ERR_STORE_MEMORY
	}
	vm.funcInfo.actIndex = uint64(idx)

	idx, err = vm.StorageData(param)
	if err != nil {
		return ERR_STORE_MEMORY
	}
	vm.funcInfo.argIndex = uint64(idx)

	return nil
}

//reference to wasm-run
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

//GetWasmVersion is to get wasm version
func GetWasmVersion(ctx *contract.Context) uint32 {
	accountObj, err := ctx.RoleIntf.GetAccount(ctx.Trx.Contract)
	if err != nil {
		log.Infof("*ERROR* Failed to get account by name !!! ", err.Error())
		return 0
	}

	return binary.LittleEndian.Uint32(accountObj.CodeVersion.Bytes())
}


// NewWASM Search the CTX infor at the database according to applyContext
func NewWASM(ctx *contract.Context) *VM {
	var err         error
	var wasmCode    []byte
	var codeVersion uint32

	accountObj, err := ctx.RoleIntf.GetAccount(ctx.Trx.Contract)
	if err != nil {
		log.Infof("*ERROR* Failed to get account by name !!! ", err.Error())
		return nil
	}
	codeVersion = binary.LittleEndian.Uint32(accountObj.CodeVersion.Bytes())
	wasmCode = accountObj.ContractCode

	module, err := wasm.ReadModule(bytes.NewBuffer(wasmCode), importer)
	if err != nil {
		log.Infof("*ERROR* Failed to parse the wasm module !!! " + err.Error())
		return nil
	}
	if module.Export == nil {
		log.Infof("*ERROR* Failed to find export method from wasm module !!!")
		return nil
	}

	vm, err := NewVM(module)
	if err != nil {
		return nil
	}

	vm.codeVersion = codeVersion

	return vm
}

//Search the CTX infor at the database according to apply_context
func NewWASMTst ( ctx *contract.Context ) *VM {

	fmt.Println("NewWASM")

	var err       error
	var wasm_code []byte

	//if non-Test condition , get wasm_code from Accout
	var codeVersion uint32 = 0
	if !TST {
		//db handler will be invoked from Msg struct
		accountObj, err := ctx.RoleIntf.GetAccount(ctx.Trx.Contract)
		if err != nil {
			log.Infof("*ERROR* Failed to get account by name !!! ", err.Error())
			return nil
		}

		/*
		if ctx.Trx.Version != accountObj.CodeVersion{
			//check wasm file's hash
			//err = errors.New("*ERROR* Fail to match account's information !!!")

			return nil
		}
		*/

		codeVersion = binary.LittleEndian.Uint32(accountObj.CodeVersion.Bytes())
		wasm_code   = accountObj.ContractCode
	} else {
		var wasm_file string
		if ctx.Trx.Contract == "sub" {
			wasm_file = SUB_WASM_FILE
		} else {
			wasm_file = CTX_WASM_FILE
		}

		wasm_code, err = ioutil.ReadFile(wasm_file)
		if err != nil {
			log.Infof("*ERROR*  error in read file", err.Error())
			return nil
		}
	}

	module, err := wasm.ReadModule(bytes.NewBuffer(wasm_code), importer)
	if err != nil {
		log.Infof("*ERROR* Failed to parse the wasm module !!! " + err.Error())
		return nil
	}

	if module.Export == nil  {
		log.Infof("*ERROR* Failed to find export method from wasm module !!!")
		return nil
	}

	vm , err := NewVM(module)
	if err != nil {
		return nil
	}

	vm.codeVersion = codeVersion

	return vm
}

func (engine *wasmEngine)  GetWasteVM() *vmInstance {
	var tmp        float64 = 0
	var updateTime float64 = 0
	var vmi        *vmInstance = nil
	for _ , vmInit := range engine.vmMap {
		updateTime = time.Now().Sub(vmInit.updateTime).Seconds()
		if tmp < updateTime {
			vmi = vmInit
		}
	}
	return vmi
}

func (engine *wasmEngine) Init() error {
	return nil
}

func (engine *wasmEngine) Start(ctx *contract.Context, executionTime uint32, receivedBlock bool) ([]*types.Transaction, error) {
	return engine.Process(ctx, 1, executionTime, receivedBlock)
}

// Process the function is to be used for direct parameter insert
func (engine *wasmEngine) Process(ctx *contract.Context, depth uint8, executionTime uint32, receivedBlock bool) ([]*types.Transaction, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err.(string))
			log.Infof(err.(string))
			return
		}
	}()

	var pos        int
	var err        error
	var updateTime time.Time

	//search matched VM struct according to CTX
	var vm  *VM = nil
	vmInst, ok := engine.vmMap[ctx.Trx.Contract]
	if !ok {
		vm = NewWASM(ctx)
		if vm == nil {
			return nil, ERR_CREATE_VM
		}

		updateTime = time.Now()
		engine.vmMap[ctx.Trx.Contract] = &vmInstance{
			vm:         vm,
			createTime: updateTime,
			updateTime: updateTime,
		}

		vm.SetContract(ctx)
		vm.SetChannel(engine.vmChannel)

	} else {
		vm = vmInst.vm
		if vm == nil {
			return nil, ERR_GET_VM
		}

		//if code's version in local memory is differsnt with the code's version , delete old one and update it
		if vm.codeVersion != ctx.Trx.Version {
			delete(engine.vmMap , ctx.Trx.Contract)
			vm = NewWASM(ctx)
			if vm == nil {
				return nil, ERR_CREATE_VM
			}

			updateTime = time.Now()
			engine.vmMap[ctx.Trx.Contract] = &vmInstance{
				vm:         vm,
				createTime: updateTime,
				updateTime: updateTime,
			}
			vm.SetChannel(engine.vmChannel)
		}

		vm.SetContract(ctx)
	}

	method        := ENTRY_FUNCTION
	funcEntry, ok := vm.module.Export.Entries[method]
	if ok == false {
		return nil, ERR_FIND_VM_METHOD
	}

	findex := funcEntry.Index
	ftype  := vm.module.Function.Types[int(findex)]

	funcParams   := make([]interface{}, 1)
	if pos, err = vm.StorageData(ctx.Trx.Method); err != nil {
		return nil, ERR_STORE_MEMORY
	}
	//Get Pos of function's string in memory
	funcParams[0] = pos

	paramLength  := len(funcParams)
	parameters   := make([]uint64, paramLength)

	if paramLength != len(vm.module.Types.Entries[int(ftype)].ParamTypes) {
		return nil, ERR_PARAM_COUNT
	}

	// just handle parameter for entry function
	for i, param := range funcParams {
		switch param.(type) {
		case int:
			parameters[i] = uint64(param.(int))
		case []byte:
			offset, err := vm.storageMemory(param.([]byte), Int8)
			if err != nil {
				return nil, err
			}
			parameters[i] = uint64(offset)
		case string:
			if pos, err = vm.StorageData(param.(string)); err != nil {
				return nil, ERR_STORE_MEMORY
			}
			parameters[i] = uint64(pos)
		default:
			return nil, ERR_UNSUPPORT_TYPE
		}
	}

	res, err := vm.ExecCode(int64(findex), parameters...)
	if err != nil {
		return nil, errors.New("*ERROR* Invalid result !" + err.Error())
	}

	var result uint32
	switch val := res.(type) {
	case uint32:
		result = val
	default:
		return nil, ERR_UNSUPPORT_TYPE
	}

	if result != 0 {
		//Todo failed to execute the crx , any handle operation
		return nil, errors.New("*ERROR* Failed to execute the contract !!! contract name: " + vm.contract.Trx.Contract)
	}

	value := make([]*types.Transaction, len(vm.subTrxLst))
	copy(value, vm.subTrxLst)
	vm.subTrxLst = vm.subTrxLst[:0]

	return value, err
}
