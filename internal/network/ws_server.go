package network

import (
	"log"
	"net/http"
	"sync"
	"time"

	"ZTWssProxy/configs"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	router     *mux.Router
	upgrader   *websocket.Upgrader
	secureFlag bool
	connMu     sync.Mutex
	conns      WSConnSet
}

func NewWSServer(secureFlag bool) *WSServer {
	return &WSServer{
		router: mux.NewRouter(),
		upgrader: &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
			return true
		}},
		secureFlag: secureFlag,
	}
}

func (ws *WSServer) AddRoute(path string, handlerFunc http.HandlerFunc) {
	ws.router.HandleFunc(path, handlerFunc)
}

func (ws *WSServer) UpgradeConnection(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	conn, err := ws.upgrader.Upgrade(w, r, responseHeader)
	return conn, err
}

func (ws *WSServer) AddConn(conn *WSConn) {
	ws.connMu.Lock()
	ws.conns[conn] = struct{}{}
	ws.connMu.Unlock()
}

func (ws *WSServer) Run() {

	server := &http.Server{
		Addr:         configs.ClientConnAddr, // 监听的地址和端口
		Handler:      ws.router,              // 设置请求处理器
		ReadTimeout:  10 * time.Second,       // 设置读取超时
		WriteTimeout: 10 * time.Second,       // 设置写入超时
		IdleTimeout:  60 * time.Second,       // 设置空闲超时
	}

	go func() {
		if ws.secureFlag {
			if err := server.ListenAndServeTLS(configs.CertFile, configs.KeyFile); err != nil {
				log.Fatal("Websocket-secure server failed to start:", err.Error())
			}
		} else {
			if err := server.ListenAndServe(); err != nil {
				log.Fatal("Websocket server failed to start:", err.Error())
			}
		}
	}()

	log.Println("Websocket server run, secure:", ws.secureFlag)
}

func (ws *WSServer) Close() {
	ws.connMu.Lock()
	for conn := range ws.conns {
		conn.Close()
	}
	ws.conns = nil
	ws.connMu.Unlock()

	log.Println("Websocket server close")
}
