# cmd

```
$ cmd
Usage:
  newaccount -name NAME -pubkey PUBKEY - Create a New account
  transfer -from FROM -to TO -amount AMOUNT - transfer bottos from FROM account to TO
  deploycode -contract NAME -wasm PATH - deploy contract NAME from .wasm file

```

## newaccount
```
$ cmd newaccount -name xxxxx -pubkey xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## transfer
```
$ cmd transfer -from user1 -to user2 -amount 1000
```

## deploycode
```
$ cmd deploycode -contract xxxx -wasm wasm_file_path
```
