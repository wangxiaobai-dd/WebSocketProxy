package network

import (
	"github.com/gorilla/websocket"
	"sync"
)

type WSClient struct {
	sync.Once
	addr   string
	dialer *websocket.Dialer
	conn   *websocket.Conn
}

func NewWSClient(addr string) *WSClient {
	return &WSClient{addr: addr, dialer: &websocket.Dialer{}}
}

func (c *WSClient) Connect() (*websocket.Conn, error) {
	conn, _, err := c.dialer.Dial(c.addr, nil)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return conn, nil
}

func (c *WSClient) Close() {
	c.Do(func() {
		c.conn.Close()
	})
}
