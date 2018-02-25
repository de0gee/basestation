package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	// "github.com/schollz/patchitup/patchitup"
)

const logLevel = log.DebugLevel
const adapterID = "hci0"

var addressOfDevice = ""

// Configuration is the specific configuration for this de0gee base station
type Configuration struct {
	Username string
}

var config Configuration

func getConfiguration() (c Configuration, err error) {
	if !Exists(path.Join(UserHomeDir(), ".de0gee")) {
		os.MkdirAll(path.Join(UserHomeDir(), ".de0gee"), 0777)
	}
	configFile := path.Join(UserHomeDir(), ".de0gee", "config.toml")
	if !Exists(configFile) {
		// create new configuraiton
		c = Configuration{
			Username: RandomString(4),
		}
		// save the configuration
		buf := new(bytes.Buffer)
		err = toml.NewEncoder(buf).Encode(c)
		if err != nil {
			return
		}
		err = ioutil.WriteFile(configFile, buf.Bytes(), 0755)
		return
	}

	// load configuration
	bConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}
	err = toml.Unmarshal(bConfig, &c)
	return
}

func main() {
	var (
		doDebug    bool
		justServer bool
	)
	flag.BoolVar(&doDebug, "debug", false, "enable debugging")
	flag.BoolVar(&justServer, "serve", false, "enable just the server")
	flag.Parse()

	if doDebug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var err error

	// setup config
	config, err = getConfiguration()
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("running for %s", config.Username)

	if justServer {
		err := startServer()
		if err != nil {
			panic(err)
		}
		return
	}

	// start server
	go func() {
		err := startServer()
		if err != nil {
			panic(err)
		}
	}()

	// periodic dumps
	// go func() {
	// 	for {
	// 		log.Info("dumping the latest")
	// 		db, _ := Open("sensors.db")
	// 		err := db.Dump()
	// 		if err != nil {
	// 			log.Error(err)
	// 		}
	// 		db.Close()

	// 		// patch it up to the server
	// 		patchitup.SetLogLevel("critical")
	// 		err = patchitup.PatchUp("https://data.de0gee.com", config.Username, "sensors.db.sql")
	// 		if err != nil {
	// 			log.Error(err)
	// 		}
	// 		time.Sleep(10 * time.Minute)
	// 	}
	// }()

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
