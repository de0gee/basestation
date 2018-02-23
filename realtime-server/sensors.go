package main

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/muka/go-bluetooth/bluez/profile"

	"github.com/pkg/errors"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
)

// Define characteristics
var characteristicDefinitions = map[string]CharacteristicDefinition{
	"00002a6e-0000-1000-8000-00805f9b34fb": {
		Name: "temperature", ValueType: "uint16_t", ID: 0,
		SkipSteps: 100,
	},
	"00002a6f-0000-1000-8000-00805f9b34fb": {
		Name: "humidity", ValueType: "uint8_t", ID: 1,
		SkipSteps: 100,
	},
	"c24229aa-d7e4-4438-a328-c2c548564643": {
		Name: "ambient_light", ValueType: "uint32_t", ID: 2,
		SkipSteps: 2,
	},
	// "61bf1164-529c-4140-9c61-3f5e4fb4c0c1": CharacteristicDefinition{
	// 	Name: "uv_light", ValueType: "uint32_t",
	// },
	"2f256c42-cdef-4378-8e78-694ea0f53ea8": {
		Name: "pressure", ValueType: "uint16_t", ID: 3,
		SkipSteps: 100,
	},
	"15e438b8-558e-4b1f-992f-23f90a8c129b": {
		Name: "motion", ValueType: "uint16_t", ID: 4,
		SkipSteps: 1,
	},
	"00002a19-0000-1000-8000-00805f9b34fb": {
		Name: "battery", ValueType: "uint8_t", ID: 5,
		SkipSteps: 50,
	},
}

type CharacteristicDefinition struct {
	Name           string
	ValueType      string
	ID             int
	SkipSteps      float64
	characteristic *profile.GattCharacteristic1
}

type SensorData struct {
	Name string `json:"name"`
	Data int    `json:"data"`
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

	characteristics := make(map[string]CharacteristicDefinition)
	for uuid := range characteristicDefinitions {
		c, err2 := dev.GetCharByUUID(uuid)
		if err2 != nil {
			err = errors.Wrap(err2, "uuid: "+uuid)
			return
		}
		characteristics[uuid] = CharacteristicDefinition{
			Name:           characteristicDefinitions[uuid].Name,
			ValueType:      characteristicDefinitions[uuid].ValueType,
			ID:             characteristicDefinitions[uuid].ID,
			SkipSteps:      characteristicDefinitions[uuid].SkipSteps,
			characteristic: c,
		}
	}

	// read the values forever
	options := make(map[string]dbus.Variant)
	step := float64(0)
	for {
		step++
		for uuid := range characteristics {
			if math.Mod(step, characteristics[uuid].SkipSteps) != 0 {
				continue
			}
			b, err2 := characteristics[uuid].characteristic.ReadValue(options)
			if err2 != nil {
				err = errors.Wrap(err2, "problem reading value for "+characteristics[uuid].Name)
				return
			}
			data := 0
			log.Debugf("%s data: %+v", characteristics[uuid].Name, b)
			if len(b) == 0 {
				continue
			}
			switch characteristics[uuid].ValueType {
			case "uint8_t":
				data = int(b[0])
			case "uint16_t":
				data = int(binary.LittleEndian.Uint16(b))
			case "uint32_t":
				data = int(binary.LittleEndian.Uint32(b))
			}
			bPayload, err2 := json.Marshal(SensorData{
				Name: characteristics[uuid].Name,
				Data: data,
			})
			if err2 != nil {
				return errors.Wrap(err2, "could not encode "+characteristics[uuid].Name)
			}
			Broadcast(bPayload)
			db, err2 := Open("sensors.db")
			if err2 != nil {
				return errors.Wrap(err2, "could not open db")
			}
			err2 = db.Add("sensor", characteristics[uuid].ID, data)
			db.Close()
			if err2 != nil {
				return errors.Wrap(err2, "could not add sensor")
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

}
