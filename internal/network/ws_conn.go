package network

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSConnSet map[*WSConn]struct{}

type WSConn struct {
	sync.Mutex
	conn      *websocket.Conn
	closeFlag bool
}

func (conn *WSConn) Read(p []byte) (n int, err error) {
	return conn.Read(p)
}

func (conn *WSConn) Write(p []byte) (n int, err error) {
	return conn.Write(p)
}

func (conn *WSConn) Close() {
	conn.Lock()
	defer conn.Unlock()

	if conn.closeFlag {
		return
	}
	conn.closeFlag = true
	conn.conn.Close()
}

func NewWSConn(conn *websocket.Conn) *WSConn {
	wsConn := &WSConn{}
	wsConn.conn = conn
	return wsConn
}
