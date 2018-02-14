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
echo "power on" | sudo bluetoothctl
echo "pair 00:0B:57:1B:8C:77" | sudo bluetoothctl
echo "trust 00:0B:57:1B:8C:77" | sudo bluetoothctl
echo "connect 00:0B:57:1B:8C:77" | sudo bluetoothctl
```

# Disconnect

```
echo "disconnect 00:0B:57:1B:8C:77" | sudo bluetoothctl
```

# Read data

Script basis: https://github.com/schollz/gatt-python


```
sudo python3 -m pip install gatt
sudo apt-get install python3-dbus
python3 run.py
```
