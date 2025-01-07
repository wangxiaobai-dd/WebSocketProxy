package network

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"ZTWssProxy/configs"
)

type HttpServer struct {
	router *mux.Router
}

func NewHttpServer() *HttpServer {
	return &HttpServer{router: mux.NewRouter()}
}

func (hs *HttpServer) AddRoute(path string, handlerFunc http.HandlerFunc) {
	hs.router.HandleFunc(path, handlerFunc)
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
