package proxyserver

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"websocket_proxy/network"
	"websocket_proxy/options"
	"websocket_proxy/registry"
	"websocket_proxy/util"

	"github.com/gorilla/mux"
)

type ProxyServer struct {
	*options.ServerOptions
	*options.TokenOptions
	*options.WSClientOptions
	options.IRegistryOptions
	registry       registry.IRegistry
	tokenManager   *TokenManager
	wsServer       *network.WSServer
	httpServer     *network.HttpServer
	connCtxManager *ConnContextManager
	lastConnNum    int
	connWg         sync.WaitGroup
	taskWg         sync.WaitGroup
}

func NewProxyServer(serverID int, opts *options.Options) *ProxyServer {
	serverOpts := opts.GetServerOptions(serverID)
	registryOpts := opts.GetRegistryOptions()
	return &ProxyServer{
		ServerOptions:    serverOpts,
		TokenOptions:     opts.Token,
		IRegistryOptions: registryOpts,
		WSClientOptions:  opts.WSClient,
		registry:         registry.NewRegistry(registryOpts),
		tokenManager:     &TokenManager{},
		wsServer:         network.NewWSServer(serverOpts),
		httpServer:       network.NewHttpServer(serverOpts),
		connCtxManager:   NewConnContextManager(),
	}
}

func (ps *ProxyServer) registerHandlers() {
	ps.httpServer.AddRoute("/token", ps.handleGameSrvToken, "POST")
	ps.wsServer.AddRoute("/{loginTempID}", ps.handleClientConnect)
}

// 接收游戏服务器token
func (ps *ProxyServer) handleGameSrvToken(w http.ResponseWriter, r *http.Request) {
	token, err := ps.tokenManager.createTokenWithRequest(r, ps.TokenValidTime)
	if err != nil {
		util.RespondWithError(w, "TOKEN CREATE FAIL", err)
		return
	}

	if err = token.check(); err != nil {
		util.RespondWithError(w, "TOKEN INVALID", err)
		return
	}

	if _, exist := ps.tokenManager.get(token.LoginTempID); exist {
		util.RespondWithError(w, "TOKEN EXISTS", err)
		return
	}
	ps.tokenManager.add(token)
	util.RespondWithError(w, "OK", nil)
}

// url : https://website/token
func (ps *ProxyServer) handleClientConnect(w http.ResponseWriter, r *http.Request) {
	ps.connWg.Add(1)
	defer ps.connWg.Done()

	token, err := ps.verifyConnection(w, r)
	if err != nil {
		log.Printf("Failed to verify connection: %s", err)
		return
	}
	log.Printf("Connecting to zone, %s", token)

	clientConn, err := ps.upgradeConnection(w, r)
	if err != nil {
		log.Printf("Failed to upgrade connection: %s", token)
		return
	}

	gateConn, err := ps.connectToGateway(token)
	if err != nil {
		clientConn.Close()
		log.Printf("Failed to connect to gateway: %s,%s", err, token)
		return
	}

	connCtx := ps.manageConnection(clientConn, gateConn, token)

	go ps.forwardWSMessage(connCtx)

	log.Printf("Connected to GateServer, begin forward, %s,client:%s", token, clientConn.RemoteAddr())
}

func (ps *ProxyServer) verifyConnection(w http.ResponseWriter, r *http.Request) (*Token, error) {
	vars := mux.Vars(r)
	loginTempID, err := util.StrToUint32(vars["loginTempID"])
	if err != nil {
		return nil, err
	}
	t, exist := ps.tokenManager.get(loginTempID)
	if !exist {
		return nil, fmt.Errorf("failed to find token: %d", loginTempID)
	}
	ps.tokenManager.delete(loginTempID)
	return t, nil
}

func (ps *ProxyServer) upgradeConnection(w http.ResponseWriter, r *http.Request) (*network.WSConn, error) {
	conn, err := ps.wsServer.UpgradeConnection(w, r, nil)
	if err != nil {
		return nil, err
	}
	return network.NewWSConn(conn, ps.WSClientOptions.MsgType), nil
}

func (ps *ProxyServer) connectToGateway(token *Token) (*network.WSConn, error) {
	gateAddr := fmt.Sprintf("ws://%s:%d", token.GateIp, token.GatePort)
	wsClient := network.NewWSClient(ps.WSClientOptions, gateAddr)
	gateConn, err := wsClient.Connect()
	if err != nil {
		return nil, err
	}
	return gateConn, nil
}

func (ps *ProxyServer) manageConnection(clientConn, gateConn *network.WSConn, token *Token) *ConnContext {
	connCtx := NewConnContext(clientConn, gateConn, token)
	ps.connCtxManager.Add(connCtx)
	return connCtx
}

func (ps *ProxyServer) forwardWSMessage(connCtx *ConnContext) {
	ps.taskWg.Add(1)
	defer ps.taskWg.Done()

	var size int
	if ps.BufferSize == 0 {
		size = 128
	} else {
		size = ps.BufferSize
	}
	buf := make([]byte, size*1024)

	var wg sync.WaitGroup
	wg.Add(2)

	gateConn := connCtx.gateConn
	clientConn := connCtx.clientConn

	// Client -> Gate
	go func() {
		defer gateConn.Close()
		defer wg.Done()
		_, err := io.CopyBuffer(gateConn, clientConn, buf)
		if err != nil && !util.IsClosedNetworkError(err) {
			log.Printf("Error copying from client to gateway: %s", err)
		}
	}()

	// Gate -> Client
	go func() {
		defer clientConn.Close()
		defer wg.Done()
		_, err := io.CopyBuffer(clientConn, gateConn, buf)
		if err != nil && !util.IsClosedNetworkError(err) {
			log.Printf("Error copying from gateway to client: %s", err)
		}
	}()

	wg.Wait()
	connCtx.Close()
	ps.connCtxManager.Remove(connCtx)
}

func (ps *ProxyServer) updateToRegistry(start bool) {
	connNum := ps.connCtxManager.GetConnNum()
	if !start && ps.lastConnNum == connNum {
		return
	}
	ps.lastConnNum = connNum
	info := registry.ServerInfo{
		ServerID:     ps.ServerID,
		ServerIP:     ps.ServerIP,
		ServerDomain: ps.ServerDomain,
		TokenPort:    ps.TokenPort,
		ClientPort:   ps.ClientPort,
		ConnNum:      connNum,
		SecureFlag:   ps.SecureFlag,
	}
	err := ps.registry.PutServer(ps.GetKeyPrefix(), info, ps.GetKeyExpireTime())
	if err != nil {
		log.Printf("Failed to update proxy server info: %s", err)
	}
}

func (ps *ProxyServer) runTicker(ctx context.Context) {
	go func() {
		tickerToken := time.NewTicker(time.Second * time.Duration(ps.CheckTokenDuration))
		defer tickerToken.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerToken.C:
				ps.tokenManager.cleanExpiredTokens()
			}
		}
	}()

	go func() {
		tickerReg := time.NewTicker(time.Second * time.Duration(ps.GetUpdateDuration()))
		defer tickerReg.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tickerReg.C:
				ps.updateToRegistry(false)
			}
		}
	}()
}

func (ps *ProxyServer) wait() {
	ps.connWg.Wait()
	ps.taskWg.Wait()
}

func (ps *ProxyServer) Run() {
	//todo 游戏服务器使用 redis zrange 得到最少连接的服务器 ，找不到服务器信息就zrem

	ctx, cancel := context.WithCancel(context.Background())

	ps.registerHandlers()
	ps.httpServer.Run()
	ps.wsServer.Run()
	ps.updateToRegistry(true)
	ps.runTicker(ctx)

	log.Printf("Server run success, %s,%s,%s", ps.ServerOptions, ps.TokenOptions, ps.IRegistryOptions)

	c := make(chan os.Signal, 1)
	//go func() {
	//	time.Sleep(2 * time.Second)
	//	fmt.Println("Simulating Ctrl+C (SIGINT)...")
	//	c <- syscall.SIGINT
	//}()

	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	sig := <-c
	log.Printf("Server closing down, signal: %s", sig)

	cancel()
	ps.Close()
}

func (ps *ProxyServer) Close() {
	err := ps.registry.DeleteServer(ps.GetKeyPrefix(), ps.ServerID)
	if err != nil {
		log.Println(err)
	}
	ps.registry.Close()
	ps.connCtxManager.Destroy()
	ps.httpServer.Close()
	ps.wsServer.Close()
	ps.wait()
	log.Println("Server close success")
}
