package main

import (
	"encoding/binary"
	"math"
	"time"

	log "github.com/cihub/seelog"
	cloud "github.com/de0gee/de0gee-cloud/src"

	"github.com/pkg/errors"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile"
)

type characteristicDefinitionInternal struct {
	info cloud.CharacteristicDefinition
	gatt *profile.GattCharacteristic1
}

func CollectData(address string) (err error) {
	dev, err := api.GetDeviceByAddress(address)
	if err != nil {
		err = errors.Wrap(err, "get device by address")
		return
	}
	if dev == nil {
		err = errors.New("device not found")
		return
	}

	characteristics := make(map[string]characteristicDefinitionInternal)
	for i := range cloud.CharacteristicDefinitions {
		c, err2 := dev.GetCharByUUID(cloud.CharacteristicDefinitions[i].UUID)
		if err2 != nil {
			err = errors.Wrap(err2, "uuid: "+cloud.CharacteristicDefinitions[i].UUID)
			return
		}
		if cloud.CharacteristicDefinitions[i].ValueType == "" {
			continue
		}
		characteristics[cloud.CharacteristicDefinitions[i].UUID] = characteristicDefinitionInternal{
			gatt: c,
			info: cloud.CharacteristicDefinitions[i],
		}
	}

	// read the values forever
	options := make(map[string]dbus.Variant)
	step := float64(0)
	for {
		step++
		for uuid := range characteristics {
			if math.Mod(step, float64(characteristics[uuid].info.SkipSteps)) != 0 {
				continue
			}
			b, err2 := characteristics[uuid].gatt.ReadValue(options)
			if err2 != nil {
				err = errors.Wrap(err2, "problem reading value for "+characteristics[uuid].info.Name)
				return
			}
			log.Debugf("%s data: %+v", characteristics[uuid].info.Name, b)
			if len(b) == 0 {
				continue
			}
			packet := cloud.PostSensorData{
				APIKey:    APIKey,
				Timestamp: time.Now().UTC().UnixNano() / int64(time.Millisecond),
				SensorID:  characteristics[uuid].info.ID,
			}
			switch characteristics[uuid].info.ValueType {
			case "uint8_t":
				packet.SensorValue = int(b[0])
			case "uint16_t":
				packet.SensorValue = int(binary.LittleEndian.Uint16(b))
			case "uint32_t":
				packet.SensorValue = int(binary.LittleEndian.Uint32(b))
			case "special":
				continue
			}
			log.Debugf("%+v", packet)
			err = wireData(packet)
			if err != nil {
				log.Error(err)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
