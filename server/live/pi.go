package live

import (
	"github.com/0xJacky/zigbee-pi/server/settings"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

func PiWsHandler(c *gin.Context) {
	var conn *websocket.Conn
	var err error

	conn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	var json struct {
		Token string `json:"token"`
	}

	err = conn.WriteJSON(gin.H{
		"message": "ok",
	})

	if err != nil {
		log.Println(err)
		return
	}

	err = conn.ReadJSON(&json)

	if err != nil {
		log.Println("read json err:", err)
		_ = conn.Close()
		return
	}

	if json.Token != settings.AppSettings.TrustedToken {
		log.Println("token auth fail")
		_ = conn.Close()
		return
	}

	for {
		buffer.mux.Lock()
		err = conn.ReadJSON(&buffer)
		if err != nil {
			buffer.mux.Unlock()
			log.Println(err)
			_ = conn.Close()
			return
		}
		log.Printf("T:%s,H:%s\n", buffer.Temperature, buffer.Humidity)
		ch <- 1
		buffer.mux.Unlock()
		err = conn.WriteJSON(gin.H{
			"message": "ok",
		})
		if err != nil {
			log.Println(err)
			_ = conn.Close()
			return
		}
	}
}
