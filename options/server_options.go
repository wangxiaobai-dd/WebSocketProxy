package options

import "fmt"

type ServerOptions struct {
	ServerID     int    `yaml:"ServerID"`
	ServerIP     string `yaml:"ServerIP"`
	ServerDomain string `yaml:"ServerDomain"`
	TokenPort    int    `yaml:"TokenPort"`
	ClientPort   int    `yaml:"ClientPort"`
	BufferSize   int    `yaml:"BufferSize"`
	SecureFlag   bool   `yaml:"SecureFlag"`
	CertFile     string `yaml:"CertFile"`
	KeyFile      string `yaml:"KeyFile"`
}

func (opts ServerOptions) String() string {
	return fmt.Sprintf("ServerID:%d,ServerIP:%s,ServerDomain:%s,TokenPort:%d,ClientPort:%d,"+
		"BufferSize:%dKB,SecureFlag:%v", opts.ServerID, opts.ServerIP, opts.ServerDomain, opts.TokenPort,
		opts.ClientPort, opts.BufferSize, opts.SecureFlag)
}
