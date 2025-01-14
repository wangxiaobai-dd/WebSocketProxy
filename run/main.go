package main

import (
	"log"

	"github.com/spf13/pflag"
	"websocket_proxy/options"

	"websocket_proxy/proxyserver"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	serverID := pflag.IntP("serverID", "i", 1, "Server ID to select configuration")
	optionFile := pflag.StringP("option", "o", "configs/options.yaml", "Path to the JSON configuration file")
	pflag.Parse()

	opts, err := options.Load(*optionFile)
	if err != nil {
		log.Fatal("Load configuration failed: ", err)
	}
	server := proxyserver.NewProxyServer(*serverID, opts)
	server.Run()
	server.Close()
}
