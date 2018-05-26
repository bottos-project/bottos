echo cmd newaccount -name ${1} -pubkey 7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF

./cmd newaccount -name ${1} -pubkey 7QBxKhpppiy7q4AcNYKRY2ofb3mR5RP8ssMAX65VEWjpAgaAnF

echo cmd deploycode -contract ${1} -wasm ./${1}.wasm

./cmd deploycode -contract ${1} -wasm ./${1}.wasm
