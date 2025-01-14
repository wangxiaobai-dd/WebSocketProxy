package network

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"websocket_proxy/options"
)

type HttpServer struct {
	addr   string
	router *mux.Router
}

func NewHttpServer(opts *options.ServerOptions) *HttpServer {
	addr := fmt.Sprintf("%s:%d", opts.ServerIP, opts.TokenPort)
	return &HttpServer{
		addr:   addr,
		router: mux.NewRouter(),
	}
}

func (hs *HttpServer) AddRoute(path string, handlerFunc http.HandlerFunc, method string) {
	hs.router.HandleFunc(path, handlerFunc).Methods(method)
}

func (hs *HttpServer) Run() {
	httpServer := &http.Server{
		Addr:    hs.addr,
		Handler: hs.router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("HTTP server failed to start:", err.Error())
		}
	}()

	log.Println("HTTP server run")
}

func (hs *HttpServer) Close() {
	log.Println("HTTP server close")
}
