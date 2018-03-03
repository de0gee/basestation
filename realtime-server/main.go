package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
	cloud "github.com/de0gee/de0gee-cloud/src"
)

const adapterID = "hci0"

var addressOfDevice = ""

var (
	doDebug     bool
	Username    string
	Password    string
	CloudServer string
	APIKey      string
)

func main() {
	defer log.Flush()
	var err error
	flag.BoolVar(&doDebug, "debug", false, "enable debugging")
	flag.StringVar(&Username, "user", "", "username")
	flag.StringVar(&Password, "pass", "", "passphrase")
	flag.StringVar(&CloudServer, "cloud", "http://localhost:8002", "address of cloud server")
	flag.Parse()

	if doDebug {
		SetLogLevel("debug")
	} else {
		SetLogLevel("info")
	}

	// log into the cloud
	payloadBytes, _ := json.Marshal(cloud.LoginJSON{
		Username: Username,
		Password: Password,
	})
	target, err := uploadToServer(payloadBytes, "login")
	if err != nil {
		log.Error(err)
		return
	}
	APIKey = target.Message

	err = setupWebsockets()
	if err != nil {
		log.Error(err)
		return
	}

	addressOfDevice, err = DiscoverDevice("BlueSense")
	if err != nil {
		log.Error(err)
	}
	log.Infof("found BlueSense: %s", addressOfDevice)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		DisconnectToBluetooth(addressOfDevice)
		return
	}()
	for {
		err = connectAndRetrieveData()
		if err != nil {
			log.Warn(err)
		}
		time.Sleep(3 * time.Second)
	}
}

func connectAndRetrieveData() (err error) {
	defer DisconnectToBluetooth(addressOfDevice)
	ConnectToBluetooth(addressOfDevice)
	time.Sleep(3 * time.Second)
	err = ConnectToBluetooth(addressOfDevice)
	if err != nil {
		return
	}
	log.Infof("connected to %s", addressOfDevice)
	time.Sleep(3 * time.Second)

	log.Infof("collecting data")
	err = CollectData(addressOfDevice)
	return
}
