package proxyserver

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"ZTWssProxy/internal/network"
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
		wsServer:     network.NewWSServer(),
		httpServer:   network.NewHttpServer(),
	}
}

func (ps *ProxyServer) registerHandlers() {
	ps.httpServer.AddRoute("/token/{tempID}", ps.handleGameSrvToken)
	ps.wsServer.AddRoute("/ws/{zoneID}/{gatePort}/{token}", ps.handleClientConnect)
}

// 接收游戏服务器token
func (ps *ProxyServer) handleGameSrvToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		return
	}

	vars := mux.Vars(r)
	tempID := vars["tempID"]
	log.Println(tempID)

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data", err)
		return
	}

	t := ps.tokenManager.createTokenWithRequest(r)
	if ps.tokenManager.contains(t) {
		log.Println("Token already exists", t.info())
		return
	}
	ps.tokenManager.add(t)
}

// url : https://website/zoneID/GatePort/token
func (ps *ProxyServer) handleClientConnect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	zoneID := vars["zoneID"]
	gatePort := vars["gatePort"]
	token := vars["token"]

	log.Println("Connecting to zone:", zoneID, "gate:", gatePort, "token:", token)

	conn, err := ps.wsServer.UpgradeConnection(w, r, nil)

	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

}

func (ps *ProxyServer) checkOverdueToken() {

}

func (ps *ProxyServer) Run() {
	ps.registerHandlers()
	ps.httpServer.Run()
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

	ps.Close()
}

func (ps *ProxyServer) Close() {
	ps.httpServer.Close()
	ps.wsServer.Close()
	log.Println("Server close")
}
