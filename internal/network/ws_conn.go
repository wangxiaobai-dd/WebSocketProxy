package network

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSConnSet map[*WSConn]struct{}

type WSConn struct {
	sync.Once
	conn *websocket.Conn
}

func (conn *WSConn) Read(p []byte) (n int, err error) {
	_, message, err := conn.conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	n = copy(p, message)
	return n, nil
}

func (conn *WSConn) Write(p []byte) (n int, err error) {
	err = conn.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (conn *WSConn) Close() {
	conn.Do(func() {
		conn.conn.Close()
	})
}

func NewWSConn(conn *websocket.Conn) *WSConn {
	wsConn := &WSConn{}
	wsConn.conn = conn
	return wsConn
}
