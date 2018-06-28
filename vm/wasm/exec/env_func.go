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
	"errors"
	log "github.com/cihub/seelog"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/contract"
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

	envFunc.Register("printi",           printi)
	envFunc.Register("prints",           prints)
	envFunc.Register("get_str_value",    getStrValue)
	envFunc.Register("set_str_value",    setStrValue)
	envFunc.Register("remove_str_value", removeStrValue)
	envFunc.Register("get_param",        getParam)
	envFunc.Register("call_trx",         callTrx)
	envFunc.Register("assert",           assert)
	envFunc.Register("getCtxName",       getCtxName)

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

func getStrValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params  := envFunc.envFuncParam
	if len(params) != 8 {
		return false, errors.New("parameter count error while call getStrValue")
	}
	contractPos := int(params[0])
	contractLen := int(params[1])
	objectPos   := int(params[2])
	objectLen   := int(params[3])
	keyPos      := int(params[4])
	keyLen      := int(params[5])
	valueBufPos := int(params[6])
	valueBufLen := int(params[7])

	// length check

	contract := make([]byte, contractLen)
	copy(contract, vm.memory[contractPos:contractPos+contractLen])

	object   := make([]byte, objectLen)
	copy(object, vm.memory[objectPos:objectPos+objectLen])

	key      := make([]byte, keyLen)
	copy(key, vm.memory[keyPos:keyPos+keyLen])

	log.Infof(string(contract), len(contract), string(object), len(object), string(key), len(key))
	value, err := contractCtx.ContractDB.GetStrValue(string(contract), string(object), string(key))

	valueLen := 0
	if err == nil {
		valueLen = len(value)
		// check buf len
		if valueLen <= valueBufLen {
			copy(vm.memory[valueBufPos:valueBufPos+valueLen], []byte(value))
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

	log.Infof("VM: from contract:%v, method:%v, func get_test_str:(contract=%v, objname=%v, key=%v, value=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, contract, object, key, value)

	return true, nil
}

func setStrValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 6 {
		return false, errors.New("parameter count error while call setStrValue")
	}
	objectPos := int(params[0])
	objectLen := int(params[1])
	keyPos    := int(params[2])
	keyLen    := int(params[3])
	valuePos  := int(params[4])
	valueLen  := int(params[5])

	// length check

	object := make([]byte, objectLen)
	copy(object, vm.memory[objectPos:objectPos+objectLen])

	key    := make([]byte, keyLen)
	copy(key, vm.memory[keyPos:keyPos+keyLen])

	value  := make([]byte, valueLen)
	copy(value, vm.memory[valuePos:valuePos+valueLen])

	log.Infof(string(object), len(object), string(key), len(key), string(value), len(value))
	err := contractCtx.ContractDB.SetStrValue(contractCtx.Trx.Contract, string(object), string(key), string(value))

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

	log.Infof("VM: from contract:%v, method:%v, func setStrValue:(objname=%v, key=%v, value=%v)\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, object, key, value)

	return true, nil
}

func removeStrValue(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params := envFunc.envFuncParam
	if len(params) != 4 {
		return false, errors.New("parameter count error while call removeStrValue")
	}
	objectPos := int(params[0])
	objectLen := int(params[1])
	keyPos    := int(params[2])
	keyLen    := int(params[3])

	// length check

	object := make([]byte, objectLen)
	copy(object, vm.memory[objectPos:objectPos+objectLen])

	key := make([]byte, keyLen)
	copy(key, vm.memory[keyPos:keyPos+keyLen])

	log.Infof(string(object), len(object), string(key), len(key))
	err := contractCtx.ContractDB.RemoveStrValue(contractCtx.Trx.Contract, string(object), string(key))

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

func printi(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()
	value := vm.envFunc.envFuncParam[0]
	fmt.Printf("VM: from contract: %v, method: %v, func printi: %v\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, value)
	log.Infof("VM: from contract:%v, method:%v, func printi: %v\n", contractCtx.Trx.Contract, contractCtx.Trx.Method, value)

	return true, nil
}

func prints(vm *VM) (bool, error) {

	pos := vm.envFunc.envFuncParam[0]
	len := vm.envFunc.envFuncParam[1]

	value := make([]byte, len)
	copy(value, vm.memory[pos:pos+len])

	BytesToString(value)
	param := string(value)
	fmt.Println("VM: func prints: ", param)
	log.Infof("VM: func prints: %v\n", param)
	return true, nil

}

func getParam(vm *VM) (bool, error) {
	contractCtx := vm.GetContract()

	envFunc := vm.envFunc
	params  := envFunc.envFuncParam
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

	return true, nil
}

func callTrx(vm *VM) (bool, error) {

	envFunc := vm.envFunc
	params  := envFunc.envFuncParam

	if len(params) != 4 {
		return false, errors.New("*ERROR* Parameter count error while call memcpy")
	}

	cPos := int(params[0])
	mPos := int(params[1])
	pPos := int(params[2])
	pLen := int(params[3])

	contrx := BytesToString(vm.memory[cPos : cPos+vm.memType[uint64(cPos)].Len-1])
	method := BytesToString(vm.memory[mPos : mPos+vm.memType[uint64(mPos)].Len-1])
	//the bytes after msgpack.Marshal
	param := vm.memory[pPos : pPos+pLen]
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
		vm.pushUint32(uint32(0))
	}

	return true, nil
}

func assert(vm *VM) (bool, error) {
	envFunc := vm.envFunc
	params := envFunc.envFuncParam

	cond := int(params[0])
	if cond != 1 {
		errStr := "*ERROR* Failed to execute contract code !!!"
		panic(errStr)
	}

	return true, nil
}

func getCtxName(vm *VM) (bool, error) {

	ctxName    := vm.contract.Trx.Contract
	ctxNameLen := uint64(len(ctxName))

	pos := vm.envFunc.envFuncParam[0]
	len := vm.envFunc.envFuncParam[1]
	if len < ctxNameLen {
		log.Infof("*ERROR* Invaild string length \n")
		if vm.envFunc.envFuncRtn {
			vm.pushInt32(int32(1))
		}
	}

	copy(vm.memory[pos:pos+len], []byte(ctxName))

	if vm.envFunc.envFuncRtn {
		vm.pushInt32(int32(0))
	}

	return true, nil
}