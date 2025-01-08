package main

import (
	"log"

	"ZTWssProxy/internal/proxyserver"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	server := proxyserver.NewProxyServer()
	server.Run()
}
