package exec

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"fmt"
	"github.com/bottos-project/core/contract"
	"github.com/bottos-project/core/common/types"
	"github.com/bottos-project/core/contract/msgpack"
)

type EnvFunc struct {
	envFuncMap      map[string]func(vm *VM) (bool, error)

	envFuncCtx      context
	envFuncParam    []uint64
	envFuncRtn      bool

	envFuncParamIdx int
	envMethod       string
}

func NewEnvFunc() *EnvFunc {
	env_func := EnvFunc{
		envFuncMap:      make(map[string]func(*VM) (bool, error)),
		envFuncParamIdx: 0,
	}

	env_func.Register("strcmp", stringcmp)
	env_func.Register("malloc", malloc)
	env_func.Register("arrayLen", arrayLen)
	env_func.Register("memcpy", memcpy)
	env_func.Register("JsonUnmashal", jsonUnmashal)
	env_func.Register("JsonMashal", jsonMashal)
	env_func.Register("memset", memset)

	env_func.Register("printi" , printi)
	env_func.Register("prints" , prints)
	env_func.Register("get_str_value" , get_str_value)
	env_func.Register("set_str_value" , set_str_value)
	env_func.Register("remove_str_value" , remove_str_value)
	env_func.Register("get_param" , get_param)
	env_func.Register("call_trx" , call_trx)
	env_func.Register("recv_trx" , recv_trx)
	env_func.Register("parse_param" , parse_param)

	return &env_func
}

func Bytes2String(bytes []byte) string {

	for i, b := range bytes {
		if b == 0 {
			return string(bytes[:i])
		}
	}
	return string(bytes)

}

func (env *EnvFunc) Register(method string, handler func(*VM) (bool, error)) {
	if _, ok := env.envFuncMap[method]; !ok {
		env.envFuncMap[method] = handler
	}
}

func (env *EnvFunc) Invoke(method string, vm *VM) (bool, error) {

	fc, ok := env.envFuncMap[method]
	if !ok {
		return false, errors.New("*ERROR* Failed to find method : " + method)
	}

	return fc(vm)
}

func (env *EnvFunc) GetEnvFuncMap() map[string]func(*VM) (bool, error) {
	return env.envFuncMap
}

func calloc(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params  := envFunc.envFuncParam

	if len(params) != 2 {
		return false, errors.New("*ERROR* Invalid parameter count during call calloc !!! ")
	}
	count  := int(params[0])
	length := int(params[1])
	//we don't know whats the alloc type here
	index, err := vm.getStoragePos((count*length), Unknown)
		if err != nil {
		return false, err
	}

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(index))
	}
	return true, nil
}
func malloc(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 1 {
		return false, errors.New("parameter count error while call calloc")
	}
	size := int(params[0])
	//we don't know whats the alloc type here
	index, err := vm.getStoragePos(size, Unknown)
	if err != nil {
		return false, err
	}

	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(index))
	}
	return true, nil

}

//use arrayLen to replace 'sizeof'
func arrayLen(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 1 {
		return false, errors.New("parameter count error while call arrayLen")
	}

	pointer := params[0]

	tl, ok := vm.memType[pointer]

	var result uint64
	if ok {
		switch tl.Type {
		case Int8, String:
			result = uint64(tl.Len / 1)
		case Int16:
			result = uint64(tl.Len / 2)
		case Int32, Float32:
			result = uint64(tl.Len / 4)
		case Int64, Float64:
			result = uint64(tl.Len / 8)
		case Unknown:
			//FIXME assume it's byte
			result = uint64(tl.Len / 1)
		default:
			result = uint64(0)
		}

	} else {
		result = uint64(0)
	}
	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(result))
	}
	return true, nil

}

func memcpy(vm *VM) (bool, error) {

	fmt.Println("VM::memcpy")

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 3 {
		return false, errors.New("parameter count error while call memcpy")
	}
	dest := int(params[0])
	src := int(params[1])
	length := int(params[2])

	if dest < src && dest+length > src {
		return false, errors.New("memcpy overlapped")
	}

	copy(vm.memory[dest:dest+length], vm.memory[src:src+length])

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(1))
	}

	return true, nil //this return will be dropped in wasm
}

func memset(vm *VM) (bool, error) {

	fmt.Println("VM::memset()")

	params := vm.envFunc.envFuncParam
	if len(params) != 3 {
		return false, errors.New("parameter count error while call memcpy")
	}
	dest := int(params[0])
	char := int(params[1])
	cnt := int(params[2])

	tmp := make([]byte, cnt)
	for i := 0; i < cnt; i++ {
		tmp[i] = byte(char)
	}

	copy(vm.memory[dest:dest+cnt], tmp)

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	//engine.vm.RestoreCtx()
	if vm.envFunc != nil {
		vm.ctx = vm.envFunc.envFuncCtx
	}

	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(1))
	}

	return true, nil //this return will be dropped in wasm
}

func readMessage(vm *VM) (bool, error) {

	fmt.Println("VM::readMessage")

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 2 {
		return false, errors.New("parameter count error while call readMessage")
	}

	addr := int(params[0])
	length := int(params[1])


	msgBytes, err := vm.GetMsgBytes()
	if err != nil {
		return false, err
	}


	if length != len(msgBytes) {
		return false, errors.New("readMessage length error")
	}
	copy(vm.memory[addr:addr+length], msgBytes[:length])
	vm.memType[uint64(addr)] = &typeInfo{Type: Unknown, Len: length}

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(length))
	}

	return true, nil
}

func readInt32Param(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 1 {
		return false, errors.New("parameter count error while call readInt32Param")
	}

	addr := params[0]
	paramBytes, err := vm.GetData(addr)
	if err != nil {
		return false, err
	}

	pidx := vm.envFunc.envFuncParamIdx

	if pidx+4 > len(paramBytes) {
		return false, errors.New("read params error")
	}

	retInt := binary.LittleEndian.Uint32(paramBytes[pidx : pidx+4])
	vm.envFunc.envFuncParamIdx += 4

	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(retInt))
	}
	return true, nil
}

func readInt64Param(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 1 {
		return false, errors.New("parameter count error while call readInt64Param")
	}

	addr := params[0]
	paramBytes, err := vm.GetData(addr)
	if err != nil {
		return false, err
	}

	pidx := vm.envFunc.envFuncParamIdx

	if pidx+8 > len(paramBytes) {
		return false, errors.New("read params error")
	}

	retInt := binary.LittleEndian.Uint64(paramBytes[pidx : pidx+8])
	vm.envFunc.envFuncParamIdx += 8

	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(retInt)
	}
	return true, nil
}

func readStringParam(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 1 {
		return false, errors.New("parameter count error while call readStringParam")
	}

	addr := params[0]
	paramBytes, err := vm.GetData(addr)
	if err != nil {
		return false, err
	}
	var length int

	pidx := vm.envFunc.envFuncParamIdx
	switch paramBytes[pidx] {
	case 0xfd: //uint16
		if pidx+3 > len(paramBytes) {
			return false, errors.New("read string failed")
		}
		length = int(binary.LittleEndian.Uint16(paramBytes[pidx+1 : pidx+3]))
		pidx += 3
	case 0xfe: //uint32
		if pidx+5 > len(paramBytes) {
			return false, errors.New("read string failed")
		}
		length = int(binary.LittleEndian.Uint16(paramBytes[pidx+1 : pidx+5]))
		pidx += 5
	case 0xff:
		if pidx+9 > len(paramBytes) {
			return false, errors.New("read string failed")
		}
		length = int(binary.LittleEndian.Uint16(paramBytes[pidx+1 : pidx+9]))
		pidx += 9
	default:
		length = int(paramBytes[pidx])
	}

	if pidx+length > len(paramBytes) {
		return false, errors.New("read string failed")
	}
	pidx += length + 1

	stringbytes := paramBytes[vm.envFunc.envFuncParamIdx+1 : pidx]

	retidx, err := vm.StorageData(stringbytes)
	if err != nil {
		return false, errors.New("set memory failed")
	}

	vm.envFunc.envFuncParamIdx = pidx
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(retidx))
	}
	return true, nil
}

func rawUnmashal(vm *VM) (bool, error) {

	fmt.Println("VM::rawUnmashal")
	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 3 {
		return false, errors.New("parameter count error while call jsonUnmashal")
	}

	pos := params[0]

	rawAddr := params[2]
	rawBytes, err := vm.GetData(rawAddr)
	if err != nil {
		return false, err
	}

	copy(vm.memory[pos:int(pos)+len(rawBytes)], rawBytes)

	return true, nil
}

func jsonUnmashal(vm *VM) (bool, error) {
	fmt.Println("VM::jsonUnmashal")
	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 3 {
		return false, errors.New("parameter count error while call jsonUnmashal")
	}

	addr := params[0]
	size := int(params[1])

	jsonaddr := params[2]
	jsonbytes, err := vm.GetData(jsonaddr)
	if err != nil {
		return false, err
	}
	paramList := &ParamList{}
	err = json.Unmarshal(jsonbytes, paramList)

	if err != nil {
		return false, err
	}

	buff := bytes.NewBuffer(nil)
	for _, param := range paramList.Params {
		switch strings.ToLower(param.Type) {
		case "int":
			tmp := make([]byte, 4)
			val, err := strconv.Atoi(param.Val)
			if err != nil {
				return false, err
			}
			binary.LittleEndian.PutUint32(tmp, uint32(val))
			buff.Write(tmp)

		case "int64":
			tmp := make([]byte, 8)
			val, err := strconv.ParseInt(param.Val, 10, 64)
			if err != nil {
				return false, err
			}
			binary.LittleEndian.PutUint64(tmp, uint64(val))
			buff.Write(tmp)

		case "int_array":
			arr := strings.Split(param.Val, ",")
			tmparr := make([]int, len(arr))
			for i, str := range arr {
				tmparr[i], err = strconv.Atoi(str)
				if err != nil {
					return false, err
				}
			}
			idx, err := vm.StorageData(tmparr)
			if err != nil {
				return false, err
			}
			tmp := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmp, uint32(idx))
			buff.Write(tmp)

		case "int64_array":
			arr := strings.Split(param.Val, ",")
			tmparr := make([]int64, len(arr))
			for i, str := range arr {
				tmparr[i], err = strconv.ParseInt(str, 10, 64)
				if err != nil {
					return false, err
				}
			}

			idx, err := vm.StorageData(tmparr)
			if err != nil {
				return false, err
			}
			tmp := make([]byte, 8)
			binary.LittleEndian.PutUint64(tmp, uint64(idx))
			buff.Write(tmp)

		case "string":
			idx, err := vm.StorageData(param.Val)
			if err != nil {
				return false, err
			}
			tmp := make([]byte, 4)
			binary.LittleEndian.PutUint32(tmp, uint32(idx))
			buff.Write(tmp)

		default:
			return false, errors.New("unsupported type :" + param.Type)
		}

	}

	bytes := buff.Bytes()
	if len(bytes) != size {
		//return false ,errors.New("")
		//todo this case is not an error, sizeof doesn't means actual memory length,so the size parameter should be removed.
	}
	//todo add more check

	if int(addr)+len(bytes) > len(vm.memory) {
		return false, errors.New("out of memory")
	}

	copy(vm.memory[int(addr):int(addr)+len(bytes)], bytes)
	vm.ctx = envFunc.envFuncCtx

	return true, nil
}

func jsonMashal(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam

	if len(params) != 2 {
		return false, errors.New("parameter count error while call jsonUnmashal")
	}

	val := params[0]
	ptype := params[1] //type
	tpstr, err := vm.GetData(ptype)
	if err != nil {
		return false, err
	}

	ret := &Rtn{}
	pstype := strings.ToLower(BytesToString(tpstr))
	ret.Type = pstype
	switch pstype {
	case "int":
		res := int(val)
		ret.Val = strconv.Itoa(res)

	case "int64":
		res := int64(val)
		ret.Val = strconv.FormatInt(res, 10)

	case "string":
		tmp, err := vm.GetData(val)
		if err != nil {
			return false, err
		}
		ret.Val = string(tmp)

	case "int_array":
		tmp, err := vm.GetData(val)
		if err != nil {
			return false, err
		}
		length := len(tmp) / 4
		retArray := make([]string, length)
		for i := 0; i < length; i++ {
			retArray[i] = strconv.Itoa(int(binary.LittleEndian.Uint32(tmp[i : i+4])))
		}
		ret.Val = strings.Join(retArray, ",")

	case "int64_array":
		tmp, err := vm.GetData(val)
		if err != nil {
			return false, err
		}
		length := len(tmp) / 8
		retArray := make([]string, length)
		for i := 0; i < length; i++ {
			retArray[i] = strconv.FormatInt(int64(binary.LittleEndian.Uint64(tmp[i:i+8])), 10)
		}
		ret.Val = strings.Join(retArray, ",")
	}

	jsonstr, err := json.Marshal(ret)
	if err != nil {
		return false, err
	}

	offset, err := vm.StorageData(string(jsonstr))
	if err != nil {
		return false, err
	}

	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(offset))
	}

	return true, nil
}

func stringcmp(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 2 {
		return false, errors.New("parameter count error while call strcmp")
	}

	var ret int

	addr1 := params[0]
	addr2 := params[1]

	fmt.Println("strcmp", addr1, addr2)

	if addr1 == addr2 {
		ret = 0
	} else {
		bytes1, err := vm.GetData(addr1)
		if err != nil {
			return false, err
		}

		bytes2, err := vm.GetData(addr2)
		if err != nil {
			return false, err
		}

		if BytesToString(bytes1) == BytesToString(bytes2) {
			ret = 0
		} else {
			ret = 1
		}
	}
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(ret))
	}
	return true, nil
}

func get_str_value(vm *VM) (bool, error) {
	contractCtx := vm.GetContract();

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 8 {
		return false, errors.New("parameter count error while call get_str_value")
	}
	contractPos := int(params[0])
	contractLen := int(params[1])
	objectPos := int(params[2])
	objectLen := int(params[3])
	keyPos := int(params[4])
	keyLen := int(params[5])
	valueBufPos := int(params[6])
	valueBufLen := int(params[7])

	// length check

	contract := make([]byte, contractLen)
	copy(contract, vm.memory[contractPos:contractPos+contractLen])

	object := make([]byte, objectLen)
	copy(object, vm.memory[objectPos:objectPos+objectLen])

	key := make([]byte, keyLen)
	copy(key, vm.memory[keyPos:keyPos+keyLen])

	fmt.Println(string(contract), len(contract), string(object), len(object), string(key), len(key))
	value, err := contractCtx.ContractDB.GetStrValue(string(contract), string(object), string(key))

	valueLen := 0
	if err == nil {
		valueLen = len(value);
		// check buf len
		if valueLen <= valueBufLen {
			copy(vm.memory[valueBufPos:valueBufPos+valueLen], []byte(value))
		} else {
			valueLen = 0;
		}
	}

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(valueLen))
	}

	fmt.Printf("VM: from contract:%v, method:%v, func get_test_str:(contract=%v, objname=%v, key=%v, value=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, contract, object, key, value);

	return true , nil
}

func set_str_value(vm *VM) (bool, error) {
	contractCtx := vm.GetContract();

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 6 {
		return false, errors.New("parameter count error while call set_str_value")
	}
	objectPos := int(params[0])
	objectLen := int(params[1])
	keyPos := int(params[2])
	keyLen := int(params[3])
	valuePos := int(params[4])
	valueLen := int(params[5])

	// length check

	object := make([]byte, objectLen)
	copy(object, vm.memory[objectPos:objectPos+objectLen])

	key := make([]byte, keyLen)
	copy(key, vm.memory[keyPos:keyPos+keyLen])

	value := make([]byte, valueLen)
	copy(value, vm.memory[valuePos:valuePos+valueLen])

	fmt.Println(string(object), len(object), string(key), len(key), string(value), len(value))
	err := contractCtx.ContractDB.SetStrValue(contractCtx.Trx.Contract, string(object), string(key), string(value))

	result := 1
	if err != nil {
		result = 0;
	}

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(result))
	}

	fmt.Printf("VM: from contract:%v, method:%v, func set_str_value:(objname=%v, key=%v, value=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, object, key, value);

	return true , nil
}

func remove_str_value(vm *VM) (bool, error) {
	contractCtx := vm.GetContract();

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 4 {
		return false, errors.New("parameter count error while call remove_str_value")
	}
	objectPos := int(params[0])
	objectLen := int(params[1])
	keyPos := int(params[2])
	keyLen := int(params[3])

	// length check

	object := make([]byte, objectLen)
	copy(object, vm.memory[objectPos:objectPos+objectLen])

	key := make([]byte, keyLen)
	copy(key, vm.memory[keyPos:keyPos+keyLen])

	fmt.Println(string(object), len(object), string(key), len(key))
	err := contractCtx.ContractDB.RemoveStrValue(contractCtx.Trx.Contract, string(object), string(key))

	result := 1
	if err != nil {
		result = 0;
	}

	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(result))
	}

	fmt.Printf("VM: from contract:%v, method:%v, func remove_str_value:(objname=%v, key=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, object, key);

	return true , nil
}

func printi(vm *VM) (bool, error) {
	contractCtx := vm.GetContract();
	value := vm.envFunc.envFuncParam[0]
	fmt.Printf("VM: from contract:%v, method:%v, func printi: %v\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, value);

	return true , nil
}

func prints(vm *VM) (bool, error) {
	//contractCtx := vm.GetContract();

	pos := vm.envFunc.envFuncParam[0]
	len := vm.envFunc.envFuncParam[1]

	value := make([]byte, len)
	copy(value, vm.memory[pos:pos+len])
	param := string(value)

	fmt.Printf("VM: prints: %v\n", param);
	return true , nil
}

func get_param(vm *VM) (bool, error) {
	contractCtx := vm.GetContract();

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 2 {
		return false, errors.New("parameter count error while call memcpy")
	}

	bufPos   := int(params[0])
	bufLen   := int(params[1])
	paramLen := len(contractCtx.Trx.Param)

	if bufLen <= paramLen {
		return false, errors.New("buffer not enough")
	}

	copy(vm.memory[int(bufPos):int(bufPos)+paramLen], contractCtx.Trx.Param)

	vm.ctx = vm.envFunc.envFuncCtx
	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(paramLen))
	}

	fmt.Printf("VM: from contract:%v, method:%v, func get_param:(param=%x)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, contractCtx.Trx.Param)

	return true , nil
}

func start_wasm (vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 1 {
		return false, errors.New("parameter count error while call memcpy")
	}

	return true,nil
}

func call_trx (vm *VM) (bool, error) {

	//check max call limit
	if vm.callWid > CALL_DEP_LIMIT {
		return false, errors.New("*ERROR* Too much the number of new contract execution(wid) !!!")
	}

	fmt.Println("VM::call_trx")
	envFunc := vm.envFunc
	params  := envFunc.envFuncParam

	if len(params) != 4 {
		return false, errors.New("*ERROR* Parameter count error while call memcpy")
	}

	c_pos  := int(params[0])
	m_pos  := int(params[1])
	p_pos  := int(params[2])
	p_len  := int(params[3])

	contrx   := BytesToString(vm.memory[c_pos:c_pos+vm.memType[uint64(c_pos)].Len-1])
	method   := BytesToString(vm.memory[m_pos:m_pos+vm.memType[uint64(m_pos)].Len-1])
	//the bytes after msgpack.Marshal
	param    := vm.memory[p_pos:p_pos + p_len]

	//below codes is just for test
	/*
	type transferparam struct {
		To			string
		Amount		uint32
	}

	var tf transferparam

	msgpack.Unmarshal(param , &tf)

	fmt.Println("VM::call_trx() param from contract: ",tf)
	*/

	trx := &types.Transaction{
		Version        : 1,
		CursorNum      : 1,
		CursorLabel    : 1,
		Lifetime       : 1,
		Sender         : vm.GetContract().Trx.Sender,
		Contract       : contrx,
		Method         : method,
		Param          : param,       //the bytes after msgpack.Marshal
		SigAlg         : 1,
		Signature      : []byte{},
	}
	ctx := &contract.Context{ Trx:trx}

	b_ctx , err := json.Marshal(ctx)
	if err != nil {
		return false , err
	}

	//Todo thread synchronization
	vm.callWid++

	vm.vm_channel <- b_ctx
	fmt.Println("Send Sem !!!")

	return true,nil
}

func recv_trx (vm *VM) (bool, error) {

	fmt.Println("VM::recv_trx")
	envFunc := vm.envFunc
	params  := envFunc.envFuncParam
	if len(params) != 2 {
		return false, errors.New("*ERROR* parameter count error while call memcpy")
	}

	crx_pos  := int(params[0])
	crx_len  := int(params[1])

	b_crx := vm.memory[crx_pos:crx_pos+crx_len]

	var crx contract.Context

	if err := json.Unmarshal(b_crx, &crx) ; err != nil{
		fmt.Println("Unmarshal: ", err.Error())
		return false , nil
	}


	vm.vm_channel <- b_crx
	fmt.Println("Send Sem !!!")

	return true,nil
}

func parse_param (vm *VM) (bool, error) {

	fmt.Println("VM::parse_param")
	envFunc := vm.envFunc
	params  := envFunc.envFuncParam

	if len(params) != 2 {
		return false, errors.New("*ERROR* Parameter count error while call memcpy")
	}

	param_pos  := int(params[0])
	param_len  := int(params[1])
	param      := vm.memory[param_pos:param_pos + param_len]

	type transferparam struct {
		To			string
		Amount		uint32
	}

	var tf transferparam
	msgpack.Unmarshal(param , &tf)

	fmt.Println("VM::parse_param() param from contract param: ", param)
	fmt.Println("VM::parse_param() param from contract tf:    ", tf)

	return true , nil
}