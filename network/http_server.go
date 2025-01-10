package network

import (
	"log"
	"net/http"

	"ZTWssProxy/configs"
	"github.com/gorilla/mux"
)

type HttpServer struct {
	router *mux.Router
}

func NewHttpServer() *HttpServer {
	return &HttpServer{router: mux.NewRouter()}
}

func (hs *HttpServer) AddRoute(path string, handlerFunc http.HandlerFunc, method string) {
	hs.router.HandleFunc(path, handlerFunc).Methods(method)
}

func (hs *HttpServer) Run() {
	httpServer := &http.Server{
		Addr:    configs.GameTokenAddr,
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