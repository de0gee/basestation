#!/bin/bash

# add path variables
export PATH=$PATH:/usr/local/go/bin:/home/pi/go/bin
export GOPATH=/home/pi/go

# restart bluetooth
sudo service bluetooth restart

# update the basestation
cd $GOPATH/src/github.com/de0gee/basestation
git pull
go build
nohup sudo ./basestation >/tmp/server.log 2>&1 &
