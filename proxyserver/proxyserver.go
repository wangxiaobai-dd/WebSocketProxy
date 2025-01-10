package proxyserver

import (
	"ZTWssProxy/network"
	"ZTWssProxy/registry"
	"ZTWssProxy/util"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"ZTWssProxy/configs"
	"github.com/gorilla/mux"
)

type ServerInfo struct {
	ServerID   int    `json:"serverID"`
	ServerIP   string `json:"serverIP"`
	TokenPort  int    `json:"tokenPort"`
	ClientPort int    `json:"clientPort"`
	ConnNum    int    `json:"connNum"`
}

type ProxyServer struct {
	configs.ProxyConfig
	etcdClient *registry.EtcdClient

	tokenManager *TokenManager
	wsServer     *network.WSServer
	httpServer   *network.HttpServer
	connWg       sync.WaitGroup
	gateManager  *network.WSClientManager
}

func NewProxyServer(config configs.ProxyConfig) *ProxyServer {
	return &ProxyServer{
		ProxyConfig:  config,
		etcdClient:   registry.NewEtcdClient(config.EtcdEndPoints),
		tokenManager: &TokenManager{},
		wsServer:     network.NewWSServer(configs.ClientConnAddr, true), //todo
		httpServer:   network.NewHttpServer(),
		gateManager:  network.NewWSClientManager(),
	}
}

func (ps *ProxyServer) registerHandlers() {
	ps.httpServer.AddRoute("/token", ps.handleGameSrvToken, "POST")
	ps.wsServer.AddRoute("/connect/{loginTempID}", ps.handleClientConnect)
}

// 接收游戏服务器token
func (ps *ProxyServer) handleGameSrvToken(w http.ResponseWriter, r *http.Request) {
	t, err := ps.tokenManager.createTokenWithRequest(r)
	if err != nil {
		log.Println(err)
		return
	}

	if err = t.check(); err != nil {
		log.Println(err)
		return
	}

	if _, exist := ps.tokenManager.get(t.LoginTempID); exist {
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

	t, exist := ps.tokenManager.get(loginTempID)
	if !exist {
		log.Println("Failed to find token:", loginTempID)
		return
	}
	//ps.tokenManager.delete(loginTempID)

	log.Println("Connecting to zone:", t.info())

	conn, err := ps.wsServer.UpgradeConnection(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err, t.info())
		return
	}

	ps.connWg.Add(1)
	defer ps.connWg.Done()

	gateAddr := fmt.Sprintf("ws://%s:%d", t.GateIp, t.GatePort)
	wsClient := network.NewWSClient(gateAddr)
	gateConn, err := wsClient.Connect()
	if err != nil {
		log.Println("Failed to connect to GatewayServer:", err, t.info())
		conn.Close()
		return
	}
	ps.gateManager.Add(wsClient)

	clientConn := network.NewWSConn(conn)
	ps.wsServer.AddConn(clientConn)

	log.Println("Connected to gateServer, begin forward:", t.info(), "remote:", conn.RemoteAddr())
	ps.forwardWSMessage(clientConn, gateConn)
}

func (ps *ProxyServer) forwardWSMessage(clientConn, gateConn *network.WSConn) {
	buf := make([]byte, 128*1024)

	var wg sync.WaitGroup
	wg.Add(2)

	// 客户端 -> 网关
	go func() {
		defer gateConn.Close()
		defer wg.Done()
		_, err := io.CopyBuffer(gateConn, clientConn, buf)
		if err != nil {
			log.Println("Error copying from client to gateway:", err)
		} else {
			log.Println("No Copy from client to gateway")
		}
	}()

	// 网关 -> 客户端
	go func() {
		defer clientConn.Close()
		defer wg.Done()
		_, err := io.CopyBuffer(clientConn, gateConn, buf)
		if err != nil {
			log.Println("Error copying from gateway to client:", err)
		} else {
			log.Println("No Copy from gateway to client")
		}
	}()

	wg.Wait()
	ps.gateManager.RemoveByConn(gateConn)
	ps.wsServer.RemoveConn(clientConn)
}

func (ps *ProxyServer) Run() {
	ps.registerHandlers()
	ps.httpServer.Run()
	ps.wsServer.Run()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			ps.tokenManager.cleanExpiredTokens()
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
	ps.gateManager.Destroy()
	ps.connWg.Wait()
	log.Println("Server close")
}

func (ps *ProxyServer) UpdateToEtcd() {
	info := ServerInfo{
		ServerID:   ps.ServerID,
		ServerIP:   ps.ServerIP,
		TokenPort:  ps.TokenPort,
		ClientPort: ps.ClientPort,
		ConnNum:    ps.gateManager.GetConnNum(),
	}
	ps.etcdClient.PutData("", info)
}
