package main

import (
	"encoding/json"
	"flag"
	"time"

	log "github.com/cihub/seelog"
	cloud "github.com/de0gee/de0gee-cloud/src"
)

const adapterID = "hci0"

var (
	doDebug         bool
	Username        string
	Password        string
	CloudServer     string
	APIKey          string
	addressOfDevice string
)

func main() {
	defer log.Flush()
	var err error
	flag.BoolVar(&doDebug, "debug", false, "enable debugging")
	flag.StringVar(&addressOfDevice, "device", "", "address of BlueSense device")
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

	for {
		time.Sleep(1 * time.Second)
		err = startBluetooth("BlueSense")
		if err == nil {
			break
		}
		log.Error(err)
	}
}
