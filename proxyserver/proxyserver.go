package proxyserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"websocket_proxy/network"
	"websocket_proxy/options"
	"websocket_proxy/registry"
	"websocket_proxy/util"

	"github.com/gorilla/mux"
)

type ServerInfo struct {
	ServerID   int    `json:"serverID"`
	ServerIP   string `json:"serverIP"`
	TokenPort  int    `json:"tokenPort"`
	ClientPort int    `json:"clientPort"`
	ConnNum    int    `json:"connNum"`
	SecureFlag bool   `json:"secureFlag"`
}

type ProxyServer struct {
	*options.ServerOptions
	*options.TokenOptions
	*options.SecureOptions
	*options.WSClientOptions
	options.IRegistryOptions
	registry     registry.IRegistry
	tokenManager *TokenManager
	wsServer     *network.WSServer
	httpServer   *network.HttpServer
	connWg       sync.WaitGroup
	gateManager  *network.WSClientManager
}

func NewProxyServer(serverID int, opts *options.Options) *ProxyServer {
	serverOpts := opts.GetServerOptions(serverID)
	registryOpts := opts.GetRegistryOptions()
	return &ProxyServer{
		ServerOptions:    serverOpts,
		TokenOptions:     opts.Token,
		IRegistryOptions: registryOpts,
		SecureOptions:    opts.Secure,
		WSClientOptions:  opts.WSClient,
		registry:         registry.NewRegistry(registryOpts),
		tokenManager:     &TokenManager{},
		wsServer:         network.NewWSServer(serverOpts, opts.Secure),
		httpServer:       network.NewHttpServer(serverOpts),
		gateManager:      network.NewWSClientManager(),
	}
}

func (ps *ProxyServer) registerHandlers() {
	ps.httpServer.AddRoute("/token", ps.handleGameSrvToken, "POST")
	ps.wsServer.AddRoute("/connect/{loginTempID}", ps.handleClientConnect)
}

func (ps *ProxyServer) respond(w http.ResponseWriter, message string, err error) {
	fmt.Fprintf(w, message)
	if err != nil {
		log.Printf("Error: %v", err) // 打印详细的错误日志
	}
}

// 接收游戏服务器token
func (ps *ProxyServer) handleGameSrvToken(w http.ResponseWriter, r *http.Request) {
	t, err := ps.tokenManager.createTokenWithRequest(r, ps.TokenValidTime)
	if err != nil {
		ps.respond(w, "TOKEN CREATE FAIL", err)
		return
	}

	if err = t.check(); err != nil {
		ps.respond(w, "TOKEN INVALID", err)
		return
	}

	if _, exist := ps.tokenManager.get(t.LoginTempID); exist {
		ps.respond(w, "TOKEN EXISTS", err)
		return
	}
	ps.tokenManager.add(t)
	ps.respond(w, "OK", nil)
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
	ps.tokenManager.delete(loginTempID)

	log.Println("Connecting to zone:", t.info())

	conn, err := ps.wsServer.UpgradeConnection(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err, t.info())
		return
	}

	ps.connWg.Add(1)
	defer ps.connWg.Done()

	gateAddr := fmt.Sprintf("ws://%s:%d", t.GateIp, t.GatePort)
	wsClient := network.NewWSClient(ps.WSClientOptions, gateAddr)
	gateConn, err := wsClient.Connect()
	if err != nil {
		log.Println("Failed to connect to GatewayServer:", err, t.info())
		conn.Close()
		return
	}
	ps.gateManager.Add(wsClient)

	clientConn := network.NewWSConn(conn, ps.WSClientOptions.MsgType)
	ps.wsServer.AddConn(clientConn)

	log.Println("Connected to gateServer, begin forward:", t.info(), "remote:", conn.RemoteAddr())
	go ps.forwardWSMessage(clientConn, gateConn)
}

func (ps *ProxyServer) forwardWSMessage(clientConn, gateConn *network.WSConn) {
	var size int
	if ps.BufferSize == 0 {
		size = 128
	} else {
		size = ps.BufferSize
	}
	buf := make([]byte, size*1024)

	var wg sync.WaitGroup
	wg.Add(2)

	// Client -> Gate
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

	// Gate -> Client
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

	tickerToken := time.NewTicker(time.Second * time.Duration(ps.CheckTokenDuration))
	defer tickerToken.Stop()
	go func() {
		for range tickerToken.C {
			ps.tokenManager.cleanExpiredTokens()
		}
	}()

	tickerReg := time.NewTicker(time.Second * time.Duration(ps.GetUpdateDuration()))
	defer tickerReg.Stop()
	go func() {
		for range tickerReg.C {
			ps.updateToRegistry()
		}
	}()

	log.Printf("Server run success, Server:%s,Token:%s,Registry:%s,Secure:%s", ps.ServerOptions, ps.TokenOptions, ps.IRegistryOptions, ps.SecureOptions)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Println("Server closing down, signal", sig)
}

func (ps *ProxyServer) Close() {
	err := ps.registry.DeleteData(ps.GetKey())
	if err != nil {
		log.Println("Server close", err)
	}
	ps.registry.Close()
	ps.httpServer.Close()
	ps.wsServer.Close()
	ps.gateManager.Destroy()
	ps.connWg.Wait()
	log.Println("Server close success")
}

func (ps *ProxyServer) updateToRegistry() {
	info := ServerInfo{
		ServerID:   ps.ServerID,
		ServerIP:   ps.ServerIP,
		TokenPort:  ps.TokenPort,
		ClientPort: ps.ClientPort,
		ConnNum:    ps.gateManager.GetConnNum(),
		SecureFlag: ps.SecureFlag,
	}
	err := ps.registry.PutDataWithTTL(ps.GetKey(), info, ps.GetKeyExpireTime())
	if err != nil {
		log.Println("Failed to update proxy server info:", err)
	}
}

func (ps *ProxyServer) GetKey() string {
	return fmt.Sprintf("%s%d", ps.GetKeyPrefix(), ps.ServerID)
}
