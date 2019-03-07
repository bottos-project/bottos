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
    name = './vm/duktape/js_lib' + id + '.js';
    src = load_js(name);
    if (typeof src === 'string') {
        print('loaded ECMAScript:', name);
        found = true;
    };
    if (!found) {
        throw new Error('module not found: ' + id);
    };
    return src;
}
