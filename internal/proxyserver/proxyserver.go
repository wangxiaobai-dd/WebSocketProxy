package proxyserver

import (
	"ZTWssProxy/configs"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"ZTWssProxy/internal/network"
	"ZTWssProxy/pkg/util"
	"github.com/gorilla/mux"
)

type ProxyServer struct {
	tokenManager *TokenManager
	wsServer     *network.WSServer
	httpServer   *network.HttpServer

	gateMu    sync.Mutex
	gateConns network.WSConnSet
	wg        sync.WaitGroup
}

func NewProxyServer() *ProxyServer {
	return &ProxyServer{
		tokenManager: &TokenManager{},
		wsServer:     network.NewWSServer(configs.ClientConnAddr, true),
		httpServer:   network.NewHttpServer(),
	}
}

func (ps *ProxyServer) registerHandlers() {
	ps.httpServer.AddRoute("/token", ps.handleGameSrvToken, "POST")
	ps.wsServer.AddRoute("/connect/{loginTempID}", ps.handleClientConnect)
}

// 接收游戏服务器token
func (ps *ProxyServer) handleGameSrvToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		return
	}

	t, err := ps.tokenManager.createTokenWithRequest(r)
	if err != nil {
		log.Println(err)
		return
	}

	if err = t.check(); err != nil {
		log.Println(err)
		return
	}

	if ps.tokenManager.contains(t) {
		log.Println("Token already exists", t.info())
		return
	}
	ps.tokenManager.add(t)
}

// url : https://website/token
func (ps *ProxyServer) handleClientConnect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loginTempID, err := util.StrToUint32(vars["loginTempID"])
	if err != nil {
		log.Println(err)
		return
	}

	t, canLogin := ps.tokenManager.get(loginTempID)
	if canLogin == false {
		log.Println("Failed to find token:", loginTempID)
		return
	}

	log.Println("Connecting to zone:", t.info())

	clientConn, err := ps.wsServer.UpgradeConnection(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	ps.wg.Add(1)
	defer ps.wg.Done()

	gateAddr := fmt.Sprintf("%s:%s", t.GateIp, t.GatePort)
	wsClient := network.NewWSClient(gateAddr)
	gateConn, err := wsClient.Connect()
	if err != nil {
		log.Println("Failed to connect to GatewayServer:", err)
		clientConn.Close()
		return
	}

	wsClientConn := network.NewWSConn(clientConn)
	ps.wsServer.AddConn(wsClientConn)

	wsGateConn := network.NewWSConn(gateConn)
	ps.gateMu.Lock()
	ps.gateConns[wsGateConn] = struct{}{}
	ps.gateMu.Unlock()

	forwardWSMessage(wsClientConn, wsGateConn)

	log.Println("Connected to gateServer, begin forward")
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

	ps.gateMu.Lock()
	for conn := range ps.gateConns {
		conn.Close()
	}
	ps.gateConns = nil
	ps.gateMu.Unlock()

	ps.wg.Wait()
	log.Println("Server close")
}
