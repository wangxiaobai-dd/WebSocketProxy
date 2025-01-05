package network

import (
	"ZTWssProxy/configs"
	"log"
	"net/http"
)

type HttpServer struct {
}
type HttpHandler struct {
}

func (handler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (s *HttpServer) Run() {

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
