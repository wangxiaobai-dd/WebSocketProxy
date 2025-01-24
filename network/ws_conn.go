package network

import (
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

type WSConnSet map[*WSConn]struct{}

type WSConn struct {
	msgType int
	conn    *websocket.Conn
	sync.Once
}

func NewWSConn(conn *websocket.Conn, msgType int) *WSConn {
	wsConn := &WSConn{}
	wsConn.conn = conn
	if msgType == 0 {
		wsConn.msgType = websocket.TextMessage
	} else {
		wsConn.msgType = msgType
	}
	return wsConn
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
	err = conn.conn.WriteMessage(conn.msgType, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (conn *WSConn) Close() {
	conn.Do(func() {
		//log.Printf("Close Conn,addr:%s", conn.conn.RemoteAddr())
		conn.conn.Close()
	})
}

func (conn *WSConn) RemoteAddr() net.Addr {
	return conn.conn.RemoteAddr()
}
