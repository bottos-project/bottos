package errors

// import "fmt"


type ErrCode uint32

const (
	ErrNoError                ErrCode = 0
	
	ErrTrxPendingNumLimit     ErrCode = 10001
	ErrTrxSignError           ErrCode = 10002
	ErrTrxAccountError        ErrCode = 10003
	ErrTrxLifeTimeError       ErrCode = 10004
	ErrTrxUniqueError         ErrCode = 10005
	ErrTrxChainMathError      ErrCode = 10006
	ErrTrxContractHanldeError ErrCode = 10007

	ErrApiTrxNotFound		  ErrCode = 10101
	ErrApiBlockNotFound		  ErrCode = 10102
	ErrApiQueryChainInfoError ErrCode = 10103
	ErrApiAccountNotFound	  ErrCode = 10104
	ErrApiObjectNotFound	  ErrCode = 10105

	ErrInvalid              ErrCode = 0xFFFFFFFF
)




var (

	aaa = map[ErrCode]string{
		ErrTrxPendingNumLimit     : "push trx: "    		+    "check Pending pool max num error",
		ErrTrxSignError           : "push trx: "    		+    "check signature error",
		ErrTrxAccountError        : "push trx: "    		+    "check account valid error",
		ErrTrxLifeTimeError       : "push trx: "    		+    "check life time error",
		ErrTrxUniqueError         : "push trx: "    		+    "check trx unique error",
		ErrTrxChainMathError      : "push trx: "    		+    "check match chain error",
		ErrTrxContractHanldeError : "push trx: "    		+    "process contract error",

		ErrApiTrxNotFound         : "query trx: "   		+    "trx not found",
		ErrApiBlockNotFound	      : "query block: " 		+    "block not found",
		ErrApiQueryChainInfoError : "query chain info: " 	+    "error",
		ErrApiAccountNotFound     : "query account: " 		+	 "account not found",
		ErrApiObjectNotFound      : "query object: " 		+	 "object not found",

	  }
)


func GetCodeString(errorCode ErrCode) string {
	return  aaa[errorCode]
}



  
