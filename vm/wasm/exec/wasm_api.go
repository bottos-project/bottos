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
	"io/ioutil"
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
	// INVOKE_FUNCTION config the function name
	INVOKE_FUNCTION = "invoke"
	// ENTRY_FUNCTION config the entry name
	ENTRY_FUNCTION = "start"

	// CTX_WASM_FILE config ctx wasm file
	CTX_WASM_FILE = "/opt/bin/go/usermng.wasm"
	// SUB_WASM_FILE config sub wasm file
	SUB_WASM_FILE = "/opt/bin/go/sub.wasm"

	// TST Test status
	TST = false

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
	Auth         Authorization
	MethodParam []byte //parameter
}

type FuncInfo struct {
	funcIndex int64
	actIndex  uint64
	argIndex  uint64

	funcEntry wasm.ExportEntry
	funcType  wasm.FunctionSig
}

type SUB_CRX_MSG struct {
	ctx      *contract.Context
	callDep int
}

var wasmEngine *WASM_ENGINE

// VM_INSTANCE it means a VM instance , include its created time , end time and status
type VM_INSTANCE struct {
	vm          *VM       //it means a vm , it is a WASM module/file
	createTime time.Time //vm instance's created time
	endTime    time.Time //vm instance's deadline
}

// WASM_ENGINE struct wasm is a executable environment for other caller
type WASM_ENGINE struct {
	//the string type need be modified
	vmMap         map[string]*VM_INSTANCE
	vmEngineLock  *sync.Mutex

	//the channel is to communicate with each vm
	vmChannel chan []byte
}

type wasm_interface interface {
	Init() error
	//　a wrap for VM_Call
	Apply(ctx ApplyContext, execution_time uint32, received_block bool) interface{}
	Start(ctx *contract.Context, execution_time uint32, received_block bool) (uint32, error)
	Process(ctx *contract.Context, depth uint8, execution_time uint32, received_block bool) (uint32, error)
	GetFuncInfo(module wasm.Module, entry wasm.ExportEntry) error
}

type VM_RUNTIME struct {
	vm_list []VM_INSTANCE
}

func GetInstance() *WASM_ENGINE {

	if wasmEngine == nil {
		wasmEngine = &WASM_ENGINE{
			vmMap:         make(map[string]*VM_INSTANCE),
			vmEngineLock:  new(sync.Mutex),
			vmChannel:     make(chan []byte, 10),
		}
		wasmEngine.Init()
	}

	return wasmEngine
}

func (vm *VM) GetFuncInfo(method string, param []byte) error {

	index      := vm.funcInfo.funcEntry.Index
	type_index := vm.module.Function.Types[int(index)]

	vm.funcInfo.funcType = vm.module.Types.Entries[int(type_index)]
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

func GetWasmVersion(ctx *contract.Context) uint32 {
	accountObj, err := ctx.RoleIntf.GetAccount(ctx.Trx.Contract)
	if err != nil {
		fmt.Println("*ERROR* Failed to get account by name !!! ", err.Error())
		return 0
	}

	return binary.LittleEndian.Uint32(accountObj.CodeVersion.Bytes())
}

// NewWASM Search the CTX infor at the database according to apply_context
func NewWASM(ctx *contract.Context) *VM {

	var err error
	var wasmCode []byte

	//if non-Test condition , get wasm_code from Accout
	var codeVersion uint32 = 0
	if !TST {
		//db handler will be invoked from Msg struct
		accountObj, err := ctx.RoleIntf.GetAccount(ctx.Trx.Contract)
		if err != nil {
			fmt.Println("*ERROR* Failed to get account by name !!! ", err.Error())
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
		wasmCode = accountObj.ContractCode
	} else {
		var wasmFile string
		if ctx.Trx.Contract == "sub" {
			wasmFile = SUB_WASM_FILE
		} else {
			wasmFile = CTX_WASM_FILE
		}

		wasmCode, err = ioutil.ReadFile(wasmFile)
		if err != nil {
			fmt.Println("*ERROR*  error in read file", err.Error())
			return nil
		}
	}

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

//as a goruntine to watch vm instance in wasm engine , it will be called by outer
func (engine *WASM_ENGINE) watchVm() error {

	for {
		for contractName, vmInstance := range engine.vmMap {

			if time.Now().After(vmInstance.endTime) {
				//engine.vm_engine_lock.Lock()

				delete(engine.vmMap, contractName)

				//engine.vm_engine_lock.Unlock()
			}
		}

		time.Sleep(time.Second * WAIT_TIME)
	}

	return nil
}

func (engine *WASM_ENGINE) Find(contractName string) (*VM_INSTANCE, error) {
	if len(engine.vmMap) == 0 {
		return nil, errors.New("*WARN* Can't find the vm instance !!!")
	}

	vmInstance, ok := engine.vmMap[contractName]
	if !ok {
		return nil, errors.New("*WARN* Can't find the vm instance !!!")
	}

	return vmInstance, nil
}

func (engine *WASM_ENGINE) startSubCrx(event []byte) error {
	if event == nil {
		return errors.New("*ERROR* empty parameter !!!")
	}

	//Todo verify if event is a valid crx
	//github.com/asaskevich/govalidator

	//unpack the crx from byte to struct
	var subCrx contract.Context

	if err := json.Unmarshal(event, &subCrx); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return errors.New("*ERROR* Failed to unpack contract from byte array to struct !!!")
	}

	//check recursion limit
	/*
		if sub_crx.Trx.RecursionLayer > RECURSION_CALL_LIMIT {
			return errors.New("*ERROR* Exceeds maximum call number !!!")
		}
	*/

	//execute a new sub wasm crx
	go engine.Start(&subCrx, 1, false)

	return nil
}

//the function is called as a goruntine and to handle new vm request or other request
func (engine *WASM_ENGINE) StartHandler() error {

	var event []byte //it means a MSG struct from ctx execution
	var ok    bool

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

func (engine *WASM_ENGINE) StopHandler() error {
	engine.vmChannel <- []byte{0}
	return nil
}

func (engine *WASM_ENGINE) Init() error {
	//ToDo load some initial operation
	return nil
}

// Apply the function is to be used for json parameter
func (engine *WASM_ENGINE) Apply(ctx *contract.Context, executionTime uint32, receivedBlock bool) (interface{}, error) {

	var divisor  time.Duration
	var deadline time.Time

	//search matched VM struct according to CTX
	var vm *VM = nil
	vmInstance, ok := engine.vmMap[ctx.Trx.Contract]
	if !ok {
		vm = NewWASM(ctx)

		divisor, _ = time.ParseDuration(VM_PERIOD_OF_VALIDITY)
		deadline   = time.Now().Add(divisor)

		engine.vmMap[ctx.Trx.Contract] = &VM_INSTANCE{
			vm:          vm,
			createTime:  time.Now(),
			endTime:     deadline,
		}

		vm.SetContract(ctx)
		vm.SetChannel(engine.vmChannel)

	} else {

		version := GetWasmVersion(ctx)
		//if version in local memory is different with the latest version in db , it need to update a new vm
		if version != vm.codeVersion {
			//create a new vm instance because of different code version
			vm = NewWASM(ctx)
		} else {
			vm = vmInstance.vm
		}

		//to set a new context for a existing VM instance
		vm.SetContract(ctx)
	}

	vm.funcInfo.funcEntry, ok = vm.module.Export.Entries[INVOKE_FUNCTION]
	if ok == false {
		return nil, errors.New("*ERROR* Failed to find invoke method from wasm module !!!")
	}

	if err := vm.GetFuncInfo(ctx.Trx.Method, ctx.Trx.Param); err != nil {
		return nil, err
	}

	output, err := vm.VM_Call()
	if err != nil {
		return nil, err
	}

	res, err := vm.GetData(uint64(binary.LittleEndian.Uint32(output)))
	if err != nil {
		return nil, err
	}

	result := &Rtn{}
	json.Unmarshal(res, result)

	fmt.Println("result = ", result.Val)

	return nil, nil
}

func (vm *VM) VM_Call() ([]byte, error) {

	func_params   := make([]uint64, 2)
	func_params[0] = vm.funcInfo.actIndex
	func_params[1] = vm.funcInfo.argIndex

	res, err := vm.ExecCode(vm.funcInfo.funcIndex, func_params...)
	if err != nil {
		return nil, err
	}

	if res != 0 {
		//Todo failed to execute the crx , any handle operation
		return nil, errors.New("*ERROR* Failed to execute the contract !!! contract name: " + vm.contract.Trx.Contract)
	}

	switch vm.funcInfo.funcType.ReturnTypes[0] {
	case wasm.ValueTypeI32:
		return I32ToBytes(res.(uint32)), nil
	case wasm.ValueTypeI64:
		return I64ToBytes(res.(uint64)), nil
	case wasm.ValueTypeF32:
		return F32ToBytes(res.(float32)), nil
	case wasm.ValueTypeF64:
		return F64ToBytes(res.(float64)), nil
	default:
		return nil, errors.New("*ERROR* the type of return value can't be supported")
	}
}

func (engine *WASM_ENGINE) Start(ctx *contract.Context, executionTime uint32, receivedBlock bool) ([]*types.Transaction, error) {
	return engine.Process(ctx, 1, executionTime, receivedBlock)
}

// Process the function is to be used for direct parameter insert
func (engine *WASM_ENGINE) Process(ctx *contract.Context, depth uint8, executionTime uint32, receivedBlock bool) ([]*types.Transaction, error) {

	var pos int
	var err error
	var divisor time.Duration
	var deadline time.Time

	//search matched VM struct according to CTX
	var vm *VM = nil
	vmInstance, ok := engine.vmMap[ctx.Trx.Contract]
	if !ok {
		vm = NewWASM(ctx)

		divisor, _ = time.ParseDuration(VM_PERIOD_OF_VALIDITY)
		deadline   = time.Now().Add(divisor)

		engine.vmMap[ctx.Trx.Contract] = &VM_INSTANCE{
			vm:          vm,
			createTime:  time.Now(),
			endTime:     deadline,
		}

		vm.SetContract(ctx)
		vm.SetChannel(engine.vmChannel)

	} else {
		/*
			version := GetWasmVersion(ctx)
			//if version in local memory is different with the latest version in db , it need to update a new vm
			if version != vm_instance.vm.codeVersion {
				//create a new vm instance because of different code version
				vm = NewWASM(ctx)
				vm_instance.vm = vm
			} else {
		*/
		vm = vmInstance.vm
		//}

		//to set a new context for a existing VM instance
		vm.SetContract(ctx)
	}

	method := ENTRY_FUNCTION
	funcEntry, ok := vm.module.Export.Entries[method]
	if ok == false {
		return nil, errors.New("*ERROR* Failed to find the method from the wasm module !!!")
	}

	findex := funcEntry.Index
	ftype  := vm.module.Function.Types[int(findex)]

	funcParams := make([]interface{}, 1)
	//Get function's string first char
	funcParams[0] = int([]byte(ctx.Trx.Method)[0])

	param_length := len(funcParams)
	parameters := make([]uint64, param_length)

	if param_length != len(vm.module.Types.Entries[int(ftype)].ParamTypes) {
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
