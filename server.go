package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/cihub/seelog"
	"github.com/de0gee/de0gee-cloud/src"

	"github.com/gin-gonic/gin"
)

func startServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Message": "",
		})
	})
	router.POST("/", func(c *gin.Context) {
		data := cloud.LoginJSON{
			Username: c.PostForm("inputEmail"),
			Password: c.PostForm("inputPassword"),
		}
		payloadBytes, _ := json.Marshal(data)
		target, err := uploadToServer(payloadBytes, "login")
		if target.Success == false {
			err = errors.New(target.Message)
		}
		log.Debugf("%+v", target)
		if err != nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"Message": err.Error(),
			})
		} else {
			log.Debugf("redirecting to %s", CloudServer+"/realtime?apikey="+target.Message)
			ioutil.WriteFile("apikey", []byte(target.Message), 0755)
			c.Redirect(http.StatusMovedPermanently, CloudServer+"/realtime?apikey="+target.Message)
		}
	})
	router.Run(":" + serverPort)
}
