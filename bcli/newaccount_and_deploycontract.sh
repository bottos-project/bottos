echo bcli newaccount -name ${1} -pubkey 7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF

./bcli newaccount -name ${1} -pubkey 7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF

echo bcli deploycode -contract ${1} -wasm ./${1}.wasm

./bcli deploycode -contract ${1} -wasm ./${1}.wasm
