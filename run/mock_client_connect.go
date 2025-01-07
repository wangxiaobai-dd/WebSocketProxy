package main

import (
	"ZTWssProxy/configs"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	url := fmt.Sprintf("wss://%s/ws/zone:69000/5240/1000", configs.ClientConnAddr)
	//url := fmt.Sprintf("https://%s/ws/zone:69000/5240/1000", configs.ClientConnAddr)

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
	err = conn.WriteMessage(websocket.TextMessage, message)
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
