package main

import (
	"encoding/json"
	"flag"
	"os"
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
	serverPort      string
)

func main() {
	defer log.Flush()
	var err error
	flag.BoolVar(&doDebug, "debug", false, "enable debugging")
	flag.StringVar(&addressOfDevice, "device", "", "address of BlueSense device")
	flag.StringVar(&Username, "user", "", "username")
	flag.StringVar(&Password, "pass", "", "passphrase")
	flag.StringVar(&CloudServer, "cloud", "https://cloud.de0gee.com", "address of cloud server")
	flag.StringVar(&serverPort, "port", "8005", "port of login server")
	flag.Parse()

	if doDebug {
		SetLogLevel("debug")
	} else {
		SetLogLevel("info")
	}

	log.Debug("starting server")
	go startServer()

	// wait until an API key exists
	for {
		if _, err := os.Stat("apikey"); err == nil {
			// path/to/whatever exists
			break
		}
		time.Sleep(1 * time.Second)
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
