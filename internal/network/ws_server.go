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
	addr       string
	router     *mux.Router
	upgrader   *websocket.Upgrader
	secureFlag bool
	connMu     sync.Mutex
	conns      WSConnSet
}

func NewWSServer(addr string, secureFlag bool) *WSServer {
	return &WSServer{
		addr:   addr,
		router: mux.NewRouter(),
		upgrader: &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
			return true
		}},
		secureFlag: secureFlag,
	}
}

func (s *WSServer) AddRoute(path string, handlerFunc http.HandlerFunc) {
	s.router.HandleFunc(path, handlerFunc)
}

func (s *WSServer) UpgradeConnection(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	conn, err := s.upgrader.Upgrade(w, r, responseHeader)
	return conn, err
}

func (s *WSServer) AddConn(conn *WSConn) {
	s.connMu.Lock()
	s.conns[conn] = struct{}{}
	s.connMu.Unlock()
}

func (s *WSServer) Run() {
	server := &http.Server{
		Addr:         s.addr,           // 监听的地址和端口
		Handler:      s.router,         // 设置请求处理器
		ReadTimeout:  10 * time.Second, // 设置读取超时
		WriteTimeout: 10 * time.Second, // 设置写入超时
		IdleTimeout:  60 * time.Second, // 设置空闲超时
	}

	go func() {
		if s.secureFlag {
			if err := server.ListenAndServeTLS(configs.CertFile, configs.KeyFile); err != nil {
				log.Fatal("Websocket-secure server failed to start:", err.Error())
			}
		} else {
			if err := server.ListenAndServe(); err != nil {
				log.Fatal("Websocket server failed to start:", err.Error())
			}
		}
	}()

	log.Println("Websocket server run, secure:", s.secureFlag)
}

func (s *WSServer) Close() {
	s.connMu.Lock()
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
	s.connMu.Unlock()

	log.Println("Websocket server close")
}
