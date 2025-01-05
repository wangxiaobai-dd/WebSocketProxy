package proxyserver

import (
	"ZTWssProxy/internal/network"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type connWrapper struct {
}

func (p *connWrapper) Verify() {

}

func (p *connWrapper) Run() {

}

type ProxyServer struct {
	tokenManager *TokenManager
	wsServer     *network.WSServer
	httpServer   *network.HttpServer
}

func NewProxyServer() *ProxyServer {
	return &ProxyServer{
		tokenManager: &TokenManager{},
		wsServer:     &network.WSServer{},
		httpServer:   &network.HttpServer{},
	}
}

func (ps *ProxyServer) RegisterHandlers() {

}

func (ps *ProxyServer) handleGameSrvToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data", err)
		return
	}

}

// url : https://website/zoneID/GatePort/token
func (ps *ProxyServer) handleClientConnect(w http.ResponseWriter, r *http.Request) {
	//conn, err := ps.upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	log.Println("Failed to upgrade connection:", err)
	//	return
	//}
	//defer conn.Close()
}

func (ps *ProxyServer) checkOverdueToken() {

}

func (ps *ProxyServer) Run() {

	ps.wsServer = &network.WSServer{}
	ps.wsServer.Run()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			ps.checkOverdueToken()
		}
	}()

	log.Println("Server run success")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Println("Server closing down, signal", sig)

	ps.close()
}

func (ps *ProxyServer) close() {
	log.Println("Server close")
}
