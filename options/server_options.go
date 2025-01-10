package options

import "fmt"

type ServerOptions struct {
	ServerID   int    `json:"ServerID" yaml:"ServerID"`
	ServerIP   string `json:"ServerIP" yaml:"ServerIP"`
	TokenPort  int    `json:"TokenPort" yaml:"TokenPort"`
	ClientPort int    `json:"ClientPort" yaml:"ClientPort"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{}
}

func (opts *ServerOptions) String() string {
	return fmt.Sprintf("ServerID:%d, ServerIP:%s, TokenPort:%d, ClientPort:%d", opts.ServerID, opts.ServerIP, opts.TokenPort, opts.ClientPort)
}
