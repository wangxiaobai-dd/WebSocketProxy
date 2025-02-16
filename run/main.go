package main

import (
	"log"
	"strconv"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/pflag"
	"websocket_proxy/options"
	"websocket_proxy/proxyserver"
)

var serverID *int
var opts *options.Options

func init() {
	serverID = pflag.IntP("serverID", "i", 1, "Server ID to select configuration")
	optionFile := pflag.StringP("option", "o", "configs/options.yaml", "Path to the JSON configuration file")
	pflag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	opts, err = options.Load(*optionFile)
	if err != nil {
		log.Fatal("Load configuration failed: ", err)
	}

	if opts.Log.Console == false {
		prefix := opts.Log.Path + opts.Log.LinkName + strconv.Itoa(*serverID)
		writer, _ := rotatelogs.New(
			prefix+".log.%Y%m%d-%H",
			rotatelogs.WithLinkName(prefix+".log"),
			rotatelogs.WithRotationTime(time.Hour),
		)
		log.SetOutput(writer)
	}
}

func main() {
	server := proxyserver.NewProxyServer(*serverID, opts)
	server.Run()
}
