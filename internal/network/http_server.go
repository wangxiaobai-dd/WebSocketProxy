package network

import (
	"log"
	"net/http"

	"ZTWssProxy/configs"
)

type HttpServer struct {
	mux *http.ServeMux
}

func NewHttpServer() *HttpServer {
	return &HttpServer{mux: http.NewServeMux()}
}

func (hs *HttpServer) AddRoute(path string, handlerFunc http.HandlerFunc) {
	hs.mux.HandleFunc(path, handlerFunc)
}

func (hs *HttpServer) Run() {

	httpServer := &http.Server{
		Addr: configs.GameTokenAddr,
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
