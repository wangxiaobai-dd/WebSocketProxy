package main

import (
	"ZTWssProxy/configs"
	"ZTWssProxy/proxyserver"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	config := configs.ProxyConfig{}
	server := proxyserver.NewProxyServer(config)
	server.Run()
}
