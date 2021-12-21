package live

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

var (
	clients = make(map[*websocket.Conn]bool)
	mux     sync.Mutex
)

var buffer struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	mux         sync.Mutex
}

var ch chan interface{}

var upGrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 20 * time.Second,
	// 取消 ws 跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	ch = make(chan interface{})
	go BufferWatcher()
}
