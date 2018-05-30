# bcli

```
$ bcli
Usage:
  newaccount -name NAME -pubkey PUBKEY         - Create a New account
  getaccount -name NAME                        - Get account balance
  transfer -from FROM -to TO -amount AMOUNT    - Transfer BTO from FROM account to TO
  deploycode -contract NAME -wasm PATH         - Deploy contract NAME from .wasm file
  deployabi -contract NAME -abi PATH           - Deploy contract ABI from .abi file

```

## newaccount
```
$ bcli newaccount -name xxxxx -pubkey xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## transfer
```
$ bcli transfer -from user1 -to user2 -amount 1000
```

## deploycode
```
$ bcli deploycode -contract xxxx -wasm wasm_file_path
```
