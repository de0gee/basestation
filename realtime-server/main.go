package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

const logLevel = log.DebugLevel
const adapterID = "hci0"

var addressOfDevice = ""

func main() {
	log.SetLevel(log.DebugLevel)
	var err error

	go func() {
		err := startServer()
		if err != nil {
			panic(err)
		}
	}()

	// periodic dumps
	go func() {
		for {
			log.Info("dumping the latest")
			db, _ := Open("sensors.db")
			db.Dump()
			db.Close()
			time.Sleep(10 * time.Minute)
		}
	}()

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
		os.Exit(1)
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
