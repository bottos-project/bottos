var BSUnpack = require('BasicUnpack')
var Types = require('Types')

var BTUnpack = function(buf,abi){
    var obj = {}
    var keys = Object.keys(abi)
    var index = 3
    for(var i = 0;i<keys.length; i++){
        var key = keys[i]
        var keyAbi = abi[key]
        var abiType = keyAbi.type
        var currentBuf = []
        if(buf.length > index){
            currentBuf = buf.slice(index,buf.length)
        }
        type = currentBuf[0]
        if(abiType==Types.uint8){
            obj[key] = BSUnpack.UnpackUint8(currentBuf)
            index += 2
        }else if(abiType==Types.uint16){
            obj[key] = BSUnpack.UnpackUint16(currentBuf)
            index += 3
        }else if(abiType==Types.uint32){
            obj[key] = BSUnpack.UnpackUint32(currentBuf)
            index += 5
        }else if(abiType==Types.uint64){
            obj[key] = BSUnpack.UnpackUint64(currentBuf)
            index += 9
        }else if(abiType==Types.string){
            obj[key] = BSUnpack.UnpackStr16(currentBuf)
            index = currentBuf[2] + 3 + index
        }else if(abiType==Types.array){
            unPackBin16 = BSUnpack.UnpackBin16(currentBuf)
            obj[key] = unPackBin16
            index += currentBuf.byteLength
        }else if(abiType==Types.object){
            obj[key] = BTUnpack(currentBuf,keyAbi)
            index += 3
        }else{
            // console.log({type,abiType})
            print('Invalid type')
        }
    }
    return obj
}

module.exports = BTUnpack