package main

import (
	"ZTWssProxy/configs"
	"ZTWssProxy/internal/network"
	"ZTWssProxy/pkg/util"
	"log"
	"os"
	"os/signal"
)

func main() {
	gateServer := network.NewWSServer(configs.TestGateIp+util.Uint32ToStr(configs.TestGatePort), false)
	gateServer.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Println("Server closing down, signal", sig)

	gateServer.Close()
}
