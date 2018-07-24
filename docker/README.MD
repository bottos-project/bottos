Build eos image
git clone https://github.com/bottos-project/bottos.git --recursive  --depth 1
cd bottos/docker
docker build . -t bottos/bottoschain

Start botNode docker
docker run -p 8080:8080 -t bottos/bottoschain


Get chain info
curl -s -d "service=bottos&method=Chain.GetInfo&request={}" -H'content-type: application/x-www-form-urlencoded'  -X POST  http://127.0.0.1:8080/rpc 
