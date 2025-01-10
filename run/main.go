package main

import (
	"log"

	"ZTWssProxy/options"
	"github.com/spf13/pflag"

	"ZTWssProxy/proxyserver"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	serverID := pflag.IntP("serverID", "i", 1, "Server ID to select configuration")
	optionFile := pflag.StringP("option", "o", "configs/options.yaml", "Path to the JSON configuration file")
	pflag.Parse()

	ops, err := options.Load(*optionFile)
	if err != nil {
		log.Fatal("Load configuration failed: ", err)
	}
	server := proxyserver.NewProxyServer(*serverID, ops)
	server.Run()
	server.Close()
}
