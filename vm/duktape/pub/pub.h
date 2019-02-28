typedef struct  {
    char *sender;
    char *contractName;
    char *contractCode;
    char *contractAbi;
    char *method;
    char *params;
}Context;

typedef struct  {
    char *contractName;
    char *method;
    char *params;
}SubTrx;


#define MAX_SUB_TRX_NUM (10)

SubTrx SubTrxList[MAX_SUB_TRX_NUM];




/* contract exe error code */

typedef enum  {
        VM_EXE_SUCC = 0,
        /* CONTRACT_EXEC_ERROR START: 1 */

        /* CONTRACT_EXEC_ERROR START: 0xfff */
        VM_CTX_INIT_FAIL = 0x2001,
        VM_LOAD_CODE_FAIL,
        VM_JS_CALL_FAIL, 
        VM_JS_RTN_FAIL,
        VM_ADD_SUB_TRX_FAIL
}JSExeErrorCode;

JSExeErrorCode  process(Context *contractCtx);
