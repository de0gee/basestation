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
	connectedAddress, _ := CurrentConnection()
	if connectedAddress == address {
		return
	}

	script := strings.Replace(`
service bluetooth start
echo "power on" | bluetoothctl
echo "pair XX" | bluetoothctl
echo "trust XX" | bluetoothctl
echo "connect XX" | bluetoothctl
`, "XX", address, -1)
	ioutil.WriteFile("run.sh", []byte(script), 0777)
	log.Debug(script)
	stdOut, stdErr := RunCommand(3*time.Second, "sh run.sh")
	log.Debug("stdOut:%s", stdOut)
	log.Debug("stdErr:%s", stdErr)

	connectedAddress, err = CurrentConnection()
	if err != nil {
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
	log.Debugf("checking current connections")
	out, _ := RunCommand(1*time.Minute, "hcitool con")
	validMac := regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`)
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
	stdOut, stdErr := RunCommand(3*time.Second, "sh run.sh")
	log.Debug("stdOut:%s", stdOut)
	log.Debug("stdErr:%s", stdErr)
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
