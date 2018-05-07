package exec

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"fmt"
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

	env_func.Register("myprint" , myprint)
	env_func.Register("get_str_value" , get_str_value)
	env_func.Register("get_usermng_reg_user" , get_usermng_reg_user)
	env_func.Register("get_usermng_user_login" , get_usermng_user_login)
	env_func.Register("set_str_value" , set_str_value)
	env_func.Register("printi" , printi)
	env_func.Register("prints" , prints)
	env_func.Register("get_test_str" , get_test_str)
	env_func.Register("get_param" , get_param)
	env_func.Register("set_test_byte" , set_test_byte)

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

func myprint(vm *VM) (bool, error) {
	pos := vm.envFunc.envFuncParam[0]
	len := vm.envFunc.envFuncParam[1]

	param := Bytes2String(vm.memory[pos:pos+len])
	fmt.Println("VM::myprint() param = ",param," , vm.memType = ",vm.memType[pos])

	return true , nil
}

func get_str_value(vm *VM) (bool, error) {
	fmt.Println("VM::get_str_value")
	return true , nil
}


func get_usermng_reg_user(vm *VM) (bool, error) {

	pos    := vm.envFunc.envFuncParam[0]
	length := vm.envFunc.envFuncParam[1]

	str := "Hello,World !!!"
	if length < uint64(len(str)) {
		return false , errors.New("*ERROR* out of the border of the memory !!!")
	}

	fmt.Println("VM::get_usermng_reg_user() pos = ",pos," , len = ",length)

	// Test to return a string value
	buff := bytes.NewBuffer(nil)
	buff.Write([]byte(str))
	bytes := buff.Bytes()

	if int(pos)+len(bytes) > len(vm.memory) {
		return false, errors.New("out of memory")
	}

	copy(vm.memory[int(pos):int(pos)+len(bytes)], bytes)
	vm.ctx = vm.envFunc.envFuncCtx

	// Test to return a int value
	/*
	ret := 123
	vm.ctx = vm.envFunc.envFuncCtx
	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(ret))
	}
	*/
	return true , nil
}

func get_usermng_user_login(vm *VM) (bool, error) {
	fmt.Println("VM::get_usermng_user_login")
	return true , nil
}

func set_str_value(vm *VM) (bool, error) {
	fmt.Println("VM::set_str_value")
	return true , nil
}

func printi(vm *VM) (bool, error) {

	//val := vm.envFunc.envFuncParam[0]
	//fmt.Println("VM::printi val = ",val)

	return true , nil
}

func prints(vm *VM) (bool, error) {

	var len uint64

	pos := vm.envFunc.envFuncParam[0]
	if _ , ok := vm.memType[pos]; ok {
		len = uint64(vm.memType[pos].Len)
	}else{
		len = vm.envFunc.envFuncParam[1]
	}

	param := Bytes2String(vm.memory[pos:pos+len])

	fmt.Println("VM::prints param = ",param)


	return true , nil
}

func get_test_str(vm *VM) (bool, error) {

	str := "string from get_test_str !!!"
	pos:=0
	// Test to return a string value
	buff := bytes.NewBuffer(nil)
	buff.Write([]byte(str))
	bytes := buff.Bytes()

	if int(pos)+len(bytes) > len(vm.memory) {
		return false, errors.New("*ERROR* out of memory")
	}

	copy(vm.memory[int(pos):int(pos)+len(bytes)], bytes)
	vm.ctx = vm.envFunc.envFuncCtx

	return true , nil
}

func get_param(vm *VM) (bool, error) {
	pos    := vm.envFunc.envFuncParam[0]
	//length := vm.envFunc.envFuncParam[1]

	str   := "Too Young , Too Simple !!!"
	bytes := []byte(str)
	if int(pos)+len(bytes) > len(vm.memory) {
		return false, errors.New("*ERROR* out of memory")
	}

	fmt.Println("VM::get_param() len(str) = ",len(str))

	copy(vm.memory[int(pos):int(pos)+len(bytes)], bytes)
	ret := len(str)
	vm.ctx = vm.envFunc.envFuncCtx
	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(ret))
	}

	return true , nil
}

func set_test_byte(vm *VM) (bool, error) {

	var len uint64

	pos := vm.envFunc.envFuncParam[0]
	if _ , ok := vm.memType[pos]; ok {
		len = uint64(vm.memType[pos].Len)
	}else{
		len = vm.envFunc.envFuncParam[1]
	}

	if len <= 0 {
		return false , errors.New("*ERROR* out of memory")
	}

	param := Bytes2String(vm.memory[pos : pos + len])
	fmt.Println("VM::set_test_byte() pos = ",pos," , len = ",len , " ,param = ",param)
	//fmt.Println("vm.memory = ",vm.memory)

	return true , nil
}

//ToDo move to conv.go
func BytesToString(bytes []byte) string {

	for i, b := range bytes {
		if b == 0 {
			return string(bytes[:i])
		}
	}
	return string(bytes)

}
