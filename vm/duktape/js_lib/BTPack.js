var BSPack = require('./BasicPack')

var BTPack = function(obj,abi){
  var keys = Object.keys(obj)
  var size = keys.length
  var packBuf = []
  packBuf = packBuf.concat(BSPack.PackArraySize(size).toArray())
  for(var i = 0;i < keys.length;i++){
    var key = keys[i]
    var keyAbi = abi[key]
    var value = obj[key]
    var type = keyAbi && keyAbi.type
    if(type=='string'){
      // console.log('string')
      packBuf = packBuf.concat(BSPack.PackStr16(value).toArray())
    }else if(type=='uint8'){
      // console.log('uint8')
      packBuf = packBuf.concat(BSPack.PackUint8(value).toArray())
    }else if(type=='uint16'){
      packBuf = packBuf.concat(BSPack.PackUint16(value).toArray())
    }else if(type=='uint32'){
      // console.log('uint32')
      packBuf = packBuf.concat(BSPack.PackUint32(value).toArray())
    }else if(type == 'uint64'){
      // console.log('uint64')
      packBuf = packBuf.concat(BSPack.PackUint64(value).toArray())
    }else if(type=='uint256'){
      packBuf = packBuf.concat(BSPack.PackUint256(value).toArray())
    }else if(type=='array'){
      // console.log('array')
      packBuf = packBuf.concat(BSPack.PackBin16(value).toArray())
    }else if(type=='object'){
      // console.log('object')
      packBuf = packBuf.concat(BTPack(value,keyAbi).toArray())
    }else{
      print('Invalid type')
    }
  }
  return packBuf
}

module.exports = BTPack
