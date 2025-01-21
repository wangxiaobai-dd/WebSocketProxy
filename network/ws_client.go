package network

import (
	"github.com/gorilla/websocket"
	"websocket_proxy/options"
)

type WSClient struct {
	addr    string
	dialer  *websocket.Dialer
	conn    *WSConn // 与服务器连接
	msgType int
}

func NewWSClient(opts *options.WSClientOptions, addr string) *WSClient {
	return &WSClient{msgType: opts.MsgType, addr: addr, dialer: &websocket.Dialer{}}
}

func (c *WSClient) Connect() (*WSConn, error) {
	conn, _, err := c.dialer.Dial(c.addr, nil)
	if err != nil {
		return nil, err
	}
	c.conn = NewWSConn(conn, c.msgType)
	return c.conn, nil
}

func (c *WSClient) Close() {
	c.conn.Close()
}
