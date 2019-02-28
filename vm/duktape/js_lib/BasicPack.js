/*
  Copyright 2017~2022 The Bottos Authors
  This file is part of the Bottos Data Exchange Client
  Created by Developers Team of Bottos.

  This program is free software: you can distribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with Bottos. If not, see <http://www.gnu.org/licenses/>.
*/

var Long = require('long')
var BigNumber = require('bignumber')

var BIN16 = 0xc5  //197
var UINT8  = 0xcc  //204
var UINT16 = 0xcd  //205
var UINT32 = 0xce  //206
var UINT64 = 0xcf  //207
var STR16   = 0xda  //218
var ARRAY16 = 0xdc  //220

var LEN_INT32 = 4
var LEN_INT64 = 8

var MAX16BIT = 2 << (16 - 1)

var REGULAR_UINT7_MAX  = 2 << (7 - 1)
var REGULAR_UINT8_MAX  = 2 << (8 - 1)
var REGULAR_UINT16_MAX = 2 << (16 - 1)
var REGULAR_UINT32_MAX = 2 << (32 - 1)

var SPECIAL_INT8  = 32
var SPECIAL_INT16 = 2 << (8 - 2)
var SPECIAL_INT32 = 2 << (16 - 2)
var SPECIAL_INT64 = 2 << (32 - 2)

var PackUint8 = function(value){
    var buf = new Uint8Array(2)
    buf[0]=UINT8;
    buf[1]=value;
    return buf
}

var PackUint16 = function(value){
    var buf = new Uint8Array(3)
    buf[0]=UINT16
    buf[1]=value>>8
    buf[2]=value
    return buf
}

var PackUint32 = function(value){
    var buf = new Uint8Array(5)
    buf[0]=UINT32
    buf[1]=value>>24
    buf[2]=value>>16
    buf[3]=value>>8
    buf[4]=value
    return buf
}

var PackUint64 = function(n){
  var num = Long.fromNumber(Number(n))
  var h = num.getHighBitsUnsigned()
  var l = num.getLowBitsUnsigned()

  var buf = new Uint8Array(9)
  buf[0]=UINT64
  buf[1]=h>>24
  buf[2]=h>>16
  buf[3]=h>>8
  buf[4]=h

  buf[5]=l>>24
  buf[6]=l>>16
  buf[7]=l>>8
  buf[8]=l
  return buf
}

var PackUint256 = function(num){
  var bigNumber = new BigNumber(num)
  var hexNumber = bigNumber.toString(16)
  var hexNumberArr = hexNumber.split('')
  var hexLength = hexNumber.length
  var zeroLength = 64 - hexLength
  var zeroArr = []
  zeroArr.fillZero(zeroLength)    // 前面补0
  var dataArr = zeroArr.concat(hexNumberArr)
  var result = arrayChunk(dataArr,2)
  var arrayBuf = new ArrayBuffer(32)
  for(var j = 0;j<result.length;j++){
    arrayBuf[j] = result[j]
  }

  var arrbuf = PackBin16(arrayBuf)
  var length = arrbuf.byteLength
  var buf = new Uint8Array(length)
  for(var i = 0;i<length;i++){
    buf[i] = arrbuf[i]
  }

  return buf
}

var PackBin16 = function(byteArray){
  var byteLen = byteArray.byteLength
  var len = byteLen + 3
  var bytes = new ArrayBuffer(len)
  bytes[0] = BIN16
  bytes[1] = byteLen>>8
  bytes[2] = byteLen
  for(var i = 0;i<byteLen;i++){
    bytes[i+3] = byteArray[i]
  }
  return bytes
}

var PackStr16 = function(str){
  str = convertUnicode2Utf8(str)
  var len = str.length
  var byteLen = len + 3
  var bytes = new Uint8Array(byteLen)
  bytes[0] = STR16
  bytes[1] = len >> 8
  bytes[2] = len
  for(var i = 0;i<len;i++){
    bytes[i+3] = str[i]
  }
  return bytes
}

var PackArraySize = function(length){
  var size = new Uint8Array(3)
  size[0] = ARRAY16
  size[1] = length>>8
  size[2] = length
  return size
}

var convertUnicode2Utf8 = function(str){
  var isGetBytes=true
  var back = [];
  var byteSize = 0;
  for (var i = 0; i < str.length; i++) {
      var code = str.charCodeAt(i);
      if (0x00 <= code && code <= 0x7f) {
            byteSize += 1;
            back.push(code);
      } else if (0x80 <= code && code <= 0x7ff) {
            byteSize += 2;
            back.push((192 | (31 & (code >> 6))));
            back.push((128 | (63 & code)))
      } else if ((0x800 <= code && code <= 0xd7ff)
              || (0xe000 <= code && code <= 0xffff)) {
            byteSize += 3;
            back.push((224 | (15 & (code >> 12))));
            back.push((128 | (63 & (code >> 6))));
            back.push((128 | (63 & code)))
      }
    }
    for (i = 0; i < back.length; i++) {
        back[i] &= 0xff;
    }
    if (isGetBytes) {
        return back
    }
    if (byteSize <= 0xff) {
        return [0, byteSize].concat(back);
    } else {
        return [byteSize >> 8, byteSize & 0xff].concat(back);
    }
}

// 将数组每两个元素组成一个新的字符串作为一个元素存储
var arrayChunk = function(dataArr,colomns){
  var result = []
  for( var i = 0, len = dataArr.length; i < len; i += colomns ) {
    var tempArr = dataArr.slice(i,i+colomns)
    var tempStr = tempArr.join('')
    result.push(Number.parseInt(tempStr,16))
  }
  return result
}

var BasicPack = {
  PackUint8,
  PackUint16,
  PackUint32,
  PackUint64,
  PackBin16,
  PackStr16,
  PackArraySize,
  PackUint256
}

module.exports = BasicPack