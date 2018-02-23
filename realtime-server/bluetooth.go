package main

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

// ConnectToBluetooth will connect to the corresponding address
func ConnectToBluetooth(address string) (err error) {
	log.Infof("connecting to %s", address)
	connectedAddress, _ := CurrentConnection()
	if connectedAddress == address {
		log.Infof("connected to %s", address)
		return
	}

	script := `
#!/usr/bin/expect -f

set prompt "#"
set address [lindex $argv 0]

spawn bluetoothctl
expect -re $prompt
send "remove $address\r"
sleep 1
expect -re $prompt
send "scan on\r"
send_user "\nSleeping\r"
sleep 5
send_user "\nDone sleeping\r"
send "scan off\r"
expect "Controller"
send "trust $address\r"
sleep 2
send "pair $address\r"
sleep 2
send "connect $address\r"
sleep 2
send "0000\r"
sleep 3
send_user "\nShould be paired now.\r"
send "quit\r"
expect eof
`
	ioutil.WriteFile("run.sh", []byte(script), 0777)
	log.Debug(script)
	RunCommand(30*time.Second, "expect run.sh "+address)

	connectedAddress, err = CurrentConnection()
	log.Infof("current connected: %s", connectedAddress)
	if err != nil {
		log.Warn(err)
		return
	}
	if connectedAddress != address {
		err = errors.New("problem connecting")
	} else {
		log.Infof("connected to %s", address)
	}
	return
}

func CurrentConnection() (mac string, err error) {
	log.Infof("checking current connections")
	out, _ := RunCommand(1*time.Minute, "hcitool con")
	validMac := regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)
	log.Infof("hcitool con: %s", out)
	for _, line := range strings.Split(out, "\n") {
		macs := validMac.FindAllString(line, 1)
		if len(macs) > 0 {
			mac = macs[0]
			return
		}
	}
	err = errors.New("no devices connected")
	return
}

// DisconnectToBluetooth will disconnect to the corresponding address
func DisconnectToBluetooth(address string) (err error) {
	script := strings.Replace(`
echo "disconnect XX" | bluetoothctl
`, "XX", address, -1)
	ioutil.WriteFile("run.sh", []byte(script), 0777)
	log.Debug(script)
	RunCommand(3*time.Second, "sh run.sh")
	// log.Debug("stdOut:%s", stdOut)
	// log.Debug("stdErr:%s", stdErr)
	log.Infof("disconnected from %s", address)
	return
}

// DiscoverDevice will scan and find the mac address of the device with the given name
func DiscoverDevice(deviceName string) (address string, err error) {
	log.Infof("scanning for %s bluetooth devices", deviceName)
	out, _ := RunCommand(1*time.Minute, "btmgmt find")
	validMac := regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)
	currentMac := ""
	for _, line := range strings.Split(out, "\n") {
		macs := validMac.FindAllString(line, 1)
		if len(macs) > 0 {
			currentMac = macs[0]
		}
		if strings.Contains(line, deviceName) {
			address = currentMac
			return
		}
	}
	err = errors.New("could not find " + deviceName + " device")
	return
}
