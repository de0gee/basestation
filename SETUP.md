# Instructions for making a turn-key image

```
sudo apt-get update
sudo apt-get dist-upgrade -y
sudo apt-get install -y dnsmasq hostapd vim g++ sqlite3 python3-flask
```

## Install node 

```
wget https://nodejs.org/dist/v8.9.4/node-v8.9.4-linux-armv6l.tar.xz
sudo mkdir /usr/lib/nodejs
sudo tar -xJvf node-v8.9.4-linux-armv6l.tar.xz -C /usr/lib/nodejs 
rm -rf node-v8.9.4-linux-armv6l.tar.xz
sudo mv /usr/lib/nodejs/node-v8.9.4-linux-armv6l /usr/lib/nodejs/node-v8.9.4
echo 'export NODEJS_HOME=/usr/lib/nodejs/node-v8.9.4' >> ~/.profile
echo 'export PATH=$NODEJS_HOME/bin:$PATH' >> ~/.profile
source ~/.profile
```

## Install Go

```
# install Go
wget https://dl.google.com/go/go1.9.4.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.9.4.*
rm go1.9.*
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >>  ~/.profile
echo 'export GOPATH=$HOME/go' >>  ~/.profile
source ~/.profile
```

# Install Hostapd

```
sudo systemctl stop dnsmasq && sudo systemctl stop hostapd

echo 'interface wlan0
static ip_address=192.168.4.1/24' | sudo tee --append /etc/dhcpcd.conf
sudo mv /etc/dnsmasq.conf /etc/dnsmasq.conf.orig  

sudo systemctl daemon-reload
sudo systemctl restart dhcpcd

echo 'interface=wlan0
dhcp-range=192.168.4.2,192.168.4.20,255.255.255.0,24h' | sudo tee --append /etc/dnsmasq.conf

echo 'interface=wlan0
driver=nl80211
ssid=NameOfNetwork
hw_mode=g
channel=7
wmm_enabled=0
macaddr_acl=0
auth_algs=1
ignore_broadcast_ssid=0
wpa=2
wpa_passphrase=AardvarkBadgerHedgehog
wpa_key_mgmt=WPA-PSK
wpa_pairwise=TKIP
rsn_pairwise=CCMP' | sudo tee --append /etc/hostapd/hostapd.conf

echo 'DAEMON_CONF="/etc/hostapd/hostapd.conf"' | sudo tee --append /etc/default/hostapd

sudo systemctl start hostapd && sudo systemctl start dnsmasq
git clone https://github.com/de0gee/basestation.git
```

Add to cron

```
cd /home/pi/basestation/turnkey && /usr/bin/sudo /usr/bin/python3 server.py
```

