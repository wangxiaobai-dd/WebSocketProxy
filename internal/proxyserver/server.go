package proxyserver

import (
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	proxyAddr     = "localhost:8080" // 代理服务器监听地址
	recvTokenAddr = "localhost:8081" // 代理服务器接收游戏服务器发送的 token
	certFile      = "server_ssl.crt"
	keyFile       = "server_ssl.key"
)

type ProxyServer struct {
	toLoginTokens sync.Map
}

func NewProxyServer() *ProxyServer {
	return &ProxyServer{}
}

func (p *ProxyServer) handleGameSrvToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data", err)
		return
	}

	t := createTokenWithRequest(r)
	if t.loginTempID == 0 {
		log.Println("Error token", "zoneID:", t.zoneID, "accid", t.accid)
		return
	}

	_, ok := p.toLoginTokens.Load(t.loginTempID)
	if ok {
		log.Println("Token already exists", t.info())
		return
	}

	p.toLoginTokens.Store(t.loginTempID, t)
	log.Println("Save token", t.info())
}

// url : https://website/zoneID/GatePort/token
func (p *ProxyServer) handleClientConnect(w http.ResponseWriter, r *http.Request) {

}

func (p *ProxyServer) checkOverdueToken() {

}

func (p *ProxyServer) Run() {
	http.HandleFunc("/token", p.handleGameSrvToken)
	http.HandleFunc("/wss", p.handleClientConnect)

	var eg errgroup.Group
	if err := eg.Wait(); err != nil {

	}

	go func() {
		err := http.ListenAndServe(recvTokenAddr, nil)
		if err != nil {
			log.Fatal("HTTP server failed to start:", err)
		}
	}()

	go func() {
		err := http.ListenAndServeTLS(proxyAddr, certFile, keyFile, nil)
		if err != nil {
			log.Fatal("WSS server failed to start:", err)
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			p.checkOverdueToken()
		}
	}()
}
