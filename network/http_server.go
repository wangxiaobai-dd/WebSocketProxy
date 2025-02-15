package network

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"websocket_proxy/options"
	"websocket_proxy/util"
)

type HttpServer struct {
	addr       string
	secureFlag bool
	certFile   string
	keyFile    string
	router     *mux.Router
	server     *http.Server
}

func NewHttpServer(opts *options.ServerOptions) *HttpServer {
	addr := fmt.Sprintf("%s:%d", opts.ServerIP, opts.TokenPort)
	return &HttpServer{
		addr:       addr,
		router:     mux.NewRouter(),
		secureFlag: opts.SecureFlag,
		certFile:   opts.CertFile,
		keyFile:    opts.KeyFile,
	}
}

func (hs *HttpServer) AddRoute(path string, handlerFunc http.HandlerFunc, method string) {
	hs.router.HandleFunc(path, handlerFunc).Methods(method)
}

func (hs *HttpServer) Run() {
	hs.server = &http.Server{
		Addr:    hs.addr,
		Handler: hs.router,
	}

	go func() {
		if hs.secureFlag {
			if err := hs.server.ListenAndServeTLS(hs.certFile, hs.keyFile); !util.IsClosedServerError(err) {
				log.Fatalf("HTTPS server failed to start, %s", err)
			}
		} else {
			if err := hs.server.ListenAndServe(); !util.IsClosedServerError(err) {
				log.Fatalf("HTTP server failed to start, %s", err)
			}
		}
	}()

	log.Println("HTTP server run")
}

func (hs *HttpServer) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := hs.server.Shutdown(ctx)
	if err != nil {
		log.Printf("HTTP server failed to shutdown: %s", err)
	}

	log.Println("HTTP server close")
}
