package main

import (
	"github.com/0xJacky/zigbee-pi/server/live"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()

		c.Next()
	}
}

func InitRoute() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	r.Use(recovery())

	r.GET("/monitor", live.ClientWsHandler)
	r.GET("/pi", live.PiWsHandler)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "page not found",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	return r
}
