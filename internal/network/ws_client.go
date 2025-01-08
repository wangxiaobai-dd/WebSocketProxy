package network

import (
	"github.com/gorilla/websocket"
)

type WSClient struct {
	addr   string
	dialer *websocket.Dialer
	conn   *WSConn // 与服务器连接
}

func NewWSClient(addr string) *WSClient {
	return &WSClient{addr: addr, dialer: &websocket.Dialer{}}
}

func (c *WSClient) Connect() (*WSConn, error) {
	conn, _, err := c.dialer.Dial(c.addr, nil)
	if err != nil {
		return nil, err
	}
	c.conn = NewWSConn(conn)
	return c.conn, nil
}

func (c *WSClient) Close() {
	c.conn.Close()
}
