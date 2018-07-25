
### Build Bottos Chain image

git clone https://github.com/bottos-project/bottos.git --recursive  --depth 1


cd bottos/docker

docker build . -t bottos/bottoschain

### Start botNode docker

docker run -p 8080:8080 -t bottos/bottoschain


### Get chain info

curl -s -d "service=bottos&method=Chain.GetInfo&request={}" -H'content-type: application/x-www-form-urlencoded'  -X POST  http://127.0.0.1:8080/rpc 
#### Sample output:

{"result":{"head_block_num":45910,"head_block_hash":"fd7f6c05de590e334876183011b5791918ab99cee62bdb73b66d8854ddf4de2d","head_block_time":1532483234,"head_block_delegate":"delta","cursor_label":3723812397,"last_consensus_block_num":45885}}
