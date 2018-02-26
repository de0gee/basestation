package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/schollz/patchitup-encrypted/patchitup"
)

func startServer() (err error) {
	// setup gin server
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	// Standardize logs
	r.Use(middleWareHandler(), gin.Recovery())

	r.HEAD("/", handlerOK)
	r.GET("/ws", wshandler) // handler for the web sockets (see websockets.go)
	// r.GET("/activity", handlerActivity)
	r.GET("/activity", handlerGetActivity)
	r.POST("/activity", handlerPostActivity)
	r.OPTIONS("/activity", handlerOK)

	// get the username
	r.GET("/username", handlerGetUsername)
	r.OPTIONS("/username", handlerOK)
	log.Infof("Running on 0.0.0.0:%s", "8002")

	err = r.Run(":8002") // listen and serve on 0.0.0.0:8080
	return
}

func handlerOK(c *gin.Context) { // handler for the uptime robot
	c.String(http.StatusOK, "OK")
}

func handlerGetUsername(c *gin.Context) {
	patchitup.DataFolder = "."
	p, _ := patchitup.New(patchitup.Configuration{
		ServerAddress: "https://data.de0gee.com",
		PathToFile:    "sensors.db.sql",
	})
	public, private := p.KeyPair()

	c.JSON(http.StatusOK, gin.H{
		"message": public + "-" + private,
		"success": true,
	})
}

func handlerPostActivity(c *gin.Context) {
	message, err := func(c *gin.Context) (message string, err error) {
		type PostActivity struct {
			Activity string `json:"activity" binding:required`
		}
		var postedJSON PostActivity
		err = c.ShouldBindJSON(&postedJSON)
		if err != nil {
			return
		}

		db, err := Open("sensors.db")
		if err != nil {
			err = errors.Wrap(err, "could not open db")
			return
		}
		defer db.Close()
		id := 0
		for i, activity := range possibleActivities {
			if activity == postedJSON.Activity {
				id = i
				break
			}
		}
		err = db.Add("activity", id)
		if err != nil {
			return
		}
		message = fmt.Sprintf("set activity to '%s'", postedJSON.Activity)
		return
	}(c)
	if err != nil {
		message = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"success": err == nil,
	})
}

func handlerGetActivity(c *gin.Context) {
	message, err := func(c *gin.Context) (message string, err error) {
		db, err := Open("sensors.db")
		if err != nil {
			err = errors.Wrap(err, "could not open db")
			return
		}
		defer db.Close()

		message, err = db.GetLatestActivity()
		return
	}(c)
	if err != nil {
		message = err.Error()
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"success": err == nil,
	})
}

func addCORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
}

func middleWareHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// Add base headers
		addCORS(c)
		// Run next function
		c.Next()
		// Log request
		log.Infof("%v %v %v %s", c.Request.RemoteAddr, c.Request.Method, c.Request.URL, time.Since(t))
	}
}
