#!/bin/bash

# add path variables
export NODEJS_HOME=/usr/lib/nodejs/node-v8.9.4
export PATH=$NODEJS_HOME/bin:$PATH
export PATH=$PATH:/usr/local/go/bin:/home/pi/go/bin
export GOPATH=/home/pi/go

# pull latest
cd $GOPATH/src/github.com/de0gee/basestation
git pull

# restart bluetooth
sudo service bluetooth restart

cd $GOPATH/src/github.com/de0gee/basestation/realtime-client 
yarn install
nohup yarn run start >/tmp/client.log 2>&1 &

cd $GOPATH/src/github.com/de0gee/basestation/realtime-server
go get -u -v ./...
go build
nohup sudo ./realtime-server >/tmp/server.log 2>&1 &
