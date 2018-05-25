
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
 * @Last Modified by:
 * @Last Modified time:
 */

package exec

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"encoding/binary"
	"bytes"
	"errors"
	"sync"
	"time"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/vm/wasm/wasm"
	"github.com/bottos-project/bottos/vm/wasm/validate"
	"github.com/bottos-project/bottos/contract"
)

var account_name uint64
const (
	INVOKE_FUNCTION       = "invoke"
	ENTRY_FUNCTION        = "start"

	CTX_WASM_FILE = "/opt/bin/go/usermng.wasm"
	SUB_WASM_FILE = "/opt/bin/go/sub.wasm"

	VM_PERIOD_OF_VALIDITY     = "1h"
	WAIT_TIME                 = 4

	BOT_INVALID_CODE          = 1

	CALL_DEP_LIMIT            = 5
	CALL_WID_LIMIT            = 10
)

type ParamList struct {
	Params []ParamInfo
}

type ParamInfo struct {
	Type string
	Val  string
}

type Rtn struct {
	Type string
	Val  string
}

type Apply_context struct {
	Msg         Message
}

type Authorization struct {
	Accout        string
	CodeVersion	  common.Hash
}

type Message struct {
	Wasm_name    string           //crx name
	Method_name  string           //method name
	Auth         Authorization
	Method_param []byte           //parameter
}

type FuncInfo struct {
	func_index int64
	act_index  uint64
	arg_index  uint64

	func_entry wasm.ExportEntry
	func_type  wasm.FunctionSig
}

type SUB_CRX_MSG struct {
	ctx        *contract.Context
	call_dep   int
}

var wasm_engine *WASM_ENGINE

//it means a VM instance , include its created time , end time and status
type VM_INSTANCE struct {
	vm            *VM        //it means a vm , it is a WASM module/file
	create_time   time.Time  //vm instance's created time
	end_time      time.Time  //vm instance's deadline
}

//struct wasm is a executable environment for other caller
type WASM_ENGINE struct {
	//the string type need be modified
	vm_map            map[string]*VM_INSTANCE
	vm_engine_lock    *sync.Mutex

	//the channel is to communicate with each vm
	vm_channel        chan []byte
}

type wasm_interface interface {
	Init() error
	//　a wrap for VM_Call
	Apply( ctx Apply_context ,execution_time uint32, received_block bool ) interface{}
	Start( ctx *contract.Context  , execution_time uint32, received_block bool ) (uint32 , error)
	Process( ctx *contract.Context , depth uint8 , execution_time uint32, received_block bool ) (uint32 , error)
	GetFuncInfo(module wasm.Module , entry wasm.ExportEntry) error
}

type VM_RUNTIME struct {
	vm_list []VM_INSTANCE
}

func GetInstance() *WASM_ENGINE {

	if wasm_engine == nil {
		wasm_engine = &WASM_ENGINE{
			vm_map        : make(map[string]*VM_INSTANCE),
			vm_engine_lock: new(sync.Mutex),
			vm_channel    : make(chan []byte , 10),
		}
		wasm_engine.Init()
	}

	return wasm_engine
}

func (vm *VM) GetFuncInfo(method string , param []byte) error {

	index := vm.funcInfo.func_entry.Index
	type_index := vm.module.Function.Types[int(index)]

	vm.funcInfo.func_type = vm.module.Types.Entries[int(type_index)]
	vm.funcInfo.func_index = int64(index)

	var err error
	var idx int

	idx , err = vm.StorageData(method)
	if err != nil {
		return errors.New("*ERROR* Failed to store the method name at the memory !!!")
	}
	vm.funcInfo.act_index = uint64(idx)

	idx , err = vm.StorageData(param)
	if err != nil {
		return errors.New("*ERROR* Failed to store the method arguments at the memory !!!")
	}
	vm.funcInfo.arg_index = uint64(idx)

	return nil
}

//reference to wasm-run
func importer ( name string ) (*wasm.Module, error) {
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


//Search the CTX infor at the database according to apply_context
func NewWASM ( ctx *contract.Context ) *VM {

	fmt.Println("NewWASM")

	var err       error
	var wasm_code []byte

	TST := false
	//if non-Test condition , get wasm_code from Accout
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

		wasm_code = accountObj.ContractCode
	} else {
		var wasm_file string
		if ctx.Trx.Contract == "sub" {
			wasm_file = SUB_WASM_FILE
		} else {
			wasm_file = CTX_WASM_FILE
		}

		wasm_code, err = ioutil.ReadFile(wasm_file)
		if err != nil {
			fmt.Println("*ERROR*  error in read file", err.Error())
			return nil
		}
	}

	module, err := wasm.ReadModule(bytes.NewBuffer(wasm_code), importer)
	if err != nil {
		fmt.Println("*ERROR* Failed to parse the wasm module !!! " + err.Error())
		return nil
	}

	if module.Export == nil  {
		fmt.Println("*ERROR* Failed to find export method from wasm module !!!")
		return nil
    }

	vm , err := NewVM(module)
	if err != nil {
		return nil
	}

	return vm
}

//as a goruntine to watch vm instance in wasm engine , it will be called by outer
func (engine *WASM_ENGINE) watch_vm () error {

	for {
		for contract_name , vm_instance := range engine.vm_map {

			if time.Now().After(vm_instance.end_time) {
				//engine.vm_engine_lock.Lock()

				delete(engine.vm_map, contract_name)

				//engine.vm_engine_lock.Unlock()
			}
		}

		time.Sleep(time.Second * WAIT_TIME)
	}

	return nil
}

func (engine *WASM_ENGINE) Find (contract_name string) (*VM_INSTANCE , error) {
	if len(engine.vm_map) == 0 {
		return nil , errors.New("*WARN* Can't find the vm instance !!!")
	}

	vm_instance , ok := engine.vm_map[contract_name];
	if !ok {
		return nil , errors.New("*WARN* Can't find the vm instance !!!")
	}

	return vm_instance , nil
}

func (engine *WASM_ENGINE) startSubCrx (event []byte) error {
	if event == nil {
		return errors.New("*ERROR* empty parameter !!!")
	}

	//Todo verify if event is a valid crx
	//github.com/asaskevich/govalidator

	//unpack the crx from byte to struct
	var sub_crx contract.Context

	if err := json.Unmarshal(event, &sub_crx) ; err != nil{
		fmt.Println("Unmarshal: ", err.Error())
		return errors.New("*ERROR* Failed to unpack contract from byte array to struct !!!")
	}

	fmt.Println("WASM_ENGINE::startSubCrx sub_crx = ",sub_crx)
	//check recursion limit
	/*
	if sub_crx.Trx.RecursionLayer > RECURSION_CALL_LIMIT {
		return errors.New("*ERROR* Exceeds maximum call number !!!")
	}
	*/

	//execute a new sub wasm crx
	go engine.Start(&sub_crx , 1 , false)

	return nil
}

//the function is called as a goruntine and to handle new vm request or other request
func (engine *WASM_ENGINE) StartHandler () error {

	fmt.Println("WASM_ENGINE::StartHandler")
	var event []byte  //it means a MSG struct from ctx execution
	var ok    bool

	for {
		event , ok = <- engine.vm_channel
		if ! ok {
			continue
		}

		if len(event) == 1 && event[0] == 0  {
			break
		}
		engine.startSubCrx(event)
	}

	return nil
}

func (engine *WASM_ENGINE) StopHandler () error {
	engine.vm_channel <- []byte{0}
	return nil
}

func (engine *WASM_ENGINE) Init() error {
	fmt.Println("Init")
	//ToDo load some initial operation
	return nil
}

//the function is to be used for json parameter
func (engine *WASM_ENGINE) Apply ( ctx *contract.Context  ,execution_time uint32, received_block bool ) (interface{} , error){

	fmt.Println("WASM_ENGINE::Apply() ")

	var divisor  time.Duration
	var deadline time.Time

	//search matched VM struct according to CTX
	var vm *VM = nil
	vm_instance , ok := engine.vm_map[ctx.Trx.Contract];
	if !ok {
		vm = NewWASM(ctx)

		divisor, _ = time.ParseDuration(VM_PERIOD_OF_VALIDITY)
		deadline   = time.Now().Add(divisor)

		engine.vm_map[ctx.Trx.Contract] = &VM_INSTANCE{
			vm:          vm,
			create_time: time.Now(),
			end_time:    deadline,
		}

		vm.SetContract(ctx)
		vm.SetChannel(engine.vm_channel)

	}else{
		vm = vm_instance.vm
		//to set a new context for a existing VM instance
		vm.SetContract(ctx)
	}

	//avoid that vm instance is deleted because of deadline
	//vm.vm_lock.Lock()

	vm.funcInfo.func_entry , ok = vm.module.Export.Entries[INVOKE_FUNCTION]
	if ok == false {
		return nil , errors.New("*ERROR* Failed to find invoke method from wasm module !!!")
	}

	if err := vm.GetFuncInfo(ctx.Trx.Method,ctx.Trx.Param); err != nil {
		return nil , err
	}

	output , err := vm.VM_Call()
	if err != nil {
		return nil , err
	}

	res, err := vm.GetData(uint64(binary.LittleEndian.Uint32(output)))
	if err != nil {
		return nil , err
	}

	result := &Rtn{}
	json.Unmarshal(res, result)

	//vm.vm_lock.Unlock()

	fmt.Println("result = ",result.Val)

	return nil , nil
}

func (vm *VM) VM_Call() ([]byte , error)  {

	func_params := make([]uint64, 2)
	func_params[0] = vm.funcInfo.act_index
	func_params[1] = vm.funcInfo.arg_index

	fmt.Println("VM::VM_Call() ")

	res, err := vm.ExecCode( vm.funcInfo.func_index , func_params ...)
	if err != nil {
		return nil , err
	}

	if res != 0 {
		//Todo failed to execute the crx , any handle operation
		return nil , errors.New("*ERROR* Failed to execute the contract !!! contract name: "+vm.contract.Trx.Contract)
	}

	switch vm.funcInfo.func_type.ReturnTypes[0] {
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

func (engine *WASM_ENGINE) Start ( ctx *contract.Context ,  execution_time uint32, received_block bool ) (uint32 , error) {
	fmt.Println("WASM_ENGINE::Start")
	return engine.Process(ctx , 1 ,execution_time , received_block )
}

//the function is to be used for direct parameter insert
func (engine *WASM_ENGINE) Process ( ctx *contract.Context , depth uint8 , execution_time uint32, received_block bool ) (uint32 , error) {

	fmt.Println("WASM_ENGINE::Process")

	var pos      int
	var err      error
	var divisor  time.Duration
	var deadline time.Time

	//search matched VM struct according to CTX
	var vm *VM = nil
	vm_instance , ok := engine.vm_map[ctx.Trx.Contract];
	if !ok {
		vm = NewWASM(ctx)

		divisor, _ = time.ParseDuration(VM_PERIOD_OF_VALIDITY)
		deadline   = time.Now().Add(divisor)

		engine.vm_map[ctx.Trx.Contract] = &VM_INSTANCE{
			vm:          vm,
			create_time: time.Now(),
			end_time:    deadline,
		}

		vm.SetContract(ctx)
		vm.SetChannel(engine.vm_channel)

	} else {
		vm = vm_instance.vm
		//to set a new context for a existing VM instance
		vm.SetContract(ctx)
	}

	//avoid that vm instance is deleted because of deadline
	//vm.vm_lock.Lock()
	method := ENTRY_FUNCTION
	func_entry , ok := vm.module.Export.Entries[method]
	if ok == false {
		return BOT_INVALID_CODE , errors.New("*ERROR* Failed to find the method from the wasm module !!!")
	}

	findex := func_entry.Index
	ftype  := vm.module.Function.Types[int(findex)]

	func_params    := make([]interface{}, 1)
	//Get function's string first char
	func_params[0]  = int([]byte(ctx.Trx.Method)[0])

	param_length := len(func_params)
	parameters   := make([]uint64, param_length)

	if param_length != len(vm.module.Types.Entries[int(ftype)].ParamTypes) {
		return BOT_INVALID_CODE , errors.New("*ERROR* Parameters count is not right")
	}

	// just handle parameter for entry function
	for i, param := range func_params {
		switch param.(type) {
		case int:
			parameters[i] = uint64(param.(int))
		case []byte:
			offset, err := vm.storageMemory(param.([]byte) , Int8)
			if err != nil {
				return BOT_INVALID_CODE, err
			}
			parameters[i] = uint64(offset)
		case string:
			if pos , err = vm.StorageData(param.(string)); err != nil {
				return BOT_INVALID_CODE , errors.New("*ERROR* Failed to storage data to memory !!!")
			}
			parameters[i] = uint64(pos)
		default:
			return BOT_INVALID_CODE , errors.New("*ERROR* parameter is unsupport type !!!")
		}
	}

	res, err := vm.ExecCode(int64(findex), parameters...)
	if err != nil {
		return BOT_INVALID_CODE , errors.New("*ERROR* Invalid result !" + err.Error())
	}

	var result uint32
	switch val := res.(type) {
	case uint32:
		result = val
	default:
		return BOT_INVALID_CODE , errors.New("*ERROR* unsupported type !!!")
	}

	if result != 0 {
		//Todo failed to execute the crx , any handle operation
		return result , errors.New("*ERROR* Failed to execute the contract !!! contract name: "+vm.contract.Trx.Contract)
	}

	if len(vm.sub_trx_lst) == 0 {
		return result , nil
	}


	if depth + 1 >= CALL_DEP_LIMIT {
		return BOT_INVALID_CODE , errors.New("*ERROR* Too much the number of new contract execution(dep) !!!")
	}

	//recursive call sub-trx
	for i , sub_trx := range vm.sub_trx_lst {

		if i + 1 > CALL_WID_LIMIT {
			return BOT_INVALID_CODE , errors.New("*ERROR* Too much the number of new contract execution(wid) !!!")
		}

		if result , err = engine.Process(sub_trx , depth + 1 ,  execution_time , received_block); err != nil {
			return result , err
		}
	}

	//clean
	vm.sub_trx_lst = vm.sub_trx_lst[:0]
	//vm.vm_lock.Unlock()

	return result , nil
}


