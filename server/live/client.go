package live

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

func BufferWatcher() {
	log.Println("BufferWatcher")
	for {
		select {
		case _ = <-ch:
			log.Println("send buffer")
			SendMessage()
		}
	}
}

// ClientWsHandler 处理ws请求
func ClientWsHandler(c *gin.Context) {
	var conn *websocket.Conn
	var err error

	conn, err = upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	// 将上下文保存到字典
	addClient(conn)
}

func SendMessage() {
	setMessage(gin.H{
		"temperature": buffer.Temperature,
		"humidity":    buffer.Humidity,
	})
}

func addClient(conn *websocket.Conn) {
	mux.Lock()
	clients[conn] = true
	mux.Unlock()
}

func getClients() (conns []*websocket.Conn) {
	mux.Lock()

	for k := range clients {
		conns = append(conns, k)
	}

	mux.Unlock()
	return
}

func deleteClient(conn *websocket.Conn) {
	mux.Lock()
	_ = conn.Close()
	delete(clients, conn)
	mux.Unlock()
}

func setMessage(content interface{}) {
	conns := getClients()
	for i := range conns {
		i := i
		err := conns[i].WriteJSON(content)
		if err != nil {
			log.Println("write json err:", err)
			deleteClient(conns[i])
		}
	}
}
