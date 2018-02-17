# Instructions for making a turn-key image

```
sudo apt-get update
sudo apt-get dist-upgrade -y
sudo apt-get install -y dnsmasq hostapd vim g++ sqlite3
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

sudo apt-get install -y python3-flask
```

Add to cron

```
cd /home/pi/basestation/turnkey && /usr/bin/python3 server.py
```


# Setup server (this should be scripted)

```
go get github.com/mholt/caddy/caddy

# Disable:

sudo sed -i '/DAEMON_CONF="\/etc/s/^/#/g' /etc/default/hostapd
sudo sed -i '/interface wlan0/s/^/#/g' /etc/dhcpcd.conf
sudo sed -i '/static ip_address=192.168.4.1\/24/s/^/#/g' /etc/dhcpcd.conf
sudo sed -i '/interface=wlan0/s/^/#/g' /etc/dnsmasq.conf
sudo sed -i '/dhcp-range=/s/^/#/g' /etc/dnsmasq.conf

# Enable:
sudo sed -i '/DAEMON_CONF="\/etc/s/^#//g' /etc/default/hostapd
sudo sed -i '/interface wlan0/s/^#//g' /etc/dhcpcd.conf
sudo sed -i '/static ip_address=192.168.4.1\/24/s/^#//g' /etc/dhcpcd.conf
sudo sed -i '/interface=wlan0/s/^#//g' /etc/dnsmasq.conf
sudo sed -i '/dhcp-range=/s/^#//g' /etc/dnsmasq.conf

```

with password, run:

```
country=GB
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
network={
    ssid="SSID"
    psk="PSK"
}
```

Then reload with:

```
sudo reboot now
```

Testing WiFi: https://raspberrypi.stackexchange.com/questions/61131/is-there-a-way-to-test-a-wifi-password-from-the-command-line-before-connecting-t

```
sudo ifconfig wlan0 down
sudo ifconfig wlan0 up
```

```
wpa_passphrase "ssid" "password" > out
sudo wpa_supplicant -Dwext -cout -iwlan0
```

WITHOUT being connected:

eth0      no wireless extensions.

wlan0     IEEE 802.11  ESSID:off/any  
          Mode:Managed  Access Point: Not-Associated   Tx-Power=31 dBm   
          Retry short limit:7   RTS thr:off   Fragment thr:off
          Power Management:on
          
lo        no wireless extensions.

WITH BEING CONNECTED:

eth0      no wireless extensions.

wlan0     IEEE 802.11  ESSID:"R"  
          Mode:Managed  Frequency:2.437 GHz  Access Point: 88:D7:F6:A7:2A:48   
          Bit Rate=65 Mb/s   Tx-Power=31 dBm   
          Retry short limit:7   RTS thr:off   Fragment thr:off
          Power Management:on
          Link Quality=70/70  Signal level=-35 dBm  
          Rx invalid nwid:0  Rx invalid crypt:0  Rx invalid frag:0
          Tx excessive retries:0  Invalid misc:0   Missed beacon:0

lo        no wireless extensions.

















