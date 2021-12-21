package main

import (
	"github.com/0xJacky/zigbee-pi/server/settings"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(settings.ServerSettings.RunMode)

	r := InitRoute()

	err := r.Run(":" + settings.ServerSettings.HttpPort)

	if err != nil {
		log.Fatalln(err)
	}
}
