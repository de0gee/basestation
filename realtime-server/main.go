package main

import (
	log "github.com/Sirupsen/logrus"
)

const logLevel = log.DebugLevel
const adapterID = "hci0"

func main() {
	var err error
	go func() {
		err := startServer()
		if err != nil {
			panic(err)
		}
	}()

	log.SetLevel(log.DebugLevel)
	// address, err := DiscoverDevice("BlueSense")
	// if err != nil {
	// 	log.Error(err)
	// }
	// log.Infof("found BlueSense: %s", address)

	address := "00:0B:57:1B:8C:77"
	err = ConnectToBluetooth(address)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("connected to %s", address)
	log.Infof("collecting data")
	err = CollectData(address)
	if err != nil {
		log.Error(err)
		return
	}
	defer DisconnectToBluetooth(address)
}
