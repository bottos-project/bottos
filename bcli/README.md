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
