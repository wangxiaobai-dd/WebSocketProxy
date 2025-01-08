package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"ZTWssProxy/configs"
	"github.com/gorilla/websocket"
)

func main() {
	url := fmt.Sprintf("wss://%s/connect/%d", configs.ClientConnAddr, configs.TestLoginTempID)

	// 创建 HTTP 客户端
	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer conn.Close()

	// 如果需要发送数据，可以使用 conn.WriteMessage
	message := []byte(`{"key":"value"}`)
	err = conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		log.Printf("Failed to send WebSocket message: %v", err)
		return
	}

	// 接收服务器的响应
	_, response, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read WebSocket message: %v", err)
		return
	}

	// 打印响应内容
	fmt.Printf("Response from server: %s\n", string(response))
}
