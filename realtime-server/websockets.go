package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/satori/go.uuid"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Websockets struct {
	connections map[string]*websocket.Conn
	sync.Mutex
}

var (
	ws Websockets
)

func init() {
	ws.Lock()
	defer ws.Unlock()
	ws.connections = make(map[string]*websocket.Conn)
}

func wshandler(c *gin.Context) {
	var w http.ResponseWriter = c.Writer
	var r *http.Request = c.Request

	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	ws.Lock()
	id := fmt.Sprintf("%s", uuid.NewV4())
	log.Infof("add connection %s", id)
	ws.connections[id] = conn
	ws.Unlock()
}

func Broadcast(msg []byte) (err error) {
	ws.Lock()
	defer ws.Unlock()
	removeConns := make(map[string]struct{})
	for conn := range ws.connections {
		err = ws.connections[conn].WriteMessage(1, msg)
		if err != nil {
			removeConns[conn] = struct{}{}
		}
	}
	for conn := range removeConns {
		delete(ws.connections, conn)
	}
	return
}
