package network

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	sync.Mutex
	conn      *websocket.Conn
	writeChan chan []byte
	isClose   bool
}
