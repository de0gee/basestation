#!/bin/bash

# add path variables
export PATH=$PATH:/usr/local/go/bin:/home/pi/go/bin
export GOPATH=/home/pi/go

# restart bluetooth
sudo service bluetooth restart

# update the basestation
cd $GOPATH/src/github.com/de0gee/basestation
git pull
rm -rf basestation-linux*
curl -s https://api.github.com/repos/de0gee/basestation/releases/latest  | grep 'arm' | grep '.tar.gz' | grep 'http' | cut -d '"' -f 4 | wget -qi -
tar -xvzf basestation-linux-arm-6.tar.gz
chmod +x basestation-linux-arm-6
nohup sudo ./basestation-linux-arm-6 >/tmp/server.log 2>&1 &
