Array.prototype.toUint8Array = function(){
    return new Uint8Array(this)
}
Array.prototype.from = function(arrayBuffer){
    return arrayBuffer.toArray()
}

Array.prototype.fillZero = function(n){
    var arr = []
    for(i=0;i<n;i++){
        arr[i] = 0
    }
    return arr;
}

Uint8Array.prototype.toArray = function(){
    var arr = []
    var obj = this
    Object.keys(obj).forEach(function(key){
        var value = obj[key]
        arr.push(value)
    })
    return arr
}

Uint8Array.prototype.slice = function(i,j){
    var arr = this.toArray()
    var arrTemp = arr.slice(i,j)
    return new Uint8Array(arrTemp)
}


Duktape.modSearch = function (id) {
    var name;
    var src;
    var found = false;

    /* ECMAScript check. */
    name = './vm/duktape/js_lib/' + id + '.js';
    /* name = id + '.js'; */
    src = load_js(name);
    if (typeof src === 'string') {
        print('loaded ECMAScript:', name);
        found = true;
    }

    /* Must find either a DLL or an ECMAScript file (or both) */
    if (!found) {
        throw new Error('module not found: ' + id);
    }

    /* For pure C modules, 'src' may be undefined which is OK. */
    return src;
}

var BTPack = require('BTPack')
var BTUnpack = require('BTUnpack')
var console = require('console')

var previousabi= {
    account:{type:'string'},
    age:{type:'uint32'}
}

/* 
V1: "12345",
V2: 888,
V3: "6789",
V4: 999,
:
dc0004da00053132333435ce00000378da000436373839ce000003e7 */
function reguser(){

    var accountBuf = getParam()
    var abiget = getAbi()

    console.log(accountBuf)
    console.log(JSON.stringify(abiget))

    var accountInfo = BTUnpack(accountBuf,JSON.parse(abiget))

    console.log(JSON.stringify(accountInfo))

    var table = 'usermng'
    var key = 'userreginfo'

    var accountStr = accountBuf.toString();

    setBinValue(table,table.length,key,key.length,accountStr,accountStr.length)

    console.log('regist user complate\n');


    return 0
}


function login(accountBuf){
    var contract = "usermng"
    var key = "keytest"
    var obj = "userreginfo"
    var result = getBinValue(contract,contract.length,key,key.length,obj,obj.length)
    var arrayBuf = result.split(",")
    var unpackObj = BTUnpack(arrayBuf,abi)
    console.log(JSON.stringify(unpackObj))
    console.log('login complate\n');
    return 0
}
