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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

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
	callDep int
}

var wasmEng *wasmEngine

// VM_INSTANCE it means a VM instance , include its created time , end time and status
type vmInstance struct {
	vm         *VM       //it means a vm , it is a WASM module/file
	createTime time.Time //vm instance's created time
	endTime    time.Time //vm instance's deadline
}

// WASM_ENGINE struct wasm is a executable environment for other caller
type wasmEngine struct {
	//the string type need be modified
	vmMap        map[string]*vmInstance
	vmEngineLock *sync.Mutex

	//the channel is to communicate with each vm
	vmChannel chan []byte
}

type wasmInterface interface {
	Init() error
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

	index := vm.funcInfo.funcEntry.Index
	typeIndex := vm.module.Function.Types[int(index)]

	vm.funcInfo.funcType = vm.module.Types.Entries[int(typeIndex)]
	vm.funcInfo.funcIndex = int64(index)

	var err error
	var idx int

	idx, err = vm.StorageData(method)
	if err != nil {
		return errors.New("*ERROR* Failed to store the method name at the memory !!!")
	}
	vm.funcInfo.actIndex = uint64(idx)

	idx, err = vm.StorageData(param)
	if err != nil {
		return errors.New("*ERROR* Failed to store the method arguments at the memory !!!")
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
		fmt.Println("*ERROR* Failed to get account by name !!! ", err.Error())
		return 0
	}

	return binary.LittleEndian.Uint32(accountObj.CodeVersion.Bytes())
}

// NewWASM Search the CTX infor at the database according to applyContext
func NewWASM(ctx *contract.Context) *VM {
	var err error
	var wasmCode []byte

	var codeVersion uint32 = 0
		accountObj, err := ctx.RoleIntf.GetAccount(ctx.Trx.Contract)
		if err != nil {
			fmt.Println("*ERROR* Failed to get account by name !!! ", err.Error())
			return nil
		}
		codeVersion = binary.LittleEndian.Uint32(accountObj.CodeVersion.Bytes())
		wasmCode = accountObj.ContractCode

	module, err := wasm.ReadModule(bytes.NewBuffer(wasmCode), importer)
	if err != nil {
		fmt.Println("*ERROR* Failed to parse the wasm module !!! " + err.Error())
		return nil
	}

	if module.Export == nil {
		fmt.Println("*ERROR* Failed to find export method from wasm module !!!")
		return nil
	}

	vm, err := NewVM(module)
	if err != nil {
		return nil
	}

	vm.codeVersion = codeVersion

	return vm
}

func (engine *wasmEngine) Find(contractName string) (*vmInstance, error) {
	if len(engine.vmMap) == 0 {
		return nil, errors.New("*WARN* Can't find the vm instance !!!")
	}

	vmInst, ok := engine.vmMap[contractName]
	if !ok {
		return nil, errors.New("*WARN* Can't find the vm instance !!!")
	}

	return vmInst, nil
}

func (engine *wasmEngine) startSubCrx(event []byte) error {
	if event == nil {
		return errors.New("*ERROR* empty parameter !!!")
	}

	//unpack the crx from byte to struct
	var subCrx contract.Context

	if err := json.Unmarshal(event, &subCrx); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return errors.New("*ERROR* Failed to unpack contract from byte array to struct !!!")
	}

	//execute a new sub wasm crx
	go engine.Start(&subCrx, 1, false)

	return nil
}

//the function is called as a goruntine and to handle new vm request or other request
func (engine *wasmEngine) StartHandler() error {

	var event []byte //it means a MSG struct from ctx execution
	var ok bool

	for {
		event, ok = <-engine.vmChannel
		if !ok {
			continue
		}

		if len(event) == 1 && event[0] == 0 {
			break
		}
		engine.startSubCrx(event)
	}

	return nil
}

func (engine *wasmEngine) StopHandler() error {
	engine.vmChannel <- []byte{0}
	return nil
}

func (engine *wasmEngine) Init() error {
	return nil
}

func (engine *wasmEngine) Start(ctx *contract.Context, executionTime uint32, receivedBlock bool) ([]*types.Transaction, error) {
	return engine.Process(ctx, 1, executionTime, receivedBlock)
}

// Process the function is to be used for direct parameter insert
func (engine *wasmEngine) Process(ctx *contract.Context, depth uint8, executionTime uint32, receivedBlock bool) ([]*types.Transaction, error) {

	var pos int
	var err error
	var divisor time.Duration
	var deadline time.Time

	//search matched VM struct according to CTX
	var vm *VM = nil
	vmInst, ok := engine.vmMap[ctx.Trx.Contract]
	if !ok {
		vm = NewWASM(ctx)

		divisor, _ = time.ParseDuration(VM_PERIOD_OF_VALIDITY)
		deadline = time.Now().Add(divisor)

		engine.vmMap[ctx.Trx.Contract] = &vmInstance{
			vm:         vm,
			createTime: time.Now(),
			endTime:    deadline,
		}

		vm.SetContract(ctx)
		vm.SetChannel(engine.vmChannel)

	} else {
		vm = vmInst.vm
		vm.SetContract(ctx)
	}

	method := ENTRY_FUNCTION
	funcEntry, ok := vm.module.Export.Entries[method]
	if ok == false {
		return nil, errors.New("*ERROR* Failed to find the method from the wasm module !!!")
	}

	findex := funcEntry.Index
	ftype := vm.module.Function.Types[int(findex)]

	funcParams := make([]interface{}, 1)
	//Get function's string first char
	funcParams[0] = int([]byte(ctx.Trx.Method)[0])

	paramLength := len(funcParams)
	parameters := make([]uint64, paramLength)

	if paramLength != len(vm.module.Types.Entries[int(ftype)].ParamTypes) {
		return nil, errors.New("*ERROR* Parameters count is not right")
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
				return nil, errors.New("*ERROR* Failed to storage data to memory !!!")
			}
			parameters[i] = uint64(pos)
		default:
			return nil, errors.New("*ERROR* parameter is unsupport type !!!")
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
		return nil, errors.New("*ERROR* unsupported type !!!")
	}

	if result != 0 {
		//Todo failed to execute the crx , any handle operation
		return nil, errors.New("*ERROR* Failed to execute the contract !!! contract name: " + vm.contract.Trx.Contract)
	}

	value := make([]*types.Transaction, len(vm.subTrxLst))
	copy(value, vm.subTrxLst)
	vm.subTrxLst = vm.subTrxLst[:0]

	return value, nil
}
