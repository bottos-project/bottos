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

/*
 * file description:  convert variable
 * @Author: Stewart Li
 * @Date:   2017-12-07
 * @Last Modified by:
 * @Last Modified time:
 */

package exec

import (
	"fmt"
	"strings"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract"
	log "github.com/cihub/seelog"
)

// EnvFunc defines env for func execution
type EnvFunc struct {
	envFuncMap map[string]func(vm *VM) (bool, error)

	envFuncCtx   context
	envFuncParam []uint64
	envFuncRtn   bool

	envFuncParamIdx int
	envMethod       string
}

// NewEnvFunc new an EnvFunc
func NewEnvFunc() *EnvFunc {
	envFunc := EnvFunc{
		envFuncMap:      make(map[string]func(*VM) (bool, error)),
		envFuncParamIdx: 0,
	}

	//env function for C/C++
	envFunc.Register("printi",           printi)
	envFunc.Register("prints",           prints)
	envFunc.Register("getStrValue",      getStrValue)
	envFunc.Register("setStrValue",      setStrValue)
	envFunc.Register("removeStrValue",   removeStrValue)
	envFunc.Register("getStringValue",   getStrValue)
	envFunc.Register("setStringValue",   setStrValue)
	envFunc.Register("removeStringValue",removeStrValue)
	envFunc.Register("getParam",         getParam)
	envFunc.Register("getMethod",        getMethod)
	envFunc.Register("callTrx",          callTrx)
	envFunc.Register("assert",           assert)
	envFunc.Register("getCtxName",       getCtxName)
	envFunc.Register("getSender",        getSender)
	envFunc.Register("malloc",           malloc)
	envFunc.Register("memset",           memset)
	envFunc.Register("memcpy",           memcpy)
	envFunc.Register("strcat_s",         strcat_s)
	envFunc.Register("strcpy_s",         strcpy_s)
	envFunc.Register("isAccountExist",   isAccountExist)

	envFunc.Register("getMethodJs",       getMethodJs)

	return &envFunc
}

// Register register a method in VM
func (env *EnvFunc) Register(method string, handler func(*VM) (bool, error)) {
	if _, ok := env.envFuncMap[method]; !ok {
		env.envFuncMap[method] = handler
	}
}

// GetEnvFuncMap retrieve a method from FuncMap
func (env *EnvFunc) GetEnvFuncMap() map[string]func(*VM) (bool, error) {
	return env.envFuncMap
}

//uint32_t getStrValue(unsigned char * contract, uint32_t contractlen, unsigned char * object, uint32_t objlen, unsigned char * key,   uint32_t keylen, unsigned char * value_buf, uint32_t value_buf_len);
func getStrValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params  := envFunc.envFuncParam
	if len(params) != 8 {
		return false, ERR_PARAM_COUNT
	}
	contractPos := params[0]
	contractLen := params[1]
	objectPos   := params[2]
	objectLen   := params[3]
	keyPos      := params[4]
	keyLen      := params[5]
	valueBufPos := params[6]
	valueBufLen := params[7]
	vmLen       := uint64(len(vm.memory))

	if valueBufPos >= vmLen || valueBufPos + valueBufLen >= vmLen {
		fmt.Println("VM::getStrValue *ERROR* Out of bound")
		log.Infof("*ERROR* Out of bound \n")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	contract , err := Convert(vm , contractPos , contractLen)
	if err != nil {
		return true, nil
	}
	object   , err := Convert(vm , objectPos , objectLen)
	if err != nil {
		return true, nil
	}
	key      , err := Convert(vm , keyPos , keyLen)
	if err != nil {
		return true, nil
	}

	log.Infof(string(contract), len(contract), string(object), len(object), string(key), len(key))
	value, err := contractCtx.ContractDB.GetStrValue(string(contract), string(object), string(key))

	var valueLen uint64 = 0
	if err == nil {
		valueLen = uint64(len(value))
		// check buf len
		if valueLen <= valueBufLen {
			copy(vm.memory[valueBufPos:valueBufPos+valueLen], []byte(value))
		} else {
			valueLen = 0
		}
		vm.memory[valueBufPos+valueLen] = 0
	}

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(valueLen))
	}

	log.Infof("VM: from contract:%v, method:%v, func get_test_str:(contract=%v, objname=%v, key=%v, value=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, contract, object, key, value)

	return true, nil
}

//uint32_t setStrValue(unsigned char * object,   uint32_t objlen,      unsigned char * key,    uint32_t keylen, unsigned char * value, uint32_t vallen);
func setStrValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 6 {
		return false, ERR_PARAM_COUNT
	}
	objectPos := params[0]
	objectLen := params[1]
	keyPos    := params[2]
	keyLen    := params[3]
	valuePos  := params[4]
	valueLen  := params[5]

	object , err := Convert(vm , objectPos , objectLen)
	if err != nil {
		return true, nil
	}
	key     , err := Convert(vm , keyPos , keyLen)
	if err != nil {
		return true, nil
	}
	value   , err := Convert(vm , valuePos , valueLen)
	if err != nil {
		return true, nil
	}

	log.Infof(string(object), len(object), string(key), len(key), string(value), len(value))
	result := 1
	err = contractCtx.ContractDB.SetStrValue(contractCtx.Trx.Contract, string(object), string(key), string(value))
	if err != nil {
		result = 0
	}

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(result))
	}

	log.Infof("VM: from contract:%v, method:%v, func setStrValue:(objname=%v, key=%v, value=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, object, key, value)

	return true, nil
}

//uint32_t removeStrValue(unsigned char * object, uint32_t objlen, unsigned char * key, uint32_t keylen);
func removeStrValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 4 {
		return false, ERR_PARAM_COUNT
	}
	objectPos := params[0]
	objectLen := params[1]
	keyPos    := params[2]
	keyLen    := params[3]

	object , err := Convert(vm , objectPos , objectLen)
	if err != nil {
		return true, nil
	}
	key     , err := Convert(vm , keyPos , keyLen)
	if err != nil {
		return true, nil
	}

	log.Infof(string(object), len(object), string(key), len(key))
	err = contractCtx.ContractDB.RemoveStrValue(contractCtx.Trx.Contract, string(object), string(key))

	result := 1
	if err != nil {
		result = 0
	}

	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(result))
	}

	log.Infof("VM: from contract:%v, method:%v, func removeStrValue:(objname=%v, key=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, object, key)

	return true, nil
}

func getBinValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 8 {
		return false, ERR_PARAM_COUNT
	}
	contractPos := params[0]
	contractLen := params[1]
	objectPos   := params[2]
	objectLen   := params[3]
	keyPos      := params[4]
	keyLen      := params[5]
	valueBufPos := params[6]
	valueBufLen := params[7]

	contract , err := Convert(vm , contractPos , contractLen)
	if err != nil {
		return true, nil
	}
	object   , err := Convert(vm , objectPos , objectLen)
	if err != nil {
		return true, nil
	}
	key      , err := Convert(vm , keyPos , keyLen)
	if err != nil {
		return true, nil
	}

	log.Infof(string(contract), len(contract), string(object), len(object), string(key), len(key))
	var valueLen uint64 = 0
	value, err := contractCtx.ContractDB.GetBinValue(string(contract), string(object), string(key))
	if err == nil {
		valueLen = uint64(len(value))
		// check buf len
		if valueLen <= valueBufLen {
			copy(vm.memory[valueBufPos:valueBufPos+valueLen], value)
		} else {
			valueLen = 0
		}
	}

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(valueLen))
	}

	log.Infof("VM: from contract:%v, method:%v, func get_bin_value:(contract=%v, objname=%v, key=%v, value=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, contract, object, key, value)

	return true, nil
}

func setBinValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 6 {
		return false, ERR_PARAM_COUNT
	}
	objectPos := params[0]
	objectLen := params[1]
	keyPos    := params[2]
	keyLen    := params[3]
	valuePos  := params[4]
	valueLen  := params[5]

	object   , err := Convert(vm , objectPos , objectLen)
	if err != nil {
		return true, nil
	}
	key      , err := Convert(vm , keyPos , keyLen)
	if err != nil {
		return true, nil
	}
	value    , err := Convert(vm , valuePos , valueLen)
	if err != nil {
		return true, nil
	}

	log.Infof(string(object), len(object), string(key), len(key), string(value), len(value))
	err = contractCtx.ContractDB.SetBinValue(contractCtx.Trx.Contract, string(object), string(key), value)

	result := 1
	if err != nil {
		result = 0
	}

	//1. recover the vm context
	//2. if the call returns value,push the result to the stack
	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(result))
	}

	log.Infof("VM: from contract:%v, method:%v, func setBinValue:(objname=%v, key=%v, value=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, object, key, value)

	return true, nil
}

func removeBinValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params  := envFunc.envFuncParam
	if len(params) != 4 {
		return false, ERR_PARAM_COUNT
	}
	objectPos := params[0]
	objectLen := params[1]
	keyPos    := params[2]
	keyLen    := params[3]

	object   , err := Convert(vm , objectPos , objectLen)
	if err != nil {
		return true, nil
	}
	key      , err := Convert(vm , keyPos , keyLen)
	if err != nil {
		return true, nil
	}

	log.Infof(string(object), len(object), string(key), len(key))
	err = contractCtx.ContractDB.RemoveBinValue(contractCtx.Trx.Contract, string(object), string(key))

	result := 1
	if err != nil {
		result = 0
	}

	vm.ctx = envFunc.envFuncCtx
	if envFunc.envFuncRtn {
		vm.pushUint64(uint64(result))
	}

	log.Infof("VM: from contract:%v, method:%v, func removeBinValue:(objname=%v, key=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, object, key)

	return true, nil
}

//void     printi(uint32_t value);
func printi(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()
	value       := vm.envFunc.envFuncParam[0]

	fmt.Printf("VM: from contract: %v, method: %v, func printi: %v\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, value)
	log.Infof("VM: from contract:%v, method:%v, func printi: %v\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, value)

	return true, nil
}

func printi64(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()
	value       := vm.envFunc.envFuncParam[0]
	fmt.Printf("VM: from contract: %v, method: %v, func printi64: %v\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, value)
	log.Infof("VM: from contract:%v, method:%v, func printi64: %v\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, value)

	return true, nil
}

//void     prints(unsigned char * str, uint32_t len);
func prints(vm *VM) (bool, error) {
	pos := vm.envFunc.envFuncParam[0]
	len := vm.envFunc.envFuncParam[1]

	value , err := Convert(vm , pos , len)
	if err != nil {
		fmt.Println("*ERROR* vm::prints failed to convert parameter in prints , err: ", err)
		log.Infof("*ERROR* vm::prints failed to convert parameter in prints , err: ", err)
		return true, nil
	}

	param := string(value)
	fmt.Println("VM: func prints: ", param)
	log.Infof("VM: func prints: %v\n", param)
	return true, nil
}

//uint32_t getMethod(unsigned char * param, uint32_t buf_len);
func getMethod(vm *VM) (bool, error) {
	params := vm.envFunc.envFuncParam
	if len(params) != 2 {
		return false, ERR_PARAM_COUNT
	}

	pos    := params[0]
	length := params[1]
	vmLen  := uint64(len(vm.memory))
    if pos >= vmLen || pos + length >= vmLen {
		fmt.Println("VM::getMethod *ERROR* Out of bound")
		log.Infof("*ERROR* Out of bound \n")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	contractCtx := vm.GetContract()
	methodLen   := uint64(len(contractCtx.Trx.Method))
	if methodLen > length {
		log.Infof("*ERROR* Invaild string length \n")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	if uint64(copy(vm.memory[pos:pos + methodLen], []byte(contractCtx.Trx.Method))) != methodLen {
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(methodLen))
	}

	return true, nil
}

//uint32_t getParam (unsigned char * param, uint32_t buf_len);
func getParam(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params  := envFunc.envFuncParam
	if len(params) != 2 {
		return false, ERR_PARAM_COUNT
	}

	bufPos   := params[0]
	bufLen   := params[1]
	vmLen    := uint64(len(vm.memory))
	paramLen := uint64(len(contractCtx.Trx.Param))
	if bufPos >= vmLen || bufPos + bufLen >= vmLen {
		fmt.Println("VM::getParam *ERROR* Out of bound")
		log.Infof("*ERROR* Out of bound \n")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	if bufLen <= paramLen {
		log.Infof("*ERROR* Invaild string length \n")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	copy(vm.memory[bufPos : bufPos + paramLen], contractCtx.Trx.Param)

	vm.ctx = vm.envFunc.envFuncCtx
	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(paramLen))
	}

	return true, nil
}

//uint32_t callTrx(unsigned char * contract , unsigned char * method , unsigned char * buf , uint32_t buf_len );
func callTrx(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params := envFunc.envFuncParam

	if len(params) != 4 {
		return false, ERR_PARAM_COUNT
	}

	cPos := params[0]
	mPos := params[1]
	pPos := params[2]
	pLen := params[3]

	contrxByte, err := Convert(vm, cPos, vm.StrLen(cPos))
	if err != nil {
		return true, nil
	}
	methodByte, err := Convert(vm, mPos, vm.StrLen(mPos))
	if err != nil {
		return true, nil
	}

	contrx := BytesToString(contrxByte)
	method := BytesToString(methodByte)

	var param []byte
	//the bytes after msgpack.Marshal
	if vm.sourceFile == CPP {
		param = vm.memory[pPos: pPos+pLen]
	} else if vm.sourceFile == JS {
		param , err = PackStrToByteArray(vm, pPos, vm.StrLen(pPos))
		if err != nil {
			return true, nil
		}
	}

	value := make([]byte, len(param))
	copy(value, param)

	trx := &types.Transaction{
		Version:     vm.contract.Trx.Version,
		CursorNum:   vm.contract.Trx.CursorNum,
		CursorLabel: vm.contract.Trx.CursorLabel,
		Lifetime:    vm.contract.Trx.Lifetime,
		Sender:      vm.contract.Trx.Contract,
		Contract:    contrx,
		Method:      method,
		Param:       value, //the bytes after msgpack.Marshal
		SigAlg:      vm.contract.Trx.SigAlg,
		Signature:   []byte{},
	}
	ctx := &contract.Context{RoleIntf: vm.GetContract().RoleIntf, ContractDB: vm.GetContract().ContractDB, Trx: trx}

	//Todo thread synchronization
	vm.callWid++

	vm.subCtnLst = append(vm.subCtnLst, ctx)
	vm.subTrxLst = append(vm.subTrxLst, trx)

	if vm.envFunc.envFuncRtn {
		vm.pushUint32(uint32(VM_NOERROR))
	}

	return true, nil
}

//uint32_t assert (bool condition);
func assert(vm *VM) (bool, error) {
	envFunc := vm.envFunc
	params  := envFunc.envFuncParam

	cond := params[0]
	if cond != 1 {
		errStr := "*ERROR* failed to execute safe-function !!!"
		log.Infof(errStr)
		panic(errStr)
	}

	return true, nil
}

//uint32_t getCtxName(unsigned char * str , uint32_t len);
func getCtxName(vm *VM) (bool, error) {

	ctxName    := vm.contract.Trx.Contract
	ctxNameLen := uint64(len(ctxName))

	pos    := vm.envFunc.envFuncParam[0]
	length := vm.envFunc.envFuncParam[1]
	vmLen  := uint64(len(vm.memory))
	if pos >= vmLen || pos + length >= vmLen {
		fmt.Println("VM::getCtxName *ERROR* Out of bound")
		log.Infof("*ERROR* Out of bound \n")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	if length < ctxNameLen + 1 {
		log.Infof("*ERROR* Invaild string length \n")
		if vm.envFunc.envFuncRtn {
			vm.pushInt32(int32(VM_NULL))
		}
		return true, nil
	}

	copy(vm.memory[pos:pos+ctxNameLen], []byte(ctxName))
	vm.memory[pos+ctxNameLen] = 0
	if vm.envFunc.envFuncRtn {
		vm.pushInt32(int32(ctxNameLen))
	}

	return true, nil
}

//uint32_t getSender (unsigned char * str , uint32_t len);
func getSender(vm *VM) (bool, error) {

	senderName := vm.contract.Trx.Sender
	senderNameLen := uint64(len(senderName))

	pos    := vm.envFunc.envFuncParam[0]
	length := vm.envFunc.envFuncParam[1]
	vmLen  := uint64(len(vm.memory))
	if pos >= vmLen || pos + length >= vmLen {
		fmt.Println("VM::getSender *ERROR* Out of bound")
		log.Infof("*ERROR* Out of bound \n")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	if length < senderNameLen + 1 {
		log.Infof("*ERROR* Invaild string length \n")
		if vm.envFunc.envFuncRtn {
			vm.pushInt32(int32(VM_NULL))
		}
		return true, nil
	}

	copy(vm.memory[pos:pos+senderNameLen], []byte(senderName))
	vm.memory[pos+senderNameLen] = 0
	if vm.envFunc.envFuncRtn {
		vm.pushInt32(int32(senderNameLen))
	}

	return true, nil
}

//void    *memset(void * ptr, int value, size_t num);
func memset(vm *VM) (bool, error) {
	params  := vm.envFunc.envFuncParam
	if len(params) != 3 {
		fmt.Println("*ERROR* Invalid parameter count when call memset !!!")
		return false, ERR_PARAM_COUNT
	}

	pos     := vm.envFunc.envFuncParam[0]
	element := vm.envFunc.envFuncParam[1]
	count   := vm.envFunc.envFuncParam[2]
	vmLen   := uint64(len(vm.memory))
	if pos >= vmLen || pos + count >= vmLen {
		fmt.Println("VM::memset *ERROR* Out of bound")
		log.Infof("*ERROR* Out of bound \n")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	tempMem := make([]byte, count)
	var i uint64 = 0
	for ; i < count; i++ {
		tempMem[i] = byte(element)
	}

	copy(vm.memory[pos:pos + count], tempMem)

	if vm.envFunc.envFuncRtn {
		vm.pushInt32(int32(pos))
	}

	return true, nil
}

//void    *memcpy(void * destination, const void * source, size_t num);
func memcpy(vm *VM) (bool, error) {
	params := vm.envFunc.envFuncParam
	if len(params) != 3 {
		return false, ERR_PARAM_COUNT
	}

	dst    := params[0]
	src    := params[1]
	length := params[2]
	vmLen  := uint64(len(vm.memory))
	if dst >= vmLen || src >= vmLen || src + length >= vmLen {
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	if dst < src && dst + length > src {
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	copy(vm.memory[dst:dst + length], vm.memory[src:src + length])
	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(dst))
	}

	return true, nil
}

//uint32_t strcat_s(unsigned char * strDes, uint32_t size ,const unsigned char * strSrc);
func strcat_s(vm *VM) (bool, error) {
	params := vm.envFunc.envFuncParam
	if len(params) != 3 {
		return false, ERR_PARAM_COUNT
	}

	dst      := params[0]
	totalLen := params[1]
	src      := params[2]

	dstLen    := vm.StrLen(dst)
	srcLen    := vm.StrLen(src)
	dstPoint  := dst      + dstLen
	remindLen := totalLen - dstLen

	if remindLen < srcLen + 1 {
		if vm.envFunc.envFuncRtn {
			vm.pushUint32(uint32(VM_ERROR_OUT_OF_MEMORY))
		}

		return true, nil
	}

	copy(vm.memory[dstPoint:dstPoint + srcLen],vm.memory[src:src + srcLen])
	vm.memory[dstPoint + srcLen] = 0
	if vm.envFunc.envFuncRtn {
		vm.pushUint32(uint32(VM_NOERROR))
	}

	return true, nil
}

//uint32_t strcpy_s(unsigned char * strDes, uint32_t size ,const unsigned char * strSrc);
func strcpy_s(vm *VM) (bool, error) {
	params := vm.envFunc.envFuncParam
	if len(params) != 3 {
		return false, ERR_PARAM_COUNT
	}

	dst      := params[0]
	totalLen := params[1]
	src      := params[2]

	srcLen    := vm.StrLen(src)
	if totalLen < srcLen + 1 {
		if vm.envFunc.envFuncRtn {
			vm.pushUint32(uint32(VM_ERROR_OUT_OF_MEMORY))
		}

		return true, nil
	}

	copy(vm.memory[dst:dst + srcLen],vm.memory[src:src + srcLen])
	vm.memory[dst + srcLen] = 0

	if vm.envFunc.envFuncRtn {
		vm.pushUint32(uint32(VM_NOERROR))
	}

	return true, nil
}

//bool     isAccountExist(unsigned char * account);
func isAccountExist(vm *VM) (bool, error) {
	params := vm.envFunc.envFuncParam
	if len(params) != 1 {
		return false, ERR_PARAM_COUNT
	}

	contractCtx := vm.GetContract()
	pos         := uint64(params[0])
	length      := vm.StrLen(pos)
	accountNameByte , err := Convert(vm , uint64(pos) , uint64(length))
	if err != nil {
		return true, nil
	}
	accountName := BytesToString(accountNameByte)

	if contractCtx == nil || contractCtx.RoleIntf == nil {
		log.Infof("*ERROR* param is empty when call isAccountExist !!! ")
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_FALSE))
		}
		return true, nil
	}

	accountObj, err := contractCtx.RoleIntf.GetAccount(accountName)
	if err != nil {
		log.Infof("*ERROR* Failed to get account by name !!! ", err.Error())
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_FALSE))
		}
		return true, nil
	}

	if strings.Compare(accountObj.AccountName,accountName) != 0 {
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_FALSE))
		}
		return true, nil
	}

	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(VM_TRUE))
	}
	return true, nil
}

//void    *malloc(size_t size);
func malloc(vm *VM) (bool, error) {
	params  := vm.envFunc.envFuncParam
	if len(params) != 1 {
		return false, ERR_PARAM_COUNT
	}

	size := uint64(params[0])

	index, err := vm.getStoragePos(size, Unknown)
	if err != nil {
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}

	vm.RecoverContext()
	if vm.envFunc.envFuncRtn {
		vm.pushUint64(uint64(index))
	}

	return true, nil
}

func getMethodJs(vm *VM) (bool, error) {
	//
	envFunc := vm.envFunc
	params  := envFunc.envFuncParam
	if len(params) != 1 {
		return false, ERR_PARAM_COUNT
	}

	pos     := uint64(params[0])
	fmt.Println("vm::getMethodJs pos: = ",pos)
	/*
	var pos uint64 = 0
	var err error
	contractCtx := vm.GetContract()
	if pos, err = vm.StorageData(contractCtx.Trx.Method); err != nil {
		if vm.envFunc.envFuncRtn {
			vm.pushUint64(uint64(VM_NULL))
		}
		return true, nil
	}
	fmt.Println("VM::getMethodJs contractCtx.Trx.Method: ",string(contractCtx.Trx.Method)," , pos: ",pos)
	if vm.envFunc.envFuncRtn {
		vm.pushUint64(pos)
	}
	*/

	return true, nil
}