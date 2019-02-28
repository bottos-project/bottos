package duktape

import (
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/config"
	log "github.com/cihub/seelog"
)

/*
#include <stdlib.h>
#include <stdint.h>
#cgo CFLAGS: -Icommon/
#include "types.h"
*/
import "C"

//export StorageSaveFunc
func StorageSaveFunc(contract *C.char, object *C.char, key *C.char, value *C.char) uint32 {
	// _, storage := getEngineByStorageHandler(uint64(uintptr(handler)))
	// if storage == nil {
	// 	logging.VLog().Error("Failed to get storage handler.")
	// 	return nil
	// }

	contractStr := C.GoString(contract)
	objectStr := C.GoString(object)
	keyStr := C.GoString(key)
	valueStr := C.GoString(value)

	temp, _ := common.HexToBytes(valueStr)

	err := roleIntf.SetBinValue(contractStr, objectStr, keyStr, temp)

	log.Infof("go: StorageSaveFunc contractStr %s, obj %s, key %s, %x\n", contractStr, objectStr, keyStr, temp)

	var result uint32 = 1
	if err != nil {
		result = 0
	}

	return result
}

//export StorageReadFunc
func StorageReadFunc(contract *C.char, object *C.char, key *C.char, binResult *C.BinResult) uint32 {
	// _, storage := getEngineByStorageHandler(uint64(uintptr(handler)))
	// if storage == nil {
	// 	logging.VLog().Error("Failed to get storage handler.")
	// 	return nil
	// }

	contractStr := C.GoString(contract)
	objectStr := C.GoString(object)
	keyStr := C.GoString(key)
	//valueStr := C.GoString(value);

	value, err := roleIntf.GetBinValue(contractStr, objectStr, keyStr)
	log.Infof("go: StorageReadFunc contractStr %s, objectStr %s, keyStr %s\n", contractStr, objectStr, keyStr)

	if err == nil {
		/* var templen C.uint32_t = 3 */ /* (uint32)(len(value)) */
		valueString := common.BytesToHex(value)
		binResult.valueLen = C.uint32_t(len(valueString))
		binResult.binValue = C.CString(valueString)
		log.Infof("storage read ok\n")
		return 1
	} else {
		log.Infof("storage read error\n")
		return 0
	}
}

//export AddSubTrxFunc
func AddSubTrxFunc(contractName *C.char, method *C.char, param *C.char) uint32 {

	contractNameStr := C.GoString(contractName)
	methodStr := C.GoString(method)
	paramStr := C.GoString(param)

	trx := SubTrx{
		Contract: contractNameStr,
		Method:   methodStr,
		Param:    paramStr,
	}

	if uint32(len(subTrxSlice)) < config.DEFAUL_MAX_SUB_CONTRACT_NUM {
		subTrxSlice = append(subTrxSlice, trx)
		log.Infof("go:AddSubTrxFunc add sub trx succ, current len %d", len(subTrxSlice))
		return 0
	} else {
		log.Infof("go:AddSubTrxFuncadd sub trx fail, current len %d", len(subTrxSlice))
		return 1
	}
}
