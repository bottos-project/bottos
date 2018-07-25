#!/bin/sh
cd /go/bin
echo "Start conusl..."

nohup consul agent -dev > consul.log 2>&1 &
	sleep 3
echo "Start micro..."
nohup micro api > micro.log 2>&1 &

cd /go/src/github.com/bottos-project/bottos/

echo "Start bottosChain..."
echo "It will take thirty sencond..."
nohup ./bottos > bottos.log 2>&1 &
tail -F bottos.log
#tail -F core.log
