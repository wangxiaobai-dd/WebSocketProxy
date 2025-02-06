package proxyserver

import (
	"sync"

	"websocket_proxy/network"
)

type ConnContext struct {
	clientConn *network.WSConn
	gateConn   *network.WSConn
	token      *Token
	sync.Once
}

func NewConnContext(clientConn *network.WSConn, gateConn *network.WSConn, token *Token) *ConnContext {
	return &ConnContext{clientConn: clientConn, gateConn: gateConn, token: token}
}

func (c *ConnContext) Close() {
	c.Do(func() {
		c.clientConn.Close()
		c.gateConn.Close()
		//log.Printf("Close conn, loginTempID:%d", c.token.LoginTempID)
	})
}
