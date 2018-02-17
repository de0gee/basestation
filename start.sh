#!/bin/bash

go get -u -v github.com/de0gee/basestation/...

sudo service bluetooth restart
sleep 3

cd $GOPATH/src/github.com/de0gee/basestation/realtime-client 
npm install
nohup npm run start >/tmp/client.log 2>&1 &

cd $GOPATH/src/github.com/de0gee/basestation/realtime-server
go build
nohup sudo ./realtime-server >/tmp/server.log 2>&1 &
