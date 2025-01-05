package network

import (
	"ZTWssProxy/configs"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type WSServer struct {
	onMessage func(conn *websocket.Conn, message []byte)
	onConnect func(conn *websocket.Conn)
	onClose   func(conn *websocket.Conn)
}

func (ws *WSServer) SetOnMessageCallback(cb func(conn *websocket.Conn, message []byte)) {
	ws.onMessage = cb
}

func (ws *WSServer) SetOnConnectCallback(cb func(conn *websocket.Conn)) {
	ws.onConnect = cb
}

func (ws *WSServer) SetOnCloseCallback(cb func(conn *websocket.Conn)) {
	ws.onClose = cb
}

type WSHandler struct {
	upgrader websocket.Upgrader
	connMu   sync.Mutex
	conns    WSConnSet
	wg       sync.WaitGroup
	wrapper  func(conn *WSConn) Wrapper
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("WSHandler, ServeHTTP")
}

func (ws *WSServer) Run() {

	handler := &WSHandler{}

	server := &http.Server{
		Addr:         configs.ClientConnAddr, // 监听的地址和端口
		Handler:      handler,                // 设置请求处理器
		ReadTimeout:  10 * time.Second,       // 设置读取超时
		WriteTimeout: 10 * time.Second,       // 设置写入超时
		IdleTimeout:  60 * time.Second,       // 设置空闲超时
	}

	go func() {
		if err := server.ListenAndServeTLS(configs.CertFile, configs.KeyFile); err != nil {
			log.Fatal("WSS server failed to start:", err.Error())
		}
	}()

	log.Println("WSS server run")
}

func (s *WSServer) Close() {
	log.Println("WSS server close")
}
