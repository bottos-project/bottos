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
var console = require('console')
var Storage = require('Storage')
var BTPack = require('BTPack')
var BTUnpack = require("BTUnpack")

var Lib = {}
Lib.getParams = function(){
    var paramBuf = getParam()
    var ABI = getAbi()
    var Abi = JSON.parse(ABI)
    var method = getMethod()
    var abi = Abi[method]
    return BTUnpack(paramBuf,abi)
}

Lib.getAbi = function(){
    var ABI = getAbi()
    var Abi = JSON.parse(ABI)
    var method = getMethod()
    var abi = Abi[method]
    return abi
}

Lib.getPack = function(obj){
    var abi = getCurrentAbi()
    var packBuf = BTPack(obj,abi)
    return buf2hex(packBuf)
}

Lib.getUnpack = function(packstr){
    var abi = getCurrentAbi()
    var buf = stringToBuffer(packstr)
    return BTUnpack(buf,abi)
}

var getCurrentAbi = function(){
    var ABI = getAbi()
    var Abi = JSON.parse(ABI)
    var method = getMethod()
    var abi = Abi[method]
    return abi
}

var buf2hex = function buf2hex(buffer) {
    return Array.prototype.map.call(new Uint8Array(buffer), function (x) {
        return ('00' + x.toString(16)).slice(-2);
    }).join('');
};

module.exports = {
    console:console,
    Lib:Lib,
    Storage:Storage
}