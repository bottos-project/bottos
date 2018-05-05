package exec

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"encoding/binary"
	"bytes"
	"errors"

	"github.com/bottos-project/core/vm/wasm/wasm"
	"github.com/bottos-project/core/vm/wasm/validate"
	"github.com/bottos-project/core/contract"
)


const (
	INVOKE_FUNCTION = "invoke"
	CTX_WASM_FILE = "C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\vm_bak\\testcase\\test_data2\\contract.wasm"
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


type FuncInfo struct {
	func_index int64
	act_index  uint64
	arg_index  uint64

	func_entry wasm.ExportEntry
	func_type  wasm.FunctionSig
}

var wasm_engine *WASM_ENGINE

//struct wasm is a executable environment for other caller
type WASM_ENGINE struct {
	vm      *VM             //it will be inited at NewVM() , one VM struct is on behalf of one wasm module
	vm_map  map[string]*VM  //the string type need be modified
}

type wasm_interface interface {


	Init() error
	//　a wrap for VM_Call
	Apply( ctx *contract.Context ,execution_time uint32, received_block bool ) interface{}

	GetFuncInfo(module wasm.Module , entry wasm.ExportEntry) error
}

func GetInstance() *WASM_ENGINE {

	if wasm_engine == nil {
		wasm_engine = &WASM_ENGINE{
			vm_map: make(map[string]*VM),
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


//Search the CTX infor at the database according to apply_context
func NewWASM ( ctx *contract.Context ) *VM {

	fmt.Println("NewWASM")

	var err error
	var wasm_code []byte

	TST := true
	//if non-Test condition , get wasm_code from Accout
	if !TST {
		//db handler will be invoked from Msg struct
		accountObj, err := ctx.RoleIntf.GetAccount(ctx.Trx.Contract)
		if err != nil {
			fmt.Println("*ERROR* Failed to get account by name !!! ", err.Error())
			return nil
		}

		/*
		// TODO
		if ctx.Msg.Auth.CodeVersion !=  account_name.CodeVersion{
			//check wasm file's hash
			//err = errors.New("*ERROR* Fail to match account's information !!!")

			return nil
		}
		*/
		wasm_code = accountObj.ContractCode
	} else {
		wasm_code, err = ioutil.ReadFile(CTX_WASM_FILE)
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


func (engine *WASM_ENGINE) Init() error {
	fmt.Println("Init")
	//load some initial operation
	return nil
}


func (engine *WASM_ENGINE) Apply ( ctx *contract.Context, execution_time uint32, received_block bool ) (interface{} , error){
	fmt.Println("Apply")

	//search matched VM struct according to CTX
	vm , ok := engine.vm_map[ctx.Trx.Contract];
	if !ok {
		vm = NewWASM(ctx)
		engine.vm_map[ctx.Trx.Contract] = vm
	}

	vm.funcInfo.func_entry , ok = vm.module.Export.Entries[INVOKE_FUNCTION]
	if ok == false {
		return nil , errors.New("*ERROR* Failed to find invoke method from wasm module !!!")
	}

	if err := vm.GetFuncInfo(ctx.Trx.Method, ctx.Trx.Param); err != nil {
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

	fmt.Println("result = ",result.Val)

	return nil , nil
}

func (vm *VM) VM_Call()  ([]byte , error)  {

	func_params := make([]uint64, 2)
	func_params[0] = vm.funcInfo.act_index
	func_params[1] = vm.funcInfo.arg_index

	res, err := vm.ExecCode( vm.funcInfo.func_index , func_params ...)
	if err != nil {
		return nil , err
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



