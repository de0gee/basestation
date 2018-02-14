package main

import (
	"encoding/binary"
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
)

type CharacteristicDefinition struct {
	Name      string
	ValueType string
}

var CharactersticDefinitions = map[string]CharacteristicDefinition{
	"15e438b8-558e-4b1f-992f-23f90a8c129b": CharacteristicDefinition{
		Name: "motion", ValueType: "uint16_t",
	},
}

type SensorData struct {
	Name string `json:"name"`
	Data int    `json:"data"`
}

func CollectData(address string) {
	dev, err := api.GetDeviceByAddress(address)
	if err != nil {
		panic(err)
	}
	if dev == nil {
		panic("Device not found")
	}

	characteristic, err := dev.GetCharByUUID("15e438b8-558e-4b1f-992f-23f90a8c129b")
	if err != nil {
		panic(err)
	}
	light, err := dev.GetCharByUUID("c24229aa-d7e4-4438-a328-c2c548564643")
	options := make(map[string]dbus.Variant)
	for {
		b, err := (characteristic.ReadValue(options))
		if err != nil {
			log.Warn(err)
			continue
		}
		u16 := binary.LittleEndian.Uint16(b)
		bPayload, err := json.Marshal(SensorData{
			Name: "motion",
			Data: int(u16),
		})
		if err != nil {
			log.Warn(err)
			continue
		}
		Broadcast(bPayload)

		b, err = (light.ReadValue(options))
		if err != nil {
			log.Warn(err)
			continue
		}
		u32 := binary.LittleEndian.Uint32(b)
		bPayload, err = json.Marshal(SensorData{
			Name: "ambient_light",
			Data: int(u32),
		})
		if err != nil {
			log.Warn(err)
			continue
		}
		Broadcast(bPayload)
		time.Sleep(10 * time.Millisecond)
	}

}
