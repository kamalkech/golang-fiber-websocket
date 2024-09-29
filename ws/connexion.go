package ws

import (
	"github.com/gofiber/contrib/websocket"
	"sync"
)

var (
	WsConnections      = make(map[*websocket.Conn]struct{})
	WsConnectionsMutex sync.RWMutex
)

type Payload struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
