package main

// func getTargetAddress() (targetAddress string) {
// 	log.SetLevel(logLevel)

// 	//clean up connection on exit
// 	defer api.Exit()

// 	log.Debugf("Reset bluetooth device")
// 	a := linux.NewBtMgmt(adapterID)
// 	err := a.Reset()
// 	if err != nil {
// 		log.Error(err)
// 		os.Exit(1)
// 	}

// 	devices, err := api.GetDevices()
// 	if err != nil {
// 		log.Error(err)
// 		os.Exit(1)
// 	}

// 	log.Infof("Cached devices:")
// 	for _, dev := range devices {
// 		name, address := getDeviceInfo(&dev)
// 		log.Println(name, address)
// 		if name == "BlueSense" {
// 			targetAddress = address
// 			return
// 		}
// 	}

// 	log.Infof("Discovered devices:")
// 	err = discoverDevices(adapterID)
// 	if err != nil {
// 		log.Error(err)
// 		os.Exit(1)
// 	}

// 	select {}
// 	return
// }

// func discoverDevices(adapterID string) error {

// 	err := api.StartDiscovery()
// 	if err != nil {
// 		return err
// 	}

// 	log.Debugf("Started discovery")
// 	err = api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {
// 		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
// 		dev := discoveryEvent.Device
// 		log.Println(getDeviceInfo(dev))
// 	}))

// 	return err
// }

// func getDeviceInfo(dev *api.Device) (name, address string) {
// 	if dev == nil {
// 		return
// 	}
// 	props, err := dev.GetProperties()
// 	if err != nil {
// 		log.Errorf("%s: Failed to get properties: %s", dev.Path, err.Error())
// 		return
// 	}
// 	log.Infof("name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)
// 	name = props.Name
// 	address = props.Address
// 	return
// }

// // example of reading temperature from a TI sensortag
// func readValue(tagAddress string) {

// 	dev, err := api.GetDeviceByAddress(tagAddress)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if dev == nil {
// 		panic("Device not found")
// 	}

// 	properties, err := dev.GetProperties()
// 	if err != nil {
// 		panic("problem getting properties")
// 	}
// 	log.Printf("%+v", properties)

// 	err = dev.Connect()
// 	if err != nil {
// 		panic(err)
// 	}

// 	characteristic, err := dev.GetCharByUUID("00002a29-0000-1000-8000-00805f9b34fb")
// 	if err != nil {
// 		panic(err)
// 	}
// 	options := make(map[string]dbus.Variant)
// 	log.Println(characteristic.ReadValue(options))
// }
