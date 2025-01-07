package proxyserver

import (
	"io"
	"log"

	"ZTWssProxy/internal/network"
)

func forwardWSMessage(clientConn, gateConn *network.WSConn) {
	buf := make([]byte, 128*1024)

	// 客户端 -> 网关
	go func() {
		_, err := io.CopyBuffer(gateConn, clientConn, buf)
		if err != nil {
			log.Printf("Error copying from client to gateway: %v", err)
		}
		gateConn.Close()
	}()

	// 网关 -> 客户端
	go func() {
		_, err := io.CopyBuffer(clientConn, gateConn, buf)
		if err != nil {
			log.Printf("Error copying from gateway to client: %v", err)
		}
		clientConn.Close()
	}()
}
