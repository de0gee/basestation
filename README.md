# Get started

This is the basic install instructions for Raspberry Pi / Linux. Follow these instructions step-by-step to get the Pi working with reading data directly off BlueSense via Bluetooth and visualizing the results in a web browser and saving the data to a database.

In the near-future, these instructions will be compiled into a single script.

## Image

Start from https://github.com/schollz/raspberry-pi-turnkey

```
$ go get -u -v github.com/de0gee/basestation/...
```

```
$ sudo apt-get install -y expect
$ sudo npm install -g yarn
```

Pull latest 

```
$ cd /home/pi/raspberry-pi-turnkey && git pull
$ cd $GOPATH/src/github.com/de0gee/basestation && git pull
$ cd $GOPATH/src/github.com/de0gee/basestation/realtime-client && yarn install
```

(save as de0gee-intermediate.img)

### Add start script

```
$ cp $GOPATH/src/github.com/de0gee/basestation/startup.sh ~/raspberry-pi-turnkey/start.sh
```

Make executable

```
$ chmod +x ~/raspberry-pi-turnkey/startup.sh
```

### Startup server on boot

Open up the `rc.local`

```
$ sudo nano /etc/rc.local
```

And add the following line before `exit 0`:

```
su pi -c '/usr/bin/sudo /usr/bin/python3 /home/pi/raspberry-pi-turnkey/startup.py &'
```


## Update the Raspberry Pi

Run these commands on the Pi (patient, as each takes a few minutes):

```
# install base libraries
sudo apt-get update
sudo apt-get dist-upgrade -y
sudo apt-get install -y vim htop git g++ sqlite3
sudo apt-get install -y pi-bluetooth # pi only
sudo apt-get install -y --no-install-recommends bluetooth
```

## Install node 

Raspberry Pi only:

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

Linux only:

```
# install node
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs
```

## Install Go

Raspberry Pi only:

```
wget https://dl.google.com/go/go1.10.linux-armv6l.tar.gz

```

Linux only:

```
wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz
```

Then, for both Raspberry Pi or Linux:

```
sudo tar -C /usr/local -xzf go*gz
rm go*gz
echo 'export PATH=$PATH:/usr/local/go/bin' >>  ~/.profile
echo 'export GOPATH=$HOME/go' >>  ~/.profile
source ~/.profile
```

## Download the source code for the base station


```
go get -u -v github.com/de0gee/basestation/...
```

## Open up two terminals that are SSHed into the Pi

In the first terminal, do

```
cd /home/pi/go/src/github.com/de0gee/basestation/realtime-client
npm install
npm run start
```

this will compile the client server for viewing the web data.

In the second terminal, do

```
cd /home/pi/go/src/github.com/de0gee/basestation/realtime-server
go build
sudo ./realtime-server
```

The server will automatically scan for Bluesense and pair with it and start pulling data.

## Open up a web browser

Goto your Pi's address `X.X.X.X:3000` to see the data in realtime! All the data is simultaneously being saved to a database.

# Command-line Bluetooth manipulation

## Find device

```
$ sudo btmgmt find
Discovery started
hci0 type 7 discovering on
hci0 dev_found: 00:0B:57:1B:8C:77 type LE Public rssi -45 flags 0x0000 
AD flags 0x06 
name BlueSense
hci0 type 7 discovering off
```

## Connect to device

```
sudo service bluetooth start
echo "power on" | sudo bluetoothctl
echo "pair 00:0B:57:1B:8C:77" | sudo bluetoothctl
echo "trust 00:0B:57:1B:8C:77" | sudo bluetoothctl
echo "connect 00:0B:57:1B:8C:77" | sudo bluetoothctl
```

## Discover services 

```
sudo apt-get install python3-dbus python3-pip
sudo python3 -m pip install gatt
gattctl --connect 00:0B:57:1B:8C:77
```

# Disconnect

```
echo "disconnect 00:0B:57:1B:8C:77" | sudo bluetoothctl
```


# BlueSense GATT Profile

```
Last updated: 2018-FEB-13

Service 1: [Device Information]
sourceId=org.bluetooth.service.device_information
uuid = 180A
    Characteristic 1: [Manufacturer Name String]
    sourceId = org.bluetooth.characteristic.manufacturer_name_string
    uuid = 00002a29-0000-1000-8000-00805f9b34fb
    value type =  utf-8
    length (bytes) = 3

    Characteristic 2: [Model Number String]
    sourceId = org.bluetooth.characteristic.model_number_string
    uuid = 2A24
    value type =  utf-8
    length (bytes) = 9
    
    Characteristic 3: [Hardware Revision String]
    sourceId = org.bluetooth.characteristic.hardware_revision_string
    uuid = 2A27
    value type =  utf-8
    length (bytes) = 3

    Characteristic 4: [Firmware Revision String]
    sourceId = org.bluetooth.characteristic.firmware_revision_string
    uuid = 2A26
    value type =  utf-8
    length (bytes) = 5

Service 2: [Battery Service]
sourceId=org.bluetooth.service.battery_service
uuid = 180F
    Characteristic 1: [Battery Level]
    sourceId = org.bluetooth.characteristic.battery_level
    uuid = 2A19
    value type =  uint8_t
    length (bytes) = 1

Service 3: [Environmental Sensors]
sourceId=custom.type
uuid = c355c42e-b56c-458e-bacb-9248717bbac2
    Characteristic 1: [Temperature]
    sourceId = org.bluetooth.characteristic.temperature
    uuid = 00002a6e-0000-1000-8000-00805f9b34fb
    value type =  int16_t
    length (bytes) = 2

Characteristic 2: [Humidity]
    sourceId = org.bluetooth.characteristic.humidity
    uuid = 2A6F
    value type =  uint8_t
    length (bytes) = 1

    Characteristic 3: [Ambient Light]
    sourceId = custom.type
    uuid = c24229aa-d7e4-4438-a328-c2c548564643
    value type =  uint32_t
    length (bytes) = 4

    Characteristic 4: [UV Light]
    sourceId = custom.type
    uuid = 61bf1164-529c-4140-9c61-3f5e4fb4c0c1
    value type =  uint8_t
    length (bytes) = 1

    Characteristic 5: [Pressure]
    sourceId = custom.type
    uuid = 2f256c42-cdef-4378-8e78-694ea0f53ea8
    value type =  uint16_t
    length (bytes) = 2

Service 4: [IMU]
sourceId=custom.type
uuid = 5b2c25e7-7c43-4a15-a4c6-7cf2d81e1b40
    Characteristic 1: [Motion]
    sourceId = custom.type
    uuid = 15e438b8-558e-4b1f-992f-23f90a8c129b
    value type =  uint16_t
    length (bytes) = 2
```
