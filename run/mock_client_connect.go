package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"

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

	go func() {
		defer conn.Close()
		defer os.Exit(0)
		for {
			_, response, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Failed to read WebSocket message: %v", err)
				break
			}
			// 打印服务器返回的消息
			fmt.Printf("Response from server: %s\n", string(response))
		}
		log.Println("Closing WebSocket connection")
	}()

	// 从控制台读取用户输入并发送
	reader := bufio.NewReader(os.Stdin)
	for {
		// 提示用户输入
		fmt.Print("Enter message to send (or 'exit' to quit): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}

		input = input[:len(input)-1]
		if input == "exit" {
			fmt.Println("Exiting...")
			break
		}

		// 发送消息到 WebSocket 服务器
		err = conn.WriteMessage(websocket.TextMessage, []byte(input))
		if err != nil {
			log.Printf("Failed to send WebSocket message: %v", err)
		}
	}
}
