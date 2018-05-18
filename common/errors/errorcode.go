package errors

// import "fmt"


type ErrCode uint32

const (
	ErrNoError                ErrCode = 0
	
	ErrTrxPendingNumLimit     ErrCode = 45001
	ErrTrxSignError           ErrCode = 45002
	ErrTrxAccountError        ErrCode = 45003
	ErrTrxLifeTimeError       ErrCode = 45004
	ErrTrxUniqueError         ErrCode = 45005
	ErrTrxChainMathError      ErrCode = 45006
	ErrTrxContractHanldeError ErrCode = 45007

	ErrInvalid              ErrCode = 0xFFFFFFFF
)




var (

	aaa = map[ErrCode]string{
		ErrTrxPendingNumLimit     : "push trx: "    +    "check Pending pool max num error",
		ErrTrxSignError           : "push trx: "    +    "check signature error",
		ErrTrxAccountError        : "push trx: "    +    "check account valid error",
		ErrTrxLifeTimeError       : "push trx: "    +    "check life time error",
		ErrTrxUniqueError         : "push trx: "    +    "check trx unique error",
		ErrTrxChainMathError      : "push trx: "    +    "check match chain error",
		ErrTrxContractHanldeError : "push trx: "    +    "process contract error",

	  }
)


func GetCodeString(errorCode ErrCode) string {
	return  aaa[errorCode]
}



  