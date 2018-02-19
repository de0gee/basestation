#!/bin/bash

cd $GOPATH/src/github.com/de0gee/basestation
git pull

sudo service bluetooth restart

cd $GOPATH/src/github.com/de0gee/basestation/realtime-client 
yarn install
nohup yarn run start >/tmp/client.log 2>&1 &

cd $GOPATH/src/github.com/de0gee/basestation/realtime-server
go build
nohup sudo ./realtime-server >/tmp/server.log 2>&1 &
