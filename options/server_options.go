package options

import "fmt"

type ServerOptions struct {
	ServerID   int    `yaml:"ServerID"`
	ServerIP   string `yaml:"ServerIP"`
	TokenPort  int    `yaml:"TokenPort"`
	ClientPort int    `yaml:"ClientPort"`
}

func (opts ServerOptions) String() string {
	return fmt.Sprintf("ServerID:%d,ServerIP:%s,TokenPort:%d,ClientPort:%d", opts.ServerID, opts.ServerIP, opts.TokenPort, opts.ClientPort)
}
