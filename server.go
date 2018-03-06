package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func startServer() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Bad": false,
		})
	})
	router.POST("/", func(c *gin.Context) {
		email := c.PostForm("inputEmail")
		passphrase := c.PostForm("inputEmail")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	router.Run(":" + serverPort)
}
