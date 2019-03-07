function set(table,key,value){
    if(typeof table != 'string'){
        console.log('table mast be string')
        return
    }

    if(typeof key != 'string'){
        console.log('key mast be string')
        return
    }

    if(typeof value != 'string'){
        console.log('value mast be string')
        return
    }
    return setBinValue(table,table.length,key,key.length,value,value.length)
}

function get(contract,table,key){
    if(typeof contract != 'string'){
        console.log('contract mast be string')
        return
    }

    if(typeof table != 'string'){
        console.log('table mast be string')
        return
    }

    if(typeof key != 'string'){
        console.log('key mast be string')
        return
    }

    return getBinValue(contract,contract.length,table,table.length,key,key.length)
}

function del(table,key){
    if(typeof table != 'string'){
        throw 'table mast be string'
        return
    }

    if(typeof key != 'string'){
        throw 'key mast be string'
        return
    }

    return removeBinValue(table,table.length,key,key.length)
}

module.exports = {
    set:set,
    get:get,
    del:del
}