
var objstr = 'johntest'
var keystr = 'keyname111'
var valuestr = 'valuename222'

var obj = {
    obj:objstr,
    objLen:objstr.length,
    key:keystr,
    keyLen:keystr.length,
    value:valuestr,
    valueLen:valuestr.length
}

var abi = {
    type:'object',
    obj:{type:'string'},
    objLen:{type:'uint8'},
    key:{type:'string'},
    keyLen:{type:'uint8'},
    value:{type:'string'},
    valueLen:{type:'uint8'}
}
function userreg(param1,param2){
    var packBuf = BTPack(obj,abi)
    var unpack = BTUnpack(packBuf,abi)
    // print(JSON.stringify(unpack))

    setBinValue(unpack.obj,unpack.objLen,unpack.key,unpack.keyLen,unpack.value,unpack.valueLen)
    print('========user regist complate\n')
    return param1 + param2
}
