package options

import "fmt"

type RedisOptions struct {
	RegistryOptions `yaml:"Registry"`
	Addr            string `yaml:"Addr"`
	KeepAlive       int    `yaml:"KeepAlive"`
}

func (opts RedisOptions) String() string {
	return fmt.Sprintf("Addr:%s,%s,KeepAlive:%d", opts.Addr, opts.RegistryOptions, opts.KeepAlive)
}
