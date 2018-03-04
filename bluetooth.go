package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/cihub/seelog"
	cloud "github.com/de0gee/de0gee-cloud/src"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
)

var (
	device                 = "default"
	addr                   = ""
	sub                    = 0
	sd                     = 10 * time.Second
	definedCharacteristics map[string]charDefinitions
)

type charDefinitions struct {
	info cloud.CharacteristicDefinition
}

func startBluetooth(name string) (err error) {
	// define characteristics
	definedCharacteristics = make(map[string]charDefinitions)
	for i := range cloud.CharacteristicDefinitions {
		if cloud.CharacteristicDefinitions[i].ValueType == "" {
			continue
		}
		uuidName := strings.ToLower(strings.Replace(cloud.CharacteristicDefinitions[i].UUID, "-", "", -1))
		if strings.HasPrefix(cloud.CharacteristicDefinitions[i].UUID, "0000") {
			uuidName = strings.Split(cloud.CharacteristicDefinitions[i].UUID, "-")[0][4:]
		}
		definedCharacteristics[uuidName] = charDefinitions{
			info: cloud.CharacteristicDefinitions[i],
		}
	}

	d, err := dev.NewDevice(device)
	if err != nil {
		log.Errorf("can't new device : %s", err)
		return
	}
	defer d.Stop()

	ble.SetDefaultDevice(d)
	// Default to search device with name of Gopher (or specified by user).
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.LocalName()) == strings.ToUpper(name)
	}

	// If addr is specified, search for addr instead.
	if len(addr) != 0 {
		filter = func(a ble.Advertisement) bool {
			return strings.ToUpper(a.Address().String()) == strings.ToUpper(addr)
		}
	}

	// Scan for specified durantion, or until interrupted by user.
	log.Debug("Scanning for %s...", sd.String())
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), sd))
	cln, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Errorf("can't connect : %s", err)
		return
	}

	// Make sure we had the chance to print out the message.
	done := make(chan struct{})
	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-cln.Disconnected()
		log.Debugf("[ %s ] is disconnected ", cln.Address())
		close(done)
	}()

	log.Debugf("Discovering profile...")
	p, err := cln.DiscoverProfile(true)
	if err != nil {
		log.Errorf("can't discover profile: %s", err)
		return
	}

	// Start the exploration.
	quit := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(quit)
	}()

	t := time.NewTicker(50 * time.Millisecond)
	counter := float64(0)
loop:
	for {
		err = explore(cln, p, counter)
		if err != nil {
			break loop
		}
		select {
		case <-quit:
			break loop
		case <-t.C:
			counter++
		}
	}

	// Disconnect the connection. (On OS X, this might take a while.)
	log.Debugf("Disconnecting [ %s ]... (this might take up to few seconds on OS X)", cln.Address())
	cln.CancelConnection()

	<-done
	return err
}

func explore(cln ble.Client, p *ble.Profile, counter float64) error {
	websocketPacket := cloud.PostWebsocket{
		Timestamp: time.Now().UTC().UnixNano() / int64(time.Millisecond),
		Sensors:   make(map[int]int),
	}

	for _, s := range p.Services {
		// log.Debugf("    Service: %s %s, Handle (0x%02X)", s.UUID, ble.Name(s.UUID), s.Handle)
		for _, c := range s.Characteristics {
			if _, ok := definedCharacteristics[c.UUID.String()]; !ok {
				// log.Debugf("dont have %s", c.UUID)
				continue
			}
			if math.Mod(counter, float64(definedCharacteristics[c.UUID.String()].info.SkipSteps)) != 0 {
				continue
			}
			// log.Debugf("      Characteristic: %s %s",
			// 	c.UUID, ble.Name(c.UUID))
			if (c.Property & ble.CharRead) != 0 {
				b, err := cln.ReadCharacteristic(c)
				if err != nil {
					log.Debugf("Failed to read characteristic: %s", err)
					return err
				}
				log.Debugf("%s: %x %d", definedCharacteristics[c.UUID.String()].info.Name, b, len(b))

				// parse the read data
				if len(b) == 0 {
					continue
				}
				switch definedCharacteristics[c.UUID.String()].info.ValueType {
				case "uint8_t":
					websocketPacket.Sensors[definedCharacteristics[c.UUID.String()].info.ID] = int(b[0])
				case "uint16_t":
					websocketPacket.Sensors[definedCharacteristics[c.UUID.String()].info.ID] = int(binary.LittleEndian.Uint16(b))
				case "uint32_t":
					websocketPacket.Sensors[definedCharacteristics[c.UUID.String()].info.ID] = int(binary.LittleEndian.Uint32(b))
				case "special":
					var val int16
					id := definedCharacteristics[c.UUID.String()].info.ID
					binary.Read(bytes.NewBuffer(b[0:2]), binary.LittleEndian, &val)
					websocketPacket.Sensors[id] = int(val)

					id++
					binary.Read(bytes.NewBuffer(b[2:4]), binary.LittleEndian, &val)
					websocketPacket.Sensors[id] = int(val)

					id++
					binary.Read(bytes.NewBuffer(b[4:6]), binary.LittleEndian, &val)
					websocketPacket.Sensors[id] = int(val)

					continue
				}

			}

		}
	}
	// log.Debugf("%+v", packet)
	if len(websocketPacket.Sensors) > 0 {
		err := wireData2(websocketPacket)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}
