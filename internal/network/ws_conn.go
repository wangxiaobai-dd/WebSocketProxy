package network

import (
	"github.com/gorilla/websocket"
	"sync"
)

type WSConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	sync.Mutex
	conn      *websocket.Conn
	writeChan chan []byte
	isClose   bool
}
