package options

import "fmt"

type ServerOptions struct {
	ServerID   int    `yaml:"ServerID"`
	ServerIP   string `yaml:"ServerIP"`
	TokenPort  int    `yaml:"TokenPort"`
	ClientPort int    `yaml:"ClientPort"`
	BufferSize int    `yaml:"BufferSize"`
}

func (opts ServerOptions) String() string {
	return fmt.Sprintf("ServerID:%d,ServerIP:%s,TokenPort:%d,ClientPort:%d,BufferSize:%dKB", opts.ServerID, opts.ServerIP, opts.TokenPort, opts.ClientPort, opts.BufferSize)
}
