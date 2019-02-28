package duktape
/*
#include <stdlib.h>
#include <stdint.h>
#cgo CFLAGS: -Ipub/
#cgo CFLAGS: -Icommon/
#cgo LDFLAGS: -L ../../lib/ -lbottosduktape
#cgo LDFLAGS: -lm

#include "pub.h"
#include "types.h"
extern JSExeErrorCode  process(Context *contractCtx);

extern uint32_t StorageSaveFunc(char *a,char *b, char *c, char *d);

uint32_t StorageSaveFuncCgo(char *a,char *b, char *c, char *d)
{	
	return StorageSaveFunc(a, b, c ,d);
}

extern uint32_t StorageReadFunc(char *a,char *b, char *c, BinResult *binResult);
uint32_t StorageReadFuncCgo(char *a,char *b, char *c, BinResult *binResult)
{	
	return StorageReadFunc(a, b, c ,binResult);
}

extern uint32_t AddSubTrxFunc(char *a,char *b, char *c);
uint32_t AddSubTrxCgo(char* contractName, char* method, char* param)
{	
	return AddSubTrxFunc(contractName, method, param);
}

*/
import "C"

import (
	//"errors"
	//"strconv"
	"github.com/bottos-project/bottos/role"
	"github.com/bottos-project/bottos/common"
	"github.com/bottos-project/bottos/common/types"
	"github.com/bottos-project/bottos/config"
	log "github.com/cihub/seelog"
)

// TrxApplyService is to define a service for apply a transaction
var roleIntf role.RoleInterface	


type SubTrx struct {
	Contract    string // max length 21
	Method      string // max length 21
	Param       string
}

//var subTrxArray [config.DEFAUL_MAX_SUB_CONTRACT_NUM]types.Transaction;

var subTrxSlice []SubTrx = make([]SubTrx, config.DEFAUL_MAX_SUB_CONTRACT_NUM)
func InitDuktapeVm(roleIntfInput role.RoleInterface) {

	roleIntf =  roleIntfInput
}

func Process(contractCode []byte, contractAbi []byte, trx *types.Transaction) (uint32, []*types.Transaction){
	contractNameCtring := C.CString(trx.Contract)
	methodCtring := C.CString(trx.Method)
	paramsCtring := C.CString(common.BytesToHex(trx.Param[:]))
	contractCtring := C.CString(string(contractCode[:]))
	//contractCtring := C.CString(common.BytesToHex(contractCode[:]))
	contractAbiCtring := C.CString(string(contractAbi[:]))
	senderCtring := C.CString(trx.Sender)
	var employee C.Context = C.Context{senderCtring, contractNameCtring, contractCtring, contractAbiCtring, methodCtring, paramsCtring}
	
	subTrxSlice = subTrxSlice[:0]

	var rtnValue C.JSExeErrorCode = C.process(&employee);

	if (rtnValue != C.VM_EXE_SUCC) {
		log.Errorf("exe failed error code %d", rtnValue)
		//var errString string = strconv.Itoa((int)(rtnValue))		
		//return  errors.New(errString), nil
		return  uint32(rtnValue), nil
	}

	if (len(subTrxSlice) == 0) {
		return 0, nil
	} else {
		log.Infof("slice num is %d", len(subTrxSlice))

		value := make([]*types.Transaction, len(subTrxSlice))
		value = value[:0]
		for i, v:= range subTrxSlice {
			log.Infof("slice[%d] = %v", i, v)

			trx := &types.Transaction{
				Version:     trx.Version,
				CursorNum:   trx.CursorNum,
				CursorLabel: trx.CursorLabel,
				Lifetime:    trx.Lifetime,
				Sender:      trx.Contract,
				Contract:    v.Contract,
				Method:      v.Method,
				Param:       ([]byte)(v.Param), //the bytes after msgpack.Marshal
				SigAlg:      trx.SigAlg,
				Signature:   []byte{},
			}

			value = append(value, trx)
		}
		
		return 0, value
	}
}

// gcc vm/duktape/duktapesrc/extras/module-duktape/duk_module_duktape.c  vm/duktape/duktapesrc/src/duktape.c  vm/duktape/engine/db.c vm/duktape/engine/process.c -Ivm/duktape/common -Ivm/duktape/duktape/extras/module-duktape -Ivm/duktape/duktapesrc/src -Ivm/duktape/engine -Ivm/duktape/pub  -fPIC -shared -o libduktapebottos.so -lm