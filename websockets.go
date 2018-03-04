package main

import (
	"encoding/json"
	"strings"

	log "github.com/cihub/seelog"
	cloud "github.com/de0gee/de0gee-cloud/src"
	"github.com/gorilla/websocket"
)

var conn *websocket.Conn

func setupWebsockets() (err error) {
	conn, _, err = websocket.DefaultDialer.Dial(strings.Replace(CloudServer, "http", "ws", -1)+"/ws2?apikey="+APIKey, nil)
	return
}

func wireData(sensorData cloud.PostSensorData) (err error) {
	log.Debugf("data: %+v", sensorData)
	data, err := json.Marshal(sensorData)
	if err != nil {
		return
	}
	if conn == nil {
		err = setupWebsockets()
		if err != nil {
			return
		}
	}
	errWrite := conn.WriteMessage(websocket.TextMessage, data)
	if errWrite != nil {
		conn.Close()
		conn, _, err = websocket.DefaultDialer.Dial(strings.Replace(CloudServer, "http", "ws", -1)+"/ws2?apikey="+APIKey, nil)
		if err != nil {
			return
		} else {
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}
		}
	}
	return
}

func wireData2(sensorData cloud.PostWebsocket) (err error) {
	log.Debugf("data: %+v", sensorData)
	data, err := json.Marshal(sensorData)
	if err != nil {
		return
	}
	if conn == nil {
		err = setupWebsockets()
		if err != nil {
			return
		}
	}
	errWrite := conn.WriteMessage(websocket.TextMessage, data)
	if errWrite != nil {
		conn.Close()
		conn, _, err = websocket.DefaultDialer.Dial(strings.Replace(CloudServer, "http", "ws", -1)+"/ws2?apikey="+APIKey, nil)
		if err != nil {
			return
		} else {
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				return
			}
		}
	}
	return
}
