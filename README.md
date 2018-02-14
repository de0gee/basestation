# Find device

```
$ sudo btmgmt find
Discovery started
hci0 type 7 discovering on
hci0 dev_found: 00:0B:57:1B:8C:77 type LE Public rssi -45 flags 0x0000 
AD flags 0x06 
name BlueSense
hci0 type 7 discovering off
```

# Connect to device

```
sudo service bluetooth start
echo "power on" | sudo bluetoothctl
echo "pair 00:0B:57:1B:8C:77" | sudo bluetoothctl
echo "trust 00:0B:57:1B:8C:77" | sudo bluetoothctl
echo "connect 00:0B:57:1B:8C:77" | sudo bluetoothctl
```

# Discover services 

```
$ gattctl --connect 00:0B:57:1B:8C:77
```

# Disconnect

```
echo "disconnect 00:0B:57:1B:8C:77" | sudo bluetoothctl
```

# Pre-requisities

This is the basic install script for Raspberry Pi / Linux
```
sudo apt-get update
sudo apt-get dist-upgrade -y
sudo apt-get install -y vim zsh htop git g++ sqlite3
sudo apt-get install -y pi-bluetooth # pi only
sudo apt-get install -y --no-install-recommends bluetooth

# install node
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs

# install Go
wget https://dl.google.com/go/go1.9.4.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.9.4.*
rm go1.9.*
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
echo 'export GOPATH=$HOME/go' >> ~/.zshrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc

go get -u -v github.com/de0gee/basestation/...
cd /home/pi/go/src/github.com/de0gee/basestation/realtime-client &&  npm install
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
