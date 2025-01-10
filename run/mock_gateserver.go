package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"ZTWssProxy/configs"
	"ZTWssProxy/network"
	"ZTWssProxy/options"
	"github.com/gorilla/websocket"
)

func main() {
	serverOpts := &options.ServerOptions{
		ServerIP:   configs.TestGateIp,
		ClientPort: configs.TestGatePort,
	}
	sslOpts := &options.SSLOptions{}

	gateServer := network.NewWSServer(serverOpts, sslOpts)
	gateServer.AddRoute("/", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := gateServer.UpgradeConnection(writer, request, nil)
		if err != nil {
			log.Println("request", err)
		}

		go func() {
			defer conn.Close()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Error reading WebSocket message: %v", err)
					break
				}
				fmt.Printf("Received message: %s\n", string(message))

				// 发送响应消息回客户端
				responseMessage := fmt.Sprintf("Received your message: %s", string(message))
				err = conn.WriteMessage(websocket.TextMessage, []byte(responseMessage))
				if err != nil {
					log.Printf("Error writing WebSocket message: %v", err)
					break
				}
			}
			log.Println("WebSocket connection closed")
		}()
	})
	gateServer.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Println("Server closing down, signal", sig)

	gateServer.Close()
}
