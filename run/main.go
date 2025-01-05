package main

import (
	"ZTWssProxy/internal/proxyserver"
)

func main() {
	server := proxyserver.NewProxyServer()
	server.Run()
}
