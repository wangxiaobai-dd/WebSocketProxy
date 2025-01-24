package network

import (
	"github.com/gorilla/websocket"
	"websocket_proxy/options"
)

type WSClient struct {
	addr    string
	msgType int
	dialer  *websocket.Dialer
}

func NewWSClient(opts *options.WSClientOptions, addr string) *WSClient {
	return &WSClient{msgType: opts.MsgType, addr: addr, dialer: &websocket.Dialer{}}
}

func (c *WSClient) Connect() (*WSConn, error) {
	conn, _, err := c.dialer.Dial(c.addr, nil)
	if err != nil {
		return nil, err
	}
	return NewWSConn(conn, c.msgType), nil
}
