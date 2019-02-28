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

var UnpackUint8 = function(buf){
    var v1 = buf[1]
    return v1;
}

var UnpackUint16 = function(buf){
    var v1 = buf[1] << 8
    var v2 = buf[2]
    return v1+v2
}

var UnpackUint32 = function(buf){
    var v1 = buf[1] << 24
    var v2 = buf[2] << 16
    var v3 = buf[3] << 8
    var v4 = buf[4]
    return v1 + v2 + v3 + v4
}

var UnpackUint64 = function(buf){
    var v1 = buf[1] << 24
    var v2 = buf[2] << 16
    var v3 = buf[3] << 8
    var v4 = buf[4]
    var h = (v1 + v2 + v3 + v4)<<32

    var v5 = buf[5] << 24
    var v6 = buf[6] << 16
    var v7 = buf[7] << 8
    var v8 = buf[8]
    var l = v5 + v6 + v7 + v8

    // var value = h + l
    var longValue = Long.fromBits(l,h,true)
    var value = longValue.toString()
    return value
}

var UnpackStr16 = function(buf){
    var v1 = buf[1] << 8 
    var v2 = buf[2]
    var length = v1 + v2
    var str = ''
    for(var i = 3;i<length+3;i++){
        var s = String.fromCharCode(buf[i])
        str += s
    }
    return str
}

var UnpackBin16 = function(buf){
    var v1 = buf[1] << 8
    var v2 = buf[2]
    var length = v1 + v2
    var bytes = []
    for(var i = 3;i<length;i++){
        bytes.push(buf[i])
    }
    return bytes
}

var UnpackArraySize = function(buf){
    var v1 = buf[1] << 8
    var v2 = buf[2]
    return  v1 + v2
}

module.exports = {
    UnpackUint8,
    UnpackUint16,
    UnpackUint32,
    UnpackUint64,
    UnpackStr16,
    UnpackBin16,
    UnpackArraySize
}